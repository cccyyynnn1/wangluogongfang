package utils

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// ResponseData 响应数据结构
type ResponseData struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Time    string      `json:"time"`
	Data    interface{} `json:"data"`
}

// SuccessResponse 成功响应
func SuccessResponse(c *gin.Context, message string, data interface{}) {
	c.JSON(http.StatusOK, ResponseData{
		Code:    200,
		Message: message,
		Time:    time.Now().Format(time.RFC3339),
		Data:    data,
	})
}

// ErrorResponse 错误响应
func ErrorResponse(c *gin.Context, statusCode int, message string) {
	c.JSON(statusCode, ResponseData{
		Code:    statusCode,
		Message: message,
		Time:    time.Now().Format(time.RFC3339),
		Data:    gin.H{},
	})
}

// BadRequestResponse 400错误响应
func BadRequestResponse(c *gin.Context, message string) {
	ErrorResponse(c, http.StatusBadRequest, message)
}

// NotFoundResponse 404错误响应
func NotFoundResponse(c *gin.Context, message string) {
	ErrorResponse(c, http.StatusNotFound, message)
}

// InternalServerErrorResponse 500错误响应
func InternalServerErrorResponse(c *gin.Context, message string) {
	ErrorResponse(c, http.StatusInternalServerError, message)
}

// UnauthorizedResponse 401错误响应
func UnauthorizedResponse(c *gin.Context, message string) {
	ErrorResponse(c, http.StatusUnauthorized, message)
}

// ForbiddenResponse 403错误响应
func ForbiddenResponse(c *gin.Context, message string) {
	ErrorResponse(c, http.StatusForbidden, message)
}

// ValidationErrorResponse 验证错误响应
func ValidationErrorResponse(c *gin.Context, message string) {
	BadRequestResponse(c, "请求参数错误: "+message)
}

// FileNotFoundResponse 文件不存在响应
func FileNotFoundResponse(c *gin.Context) {
	NotFoundResponse(c, "文件不存在")
}

// DirectoryNotFoundResponse 目录不存在响应
func DirectoryNotFoundResponse(c *gin.Context) {
	NotFoundResponse(c, "目录不存在")
}

// ProcessNotFoundResponse 进程不存在响应
func ProcessNotFoundResponse(c *gin.Context) {
	NotFoundResponse(c, "进程不存在")
}

// UserNotFoundResponse 用户不存在响应
func UserNotFoundResponse(c *gin.Context) {
	NotFoundResponse(c, "用户不存在")
}

// RegistryKeyNotFoundResponse 注册表键不存在响应
func RegistryKeyNotFoundResponse(c *gin.Context) {
	NotFoundResponse(c, "注册表键不存在")
}

// NetworkConnectionNotFoundResponse 网络连接不存在响应
func NetworkConnectionNotFoundResponse(c *gin.Context) {
	NotFoundResponse(c, "网络连接不存在")
}

// ScanSuccessResponse 扫描成功响应
func ScanSuccessResponse(c *gin.Context, data interface{}) {
	SuccessResponse(c, "扫描完成", data)
}

// ScanFailedResponse 扫描失败响应
func ScanFailedResponse(c *gin.Context, err error) {
	InternalServerErrorResponse(c, "扫描失败: "+err.Error())
}

// OperationSuccessResponse 操作成功响应
func OperationSuccessResponse(c *gin.Context, operation string) {
	SuccessResponse(c, operation+"成功", gin.H{})
}

// OperationFailedResponse 操作失败响应
func OperationFailedResponse(c *gin.Context, operation string, err error) {
	InternalServerErrorResponse(c, operation+"失败: "+err.Error())
}

// ListSuccessResponse 列表获取成功响应
func ListSuccessResponse(c *gin.Context, data interface{}, count int) {
	SuccessResponse(c, "获取列表成功", gin.H{
		"data":  data,
		"count": count,
	})
}

// InfoSuccessResponse 信息获取成功响应
func InfoSuccessResponse(c *gin.Context, data interface{}) {
	SuccessResponse(c, "获取信息成功", data)
}

// StatusSuccessResponse 状态获取成功响应
func StatusSuccessResponse(c *gin.Context, data interface{}) {
	SuccessResponse(c, "获取状态成功", data)
}

// HealthCheckResponse 健康检查响应
func HealthCheckResponse(c *gin.Context) {
	c.JSON(http.StatusOK, ResponseData{
		Code:    200,
		Message: "服务正常",
		Time:    time.Now().Format(time.RFC3339),
		Data: gin.H{
			"status":  "healthy",
			"version": "1.0.0",
		},
	})
}

// MetricsResponse 指标响应
func MetricsResponse(c *gin.Context, metrics interface{}) {
	c.JSON(http.StatusOK, ResponseData{
		Code:    200,
		Message: "获取指标成功",
		Time:    time.Now().Format(time.RFC3339),
		Data:    metrics,
	})
}
