package middleware

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// TimeoutConfig 超时配置
type TimeoutConfig struct {
	DefaultTimeout time.Duration
	MaxTimeout     time.Duration
	EnableDegrade  bool
	DegradeTimeout time.Duration
}

// DefaultTimeoutConfig 默认超时配置
func DefaultTimeoutConfig() *TimeoutConfig {
	return &TimeoutConfig{
		DefaultTimeout: 30 * time.Second,
		MaxTimeout:     60 * time.Second,
		EnableDegrade:  true,
		DegradeTimeout: 10 * time.Second,
	}
}

// Timeout 超时中间件
func Timeout(config *TimeoutConfig, logger *logrus.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 跳过需长时间处理的接口（规则热重载）
		if c.Request.Method == http.MethodPost && c.FullPath() == "/api/v1/security/reload-rules" {
			c.Next()
			return
		}

		// 从请求头获取超时时间
		timeoutStr := c.GetHeader("X-Timeout")
		var timeout time.Duration

		if timeoutStr != "" {
			if parsed, err := time.ParseDuration(timeoutStr); err == nil {
				timeout = parsed
				if timeout > config.MaxTimeout {
					timeout = config.MaxTimeout
				}
			} else {
				timeout = config.DefaultTimeout
			}
		} else {
			timeout = config.DefaultTimeout
		}

		// 创建带超时的上下文
		ctx, cancel := context.WithTimeout(c.Request.Context(), timeout)
		defer cancel()

		// 替换请求上下文
		c.Request = c.Request.WithContext(ctx)

		// 创建响应通道
		done := make(chan bool, 1)
		var err error

		// 在协程中处理请求
		go func() {
			defer func() {
				done <- true
			}()

			// 继续处理请求
			c.Next()

			// 检查是否有错误
			if len(c.Errors) > 0 {
				err = c.Errors.Last().Err
			}
		}()

		// 等待处理完成或超时
		select {
		case <-done:
			// 请求正常完成
			if err != nil {
				logger.Errorf("请求处理出错: %v", err)
			}
		case <-ctx.Done():
			// 请求超时
			logger.Warnf("请求超时: %s %s", c.Request.Method, c.Request.URL.Path)

			// 如果启用降级，使用更短的超时重试
			if config.EnableDegrade {
				c.JSON(http.StatusRequestTimeout, gin.H{
					"code":    408,
					"message": "请求超时，正在降级处理",
					"time":    time.Now(),
				})

				// 使用降级超时重试
				degradeCtx, degradeCancel := context.WithTimeout(context.Background(), config.DegradeTimeout)
				defer degradeCancel()

				// 这里可以实现降级逻辑，比如返回缓存数据或简化处理
				select {
				case <-degradeCtx.Done():
					c.JSON(http.StatusRequestTimeout, gin.H{
						"code":    408,
						"message": "降级处理也超时",
						"time":    time.Now(),
					})
				case <-time.After(config.DegradeTimeout):
					c.JSON(http.StatusRequestTimeout, gin.H{
						"code":    408,
						"message": "服务暂时不可用，请稍后重试",
						"time":    time.Now(),
					})
				}
			} else {
				c.JSON(http.StatusRequestTimeout, gin.H{
					"code":    408,
					"message": "请求超时",
					"time":    time.Now(),
				})
			}

			c.Abort()
		}
	}
}

// DegradeHandler 降级处理器
type DegradeHandler struct {
	logger *logrus.Logger
	cache  map[string]interface{}
}

// NewDegradeHandler 创建降级处理器
func NewDegradeHandler(logger *logrus.Logger) *DegradeHandler {
	return &DegradeHandler{
		logger: logger,
		cache:  make(map[string]interface{}),
	}
}

// HandleDegrade 处理降级逻辑
func (h *DegradeHandler) HandleDegrade(c *gin.Context, operation string) {
	// 根据操作类型返回降级响应
	switch operation {
	case "file_scan":
		c.JSON(http.StatusOK, gin.H{
			"code":    200,
			"message": "文件扫描服务暂时不可用，返回基础信息",
			"data": gin.H{
				"file_path": c.Param("path"),
				"status":    "degraded",
				"message":   "使用降级模式，仅提供基础文件信息",
			},
			"time": time.Now(),
		})
	case "process_list":
		c.JSON(http.StatusOK, gin.H{
			"code":    200,
			"message": "进程列表服务暂时不可用",
			"data": gin.H{
				"processes": []interface{}{},
				"count":     0,
				"status":    "degraded",
			},
			"time": time.Now(),
		})
	case "network_connections":
		c.JSON(http.StatusOK, gin.H{
			"code":    200,
			"message": "网络连接服务暂时不可用",
			"data": gin.H{
				"connections": []interface{}{},
				"count":       0,
				"status":      "degraded",
			},
			"time": time.Now(),
		})
	default:
		c.JSON(http.StatusServiceUnavailable, gin.H{
			"code":    503,
			"message": "服务暂时不可用",
			"time":    time.Now(),
		})
	}
}
