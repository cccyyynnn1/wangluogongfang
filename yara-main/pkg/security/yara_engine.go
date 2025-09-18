package security

import (
	"crypto/md5"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/sirupsen/logrus"
)

// YaraString Yara字符串定义
type YaraString struct {
	Identifier string
	Value      string
	Type       string   // "text", "hex", "regex"
	Modifiers  []string // "nocase", "wide", "ascii", etc.
}

// YaraCondition Yara条件定义
type YaraCondition struct {
	Expression string
	Type       string // "and", "or", "not", "simple"
}

// YaraRule Yara规则结构
type YaraRule struct {
	Name        string
	Description string
	Severity    string
	Category    string
	Tags        []string
	Meta        map[string]string
	Strings     []*YaraString
	Condition   *YaraCondition
	IsEnabled   bool
	SourceFile  string // 新增：来源文件
	Priority    int    // 新增：优先级
}

// YaraMatch Yara匹配结果
type YaraMatch struct {
	Rule           *YaraRule
	MatchedAt      []int64
	MatchedStrings []string
	Score          int
	Confidence     float64 // 新增：置信度
}

// YaraEngine Yara规则引擎
type YaraEngine struct {
	rules       []*YaraRule
	rulesDir    string
	logger      *logrus.Logger
	ruleManager *RuleManager               // 新增：规则管理器
	performance *PerformanceConfig         // 新增：性能配置
	scanCache   map[string]*ScanCacheEntry // 新增：扫描缓存
	cacheMutex  sync.RWMutex
	lastScanData []byte                    // 新增：最近一次扫描数据（仅用于近似 pe.imports 判断）
}

// PerformanceConfig 性能配置
type PerformanceConfig struct {
	MaxStringsPerRule      int
	MaxConditionsPerRule   int
	EnableCaching          bool
	CacheSize              int
	RulePriorityWeight     float64
	EnableParallelScanning bool
	EnableRuleOptimization bool
}

// ScanCacheEntry 扫描缓存条目
type ScanCacheEntry struct {
	Hash        string
	Result      []*YaraMatch
	Timestamp   time.Time
	AccessCount int
}

// NewYaraEngine 创建Yara规则引擎
func NewYaraEngine(rulesDir string, logger *logrus.Logger) *YaraEngine {
	engine := &YaraEngine{
		rules:     make([]*YaraRule, 0),
		rulesDir:  rulesDir,
		logger:    logger,
		scanCache: make(map[string]*ScanCacheEntry),
	}

	// 创建规则管理器
	ruleManager, err := NewRuleManager(rulesDir, logger)
	if err != nil {
		logger.Warnf("创建规则管理器失败: %v，使用传统模式", err)
		// 回退到传统模式
		engine.loadDefaultRules()
		engine.loadExternalRules()
	} else {
		engine.ruleManager = ruleManager
		engine.performance = &PerformanceConfig{
			MaxStringsPerRule:      200,
			MaxConditionsPerRule:   30,
			EnableCaching:          true,
			CacheSize:              3000,
			RulePriorityWeight:     0.9,
			EnableParallelScanning: true,
			EnableRuleOptimization: true,
		}

		// 从规则管理器加载规则
		engine.loadRulesFromManager()
	}

	return engine
}

// NewYaraEngineWithManager 使用已存在的规则管理器创建Yara规则引擎
func NewYaraEngineWithManager(rulesDir string, logger *logrus.Logger, ruleManager *RuleManager) *YaraEngine {
	engine := &YaraEngine{
		rules:       make([]*YaraRule, 0),
		rulesDir:    rulesDir,
		logger:      logger,
		ruleManager: ruleManager,
		scanCache:   make(map[string]*ScanCacheEntry),
		performance: &PerformanceConfig{
			MaxStringsPerRule:      200,
			MaxConditionsPerRule:   30,
			EnableCaching:          true,
			CacheSize:              3000,
			RulePriorityWeight:     0.9,
			EnableParallelScanning: true,
			EnableRuleOptimization: true,
		},
	}

	// 从规则管理器加载规则
	engine.loadRulesFromManager()

	return engine
}

// loadRulesFromManager 从规则管理器加载规则
func (engine *YaraEngine) loadRulesFromManager() {
	if engine.ruleManager == nil {
		engine.logger.Warn("规则管理器为空，无法加载规则")
		return
	}

	engine.logger.Info("从规则管理器加载规则...")

	// 获取所有启用的规则文件
	enabledRules := engine.ruleManager.GetEnabledRules()
	engine.logger.Infof("找到 %d 个启用的规则文件", len(enabledRules))

	if len(enabledRules) == 0 {
		engine.logger.Warn("没有找到启用的规则文件，回退到扫描规则目录加载 .yar 文件")
		engine.loadExternalRules()
		return
	}

	totalRulesLoaded := 0

	for _, ruleConfig := range enabledRules {
		engine.logger.Debugf("处理规则文件: %s (优先级: %d)", ruleConfig.File, ruleConfig.Priority)
		
		// 获取规则内容
		content, err := engine.ruleManager.GetRuleContent(ruleConfig.File)
		if err != nil {
			engine.logger.Warnf("获取规则内容失败 %s: %v", ruleConfig.File, err)
			continue
		}

		engine.logger.Debugf("规则文件 %s 内容长度: %d bytes", ruleConfig.File, len(content))

		// 解析规则
		rules, err := engine.parseYaraRules(string(content))
		if err != nil {
			engine.logger.Warnf("解析规则文件失败 %s: %v", ruleConfig.File, err)
			continue
		}

		engine.logger.Infof("从 %s 解析出 %d 条规则", ruleConfig.File, len(rules))

		if len(rules) == 0 {
			engine.logger.Warnf("规则文件 %s 解析后没有规则，可能格式不正确", ruleConfig.File)
			continue
		}

		// 设置规则来源和优先级
		for _, rule := range rules {
			rule.SourceFile = ruleConfig.File
			rule.Priority = ruleConfig.Priority
			engine.applyRuleDefaults(rule)
			engine.logger.Debugf("设置规则 %s 来源: %s, 优先级: %d", rule.Name, rule.SourceFile, rule.Priority)
		}

		engine.rules = append(engine.rules, rules...)
		totalRulesLoaded += len(rules)
		engine.logger.Infof("从 %s 加载了 %d 条规则", ruleConfig.File, len(rules))
	}

	// 按优先级排序规则
	engine.sortRulesByPriority()

	engine.logger.Infof("规则管理器模式：总共加载了 %d 条规则", len(engine.rules))
	
	// 输出规则摘要
	if len(engine.rules) > 0 {
		engine.logger.Info("已加载的规则摘要:")
		for i, rule := range engine.rules {
			if i < 5 { // 只显示前5条规则
				engine.logger.Infof("  %d. %s (来源: %s, 优先级: %d, 字符串数: %d)", 
					i+1, rule.Name, rule.SourceFile, rule.Priority, len(rule.Strings))
			}
		}
		if len(engine.rules) > 5 {
			engine.logger.Infof("  ... 还有 %d 条规则", len(engine.rules)-5)
		}
	}
}

// sortRulesByPriority 按优先级排序规则
func (engine *YaraEngine) sortRulesByPriority() {
	// 使用稳定的排序，保持相同优先级规则的相对顺序
	sort.SliceStable(engine.rules, func(i, j int) bool {
		return engine.rules[i].Priority < engine.rules[j].Priority
	})
}

// loadDefaultRules 加载默认规则（传统模式）
func (engine *YaraEngine) loadDefaultRules() {
	// 恶意软件检测规则
	engine.rules = append(engine.rules, &YaraRule{
		Name:        "Malware_Generic",
		Description: "通用恶意软件检测规则2",
		Severity:    "high",
		Category:    "malware2",
		Tags:        []string{"generic", "malware"},
		Meta: map[string]string{
			"author": "security-team",
			"date":   time.Now().Format("2006-01-02"),
		},
		Strings: []*YaraString{
			{Identifier: "$s1", Value: "CreateRemoteThread", Type: "text", Modifiers: []string{"nocase"}},
			{Identifier: "$s2", Value: "VirtualAllocEx", Type: "text", Modifiers: []string{"nocase"}},
			{Identifier: "$s3", Value: "WriteProcessMemory", Type: "text", Modifiers: []string{"nocase"}},
			{Identifier: "$s4", Value: "cmd.exe /c", Type: "text", Modifiers: []string{"nocase"}},
			{Identifier: "$s5", Value: "powershell.exe", Type: "text", Modifiers: []string{"nocase"}},
		},
		Condition: &YaraCondition{
			Expression: "any of ($s*)",
			Type:       "simple",
		},
		IsEnabled:  true,
		SourceFile: "default",
		Priority:   1,
	})

	engine.logger.Infof("传统模式：加载了 %d 个基础Yara规则", len(engine.rules))
}

// loadExternalRules 加载外部Yara规则文件（传统模式）
func (engine *YaraEngine) loadExternalRules() {
	if engine.rulesDir == "" {
		return
	}

	// 遍历规则目录
	err := filepath.Walk(engine.rulesDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() || !strings.HasSuffix(path, ".yar") {
			return nil
		}

		engine.logger.Infof("加载Yara规则文件: %s", path)
		if err := engine.LoadRulesFromFile(path); err != nil {
			engine.logger.Warnf("加载规则文件失败 %s: %v", path, err)
		}

		return nil
	})

	if err != nil {
		engine.logger.Warnf("遍历规则目录失败: %v", err)
	}
}

// LoadRulesFromFile 从文件加载Yara规则
func (engine *YaraEngine) LoadRulesFromFile(filePath string) error {
	content, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("读取规则文件失败: %w", err)
	}

	rules, err := engine.parseYaraRules(string(content))
	if err != nil {
		return fmt.Errorf("解析规则文件失败: %w", err)
	}

	// 设置规则来源
	filename := filepath.Base(filePath)
	for _, rule := range rules {
		rule.SourceFile = filename
		rule.Priority = 999 // 默认优先级
		engine.applyRuleDefaults(rule)
	}

	engine.rules = append(engine.rules, rules...)
	engine.logger.Infof("成功加载 %d 条规则", len(rules))
	return nil
}

// ReloadRules 重新加载规则
func (engine *YaraEngine) ReloadRules() error {
	engine.logger.Info("重新加载Yara规则...")

	if engine.ruleManager != nil {
		// 使用规则管理器重新加载
		if err := engine.ruleManager.ReloadRules(); err != nil {
			return fmt.Errorf("规则管理器重新加载失败: %w", err)
		}

		// 清空现有规则
		engine.rules = make([]*YaraRule, 0)

		// 重新加载规则
		engine.loadRulesFromManager()
	} else {
		// 传统模式重新加载
		engine.rules = make([]*YaraRule, 0)
		engine.loadDefaultRules()
		engine.loadExternalRules()
	}

	// 清空扫描缓存
	engine.clearScanCache()

	engine.logger.Info("Yara规则重新加载完成")
	return nil
}

// clearScanCache 清空扫描缓存
func (engine *YaraEngine) clearScanCache() {
	engine.cacheMutex.Lock()
	defer engine.cacheMutex.Unlock()

	engine.scanCache = make(map[string]*ScanCacheEntry)
	engine.logger.Info("扫描缓存已清空")
}

// parseYaraRules 解析Yara规则内容
func (engine *YaraEngine) parseYaraRules(content string) ([]*YaraRule, error) {
	var rules []*YaraRule
	lines := strings.Split(content, "\n")

	var currentRule *YaraRule
	var inStrings bool
	var inCondition bool
	var inMeta bool
	var currentString *YaraString

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "//") || strings.HasPrefix(line, "/*") {
			continue
		}

		// 检测规则开始
		if strings.HasPrefix(line, "rule ") {
			if currentRule != nil {
				rules = append(rules, currentRule)
			}

			ruleName := strings.TrimSpace(strings.TrimSuffix(strings.TrimPrefix(line, "rule "), " {"))
			currentRule = &YaraRule{
				Name:      ruleName,
				Meta:      make(map[string]string),
				Tags:      make([]string, 0),
				Strings:   make([]*YaraString, 0),
				IsEnabled: true,
				Priority:  999, // 默认优先级
			}
			inStrings = false
			inCondition = false
			inMeta = false
			currentString = nil
			continue
		}

		if currentRule == nil {
			continue
		}

		// 检测元数据段
		if line == "meta:" {
			inMeta = true
			inStrings = false
			inCondition = false
			continue
		}

		// 检测字符串段
		if line == "strings:" {
			inStrings = true
			inMeta = false
			inCondition = false
			continue
		}

		// 检测条件段
		if line == "condition:" {
			inCondition = true
			inMeta = false
			inStrings = false
			continue
		}

		// 检测规则结束
		if line == "}" {
			if currentRule != nil {
				// 添加最后一个字符串（如果有的话）
				if currentString != nil {
					currentRule.Strings = append(currentRule.Strings, currentString)
					currentString = nil
				}
				engine.applyRuleDefaults(currentRule)
				rules = append(rules, currentRule)
				currentRule = nil
			}
			continue
		}

		// 解析元数据
		if inMeta && strings.Contains(line, "=") {
			parts := strings.SplitN(line, "=", 2)
			if len(parts) == 2 {
				key := strings.TrimSpace(parts[0])
				value := strings.Trim(strings.TrimSpace(parts[1]), `"`)

				switch key {
				case "description":
					currentRule.Description = value
				case "severity":
					currentRule.Severity = value
				case "category":
					currentRule.Category = value
				case "author":
					currentRule.Meta["author"] = value
				case "date":
					currentRule.Meta["date"] = value
				case "version":
					currentRule.Meta["version"] = value
				case "priority":
					if priority, err := strconv.Atoi(value); err == nil {
						currentRule.Priority = priority
					}
				case "tags":
					// 处理标签，格式如: tags = "malware,generic,enhanced"
					tags := strings.Split(value, ",")
					for _, tag := range tags {
						tag = strings.TrimSpace(tag)
						if tag != "" {
							currentRule.Tags = append(currentRule.Tags, tag)
						}
					}
				default:
					currentRule.Meta[key] = value
				}
			}
			continue
		}

		// 解析字符串模式
		if inStrings {
			// 检查是否是新的字符串定义
			if strings.Contains(line, "=") && !strings.HasPrefix(line, " ") && !strings.HasPrefix(line, "\t") {
				// 添加前一个字符串（如果有的话）
				if currentString != nil {
					currentRule.Strings = append(currentRule.Strings, currentString)
				}

				parts := strings.SplitN(line, "=", 2)
				if len(parts) == 2 {
					identifier := strings.TrimSpace(parts[0])
					value := strings.TrimSpace(parts[1])

					// 解析字符串类型和修饰符
					stringType := "text"
					var modifiers []string

					// 检查十六进制模式 { ... }
					if strings.Contains(value, "{") && strings.Contains(value, "}") {
						stringType = "hex"
						// 提取十六进制内容
						start := strings.Index(value, "{")
						end := strings.LastIndex(value, "}")
						if start != -1 && end != -1 && end > start {
							value = strings.TrimSpace(value[start+1 : end])
						}
					} else if strings.Contains(value, "/") && strings.Count(value, "/") >= 2 {
						// 检查正则表达式模式 / ... /
						stringType = "regex"
						// 提取正则表达式内容
						start := strings.Index(value, "/")
						end := strings.LastIndex(value, "/")
						if start != -1 && end != -1 && end > start {
							value = strings.TrimSpace(value[start+1 : end])
						}
					} else {
						// 普通文本字符串
						value = strings.Trim(value, `"`)
					}

					// 检查修饰符
					if strings.Contains(line, "nocase") {
						modifiers = append(modifiers, "nocase")
					}
					if strings.Contains(line, "wide") {
						modifiers = append(modifiers, "wide")
					}
					if strings.Contains(line, "ascii") {
						modifiers = append(modifiers, "ascii")
					}
					if strings.Contains(line, "fullword") {
						modifiers = append(modifiers, "fullword")
					}
					if strings.Contains(line, "private") {
						modifiers = append(modifiers, "private")
					}

					currentString = &YaraString{
						Identifier: identifier,
						Value:      value,
						Type:       stringType,
						Modifiers:  modifiers,
					}
				}
			} else if currentString != nil && (strings.HasPrefix(line, " ") || strings.HasPrefix(line, "\t")) {
				// 多行字符串定义，追加到当前字符串
				line = strings.TrimSpace(line)
				if line != "" {
					currentString.Value += " " + line
				}
			}
			continue
		}

		// 解析条件
		if inCondition {
			if currentRule.Condition == nil {
				currentRule.Condition = &YaraCondition{
					Expression: line,
					Type:       engine.parseConditionType(line),
				}
			} else {
				// 多行条件，追加到现有表达式
				currentRule.Condition.Expression += " " + line
			}
		}
	}

	// 添加最后一个规则
	if currentRule != nil {
		// 添加最后一个字符串（如果有的话）
		if currentString != nil {
			currentRule.Strings = append(currentRule.Strings, currentString)
		}
		engine.applyRuleDefaults(currentRule)
		rules = append(rules, currentRule)
	}

	engine.logger.Debugf("解析YARA规则完成，共解析出 %d 条规则", len(rules))
	return rules, nil
}

// parseConditionType 解析条件类型
func (engine *YaraEngine) parseConditionType(condition string) string {
	if strings.Contains(condition, " any of ") || strings.Contains(condition, " any of(") || strings.Contains(condition, "any of ($") {
		return "any_of"
	} else if strings.Contains(condition, " all of ") || strings.Contains(condition, "all of ($") {
		return "all_of"
	} else if strings.Contains(condition, "1 of (") {
		return "one_of"
	} else if strings.Contains(condition, " and ") {
		return "and"
	} else if strings.Contains(condition, " or ") {
		return "or"
	} else if strings.Contains(condition, "not ") {
		return "not"
	}
	return "simple"
}

// ScanFile 扫描文件
func (engine *YaraEngine) ScanFile(filePath string) ([]*YaraMatch, error) {
	// 检查缓存
	if engine.performance != nil && engine.performance.EnableCaching {
		if cached := engine.getCachedResult(filePath); cached != nil {
			engine.logger.Debugf("使用缓存结果: %s", filePath)
			return cached, nil
		}
	}

	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("读取文件失败: %w", err)
	}

	matches, err := engine.ScanData(data, filePath)
	if err != nil {
		return nil, err
	}

	// 缓存结果
	if engine.performance != nil && engine.performance.EnableCaching {
		engine.cacheResult(filePath, data, matches)
	}

	return matches, nil
}

// ScanData 扫描数据
func (engine *YaraEngine) ScanData(data []byte, identifier string) ([]*YaraMatch, error) {
	var matches []*YaraMatch

	engine.setLastScanData(data)

	for _, rule := range engine.rules {
		if !rule.IsEnabled {
			continue
		}

		if match := engine.evaluateRule(rule, data, identifier); match != nil {
			matches = append(matches, match)
		}
	}

	// 按优先级排序匹配结果
	engine.sortMatchesByPriority(matches)

	return matches, nil
}

// sortMatchesByPriority 按优先级排序匹配结果
func (engine *YaraEngine) sortMatchesByPriority(matches []*YaraMatch) {
	sort.Slice(matches, func(i, j int) bool {
		return matches[i].Rule.Priority < matches[j].Rule.Priority
	})
}

// getCachedResult 获取缓存结果
func (engine *YaraEngine) getCachedResult(filePath string) []*YaraMatch {
	engine.cacheMutex.RLock()
	defer engine.cacheMutex.RUnlock()

	// 计算文件哈希作为缓存键
	fileInfo, err := os.Stat(filePath)
	if err != nil {
		return nil
	}

	// 使用文件大小和修改时间作为简单的缓存键
	cacheKey := fmt.Sprintf("%s_%d_%d", filePath, fileInfo.Size(), fileInfo.ModTime().Unix())

	if entry, exists := engine.scanCache[cacheKey]; exists {
		// 检查缓存是否过期（24小时）
		if time.Since(entry.Timestamp) < 24*time.Hour {
			entry.AccessCount++
			return entry.Result
		}
		// 删除过期缓存
		delete(engine.scanCache, cacheKey)
	}

	return nil
}

// cacheResult 缓存扫描结果
func (engine *YaraEngine) cacheResult(filePath string, data []byte, matches []*YaraMatch) {
	engine.cacheMutex.Lock()
	defer engine.cacheMutex.Unlock()

	// 检查缓存大小限制
	if len(engine.scanCache) >= engine.performance.CacheSize {
		engine.evictOldestCache()
	}

	// 计算缓存键
	fileInfo, err := os.Stat(filePath)
	if err != nil {
		return
	}

	cacheKey := fmt.Sprintf("%s_%d_%d", filePath, fileInfo.Size(), fileInfo.ModTime().Unix())

	engine.scanCache[cacheKey] = &ScanCacheEntry{
		Hash:        fmt.Sprintf("%x", md5.Sum(data)),
		Result:      matches,
		Timestamp:   time.Now(),
		AccessCount: 1,
	}
}

// evictOldestCache 淘汰最旧的缓存
func (engine *YaraEngine) evictOldestCache() {
	var oldestKey string
	var oldestTime time.Time

	for key, entry := range engine.scanCache {
		if oldestKey == "" || entry.Timestamp.Before(oldestTime) {
			oldestKey = key
			oldestTime = entry.Timestamp
		}
	}

	if oldestKey != "" {
		delete(engine.scanCache, oldestKey)
		engine.logger.Debugf("淘汰最旧缓存: %s", oldestKey)
	}
}

// evaluateRule 评估单个规则
func (engine *YaraEngine) evaluateRule(rule *YaraRule, data []byte, identifier string) *YaraMatch {
	// 检查PE文件头（如果规则要求）
	if engine.requiresPEHeader(rule) {
		if !engine.hasPEHeader(data) {
			return nil
		}
	}

	// 评估字符串匹配
	stringMatches := engine.evaluateStrings(rule, data)
	if len(stringMatches) == 0 {
		return nil
	}

	// 评估条件
	if !engine.evaluateCondition(rule, stringMatches) {
		return nil
	}

	// 计算匹配分数和置信度
	score := engine.calculateMatchScore(rule, stringMatches)
	confidence := engine.calculateConfidence(rule, stringMatches)

	return &YaraMatch{
		Rule:           rule,
		MatchedAt:      engine.getMatchPositions(stringMatches),
		MatchedStrings: engine.getMatchedStrings(stringMatches),
		Score:          score,
		Confidence:     confidence,
	}
}

// calculateConfidence 计算置信度
func (engine *YaraEngine) calculateConfidence(rule *YaraRule, stringMatches map[string][]int) float64 {
	// 基础置信度：匹配的字符串数量 / 总字符串数量
	baseConfidence := float64(len(stringMatches)) / float64(len(rule.Strings))

	// 根据严重程度调整置信度
	switch rule.Severity {
	case "critical":
		baseConfidence *= 1.2
	case "high":
		baseConfidence *= 1.1
	case "medium":
		baseConfidence *= 1.0
	case "low":
		baseConfidence *= 0.9
	}

	// 根据优先级调整置信度
	if engine.performance != nil {
		baseConfidence *= engine.performance.RulePriorityWeight
	}

	// 限制置信度在0.0-1.0范围内
	if baseConfidence > 1.0 {
		baseConfidence = 1.0
	}
	if baseConfidence < 0.0 {
		baseConfidence = 0.0
	}

	return baseConfidence
}

// requiresPEHeader 检查规则是否需要PE文件头
func (engine *YaraEngine) requiresPEHeader(rule *YaraRule) bool {
	if rule.Condition == nil {
		return false
	}
	return strings.Contains(rule.Condition.Expression, "uint16(0) == 0x5a4d")
}

// hasPEHeader 检查数据是否有PE文件头
func (engine *YaraEngine) hasPEHeader(data []byte) bool {
	if len(data) < 2 {
		return false
	}
	return data[0] == 0x4D && data[1] == 0x5A // MZ
}

// evaluateStrings 评估字符串匹配
func (engine *YaraEngine) evaluateStrings(rule *YaraRule, data []byte) map[string][]int {
	matches := make(map[string][]int)

	for _, str := range rule.Strings {
		positions := engine.findString(str, data)
		if len(positions) > 0 {
			matches[str.Identifier] = positions
		}
	}

	return matches
}

// findString 查找字符串
func (engine *YaraEngine) findString(str *YaraString, data []byte) []int {
	var positions []int

	switch str.Type {
	case "text":
		positions = engine.findTextString(str, data)
	case "hex":
		positions = engine.findHexString(str, data)
	case "regex":
		positions = engine.findRegexString(str, data)
	}

	return positions
}

// findTextString 查找文本字符串
func (engine *YaraEngine) findTextString(str *YaraString, data []byte) []int {
	var positions []int
	searchValue := str.Value

	// 处理nocase修饰符
	if engine.hasModifier(str, "nocase") {
		searchValue = strings.ToLower(searchValue)
	}

	// 处理wide修饰符
	if engine.hasModifier(str, "wide") {
		positions = engine.findWideString(searchValue, data)
	} else {
		positions = engine.findSimpleString(searchValue, data)
	}

	return positions
}

// findHexString 查找十六进制字符串
func (engine *YaraEngine) findHexString(str *YaraString, data []byte) []int {
	var positions []int
	hexBytes := engine.parseHexString(str.Value)

	if len(hexBytes) == 0 {
		return positions
	}

	for i := 0; i <= len(data)-len(hexBytes); i++ {
		if engine.bytesEqual(data[i:i+len(hexBytes)], hexBytes) {
			positions = append(positions, i)
		}
	}

	return positions
}

// findRegexString 查找正则表达式字符串
func (engine *YaraEngine) findRegexString(str *YaraString, data []byte) []int {
	var positions []int

	regex, err := regexp.Compile(str.Value)
	if err != nil {
		engine.logger.Warnf("正则表达式编译失败: %v", err)
		return positions
	}

	// 查找所有匹配
	allMatches := regex.FindAllIndex(data, -1)
	for _, match := range allMatches {
		positions = append(positions, match[0])
	}

	return positions
}

// findWideString 查找宽字符串
func (engine *YaraEngine) findWideString(searchValue string, data []byte) []int {
	var positions []int

	// 宽字符串是每个字符后跟一个空字节
	for i := 0; i < len(data)-len(searchValue)*2; i++ {
		matched := true
		for j := 0; j < len(searchValue); j++ {
			if data[i+j*2] != searchValue[j] || data[i+j*2+1] != 0 {
				matched = false
				break
			}
		}
		if matched {
			positions = append(positions, i)
		}
	}

	return positions
}

// findSimpleString 查找简单字符串
func (engine *YaraEngine) findSimpleString(searchValue string, data []byte) []int {
	var positions []int
	searchBytes := []byte(searchValue)

	for i := 0; i <= len(data)-len(searchBytes); i++ {
		if engine.bytesEqual(data[i:i+len(searchBytes)], searchBytes) {
			positions = append(positions, i)
		}
	}

	return positions
}

// parseHexString 解析十六进制字符串
func (engine *YaraEngine) parseHexString(hexStr string) []byte {
	var result []byte
	hexStr = strings.ReplaceAll(hexStr, " ", "")

	for i := 0; i < len(hexStr); i += 2 {
		if i+1 < len(hexStr) {
			if b, err := strconv.ParseUint(hexStr[i:i+2], 16, 8); err == nil {
				result = append(result, byte(b))
			}
		}
	}

	return result
}

// hasModifier 检查字符串是否有指定修饰符
func (engine *YaraEngine) hasModifier(str *YaraString, modifier string) bool {
	for _, mod := range str.Modifiers {
		if mod == modifier {
			return true
		}
	}
	return false
}

// bytesEqual 比较字节数组
func (engine *YaraEngine) bytesEqual(a, b []byte) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

// evaluateCondition 评估条件
func (engine *YaraEngine) evaluateCondition(rule *YaraRule, stringMatches map[string][]int) bool {
	if rule.Condition == nil {
		return len(stringMatches) > 0
	}

	expr := rule.Condition.Expression
	if expr == "" {
		return len(stringMatches) > 0
	}

	// 处理 any/all/1 of ($prefix*) 分组语义
	groupPattern := regexp.MustCompile(`(?i)(any|all|1)\s+of\s*\(\s*\$([a-zA-Z0-9_]+)\*\s*\)`) // 捕获 like: any of ($sys_app*)
	replaced := expr
	matches := groupPattern.FindAllStringSubmatch(expr, -1)
	computed := map[string]bool{}
	placeholderIdx := 0
	for _, m := range matches {
		if len(m) < 3 {
			continue
		}
		quantifier := strings.ToLower(m[1])
		prefix := "$" + m[2]
		// 统计该前缀的匹配数与总数
		matchedCount := 0
		total := 0
		for _, ys := range rule.Strings {
			if strings.HasPrefix(ys.Identifier, prefix) {
				total++
				if _, ok := stringMatches[ys.Identifier]; ok {
					matchedCount++
				}
			}
		}

		okVal := false
		switch quantifier {
		case "any", "1":
			okVal = matchedCount >= 1
		case "all":
			if total == 0 {
				okVal = false
			} else {
				okVal = matchedCount == total
			}
		}

		placeholder := fmt.Sprintf("__GROUP_%d__", placeholderIdx)
		placeholderIdx++
		computed[placeholder] = okVal
		replaced = strings.ReplaceAll(replaced, m[0], placeholder)
	}

	// 处理 pe.imports("dll","func") 近似支持
	importsPattern := regexp.MustCompile(`(?i)pe\.imports\(\s*"([^"]+)"\s*,\s*"([^"]+)"\s*\)`) // 捕获 DLL 与函数
	importsMatches := importsPattern.FindAllStringSubmatch(replaced, -1)
	for _, m := range importsMatches {
		if len(m) < 3 {
			continue
		}
		dll := m[1]
		fn := m[2]
		okVal := engine.approxHasImport(engine.lastScanData, dll, fn)
		if okVal {
			replaced = strings.ReplaceAll(replaced, m[0], "true")
		} else {
			replaced = strings.ReplaceAll(replaced, m[0], "false")
		}
	}

	// 将所有已匹配的标识符替换为 true，未匹配的（简单变量形式）替换为 false
	identifierPattern := regexp.MustCompile(`\$[a-zA-Z0-9_]+`)
	replaced = identifierPattern.ReplaceAllStringFunc(replaced, func(tok string) string {
		if strings.HasPrefix(tok, "__GROUP_") { // 跳过占位符
			return tok
		}
		if _, ok := computed[tok]; ok {
			if computed[tok] {
				return "true"
			}
			return "false"
		}
		if _, ok := stringMatches[tok]; ok {
			return "true"
		}
		return "false"
	})

	// 将占位符替换为 true/false
	for k, v := range computed {
		if v {
			replaced = strings.ReplaceAll(replaced, k, "true")
		} else {
			replaced = strings.ReplaceAll(replaced, k, "false")
		}
	}

	// 简单的布尔表达式求值（支持 and / or / not）
	result := engine.evalBooleanExpression(replaced)
	return result
}

// evalBooleanExpression 评估由 true/false 与 and/or/not 组成的表达式（忽略大小写）
func (engine *YaraEngine) evalBooleanExpression(expr string) bool {
	// 归一化
	s := strings.ToLower(expr)
	// 移除多余括号（不实现括号优先，仅做简单清洗）
	s = strings.ReplaceAll(s, "(", " ")
	s = strings.ReplaceAll(s, ")", " ")
	// 标准化空白
	fields := strings.Fields(s)
	if len(fields) == 0 {
		return false
	}

	// 把 true/false 与 and/or/not 顺序计算：not > and > or
	// 先处理 not
	var tmp []string
	for i := 0; i < len(fields); i++ {
		if fields[i] == "not" && i+1 < len(fields) {
			val := fields[i+1] == "true"
			if val {
				tmp = append(tmp, "false")
			} else {
				tmp = append(tmp, "true")
			}
			i++
		} else {
			tmp = append(tmp, fields[i])
		}
	}

	// 处理 and
	var tmp2 []string
	i := 0
	for i < len(tmp) {
		if i+2 < len(tmp) && tmp[i+1] == "and" {
			left := tmp[i] == "true"
			right := tmp[i+2] == "true"
			val := left && right
			if val {
				tmp2 = append(tmp2, "true")
			} else {
				tmp2 = append(tmp2, "false")
			}
			i += 3
		} else {
			tmp2 = append(tmp2, tmp[i])
			i++
		}
	}

	// 处理 or
	res := false
	opOrSeen := false
	for j := 0; j < len(tmp2); j++ {
		token := tmp2[j]
		if token == "or" {
			opOrSeen = true
			continue
		}
		val := token == "true"
		if !opOrSeen {
			// 第一个值
			res = val
			opOrSeen = true // 标记已设置初值
		} else {
			res = res || val
		}
	}
	return res
}

// applyRuleDefaults 为常见规则设置默认描述/严重性/类别
func (engine *YaraEngine) applyRuleDefaults(rule *YaraRule) {
	if rule == nil {
		return
	}
	if rule.Description != "" && rule.Severity != "" && rule.Category != "" {
		return
	}

	name := strings.ToLower(rule.Name)
	switch name {
	case "disable_antivirus":
		if rule.Description == "" {
			rule.Description = "检测禁用或规避本地安全产品行为"
		}
		if rule.Severity == "" {
			rule.Severity = "high"
		}
		if rule.Category == "" {
			rule.Category = "evasion"
		}
	case "anti_dbgtools":
		if rule.Description == "" {
			rule.Description = "检测反调试与调试工具自检行为"
		}
		if rule.Severity == "" {
			rule.Severity = "medium"
		}
		if rule.Category == "" {
			rule.Category = "anti_debug"
		}
	case "sandbox_detection":
		if rule.Description == "" {
			rule.Description = "检测沙箱/虚拟化环境识别行为"
		}
		if rule.Severity == "" {
			rule.Severity = "medium"
		}
		if rule.Category == "" {
			rule.Category = "anti_vm"
		}
	case "inject_thread":
		if rule.Description == "" {
			rule.Description = "检测远程线程注入相关API使用"
		}
		if rule.Severity == "" {
			rule.Severity = "high"
		}
		if rule.Category == "" {
			rule.Category = "injection"
		}
	case "create_process":
		if rule.Description == "" {
			rule.Description = "检测可疑进程创建行为/相关API"
		}
		if rule.Severity == "" {
			rule.Severity = "medium"
		}
		if rule.Category == "" {
			rule.Category = "process"
		}
	case "persistence":
		if rule.Description == "" {
			rule.Description = "检测常见持久化注册表位置与机制"
		}
		if rule.Severity == "" {
			rule.Severity = "high"
		}
		if rule.Category == "" {
			rule.Category = "persistence"
		}
	case "hijack_network":
		if rule.Description == "" {
			rule.Description = "检测网络劫持/代理/hosts 篡改等"
		}
		if rule.Severity == "" {
			rule.Severity = "medium"
		}
		if rule.Category == "" {
			rule.Category = "network"
		}
	case "create_service":
		if rule.Description == "" {
			rule.Description = "检测恶意服务创建与控制相关API"
		}
		if rule.Severity == "" {
			rule.Severity = "high"
		}
		if rule.Category == "" {
			rule.Category = "service"
		}
	case "create_com_service":
		if rule.Description == "" {
			rule.Description = "检测COM注册/加载典型入口"
		}
		if rule.Severity == "" {
			rule.Severity = "medium"
		}
		if rule.Category == "" {
			rule.Category = "com"
		}
	case "network_udp_sock":
		if rule.Description == "" {
			rule.Description = "检测UDP通信相关API/类库使用"
		}
		if rule.Severity == "" {
			rule.Severity = "low"
		}
		if rule.Category == "" {
			rule.Category = "network"
		}
	case "network_tcp_listen":
		if rule.Description == "" {
			rule.Description = "检测TCP监听/接受连接相关API使用"
		}
		if rule.Severity == "" {
			rule.Severity = "low"
		}
		if rule.Category == "" {
			rule.Category = "network"
		}
	default:
		// 其他规则使用默认中等严重性
		if rule.Severity == "" {
			rule.Severity = "medium"
		}
		if rule.Category == "" {
			rule.Category = "unknown"
		}
		if rule.Description == "" {
			rule.Description = "检测到可疑特征"
		}
	}
}

// GetRules 获取所有规则
func (engine *YaraEngine) GetRules() []*YaraRule {
	return engine.rules
}

// GetRuleByName 根据名称获取规则
func (engine *YaraEngine) GetRuleByName(name string) *YaraRule {
	for _, rule := range engine.rules {
		if rule.Name == name {
			return rule
		}
	}
	return nil
}

// GetRulesCount 获取规则数量
func (engine *YaraEngine) GetRulesCount() int {
	return len(engine.rules)
}

// EnableRule 启用规则
func (engine *YaraEngine) EnableRule(name string) error {
	rule := engine.GetRuleByName(name)
	if rule == nil {
		return fmt.Errorf("规则不存在: %s", name)
	}
	rule.IsEnabled = true
	return nil
}

// DisableRule 禁用规则
func (engine *YaraEngine) DisableRule(name string) error {
	rule := engine.GetRuleByName(name)
	if rule == nil {
		return fmt.Errorf("规则不存在: %s", name)
	}
	rule.IsEnabled = false
	return nil
}

// GetRuleManager 获取规则管理器
func (engine *YaraEngine) GetRuleManager() *RuleManager {
	return engine.ruleManager
}

// GetPerformanceStats 获取性能统计
func (engine *YaraEngine) GetPerformanceStats() map[string]interface{} {
	stats := map[string]interface{}{
		"total_rules":    len(engine.rules),
		"enabled_rules":  0,
		"cache_size":     0,
		"cache_hit_rate": 0.0,
	}

	// 统计启用的规则
	for _, rule := range engine.rules {
		if rule.IsEnabled {
			stats["enabled_rules"] = stats["enabled_rules"].(int) + 1
		}
	}

	// 统计缓存信息
	if engine.performance != nil && engine.performance.EnableCaching {
		stats["cache_size"] = len(engine.scanCache)
		stats["cache_hit_rate"] = engine.calculateCacheHitRate()
	}

	return stats
}

// calculateCacheHitRate 计算缓存命中率
func (engine *YaraEngine) calculateCacheHitRate() float64 {
	engine.cacheMutex.RLock()
	defer engine.cacheMutex.RUnlock()

	totalAccess := 0
	for _, entry := range engine.scanCache {
		totalAccess += entry.AccessCount
	}

	if totalAccess == 0 {
		return 0.0
	}

	// 简单的命中率计算：缓存条目数 / 总访问数
	return float64(len(engine.scanCache)) / float64(totalAccess)
}

// approxHasImport 近似判断数据中是否包含 DLL/函数名（ASCII 或宽字节）
func (engine *YaraEngine) approxHasImport(dataCtx []byte, dll, fn string) bool {
	if len(dataCtx) == 0 {
		return false
	}
	data := dataCtx
	low := strings.ToLower(string(data))
	if strings.Contains(low, strings.ToLower(fn)) {
		return true
	}
	// 宽字节搜索函数名
	wide := make([]byte, 0, len(fn)*2)
	for i := 0; i < len(fn); i++ {
		wide = append(wide, fn[i], 0)
	}
	return bytesIndex(data, wide) >= 0
}

// bytesIndex 在 a 中查找子切片 b 的起始索引，找不到返回 -1
func bytesIndex(a, b []byte) int {
	if len(b) == 0 {
		return 0
	}
	for i := 0; i+len(b) <= len(a); i++ {
		matched := true
		for j := 0; j < len(b); j++ {
			if a[i+j] != b[j] {
				matched = false
				break
			}
		}
		if matched {
			return i
		}
	}
	return -1
}

// 在 ScanData/evaluateRule 期间保存最近一次被评估的数据供 approxHasImport 使用
var _lastDataMu sync.Mutex

// lastScanData 保存最近一次扫描数据的副本（谨慎使用，避免大量分配）
func (engine *YaraEngine) setLastScanData(data []byte) {
	_lastDataMu.Lock()
	defer _lastDataMu.Unlock()
	engine.lastScanData = data
}

// calculateMatchScore 计算匹配分数
func (engine *YaraEngine) calculateMatchScore(rule *YaraRule, stringMatches map[string][]int) int {
	score := len(stringMatches) * 10

	// 根据严重程度调整分数
	switch strings.ToLower(rule.Severity) {
	case "critical":
		score *= 3
	case "high":
		score *= 2
	case "medium":
		// 保持不变
	case "low":
		score = score / 2
	}

	// 根据优先级调整分数
	if engine.performance != nil {
		score = int(float64(score) * engine.performance.RulePriorityWeight)
	}
	return score
}

// getMatchPositions 获取匹配位置
func (engine *YaraEngine) getMatchPositions(stringMatches map[string][]int) []int64 {
	var positions []int64
	for _, pos := range stringMatches {
		for _, p := range pos {
			positions = append(positions, int64(p))
		}
	}
	return positions
}

// getMatchedStrings 获取匹配的字符串
func (engine *YaraEngine) getMatchedStrings(stringMatches map[string][]int) []string {
	var result []string
	for identifier := range stringMatches {
		result = append(result, identifier)
	}
	return result
}
