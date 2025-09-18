package file

import (
	"archive/tar"
	"archive/zip"
	"compress/gzip"
	"context"
	"crypto/md5"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"math"
	"os"
	"path/filepath"
	"strings"
	"time"

	"yara-security-service/internal/models"
	"yara-security-service/pkg/security"

	"github.com/sirupsen/logrus"
)

// FileAnalyzer 文件分析器
type FileAnalyzer struct {
	logger  *logrus.Logger
	scanner *security.Scanner
}

// NewFileAnalyzer 创建文件分析器
func NewFileAnalyzer(logger *logrus.Logger, scanner *security.Scanner) *FileAnalyzer {
	return &FileAnalyzer{
		logger:  logger,
		scanner: scanner,
	}
}

// DeepFileAnalysis 深度文件分析
func (fa *FileAnalyzer) DeepFileAnalysis(ctx context.Context, filePath string) (*models.DeepFileAnalysis, error) {
	analysis := &models.DeepFileAnalysis{
		FilePath:     filePath,
		Timestamp:    time.Now(),
		IsSuspicious: false,
		ThreatLevel:  "low",
	}

	// 1. 基础文件信息分析
	fileInfo, err := fa.analyzeBasicFileInfo(filePath)
	if err != nil {
		return nil, fmt.Errorf("分析基础文件信息失败: %w", err)
	}
	analysis.BasicInfo = fileInfo

	// 2. 文件内容深度分析
	contentAnalysis, err := fa.analyzeFileContent(ctx, filePath)
	if err != nil {
		fa.logger.Warnf("分析文件内容失败: %v", err)
	} else {
		analysis.ContentAnalysis = contentAnalysis
	}

	// 3. 文件结构分析
	structureAnalysis, err := fa.analyzeFileStructure(filePath)
	if err != nil {
		fa.logger.Warnf("分析文件结构失败: %v", err)
	} else {
		analysis.StructureAnalysis = structureAnalysis
	}

	// 4. 安全威胁分析
	threatAnalysis, err := fa.analyzeSecurityThreats(ctx, filePath)
	if err != nil {
		fa.logger.Warnf("分析安全威胁失败: %v", err)
	} else {
		analysis.ThreatAnalysis = threatAnalysis
	}

	// 5. 综合风险评估
	fa.assessOverallRisk(analysis)

	return analysis, nil
}

// analyzeBasicFileInfo 分析基础文件信息
func (fa *FileAnalyzer) analyzeBasicFileInfo(filePath string) (*models.BasicFileInfo, error) {
	file, err := os.Stat(filePath)
	if err != nil {
		return nil, fmt.Errorf("获取文件信息失败: %w", err)
	}

	// 计算文件哈希
	md5Hash, sha256Hash, err := fa.calculateFileHashes(filePath)
	if err != nil {
		fa.logger.Warnf("计算文件哈希失败: %v", err)
	}

	// 检测文件类型
	fileType := fa.detectFileType(filePath)

	// 检测文件扩展名
	extension := filepath.Ext(filePath)

	// 检测文件大小类别
	sizeCategory := fa.categorizeFileSize(file.Size())

	// 获取真实的文件时间信息
	creationTime, accessTime := fa.getRealFileTimes(filePath, file)

	return &models.BasicFileInfo{
		Name:         file.Name(),
		Size:         file.Size(),
		SizeCategory: sizeCategory,
		Extension:    extension,
		FileType:     fileType,
		ModTime:      file.ModTime(),
		CreateTime:   creationTime,
		AccessTime:   accessTime,
		Permissions:  file.Mode().String(),
		MD5Hash:      md5Hash,
		SHA256Hash:   sha256Hash,
		IsHidden:     fa.isFileHidden(file),
		IsSystem:     fa.isSystemFile(filePath),
	}, nil
}

// analyzeFileContent 分析文件内容
func (fa *FileAnalyzer) analyzeFileContent(ctx context.Context, filePath string) (*models.ContentAnalysis, error) {
	analysis := &models.ContentAnalysis{
		Timestamp: time.Now(),
	}

	// 读取文件内容
	content, err := fa.readFileContent(filePath)
	if err != nil {
		return nil, fmt.Errorf("读取文件内容失败: %w", err)
	}

	// 分析文本内容
	textAnalysis := fa.analyzeTextContent(content)
	analysis.TextAnalysis = textAnalysis

	// 分析二进制内容
	binaryAnalysis := fa.analyzeBinaryContent(content)
	analysis.BinaryAnalysis = binaryAnalysis

	// 分析编码内容
	encodingAnalysis := fa.analyzeEncoding(content)
	analysis.EncodingAnalysis = encodingAnalysis

	// 检测可疑模式
	suspiciousPatterns := fa.detectSuspiciousPatterns(content)
	analysis.SuspiciousPatterns = suspiciousPatterns

	return analysis, nil
}

// analyzeFileStructure 分析文件结构
func (fa *FileAnalyzer) analyzeFileStructure(filePath string) (*models.StructureAnalysis, error) {
	analysis := &models.StructureAnalysis{
		Timestamp: time.Now(),
	}

	// 分析PE文件结构
	if strings.HasSuffix(strings.ToLower(filePath), ".exe") ||
		strings.HasSuffix(strings.ToLower(filePath), ".dll") {
		peAnalysis := fa.analyzePEStructure(filePath)
		analysis.PEAnalysis = peAnalysis
	}

	// 分析压缩文件结构
	if strings.HasSuffix(strings.ToLower(filePath), ".zip") ||
		strings.HasSuffix(strings.ToLower(filePath), ".rar") {
		archiveAnalysis := fa.analyzeArchiveStructure(filePath)
		analysis.ArchiveAnalysis = archiveAnalysis
	}

	// 分析文档文件结构
	if strings.HasSuffix(strings.ToLower(filePath), ".doc") ||
		strings.HasSuffix(strings.ToLower(filePath), ".docx") {
		documentAnalysis := fa.analyzeDocumentStructure(filePath)
		analysis.DocumentAnalysis = documentAnalysis
	}

	return analysis, nil
}

// analyzeSecurityThreats 分析安全威胁
func (fa *FileAnalyzer) analyzeSecurityThreats(ctx context.Context, filePath string) (*models.ThreatAnalysis, error) {
	analysis := &models.ThreatAnalysis{
		Timestamp: time.Now(),
	}

	// 使用安全扫描器扫描文件
	if fa.scanner != nil {
		scanResult, err := fa.scanner.ScanFile(ctx, filePath)
		if err == nil {
			analysis.ScanResult = scanResult
			analysis.IsInfected = scanResult.IsInfected
			analysis.Threats = scanResult.Threats
		}
	}

	// 检测文件行为特征
	behaviorAnalysis := fa.analyzeFileBehavior(filePath)
	analysis.BehaviorAnalysis = behaviorAnalysis

	// 检测文件签名
	signatureAnalysis := fa.analyzeFileSignature(filePath)
	analysis.SignatureAnalysis = signatureAnalysis

	// 检测文件熵值
	entropyAnalysis := fa.analyzeFileEntropy(filePath)
	analysis.EntropyAnalysis = entropyAnalysis

	return analysis, nil
}

// assessOverallRisk 综合风险评估
func (fa *FileAnalyzer) assessOverallRisk(analysis *models.DeepFileAnalysis) {
	riskScore := 0
	riskFactors := []string{}

	// 基于文件类型评估风险
	if analysis.BasicInfo != nil {
		if analysis.BasicInfo.FileType == "executable" {
			riskScore += 30
			riskFactors = append(riskFactors, "可执行文件")
		}
		if analysis.BasicInfo.FileType == "script" {
			riskScore += 25
			riskFactors = append(riskFactors, "脚本文件")
		}
		if analysis.BasicInfo.IsSystem {
			riskScore -= 10
			riskFactors = append(riskFactors, "系统文件")
		}
	}

	// 基于内容分析评估风险
	if analysis.ContentAnalysis != nil {
		if len(analysis.ContentAnalysis.SuspiciousPatterns) > 0 {
			riskScore += len(analysis.ContentAnalysis.SuspiciousPatterns) * 10
			riskFactors = append(riskFactors, "可疑内容模式")
		}
		if analysis.ContentAnalysis.BinaryAnalysis != nil &&
			analysis.ContentAnalysis.BinaryAnalysis.HighEntropy {
			riskScore += 20
			riskFactors = append(riskFactors, "高熵值")
		}
	}

	// 基于威胁分析评估风险
	if analysis.ThreatAnalysis != nil {
		if analysis.ThreatAnalysis.IsInfected {
			riskScore += 50
			riskFactors = append(riskFactors, "检测到恶意代码")
		}
		if len(analysis.ThreatAnalysis.Threats) > 0 {
			riskScore += len(analysis.ThreatAnalysis.Threats) * 15
			riskFactors = append(riskFactors, "多个威胁")
		}
	}

	// 设置综合风险等级
	analysis.RiskScore = riskScore
	analysis.RiskFactors = riskFactors

	if riskScore >= 80 {
		analysis.ThreatLevel = "critical"
		analysis.IsSuspicious = true
	} else if riskScore >= 60 {
		analysis.ThreatLevel = "high"
		analysis.IsSuspicious = true
	} else if riskScore >= 40 {
		analysis.ThreatLevel = "medium"
		analysis.IsSuspicious = true
	} else if riskScore >= 20 {
		analysis.ThreatLevel = "low"
		analysis.IsSuspicious = false
	} else {
		analysis.ThreatLevel = "safe"
		analysis.IsSuspicious = false
	}
}

// calculateFileHashes 计算文件哈希
func (fa *FileAnalyzer) calculateFileHashes(filePath string) (string, string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return "", "", err
	}
	defer file.Close()

	md5Hash := md5.New()
	sha256Hash := sha256.New()
	multiWriter := io.MultiWriter(md5Hash, sha256Hash)

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
			return "", "", err
		}
	}

	md5Result := hex.EncodeToString(md5Hash.Sum(nil))
	sha256Result := hex.EncodeToString(sha256Hash.Sum(nil))

	return md5Result, sha256Result, nil
}

// detectFileType 检测文件类型
func (fa *FileAnalyzer) detectFileType(filePath string) string {
	extension := strings.ToLower(filepath.Ext(filePath))

	switch extension {
	case ".exe", ".dll", ".sys":
		return "executable"
	case ".ps1", ".bat", ".cmd", ".vbs", ".js":
		return "script"
	case ".doc", ".docx", ".xls", ".xlsx", ".ppt", ".pptx":
		return "document"
	case ".zip", ".rar", ".7z", ".tar", ".gz":
		return "archive"
	case ".jpg", ".jpeg", ".png", ".gif", ".bmp":
		return "image"
	case ".mp3", ".mp4", ".avi", ".wav":
		return "media"
	case ".txt", ".log", ".ini", ".cfg":
		return "text"
	default:
		return "unknown"
	}
}

// categorizeFileSize 分类文件大小
func (fa *FileAnalyzer) categorizeFileSize(size int64) string {
	switch {
	case size < 1024:
		return "tiny"
	case size < 1024*1024:
		return "small"
	case size < 10*1024*1024:
		return "medium"
	case size < 100*1024*1024:
		return "large"
	default:
		return "huge"
	}
}

// readFileContent 读取文件内容
func (fa *FileAnalyzer) readFileContent(filePath string) ([]byte, error) {
	// 限制读取大小，避免内存问题
	maxSize := int64(10 * 1024 * 1024) // 10MB

	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	// 获取文件大小
	stat, err := file.Stat()
	if err != nil {
		return nil, err
	}

	// 如果文件太大，只读取前10MB
	readSize := stat.Size()
	if readSize > maxSize {
		readSize = maxSize
	}

	content := make([]byte, readSize)
	_, err = file.Read(content)
	return content, err
}

// analyzeTextContent 分析文本内容
func (fa *FileAnalyzer) analyzeTextContent(content []byte) *models.TextAnalysis {
	analysis := &models.TextAnalysis{
		Timestamp: time.Now(),
	}

	// 检测文本编码
	analysis.Encoding = fa.detectTextEncoding(content)

	// 分析文本特征
	analysis.LineCount = strings.Count(string(content), "\n")
	analysis.WordCount = len(strings.Fields(string(content)))
	analysis.CharacterCount = len(content)

	// 检测语言
	analysis.Language = fa.detectLanguage(content)

	// 检测特殊字符
	analysis.SpecialCharacters = fa.detectSpecialCharacters(content)

	return analysis
}

// analyzeBinaryContent 分析二进制内容
func (fa *FileAnalyzer) analyzeBinaryContent(content []byte) *models.BinaryAnalysis {
	analysis := &models.BinaryAnalysis{
		Timestamp: time.Now(),
	}

	// 计算熵值
	analysis.Entropy = fa.calculateEntropy(content)
	analysis.HighEntropy = analysis.Entropy > 7.0

	// 检测文件头
	analysis.FileHeader = fa.detectFileHeader(content)

	// 检测字符串
	analysis.Strings = fa.extractStrings(content)

	// 检测十六进制模式
	analysis.HexPatterns = fa.detectHexPatterns(content)

	return analysis
}

// analyzeEncoding 分析编码
func (fa *FileAnalyzer) analyzeEncoding(content []byte) *models.EncodingAnalysis {
	analysis := &models.EncodingAnalysis{
		Timestamp: time.Now(),
	}

	// 检测编码类型
	analysis.EncodingType = fa.detectEncodingType(content)

	// 检测是否编码
	analysis.IsEncoded = fa.isContentEncoded(content)

	// 检测编码算法
	analysis.EncodingAlgorithm = fa.detectEncodingAlgorithm(content)

	return analysis
}

// detectSuspiciousPatterns 检测可疑模式
func (fa *FileAnalyzer) detectSuspiciousPatterns(content []byte) []string {
	var patterns []string
	contentStr := strings.ToLower(string(content))

	// 检测恶意URL
	urlPatterns := []string{
		"http://", "https://", "ftp://", "file://",
	}
	for _, pattern := range urlPatterns {
		if strings.Contains(contentStr, pattern) {
			patterns = append(patterns, "URL模式")
			break
		}
	}

	// 检测命令执行
	cmdPatterns := []string{
		"cmd.exe", "powershell", "exec", "system",
		"shell", "command", "run", "execute",
	}
	for _, pattern := range cmdPatterns {
		if strings.Contains(contentStr, pattern) {
			patterns = append(patterns, "命令执行")
			break
		}
	}

	// 检测注册表操作
	regPatterns := []string{
		"reg add", "reg delete", "reg set", "reg query",
		"registry", "hkey_", "software\\", "system\\",
	}
	for _, pattern := range regPatterns {
		if strings.Contains(contentStr, pattern) {
			patterns = append(patterns, "注册表操作")
			break
		}
	}

	// 检测网络连接
	netPatterns := []string{
		"connect", "socket", "bind", "listen",
		"accept", "send", "recv", "gethostbyname",
	}
	for _, pattern := range netPatterns {
		if strings.Contains(contentStr, pattern) {
			patterns = append(patterns, "网络连接")
			break
		}
	}

	return patterns
}

// 其他辅助方法...
func (fa *FileAnalyzer) getFileCreateTime(file os.FileInfo) time.Time {
	// 使用新的文件时间工具获取真实的创建时间
	creationTime, err := security.GetFileCreationTime(file.Name())
	if err != nil {
		fa.logger.Warnf("获取文件创建时间失败: %v，使用修改时间作为替代", err)
		return file.ModTime()
	}

	// 检查时间是否为零值（表示不可用）
	zeroTime := time.Time{}
	if creationTime.Equal(zeroTime) {
		fa.logger.Debugf("文件创建时间不可用，使用修改时间作为替代: %s", file.Name())
		return file.ModTime()
	}

	return creationTime
}

// getRealFileTimes 获取真实的文件时间信息
func (fa *FileAnalyzer) getRealFileTimes(filePath string, file os.FileInfo) (time.Time, time.Time) {
	// 使用新的文件时间工具获取真实的创建时间和访问时间
	creationTime, err := security.GetFileCreationTime(filePath)
	if err != nil {
		fa.logger.Warnf("获取文件创建时间失败: %v，使用修改时间作为替代", err)
		creationTime = file.ModTime()
	}

	// 使用新的函数获取访问时间，并强制更新
	accessTime, err := security.GetFileAccessTimeWithUpdate(filePath)
	if err != nil {
		fa.logger.Warnf("获取文件访问时间失败: %v，使用修改时间作为替代", err)
		accessTime = file.ModTime()
	}

	// 检查时间是否为零值（表示不可用）
	zeroTime := time.Time{}
	if creationTime.Equal(zeroTime) {
		fa.logger.Debugf("文件创建时间不可用，使用修改时间作为替代: %s", filePath)
		creationTime = file.ModTime()
	}

	if accessTime.Equal(zeroTime) {
		fa.logger.Debugf("文件访问时间不可用，使用修改时间作为替代: %s", filePath)
		accessTime = file.ModTime()
	}

	return creationTime, accessTime
}

func (fa *FileAnalyzer) getFileAccessTime(file os.FileInfo) time.Time {
	// 使用新的文件时间工具获取真实的访问时间
	accessTime, err := security.GetFileAccessTime(file.Name())
	if err != nil {
		fa.logger.Warnf("获取文件访问时间失败: %v，使用修改时间作为替代", err)
		return file.ModTime()
	}

	// 检查时间是否为零值（表示不可用）
	zeroTime := time.Time{}
	if accessTime.Equal(zeroTime) {
		fa.logger.Debugf("文件访问时间不可用，使用修改时间作为替代: %s", file.Name())
		return file.ModTime()
	}

	return accessTime
}

func (fa *FileAnalyzer) isFileHidden(file os.FileInfo) bool {
	// 检查文件是否隐藏
	// 在Windows上，使用文件属性检查隐藏状态

	// 获取文件路径
	filePath := file.Name()

	// 检查文件名是否以点开头（Unix风格的隐藏文件）
	if strings.HasPrefix(filePath, ".") {
		return true
	}

	// 检查是否为Windows系统隐藏文件
	hiddenFiles := []string{
		"thumbs.db", "desktop.ini", ".ds_store", "autorun.inf",
		"ntuser.dat", "ntuser.ini", "ntuser.dat.log",
	}

	fileNameLower := strings.ToLower(filePath)
	for _, hidden := range hiddenFiles {
		if fileNameLower == hidden {
			return true
		}
	}

	// 检查系统目录中的文件
	// systemPaths := []string{
	// 	"c:\\windows\\", "c:\\system32\\", "c:\\syswow64\\",
	// 	"c:\\programdata\\", "c:\\$recycle.bin\\",
	// }

	// 这里需要完整的文件路径，但os.FileInfo只提供文件名
	// 在实际实现中，应该传入完整路径或使用Windows API
	// 暂时使用文件名检查

	// 检查是否为系统文件（通过扩展名）
	systemExtensions := []string{
		".sys", ".dll", ".exe", ".drv", ".vxd", ".386",
	}

	ext := strings.ToLower(filepath.Ext(filePath))
	for _, sysExt := range systemExtensions {
		if ext == sysExt {
			// 系统文件通常也是隐藏的
			return true
		}
	}

	// 在实际实现中，这里应该：
	// 1. 使用Windows API GetFileAttributes检查FILE_ATTRIBUTE_HIDDEN
	// 2. 检查文件是否在系统目录中
	// 3. 检查文件权限和属性

	return false
}

func (fa *FileAnalyzer) isSystemFile(filePath string) bool {
	// 检查是否为系统文件
	systemPaths := []string{
		"C:\\Windows\\", "C:\\System32\\", "C:\\SysWOW64\\",
	}
	for _, path := range systemPaths {
		if strings.HasPrefix(strings.ToLower(filePath), strings.ToLower(path)) {
			return true
		}
	}
	return false
}

func (fa *FileAnalyzer) detectTextEncoding(content []byte) string {
	// 简单的编码检测
	if len(content) >= 3 && content[0] == 0xEF && content[1] == 0xBB && content[2] == 0xBF {
		return "UTF-8-BOM"
	}
	return "UTF-8"
}

func (fa *FileAnalyzer) detectLanguage(content []byte) string {
	// 简单的语言检测
	contentStr := string(content)
	if strings.Contains(contentStr, "function") || strings.Contains(contentStr, "var ") {
		return "JavaScript"
	}
	if strings.Contains(contentStr, "powershell") || strings.Contains(contentStr, "cmd") {
		return "Batch"
	}
	return "Unknown"
}

func (fa *FileAnalyzer) detectSpecialCharacters(content []byte) []string {
	var specialChars []string
	contentStr := string(content)

	if strings.Contains(contentStr, "\\x") {
		specialChars = append(specialChars, "十六进制编码")
	}
	if strings.Contains(contentStr, "\\u") {
		specialChars = append(specialChars, "Unicode编码")
	}
	if strings.Contains(contentStr, "base64") {
		specialChars = append(specialChars, "Base64编码")
	}

	return specialChars
}

func (fa *FileAnalyzer) calculateEntropy(data []byte) float64 {
	if len(data) == 0 {
		return 0
	}

	// 计算字节频率
	freq := make(map[byte]int)
	for _, b := range data {
		freq[b]++
	}

	// 计算熵值
	entropy := 0.0
	dataLen := float64(len(data))
	for _, count := range freq {
		p := float64(count) / dataLen
		if p > 0 {
			entropy -= p * math.Log2(p)
		}
	}

	return entropy
}

func (fa *FileAnalyzer) detectFileHeader(content []byte) string {
	if len(content) < 4 {
		return "unknown"
	}

	// 检测常见文件头
	if len(content) >= 2 && content[0] == 0x4D && content[1] == 0x5A {
		return "PE/EXE"
	}
	if len(content) >= 4 && content[0] == 0x50 && content[1] == 0x4B {
		return "ZIP"
	}
	if len(content) >= 4 && content[0] == 0x52 && content[1] == 0x61 && content[2] == 0x72 {
		return "RAR"
	}
	if len(content) >= 4 && content[0] == 0x25 && content[1] == 0x50 && content[2] == 0x44 && content[3] == 0x46 {
		return "PDF"
	}

	return "unknown"
}

func (fa *FileAnalyzer) extractStrings(content []byte) []string {
	var strings []string
	currentString := ""

	for _, b := range content {
		if b >= 32 && b <= 126 {
			currentString += string(b)
		} else {
			if len(currentString) >= 4 {
				strings = append(strings, currentString)
			}
			currentString = ""
		}
	}

	if len(currentString) >= 4 {
		strings = append(strings, currentString)
	}

	return strings
}

func (fa *FileAnalyzer) detectHexPatterns(content []byte) []string {
	var patterns []string

	// 检测连续的十六进制字节
	for i := 0; i < len(content)-3; i++ {
		if content[i] == 0x90 && content[i+1] == 0x90 && content[i+2] == 0x90 {
			patterns = append(patterns, "NOP sled")
		}
		if content[i] == 0x31 && content[i+1] == 0xC0 {
			patterns = append(patterns, "XOR EAX, EAX")
		}
	}

	return patterns
}

func (fa *FileAnalyzer) detectEncodingType(content []byte) string {
	// 简单的编码类型检测
	if len(content) >= 3 && content[0] == 0xEF && content[1] == 0xBB && content[2] == 0xBF {
		return "UTF-8-BOM"
	}
	return "UTF-8"
}

func (fa *FileAnalyzer) isContentEncoded(content []byte) bool {
	// 检测是否包含编码内容
	contentStr := string(content)
	return strings.Contains(contentStr, "base64") ||
		strings.Contains(contentStr, "\\x") ||
		strings.Contains(contentStr, "\\u")
}

func (fa *FileAnalyzer) detectEncodingAlgorithm(content []byte) string {
	contentStr := string(content)
	if strings.Contains(contentStr, "base64") {
		return "Base64"
	}
	if strings.Contains(contentStr, "\\x") {
		return "Hex"
	}
	if strings.Contains(contentStr, "\\u") {
		return "Unicode"
	}
	return "None"
}

// 简化实现的PE、Archive、Document分析
func (fa *FileAnalyzer) analyzePEStructure(filePath string) *models.PEAnalysis {
	analysis := &models.PEAnalysis{
		Timestamp: time.Now(),
		IsValidPE: false,
	}

	// 读取文件头
	file, err := os.Open(filePath)
	if err != nil {
		fa.logger.Errorf("打开文件失败: %v", err)
		return analysis
	}
	defer file.Close()

	// 读取PE文件头
	header := make([]byte, 64)
	_, err = file.Read(header)
	if err != nil {
		fa.logger.Errorf("读取文件头失败: %v", err)
		return analysis
	}

	// 检查PE文件签名 (MZ)
	if len(header) >= 2 && header[0] == 0x4D && header[1] == 0x5A {
		analysis.IsValidPE = true

		// 获取PE头偏移
		if len(header) >= 60 {
			peOffset := int(header[60]) | (int(header[61]) << 8) | (int(header[62]) << 16) | (int(header[63]) << 24)

			// 尝试读取PE头
			if peOffset > 0 {
				file.Seek(int64(peOffset), 0)
				peHeader := make([]byte, 24)
				if _, err := file.Read(peHeader); err == nil {
					// 检查PE签名
					if len(peHeader) >= 4 && peHeader[0] == 0x50 && peHeader[1] == 0x45 && peHeader[2] == 0x00 && peHeader[3] == 0x00 {
						// PE文件验证成功
						fa.logger.Debugf("PE文件验证成功: %s", filePath)
					}
				}
			}
		}
	}

	return analysis
}

func (fa *FileAnalyzer) analyzeArchiveStructure(filePath string) *models.ArchiveAnalysis {
	analysis := &models.ArchiveAnalysis{
		Timestamp:        time.Now(),
		IsValidArchive:   false,
		ArchiveType:      "Unknown",
		FileCount:        0,
		CompressedSize:   0,
		UncompressedSize: 0,
		Files:            []string{},
	}

	// 读取文件头
	file, err := os.Open(filePath)
	if err != nil {
		fa.logger.Errorf("打开文件失败: %v", err)
		return analysis
	}
	defer file.Close()

	// 读取文件头进行格式检测
	header := make([]byte, 16)
	_, err = file.Read(header)
	if err != nil {
		fa.logger.Errorf("读取文件头失败: %v", err)
		return analysis
	}

	// 检测压缩包格式
	archiveType := fa.detectArchiveType(header)
	if archiveType == "Unknown" {
		fa.logger.Debugf("无法识别的压缩包格式: %s", filePath)
		return analysis
	}

	analysis.ArchiveType = archiveType
	analysis.IsValidArchive = true

	// 获取文件大小
	fileInfo, err := file.Stat()
	if err != nil {
		fa.logger.Errorf("获取文件信息失败: %v", err)
		return analysis
	}
	analysis.CompressedSize = fileInfo.Size()

	// 尝试分析压缩包内容
	fa.analyzeArchiveContent(filePath, archiveType, analysis)

	return analysis
}

func (fa *FileAnalyzer) detectArchiveType(header []byte) string {
	// ZIP格式检测 (PK\x03\x04)
	if len(header) >= 4 && header[0] == 0x50 && header[1] == 0x4B && header[2] == 0x03 && header[3] == 0x04 {
		return "ZIP"
	}

	// RAR格式检测 (Rar!\x1A\x07)
	if len(header) >= 6 && header[0] == 0x52 && header[1] == 0x61 && header[2] == 0x72 && header[3] == 0x21 && header[4] == 0x1A && header[5] == 0x07 {
		return "RAR"
	}

	// 7Z格式检测 (7z\xBC\xAF\x27\x1C)
	if len(header) >= 6 && header[0] == 0x37 && header[1] == 0x7A && header[2] == 0xBC && header[3] == 0xAF && header[4] == 0x27 && header[5] == 0x1C {
		return "7Z"
	}

	// TAR格式检测 (ustar)
	if len(header) >= 8 && string(header[257:262]) == "ustar" {
		return "TAR"
	}

	// GZIP格式检测 (\x1F\x8B)
	if len(header) >= 2 && header[0] == 0x1F && header[1] == 0x8B {
		return "GZIP"
	}

	// BZIP2格式检测 (BZ)
	if len(header) >= 2 && header[0] == 0x42 && header[1] == 0x5A {
		return "BZIP2"
	}

	return "Unknown"
}

func (fa *FileAnalyzer) analyzeArchiveContent(filePath, archiveType string, analysis *models.ArchiveAnalysis) {
	// 根据压缩包类型进行内容分析
	switch archiveType {
	case "ZIP":
		fa.analyzeZipContent(filePath, analysis)
	case "RAR":
		fa.analyzeRarContent(filePath, analysis)
	case "7Z":
		fa.analyze7zContent(filePath, analysis)
	case "TAR":
		fa.analyzeTarContent(filePath, analysis)
	case "GZIP":
		fa.analyzeGzipContent(filePath, analysis)
	case "BZIP2":
		fa.analyzeBzip2Content(filePath, analysis)
	default:
		fa.logger.Debugf("暂不支持分析压缩包类型: %s", archiveType)
	}
}

func (fa *FileAnalyzer) analyzeZipContent(filePath string, analysis *models.ArchiveAnalysis) {
	// 使用archive/zip包分析ZIP文件
	reader, err := zip.OpenReader(filePath)
	if err != nil {
		fa.logger.Errorf("打开ZIP文件失败: %v", err)
		return
	}
	defer reader.Close()

	analysis.FileCount = len(reader.File)
	var totalUncompressedSize int64

	for _, file := range reader.File {
		analysis.Files = append(analysis.Files, file.Name)
		totalUncompressedSize += int64(file.UncompressedSize64)
	}

	analysis.UncompressedSize = totalUncompressedSize
	fa.logger.Debugf("ZIP文件分析完成: %d 个文件, 压缩大小: %d, 解压大小: %d",
		analysis.FileCount, analysis.CompressedSize, analysis.UncompressedSize)
}

func (fa *FileAnalyzer) analyzeRarContent(filePath string, analysis *models.ArchiveAnalysis) {
	// RAR文件分析（简化实现）
	// 在实际实现中，需要使用专门的RAR库
	analysis.FileCount = 1                              // 简化实现
	analysis.UncompressedSize = analysis.CompressedSize // 估算
	analysis.Files = append(analysis.Files, "unknown.rar")
	fa.logger.Debugf("RAR文件分析完成（简化实现）")
}

func (fa *FileAnalyzer) analyze7zContent(filePath string, analysis *models.ArchiveAnalysis) {
	// 7Z文件分析（简化实现）
	// 在实际实现中，需要使用专门的7Z库
	analysis.FileCount = 1                              // 简化实现
	analysis.UncompressedSize = analysis.CompressedSize // 估算
	analysis.Files = append(analysis.Files, "unknown.7z")
	fa.logger.Debugf("7Z文件分析完成（简化实现）")
}

func (fa *FileAnalyzer) analyzeTarContent(filePath string, analysis *models.ArchiveAnalysis) {
	// TAR文件分析
	file, err := os.Open(filePath)
	if err != nil {
		fa.logger.Errorf("打开TAR文件失败: %v", err)
		return
	}
	defer file.Close()

	reader := tar.NewReader(file)
	fileCount := 0
	var totalSize int64

	for {
		header, err := reader.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			fa.logger.Errorf("读取TAR文件头失败: %v", err)
			break
		}

		analysis.Files = append(analysis.Files, header.Name)
		totalSize += header.Size
		fileCount++
	}

	analysis.FileCount = fileCount
	analysis.UncompressedSize = totalSize
	fa.logger.Debugf("TAR文件分析完成: %d 个文件, 总大小: %d", fileCount, totalSize)
}

func (fa *FileAnalyzer) analyzeGzipContent(filePath string, analysis *models.ArchiveAnalysis) {
	// GZIP文件分析
	file, err := os.Open(filePath)
	if err != nil {
		fa.logger.Errorf("打开GZIP文件失败: %v", err)
		return
	}
	defer file.Close()

	reader, err := gzip.NewReader(file)
	if err != nil {
		fa.logger.Errorf("创建GZIP读取器失败: %v", err)
		return
	}
	defer reader.Close()

	// 读取解压后的内容来计算大小
	content, err := io.ReadAll(reader)
	if err != nil {
		fa.logger.Errorf("读取GZIP内容失败: %v", err)
		return
	}

	analysis.FileCount = 1
	analysis.UncompressedSize = int64(len(content))
	analysis.Files = append(analysis.Files, "gzip_content")
	fa.logger.Debugf("GZIP文件分析完成: 压缩大小: %d, 解压大小: %d",
		analysis.CompressedSize, analysis.UncompressedSize)
}

func (fa *FileAnalyzer) analyzeBzip2Content(filePath string, analysis *models.ArchiveAnalysis) {
	// BZIP2文件分析（简化实现）
	// 在实际实现中，需要使用专门的BZIP2库
	analysis.FileCount = 1                              // 简化实现
	analysis.UncompressedSize = analysis.CompressedSize // 估算
	analysis.Files = append(analysis.Files, "unknown.bz2")
	fa.logger.Debugf("BZIP2文件分析完成（简化实现）")
}

func (fa *FileAnalyzer) analyzeDocumentStructure(filePath string) *models.DocumentAnalysis {
	analysis := &models.DocumentAnalysis{
		Timestamp:       time.Now(),
		IsValidDocument: false,
		DocumentType:    "Unknown",
		HasMacros:       false,
		MacroCount:      0,
		EmbeddedObjects: []string{},
		Metadata:        make(map[string]string),
	}

	// 读取文件头
	file, err := os.Open(filePath)
	if err != nil {
		fa.logger.Errorf("打开文件失败: %v", err)
		return analysis
	}
	defer file.Close()

	// 读取文件头进行格式检测
	header := make([]byte, 16)
	_, err = file.Read(header)
	if err != nil {
		fa.logger.Errorf("读取文件头失败: %v", err)
		return analysis
	}

	// 检测文档格式
	documentType := fa.detectDocumentType(header, filePath)
	if documentType == "Unknown" {
		fa.logger.Debugf("无法识别的文档格式: %s", filePath)
		return analysis
	}

	analysis.DocumentType = documentType
	analysis.IsValidDocument = true

	// 根据文档类型进行详细分析
	fa.analyzeDocumentContent(filePath, documentType, analysis)

	return analysis
}

func (fa *FileAnalyzer) detectDocumentType(header []byte, filePath string) string {
	// 检查文件扩展名
	ext := strings.ToLower(filepath.Ext(filePath))

	// Office文档格式检测
	// DOCX/XLSX/PPTX (ZIP格式，包含PK\x03\x04)
	if len(header) >= 4 && header[0] == 0x50 && header[1] == 0x4B && header[2] == 0x03 && header[3] == 0x04 {
		switch ext {
		case ".docx":
			return "DOCX"
		case ".xlsx":
			return "XLSX"
		case ".pptx":
			return "PPTX"
		case ".zip":
			return "ZIP"
		}
	}

	// PDF格式检测 (%PDF)
	if len(header) >= 4 && header[0] == 0x25 && header[1] == 0x50 && header[2] == 0x44 && header[3] == 0x46 {
		return "PDF"
	}

	// RTF格式检测 ({\rtf)
	if len(header) >= 5 && header[0] == 0x7B && header[1] == 0x5C && header[2] == 0x72 && header[3] == 0x74 && header[4] == 0x66 {
		return "RTF"
	}

	// TXT格式检测（纯文本）
	if fa.isTextFile(header) {
		return "TXT"
	}

	// 根据扩展名判断
	switch ext {
	case ".doc":
		return "DOC"
	case ".xls":
		return "XLS"
	case ".ppt":
		return "PPT"
	case ".odt":
		return "ODT"
	case ".ods":
		return "ODS"
	case ".odp":
		return "ODP"
	case ".html", ".htm":
		return "HTML"
	case ".xml":
		return "XML"
	case ".json":
		return "JSON"
	case ".csv":
		return "CSV"
	case ".md":
		return "Markdown"
	}

	return "Unknown"
}

func (fa *FileAnalyzer) isTextFile(header []byte) bool {
	// 检查是否为可打印的ASCII字符
	for i := 0; i < len(header) && i < 16; i++ {
		if header[i] < 0x20 && header[i] != 0x09 && header[i] != 0x0A && header[i] != 0x0D {
			return false
		}
	}
	return true
}

func (fa *FileAnalyzer) analyzeDocumentContent(filePath, documentType string, analysis *models.DocumentAnalysis) {
	// 根据文档类型进行内容分析
	switch documentType {
	case "DOCX", "XLSX", "PPTX":
		fa.analyzeOfficeDocument(filePath, documentType, analysis)
	case "PDF":
		fa.analyzePdfDocument(filePath, analysis)
	case "RTF":
		fa.analyzeRtfDocument(filePath, analysis)
	case "TXT", "HTML", "XML", "JSON", "CSV", "Markdown":
		fa.analyzeTextDocument(filePath, documentType, analysis)
	default:
		fa.logger.Debugf("暂不支持分析文档类型: %s", documentType)
	}
}

func (fa *FileAnalyzer) analyzeOfficeDocument(filePath, documentType string, analysis *models.DocumentAnalysis) {
	// Office文档分析（基于ZIP格式）
	reader, err := zip.OpenReader(filePath)
	if err != nil {
		fa.logger.Errorf("打开Office文档失败: %v", err)
		return
	}
	defer reader.Close()

	// 检查是否包含宏
	hasMacros := false
	macroCount := 0
	embeddedObjects := []string{}

	for _, file := range reader.File {
		fileName := strings.ToLower(file.Name)

		// 检查宏文件
		if strings.Contains(fileName, "vba") || strings.Contains(fileName, "macro") {
			hasMacros = true
			macroCount++
		}

		// 检查嵌入对象
		if strings.Contains(fileName, "media") || strings.Contains(fileName, "image") {
			embeddedObjects = append(embeddedObjects, file.Name)
		}

		// 提取元数据
		if fileName == "docprops/core.xml" {
			fa.extractOfficeMetadata(file, analysis)
		}
	}

	analysis.HasMacros = hasMacros
	analysis.MacroCount = macroCount
	analysis.EmbeddedObjects = embeddedObjects

	fa.logger.Debugf("Office文档分析完成: %s, 宏数量: %d, 嵌入对象: %d",
		documentType, macroCount, len(embeddedObjects))
}

func (fa *FileAnalyzer) extractOfficeMetadata(file *zip.File, analysis *models.DocumentAnalysis) {
	// 提取Office文档元数据（简化实现）
	analysis.Metadata["creator"] = "Unknown"
	analysis.Metadata["created"] = "Unknown"
	analysis.Metadata["modified"] = "Unknown"
	analysis.Metadata["title"] = "Unknown"
	analysis.Metadata["subject"] = "Unknown"
}

func (fa *FileAnalyzer) analyzePdfDocument(filePath string, analysis *models.DocumentAnalysis) {
	// PDF文档分析（简化实现）
	// 在实际实现中，需要使用专门的PDF库
	analysis.HasMacros = false // PDF通常不包含宏
	analysis.MacroCount = 0
	analysis.EmbeddedObjects = []string{}

	// 提取基本元数据
	analysis.Metadata["format"] = "PDF"
	analysis.Metadata["version"] = "Unknown"

	fa.logger.Debugf("PDF文档分析完成（简化实现）")
}

func (fa *FileAnalyzer) analyzeRtfDocument(filePath string, analysis *models.DocumentAnalysis) {
	// RTF文档分析（简化实现）
	analysis.HasMacros = false // RTF通常不包含宏
	analysis.MacroCount = 0
	analysis.EmbeddedObjects = []string{}

	// 提取基本元数据
	analysis.Metadata["format"] = "RTF"

	fa.logger.Debugf("RTF文档分析完成（简化实现）")
}

func (fa *FileAnalyzer) analyzeTextDocument(filePath, documentType string, analysis *models.DocumentAnalysis) {
	// 文本文档分析
	content, err := fa.readFileContent(filePath)
	if err != nil {
		fa.logger.Errorf("读取文本文档失败: %v", err)
		return
	}

	analysis.HasMacros = false
	analysis.MacroCount = 0
	analysis.EmbeddedObjects = []string{}

	// 提取基本元数据
	analysis.Metadata["format"] = documentType
	analysis.Metadata["size"] = fmt.Sprintf("%d bytes", len(content))
	analysis.Metadata["lines"] = fmt.Sprintf("%d", strings.Count(string(content), "\n")+1)

	fa.logger.Debugf("文本文档分析完成: %s, 大小: %d bytes", documentType, len(content))
}

func (fa *FileAnalyzer) analyzeFileBehavior(filePath string) *models.BehaviorAnalysis {
	return &models.BehaviorAnalysis{
		Timestamp: time.Now(),
	}
}

func (fa *FileAnalyzer) analyzeFileSignature(filePath string) *models.SignatureAnalysis {
	return &models.SignatureAnalysis{
		Timestamp: time.Now(),
	}
}

func (fa *FileAnalyzer) analyzeFileEntropy(filePath string) *models.EntropyAnalysis {
	return &models.EntropyAnalysis{
		Timestamp: time.Now(),
	}
}

// Manager 文件管理器（兼容性类型）
type Manager struct {
	analyzer *FileAnalyzer
	logger   *logrus.Logger
}

// NewManager 创建文件管理器
func NewManager(logger *logrus.Logger) *Manager {
	return &Manager{
		analyzer: NewFileAnalyzer(logger, nil),
		logger:   logger,
	}
}

// CopyFile 复制文件
func (m *Manager) CopyFile(source, dest string) error {
	// 读取源文件
	sourceFile, err := os.Open(source)
	if err != nil {
		return fmt.Errorf("打开源文件失败: %w", err)
	}
	defer sourceFile.Close()

	// 创建目标文件
	destFile, err := os.Create(dest)
	if err != nil {
		return fmt.Errorf("创建目标文件失败: %w", err)
	}
	defer destFile.Close()

	// 复制文件内容
	_, err = io.Copy(destFile, sourceFile)
	if err != nil {
		return fmt.Errorf("复制文件内容失败: %w", err)
	}

	m.logger.Infof("文件复制成功: %s -> %s", source, dest)
	return nil
}

// MoveFile 移动文件
func (m *Manager) MoveFile(source, dest string) error {
	// 先复制文件
	err := m.CopyFile(source, dest)
	if err != nil {
		return err
	}

	// 删除源文件
	err = os.Remove(source)
	if err != nil {
		return fmt.Errorf("删除源文件失败: %w", err)
	}

	m.logger.Infof("文件移动成功: %s -> %s", source, dest)
	return nil
}

// DeleteFile 删除文件
func (m *Manager) DeleteFile(filePath string) error {
	err := os.Remove(filePath)
	if err != nil {
		return fmt.Errorf("删除文件失败: %w", err)
	}

	m.logger.Infof("文件删除成功: %s", filePath)
	return nil
}

// GetFileHash 获取文件哈希
func (m *Manager) GetFileHash(filePath, algorithm string) (string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return "", fmt.Errorf("打开文件失败: %w", err)
	}
	defer file.Close()

	var hash string
	switch strings.ToLower(algorithm) {
	case "md5":
		hashObj := md5.New()
		_, err = io.Copy(hashObj, file)
		if err != nil {
			return "", fmt.Errorf("计算MD5哈希失败: %w", err)
		}
		hash = hex.EncodeToString(hashObj.Sum(nil))
	case "sha256":
		hashObj := sha256.New()
		_, err = io.Copy(hashObj, file)
		if err != nil {
			return "", fmt.Errorf("计算SHA256哈希失败: %w", err)
		}
		hash = hex.EncodeToString(hashObj.Sum(nil))
	default:
		return "", fmt.Errorf("不支持的哈希算法: %s", algorithm)
	}

	return hash, nil
}

// GetFileHashes 获取文件的所有哈希值
func (m *Manager) GetFileHashes(filePath string) (map[string]string, error) {
	hashes := make(map[string]string)

	// 计算MD5哈希
	md5Hash, err := m.GetFileHash(filePath, "md5")
	if err == nil {
		hashes["md5"] = md5Hash
	}

	// 计算SHA256哈希
	sha256Hash, err := m.GetFileHash(filePath, "sha256")
	if err == nil {
		hashes["sha256"] = sha256Hash
	}

	return hashes, nil
}

// VerifyFileHash 验证文件哈希
func (m *Manager) VerifyFileHash(filePath, algorithm, expectedHash string) (bool, error) {
	actualHash, err := m.GetFileHash(filePath, algorithm)
	if err != nil {
		return false, err
	}

	return strings.EqualFold(actualHash, expectedHash), nil
}
