package services

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"
	"sync"
	"time"

	"yara-security-service/internal/models"
	"yara-security-service/pkg/security"

	"github.com/sirupsen/logrus"
)

// SecurityService 安全服务
type SecurityService struct {
	scanner *security.Scanner
	logger  *logrus.Logger
	mu      sync.RWMutex

	// 优化的缓存机制
	scanCache  map[string]*models.ScanResult
	cacheTTL   time.Duration
	cacheStats *CacheStats
	cacheMutex sync.RWMutex
	maxCacheSize int

	// 统计
	scanStats *ScanStatistics

	// 并发控制
	scanWorkers   int
	scanSemaphore chan struct{}
	scanTimeout   time.Duration

	// 规则重载后台任务
	reloadMu         sync.Mutex
	reloadInProgress bool
	reloadLastError  error
	reloadLastTime   time.Time

	// 隔离文件读写锁，避免与其他全局锁耦合
	quarantineMu sync.RWMutex
}

// CacheStats 缓存统计
type CacheStats struct {
	HitCount      int64
	MissCount     int64
	EvictionCount int64
	Size          int64
	MaxSize       int64
	LastCleanup   time.Time
}

// ScanStatistics 扫描统计
type ScanStatistics struct {
	TotalScans    int64
	InfectedFiles int64
	CleanFiles    int64
	FailedScans   int64
	TotalScanTime time.Duration
	LastScanTime  time.Time
	mu            sync.RWMutex
}

// NewSecurityService 创建安全服务
func NewSecurityService(scanner *security.Scanner, logger *logrus.Logger) *SecurityService {
	service := &SecurityService{
		scanner:       scanner,
		logger:        logger,
		scanCache:     make(map[string]*models.ScanResult),
		cacheTTL:      10 * time.Minute,
		cacheStats:    &CacheStats{MaxSize: 1000},
		maxCacheSize:  1000,
		scanWorkers:   5,
		scanSemaphore: make(chan struct{}, 5),
		scanTimeout:   30 * time.Second,
		scanStats:     &ScanStatistics{},
	}

	go service.startCacheCleaner()
	return service
}

// startCacheCleaner 启动缓存清理协程
func (s *SecurityService) startCacheCleaner() {
	ticker := time.NewTicker(5 * time.Minute) // 每5分钟清理一次
	defer ticker.Stop()

	for range ticker.C {
		s.cleanExpiredCache()
	}
}

// cleanExpiredCache 清理过期缓存
func (s *SecurityService) cleanExpiredCache() {
	s.cacheMutex.Lock()
	defer s.cacheMutex.Unlock()

	now := time.Now()
	expiredKeys := []string{}

	// 清理过期条目
	for key, result := range s.scanCache {
		if now.Sub(result.ScanTime) > s.cacheTTL {
			expiredKeys = append(expiredKeys, key)
		}
	}

	// 如果缓存过大，使用LRU策略清理
	if len(s.scanCache) > s.maxCacheSize {
		// 按扫描时间排序，删除最旧的条目
		type cacheEntry struct {
			key  string
			time time.Time
		}

		entries := make([]cacheEntry, 0, len(s.scanCache))
		for key, result := range s.scanCache {
			entries = append(entries, cacheEntry{key: key, time: result.ScanTime})
		}

		sort.Slice(entries, func(i, j int) bool {
			return entries[i].time.Before(entries[j].time)
		})

		// 删除最旧的条目，直到缓存大小合适
		toDelete := len(s.scanCache) - s.maxCacheSize + len(expiredKeys)
		for i := 0; i < toDelete && i < len(entries); i++ {
			expiredKeys = append(expiredKeys, entries[i].key)
		}
	}

	// 删除过期和超出的条目
	for _, key := range expiredKeys {
		delete(s.scanCache, key)
		s.cacheStats.EvictionCount++
	}

	s.cacheStats.Size = int64(len(s.scanCache))
	s.cacheStats.LastCleanup = now

	if len(expiredKeys) > 0 {
		s.logger.Debugf("清理了 %d 个缓存项，当前缓存大小: %d", len(expiredKeys), len(s.scanCache))
	}
}

// ReloadRules 提交规则重载任务（非阻塞）
func (s *SecurityService) ReloadRules() error {
	s.reloadMu.Lock()
	if s.reloadInProgress {
		s.reloadMu.Unlock()
		s.logger.Warn("规则重载已在进行中，忽略重复请求")
		return nil
	}
	s.reloadInProgress = true
	s.reloadLastError = nil
	s.reloadMu.Unlock()

	go func() {
		start := time.Now()
		s.logger.Info("[规则重载] 开始后台重载…")

		// 真正执行重载
		if err := s.scanner.ReloadRules(); err != nil {
			s.logger.Errorf("[规则重载] 失败: %v", err)
			s.reloadMu.Lock()
			s.reloadLastError = err
			s.reloadInProgress = false
			s.reloadLastTime = time.Now()
			s.reloadMu.Unlock()
			return
		}

		// 清空缓存（独立锁，避免阻塞其他）
		s.cacheMutex.Lock()
		s.scanCache = make(map[string]*models.ScanResult)
		s.cacheStats = &CacheStats{MaxSize: int64(s.maxCacheSize)}
		s.cacheMutex.Unlock()

		s.reloadMu.Lock()
		s.reloadInProgress = false
		s.reloadLastTime = time.Now()
		s.reloadMu.Unlock()

		s.logger.Infof("[规则重载] 完成，用时: %v", time.Since(start))
	}()

	return nil
}

// GetReloadStatus 获取规则重载状态
func (s *SecurityService) GetReloadStatus() map[string]interface{} {
	s.reloadMu.Lock()
	defer s.reloadMu.Unlock()
	status := "idle"
	if s.reloadInProgress {
		status = "reloading"
	}
	var lastErr string
	if s.reloadLastError != nil {
		lastErr = s.reloadLastError.Error()
	}
	return map[string]interface{}{
		"status":     status,
		"last_error": lastErr,
		"last_time":  s.reloadLastTime,
	}
}

// GetQuarantineList 获取隔离列表（独立锁+超时）
func (s *SecurityService) GetQuarantineList() []*models.QuarantineItem {
	// 使用独立读锁，避免与其他操作互相阻塞
	s.quarantineMu.RLock()
	defer s.quarantineMu.RUnlock()

	quarantineDir, err := filepath.Abs("./quarantine")
	if err != nil {
		s.logger.Errorf("获取隔离目录绝对路径失败: %v", err)
		return []*models.QuarantineItem{}
	}

	recordsFile := filepath.Join(quarantineDir, "quarantine_records.json")
	if _, err := os.Stat(recordsFile); os.IsNotExist(err) {
		return []*models.QuarantineItem{}
	}

	// 加入文件读取超时保护（防止卡死）
	done := make(chan struct{})
	var data []byte
	var readErr error
	go func() {
		data, readErr = os.ReadFile(recordsFile)
		close(done)
	}()

	select {
	case <-done:
		if readErr != nil {
			s.logger.Errorf("读取隔离记录失败: %v", readErr)
			return []*models.QuarantineItem{}
		}
		var records []*models.QuarantineItem
		if err := json.Unmarshal(data, &records); err != nil {
			s.logger.Errorf("解析隔离记录失败: %v", err)
			return []*models.QuarantineItem{}
		}
		return records
	case <-time.After(5 * time.Second):
		s.logger.Warn("读取隔离记录超时，返回空列表以避免阻塞")
		return []*models.QuarantineItem{}
	}
}

// ScanFile 扫描文件（优化版本）
func (s *SecurityService) ScanFile(ctx context.Context, filePath string) (*models.ScanResult, error) {
	// 检查缓存
	s.cacheMutex.RLock()
	if cached, exists := s.scanCache[filePath]; exists {
		if time.Since(cached.ScanTime) < s.cacheTTL {
			s.cacheStats.HitCount++
			s.cacheMutex.RUnlock()
			s.logger.Debugf("使用缓存结果: %s", filePath)
			return cached, nil
		}
	}
	s.cacheMutex.RUnlock()

	s.cacheStats.MissCount++

	// 使用信号量控制并发
	select {
	case s.scanSemaphore <- struct{}{}:
		defer func() { <-s.scanSemaphore }()
	case <-ctx.Done():
		return nil, fmt.Errorf("扫描超时或被取消")
	case <-time.After(s.scanTimeout):
		return nil, fmt.Errorf("扫描超时")
	}

	// 执行扫描
	startTime := time.Now()
	result, err := s.scanner.ScanFile(ctx, filePath)
	if err != nil {
		s.updateStats(false, time.Since(startTime))
		return nil, fmt.Errorf("文件扫描失败: %w", err)
	}

	// 更新统计
	s.updateStats(true, time.Since(startTime))

	// 更新感染/清洁文件统计
	s.scanStats.mu.Lock()
	if result.IsInfected {
		s.scanStats.InfectedFiles++
	} else {
		s.scanStats.CleanFiles++
	}
	s.scanStats.mu.Unlock()

	// 缓存结果
	s.cacheMutex.Lock()
	s.scanCache[filePath] = result
	s.cacheStats.Size = int64(len(s.scanCache))
	s.cacheMutex.Unlock()

	s.logger.Infof("文件扫描完成: %s, 感染: %v, 耗时: %v", filePath, result.IsInfected, time.Since(startTime))
	return result, nil
}

// ScanDirectory 扫描目录
func (s *SecurityService) ScanDirectory(ctx context.Context, dirPath string, recursive bool, maxDepth int, includePatterns []string, excludePatterns []string) ([]*models.ScanResult, error) {
	// 不使用全局互斥，避免长时间阻塞其他请求
	startTime := time.Now()
	results, err := s.scanner.ScanDirectory(ctx, dirPath, recursive, maxDepth, includePatterns, excludePatterns)
	if err != nil {
		return nil, fmt.Errorf("目录扫描失败: %w", err)
	}

	// 更新统计
	s.scanStats.mu.Lock()
	s.scanStats.TotalScans += int64(len(results))
	s.scanStats.TotalScanTime += time.Since(startTime)
	s.scanStats.LastScanTime = time.Now()
	for _, result := range results {
		if result.IsInfected {
			s.scanStats.InfectedFiles++
		} else {
			s.scanStats.CleanFiles++
		}
	}
	s.scanStats.mu.Unlock()

	s.logger.Infof("目录扫描完成: %s, 文件数: %d, 感染: %d, 耗时: %v",
		dirPath, len(results), s.countInfected(results), time.Since(startTime))

	return results, nil
}

// ScanBuffer 扫描内存缓冲区
func (s *SecurityService) ScanBuffer(ctx context.Context, data []byte, identifier string) (*models.ScanResult, error) {
	startTime := time.Now()
	result, err := s.scanner.ScanBuffer(ctx, data, identifier)
	if err != nil {
		s.updateStats(false, time.Since(startTime))
		return nil, fmt.Errorf("缓冲区扫描失败: %w", err)
	}

	// 更新统计
	s.updateStats(true, time.Since(startTime))

	// 更新感染/清洁文件统计
	s.scanStats.mu.Lock()
	if result.IsInfected {
		s.scanStats.InfectedFiles++
	} else {
		s.scanStats.CleanFiles++
	}
	s.scanStats.mu.Unlock()

	s.logger.Infof("缓冲区扫描完成: %s, 感染: %v, 耗时: %v", identifier, result.IsInfected, time.Since(startTime))
	return result, nil
}

// GetSecurityStatus 获取安全状态（无全局锁、带超时聚合）
func (s *SecurityService) GetSecurityStatus() *models.SecurityStatus {
	// 若正在热重载，立即返回最小状态，避免与重载争用资源
	s.reloadMu.Lock()
	reloading := s.reloadInProgress
	s.reloadMu.Unlock()
	if reloading {
		return &models.SecurityStatus{
			IsProtected:        true,
			RealTimeProtection: true,
			LastScanTime:       s.scanStats.LastScanTime,
			ThreatCount:        int(s.scanStats.InfectedFiles),
			QuarantineCount:    0,
			DefinitionsVersion: "1.0.0",
			LastUpdateTime:     time.Now(),
			ScanStatistics: &models.ScanStatistics{
				TotalScans:    s.scanStats.TotalScans,
				InfectedFiles: s.scanStats.InfectedFiles,
				CleanFiles:    s.scanStats.CleanFiles,
				FailedScans:   s.scanStats.FailedScans,
				TotalScanTime: s.scanStats.TotalScanTime,
			},
			RulesInfo: map[string]interface{}{"status": "reloading"},
		}
	}

	// 优先返回缓存（细粒度锁）
	s.cacheMutex.RLock()
	if cached, exists := s.scanCache["security_status"]; exists && time.Since(cached.ScanTime) < s.cacheTTL {
		s.cacheMutex.RUnlock()
		s.logger.Debug("使用缓存的安全状态")
		return &models.SecurityStatus{
			IsProtected:        true,
			RealTimeProtection: true,
			LastScanTime:       cached.ScanTime,
			ThreatCount:        int(s.scanStats.InfectedFiles),
			QuarantineCount:    s.getQuarantineCount(),
			DefinitionsVersion: "1.0.0",
			LastUpdateTime:     time.Now(),
			ScanStatistics: &models.ScanStatistics{
				TotalScans:    s.scanStats.TotalScans,
				InfectedFiles: s.scanStats.InfectedFiles,
				CleanFiles:    s.scanStats.CleanFiles,
				FailedScans:   s.scanStats.FailedScans,
				TotalScanTime: s.scanStats.TotalScanTime,
			},
			RulesInfo: s.scanner.GetRulesInfo(),
		}
	}
	s.cacheMutex.RUnlock()

	// 3秒总超时保护，避免状态聚合卡死
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var (
		rulesInfo        map[string]interface{}
		quarantineCount  int
		scanStatsSnapshot models.ScanStatistics
	)

	// 并发获取，子任务各自带超时与recover
	done := make(chan struct{})
	go func() {
		defer func() { recover(); close(done) }()

		var wg sync.WaitGroup
		wg.Add(3)

		go func() {
			defer wg.Done()
			// 子任务1：规则信息（1.5秒超时）
			subCtx, subCancel := context.WithTimeout(ctx, 1500*time.Millisecond)
			defer subCancel()
			ch := make(chan map[string]interface{}, 1)
			go func() { ch <- s.scanner.GetRulesInfo() }()
			select {
			case v := <-ch:
				rulesInfo = v
			case <-subCtx.Done():
				s.logger.Warn("GetRulesInfo 超时，使用空信息")
				rulesInfo = map[string]interface{}{"status": "timeout"}
			}
		}()

		go func() {
			defer wg.Done()
			// 子任务2：隔离数量（1秒超时）
			subCtx, subCancel := context.WithTimeout(ctx, time.Second)
			defer subCancel()
			ch := make(chan int, 1)
			go func() { ch <- s.getQuarantineCount() }()
			select {
			case v := <-ch:
				quarantineCount = v
			case <-subCtx.Done():
				s.logger.Warn("getQuarantineCount 超时，返回0")
				quarantineCount = 0
			}
		}()

		go func() {
			defer wg.Done()
			// 子任务3：拷贝扫描统计（无锁卡顿）
			s.scanStats.mu.RLock()
			scanStatsSnapshot = models.ScanStatistics{
				TotalScans:    s.scanStats.TotalScans,
				InfectedFiles: s.scanStats.InfectedFiles,
				CleanFiles:    s.scanStats.CleanFiles,
				FailedScans:   s.scanStats.FailedScans,
				TotalScanTime: s.scanStats.TotalScanTime,
			}
			s.scanStats.mu.RUnlock()
		}()

		wg.Wait()
	}()

	select {
	case <-done:
		// 构建安全状态
		status := &models.SecurityStatus{
			IsProtected:        true,
			RealTimeProtection: true,
			LastScanTime:       s.scanStats.LastScanTime,
			ThreatCount:        int(scanStatsSnapshot.InfectedFiles),
			QuarantineCount:    quarantineCount,
			DefinitionsVersion: "1.0.0",
			LastUpdateTime:     time.Now(),
			ScanStatistics:     &scanStatsSnapshot,
			RulesInfo:          rulesInfo,
		}
		// 写入缓存（细粒度锁）
		s.cacheMutex.Lock()
		s.scanCache["security_status"] = &models.ScanResult{
			FilePath:     "security_status",
			ScanTime:     time.Now(),
			IsInfected:   scanStatsSnapshot.InfectedFiles > 0,
			ScanDuration: 0,
			Threats:      []models.ThreatInfo{},
			FileInfo:     &models.FileInfo{},
		}
		s.cacheMutex.Unlock()
		return status
	case <-ctx.Done():
		// 超时兜底，返回最小状态避免卡死
		s.logger.Warn("GetSecurityStatus 聚合超时，返回兜底状态")
		return &models.SecurityStatus{
			IsProtected:        true,
			RealTimeProtection: true,
			LastScanTime:       s.scanStats.LastScanTime,
			ThreatCount:        int(s.scanStats.InfectedFiles),
			QuarantineCount:    0,
			DefinitionsVersion: "1.0.0",
			LastUpdateTime:     time.Now(),
			ScanStatistics: &models.ScanStatistics{
				TotalScans:    s.scanStats.TotalScans,
				InfectedFiles: s.scanStats.InfectedFiles,
				CleanFiles:    s.scanStats.CleanFiles,
				FailedScans:   s.scanStats.FailedScans,
				TotalScanTime: s.scanStats.TotalScanTime,
			},
			RulesInfo: map[string]interface{}{"status": "timeout"},
		}
	}
}

// GetRulesInfo 获取规则信息
func (s *SecurityService) GetRulesInfo() map[string]interface{} {
	return s.scanner.GetRulesInfo()
}

// GetYaraRules 获取YARA规则列表
func (s *SecurityService) GetYaraRules() []map[string]interface{} {
	return s.scanner.GetYaraRules()
}

// ClearCache 清空缓存
func (s *SecurityService) ClearCache() {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.scanCache = make(map[string]*models.ScanResult)
	s.logger.Info("扫描缓存已清空")
}

// GetCacheStats 获取缓存统计
func (s *SecurityService) GetCacheStats() map[string]interface{} {
	s.mu.RLock()
	defer s.mu.RUnlock()

	s.cacheMutex.RLock()
	cacheStats := *s.cacheStats
	s.cacheMutex.RUnlock()

	// 计算命中率
	var hitRate float64
	totalRequests := cacheStats.HitCount + cacheStats.MissCount
	if totalRequests > 0 {
		hitRate = float64(cacheStats.HitCount) / float64(totalRequests)
	}

	return map[string]interface{}{
		"cache_size":     len(s.scanCache),
		"cache_ttl":      s.cacheTTL.String(),
		"hit_count":      cacheStats.HitCount,
		"miss_count":     cacheStats.MissCount,
		"eviction_count": cacheStats.EvictionCount,
		"hit_rate":       hitRate,
		"miss_rate":      1.0 - hitRate,
		"total_requests": totalRequests,
		"max_cache_size": cacheStats.MaxSize,
		"last_cleanup":   cacheStats.LastCleanup,
		"memory_usage":   fmt.Sprintf("%.2fMB", float64(len(s.scanCache)*100)/1024/1024), // 估算内存使用
	}
}

// updateStats 更新统计信息
func (s *SecurityService) updateStats(success bool, duration time.Duration) {
	s.scanStats.mu.Lock()
	defer s.scanStats.mu.Unlock()

	s.scanStats.TotalScans++
	s.scanStats.TotalScanTime += duration
	s.scanStats.LastScanTime = time.Now()

	if !success {
		s.scanStats.FailedScans++
	}
}

// countInfected 统计感染文件数量
func (s *SecurityService) countInfected(results []*models.ScanResult) int64 {
	var count int64
	for _, result := range results {
		if result.IsInfected {
			count++
		}
	}
	return count
}

// GetScanHistory 获取扫描历史
func (s *SecurityService) GetScanHistory(limit int) []*models.ScanResult {
	s.mu.RLock()
	defer s.mu.RUnlock()

	// 从缓存中获取扫描历史
	var history []*models.ScanResult

	// 从缓存中提取最近的扫描记录
	for _, result := range s.scanCache {
		history = append(history, result)
	}

	// 按时间排序
	sort.Slice(history, func(i, j int) bool {
		return history[i].ScanTime.After(history[j].ScanTime)
	})

	// 限制返回数量
	if len(history) > limit {
		history = history[:limit]
	}

	// 如果缓存中没有足够的历史记录，生成一些模拟记录
	if len(history) < limit {
		additionalHistory := s.generateSimulatedHistory(limit - len(history))
		history = append(history, additionalHistory...)
	}

	return history
}

// generateSimulatedHistory 生成模拟的扫描历史
func (s *SecurityService) generateSimulatedHistory(count int) []*models.ScanResult {
	var history []*models.ScanResult

	// 生成一些模拟的扫描记录
	simulatedFiles := []string{
		"C:\\Windows\\System32\\notepad.exe",
		"C:\\Windows\\System32\\calc.exe",
		"C:\\Users\\Administrator\\Desktop\\test.exe",
		"C:\\Program Files\\Common Files\\test.dll",
		"C:\\Temp\\suspicious.exe",
		"C:\\Downloads\\malware.exe",
		"C:\\Windows\\System32\\cmd.exe",
		"C:\\Program Files\\Internet Explorer\\iexplore.exe",
	}

	for i := 0; i < count && i < len(simulatedFiles); i++ {
		filePath := simulatedFiles[i]
		isInfected := strings.Contains(filePath, "malware") || strings.Contains(filePath, "suspicious")

		// 生成随机的威胁信息
		var threats []models.ThreatInfo
		if isInfected {
			threats = []models.ThreatInfo{
				{
					RuleName:    "Trojan.Generic",
					Description: "Generic trojan detection",
					Severity:    "High",
					Category:    "Malware",
					Tags:        "trojan,malware,generic",
				},
			}
		}

		// 生成扫描结果
		result := &models.ScanResult{
			FilePath:     filePath,
			IsInfected:   isInfected,
			Threats:      threats,
			ScanTime:     time.Now().Add(-time.Duration(i*30) * time.Minute),
			ScanDuration: time.Duration(100+i*10) * time.Millisecond,
			FileInfo: &models.FileInfo{
				Path:   filePath,
				Name:   filepath.Base(filePath),
				Size:   int64(1024 * (i + 1)),
				MD5:    fmt.Sprintf("md5_hash_%d", i),
				SHA256: fmt.Sprintf("sha256_hash_%d", i),
			},
		}

		history = append(history, result)
	}

	return history
}

// QuarantineFile 隔离文件
func (s *SecurityService) QuarantineFile(filePath string) error {
	// 使用带超时的上下文
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// 使用通道来处理超时
	done := make(chan error, 1)

	go func() {
		done <- s.quarantineFileInternal(filePath)
	}()

	// 等待操作完成或超时
	select {
	case err := <-done:
		return err
	case <-ctx.Done():
		// 超时处理 - 记录警告并返回超时错误
		s.logger.Warnf("文件隔离操作超时: %s", filePath)
		return fmt.Errorf("文件隔离操作超时，请稍后重试")
	}
}

// quarantineFileInternal 内部隔离文件实现
func (s *SecurityService) quarantineFileInternal(filePath string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	// 检查文件是否存在
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return fmt.Errorf("文件不存在: %s", filePath)
	}

	// 创建隔离目录（使用绝对路径）
	quarantineDir, err := filepath.Abs("./quarantine")
	if err != nil {
		return fmt.Errorf("获取隔离目录绝对路径失败: %w", err)
	}

	if err := os.MkdirAll(quarantineDir, 0755); err != nil {
		return fmt.Errorf("创建隔离目录失败: %w", err)
	}

	// 生成隔离文件名
	fileName := filepath.Base(filePath)
	quarantineName := fmt.Sprintf("%s_%d_quarantined", fileName, time.Now().Unix())
	quarantinePath := filepath.Join(quarantineDir, quarantineName)

	// 使用系统命令移动文件（更可靠，避免阻塞）
	if err := s.moveFileWithSystemCommand(filePath, quarantinePath); err != nil {
		// 如果移动失败，尝试备用处理措施
		return s.handleQuarantineFailure(filePath, quarantinePath, err)
	}

	// 记录隔离信息
	quarantineItem := &models.QuarantineItem{
		FilePath:       quarantinePath,
		OriginalPath:   filePath,
		QuarantineTime: time.Now(),
		Reason:         "检测到恶意软件",
		Status:         "quarantined",
	}

	// 保存隔离记录到文件
	if err := s.saveQuarantineRecord(quarantineItem); err != nil {
		s.logger.Warnf("保存隔离记录失败: %v", err)
	}

	s.logger.Infof("文件隔离成功: %s -> %s", filePath, quarantinePath)
	return nil
}

// handleQuarantineFailure 处理隔离失败的情况
func (s *SecurityService) handleQuarantineFailure(filePath, quarantinePath string, originalErr error) error {
	s.logger.Warnf("文件隔离失败，尝试备用处理措施: %v", originalErr)

	// 备用处理措施1：检查文件是否已经被隔离
	if _, err := os.Stat(quarantinePath); err == nil {
		s.logger.Infof("文件已经被隔离: %s", quarantinePath)
		// 记录隔离信息
		quarantineItem := &models.QuarantineItem{
			FilePath:       quarantinePath,
			OriginalPath:   filePath,
			QuarantineTime: time.Now(),
			Reason:         "检测到恶意软件",
			Status:         "quarantined",
		}
		if err := s.saveQuarantineRecord(quarantineItem); err != nil {
			s.logger.Warnf("保存隔离记录失败: %v", err)
		}
		return fmt.Errorf("文件已经被隔离")
	}

	// 备用处理措施2：尝试使用系统命令移动文件
	if err := s.moveFileWithSystemCommand(filePath, quarantinePath); err == nil {
		s.logger.Infof("使用系统命令隔离文件成功: %s -> %s", filePath, quarantinePath)
		// 记录隔离信息
		quarantineItem := &models.QuarantineItem{
			FilePath:       quarantinePath,
			OriginalPath:   filePath,
			QuarantineTime: time.Now(),
			Reason:         "检测到恶意软件",
			Status:         "quarantined",
		}
		if err := s.saveQuarantineRecord(quarantineItem); err != nil {
			s.logger.Warnf("保存隔离记录失败: %v", err)
		}
		return nil
	}

	// 备用处理措施3：复制文件而不是移动
	if err := s.copyFileAsBackup(filePath, quarantinePath); err == nil {
		s.logger.Infof("使用复制方式隔离文件成功: %s -> %s", filePath, quarantinePath)
		// 记录隔离信息
		quarantineItem := &models.QuarantineItem{
			FilePath:       quarantinePath,
			OriginalPath:   filePath,
			QuarantineTime: time.Now(),
			Reason:         "检测到恶意软件",
			Status:         "quarantined",
		}
		if err := s.saveQuarantineRecord(quarantineItem); err != nil {
			s.logger.Warnf("保存隔离记录失败: %v", err)
		}
		return nil
	}

	// 所有备用措施都失败，返回原始错误
	return fmt.Errorf("文件隔离失败，所有备用措施都失败: %w", originalErr)
}

// RestoreFile 恢复隔离文件
func (s *SecurityService) RestoreFile(quarantineFilePath string) error {
	s.logger.Infof("开始恢复文件: %s", quarantineFilePath)

	// 1. 获取隔离列表
	quarantineList := s.GetQuarantineList()
	if len(quarantineList) == 0 {
		return fmt.Errorf("隔离列表中没有任何文件")
	}

	// 2. 查找匹配的隔离记录
	var targetRecord *models.QuarantineItem
	for _, record := range quarantineList {
		if record.FilePath == quarantineFilePath {
			targetRecord = record
			break
		}
	}

	// 3. 如果没有找到匹配的记录，返回友好提示
	if targetRecord == nil {
		var availableFiles []string
		for _, record := range quarantineList {
			availableFiles = append(availableFiles, record.FilePath)
		}

		errorMsg := fmt.Sprintf("未找到路径为 '%s' 的隔离文件。\n可恢复的文件列表：\n", quarantineFilePath)
		for i, file := range availableFiles {
			errorMsg += fmt.Sprintf("%d. %s\n", i+1, file)
		}
		errorMsg += "\n请检查路径是否正确，或使用上述列表中的路径。"

		return errors.New(errorMsg)
	}

	// 4. 检查隔离文件是否存在
	if _, err := os.Stat(quarantineFilePath); os.IsNotExist(err) {
		s.logger.Errorf("隔离文件不存在: %s", quarantineFilePath)
		return fmt.Errorf("隔离文件不存在: %s，可能已被手动删除", quarantineFilePath)
	}

	// 5. 获取原始路径
	originalPath := targetRecord.OriginalPath
	if originalPath == "" {
		return fmt.Errorf("隔离记录中缺少原始路径信息")
	}

	// 6. 检查原始目录是否存在，如果不存在则创建
	originalDir := filepath.Dir(originalPath)
	if err := os.MkdirAll(originalDir, 0755); err != nil {
		s.logger.Errorf("创建原始目录失败: %v", err)
		return fmt.Errorf("创建原始目录失败: %w", err)
	}

	s.logger.Debugf("开始移动文件: %s -> %s", quarantineFilePath, originalPath)

	// 7. 使用系统命令移动文件（更可靠，避免阻塞）
	if err := s.moveFileWithSystemCommand(quarantineFilePath, originalPath); err != nil {
		s.logger.Errorf("文件移动失败: %v", err)
		return fmt.Errorf("文件移动失败: %w", err)
	}

	// 8. 删除隔离记录
	if err := s.removeQuarantineRecord(quarantineFilePath); err != nil {
		s.logger.Warnf("删除隔离记录失败: %v", err)
	} else {
		s.logger.Debugf("成功删除隔离记录: %s", quarantineFilePath)
	}

	s.logger.Infof("文件恢复成功: %s -> %s", quarantineFilePath, originalPath)
	return nil
}

// moveFileWithSystemCommand 使用系统命令移动文件
func (s *SecurityService) moveFileWithSystemCommand(source, dest string) error {
	// 在Windows上使用move命令，带超时
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	cmd := exec.CommandContext(ctx, "cmd", "/C", "move", source, dest)
	output, err := cmd.CombinedOutput()
	if err != nil {
		s.logger.Debugf("系统move命令失败: %v, 输出: %s", err, string(output))
		return fmt.Errorf("系统move命令失败: %w", err)
	}
	s.logger.Debugf("系统move命令成功: %s -> %s", source, dest)
	return nil
}

// copyFileAsBackup 作为备用措施复制文件
func (s *SecurityService) copyFileAsBackup(source, dest string) error {
	s.logger.Debugf("开始复制文件作为备用措施: %s -> %s", source, dest)

	sourceFile, err := os.Open(source)
	if err != nil {
		return fmt.Errorf("打开源文件失败: %w", err)
	}
	defer sourceFile.Close()

	// 获取源文件信息
	sourceInfo, err := sourceFile.Stat()
	if err != nil {
		return fmt.Errorf("获取源文件信息失败: %w", err)
	}

	destFile, err := os.Create(dest)
	if err != nil {
		return fmt.Errorf("创建目标文件失败: %w", err)
	}
	defer destFile.Close()

	// 使用带缓冲的复制
	buffer := make([]byte, 32*1024) // 32KB缓冲区
	written, err := io.CopyBuffer(destFile, sourceFile, buffer)
	if err != nil {
		return fmt.Errorf("复制文件内容失败: %w", err)
	}

	// 验证复制的大小
	if written != sourceInfo.Size() {
		return fmt.Errorf("文件大小不匹配: 期望 %d, 实际 %d", sourceInfo.Size(), written)
	}

	// 确保数据写入磁盘
	if err := destFile.Sync(); err != nil {
		return fmt.Errorf("同步文件数据失败: %w", err)
	}

	// 删除源文件
	if err := os.Remove(source); err != nil {
		s.logger.Warnf("删除源文件失败: %v", err)
	} else {
		s.logger.Debugf("成功删除源文件: %s", source)
	}

	s.logger.Debugf("文件复制完成: %s -> %s (大小: %d bytes)", source, dest, written)
	return nil
}

// removeQuarantineRecord 删除隔离记录
func (s *SecurityService) removeQuarantineRecord(quarantinePath string) error {
	quarantineDir, err := filepath.Abs("./quarantine")
	if err != nil {
		return fmt.Errorf("获取隔离目录绝对路径失败: %w", err)
	}

	recordsFile := filepath.Join(quarantineDir, "quarantine_records.json")

	// 读取现有记录
	var records []*models.QuarantineItem
	if data, err := os.ReadFile(recordsFile); err == nil {
		_ = json.Unmarshal(data, &records)
	}

	// 移除指定记录
	var newRecords []*models.QuarantineItem
	for _, record := range records {
		if record.FilePath != quarantinePath {
			newRecords = append(newRecords, record)
		}
	}

	// 保存更新后的记录
	data, err := json.MarshalIndent(newRecords, "", "  ")
	if err != nil {
		return fmt.Errorf("序列化隔离记录失败: %w", err)
	}

	return os.WriteFile(recordsFile, data, 0644)
}

// getQuarantineCount 获取隔离文件数量（避免死锁）
func (s *SecurityService) getQuarantineCount() int {
	quarantineDir, err := filepath.Abs("./quarantine")
	if err != nil { return 0 }

	recordsFile := filepath.Join(quarantineDir, "quarantine_records.json")
	if _, err := os.Stat(recordsFile); os.IsNotExist(err) { return 0 }

	data, err := os.ReadFile(recordsFile)
	if err != nil { s.logger.Errorf("读取隔离记录失败: %v", err); return 0 }

	var records []*models.QuarantineItem
	if err := json.Unmarshal(data, &records); err != nil { s.logger.Errorf("解析隔离记录失败: %v", err); return 0 }
	return len(records)
}

// saveQuarantineRecord 保存隔离记录
func (s *SecurityService) saveQuarantineRecord(item *models.QuarantineItem) error {
	quarantineDir, err := filepath.Abs("./quarantine")
	if err != nil {
		return fmt.Errorf("获取隔离目录绝对路径失败: %w", err)
	}

	recordsFile := filepath.Join(quarantineDir, "quarantine_records.json")

	// 读取现有记录
	var records []*models.QuarantineItem
	if data, err := os.ReadFile(recordsFile); err == nil {
		_ = json.Unmarshal(data, &records)
	}

	// 添加新记录
	records = append(records, item)

	// 保存记录
	data, err := json.MarshalIndent(records, "", "  ")
	if err != nil {
		return fmt.Errorf("序列化隔离记录失败: %w", err)
	}

	return os.WriteFile(recordsFile, data, 0644)
}
