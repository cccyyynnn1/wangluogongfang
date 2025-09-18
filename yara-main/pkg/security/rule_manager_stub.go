package security

import (
  "fmt"
  "io/ioutil"
  "path/filepath"

  "github.com/sirupsen/logrus"
)

// RuleFileConfig 简单的规则文件描述（只包含文件名和优先级）
type RuleFileConfig struct {
  File     string
  Priority int
}

// IntegrationInfo 用于版本/集成信息
type IntegrationInfo struct {
  Version string
}

// RuleConfig 为 Scanner 提供的基本配置结构（可扩展）
type RuleConfig struct {
  Files           []RuleFileConfig
  IntegrationInfo IntegrationInfo
}

// RuleManager 是一个非常轻量的 stub，实现部分被调用的方法，
// 以便在缺少完整规则管理器实现时让应用回退到文件系统加载模式。
type RuleManager struct {
  rulesDir string
  logger   *logrus.Logger
}

// NewRuleManager 创建一个最小的规则管理器实例。如果 rulesDir 无效，此函数仍会返回实例，
// 但后续方法会返回空列表，Yara 引擎会回退到外部规则目录加载。
func NewRuleManager(rulesDir string, logger *logrus.Logger) (*RuleManager, error) {
  if logger == nil {
    logger = logrus.New()
  }
  rm := &RuleManager{rulesDir: rulesDir, logger: logger}
  logger.Infof("使用最小 RuleManager（stub），rulesDir=%s", rulesDir)
  return rm, nil
}

// GetEnabledRules 返回启用的规则文件列表（优先级信息）。
// 这个 stub 会尝试读取目录下的 .yar/.yara 文件作为启用规则的列表，若读取失败则返回空。
func (r *RuleManager) GetEnabledRules() []RuleFileConfig {
  // 尝试从 rulesDir 读取文件名
  if r.rulesDir == "" {
    return []RuleFileConfig{}
  }
  files, err := ioutil.ReadDir(r.rulesDir)
  if err != nil {
    return []RuleFileConfig{}
  }
  out := make([]RuleFileConfig, 0)
  for _, f := range files {
    if f.IsDir() { continue }
    name := f.Name()
    if filepath.Ext(name) == ".yar" || filepath.Ext(name) == ".yara" {
      out = append(out, RuleFileConfig{File: name, Priority: 100})
    }
  }
  return out
}

// GetRuleContent 读取某个规则文件内容
func (r *RuleManager) GetRuleContent(filename string) (string, error) {
  if r.rulesDir == "" {
    return "", fmt.Errorf("rulesDir not configured")
  }
  path := filepath.Join(r.rulesDir, filename)
  b, err := ioutil.ReadFile(path)
  if err != nil {
    return "", err
  }
  return string(b), nil
}

// GetConfig 返回简单的配置对象
func (r *RuleManager) GetConfig() *RuleConfig {
  files := r.GetEnabledRules()
  return &RuleConfig{Files: files, IntegrationInfo: IntegrationInfo{Version: "stub-0.0.1"}}
}

// ReloadRules stub
func (r *RuleManager) ReloadRules() error { return nil }

// GetRulesInfo 返回简单统计
func (r *RuleManager) GetRulesInfo() map[string]interface{} {
  files := r.GetEnabledRules()
  return map[string]interface{}{"total_files": len(files)}
}

// CheckForUpdates stub
// CheckForUpdates stub — 改为返回 (bool, error)
func (r *RuleManager) CheckForUpdates() (bool, error) { return false, nil }

// ValidateRuleFile stub — 改为返回 error
func (r *RuleManager) ValidateRuleFile(filename string) error { return nil }

// GetPerformanceStats stub
func (r *RuleManager) GetPerformanceStats() map[string]interface{} { return map[string]interface{}{} }
