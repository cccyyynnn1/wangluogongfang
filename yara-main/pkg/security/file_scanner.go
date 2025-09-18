package security

import (
	"context"
	"fmt"
	"io/ioutil"
	"regexp"
	"time"

	"yara-security-service/internal/models"

	"github.com/sirupsen/logrus"
)

// FileScanner 文件扫描器
type FileScanner struct {
	logger *logrus.Logger
	
	// 恶意字符模式
	maliciousCharPatterns []*regexp.Regexp
	
	// 恶意数据结构模式
	maliciousDataPatterns []*regexp.Regexp
	
	// 文件类型检测器
	fileTypeDetector *FileTypeDetector
}

// FileTypeDetector 文件类型检测器
type FileTypeDetector struct {
	// 文件头魔数
	magicNumbers map[string][]byte
}

// NewFileScanner 创建文件扫描器
func NewFileScanner(logger *logrus.Logger) *FileScanner {
	scanner := &FileScanner{
		logger: logger,
	}
	
	// 初始化恶意字符模式
	scanner.initMaliciousCharPatterns()
	
	// 初始化恶意数据结构模式
	scanner.initMaliciousDataPatterns()
	
	// 初始化文件类型检测器
	scanner.fileTypeDetector = NewFileTypeDetector()
	
	return scanner
}

// NewFileTypeDetector 创建文件类型检测器
func NewFileTypeDetector() *FileTypeDetector {
	detector := &FileTypeDetector{
		magicNumbers: make(map[string][]byte),
	}
	
	// 初始化文件头魔数
	detector.initMagicNumbers()
	
	return detector
}

// initMagicNumbers 初始化文件头魔数
func (fd *FileTypeDetector) initMagicNumbers() {
	fd.magicNumbers = map[string][]byte{
		"PE":     {0x4D, 0x5A}, // MZ
		"ELF":    {0x7F, 0x45, 0x4C, 0x46}, // ELF
		"ZIP":    {0x50, 0x4B, 0x03, 0x04}, // PK
		"RAR":    {0x52, 0x61, 0x72, 0x21}, // Rar!
		"PDF":    {0x25, 0x50, 0x44, 0x46}, // %PDF
		"PNG":    {0x89, 0x50, 0x4E, 0x47}, // PNG
		"JPEG":   {0xFF, 0xD8, 0xFF}, // JPEG
		"GIF":    {0x47, 0x49, 0x46}, // GIF
		"BMP":    {0x42, 0x4D}, // BM
		"TIFF":   {0x49, 0x49, 0x2A, 0x00}, // II*
		"ICO":    {0x00, 0x00, 0x01, 0x00}, // ICO
		"EXE":    {0x4D, 0x5A}, // MZ
		"DLL":    {0x4D, 0x5A}, // MZ
		"BAT":    {0x40, 0x65, 0x63, 0x68, 0x6F}, // @echo
		"CMD":    {0x40, 0x65, 0x63, 0x68, 0x6F}, // @echo
		"PS1":    {0x23, 0x21, 0x2F, 0x75, 0x73, 0x72, 0x2F, 0x62, 0x69, 0x6E, 0x2F, 0x70, 0x77, 0x73, 0x68}, // #!/usr/bin/pwsh
		"VBS":    {0x27, 0x20, 0x56, 0x42, 0x53}, // ' VBS
		"JS":     {0x2F, 0x2F, 0x20, 0x4A, 0x61, 0x76, 0x61, 0x53, 0x63, 0x72, 0x69, 0x70, 0x74}, // // JavaScript
		"PY":     {0x23, 0x21, 0x2F, 0x75, 0x73, 0x72, 0x2F, 0x62, 0x69, 0x6E, 0x2F, 0x70, 0x79, 0x74, 0x68, 0x6F, 0x6E}, // #!/usr/bin/python
		"SH":     {0x23, 0x21, 0x2F, 0x62, 0x69, 0x6E, 0x2F, 0x73, 0x68}, // #!/bin/sh
	}
}

// initMaliciousCharPatterns 初始化恶意字符模式
func (fs *FileScanner) initMaliciousCharPatterns() {
	patterns := []string{
		// 恶意命令模式
		`(?i)(cmd\.exe|powershell\.exe|wscript\.exe|cscript\.exe)`,
		`(?i)(net\s+user|net\s+group|net\s+localgroup)`,
		`(?i)(reg\s+add|reg\s+delete|reg\s+export|reg\s+import)`,
		`(?i)(schtasks|at\s+\\|sc\s+create|sc\s+start)`,
		`(?i)(format\s+[a-z]:|del\s+/s|rd\s+/s)`,
		`(?i)(shutdown|restart|logoff|taskkill)`,
		
		// 恶意URL模式
		`(?i)(http://|https://|ftp://|file://)`,
		`(?i)(bit\.ly|goo\.gl|tinyurl\.com)`,
		`(?i)(malware|virus|trojan|backdoor)`,
		
		// 恶意IP地址
		`\b(?:[0-9]{1,3}\.){3}[0-9]{1,3}\b`,
		
		// 恶意邮箱
		`(?i)([a-z0-9._%+-]+@[a-z0-9.-]+\.[a-z]{2,})`,
		
		// 恶意文件路径
		`(?i)(c:\\windows\\system32|c:\\windows\\syswow64)`,
		`(?i)(%temp%|%tmp%|%systemroot%)`,
		
		// 恶意进程名
		`(?i)(explorer\.exe|svchost\.exe|lsass\.exe|csrss\.exe)`,
		
		// 恶意注册表路径
		`(?i)(HKEY_LOCAL_MACHINE|HKEY_CURRENT_USER)`,
		`(?i)(SOFTWARE\\Microsoft\\Windows\\CurrentVersion)`,
		`(?i)(SYSTEM\\CurrentControlSet\\Services)`,
		
		// 恶意API调用
		`(?i)(CreateProcess|CreateRemoteThread|VirtualAllocEx)`,
		`(?i)(WriteProcessMemory|ReadProcessMemory|OpenProcess)`,
		`(?i)(SetWindowsHookEx|SetThreadContext|SuspendThread)`,
		
		// 恶意字符串
		`(?i)(password|admin|root|system|config)`,
		`(?i)(hack|exploit|attack|breach|compromise)`,
		`(?i)(keylogger|spyware|backdoor|trojan)`,
	}
	
	for _, pattern := range patterns {
		if regex, err := regexp.Compile(pattern); err == nil {
			fs.maliciousCharPatterns = append(fs.maliciousCharPatterns, regex)
		}
	}
}

// initMaliciousDataPatterns 初始化恶意数据结构模式
func (fs *FileScanner) initMaliciousDataPatterns() {
	patterns := []string{
		// PE文件头
		`MZ.{0,100}PE`,
		
		// 可执行文件特征
		`\x4D\x5A.{0,100}\x50\x45`, // MZ...PE
		
		// 脚本文件特征
		`#!.*(python|perl|bash|sh)`,
		`@echo\s+off`,
		`#!/usr/bin/env`,
		
		// 压缩文件特征
		`PK\x03\x04`, // ZIP
		`Rar!\x1A\x07`,
		
		// 文档文件特征
		`%PDF`,
		`\x89PNG`,
		`GIF8`,
		
		// 恶意数据结构
		`\x90{4,}`, // NOP sled
		`\xCC{2,}`, // INT3
		`\xEB\xFE`, // JMP $-2
		
		// 加密/编码数据
		`[A-Za-z0-9+/]{50,}={0,2}`, // Base64
		`[0-9A-Fa-f]{32,}`, // Hex
		
		// 网络相关
		`(?:[0-9]{1,3}\.){3}[0-9]{1,3}:[0-9]{1,5}`, // IP:Port
		`[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}:[0-9]{1,5}`, // Domain:Port
	}
	
	for _, pattern := range patterns {
		if regex, err := regexp.Compile(pattern); err == nil {
			fs.maliciousDataPatterns = append(fs.maliciousDataPatterns, regex)
		}
	}
}

// ScanFileContent 扫描文件内容
func (fs *FileScanner) ScanFileContent(ctx context.Context, filePath string) (*models.FileScanResult, error) {
	startTime := time.Now()
	
	// 读取文件内容
	content, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("读取文件失败: %w", err)
	}
	
	// 检测文件类型
	fileType := fs.fileTypeDetector.DetectFileType(content)
	
	// 检测恶意字符
	charThreats := fs.detectMaliciousCharacters(string(content))
	
	// 检测恶意数据结构
	dataThreats := fs.detectMaliciousDataStructures(content)
	
	// 合并威胁
	allThreats := append(charThreats, dataThreats...)
	
	// 构建扫描结果
	result := &models.FileScanResult{
		FilePath:     filePath,
		FileType:     fileType,
		IsSuspicious: len(allThreats) > 0,
		Threats:      allThreats,
		ScanTime:     time.Now(),
		ScanDuration: time.Since(startTime),
		ThreatLevel:  fs.calculateThreatLevel(allThreats),
		Category:     "content_scan",
	}
	
	fs.logger.Infof("文件内容扫描完成: %s, 类型: %s, 可疑: %v, 威胁数: %d", 
		filePath, fileType, result.IsSuspicious, len(allThreats))
	
	return result, nil
}

// DetectFileType 检测文件类型
func (fd *FileTypeDetector) DetectFileType(content []byte) string {
	if len(content) == 0 {
		return "unknown"
	}
	
	// 检查文件头魔数
	for fileType, magic := range fd.magicNumbers {
		if len(content) >= len(magic) && bytesEqual(content[:len(magic)], magic) {
			return fileType
		}
	}
	
	// 检查文本文件
	if isTextFile(content) {
		return "text"
	}
	
	// 检查脚本文件
	if isScriptFile(content) {
		return "script"
	}
	
	return "unknown"
}

// bytesEqual 比较字节数组
func bytesEqual(a, b []byte) bool {
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

// isTextFile 检查是否为文本文件
func isTextFile(content []byte) bool {
	if len(content) == 0 {
		return false
	}
	
	// 检查是否包含大量可打印字符
	printableCount := 0
	for _, b := range content {
		if b >= 32 && b <= 126 || b == 9 || b == 10 || b == 13 {
			printableCount++
		}
	}
	
	// 如果可打印字符比例超过90%，认为是文本文件
	return float64(printableCount)/float64(len(content)) > 0.9
}

// isScriptFile 检查是否为脚本文件
func isScriptFile(content []byte) bool {
	if len(content) == 0 {
		return false
	}
	
	contentStr := string(content)
	
	// 检查脚本文件特征
	scriptPatterns := []string{
		`#!.*(python|perl|bash|sh|php|ruby)`,
		`@echo\s+off`,
		`#!/usr/bin/env`,
		`function\s+\w+`,
		`def\s+\w+`,
		`class\s+\w+`,
		`import\s+\w+`,
		`require\s+`,
		`include\s+`,
	}
	
	for _, pattern := range scriptPatterns {
		if matched, _ := regexp.MatchString(pattern, contentStr); matched {
			return true
		}
	}
	
	return false
}

// detectMaliciousCharacters 检测恶意字符
func (fs *FileScanner) detectMaliciousCharacters(content string) []models.ThreatInfo {
	var threats []models.ThreatInfo
	
	for i, pattern := range fs.maliciousCharPatterns {
		if pattern.MatchString(content) {
			threats = append(threats, models.ThreatInfo{
				RuleName:    fmt.Sprintf("Malicious_Char_Pattern_%d", i+1),
				Description: "检测到恶意字符模式",
				Severity:    "medium",
				Category:    "malicious_char",
				Tags:        "character,malicious",
			})
		}
	}
	
	return threats
}

// detectMaliciousDataStructures 检测恶意数据结构
func (fs *FileScanner) detectMaliciousDataStructures(content []byte) []models.ThreatInfo {
	var threats []models.ThreatInfo
	
	for i, pattern := range fs.maliciousDataPatterns {
		if pattern.Match(content) {
			threats = append(threats, models.ThreatInfo{
				RuleName:    fmt.Sprintf("Malicious_Data_Pattern_%d", i+1),
				Description: "检测到恶意数据结构",
				Severity:    "high",
				Category:    "malicious_data",
				Tags:        "data_structure,malicious",
			})
		}
	}
	
	return threats
}

// calculateThreatLevel 计算威胁等级
func (fs *FileScanner) calculateThreatLevel(threats []models.ThreatInfo) string {
	if len(threats) == 0 {
		return "clean"
	}
	
	highCount := 0
	mediumCount := 0
	lowCount := 0
	
	for _, threat := range threats {
		switch threat.Severity {
		case "high":
			highCount++
		case "medium":
			mediumCount++
		case "low":
			lowCount++
		}
	}
	
	if highCount > 0 {
		return "high"
	} else if mediumCount > 0 {
		return "medium"
	} else {
		return "low"
	}
} 