package middleware

import (
	"bytes"
	"io"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// LoggingConfig 日志配置
type LoggingConfig struct {
	Enabled     bool
	LogHeaders  bool
	LogBody     bool
	LogResponse bool
}

// DefaultLoggingConfig 默认日志配置
func DefaultLoggingConfig() *LoggingConfig {
	return &LoggingConfig{
		Enabled:     true,
		LogHeaders:  false,
		LogBody:     false,
		LogResponse: false,
	}
}

// Logging 日志中间件
func Logging(config *LoggingConfig, logger *logrus.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		if !config.Enabled {
			c.Next()
			return
		}

		start := time.Now()
		path := c.Request.URL.Path
		raw := c.Request.URL.RawQuery

		// 记录请求开始
		logger.WithFields(logrus.Fields{
			"method":     c.Request.Method,
			"path":       path,
			"raw_query":  raw,
			"client_ip":  c.ClientIP(),
			"user_agent": c.Request.UserAgent(),
		}).Info("请求开始")

		// 记录请求头（如果启用）
		if config.LogHeaders {
			logger.WithFields(logrus.Fields{
				"headers": c.Request.Header,
			}).Debug("请求头")
		}

		// 记录请求体（如果启用）
		if config.LogBody && c.Request.Body != nil {
			bodyBytes, _ := io.ReadAll(c.Request.Body)
			c.Request.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))

			logger.WithFields(logrus.Fields{
				"body": string(bodyBytes),
			}).Debug("请求体")
		}

		// 包装响应写入器以捕获响应
		var responseBody bytes.Buffer
		if config.LogResponse {
			c.Writer = &responseWriter{
				ResponseWriter: c.Writer,
				body:           &responseBody,
			}
		}

		// 处理请求
		c.Next()

		// 计算处理时间
		latency := time.Since(start)

		// 记录响应
		logger.WithFields(logrus.Fields{
			"status_code": c.Writer.Status(),
			"latency":     latency,
			"method":      c.Request.Method,
			"path":        path,
		}).Info("请求完成")

		// 记录响应体（如果启用）
		if config.LogResponse {
			logger.WithFields(logrus.Fields{
				"response_body": responseBody.String(),
			}).Debug("响应体")
		}

		// 记录错误（如果有）
		if len(c.Errors) > 0 {
			logger.WithFields(logrus.Fields{
				"errors": c.Errors.String(),
			}).Error("请求处理错误")
		}
	}
}

// responseWriter 响应写入器包装器
type responseWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

// Write 写入响应
func (w *responseWriter) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}

// WriteString 写入字符串响应
func (w *responseWriter) WriteString(s string) (int, error) {
	w.body.WriteString(s)
	return w.ResponseWriter.WriteString(s)
}

// Recovery 恢复中间件
func Recovery(logger *logrus.Logger) gin.HandlerFunc {
	return gin.CustomRecovery(func(c *gin.Context, recovered interface{}) {
		if err, ok := recovered.(string); ok {
			logger.WithFields(logrus.Fields{
				"error": err,
				"path":  c.Request.URL.Path,
				"ip":    c.ClientIP(),
			}).Error("请求处理panic")
		}

		c.JSON(500, gin.H{
			"code":    500,
			"message": "服务器内部错误",
			"data":    nil,
		})
	})
}

// RequestID 请求ID中间件
func RequestID() gin.HandlerFunc {
	return func(c *gin.Context) {
		requestID := c.GetHeader("X-Request-ID")
		if requestID == "" {
			requestID = generateRequestID()
		}

		c.Header("X-Request-ID", requestID)
		c.Set("request_id", requestID)
		c.Next()
	}
}

// generateRequestID 生成请求ID
func generateRequestID() string {
	// 这里可以使用UUID或其他方式生成请求ID
	return "req_" + time.Now().Format("20060102150405")
}
