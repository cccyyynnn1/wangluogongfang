package middleware

import (
	"context"
	"fmt"
	"net/http"
	"runtime/debug"
	"strings"
	"time"

	"yara-security-service/internal/models"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// ErrorHandler 错误处理中间件
type ErrorHandler struct {
	logger *logrus.Logger
}

// NewErrorHandler 创建错误处理器
func NewErrorHandler(logger *logrus.Logger) *ErrorHandler {
	return &ErrorHandler{
		logger: logger,
	}
}

// Recovery 恢复中间件
func (h *ErrorHandler) Recovery() gin.HandlerFunc {
	return gin.CustomRecovery(func(c *gin.Context, recovered interface{}) {
		if err, ok := recovered.(string); ok {
			h.logger.Errorf("Panic recovered: %s", err)
		} else {
			h.logger.Errorf("Panic recovered: %v", recovered)
		}

		// 记录堆栈信息
		stack := debug.Stack()
		h.logger.Errorf("Stack trace: %s", string(stack))

		// 返回错误响应
		c.JSON(http.StatusInternalServerError, models.Response{
			Code:    500,
			Message: "服务器内部错误",
			Time:    time.Now(),
		})
	})
}

// ErrorResponse 错误响应中间件
func (h *ErrorHandler) ErrorResponse() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 设置错误处理函数
		c.Set("error_handler", func(err error, statusCode int) {
			h.handleError(c, err, statusCode)
		})

		c.Next()
	}
}

// handleError 处理错误
func (h *ErrorHandler) handleError(c *gin.Context, err error, statusCode int) {
	if err == nil {
		return
	}

	// 记录错误
	h.logger.WithFields(logrus.Fields{
		"method":     c.Request.Method,
		"path":       c.Request.URL.Path,
		"status":     statusCode,
		"client_ip":  c.ClientIP(),
		"user_agent": c.Request.UserAgent(),
		"error":      err.Error(),
	}).Error("请求处理错误")

	// 根据错误类型返回不同的响应
	var message string
	switch statusCode {
	case 400:
		message = "请求参数错误"
	case 401:
		message = "未授权访问"
	case 403:
		message = "禁止访问"
	case 404:
		message = "资源不存在"
	case 405:
		message = "方法不允许"
	case 408:
		message = "请求超时"
	case 429:
		message = "请求过于频繁"
	case 500:
		message = "服务器内部错误"
	case 502:
		message = "网关错误"
	case 503:
		message = "服务不可用"
	case 504:
		message = "网关超时"
	default:
		message = "未知错误"
	}

	c.JSON(statusCode, models.Response{
		Code:    statusCode,
		Message: message,
		Time:    time.Now(),
	})
}

// Timeout 超时中间件
func (h *ErrorHandler) Timeout(timeout time.Duration) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(c.Request.Context(), timeout)
		defer cancel()

		c.Request = c.Request.WithContext(ctx)

		// 创建完成通道
		done := make(chan bool, 1)
		go func() {
			c.Next()
			done <- true
		}()

		select {
		case <-done:
			// 请求正常完成
		case <-ctx.Done():
			// 请求超时
			h.logger.WithFields(logrus.Fields{
				"method": c.Request.Method,
				"path":   c.Request.URL.Path,
				"timeout": timeout,
			}).Warn("请求超时")

			c.AbortWithStatusJSON(http.StatusRequestTimeout, models.Response{
				Code:    408,
				Message: "请求超时",
				Time:    time.Now(),
			})
		}
	}
}

// RateLimit 速率限制中间件
func (h *ErrorHandler) RateLimit(requests int, window time.Duration) gin.HandlerFunc {
	// 简单的内存速率限制器
	clients := make(map[string]*rateLimitInfo)
	
	return func(c *gin.Context) {
		clientIP := c.ClientIP()
		now := time.Now()

		// 获取或创建客户端信息
		info, exists := clients[clientIP]
		if !exists {
			info = &rateLimitInfo{
				requests: make([]time.Time, 0),
				window:   window,
			}
			clients[clientIP] = info
		}

		// 清理过期的请求记录
		info.cleanup(now)

		// 检查是否超过限制
		if len(info.requests) >= requests {
			h.logger.WithFields(logrus.Fields{
				"client_ip": clientIP,
				"method":    c.Request.Method,
				"path":      c.Request.URL.Path,
			}).Warn("请求频率过高")

			c.JSON(http.StatusTooManyRequests, models.Response{
				Code:    429,
				Message: "请求过于频繁，请稍后再试",
				Time:    time.Now(),
			})
			c.Abort()
			return
		}

		// 记录当前请求
		info.requests = append(info.requests, now)
		c.Next()
	}
}

// rateLimitInfo 速率限制信息
type rateLimitInfo struct {
	requests []time.Time
	window   time.Duration
}

// cleanup 清理过期的请求记录
func (r *rateLimitInfo) cleanup(now time.Time) {
	cutoff := now.Add(-r.window)
	
	// 移除过期的请求记录
	valid := 0
	for _, req := range r.requests {
		if req.After(cutoff) {
			r.requests[valid] = req
			valid++
		}
	}
	r.requests = r.requests[:valid]
}

// Validation 输入验证中间件
func (h *ErrorHandler) Validation() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 检查请求大小
		if c.Request.ContentLength > 10*1024*1024 { // 10MB
			h.logger.WithFields(logrus.Fields{
				"method":         c.Request.Method,
				"path":           c.Request.URL.Path,
				"content_length": c.Request.ContentLength,
			}).Warn("请求体过大")

			c.JSON(http.StatusRequestEntityTooLarge, models.Response{
				Code:    413,
				Message: "请求体过大",
				Time:    time.Now(),
			})
			c.Abort()
			return
		}

		// 检查路径注入
		path := c.Request.URL.Path
		if containsPathTraversal(path) {
			h.logger.WithFields(logrus.Fields{
				"method": c.Request.Method,
				"path":   path,
				"client_ip": c.ClientIP(),
			}).Warn("检测到路径遍历攻击")

			c.JSON(http.StatusBadRequest, models.Response{
				Code:    400,
				Message: "无效的请求路径",
				Time:    time.Now(),
			})
			c.Abort()
			return
		}

		c.Next()
	}
}

// containsPathTraversal 检查是否包含路径遍历
func containsPathTraversal(path string) bool {
	dangerousPatterns := []string{
		"..", "//", "\\", "~", "..\\", "../",
	}
	
	pathLower := strings.ToLower(path)
	for _, pattern := range dangerousPatterns {
		if strings.Contains(pathLower, pattern) {
			return true
		}
	}
	return false
}

// Logging 日志中间件
func (h *ErrorHandler) Logging() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		raw := c.Request.URL.RawQuery

		// 处理请求
		c.Next()

		// 计算处理时间
		latency := time.Since(start)
		statusCode := c.Writer.Status()
		clientIP := c.ClientIP()
		method := c.Request.Method
		userAgent := c.Request.UserAgent()

		// 记录请求日志
		fields := logrus.Fields{
			"method":     method,
			"path":       path,
			"raw_query":  raw,
			"status":     statusCode,
			"latency":    latency,
			"client_ip":  clientIP,
			"user_agent": userAgent,
		}

		// 根据状态码选择日志级别
		switch {
		case statusCode >= 500:
			h.logger.WithFields(fields).Error("服务器错误")
		case statusCode >= 400:
			h.logger.WithFields(fields).Warn("客户端错误")
		case statusCode >= 300:
			h.logger.WithFields(fields).Info("重定向")
		default:
			h.logger.WithFields(fields).Info("请求完成")
		}
	}
}

// Security 安全中间件
func (h *ErrorHandler) Security() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 设置安全头
		c.Header("X-Content-Type-Options", "nosniff")
		c.Header("X-Frame-Options", "DENY")
		c.Header("X-XSS-Protection", "1; mode=block")
		c.Header("Referrer-Policy", "strict-origin-when-cross-origin")
		c.Header("Content-Security-Policy", "default-src 'self'")

		// 检查User-Agent
		userAgent := c.Request.UserAgent()
		if userAgent == "" || len(userAgent) > 500 {
			h.logger.WithFields(logrus.Fields{
				"client_ip":  c.ClientIP(),
				"user_agent": userAgent,
			}).Warn("可疑的User-Agent")

			c.JSON(http.StatusBadRequest, models.Response{
				Code:    400,
				Message: "无效的请求头",
				Time:    time.Now(),
			})
			c.Abort()
			return
		}

		c.Next()
	}
}

// HealthCheck 健康检查中间件
func (h *ErrorHandler) HealthCheck() gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.Request.URL.Path == "/health" {
			c.JSON(http.StatusOK, models.Response{
				Code:    200,
				Message: "服务正常",
				Time:    time.Now(),
			})
			c.Abort()
			return
		}
		c.Next()
	}
}

// ErrorResponse 错误响应结构
type ErrorResponse struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Details interface{} `json:"details,omitempty"`
	Time    time.Time   `json:"time"`
}

// NewErrorResponse 创建错误响应
func NewErrorResponse(code int, message string, details interface{}) ErrorResponse {
	return ErrorResponse{
		Code:    code,
		Message: message,
		Details: details,
		Time:    time.Now(),
	}
}

// HandlePanic 处理panic
func (h *ErrorHandler) HandlePanic(c *gin.Context) {
	if r := recover(); r != nil {
		h.logger.WithFields(logrus.Fields{
			"method": c.Request.Method,
			"path":   c.Request.URL.Path,
			"panic":  fmt.Sprintf("%v", r),
		}).Error("处理请求时发生panic")

		c.JSON(http.StatusInternalServerError, models.Response{
			Code:    500,
			Message: "服务器内部错误",
			Time:    time.Now(),
		})
	}
} 