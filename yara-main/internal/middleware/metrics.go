package middleware

import (
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// Metrics 指标结构
type Metrics struct {
	mu sync.RWMutex

	// 请求统计
	TotalRequests    int64
	SuccessfulRequests int64
	FailedRequests   int64
	TimeoutRequests  int64

	// 响应时间统计
	TotalResponseTime time.Duration
	MinResponseTime   time.Duration
	MaxResponseTime   time.Duration
	AvgResponseTime   time.Duration

	// 并发统计
	CurrentConcurrency int64
	MaxConcurrency     int64

	// 错误统计
	ErrorCounts map[string]int64

	// 端点统计
	EndpointStats map[string]*EndpointStat
}

// EndpointStat 端点统计
type EndpointStat struct {
	Count         int64
	TotalTime     time.Duration
	MinTime       time.Duration
	MaxTime       time.Duration
	ErrorCount    int64
	LastAccess    time.Time
}

// MetricsCollector 指标收集器
type MetricsCollector struct {
	metrics *Metrics
	logger  *logrus.Logger
}

// NewMetricsCollector 创建指标收集器
func NewMetricsCollector(logger *logrus.Logger) *MetricsCollector {
	return &MetricsCollector{
		metrics: &Metrics{
			ErrorCounts:   make(map[string]int64),
			EndpointStats: make(map[string]*EndpointStat),
		},
		logger: logger,
	}
}

// Collect 收集指标中间件
func (mc *MetricsCollector) Collect() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		// 增加并发计数
		mc.metrics.mu.Lock()
		mc.metrics.CurrentConcurrency++
		if mc.metrics.CurrentConcurrency > mc.metrics.MaxConcurrency {
			mc.metrics.MaxConcurrency = mc.metrics.CurrentConcurrency
		}
		mc.metrics.mu.Unlock()

		// 处理请求
		c.Next()

		// 计算响应时间
		duration := time.Since(start)

		// 更新指标
		mc.metrics.mu.Lock()
		mc.metrics.TotalRequests++
		mc.metrics.TotalResponseTime += duration

		// 更新响应时间统计
		if mc.metrics.MinResponseTime == 0 || duration < mc.metrics.MinResponseTime {
			mc.metrics.MinResponseTime = duration
		}
		if duration > mc.metrics.MaxResponseTime {
			mc.metrics.MaxResponseTime = duration
		}
		mc.metrics.AvgResponseTime = mc.metrics.TotalResponseTime / time.Duration(mc.metrics.TotalRequests)

		// 更新成功/失败统计
		if c.Writer.Status() >= 200 && c.Writer.Status() < 400 {
			mc.metrics.SuccessfulRequests++
		} else {
			mc.metrics.FailedRequests++
		}

		// 更新端点统计
		endpoint := c.Request.Method + " " + c.FullPath()
		if stat, exists := mc.metrics.EndpointStats[endpoint]; exists {
			stat.Count++
			stat.TotalTime += duration
			if stat.MinTime == 0 || duration < stat.MinTime {
				stat.MinTime = duration
			}
			if duration > stat.MaxTime {
				stat.MaxTime = duration
			}
			if c.Writer.Status() >= 400 {
				stat.ErrorCount++
			}
			stat.LastAccess = time.Now()
		} else {
			mc.metrics.EndpointStats[endpoint] = &EndpointStat{
				Count:      1,
				TotalTime:  duration,
				MinTime:    duration,
				MaxTime:    duration,
				ErrorCount: 0,
				LastAccess: time.Now(),
			}
		}

		// 减少并发计数
		mc.metrics.CurrentConcurrency--
		mc.metrics.mu.Unlock()

		// 记录慢请求
		if duration > 5*time.Second {
			mc.logger.Warnf("慢请求: %s %s, 耗时: %v", c.Request.Method, c.Request.URL.Path, duration)
		}
	}
}

// GetMetrics 获取指标
func (mc *MetricsCollector) GetMetrics() *Metrics {
	mc.metrics.mu.RLock()
	defer mc.metrics.mu.RUnlock()
	
	// 返回副本以避免并发问题
	metrics := *mc.metrics
	return &metrics
}

// ResetMetrics 重置指标
func (mc *MetricsCollector) ResetMetrics() {
	mc.metrics.mu.Lock()
	defer mc.metrics.mu.Unlock()
	
	mc.metrics.TotalRequests = 0
	mc.metrics.SuccessfulRequests = 0
	mc.metrics.FailedRequests = 0
	mc.metrics.TimeoutRequests = 0
	mc.metrics.TotalResponseTime = 0
	mc.metrics.MinResponseTime = 0
	mc.metrics.MaxResponseTime = 0
	mc.metrics.AvgResponseTime = 0
	mc.metrics.CurrentConcurrency = 0
	mc.metrics.MaxConcurrency = 0
	mc.metrics.ErrorCounts = make(map[string]int64)
	mc.metrics.EndpointStats = make(map[string]*EndpointStat)
}

// ReportMetrics 报告指标
func (mc *MetricsCollector) ReportMetrics() {
	metrics := mc.GetMetrics()
	
	mc.logger.Infof("API指标报告 - 总请求: %d, 成功: %d, 失败: %d, 平均响应时间: %v, 当前并发: %d, 最大并发: %d",
		metrics.TotalRequests,
		metrics.SuccessfulRequests,
		metrics.FailedRequests,
		metrics.AvgResponseTime,
		metrics.CurrentConcurrency,
		metrics.MaxConcurrency,
	)

	// 报告端点统计
	for endpoint, stat := range metrics.EndpointStats {
		if stat.Count > 0 {
			avgTime := stat.TotalTime / time.Duration(stat.Count)
			mc.logger.Infof("端点 %s - 请求数: %d, 平均时间: %v, 错误数: %d",
				endpoint, stat.Count, avgTime, stat.ErrorCount)
		}
	}
} 