package handlers

import (
	"strings"
	"yara-security-service/internal/models"
	"yara-security-service/internal/services"
	"yara-security-service/internal/utils"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// RegistryHandler 注册表处理器
type RegistryHandler struct {
	registryService *services.RegistryService
	logger          *logrus.Logger
}

// NewRegistryHandler 创建注册表处理器
func NewRegistryHandler(registryService *services.RegistryService, logger *logrus.Logger) *RegistryHandler {
	return &RegistryHandler{
		registryService: registryService,
		logger:          logger,
	}
}

// GetRegistryKey 获取注册表键
func (h *RegistryHandler) GetRegistryKey(c *gin.Context) {
	keyPath := c.Param("path")
	if keyPath == "" {
		utils.BadRequestResponse(c, "注册表路径不能为空")
		return
	}

	// 如果路径参数不包含完整路径，尝试从查询参数获取完整路径
	if !strings.Contains(keyPath, "/") && !strings.Contains(keyPath, "\\") {
		fullPath := c.Query("fullPath")
		if fullPath != "" {
			keyPath = fullPath
		}
	}

	key, err := h.registryService.GetRegistryKey(keyPath)
	if err != nil {
		h.logger.Errorf("获取注册表键失败: %v", err)
		utils.InternalServerErrorResponse(c, "获取注册表键失败")
		return
	}

	utils.InfoSuccessResponse(c, key)
}

// CreateRegistryKey 创建注册表键
func (h *RegistryHandler) CreateRegistryKey(c *gin.Context) {
	var request models.CreateRegistryKeyRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		utils.ValidationErrorResponse(c, err.Error())
		return
	}

	err := h.registryService.CreateRegistryKey(request.Path)
	if err != nil {
		h.logger.Errorf("创建注册表键失败: %v", err)
		utils.OperationFailedResponse(c, "创建注册表键", err)
		return
	}

	utils.OperationSuccessResponse(c, "创建注册表键")
}

// DeleteRegistryKey 删除注册表键
func (h *RegistryHandler) DeleteRegistryKey(c *gin.Context) {
	keyPath := c.Param("path")
	if keyPath == "" {
		utils.BadRequestResponse(c, "注册表路径不能为空")
		return
	}

	// 如果路径参数不包含完整路径，尝试从查询参数获取完整路径
	if !strings.Contains(keyPath, "/") && !strings.Contains(keyPath, "\\") {
		fullPath := c.Query("fullPath")
		if fullPath != "" {
			keyPath = fullPath
		}
	}

	err := h.registryService.DeleteRegistryKey(keyPath)
	if err != nil {
		h.logger.Errorf("删除注册表键失败: %v", err)
		utils.OperationFailedResponse(c, "删除注册表键", err)
		return
	}

	utils.OperationSuccessResponse(c, "删除注册表键")
}

// SetRegistryValue 设置注册表值
func (h *RegistryHandler) SetRegistryValue(c *gin.Context) {
	var request models.SetRegistryValueRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		utils.ValidationErrorResponse(c, err.Error())
		return
	}

	// 将interface{}转换为string
	valueStr := ""
	if str, ok := request.Value.(string); ok {
		valueStr = str
	} else {
		utils.BadRequestResponse(c, "值必须是字符串类型")
		return
	}

	// 添加调试日志
	h.logger.Infof("设置注册表值 - Path: %s, Name: %s, Type: %s, Value: %s",
		request.Path, request.Name, request.Type, valueStr)

	err := h.registryService.SetRegistryValue(request.Path, request.Name, request.Type, valueStr)
	if err != nil {
		h.logger.Errorf("设置注册表值失败: %v", err)
		utils.OperationFailedResponse(c, "设置注册表值", err)
		return
	}

	utils.OperationSuccessResponse(c, "设置注册表值")
}

// GetRegistryValue 获取注册表值
func (h *RegistryHandler) GetRegistryValue(c *gin.Context) {
	keyPath := c.Param("path")
	valueName := c.Query("name")
	if keyPath == "" {
		utils.BadRequestResponse(c, "注册表路径不能为空")
		return
	}

	// 如果路径参数不包含完整路径，尝试从查询参数获取完整路径
	if !strings.Contains(keyPath, "/") && !strings.Contains(keyPath, "\\") {
		fullPath := c.Query("fullPath")
		if fullPath != "" {
			keyPath = fullPath
		}
	}

	value, valueType, err := h.registryService.GetRegistryValue(keyPath, valueName)
	if err != nil {
		h.logger.Errorf("获取注册表值失败: %v", err)
		utils.InternalServerErrorResponse(c, "获取注册表值失败")
		return
	}

	utils.InfoSuccessResponse(c, gin.H{
		"value": value,
		"type":  valueType,
	})
}

// GetRegistryValueByName 通过路径参数获取注册表值
func (h *RegistryHandler) GetRegistryValueByName(c *gin.Context) {
	keyPath := c.Param("path")
	valueName := c.Param("name")
	if keyPath == "" {
		utils.BadRequestResponse(c, "注册表路径不能为空")
		return
	}
	if valueName == "" {
		utils.BadRequestResponse(c, "值名称不能为空")
		return
	}

	// 如果路径参数不包含完整路径，尝试从查询参数获取完整路径
	if !strings.Contains(keyPath, "/") && !strings.Contains(keyPath, "\\") {
		fullPath := c.Query("fullPath")
		if fullPath != "" {
			keyPath = fullPath
		}
	}

	// 添加调试日志
	h.logger.Infof("通过路径参数获取注册表值 - Path: %s, Name: %s", keyPath, valueName)

	value, valueType, err := h.registryService.GetRegistryValue(keyPath, valueName)
	if err != nil {
		h.logger.Errorf("获取注册表值失败: %v", err)
		utils.InternalServerErrorResponse(c, "获取注册表值失败")
		return
	}

	utils.InfoSuccessResponse(c, gin.H{
		"value": value,
		"type":  valueType,
	})
}

// DeleteRegistryValue 删除注册表值
func (h *RegistryHandler) DeleteRegistryValue(c *gin.Context) {
	keyPath := c.Param("path")
	valueName := c.Query("name")
	if keyPath == "" {
		utils.BadRequestResponse(c, "注册表路径不能为空")
		return
	}

	// 如果路径参数不包含完整路径，尝试从查询参数获取完整路径
	if !strings.Contains(keyPath, "/") && !strings.Contains(keyPath, "\\") {
		fullPath := c.Query("fullPath")
		if fullPath != "" {
			keyPath = fullPath
		}
	}

	err := h.registryService.DeleteRegistryValue(keyPath, valueName)
	if err != nil {
		h.logger.Errorf("删除注册表值失败: %v", err)
		utils.OperationFailedResponse(c, "删除注册表值", err)
		return
	}

	utils.OperationSuccessResponse(c, "删除注册表值")
}

// DeleteRegistryValueByName 通过路径参数删除注册表值
func (h *RegistryHandler) DeleteRegistryValueByName(c *gin.Context) {
	keyPath := c.Param("path")
	valueName := c.Param("name")
	if keyPath == "" {
		utils.BadRequestResponse(c, "注册表路径不能为空")
		return
	}
	if valueName == "" {
		utils.BadRequestResponse(c, "值名称不能为空")
		return
	}

	// 如果路径参数不包含完整路径，尝试从查询参数获取完整路径
	if !strings.Contains(keyPath, "/") && !strings.Contains(keyPath, "\\") {
		fullPath := c.Query("fullPath")
		if fullPath != "" {
			keyPath = fullPath
		}
	}

	// 添加调试日志
	h.logger.Infof("通过路径参数删除注册表值 - Path: %s, Name: %s", keyPath, valueName)

	err := h.registryService.DeleteRegistryValue(keyPath, valueName)
	if err != nil {
		h.logger.Errorf("删除注册表值失败: %v", err)
		utils.OperationFailedResponse(c, "删除注册表值", err)
		return
	}

	utils.OperationSuccessResponse(c, "删除注册表值")
}

// ListRegistryKeys 列出注册表键
func (h *RegistryHandler) ListRegistryKeys(c *gin.Context) {
	keyPath := c.Param("path")
	if keyPath == "" {
		utils.BadRequestResponse(c, "注册表路径不能为空")
		return
	}

	// 如果路径参数不包含完整路径，尝试从查询参数获取完整路径
	if !strings.Contains(keyPath, "/") && !strings.Contains(keyPath, "\\") {
		fullPath := c.Query("fullPath")
		if fullPath != "" {
			keyPath = fullPath
		}
	}

	keys, err := h.registryService.ListRegistryKeys(keyPath)
	if err != nil {
		h.logger.Errorf("列出注册表键失败: %v", err)
		utils.InternalServerErrorResponse(c, "列出注册表键失败")
		return
	}

	utils.ListSuccessResponse(c, keys, len(keys))
}

// ListRegistryValues 列出注册表值
func (h *RegistryHandler) ListRegistryValues(c *gin.Context) {
	keyPath := c.Param("path")
	if keyPath == "" {
		utils.BadRequestResponse(c, "注册表路径不能为空")
		return
	}

	// 如果路径参数不包含完整路径，尝试从查询参数获取完整路径
	if !strings.Contains(keyPath, "/") && !strings.Contains(keyPath, "\\") {
		fullPath := c.Query("fullPath")
		if fullPath != "" {
			keyPath = fullPath
		}
	}

	values, err := h.registryService.ListRegistryValues(keyPath)
	if err != nil {
		h.logger.Errorf("列出注册表值失败: %v", err)
		utils.InternalServerErrorResponse(c, "列出注册表值失败")
		return
	}

	utils.ListSuccessResponse(c, values, len(values))
}

// SearchRegistry 搜索注册表
func (h *RegistryHandler) SearchRegistry(c *gin.Context) {
	var request models.SearchRegistryRequest

	// 检查请求方法，支持GET和POST
	if c.Request.Method == "GET" {
		// GET请求从查询参数获取数据
		request.Root = c.Query("root_path")
		request.Pattern = c.Query("search_term")
		request.ValuePattern = c.Query("value_pattern")

		if request.Root == "" {
			utils.BadRequestResponse(c, "root_path参数不能为空")
			return
		}
	} else {
		// POST请求从JSON获取数据
		if err := c.ShouldBindJSON(&request); err != nil {
			utils.ValidationErrorResponse(c, err.Error())
			return
		}
	}

	// 添加调试日志
	h.logger.Infof("搜索注册表 - Root: %s, Pattern: %s, ValuePattern: %s",
		request.Root, request.Pattern, request.ValuePattern)

	results, err := h.registryService.SearchRegistry(request.Root, request.Pattern)
	if err != nil {
		h.logger.Errorf("搜索注册表失败: %v", err)
		utils.InternalServerErrorResponse(c, "搜索注册表失败")
		return
	}

	utils.ListSuccessResponse(c, results, len(results))
}
