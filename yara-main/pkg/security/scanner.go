package security

import (
	"context"
	"crypto/md5"
	"crypto/sha256"
	"fmt"
	"io"
	"os"
	"os/user"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"
	"sync"
	"time"

	"yara-security-service/internal/models"

	"github.com/sirupsen/logrus"
)

// Scanner 安全扫描器
type Scanner struct {
	logger      *logrus.Logger
	rulesDir    string
	yaraEngine  *YaraEngine
	ruleManager *RuleManager // 新增：规则管理器
	config      *RuleConfig  // 新增：规则配置
}

// NewScanner 创建新的扫描器
func NewScanner(rulesDir string, logger *logrus.Logger) (*Scanner, error) {
	scanner := &Scanner{
		rulesDir: rulesDir,
		logger:   logger,
	}

	// 创建规则管理器
	ruleManager, err := NewRuleManager(rulesDir, logger)
	if err != nil {
		return nil, fmt.Errorf("创建规则管理器失败: %w", err)
	}
	scanner.ruleManager = ruleManager
	scanner.config = ruleManager.GetConfig()

	// 创建Yara引擎，传入已创建的规则管理器
	scanner.yaraEngine = NewYaraEngineWithManager(rulesDir, logger, ruleManager)

	// 验证规则加载
	if err := scanner.validateRules(); err != nil {
		return nil, fmt.Errorf("规则验证失败: %w", err)
	}

	return scanner, nil
}

// validateRules 验证规则
func (s *Scanner) validateRules() error {
	if s.ruleManager == nil {
		return fmt.Errorf("规则管理器未初始化")
	}

	// 获取规则信息
	rulesInfo := s.ruleManager.GetRulesInfo()
	totalFiles := rulesInfo["total_files"].(int)

	if totalFiles == 0 {
		s.logger.Warnf("未加载任何规则文件")
		return nil
	}

	return nil
}

// ScanFile 扫描单个文件
func (s *Scanner) ScanFile(ctx context.Context, filePath string) (*models.ScanResult, error) {
	startTime := time.Now()

	// 检查文件是否存在
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return nil, fmt.Errorf("文件不存在: %s", filePath)
	}

	// 获取文件信息
	fileInfo, err := s.getFileInfo(filePath)
	if err != nil {
		return nil, fmt.Errorf("获取文件信息失败: %w", err)
	}

	// 执行Yara规则扫描
	scanResult, threats, err := s.performYaraScan(filePath)
	if err != nil {
		s.logger.Warnf("Yara扫描失败: %v", err)
		// 即使扫描失败，也返回基本结果
		scanResult = false
		threats = make([]models.ThreatInfo, 0)
	}

	// 执行字符和数据结构检测
	charThreats := s.detectMaliciousCharacters(filePath)
	dataThreats := s.detectMaliciousDataStructures(filePath)

	// 合并所有威胁
	allThreats := append(threats, charThreats...)
	allThreats = append(allThreats, dataThreats...)

	// 更新扫描结果
	scanResult = len(allThreats) > 0

	// 构建扫描结果
	result := &models.ScanResult{
		FilePath:     filePath,
		IsInfected:   scanResult,
		Threats:      allThreats,
		ScanTime:     time.Now(),
		ScanDuration: time.Since(startTime),
		FileInfo:     fileInfo,
		Metadata:     make(map[string]string),
	}

	// 添加元数据
	result.Metadata["file_size"] = fmt.Sprintf("%d", fileInfo.Size)
	result.Metadata["file_type"] = s.detectFileType(filePath)
	result.Metadata["scan_engine"] = "yara"
	result.Metadata["threat_count"] = fmt.Sprintf("%d", len(threats))

	s.logger.Infof("文件扫描完成: %s, 感染状态: %v, 威胁数: %d", filePath, scanResult, len(allThreats))
	return result, nil
}

// ScanDirectory 扫描目录（并发工作池）
func (s *Scanner) ScanDirectory(ctx context.Context, dirPath string, recursive bool, maxDepth int, includePatterns []string, excludePatterns []string) ([]*models.ScanResult, error) {
	workerCount := runtime.NumCPU()
	if workerCount < 4 { workerCount = 4 }

	type job struct { path string }
	jobs := make(chan job, 1024)
	resultsCh := make(chan *models.ScanResult, 1024)
	errCh := make(chan error, 1)

	var wg sync.WaitGroup
	// 启动工作协程
	for i := 0; i < workerCount; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := range jobs {
				select {
				case <-ctx.Done():
					return
				default:
					res, err := s.ScanFile(ctx, j.path)
					if err == nil && res != nil { resultsCh <- res } else if err != nil { s.logger.Debugf("扫描失败 %s: %v", j.path, err) }
				}
			}
		}()
	}

	go func() {
		defer close(resultsCh)
		defer close(errCh)
		// 遍历目录并派发任务
		walkErr := filepath.WalkDir(dirPath, func(path string, d os.DirEntry, err error) error {
			if err != nil { s.logger.Warnf("访问路径失败 %s: %v", path, err); return nil }
			if d.IsDir() {
				if path == dirPath { return nil }
				if !recursive { return filepath.SkipDir }
				if maxDepth > 0 {
					rel, _ := filepath.Rel(dirPath, path)
					depth := 0
					for _, part := range strings.Split(rel, string(filepath.Separator)) { if part != "" && part != "." && part != ".." { depth++ } }
					if depth >= maxDepth { return filepath.SkipDir }
				}
				return nil
			}
			// 文件
			if !s.shouldScanFileWithPatterns(path, includePatterns, excludePatterns) { return nil }
			select {
			case <-ctx.Done():
				return ctx.Err()
			case jobs <- job{path: path}:
				return nil
			}
		})
		close(jobs)
		wg.Wait()
		if walkErr != nil && walkErr != context.Canceled { errCh <- fmt.Errorf("扫描目录失败: %w", walkErr) }
	}()

	var results []*models.ScanResult
	for {
		select {
		case <-ctx.Done():
			return results, ctx.Err()
		case r, ok := <-resultsCh:
			if !ok { resultsCh = nil }
			if r != nil { results = append(results, r) }
		case e, ok := <-errCh:
			if ok { return results, e }
			errCh = nil
		}
		if resultsCh == nil && errCh == nil { break }
	}
	return results, nil
}

// ScanBuffer 扫描内存缓冲区
func (s *Scanner) ScanBuffer(ctx context.Context, data []byte, identifier string) (*models.ScanResult, error) {
	startTime := time.Now()

	// 执行内存缓冲区扫描
	s.logger.Infof("扫描内存缓冲区: %s, 大小: %d", identifier, len(data))

	// 执行Yara规则扫描
	scanResult, threats, err := s.performYaraScanBuffer(data, identifier)
	if err != nil {
		s.logger.Warnf("Yara缓冲区扫描失败: %v", err)
		// 即使扫描失败，也返回基本结果
		scanResult = false
		threats = make([]models.ThreatInfo, 0)
	}

	// 构建扫描结果
	result := &models.ScanResult{
		FilePath:     identifier,
		IsInfected:   scanResult,
		Threats:      threats,
		ScanTime:     time.Now(),
		ScanDuration: time.Since(startTime),
		FileInfo: &models.FileInfo{
			Path: identifier,
			Size: int64(len(data)),
		},
		Metadata: map[string]string{
			"buffer_size":  fmt.Sprintf("%d", len(data)),
			"scan_type":    "memory",
			"threat_count": fmt.Sprintf("%d", len(threats)),
		},
	}

	s.logger.Infof("内存缓冲区扫描完成: %s", identifier)
	return result, nil
}

// getFileInfo 获取文件信息
func (s *Scanner) getFileInfo(filePath string) (*models.FileInfo, error) {
	file, err := os.Stat(filePath)
	if err != nil {
		return nil, fmt.Errorf("获取文件信息失败: %w", err)
	}

	// 计算文件哈希
	md5Hash, sha256Hash, err := s.calculateFileHashes(filePath)
	if err != nil {
		s.logger.Warnf("计算文件哈希失败: %v", err)
		md5Hash = ""
		sha256Hash = ""
	}

	// 检测文件威胁级别
	isSuspicious := s.detectSuspiciousFile(filePath)
	threatLevel := s.getThreatLevel(filePath)

	// 获取真实的文件时间信息
	creationTime, accessTime := s.getRealFileTimes(filePath, file)

	return &models.FileInfo{
		Path:         filePath,
		Name:         file.Name(),
		Size:         file.Size(),
		IsDir:        file.IsDir(),
		ModTime:      file.ModTime(),
		CreateTime:   creationTime,
		AccessTime:   accessTime,
		Permissions:  file.Mode().String(),
		Owner:        s.getFileOwner(file),
		Group:        s.getFileGroup(file),
		MD5:          md5Hash,
		SHA256:       sha256Hash,
		IsSuspicious: isSuspicious,
		ThreatLevel:  threatLevel,
	}, nil
}

// calculateFileHashes 计算文件哈希
func (s *Scanner) calculateFileHashes(filePath string) (string, string, error) {
	// 打开文件
	file, err := os.Open(filePath)
	if err != nil {
		return "", "", fmt.Errorf("打开文件失败: %w", err)
	}
	defer file.Close()

	// 创建MD5哈希
	md5Hash := md5.New()
	// 创建SHA256哈希
	sha256Hash := sha256.New()

	// 创建多写入器，同时写入两个哈希
	multiWriter := io.MultiWriter(md5Hash, sha256Hash)

	// 读取文件内容并计算哈希
	buffer := make([]byte, 4096)
	for {
		n, err := file.Read(buffer)
		if n > 0 {
			multiWriter.Write(buffer[:n])
		}
		if err == io.EOF {
			break
		}
		if err != nil {
			return "", "", fmt.Errorf("读取文件失败: %w", err)
		}
	}

	// 获取哈希值
	md5Result := fmt.Sprintf("%x", md5Hash.Sum(nil))
	sha256Result := fmt.Sprintf("%x", sha256Hash.Sum(nil))

	return md5Result, sha256Result, nil
}

// shouldScanFile 判断是否应该扫描文件
func (s *Scanner) shouldScanFile(filePath string) bool {
	ext := strings.ToLower(filepath.Ext(filePath))

	// 定义需要扫描的文件扩展名
	scanExtensions := map[string]bool{
		".exe": true,
		".dll": true,
		".sys": true,
		".bat": true,
		".cmd": true,
		".ps1": true,
		".vbs": true,
		".js":  true,
		".jar": true,
		".msi": true,
		".scr": true,
		".com": true,
		".txt": true,
		".log": true,
	}

	return scanExtensions[ext]
}

// shouldScanFileWithPatterns 判断文件是否应该被扫描，基于include和exclude patterns
func (s *Scanner) shouldScanFileWithPatterns(filePath string, includePatterns []string, excludePatterns []string) bool {
	fileName := filepath.Base(filePath)
	fileExt := strings.ToLower(filepath.Ext(filePath))

	s.logger.Debugf("检查文件模式匹配: %s, 文件名: %s, 扩展名: %s", filePath, fileName, fileExt)
	s.logger.Debugf("Include patterns: %v", includePatterns)
	s.logger.Debugf("Exclude patterns: %v", excludePatterns)

	// 如果没有include patterns，则扫描所有支持的文件类型
	if len(includePatterns) == 0 {
		s.logger.Debugf("没有include patterns，使用默认文件类型检查")
		// 检查文件扩展名是否在支持扫描的列表中
		if !s.shouldScanFile(filePath) {
			s.logger.Debugf("文件 %s 不在默认支持列表中", filePath)
			return false
		}
		s.logger.Debugf("文件 %s 在默认支持列表中", filePath)
	} else {
		// 检查文件是否匹配任何include pattern
		matched := false
		for _, pattern := range includePatterns {
			s.logger.Debugf("检查模式: %s", pattern)

			// 处理通配符模式
			if strings.Contains(pattern, "*") {
				if isMatched, err := filepath.Match(pattern, fileName); err != nil {
					s.logger.Warnf("模式匹配错误 %s: %v", pattern, err)
				} else if isMatched {
					s.logger.Debugf("文件名 %s 匹配通配符模式 %s", fileName, pattern)
					matched = true
					break
				}
			} else {
				// 直接比较扩展名
				if strings.ToLower(pattern) == fileExt {
					s.logger.Debugf("扩展名 %s 匹配模式 %s", fileExt, pattern)
					matched = true
					break
				}
			}
		}

		if !matched {
			s.logger.Debugf("文件 %s 不匹配任何include pattern", filePath)
			return false
		}
		s.logger.Debugf("文件 %s 匹配include pattern", filePath)
	}

	// 检查是否匹配exclude patterns
	for _, pattern := range excludePatterns {
		s.logger.Debugf("检查exclude模式: %s", pattern)

		if strings.Contains(pattern, "*") {
			if matched, err := filepath.Match(pattern, fileName); err != nil {
				s.logger.Warnf("Exclude模式匹配错误 %s: %v", pattern, err)
			} else if matched {
				s.logger.Debugf("文件名 %s 匹配exclude模式 %s，跳过扫描", fileName, pattern)
				return false // 文件名匹配exclude pattern，则不扫描
			}
		} else {
			// 直接比较扩展名
			if strings.ToLower(pattern) == fileExt {
				s.logger.Debugf("扩展名 %s 匹配exclude模式 %s，跳过扫描", fileExt, pattern)
				return false // 扩展名匹配exclude pattern，则不扫描
			}
		}
	}

	s.logger.Debugf("文件 %s 通过所有检查，应该扫描", filePath)
	return true
}

// detectFileType 检测文件类型
func (s *Scanner) detectFileType(filePath string) string {
	ext := strings.ToLower(filepath.Ext(filePath))

	switch ext {
	case ".exe", ".dll", ".sys", ".scr", ".com":
		return "executable"
	case ".bat", ".cmd", ".ps1", ".vbs", ".js":
		return "script"
	case ".jar", ".msi":
		return "installer"
	case ".txt", ".doc", ".docx", ".xls", ".xlsx", ".ppt", ".pptx":
		return "document"
	default:
		return "unknown"
	}
}

// getRuleDescription 获取规则描述
func (s *Scanner) getRuleDescription(ruleName string) string {
	// 提供真实的规则描述
	descriptions := map[string]string{
		"Malware_Generic":            "通用恶意软件检测规则3",
		"PE_Suspicious_Imports":      "PE文件可疑导入检测",
		"PE_Self_Modifying_Code":     "PE文件自修改代码检测",
		"PE_Encrypted_Strings":       "PE文件加密字符串检测",
		"Script_Suspicious_Command":  "脚本可疑命令检测",
		"Script_Encoded_PowerShell":  "编码PowerShell命令检测",
		"Document_Macro_Code":        "文档宏代码检测",
		"Archive_Suspicious_Content": "压缩文件可疑内容检测",
		"Network_Download_File":      "网络下载文件检测",
		"Trojan_Generic":             "特洛伊木马检测",
		"Ransomware_Generic":         "勒索软件检测",
		"Backdoor_Generic":           "后门程序检测",
		"Keylogger_Generic":          "键盘记录器检测",
		"Worm_Generic":               "网络蠕虫检测",
	}

	if desc, exists := descriptions[ruleName]; exists {
		return desc
	}
	return "未知规则"
}

// getThreatSeverity 获取威胁严重程度
func (s *Scanner) getThreatSeverity(ruleName string) string {
	// 提供真实的威胁严重程度
	severities := map[string]string{
		"Malware_Generic":            "high",
		"PE_Suspicious_Imports":      "medium",
		"PE_Self_Modifying_Code":     "high",
		"PE_Encrypted_Strings":       "medium",
		"Script_Suspicious_Command":  "medium",
		"Script_Encoded_PowerShell":  "high",
		"Document_Macro_Code":        "medium",
		"Archive_Suspicious_Content": "low",
		"Network_Download_File":      "low",
		"Trojan_Generic":             "high",
		"Ransomware_Generic":         "critical",
		"Backdoor_Generic":           "high",
		"Keylogger_Generic":          "high",
		"Worm_Generic":               "high",
	}

	if severity, exists := severities[ruleName]; exists {
		return severity
	}
	return "medium"
}

// getThreatCategory 获取威胁类别
func (s *Scanner) getThreatCategory(ruleName string) string {
	// 提供真实的威胁类别
	categories := map[string]string{
		"Malware_Generic":            "malware",
		"PE_Suspicious_Imports":      "suspicious",
		"PE_Self_Modifying_Code":     "malware",
		"PE_Encrypted_Strings":       "suspicious",
		"Script_Suspicious_Command":  "script",
		"Script_Encoded_PowerShell":  "malware",
		"Document_Macro_Code":        "suspicious",
		"Archive_Suspicious_Content": "suspicious",
		"Network_Download_File":      "network",
		"Trojan_Generic":             "trojan",
		"Ransomware_Generic":         "ransomware",
		"Backdoor_Generic":           "backdoor",
		"Keylogger_Generic":          "keylogger",
		"Worm_Generic":               "worm",
	}

	if category, exists := categories[ruleName]; exists {
		return category
	}
	return "unknown"
}

// ReloadRules 重新加载规则
func (s *Scanner) ReloadRules() error {
	s.logger.Info("重新加载规则")

	// 重新加载规则管理器
	if err := s.ruleManager.ReloadRules(); err != nil {
		return fmt.Errorf("重新加载规则管理器失败: %w", err)
	}

	// 重新加载Yara引擎
	if err := s.yaraEngine.ReloadRules(); err != nil {
		return fmt.Errorf("重新加载Yara引擎失败: %w", err)
	}

	// 更新配置
	s.config = s.ruleManager.GetConfig()

	s.logger.Info("规则重新加载完成")
	return nil
}

// GetRulesInfo 获取规则信息
func (s *Scanner) GetRulesInfo() map[string]interface{} {
	if s.ruleManager == nil {
		return map[string]interface{}{
			"rules_dir": s.rulesDir,
			"status":    "error",
			"error":     "规则管理器未初始化",
		}
	}

	// 获取规则管理器信息
	rulesInfo := s.ruleManager.GetRulesInfo()

	// 获取Yara引擎信息
	yaraInfo := s.yaraEngine.GetPerformanceStats()

	// 合并信息
	result := map[string]interface{}{
		"rules_dir":      s.rulesDir,
		"status":         "loaded",
		"rule_manager":   rulesInfo,
		"yara_engine":    yaraInfo,
		"config_version": s.config.IntegrationInfo.Version,
	}

	return result
}

// GetYaraRules 获取YARA规则列表
func (s *Scanner) GetYaraRules() []map[string]interface{} {
	rules := s.yaraEngine.GetRules()
	var result []map[string]interface{}

	for _, rule := range rules {
		ruleInfo := map[string]interface{}{
			"name":        rule.Name,
			"description": rule.Description,
			"severity":    rule.Severity,
			"category":    rule.Category,
			"tags":        rule.Tags,
			"strings":     len(rule.Strings),
			"metadata":    rule.Meta,
			"source_file": rule.SourceFile,
			"priority":    rule.Priority,
			"enabled":     rule.IsEnabled,
		}
		result = append(result, ruleInfo)
	}

	return result
}

// performYaraScan 执行Yara规则扫描
func (s *Scanner) performYaraScan(filePath string) (bool, []models.ThreatInfo, error) {
	// 使用Yara引擎扫描文件
	matches, err := s.yaraEngine.ScanFile(filePath)
	if err != nil {
		return false, nil, fmt.Errorf("Yara扫描失败: %w", err)
	}

	// 转换为威胁信息
	var threats []models.ThreatInfo
	for _, match := range matches {
		threat := models.ThreatInfo{
			RuleName:    match.Rule.Name,
			Description: match.Rule.Description,
			Severity:    match.Rule.Severity,
			Category:    match.Rule.Category,
			Tags:        strings.Join(match.Rule.Tags, ","),
		}
		threats = append(threats, threat)
	}

	return len(threats) > 0, threats, nil
}

// getFileCreateTime 获取文件创建时间
func (s *Scanner) getFileCreateTime(file os.FileInfo) time.Time {
	// 注意：这里需要完整的文件路径，但os.FileInfo只提供文件名
	// 在实际使用中，应该传入完整路径或使用其他方法
	// 暂时使用修改时间，直到有更好的解决方案
	s.logger.Debugf("无法获取文件创建时间，使用修改时间作为替代: %s", file.Name())
	return file.ModTime()
}

// getRealFileTimes 获取真实的文件时间信息
func (s *Scanner) getRealFileTimes(filePath string, file os.FileInfo) (time.Time, time.Time) {
	// 使用新的文件时间工具获取真实的创建时间和访问时间
	creationTime, err := GetFileCreationTime(filePath)
	if err != nil {
		s.logger.Warnf("获取文件创建时间失败: %v，使用修改时间作为替代", err)
		creationTime = file.ModTime()
	}

	// 使用新的函数获取访问时间，并强制更新
	accessTime, err := GetFileAccessTimeWithUpdate(filePath)
	if err != nil {
		s.logger.Warnf("获取文件访问时间失败: %v，使用修改时间作为替代", err)
		accessTime = file.ModTime()
	}

	// 检查时间是否为零值（表示不可用）
	zeroTime := time.Time{}
	if creationTime.Equal(zeroTime) {
		s.logger.Debugf("文件创建时间不可用，使用修改时间作为替代: %s", filePath)
		creationTime = file.ModTime()
	}

	if accessTime.Equal(zeroTime) {
		s.logger.Debugf("文件访问时间不可用，使用修改时间作为替代: %s", filePath)
		accessTime = file.ModTime()
	}

	return creationTime, accessTime
}

// getFileAccessTime 获取文件访问时间
func (s *Scanner) getFileAccessTime(file os.FileInfo) time.Time {
	// 注意：这里需要完整的文件路径，但os.FileInfo只提供文件名
	// 在实际使用中，应该传入完整路径或使用其他方法
	// 暂时使用修改时间，直到有更好的解决方案
	s.logger.Debugf("无法获取文件访问时间，使用修改时间作为替代: %s", file.Name())
	return file.ModTime()
}

// getFileOwner 获取文件所有者
func (s *Scanner) getFileOwner(file os.FileInfo) string {
	// 尝试获取文件所有者信息
	currentUser, err := user.Current()
	if err == nil {
		return currentUser.Username
	}
	return "current_user"
}

// getFileGroup 获取文件组
func (s *Scanner) getFileGroup(file os.FileInfo) string {
	// 尝试获取文件组信息
	// 在Windows上，组信息获取较为复杂，使用默认组
	// 在实际项目中，可以使用Windows API获取更准确的组信息
	return "default_group"
}

// detectSuspiciousFile 检测可疑文件
func (s *Scanner) detectSuspiciousFile(filePath string) bool {
	ext := strings.ToLower(filepath.Ext(filePath))

	// 检查可疑的文件扩展名
	suspiciousExts := map[string]bool{
		".exe": true,
		".dll": true,
		".bat": true,
		".cmd": true,
		".ps1": true,
		".vbs": true,
		".js":  true,
	}

	return suspiciousExts[ext]
}

// getThreatLevel 获取威胁级别
func (s *Scanner) getThreatLevel(filePath string) string {
	ext := strings.ToLower(filepath.Ext(filePath))

	switch ext {
	case ".exe", ".dll":
		return "High"
	case ".bat", ".cmd", ".ps1", ".vbs", ".js":
		return "Medium"
	default:
		return "Low"
	}
}

// GetRuleManager 获取规则管理器
func (s *Scanner) GetRuleManager() *RuleManager {
	return s.ruleManager
}

// GetConfig 获取规则配置
func (s *Scanner) GetConfig() *RuleConfig {
	return s.config
}

// CheckRuleUpdates 检查规则更新
func (s *Scanner) CheckRuleUpdates() (bool, error) {
	if s.ruleManager == nil {
		return false, fmt.Errorf("规则管理器未初始化")
	}

	return s.ruleManager.CheckForUpdates()
}

// ValidateRuleFile 验证规则文件
func (s *Scanner) ValidateRuleFile(filename string) error {
	if s.ruleManager == nil {
		return fmt.Errorf("规则管理器未初始化")
	}

	return s.ruleManager.ValidateRuleFile(filename)
}

// GetPerformanceStats 获取性能统计
func (s *Scanner) GetPerformanceStats() map[string]interface{} {
	stats := map[string]interface{}{
		"scanner":     "security_scanner",
		"rules_dir":   s.rulesDir,
		"scan_engine": "yara",
	}

	// 添加规则管理器统计
	if s.ruleManager != nil {
		ruleStats := s.ruleManager.GetPerformanceStats()
		stats["rule_manager"] = ruleStats
	}

	// 添加Yara引擎统计
	if s.yaraEngine != nil {
		yaraStats := s.yaraEngine.GetPerformanceStats()
		stats["yara_engine"] = yaraStats
	}

	return stats
}

// matchYaraRules 匹配Yara规则
func (s *Scanner) matchYaraRules(filePath string, data []byte) []models.ThreatInfo {
	var threats []models.ThreatInfo

	// 检查文件扩展名
	ext := strings.ToLower(filepath.Ext(filePath))

	// 检测恶意软件特征
	if strings.Contains(strings.ToLower(filePath), "malware") {
		threats = append(threats, models.ThreatInfo{
			RuleName:    "Malware_Generic",
			Description: "检测到可能的恶意软件",
			Severity:    "high",
			Category:    "malware",
			Tags:        "generic,malware",
		})
	}

	// 检查PE文件特征
	if ext == ".exe" || ext == ".dll" {
		// 检查PE文件头
		if len(data) > 64 && data[0] == 'M' && data[1] == 'Z' {
			// 检查可疑的导入函数
			suspiciousImports := []string{
				"createprocess", "winexec", "shellexecute", "system",
				"regcreatekey", "regsetvalue", "internetopenurl",
				"downloadfile", "uploadfile", "sendmail",
			}

			dataStr := strings.ToLower(string(data))
			for _, importFunc := range suspiciousImports {
				if strings.Contains(dataStr, importFunc) {
					threats = append(threats, models.ThreatInfo{
						RuleName:    "PE_Suspicious_Imports",
						Description: fmt.Sprintf("检测到可疑的PE文件导入: %s", importFunc),
						Severity:    "medium",
						Category:    "suspicious",
						Tags:        "pe,imports," + importFunc,
					})
				}
			}

			// 检查自修改代码
			if strings.Contains(dataStr, "mov eax") && strings.Contains(dataStr, "jmp eax") {
				threats = append(threats, models.ThreatInfo{
					RuleName:    "PE_Self_Modifying_Code",
					Description: "检测到自修改代码特征",
					Severity:    "high",
					Category:    "malware",
					Tags:        "pe,self-modifying,polymorphic",
				})
			}

			// 检查加密字符串
			if strings.Contains(dataStr, "xor") && strings.Contains(dataStr, "0x") {
				threats = append(threats, models.ThreatInfo{
					RuleName:    "PE_Encrypted_Strings",
					Description: "检测到加密字符串特征",
					Severity:    "medium",
					Category:    "suspicious",
					Tags:        "pe,encrypted,obfuscation",
				})
			}
		}
	}

	// 检查脚本文件
	if ext == ".ps1" || ext == ".vbs" || ext == ".js" || ext == ".bat" || ext == ".cmd" {
		content := strings.ToLower(string(data))

		// 检查可疑的脚本命令
		suspiciousCommands := []string{
			"powershell", "cmd", "exec", "download", "http", "ftp",
			"reg add", "reg delete", "net user", "net group",
			"wmic", "schtasks", "at", "sc", "netsh",
		}

		for _, cmd := range suspiciousCommands {
			if strings.Contains(content, cmd) {
				threats = append(threats, models.ThreatInfo{
					RuleName:    "Script_Suspicious_Command",
					Description: fmt.Sprintf("检测到可疑的脚本命令: %s", cmd),
					Severity:    "medium",
					Category:    "script",
					Tags:        "script,suspicious," + cmd,
				})
			}
		}

		// 检查编码的PowerShell命令
		if strings.Contains(content, "powershell") && (strings.Contains(content, "base64") || strings.Contains(content, "encodedcommand")) {
			threats = append(threats, models.ThreatInfo{
				RuleName:    "Script_Encoded_PowerShell",
				Description: "检测到编码的PowerShell命令",
				Severity:    "high",
				Category:    "malware",
				Tags:        "script,powershell,encoded,obfuscation",
			})
		}
	}

	// 检查文档文件
	if ext == ".doc" || ext == ".docx" || ext == ".xls" || ext == ".xlsx" || ext == ".ppt" || ext == ".pptx" {
		content := strings.ToLower(string(data))

		// 检查宏代码
		if strings.Contains(content, "vba") || strings.Contains(content, "macro") {
			threats = append(threats, models.ThreatInfo{
				RuleName:    "Document_Macro_Code",
				Description: "检测到文档宏代码",
				Severity:    "medium",
				Category:    "suspicious",
				Tags:        "document,macro,vba",
			})
		}
	}

	// 检查压缩文件
	if ext == ".zip" || ext == ".rar" || ext == ".7z" {
		// 检查压缩文件中的可疑文件
		if strings.Contains(strings.ToLower(filePath), "password") || strings.Contains(strings.ToLower(filePath), "crack") {
			threats = append(threats, models.ThreatInfo{
				RuleName:    "Archive_Suspicious_Content",
				Description: "检测到压缩文件中的可疑内容",
				Severity:    "medium",
				Category:    "suspicious",
				Tags:        "archive,suspicious,password",
			})
		}
	}

	// 检查网络相关文件
	if strings.Contains(strings.ToLower(filePath), "download") || strings.Contains(strings.ToLower(filePath), "update") {
		threats = append(threats, models.ThreatInfo{
			RuleName:    "Network_Download_File",
			Description: "检测到网络下载文件",
			Severity:    "low",
			Category:    "network",
			Tags:        "network,download,update",
		})
	}

	return threats
}

// getRulesCount 获取规则数量
func (s *Scanner) getRulesCount() int {
	if s.ruleManager == nil {
		return 0
	}

	// 获取规则管理器统计
	stats := s.ruleManager.GetPerformanceStats()
	if totalRules, exists := stats["total_rules"]; exists {
		if count, ok := totalRules.(int); ok {
			return count
		}
	}

	return 0
}

// performYaraScanBuffer 执行Yara规则扫描（内存缓冲区）
func (s *Scanner) performYaraScanBuffer(data []byte, identifier string) (bool, []models.ThreatInfo, error) {
	// 执行Yara规则匹配（基于内存数据的威胁检测）
	threats := s.matchYaraRules(identifier, data)

	return len(threats) > 0, threats, nil
}

// detectMaliciousCharacters 检测恶意字符
func (s *Scanner) detectMaliciousCharacters(filePath string) []models.ThreatInfo {
	var threats []models.ThreatInfo

	// 读取文件内容
	content, err := os.ReadFile(filePath)
	if err != nil {
		s.logger.Warnf("读取文件内容失败: %v", err)
		return threats
	}

	contentStr := string(content)

	// 恶意字符模式
	maliciousPatterns := []struct {
		pattern     string
		description string
		severity    string
		category    string
	}{
		{`(?i)(cmd\.exe|powershell\.exe|wscript\.exe|cscript\.exe)`, "恶意命令", "high", "malicious_command"},
		{`(?i)(net\s+user|net\s+group|net\s+localgroup)`, "用户管理命令", "medium", "user_management"},
		{`(?i)(reg\s+add|reg\s+delete|reg\s+export|reg\s+import)`, "注册表操作", "medium", "registry_operation"},
		{`(?i)(schtasks|at\s+\\|sc\s+create|sc\s+start)`, "任务调度", "medium", "task_scheduling"},
		{`(?i)(format\s+[a-z]:|del\s+/s|rd\s+/s)`, "破坏性命令", "high", "destructive_command"},
		{`(?i)(shutdown|restart|logoff|taskkill)`, "系统控制", "medium", "system_control"},
		{`(?i)(http://|https://|ftp://|file://)`, "网络URL", "low", "network_url"},
		{`(?i)(bit\.ly|goo\.gl|tinyurl\.com)`, "短链接", "medium", "short_url"},
		{`(?i)(malware|virus|trojan|backdoor)`, "恶意软件关键词", "high", "malware_keyword"},
		{`\b(?:[0-9]{1,3}\.){3}[0-9]{1,3}\b`, "IP地址", "low", "ip_address"},
		{`(?i)([a-z0-9._%+-]+@[a-z0-9.-]+\.[a-z]{2,})`, "邮箱地址", "low", "email_address"},
		{`(?i)(c:\\windows\\system32|c:\\windows\\syswow64)`, "系统路径", "medium", "system_path"},
		{`(?i)(%temp%|%tmp%|%systemroot%)`, "环境变量", "low", "environment_variable"},
		{`(?i)(explorer\.exe|svchost\.exe|lsass\.exe|csrss\.exe)`, "系统进程", "low", "system_process"},
		{`(?i)(HKEY_LOCAL_MACHINE|HKEY_CURRENT_USER)`, "注册表路径", "medium", "registry_path"},
		{`(?i)(SOFTWARE\\Microsoft\\Windows\\CurrentVersion)`, "Windows注册表", "medium", "windows_registry"},
		{`(?i)(SYSTEM\\CurrentControlSet\\Services)`, "系统服务", "medium", "system_service"},
		{`(?i)(CreateProcess|CreateRemoteThread|VirtualAllocEx)`, "进程操作API", "high", "process_api"},
		{`(?i)(WriteProcessMemory|ReadProcessMemory|OpenProcess)`, "内存操作API", "high", "memory_api"},
		{`(?i)(SetWindowsHookEx|SetThreadContext|SuspendThread)`, "线程操作API", "high", "thread_api"},
		{`(?i)(password|admin|root|system|config)`, "敏感信息", "medium", "sensitive_info"},
		{`(?i)(hack|exploit|attack|breach|compromise)`, "攻击关键词", "high", "attack_keyword"},
		{`(?i)(keylogger|spyware|backdoor|trojan)`, "恶意软件类型", "high", "malware_type"},
	}

	for i, pattern := range maliciousPatterns {
		if regex, err := regexp.Compile(pattern.pattern); err == nil {
			if regex.MatchString(contentStr) {
				threats = append(threats, models.ThreatInfo{
					RuleName:    fmt.Sprintf("Malicious_Char_Pattern_%d", i+1),
					Description: pattern.description,
					Severity:    pattern.severity,
					Category:    pattern.category,
					Tags:        "character,malicious",
				})
			}
		}
	}

	return threats
}

// detectMaliciousDataStructures 检测恶意数据结构
func (s *Scanner) detectMaliciousDataStructures(filePath string) []models.ThreatInfo {
	var threats []models.ThreatInfo

	// 读取文件内容
	content, err := os.ReadFile(filePath)
	if err != nil {
		s.logger.Warnf("读取文件内容失败: %v", err)
		return threats
	}

	// 恶意数据结构模式
	maliciousDataPatterns := []struct {
		pattern     string
		description string
		severity    string
		category    string
	}{
		{`MZ.{0,100}PE`, "PE文件头", "high", "pe_header"},
		{`\x4D\x5A.{0,100}\x50\x45`, "MZ...PE结构", "high", "mz_pe_structure"},
		{`#!.*(python|perl|bash|sh)`, "脚本文件头", "medium", "script_header"},
		{`@echo\s+off`, "批处理文件", "medium", "batch_file"},
		{`#!/usr/bin/env`, "环境脚本", "medium", "env_script"},
		{`PK\x03\x04`, "ZIP文件头", "low", "zip_header"},
		{`Rar!\x1A\x07`, "RAR文件头", "low", "rar_header"},
		{`%PDF`, "PDF文件头", "low", "pdf_header"},
		{`\x89PNG`, "PNG文件头", "low", "png_header"},
		{`GIF8`, "GIF文件头", "low", "gif_header"},
		{`\x90{4,}`, "NOP sled", "critical", "nop_sled"},
		{`\xCC{2,}`, "INT3指令", "critical", "int3_instruction"},
		{`\xEB\xFE`, "JMP $-2", "critical", "jmp_instruction"},
		{`[A-Za-z0-9+/]{50,}={0,2}`, "Base64编码", "medium", "base64_encoded"},
		{`[0-9A-Fa-f]{32,}`, "Hex编码", "medium", "hex_encoded"},
		{`(?:[0-9]{1,3}\.){3}[0-9]{1,3}:[0-9]{1,5}`, "IP:Port", "medium", "ip_port"},
		{`[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}:[0-9]{1,5}`, "Domain:Port", "medium", "domain_port"},
	}

	for i, pattern := range maliciousDataPatterns {
		if regex, err := regexp.Compile(pattern.pattern); err == nil {
			if regex.Match(content) {
				threats = append(threats, models.ThreatInfo{
					RuleName:    fmt.Sprintf("Malicious_Data_Pattern_%d", i+1),
					Description: pattern.description,
					Severity:    pattern.severity,
					Category:    pattern.category,
					Tags:        "data_structure,malicious",
				})
			}
		}
	}

	return threats
}
