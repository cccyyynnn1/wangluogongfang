package handlers

import (
	"context"
	"crypto/md5"
	"crypto/sha256"
	"fmt"
	"net/http"
	"os"
	"os/user"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"yara-security-service/internal/models"
	"yara-security-service/internal/services"
	"yara-security-service/internal/utils"
	"yara-security-service/pkg/security"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// FileHandler 文件处理器
type FileHandler struct {
	securityService *services.SecurityService
	scanner         *security.Scanner
	logger          *logrus.Logger
}

// NewFileHandler 创建新的文件处理器
func NewFileHandler(securityService *services.SecurityService, scanner *security.Scanner, logger *logrus.Logger) *FileHandler {
	return &FileHandler{
		securityService: securityService,
		scanner:         scanner,
		logger:          logger,
	}
}

// ScanFile 扫描单个文件
func (h *FileHandler) ScanFile(c *gin.Context) {
	var request struct {
		Path string `json:"path" binding:"required"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		utils.ValidationErrorResponse(c, err.Error())
		return
	}

	filePath := request.Path
	if filePath == "" {
		utils.BadRequestResponse(c, "文件路径不能为空")
		return
	}

	// 检查文件是否存在
	if !h.fileExists(filePath) {
		utils.FileNotFoundResponse(c)
		return
	}

	// 执行扫描（使用SecurityService以更新统计信息）
	result, err := h.securityService.ScanFile(c.Request.Context(), filePath)
	if err != nil {
		h.logger.Errorf("扫描文件失败: %v", err)
		utils.ScanFailedResponse(c, err)
		return
	}

	utils.ScanSuccessResponse(c, result)
}

// GetFileInfo 获取文件信息
func (h *FileHandler) GetFileInfo(c *gin.Context) {
	// 使用新的 /:path 格式，支持查询参数传递完整路径
	localPath := c.Param("path")
	if localPath == "" {
		utils.BadRequestResponse(c, "文件路径不能为空")
		return
	}

	// 如果路径参数不包含完整路径，尝试从查询参数获取完整路径
	if !strings.Contains(localPath, "/") && !strings.Contains(localPath, "\\") {
		fullPath := c.Query("fullPath")
		if fullPath != "" {
			localPath = fullPath
		}
	}

	// 获取文件信息
	fileInfo, err := h.getFileInfo(localPath)
	if err != nil {
		h.logger.Errorf("获取文件信息失败: %v", err)
		utils.OperationFailedResponse(c, "获取文件信息", err)
		return
	}

	utils.InfoSuccessResponse(c, fileInfo)
}

// ScanDirectory 扫描目录
func (h *FileHandler) ScanDirectory(c *gin.Context) {
	var request models.DirectoryScanRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		utils.ValidationErrorResponse(c, err.Error())
		return
	}

	if request.Path == "" {
		utils.BadRequestResponse(c, "目录路径不能为空")
		return
	}

	// 检查目录是否存在
	if !h.directoryExists(request.Path) {
		utils.DirectoryNotFoundResponse(c)
		return
	}

	// 设置合理的默认值
	if request.MaxDepth <= 0 {
		request.MaxDepth = 10
	}

	// 设置默认的include_patterns（如果为空，则扫描所有支持的文件类型）
	if len(request.IncludePatterns) == 0 {
		request.IncludePatterns = []string{"*.exe", "*.dll", "*.sys", "*.bat", "*.cmd", "*.ps1", "*.vbs", "*.js", "*.jar", "*.msi", "*.scr", "*.com", "*.txt", "*.log"}
	}

	// 创建带超时的上下文
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Minute)
	defer cancel()

	// 执行目录扫描（使用SecurityService以更新统计信息）
	results, err := h.securityService.ScanDirectory(
		ctx,
		request.Path,
		request.Recursive,
		request.MaxDepth,
		request.IncludePatterns,
		request.ExcludePatterns,
	)
	if err != nil {
		h.logger.Errorf("扫描目录失败: %v", err)
		utils.ScanFailedResponse(c, err)
		return
	}

	// 限制返回结果数量，避免响应过大
	maxResults := 1000
	if len(results) > maxResults {
		h.logger.Warnf("目录扫描结果过多，限制为前%d个结果", maxResults)
		results = results[:maxResults]
	}

	c.JSON(http.StatusOK, models.Response{
		Code:    200,
		Message: "目录扫描完成",
		Data: map[string]interface{}{
			"directory":        request.Path,
			"results":          results,
			"count":            len(results),
			"total_scanned":    len(results),
			"scan_time":        time.Now().Format(time.RFC3339),
			"include_patterns": request.IncludePatterns,
			"exclude_patterns": request.ExcludePatterns,
			"recursive":        request.Recursive,
			"max_depth":        request.MaxDepth,
		},
		Time: time.Now(),
	})
}

// ScanBuffer 扫描内存缓冲区
func (h *FileHandler) ScanBuffer(c *gin.Context) {
	var request models.BufferScanRequest

	// 尝试从JSON获取数据
	if err := c.ShouldBindJSON(&request); err != nil {
		// JSON解析失败，尝试从表单获取数据
		identifier := c.PostForm("identifier")
		if identifier == "" {
			c.JSON(http.StatusBadRequest, models.Response{
				Code:    400,
				Message: "标识符不能为空",
				Time:    time.Now(),
			})
			return
		}

		// 获取上传的文件数据
		file, err := c.FormFile("data")
		if err != nil {
			c.JSON(http.StatusBadRequest, models.Response{
				Code:    400,
				Message: "未找到上传的数据",
				Time:    time.Now(),
			})
			return
		}

		// 读取文件数据
		openedFile, err := file.Open()
		if err != nil {
			c.JSON(http.StatusInternalServerError, models.Response{
				Code:    500,
				Message: "读取文件数据失败",
				Time:    time.Now(),
			})
			return
		}
		defer openedFile.Close()

		// 读取所有数据
		data := make([]byte, file.Size)
		_, err = openedFile.Read(data)
		if err != nil {
			c.JSON(http.StatusInternalServerError, models.Response{
				Code:    500,
				Message: "读取数据失败",
				Time:    time.Now(),
			})
			return
		}

		request.Data = string(data)
		request.Identifier = identifier
	}

	// 验证必要参数
	if request.Identifier == "" {
		c.JSON(http.StatusBadRequest, models.Response{
			Code:    400,
			Message: "标识符不能为空",
			Time:    time.Now(),
		})
		return
	}

	if request.Data == "" {
		c.JSON(http.StatusBadRequest, models.Response{
			Code:    400,
			Message: "数据不能为空",
			Time:    time.Now(),
		})
		return
	}

	// 执行扫描（使用SecurityService以更新统计信息）
	result, err := h.securityService.ScanBuffer(c.Request.Context(), []byte(request.Data), request.Identifier)
	if err != nil {
		h.logger.Errorf("扫描缓冲区失败: %v", err)
		utils.ScanFailedResponse(c, err)
		return
	}

	utils.ScanSuccessResponse(c, result)
}

// GetFileList 获取文件列表
func (h *FileHandler) GetFileList(c *gin.Context) {
	dirPath := c.Query("path")
	if dirPath == "" {
		dirPath = "."
	}

	// 获取递归和深度参数
	recursiveStr := c.DefaultQuery("recursive", "false")
	recursive := recursiveStr == "true"

	maxDepthStr := c.DefaultQuery("max_depth", "1")
	maxDepth := 1
	if maxDepthStr != "1" {
		if parsed, err := strconv.Atoi(maxDepthStr); err == nil && parsed > 0 {
			maxDepth = parsed
		}
	}

	h.logger.Debugf("获取文件列表: 路径=%s, 递归=%v, 最大深度=%d", dirPath, recursive, maxDepth)

	// 获取文件列表
	files, err := h.getFileListWithDepth(dirPath, recursive, maxDepth)
	if err != nil {
		h.logger.Errorf("获取文件列表失败: %v", err)
		c.JSON(http.StatusInternalServerError, models.Response{
			Code:    500,
			Message: "获取文件列表失败",
			Time:    time.Now(),
		})
		return
	}

	c.JSON(http.StatusOK, models.Response{
		Code:    200,
		Message: "获取文件列表成功",
		Data: map[string]interface{}{
			"directory": dirPath,
			"files":     files,
			"count":     len(files),
			"recursive": recursive,
			"max_depth": maxDepth,
		},
		Time: time.Now(),
	})
}

// fileExists 检查文件是否存在
func (h *FileHandler) fileExists(filePath string) bool {
	_, err := os.Stat(filePath)
	return err == nil
}

// directoryExists 检查目录是否存在
func (h *FileHandler) directoryExists(dirPath string) bool {
	info, err := os.Stat(dirPath)
	if err != nil {
		return false
	}
	return info.IsDir()
}

// getFileList 获取文件列表
func (h *FileHandler) getFileList(dirPath string) ([]*models.FileInfo, error) {
	var files []*models.FileInfo

	err := filepath.WalkDir(dirPath, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return nil
		}

		// 跳过根目录
		if path == dirPath {
			return nil
		}

		// 获取文件信息
		info, err := d.Info()
		if err != nil {
			return nil
		}

		fileInfo := &models.FileInfo{
			Path:        path,
			Name:        d.Name(),
			Size:        info.Size(),
			IsDir:       d.IsDir(),
			ModTime:     info.ModTime(),
			Permissions: info.Mode().String(),
		}

		files = append(files, fileInfo)
		return nil
	})

	if err != nil {
		return nil, err
	}

	return files, nil
}

// getFileListWithDepth 获取文件列表，支持递归和深度控制
func (h *FileHandler) getFileListWithDepth(dirPath string, recursive bool, maxDepth int) ([]*models.FileInfo, error) {
	var files []*models.FileInfo

	err := filepath.WalkDir(dirPath, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return nil
		}

		// 跳过根目录
		if path == dirPath {
			return nil
		}

		// 处理目录
		if d.IsDir() {
			// 检查深度限制
			if recursive && maxDepth > 0 {
				relPath, err := filepath.Rel(dirPath, path)
				if err != nil {
					h.logger.Warnf("计算相对路径失败 %s: %v", path, err)
					return filepath.SkipDir
				}

				// 计算深度：分割路径并计算层级
				pathParts := strings.Split(relPath, string(filepath.Separator))
				// 过滤掉空字符串（可能由连续分隔符产生）
				var cleanParts []string
				for _, part := range pathParts {
					if part != "" && part != "." && part != ".." {
						cleanParts = append(cleanParts, part)
					}
				}
				depth := len(cleanParts)

				h.logger.Debugf("目录 %s 的相对路径: %s, 深度: %d, 最大深度: %d", path, relPath, depth, maxDepth)

				if depth >= maxDepth {
					h.logger.Debugf("达到最大深度限制，跳过目录: %s (深度: %d >= %d)", path, depth, maxDepth)
					return filepath.SkipDir
				}
			}

			// 如果是递归扫描，继续进入子目录
			if recursive {
				h.logger.Debugf("进入子目录: %s", path)
				return nil
			} else {
				// 非递归扫描，跳过子目录
				h.logger.Debugf("非递归模式，跳过子目录: %s", path)
				return filepath.SkipDir
			}
		}

		// 只处理文件
		if !d.IsDir() {
			// 获取文件信息
			info, err := d.Info()
			if err != nil {
				return nil
			}

			fileInfo := &models.FileInfo{
				Path:        path,
				Name:        d.Name(),
				Size:        info.Size(),
				IsDir:       d.IsDir(),
				ModTime:     info.ModTime(),
				Permissions: info.Mode().String(),
			}

			files = append(files, fileInfo)
			h.logger.Debugf("添加文件: %s", path)
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return files, nil
}

// getFileInfo 获取文件信息
func (h *FileHandler) getFileInfo(filePath string) (*models.FileInfo, error) {
	// 获取文件信息
	info, err := os.Stat(filePath)
	if err != nil {
		return nil, fmt.Errorf("获取文件信息失败: %w", err)
	}

	// 计算文件哈希
	md5Hash, sha256Hash, err := h.calculateFileHashes(filePath)
	if err != nil {
		h.logger.Warnf("计算文件哈希失败: %v", err)
		md5Hash = ""
		sha256Hash = ""
	}

	// 获取真实的文件时间信息
	creationTime, accessTime := h.getRealFileTimes(filePath, info)

	fileInfo := &models.FileInfo{
		Path:         filePath,
		Name:         info.Name(),
		Size:         info.Size(),
		IsDir:        info.IsDir(),
		ModTime:      info.ModTime(),
		CreateTime:   creationTime,
		AccessTime:   accessTime,
		Permissions:  info.Mode().String(),
		Owner:        h.getFileOwner(info),
		Group:        h.getFileGroup(info),
		MD5:          md5Hash,
		SHA256:       sha256Hash,
		IsSuspicious: h.detectSuspiciousFile(filePath),
		ThreatLevel:  h.getThreatLevel(filePath),
	}

	return fileInfo, nil
}

// calculateFileHashes 计算文件哈希
func (h *FileHandler) calculateFileHashes(filePath string) (string, string, error) {
	// 读取文件内容
	data, err := os.ReadFile(filePath)
	if err != nil {
		return "", "", fmt.Errorf("读取文件失败: %w", err)
	}

	// 计算MD5哈希
	md5Hash := fmt.Sprintf("%x", md5.Sum(data))

	// 计算SHA256哈希
	sha256Hash := fmt.Sprintf("%x", sha256.Sum256(data))

	return md5Hash, sha256Hash, nil
}

// getFileCreateTime 获取文件创建时间
func (h *FileHandler) getFileCreateTime(info os.FileInfo) time.Time {
	// 使用新的文件时间工具获取真实的创建时间
	creationTime, err := security.GetFileCreationTime(info.Name())
	if err != nil {
		h.logger.Warnf("获取文件创建时间失败: %v，使用修改时间作为替代", err)
		return info.ModTime()
	}

	// 检查时间是否为零值（表示不可用）
	zeroTime := time.Time{}
	if creationTime.Equal(zeroTime) {
		h.logger.Debugf("文件创建时间不可用，使用修改时间作为替代: %s", info.Name())
		return info.ModTime()
	}

	return creationTime
}

// getRealFileTimes 获取真实的文件时间信息
func (h *FileHandler) getRealFileTimes(filePath string, info os.FileInfo) (time.Time, time.Time) {
	// 使用新的文件时间工具获取真实的创建时间和访问时间
	creationTime, err := security.GetFileCreationTime(filePath)
	if err != nil {
		h.logger.Warnf("获取文件创建时间失败: %v，使用修改时间作为替代", err)
		creationTime = info.ModTime()
	}

	// 使用新的函数获取访问时间，并强制更新
	accessTime, err := security.GetFileAccessTimeWithUpdate(filePath)
	if err != nil {
		h.logger.Warnf("获取文件访问时间失败: %v，使用修改时间作为替代", err)
		accessTime = info.ModTime()
	}

	// 检查时间是否为零值（表示不可用）
	zeroTime := time.Time{}
	if creationTime.Equal(zeroTime) {
		h.logger.Debugf("文件创建时间不可用，使用修改时间作为替代: %s", filePath)
		creationTime = info.ModTime()
	}

	if accessTime.Equal(zeroTime) {
		h.logger.Debugf("文件访问时间不可用，使用修改时间作为替代: %s", filePath)
		accessTime = info.ModTime()
	}

	return creationTime, accessTime
}

// getFileAccessTime 获取文件访问时间
func (h *FileHandler) getFileAccessTime(info os.FileInfo) time.Time {
	// 使用新的文件时间工具获取真实的访问时间
	accessTime, err := security.GetFileAccessTime(info.Name())
	if err != nil {
		h.logger.Warnf("获取文件访问时间失败: %v，使用修改时间作为替代", err)
		return info.ModTime()
	}

	// 检查时间是否为零值（表示不可用）
	zeroTime := time.Time{}
	if accessTime.Equal(zeroTime) {
		h.logger.Debugf("文件访问时间不可用，使用修改时间作为替代: %s", info.Name())
		return info.ModTime()
	}

	return accessTime
}

// getFileOwner 获取文件所有者
func (h *FileHandler) getFileOwner(info os.FileInfo) string {
	// 尝试获取文件所有者信息
	currentUser, err := user.Current()
	if err == nil {
		return currentUser.Username
	}
	return "current_user"
}

// getFileGroup 获取文件组
func (h *FileHandler) getFileGroup(info os.FileInfo) string {
	// 尝试获取文件组信息
	// 在Windows上，组信息获取较为复杂，使用默认组
	// 在实际项目中，可以使用Windows API获取更准确的组信息
	return "default_group"
}

// detectSuspiciousFile 检测可疑文件
func (h *FileHandler) detectSuspiciousFile(filePath string) bool {
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
func (h *FileHandler) getThreatLevel(filePath string) string {
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
