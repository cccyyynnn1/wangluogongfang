package handlers

import (
	"yara-security-service/internal/models"
	"yara-security-service/internal/services"
	"yara-security-service/internal/utils"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// SecurityHandler 安全处理器
type SecurityHandler struct {
	securityService *services.SecurityService
	logger          *logrus.Logger
}

// NewSecurityHandler 创建安全处理器
func NewSecurityHandler(securityService *services.SecurityService, logger *logrus.Logger) *SecurityHandler {
	return &SecurityHandler{
		securityService: securityService,
		logger:          logger,
	}
}

// ScanFile 扫描文件
func (h *SecurityHandler) ScanFile(c *gin.Context) {
	var request models.ScanFileRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		utils.ValidationErrorResponse(c, err.Error())
		return
	}

	result, err := h.securityService.ScanFile(c.Request.Context(), request.Path)
	if err != nil {
		h.logger.Errorf("扫描文件失败: %v", err)
		utils.ScanFailedResponse(c, err)
		return
	}

	utils.ScanSuccessResponse(c, result)
}

// GetYaraRules 获取Yara规则
func (h *SecurityHandler) GetYaraRules(c *gin.Context) {
	rules := h.securityService.GetYaraRules()

	// 构建响应数据
	responseData := gin.H{
		"rules":       rules,
		"total_count": len(rules),
	}

	utils.InfoSuccessResponse(c, responseData)
}

// GetScanStatus 获取扫描状态
func (h *SecurityHandler) GetScanStatus(c *gin.Context) {
	status := h.securityService.GetSecurityStatus()

	utils.StatusSuccessResponse(c, status)
}

// GetRulesInfo 获取规则信息
func (h *SecurityHandler) GetRulesInfo(c *gin.Context) {
	info := h.securityService.GetRulesInfo()

	utils.InfoSuccessResponse(c, info)
}

// ReloadRules 重新加载规则
func (h *SecurityHandler) ReloadRules(c *gin.Context) {
	err := h.securityService.ReloadRules()
	if err != nil {
		h.logger.Errorf("重新加载规则失败: %v", err)
		utils.OperationFailedResponse(c, "重新加载规则", err)
		return
	}

	utils.OperationSuccessResponse(c, "重新加载规则任务已提交")
}

// ReloadStatus 获取规则重载状态
func (h *SecurityHandler) ReloadStatus(c *gin.Context) {
	status := h.securityService.GetReloadStatus()
	utils.InfoSuccessResponse(c, status)
}

// ClearCache 清除缓存
func (h *SecurityHandler) ClearCache(c *gin.Context) {
	h.securityService.ClearCache()

	utils.OperationSuccessResponse(c, "清除缓存")
}

// GetCacheStats 获取缓存统计
func (h *SecurityHandler) GetCacheStats(c *gin.Context) {
	stats := h.securityService.GetCacheStats()

	utils.InfoSuccessResponse(c, stats)
}

// QuarantineFile 隔离文件
func (h *SecurityHandler) QuarantineFile(c *gin.Context) {
	var request models.QuarantineFileRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		utils.ValidationErrorResponse(c, err.Error())
		return
	}

	err := h.securityService.QuarantineFile(request.FilePath)
	if err != nil {
		h.logger.Errorf("隔离文件失败: %v", err)
		utils.OperationFailedResponse(c, "隔离文件", err)
		return
	}

	utils.OperationSuccessResponse(c, "隔离文件")
}

// RestoreFile 恢复文件
func (h *SecurityHandler) RestoreFile(c *gin.Context) {
	var request models.RestoreFileRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		utils.ValidationErrorResponse(c, err.Error())
		return
	}

	err := h.securityService.RestoreFile(request.FilePath)
	if err != nil {
		h.logger.Errorf("恢复文件失败: %v", err)
		utils.OperationFailedResponse(c, "恢复文件", err)
		return
	}

	utils.OperationSuccessResponse(c, "恢复文件")
}

// GetQuarantineList 获取隔离列表
func (h *SecurityHandler) GetQuarantineList(c *gin.Context) {
	list := h.securityService.GetQuarantineList()

	utils.ListSuccessResponse(c, list, len(list))
}

// GetScanHistory 获取扫描历史
func (h *SecurityHandler) GetScanHistory(c *gin.Context) {
	history := h.securityService.GetScanHistory(50) // 默认获取50条记录

	utils.ListSuccessResponse(c, history, len(history))
}
