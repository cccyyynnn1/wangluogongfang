package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"strings"
	"syscall"
	"time"

	"yara-security-service/internal/config"
	"yara-security-service/internal/handlers"
	"yara-security-service/internal/middleware"
	"yara-security-service/internal/services"
	"yara-security-service/internal/utils"
	"yara-security-service/pkg/file"
	"yara-security-service/pkg/network"

	"yara-security-service/pkg/registry"
	"yara-security-service/pkg/security"
	"yara-security-service/pkg/user"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func resolveRulesPath(base string, logger *logrus.Logger) string {
	// 若为绝对路径，直接返回
	if filepath.IsAbs(base) {
		return base
	}
	// 优先以可执行文件所在目录为基准
	exe, err := os.Executable()
	if err == nil {
		exeDir := filepath.Dir(exe)
		p := filepath.Clean(filepath.Join(exeDir, base))
		if st, err := os.Stat(p); err == nil && st.IsDir() {
			return p
		}
	}
	// 次选：以当前工作目录为基准
	if wd, err := os.Getwd(); err == nil {
		p := filepath.Clean(filepath.Join(wd, base))
		if st, err := os.Stat(p); err == nil && st.IsDir() {
			return p
		}
	}
	// 兜底返回原始值
	logger.Warnf("无法解析规则目录为有效目录: %s，按原路径尝试", base)
	return base
}

func main() {
	// 初始化日志
	logger := logrus.New()
	logger.SetFormatter(&logrus.JSONFormatter{})
	logger.SetLevel(logrus.InfoLevel)

	// 加载配置
	cfg, err := config.LoadConfig("configs/config.yaml")
	if err != nil {
		logger.Fatalf("加载配置失败: %v", err)
	}

	// 解析 YARA 规则目录为绝对路径
	rulesDir := resolveRulesPath(cfg.Security.YaraRulesPath, logger)
	logger.Infof("使用规则目录: %s", rulesDir)

	// 初始化各个管理器
	scanner, err := security.NewScanner(rulesDir, logger)
	if err != nil {
		logger.Fatalf("初始化安全扫描器失败: %v", err)
	}

	registryManager := registry.NewManager(logger)
	networkManager := network.NewManager(logger)
	fileManager := file.NewManager(logger)
	userManager := user.NewManager(logger)

	// 初始化服务
	securityService := services.NewSecurityService(scanner, logger)
	processService := services.NewProcessService(logger)
	networkService := services.NewNetworkService(networkManager, logger)
	registryService := services.NewRegistryService(registryManager, logger)
	userService := services.NewUserService(userManager, logger)

	// 初始化处理器
	fileHandler := handlers.NewFileHandler(securityService, scanner, logger)
	processHandler := handlers.NewProcessHandler(processService, logger)
	registryHandler := handlers.NewRegistryHandler(registryService, logger)
	networkHandler := handlers.NewNetworkHandler(networkService, logger)
	userHandler := handlers.NewUserHandler(userService, logger)
	securityHandler := handlers.NewSecurityHandler(securityService, logger)

	// 设置Gin模式
	if cfg.Logging.Level == "debug" {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}

	// 创建路由
	router := gin.New()

	// 应用中间件
	router.Use(middleware.Recovery(logger))
	router.Use(middleware.RequestID())
	router.Use(middleware.Logging(middleware.DefaultLoggingConfig(), logger))
	router.Use(middleware.CORS(middleware.DefaultCORSConfig()))

	// 应用速率限制中间件
	rateLimitConfig := middleware.DefaultRateLimitConfig()
	router.Use(middleware.RateLimit(rateLimitConfig))

	// 应用超时中间件（适度放宽安全热更新等操作时间）
	timeoutConfig := middleware.DefaultTimeoutConfig()
	if timeoutConfig.DefaultTimeout < 30*time.Second {
		timeoutConfig.DefaultTimeout = 30 * time.Second
	}
	router.Use(middleware.Timeout(timeoutConfig, logger))

	// 应用指标收集中间件
	metricsCollector := middleware.NewMetricsCollector(logger)
	router.Use(metricsCollector.Collect())

	// 应用认证中间件（如果启用）
	if cfg.Security.APIKeyRequired {
		authConfig := &middleware.AuthConfig{
			Enabled:  true,
			APIKey:   cfg.Security.APIKey,
			Required: true,
		}
		router.Use(middleware.Auth(authConfig, logger))
	}

	// API路由组
	api := router.Group("/api/v1")
	{
		// 文件相关路由
		fileGroup := api.Group("/file")
		{
			fileGroup.POST("/scan", fileHandler.ScanFile)
			// 使用 /:path 捕获路径参数，支持查询参数传递完整路径
			fileGroup.GET("/info/:path", fileHandler.GetFileInfo)
			fileGroup.POST("/scan-directory", fileHandler.ScanDirectory)
			fileGroup.POST("/scan-buffer", fileHandler.ScanBuffer)
			fileGroup.GET("/list", fileHandler.GetFileList)

			// 文件管理路由
			fileGroup.POST("/copy", func(c *gin.Context) {
				var req struct {
					Source string `json:"source" binding:"required"`
					Dest   string `json:"dest" binding:"required"`
				}
				if err := c.ShouldBindJSON(&req); err != nil {
					utils.ValidationErrorResponse(c, err.Error())
					return
				}
				err := fileManager.CopyFile(req.Source, req.Dest)
				if err != nil {
					utils.OperationFailedResponse(c, "文件复制", err)
					return
				}
				utils.OperationSuccessResponse(c, "文件复制")
			})

			fileGroup.POST("/move", func(c *gin.Context) {
				var req struct {
					Source string `json:"source" binding:"required"`
					Dest   string `json:"dest" binding:"required"`
				}
				if err := c.ShouldBindJSON(&req); err != nil {
					utils.ValidationErrorResponse(c, err.Error())
					return
				}
				err := fileManager.MoveFile(req.Source, req.Dest)
				if err != nil {
					utils.OperationFailedResponse(c, "文件移动", err)
					return
				}
				utils.OperationSuccessResponse(c, "文件移动")
			})

			fileGroup.DELETE("/:path", func(c *gin.Context) {
				localPath := c.Param("path")
				// 如果路径参数不包含完整路径，尝试从查询参数获取完整路径
				if !strings.Contains(localPath, "/") && !strings.Contains(localPath, "\\") {
					fullPath := c.Query("fullPath")
					if fullPath != "" {
						localPath = fullPath
					}
				}
				err := fileManager.DeleteFile(localPath)
				if err != nil {
					utils.OperationFailedResponse(c, "文件删除", err)
					return
				}
				utils.OperationSuccessResponse(c, "文件删除")
			})

			fileGroup.GET("/hash/:path", func(c *gin.Context) {
				localPath := c.Param("path")
				// 如果路径参数不包含完整路径，尝试从查询参数获取完整路径
				if !strings.Contains(localPath, "/") && !strings.Contains(localPath, "\\") {
					fullPath := c.Query("fullPath")
					if fullPath != "" {
						localPath = fullPath
					}
				}
				algorithm := c.DefaultQuery("algorithm", "sha256")
				hash, err := fileManager.GetFileHash(localPath, algorithm)
				if err != nil {
					utils.OperationFailedResponse(c, "获取文件哈希", err)
					return
				}
				utils.SuccessResponse(c, "获取文件哈希成功", gin.H{
					"hash":      hash,
					"algorithm": algorithm,
					"file_path": localPath,
				})
			})

			fileGroup.GET("/hashes/:path", func(c *gin.Context) {
				localPath := c.Param("path")
				// 如果路径参数不包含完整路径，尝试从查询参数获取完整路径
				if !strings.Contains(localPath, "/") && !strings.Contains(localPath, "\\") {
					fullPath := c.Query("fullPath")
					if fullPath != "" {
						localPath = fullPath
					}
				}
				hashes, err := fileManager.GetFileHashes(localPath)
				if err != nil {
					utils.OperationFailedResponse(c, "获取文件哈希", err)
					return
				}
				utils.SuccessResponse(c, "获取文件所有哈希成功", gin.H{
					"hashes":    hashes,
					"file_path": localPath,
				})
			})

			fileGroup.POST("/verify-hash", func(c *gin.Context) {
				var req struct {
					FilePath     string `json:"file_path" binding:"required"`
					Algorithm    string `json:"algorithm" binding:"required"`
					ExpectedHash string `json:"expected_hash" binding:"required"`
				}
				if err := c.ShouldBindJSON(&req); err != nil {
					utils.ValidationErrorResponse(c, err.Error())
					return
				}

				// 添加调试信息
				logger.Infof("验证文件哈希: 文件路径=%s, 算法=%s, 期望哈希=%s", req.FilePath, req.Algorithm, req.ExpectedHash)

				// 先获取实际哈希值用于调试
				actualHash, err := fileManager.GetFileHash(req.FilePath, req.Algorithm)
				if err != nil {
					logger.Errorf("获取文件哈希失败: %v", err)
					utils.OperationFailedResponse(c, "获取文件哈希", err)
					return
				}

				logger.Infof("文件实际哈希: %s", actualHash)

				valid, err := fileManager.VerifyFileHash(req.FilePath, req.Algorithm, req.ExpectedHash)
				if err != nil {
					logger.Errorf("验证文件哈希失败: %v", err)
					utils.OperationFailedResponse(c, "验证文件哈希", err)
					return
				}

				logger.Infof("哈希验证结果: %t", valid)

				utils.SuccessResponse(c, "哈希验证完成", gin.H{
					"valid":         valid,
					"file_path":     req.FilePath,
					"algorithm":     req.Algorithm,
					"expected_hash": req.ExpectedHash,
					"actual_hash":   actualHash,
				})
			})
		}

		// 进程相关路由
		processGroup := api.Group("/process")
		{
			processGroup.GET("/list", processHandler.GetProcesses)
			processGroup.GET("/:pid", processHandler.GetProcessByPID)
			processGroup.POST("/start", processHandler.StartProcess)
			processGroup.DELETE("/:pid", processHandler.KillProcess)
			processGroup.PUT("/:pid/suspend", processHandler.SuspendProcess)
			processGroup.PUT("/:pid/resume", processHandler.ResumeProcess)
			processGroup.GET("/:pid/modules", processHandler.GetProcessModules)
			processGroup.GET("/:pid/connections", processHandler.GetProcessConnections)
			processGroup.GET("/:pid/memory", processHandler.GetProcessMemoryInfo)
			processGroup.GET("/:pid/running", processHandler.GetProcessRunningStatus)
			processGroup.GET("/:pid/children", processHandler.GetProcessChildren)

			// 进程监控相关路由
			processGroup.POST("/monitoring/enable", processHandler.EnableProcessMonitoring)
			processGroup.POST("/monitoring/disable", processHandler.DisableProcessMonitoring)
			processGroup.GET("/monitoring/list", processHandler.GetMonitoredProcesses)
			processGroup.GET("/statistics", processHandler.GetProcessStatistics)
		}

		// 模块相关路由（使用 /:module 避免通配符冲突）
		moduleGroup := api.Group("/process/module")
		{
			moduleGroup.GET("/list", processHandler.GetSystemModules)
			moduleGroup.GET("/:module", processHandler.GetModuleInfo)
			moduleGroup.PUT("/:module/suspend", processHandler.SuspendModule)
			moduleGroup.PUT("/:module/resume", processHandler.ResumeModule)
			moduleGroup.DELETE("/:module", processHandler.KillModule)
		}

		// 注册表相关路由
		registryGroup := api.Group("/registry")
		{
			registryGroup.GET("/key/:path", registryHandler.GetRegistryKey)
			registryGroup.POST("/key", registryHandler.CreateRegistryKey)
			registryGroup.DELETE("/key/:path", registryHandler.DeleteRegistryKey)
			registryGroup.PUT("/value", registryHandler.SetRegistryValue)
			registryGroup.GET("/keys/:path", registryHandler.ListRegistryKeys)
			registryGroup.GET("/values/:path", registryHandler.ListRegistryValues)
			registryGroup.GET("/search", registryHandler.SearchRegistry)
			registryGroup.POST("/search", registryHandler.SearchRegistry)
		}

		// 注册表值相关路由（使用/:path/:name格式避免通配符冲突）
		registryValueGroup := api.Group("/registry/value")
		{
			registryValueGroup.GET("/:path", registryHandler.GetRegistryValue)
			registryValueGroup.DELETE("/:path", registryHandler.DeleteRegistryValue)
			registryValueGroup.GET("/:path/:name", registryHandler.GetRegistryValueByName)       // 支持值名称作为路径参数
			registryValueGroup.DELETE("/:path/:name", registryHandler.DeleteRegistryValueByName) // 支持值名称作为路径参数
		}

		// 网络相关路由
		networkGroup := api.Group("/network")
		{
			networkGroup.GET("/connections", networkHandler.GetConnections)
			networkGroup.GET("/connections/tcp", networkHandler.GetTCPConnections)
			networkGroup.GET("/connections/udp", networkHandler.GetUDPConnections)
			networkGroup.GET("/connections/pid/:pid", networkHandler.GetConnectionsByPID)
			networkGroup.GET("/connections/port/:port", networkHandler.GetConnectionsByPort)
			networkGroup.GET("/connections/ip/:ip", networkHandler.GetConnectionsByIP)
			networkGroup.GET("/connection/:id", networkHandler.GetConnectionByID)
			networkGroup.POST("/connection/close", networkHandler.CloseConnection)
			networkGroup.GET("/ports/in-use", networkHandler.GetPortsInUse)
			networkGroup.GET("/port/:port/in-use", networkHandler.IsPortInUse)
			networkGroup.GET("/listening-ports", networkHandler.GetListeningPorts)
			networkGroup.GET("/established-connections", networkHandler.GetEstablishedConnections)
			networkGroup.GET("/interfaces", networkHandler.GetInterfaces)
			networkGroup.GET("/stats", networkHandler.GetNetworkStats)

			// 网络监控相关路由
			networkGroup.POST("/monitoring/enable", networkHandler.EnableNetworkMonitoring)
			networkGroup.POST("/monitoring/disable", networkHandler.DisableNetworkMonitoring)
			networkGroup.GET("/monitoring/connections", networkHandler.GetMonitoredConnections)
			networkGroup.GET("/monitoring/history", networkHandler.GetConnectionHistory)
		}

		// 安全相关路由
		securityGroup := api.Group("/security")
		{
			securityGroup.POST("/scan", securityHandler.ScanFile)
			securityGroup.GET("/rules", securityHandler.GetYaraRules)
			securityGroup.GET("/status", securityHandler.GetScanStatus)
			securityGroup.GET("/rules/info", securityHandler.GetRulesInfo)
			securityGroup.POST("/reload-rules", securityHandler.ReloadRules)
			securityGroup.GET("/reload-status", securityHandler.ReloadStatus)
			securityGroup.POST("/cache/clear", securityHandler.ClearCache)
			securityGroup.GET("/cache/stats", securityHandler.GetCacheStats)
			securityGroup.POST("/quarantine", securityHandler.QuarantineFile)
			securityGroup.POST("/restore", securityHandler.RestoreFile)
			securityGroup.GET("/quarantine/list", securityHandler.GetQuarantineList)
			securityGroup.GET("/scan-history", securityHandler.GetScanHistory)
		}

		// 用户相关路由
		userGroup := api.Group("/user")
		{
			userGroup.GET("/current", userHandler.GetCurrentUser)
			userGroup.GET("/id/:uid", userHandler.GetUserByID)
			userGroup.GET("/name/:username", userHandler.GetUserByName)
			userGroup.GET("/all", userHandler.GetAllUsers)
			userGroup.POST("/permissions/check", userHandler.CheckUserPermissions)
			userGroup.POST("/password/validate", userHandler.ValidatePassword)
			userGroup.GET("/password/policy", userHandler.GetPasswordPolicy)
			userGroup.GET("/status/:username", userHandler.CheckAccountStatus)
			userGroup.GET("/sessions/:username", userHandler.GetUserSessions)
			userGroup.DELETE("/session/:sessionId", userHandler.KillUserSession)
			userGroup.POST("/lock/:username", userHandler.LockUserAccount)
			userGroup.POST("/unlock/:username", userHandler.UnlockUserAccount)
			userGroup.POST("/password/change", userHandler.ChangeUserPassword)
			userGroup.GET("/groups/:username", userHandler.GetUserGroups)
			userGroup.GET("/history/:username", userHandler.GetUserLoginHistory)
		}
	}

	// 健康检查
	router.GET("/api/health", utils.HealthCheckResponse)

	// 指标端点
	router.GET("/api/metrics", func(c *gin.Context) {
		metrics := metricsCollector.GetMetrics()
		utils.MetricsResponse(c, metrics)
	})

	// 重置指标端点
	router.POST("/api/metrics/reset", func(c *gin.Context) {
		metricsCollector.ResetMetrics()
		utils.OperationSuccessResponse(c, "指标重置")
	})

	// 创建HTTP服务器
	server := &http.Server{
		Addr:         fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port),
		Handler:      router,
		ReadTimeout:  cfg.Server.ReadTimeout,
		WriteTimeout: cfg.Server.WriteTimeout,
	}

	// 启动服务器
	go func() {
		logger.Infof("服务器启动在 %s:%d (rules: %s)", cfg.Server.Host, cfg.Server.Port, rulesDir)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatalf("服务器启动失败: %v", err)
		}
	}()

	// 启动定期指标报告
	go func() {
		ticker := time.NewTicker(5 * time.Minute) // 每5分钟报告一次指标
		defer ticker.Stop()

		for range ticker.C {
			metricsCollector.ReportMetrics()
		}
	}()

	// 等待中断信号
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info("正在关闭服务器...")

	// 优雅关闭
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		logger.Fatalf("服务器强制关闭: %v", err)
	}

	logger.Info("服务器已关闭")
}
