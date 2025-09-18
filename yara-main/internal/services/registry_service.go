package services

import (
	"fmt"
	"sync"
	"time"

	"yara-security-service/internal/models"
	"yara-security-service/pkg/registry"

	"github.com/sirupsen/logrus"
)

// RegistryService 注册表服务
type RegistryService struct {
	manager *registry.Manager
	logger  *logrus.Logger
	mu      sync.RWMutex

	// 缓存
	registryCache map[string]*models.RegistryKey
	cacheTTL      time.Duration
}

// NewRegistryService 创建注册表服务
func NewRegistryService(manager *registry.Manager, logger *logrus.Logger) *RegistryService {
	return &RegistryService{
		manager:       manager,
		logger:        logger,
		registryCache: make(map[string]*models.RegistryKey),
		cacheTTL:      5 * time.Minute,
	}
}

// GetRegistryKey 获取注册表键
func (s *RegistryService) GetRegistryKey(path string) (*models.RegistryKey, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	// 检查缓存
	if cached, exists := s.registryCache[path]; exists {
		if time.Since(cached.LastModified) < s.cacheTTL {
			s.logger.Debugf("使用缓存结果: %s", path)
			return cached, nil
		}
		delete(s.registryCache, path)
	}

	// 获取注册表键
	key, err := s.manager.GetRegistryKey(path)
	if err != nil {
		return nil, fmt.Errorf("获取注册表键失败: %w", err)
	}

	// 缓存结果
	s.registryCache[path] = key

	s.logger.Infof("获取注册表键成功: %s", path)
	return key, nil
}

// CreateRegistryKey 创建注册表键
func (s *RegistryService) CreateRegistryKey(path string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	err := s.manager.CreateRegistryKey(path)
	if err != nil {
		return fmt.Errorf("创建注册表键失败: %w", err)
	}

	// 清除相关缓存
	delete(s.registryCache, path)

	s.logger.Infof("创建注册表键成功: %s", path)
	return nil
}

// DeleteRegistryKey 删除注册表键
func (s *RegistryService) DeleteRegistryKey(path string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	err := s.manager.DeleteRegistryKey(path)
	if err != nil {
		return fmt.Errorf("删除注册表键失败: %w", err)
	}

	// 清除相关缓存
	delete(s.registryCache, path)

	s.logger.Infof("删除注册表键成功: %s", path)
	return nil
}

// SetRegistryValue 设置注册表值
func (s *RegistryService) SetRegistryValue(path, name, valueType string, value interface{}) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	err := s.manager.SetRegistryValue(path, name, valueType, value)
	if err != nil {
		return fmt.Errorf("设置注册表值失败: %w", err)
	}

	// 清除相关缓存
	delete(s.registryCache, path)

	s.logger.Infof("设置注册表值成功: %s\\%s", path, name)
	return nil
}

// GetRegistryValue 获取注册表值
func (s *RegistryService) GetRegistryValue(path, name string) (interface{}, string, error) {
	value, valueType, err := s.manager.GetRegistryValue(path, name)
	if err != nil {
		return nil, "", fmt.Errorf("获取注册表值失败: %w", err)
	}

	s.logger.Debugf("获取注册表值成功: %s\\%s", path, name)
	return value, valueType, nil
}

// DeleteRegistryValue 删除注册表值
func (s *RegistryService) DeleteRegistryValue(path, name string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	err := s.manager.DeleteRegistryValue(path, name)
	if err != nil {
		return fmt.Errorf("删除注册表值失败: %w", err)
	}

	// 清除相关缓存
	delete(s.registryCache, path)

	s.logger.Infof("删除注册表值成功: %s\\%s", path, name)
	return nil
}

// ListRegistryKeys 列出注册表键
func (s *RegistryService) ListRegistryKeys(path string) ([]string, error) {
	keys, err := s.manager.ListRegistryKeys(path)
	if err != nil {
		return nil, fmt.Errorf("列出注册表键失败: %w", err)
	}

	s.logger.Debugf("列出注册表键成功: %s, 数量: %d", path, len(keys))
	return keys, nil
}

// ListRegistryValues 列出注册表值
func (s *RegistryService) ListRegistryValues(path string) ([]string, error) {
	values, err := s.manager.ListRegistryValues(path)
	if err != nil {
		return nil, fmt.Errorf("列出注册表值失败: %w", err)
	}

	s.logger.Debugf("列出注册表值成功: %s, 数量: %d", path, len(values))
	return values, nil
}

// SearchRegistry 搜索注册表
func (s *RegistryService) SearchRegistry(rootPath, searchTerm string) ([]*models.RegistryKey, error) {
	results, err := s.manager.SearchRegistry(rootPath, searchTerm)
	if err != nil {
		return nil, fmt.Errorf("搜索注册表失败: %w", err)
	}

	s.logger.Infof("搜索注册表成功: %s, 搜索词: %s, 结果数: %d", rootPath, searchTerm, len(results))
	return results, nil
}

// ClearCache 清空缓存
func (s *RegistryService) ClearCache() {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.registryCache = make(map[string]*models.RegistryKey)
	s.logger.Info("注册表缓存已清空")
}

// GetCacheStats 获取缓存统计
func (s *RegistryService) GetCacheStats() map[string]interface{} {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return map[string]interface{}{
		"cache_size": len(s.registryCache),
		"cache_ttl":  s.cacheTTL.String(),
	}
}

// ExportRegistryKey 导出注册表键
func (s *RegistryService) ExportRegistryKey(path string) (string, error) {
	exportData, err := s.manager.ExportRegistryKey(path)
	if err != nil {
		return "", fmt.Errorf("导出注册表键失败: %w", err)
	}

	s.logger.Infof("导出注册表键成功: %s", path)
	return exportData, nil
}

// ImportRegistryFile 导入注册表文件
func (s *RegistryService) ImportRegistryFile(filePath string) error {
	err := s.manager.ImportRegistryFile(filePath)
	if err != nil {
		return fmt.Errorf("导入注册表文件失败: %w", err)
	}

	// 清空缓存
	s.ClearCache()

	s.logger.Infof("导入注册表文件成功: %s", filePath)
	return nil
}

// GetRegistryStatistics 获取注册表统计
func (s *RegistryService) GetRegistryStatistics() map[string]interface{} {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return map[string]interface{}{
		"cache_size":    len(s.registryCache),
		"cache_ttl":     s.cacheTTL.String(),
		"last_accessed": time.Now(),
	}
}
