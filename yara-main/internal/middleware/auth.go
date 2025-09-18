package middleware

import (
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// AuthConfig 认证配置
type AuthConfig struct {
	Enabled     bool
	APIKey      string
	Required    bool
	ExemptPaths []string
}

// DefaultAuthConfig 默认认证配置
func DefaultAuthConfig() *AuthConfig {
	return &AuthConfig{
		Enabled:     false,
		APIKey:      "",
		Required:    false,
		ExemptPaths: []string{"/health", "/api/health"},
	}
}

// Auth 认证中间件
func Auth(config *AuthConfig, logger *logrus.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		if !config.Enabled {
			c.Next()
			return
		}

		// 检查是否为豁免路径
		if isExemptPath(c.Request.URL.Path, config.ExemptPaths) {
			c.Next()
			return
		}

		// 获取API密钥
		apiKey := c.GetHeader("X-API-Key")
		if apiKey == "" {
			apiKey = c.Query("api_key")
		}

		// 验证API密钥
		if !isValidAPIKey(apiKey, config.APIKey) {
			logger.Warnf("无效的API密钥访问: %s", c.ClientIP())
			c.JSON(http.StatusUnauthorized, gin.H{
				"code":    401,
				"message": "无效的API密钥",
				"data":    nil,
			})
			c.Abort()
			return
		}

		c.Next()
	}
}

// isExemptPath 检查是否为豁免路径
func isExemptPath(path string, exemptPaths []string) bool {
	for _, exemptPath := range exemptPaths {
		if strings.HasPrefix(path, exemptPath) {
			return true
		}
	}
	return false
}

// isValidAPIKey 验证API密钥
func isValidAPIKey(providedKey, expectedKey string) bool {
	if expectedKey == "" {
		return true // 如果未设置预期密钥，则允许所有请求
	}
	return providedKey == expectedKey
}

// RateLimitConfig 速率限制配置
type RateLimitConfig struct {
	RequestsPerMinute int
	BurstSize         int
	WindowSize        time.Duration
}

// DefaultRateLimitConfig 默认速率限制配置
func DefaultRateLimitConfig() *RateLimitConfig {
	return &RateLimitConfig{
		RequestsPerMinute: 60,
		BurstSize:         10,
		WindowSize:        1 * time.Minute,
	}
}

// RateLimiter 速率限制器
type RateLimiter struct {
	requests map[string][]time.Time
	mu       sync.RWMutex
	config   *RateLimitConfig
}

// NewRateLimiter 创建速率限制器
func NewRateLimiter(config *RateLimitConfig) *RateLimiter {
	return &RateLimiter{
		requests: make(map[string][]time.Time),
		config:   config,
	}
}

// Allow 检查是否允许请求
func (r *RateLimiter) Allow(clientIP string) bool {
	r.mu.Lock()
	defer r.mu.Unlock()

	now := time.Now()
	windowStart := now.Add(-r.config.WindowSize)

	// 清理过期的请求记录
	if requests, exists := r.requests[clientIP]; exists {
		var validRequests []time.Time
		for _, reqTime := range requests {
			if reqTime.After(windowStart) {
				validRequests = append(validRequests, reqTime)
			}
		}
		r.requests[clientIP] = validRequests
	}

	// 检查请求数量
	if len(r.requests[clientIP]) >= r.config.RequestsPerMinute {
		return false
	}

	// 添加新请求
	r.requests[clientIP] = append(r.requests[clientIP], now)
	return true
}

// RateLimit 速率限制中间件
func RateLimit(config *RateLimitConfig) gin.HandlerFunc {
	limiter := NewRateLimiter(config)

	return func(c *gin.Context) {
		clientIP := c.ClientIP()

		// 检查速率限制
		if !limiter.Allow(clientIP) {
			c.JSON(http.StatusTooManyRequests, gin.H{
				"code":    429,
				"message": "请求过于频繁，请稍后再试",
				"data":    nil,
			})
			c.Abort()
			return
		}

		c.Next()
	}
}

// IPWhitelist IP白名单中间件
func IPWhitelist(allowedIPs []string) gin.HandlerFunc {
	return func(c *gin.Context) {
		clientIP := c.ClientIP()

		// 检查IP是否在白名单中
		if !isIPAllowed(clientIP, allowedIPs) {
			c.JSON(http.StatusForbidden, gin.H{
				"code":    403,
				"message": "IP地址不在白名单中",
				"data":    nil,
			})
			c.Abort()
			return
		}

		c.Next()
	}
}

// isIPAllowed 检查IP是否被允许
func isIPAllowed(clientIP string, allowedIPs []string) bool {
	if len(allowedIPs) == 0 {
		return true // 如果白名单为空，则允许所有IP
	}

	for _, allowedIP := range allowedIPs {
		if allowedIP == "*" || allowedIP == clientIP {
			return true
		}
	}

	return false
}
