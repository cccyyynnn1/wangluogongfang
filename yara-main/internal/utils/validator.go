package utils

import (
	"fmt"
	"net"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
)

// Validator 验证工具
type Validator struct{}

// NewValidator 创建验证工具实例
func NewValidator() *Validator {
	return &Validator{}
}

// ValidateFilePath 验证文件路径
func (v *Validator) ValidateFilePath(filePath string) error {
	if filePath == "" {
		return fmt.Errorf("文件路径不能为空")
	}

	// 检查路径是否包含非法字符
	illegalChars := []string{"<", ">", ":", "\"", "|", "?", "*"}
	for _, char := range illegalChars {
		if strings.Contains(filePath, char) {
			return fmt.Errorf("文件路径包含非法字符: %s", char)
		}
	}

	// 检查文件是否存在
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return fmt.Errorf("文件不存在: %s", filePath)
	}

	return nil
}

// ValidateDirectoryPath 验证目录路径
func (v *Validator) ValidateDirectoryPath(dirPath string) error {
	if dirPath == "" {
		return fmt.Errorf("目录路径不能为空")
	}

	// 检查路径是否包含非法字符
	illegalChars := []string{"<", ">", ":", "\"", "|", "?", "*"}
	for _, char := range illegalChars {
		if strings.Contains(dirPath, char) {
			return fmt.Errorf("目录路径包含非法字符: %s", char)
		}
	}

	// 检查目录是否存在
	if _, err := os.Stat(dirPath); os.IsNotExist(err) {
		return fmt.Errorf("目录不存在: %s", dirPath)
	}

	return nil
}

// ValidatePID 验证进程ID
func (v *Validator) ValidatePID(pid int32) error {
	if pid <= 0 {
		return fmt.Errorf("进程ID必须大于0")
	}
	return nil
}

// ValidatePort 验证端口号
func (v *Validator) ValidatePort(port int) error {
	if port < 1 || port > 65535 {
		return fmt.Errorf("端口号必须在1-65535范围内")
	}
	return nil
}

// ValidateIPAddress 验证IP地址
func (v *Validator) ValidateIPAddress(ip string) error {
	if ip == "" {
		return fmt.Errorf("IP地址不能为空")
	}

	if net.ParseIP(ip) == nil {
		return fmt.Errorf("无效的IP地址: %s", ip)
	}

	return nil
}

// ValidateIPRange 验证IP范围
func (v *Validator) ValidateIPRange(ipRange string) error {
	if ipRange == "" {
		return fmt.Errorf("IP范围不能为空")
	}

	// 检查CIDR格式
	if strings.Contains(ipRange, "/") {
		_, _, err := net.ParseCIDR(ipRange)
		if err != nil {
			return fmt.Errorf("无效的CIDR格式: %s", ipRange)
		}
		return nil
	}

	// 检查IP地址
	return v.ValidateIPAddress(ipRange)
}

// ValidateEmail 验证邮箱地址
func (v *Validator) ValidateEmail(email string) error {
	if email == "" {
		return fmt.Errorf("邮箱地址不能为空")
	}

	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	if !emailRegex.MatchString(email) {
		return fmt.Errorf("无效的邮箱地址: %s", email)
	}

	return nil
}

// ValidateURL 验证URL
func (v *Validator) ValidateURL(url string) error {
	if url == "" {
		return fmt.Errorf("URL不能为空")
	}

	// 简单的URL验证
	if !strings.HasPrefix(url, "http://") && !strings.HasPrefix(url, "https://") {
		return fmt.Errorf("URL必须以http://或https://开头: %s", url)
	}

	return nil
}

// ValidateFileName 验证文件名
func (v *Validator) ValidateFileName(fileName string) error {
	if fileName == "" {
		return fmt.Errorf("文件名不能为空")
	}

	// 检查文件名长度
	if len(fileName) > 255 {
		return fmt.Errorf("文件名过长")
	}

	// 检查非法字符
	illegalChars := []string{"<", ">", ":", "\"", "/", "\\", "|", "?", "*"}
	for _, char := range illegalChars {
		if strings.Contains(fileName, char) {
			return fmt.Errorf("文件名包含非法字符: %s", char)
		}
	}

	// 检查保留名称
	reservedNames := []string{"CON", "PRN", "AUX", "NUL", "COM1", "COM2", "COM3", "COM4", "COM5", "COM6", "COM7", "COM8", "COM9", "LPT1", "LPT2", "LPT3", "LPT4", "LPT5", "LPT6", "LPT7", "LPT8", "LPT9"}
	fileNameUpper := strings.ToUpper(fileName)
	for _, reserved := range reservedNames {
		if fileNameUpper == reserved {
			return fmt.Errorf("文件名不能使用保留名称: %s", reserved)
		}
	}

	return nil
}

// ValidateFileSize 验证文件大小
func (v *Validator) ValidateFileSize(size int64, maxSize int64) error {
	if size < 0 {
		return fmt.Errorf("文件大小不能为负数")
	}

	if maxSize > 0 && size > maxSize {
		return fmt.Errorf("文件大小超过限制: %d > %d", size, maxSize)
	}

	return nil
}

// ValidateStringLength 验证字符串长度
func (v *Validator) ValidateStringLength(str string, minLength, maxLength int) error {
	length := len(str)
	if minLength > 0 && length < minLength {
		return fmt.Errorf("字符串长度不足: %d < %d", length, minLength)
	}

	if maxLength > 0 && length > maxLength {
		return fmt.Errorf("字符串长度超限: %d > %d", length, maxLength)
	}

	return nil
}

// ValidateInteger 验证整数
func (v *Validator) ValidateInteger(value int, min, max int) error {
	if min > 0 && value < min {
		return fmt.Errorf("数值过小: %d < %d", value, min)
	}

	if max > 0 && value > max {
		return fmt.Errorf("数值过大: %d > %d", value, max)
	}

	return nil
}

// ValidateFloat 验证浮点数
func (v *Validator) ValidateFloat(value float64, min, max float64) error {
	if min > 0 && value < min {
		return fmt.Errorf("数值过小: %f < %f", value, min)
	}

	if max > 0 && value > max {
		return fmt.Errorf("数值过大: %f > %f", value, max)
	}

	return nil
}

// ValidateRegex 验证正则表达式
func (v *Validator) ValidateRegex(pattern string) error {
	if pattern == "" {
		return fmt.Errorf("正则表达式不能为空")
	}

	_, err := regexp.Compile(pattern)
	if err != nil {
		return fmt.Errorf("无效的正则表达式: %s", err.Error())
	}

	return nil
}

// ValidateFileExtension 验证文件扩展名
func (v *Validator) ValidateFileExtension(fileName string, allowedExtensions []string) error {
	if fileName == "" {
		return fmt.Errorf("文件名不能为空")
	}

	ext := strings.ToLower(filepath.Ext(fileName))
	if ext == "" {
		return fmt.Errorf("文件没有扩展名")
	}

	if len(allowedExtensions) > 0 {
		found := false
		for _, allowedExt := range allowedExtensions {
			if strings.ToLower(allowedExt) == ext {
				found = true
				break
			}
		}
		if !found {
			return fmt.Errorf("不支持的文件扩展名: %s", ext)
		}
	}

	return nil
}

// ValidatePathDepth 验证路径深度
func (v *Validator) ValidatePathDepth(path string, maxDepth int) error {
	if path == "" {
		return fmt.Errorf("路径不能为空")
	}

	depth := strings.Count(path, string(os.PathSeparator))
	if maxDepth > 0 && depth > maxDepth {
		return fmt.Errorf("路径深度超限: %d > %d", depth, maxDepth)
	}

	return nil
}

// ValidateRegistryPath 验证注册表路径
func (v *Validator) ValidateRegistryPath(regPath string) error {
	if regPath == "" {
		return fmt.Errorf("注册表路径不能为空")
	}

	// 检查注册表路径格式
	regPathRegex := regexp.MustCompile(`^[A-Z_]+\\[\\\w\\-\.]+$`)
	if !regPathRegex.MatchString(regPath) {
		return fmt.Errorf("无效的注册表路径格式: %s", regPath)
	}

	return nil
}

// ValidateProcessName 验证进程名称
func (v *Validator) ValidateProcessName(processName string) error {
	if processName == "" {
		return fmt.Errorf("进程名称不能为空")
	}

	// 检查进程名称长度
	if len(processName) > 255 {
		return fmt.Errorf("进程名称过长")
	}

	// 检查非法字符
	illegalChars := []string{"<", ">", ":", "\"", "/", "\\", "|", "?", "*"}
	for _, char := range illegalChars {
		if strings.Contains(processName, char) {
			return fmt.Errorf("进程名称包含非法字符: %s", char)
		}
	}

	return nil
}

// ValidateScanDepth 验证扫描深度
func (v *Validator) ValidateScanDepth(depth int) error {
	if depth < 0 {
		return fmt.Errorf("扫描深度不能为负数")
	}

	if depth > 100 {
		return fmt.Errorf("扫描深度过大: %d > 100", depth)
	}

	return nil
}

// ValidateTimeout 验证超时时间
func (v *Validator) ValidateTimeout(timeout int) error {
	if timeout < 0 {
		return fmt.Errorf("超时时间不能为负数")
	}

	if timeout > 3600 {
		return fmt.Errorf("超时时间过长: %d > 3600", timeout)
	}

	return nil
}

// ValidateAPIKey 验证API密钥
func (v *Validator) ValidateAPIKey(apiKey string) error {
	if apiKey == "" {
		return fmt.Errorf("API密钥不能为空")
	}

	if len(apiKey) < 16 {
		return fmt.Errorf("API密钥长度不足: %d < 16", len(apiKey))
	}

	if len(apiKey) > 256 {
		return fmt.Errorf("API密钥过长: %d > 256", len(apiKey))
	}

	return nil
}

// ValidateHash 验证哈希值
func (v *Validator) ValidateHash(hash string, algorithm string) error {
	if hash == "" {
		return fmt.Errorf("哈希值不能为空")
	}

	// 检查哈希值格式
	hashRegex := regexp.MustCompile(`^[a-fA-F0-9]+$`)
	if !hashRegex.MatchString(hash) {
		return fmt.Errorf("无效的哈希值格式: %s", hash)
	}

	// 检查哈希值长度
	expectedLength := 0
	switch algorithm {
	case "md5":
		expectedLength = 32
	case "sha1":
		expectedLength = 40
	case "sha256":
		expectedLength = 64
	case "sha512":
		expectedLength = 128
	default:
		return fmt.Errorf("不支持的哈希算法: %s", algorithm)
	}

	if len(hash) != expectedLength {
		return fmt.Errorf("哈希值长度不正确: %d != %d", len(hash), expectedLength)
	}

	return nil
}

// ValidateNumericString 验证数字字符串
func (v *Validator) ValidateNumericString(str string) error {
	if str == "" {
		return fmt.Errorf("字符串不能为空")
	}

	_, err := strconv.Atoi(str)
	if err != nil {
		return fmt.Errorf("字符串不是有效的数字: %s", str)
	}

	return nil
}

// ValidateHexString 验证十六进制字符串
func (v *Validator) ValidateHexString(str string) error {
	if str == "" {
		return fmt.Errorf("字符串不能为空")
	}

	hexRegex := regexp.MustCompile(`^[a-fA-F0-9]+$`)
	if !hexRegex.MatchString(str) {
		return fmt.Errorf("字符串不是有效的十六进制: %s", str)
	}

	return nil
}
