package security

import (
	"context"
	"fmt"
	"net/url"
	"regexp"
	"time"

	"yara-security-service/internal/models"

	"github.com/sirupsen/logrus"
)

// NetworkScanner 网络流扫描器
type NetworkScanner struct {
	logger *logrus.Logger

	// URL检测模式
	maliciousURLPatterns []*regexp.Regexp

	// Shellcode检测模式
	shellcodePatterns []*regexp.Regexp
}

// NewNetworkScanner 创建网络流扫描器
func NewNetworkScanner(logger *logrus.Logger) *NetworkScanner {
	scanner := &NetworkScanner{
		logger: logger,
	}

	// 初始化检测模式
	scanner.initMaliciousURLPatterns()
	scanner.initShellcodePatterns()

	return scanner
}

// initMaliciousURLPatterns 初始化恶意URL模式
func (ns *NetworkScanner) initMaliciousURLPatterns() {
	patterns := []string{
		// 恶意域名
		`(?i)(malware|virus|trojan|backdoor|spyware|keylogger)`,
		`(?i)(bit\.ly|goo\.gl|tinyurl\.com)`,

		// 恶意IP地址
		`\b(?:[0-9]{1,3}\.){3}[0-9]{1,3}\b`,

		// 恶意路径
		`(?i)(/admin|/backdoor|/shell|/cmd|/exec)`,
		`(?i)(/download|/upload|/inject|/payload)`,
	}

	for _, pattern := range patterns {
		if regex, err := regexp.Compile(pattern); err == nil {
			ns.maliciousURLPatterns = append(ns.maliciousURLPatterns, regex)
		}
	}
}

// initShellcodePatterns 初始化Shellcode模式
func (ns *NetworkScanner) initShellcodePatterns() {
	patterns := []string{
		// 常见Shellcode特征
		`\x90{4,}`, // NOP sled
		`\xCC{2,}`, // INT3
		`\xEB\xFE`, // JMP $-2
		`\x31\xC0`, // XOR EAX, EAX
		`\x50`,     // PUSH EAX
		`\x68`,     // PUSH immediate
		`\x89\xE3`, // MOV EBX, ESP
		`\x89\xE1`, // MOV ECX, ESP
		`\xB0\x0B`, // MOV AL, 11
		`\xCD\x80`, // INT 80h (Linux syscall)

		// 字符串Shellcode
		`(?i)(cmd\.exe|powershell\.exe|wscript\.exe)`,
		`(?i)(/bin/sh|/bin/bash|/bin/dash)`,
		`(?i)(netcat|nc|telnet|ssh)`,
	}

	for _, pattern := range patterns {
		if regex, err := regexp.Compile(pattern); err == nil {
			ns.shellcodePatterns = append(ns.shellcodePatterns, regex)
		}
	}
}

// ScanURL 扫描URL
func (ns *NetworkScanner) ScanURL(ctx context.Context, urlStr string) (*models.NetworkScanResult, error) {
	startTime := time.Now()

	// 解析URL
	parsedURL, err := url.Parse(urlStr)
	if err != nil {
		return nil, fmt.Errorf("URL解析失败: %w", err)
	}

	// 检测恶意URL
	threats := ns.detectMaliciousURL(urlStr)

	// 构建扫描结果
	result := &models.NetworkScanResult{
		URL:          urlStr,
		IsSuspicious: len(threats) > 0,
		Threats:      threats,
		ScanTime:     time.Now(),
		ScanDuration: time.Since(startTime),
		ThreatLevel:  ns.calculateThreatLevel(threats),
		Category:     "url_scan",
		Metadata:     ns.getURLMetadata(parsedURL),
	}

	ns.logger.Infof("URL扫描完成: %s, 可疑: %v, 威胁数: %d",
		urlStr, result.IsSuspicious, len(threats))

	return result, nil
}

// ScanShellcode 扫描Shellcode
func (ns *NetworkScanner) ScanShellcode(ctx context.Context, data []byte) (*models.NetworkScanResult, error) {
	startTime := time.Now()

	// 检测Shellcode
	threats := ns.detectShellcode(data)

	// 构建扫描结果
	result := &models.NetworkScanResult{
		URL:          fmt.Sprintf("data:%d_bytes", len(data)),
		IsSuspicious: len(threats) > 0,
		Threats:      threats,
		ScanTime:     time.Now(),
		ScanDuration: time.Since(startTime),
		ThreatLevel:  ns.calculateThreatLevel(threats),
		Category:     "shellcode_scan",
		Metadata:     ns.getShellcodeMetadata(data),
	}

	ns.logger.Infof("Shellcode扫描完成: %d字节, 可疑: %v, 威胁数: %d",
		len(data), result.IsSuspicious, len(threats))

	return result, nil
}

// detectMaliciousURL 检测恶意URL
func (ns *NetworkScanner) detectMaliciousURL(urlStr string) []models.ThreatInfo {
	var threats []models.ThreatInfo

	for i, pattern := range ns.maliciousURLPatterns {
		if pattern.MatchString(urlStr) {
			threats = append(threats, models.ThreatInfo{
				RuleName:    fmt.Sprintf("Malicious_URL_Pattern_%d", i+1),
				Description: "检测到恶意URL模式",
				Severity:    "medium",
				Category:    "malicious_url",
				Tags:        "url,malicious",
			})
		}
	}

	return threats
}

// detectShellcode 检测Shellcode
func (ns *NetworkScanner) detectShellcode(data []byte) []models.ThreatInfo {
	var threats []models.ThreatInfo

	// 检测二进制Shellcode
	for i, pattern := range ns.shellcodePatterns {
		if pattern.Match(data) {
			threats = append(threats, models.ThreatInfo{
				RuleName:    fmt.Sprintf("Shellcode_Pattern_%d", i+1),
				Description: "检测到Shellcode模式",
				Severity:    "critical",
				Category:    "shellcode",
				Tags:        "shellcode,malicious",
			})
		}
	}

	return threats
}

// calculateThreatLevel 计算威胁等级
func (ns *NetworkScanner) calculateThreatLevel(threats []models.ThreatInfo) string {
	if len(threats) == 0 {
		return "clean"
	}

	highCount := 0
	mediumCount := 0
	lowCount := 0

	for _, threat := range threats {
		switch threat.Severity {
		case "critical", "high":
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

// getURLMetadata 获取URL元数据
func (ns *NetworkScanner) getURLMetadata(parsedURL *url.URL) map[string]string {
	return map[string]string{
		"scheme":   parsedURL.Scheme,
		"host":     parsedURL.Hostname(),
		"port":     parsedURL.Port(),
		"path":     parsedURL.Path,
		"query":    parsedURL.RawQuery,
		"fragment": parsedURL.Fragment,
	}
}

// getShellcodeMetadata 获取Shellcode元数据
func (ns *NetworkScanner) getShellcodeMetadata(data []byte) map[string]string {
	return map[string]string{
		"size":            fmt.Sprintf("%d", len(data)),
		"has_nulls":       fmt.Sprintf("%v", ns.containsBytes(data, []byte{0x00})),
		"printable_ratio": fmt.Sprintf("%.2f", ns.calculatePrintableRatio(data)),
	}
}

// containsBytes 检查字节数组是否包含子数组
func (ns *NetworkScanner) containsBytes(data, sub []byte) bool {
	if len(sub) > len(data) {
		return false
	}

	for i := 0; i <= len(data)-len(sub); i++ {
		if ns.bytesEqual(data[i:i+len(sub)], sub) {
			return true
		}
	}

	return false
}

// bytesEqual 比较字节数组
func (ns *NetworkScanner) bytesEqual(a, b []byte) bool {
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

// calculatePrintableRatio 计算可打印字符比例
func (ns *NetworkScanner) calculatePrintableRatio(data []byte) float64 {
	if len(data) == 0 {
		return 0
	}

	printableCount := 0
	for _, b := range data {
		if b >= 32 && b <= 126 {
			printableCount++
		}
	}

	return float64(printableCount) / float64(len(data))
}
