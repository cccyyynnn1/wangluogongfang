package handlers

import (
	"fmt"
	"strconv"
	"time"

	"yara-security-service/internal/models"
	"yara-security-service/internal/services"
	"yara-security-service/internal/utils"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// ProcessHandler 进程处理器
type ProcessHandler struct {
	processService *services.ProcessService
	logger         *logrus.Logger
}

// NewProcessHandler 创建进程处理器
func NewProcessHandler(processService *services.ProcessService, logger *logrus.Logger) *ProcessHandler {
	return &ProcessHandler{
		processService: processService,
		logger:         logger,
	}
}

// GetProcesses 获取进程列表
func (h *ProcessHandler) GetProcesses(c *gin.Context) {
	defer func() {
		if r := recover(); r != nil {
			h.logger.Errorf("GetProcesses panic: %v", r)
			utils.InternalServerErrorResponse(c, "获取进程列表时发生内部错误")
		}
	}()

	processes, err := h.processService.GetProcessList()
	if err != nil {
		h.logger.Errorf("获取进程列表失败: %v", err)
		utils.InternalServerErrorResponse(c, "获取进程列表失败")
		return
	}

	// 简化响应处理，避免复杂的类型转换
	if processes.Code == 200 {
		utils.SuccessResponse(c, "获取进程列表成功", processes.Data)
	} else {
		utils.InternalServerErrorResponse(c, processes.Message)
	}
}

// GetProcessByPID 根据PID获取进程信息
func (h *ProcessHandler) GetProcessByPID(c *gin.Context) {
	defer func() {
		if r := recover(); r != nil {
			h.logger.Errorf("GetProcessByPID panic: %v", r)
			utils.InternalServerErrorResponse(c, "获取进程信息时发生内部错误")
		}
	}()

	pidStr := c.Param("pid")
	pid, err := strconv.Atoi(pidStr)
	if err != nil {
		utils.BadRequestResponse(c, "无效的PID参数")
		return
	}

	process, err := h.processService.GetProcessByPID(int32(pid))
	if err != nil {
		h.logger.Errorf("获取进程信息失败: %v", err)
		utils.InternalServerErrorResponse(c, "获取进程信息失败")
		return
	}

	utils.InfoSuccessResponse(c, process)
}

// StartProcess 启动进程
func (h *ProcessHandler) StartProcess(c *gin.Context) {
	defer func() {
		if r := recover(); r != nil {
			h.logger.Errorf("StartProcess panic: %v", r)
			utils.InternalServerErrorResponse(c, "启动进程时发生内部错误")
		}
	}()

	var request models.StartProcessRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		utils.ValidationErrorResponse(c, err.Error())
		return
	}

	proc, err := h.processService.StartProcess(c.Request.Context(), request.Command, request.Args, request.WorkingDir)
	if err != nil {
		h.logger.Errorf("启动进程失败: %v", err)
		utils.OperationFailedResponse(c, "启动进程", err)
		return
	}

	// 简化响应处理
	if proc.Code == 200 {
		utils.SuccessResponse(c, "启动进程成功", proc.Data)
	} else {
		utils.OperationFailedResponse(c, "启动进程", fmt.Errorf(proc.Message))
	}
}

// KillProcess 结束进程
func (h *ProcessHandler) KillProcess(c *gin.Context) {
	pidStr := c.Param("pid")
	pid, err := strconv.Atoi(pidStr)
	if err != nil {
		utils.BadRequestResponse(c, "无效的PID参数")
		return
	}

	result, err := h.processService.KillProcess(int32(pid))
	if err != nil {
		h.logger.Errorf("结束进程失败: %v", err)
		utils.OperationFailedResponse(c, "结束进程", err)
		return
	}

	if result.Code == 200 {
		utils.OperationSuccessResponse(c, "结束进程")
	} else {
		utils.OperationFailedResponse(c, "结束进程", fmt.Errorf(result.Message))
	}
}

// SuspendProcess 暂停进程
func (h *ProcessHandler) SuspendProcess(c *gin.Context) {
	pidStr := c.Param("pid")
	pid, err := strconv.Atoi(pidStr)
	if err != nil {
		utils.BadRequestResponse(c, "无效的PID参数")
		return
	}

	result, err := h.processService.SuspendProcess(int32(pid))
	if err != nil {
		h.logger.Errorf("暂停进程失败: %v", err)
		utils.OperationFailedResponse(c, "暂停进程", err)
		return
	}

	if result.Code == 200 {
		utils.OperationSuccessResponse(c, "暂停进程")
	} else {
		utils.OperationFailedResponse(c, "暂停进程", fmt.Errorf(result.Message))
	}
}

// ResumeProcess 恢复进程
func (h *ProcessHandler) ResumeProcess(c *gin.Context) {
	pidStr := c.Param("pid")
	pid, err := strconv.Atoi(pidStr)
	if err != nil {
		utils.BadRequestResponse(c, "无效的PID参数")
		return
	}

	result, err := h.processService.ResumeProcess(int32(pid))
	if err != nil {
		h.logger.Errorf("恢复进程失败: %v", err)
		utils.OperationFailedResponse(c, "恢复进程", err)
		return
	}

	if result.Code == 200 {
		utils.OperationSuccessResponse(c, "恢复进程")
	} else {
		utils.OperationFailedResponse(c, "恢复进程", fmt.Errorf(result.Message))
	}
}

// GetProcessModules 获取进程模块相关进程信息
func (h *ProcessHandler) GetProcessModules(c *gin.Context) {
	defer func() {
		if r := recover(); r != nil {
			h.logger.Errorf("GetProcessModules panic: %v", r)
			utils.InternalServerErrorResponse(c, "获取进程模块相关进程信息时发生内部错误")
		}
	}()

	pidStr := c.Param("pid")
	pid, err := strconv.Atoi(pidStr)
	if err != nil {
		utils.BadRequestResponse(c, "无效的PID参数")
		return
	}

	relatedProcesses, err := h.processService.GetProcessModules(int32(pid))
	if err != nil {
		h.logger.Errorf("获取进程模块相关进程信息失败: %v", err)
		utils.InternalServerErrorResponse(c, "获取进程模块相关进程信息失败")
		return
	}

	// 检查响应状态
	if relatedProcesses.Code != 200 {
		h.logger.Errorf("获取进程模块相关进程信息失败: %s", relatedProcesses.Message)
		utils.InternalServerErrorResponse(c, relatedProcesses.Message)
		return
	}

	// 构建响应数据
	response := gin.H{
		"target_pid":        pid,
		"related_processes": relatedProcesses.Data,
		"process_count":     relatedProcesses.Data.(map[string]interface{})["process_count"],
		"description":       fmt.Sprintf("与进程PID=%d共享模块的所有相关进程信息", pid),
		"timestamp":         time.Now(),
	}

	utils.SuccessResponse(c, "获取进程模块相关进程信息成功", response)
}

// GetProcessConnections 获取进程网络连接
func (h *ProcessHandler) GetProcessConnections(c *gin.Context) {
	pidStr := c.Param("pid")
	pid, err := strconv.Atoi(pidStr)
	if err != nil {
		utils.BadRequestResponse(c, "无效的PID参数")
		return
	}

	connections, err := h.processService.GetProcessConnections(int32(pid))
	if err != nil {
		h.logger.Errorf("获取进程网络连接失败: %v", err)
		utils.InternalServerErrorResponse(c, "获取进程网络连接失败")
		return
	}

	// 简化响应处理
	if connections.Code == 200 {
		utils.SuccessResponse(c, "获取进程网络连接成功", connections.Data)
	} else {
		utils.InternalServerErrorResponse(c, connections.Message)
	}
}

// GetProcessMemoryInfo 获取进程内存信息
func (h *ProcessHandler) GetProcessMemoryInfo(c *gin.Context) {
	pidStr := c.Param("pid")
	pid, err := strconv.Atoi(pidStr)
	if err != nil {
		utils.BadRequestResponse(c, "无效的PID参数")
		return
	}

	memoryInfo, err := h.processService.GetProcessMemoryInfo(int32(pid))
	if err != nil {
		h.logger.Errorf("获取进程内存信息失败: %v", err)
		utils.InternalServerErrorResponse(c, "获取进程内存信息失败")
		return
	}

	utils.InfoSuccessResponse(c, memoryInfo)
}

// GetProcessRunningStatus 获取进程运行状态
func (h *ProcessHandler) GetProcessRunningStatus(c *gin.Context) {
	pidStr := c.Param("pid")
	pid, err := strconv.Atoi(pidStr)
	if err != nil {
		utils.BadRequestResponse(c, "无效的PID参数")
		return
	}

	status, err := h.processService.IsProcessRunning(int32(pid))
	if err != nil {
		h.logger.Errorf("获取进程运行状态失败: %v", err)
		utils.InternalServerErrorResponse(c, "获取进程运行状态失败")
		return
	}

	utils.StatusSuccessResponse(c, gin.H{"is_running": status})
}

// GetProcessChildren 获取子进程
func (h *ProcessHandler) GetProcessChildren(c *gin.Context) {
	pidStr := c.Param("pid")
	pid, err := strconv.Atoi(pidStr)
	if err != nil {
		utils.BadRequestResponse(c, "无效的PID参数")
		return
	}

	children, err := h.processService.GetProcessChildren(int32(pid))
	if err != nil {
		h.logger.Errorf("获取子进程失败: %v", err)
		utils.InternalServerErrorResponse(c, "获取子进程失败")
		return
	}

	// 构建响应数据
	responseData := gin.H{
		"pid":      pid,
		"children": children,
	}

	utils.InfoSuccessResponse(c, responseData)
}

// GetSystemModules 获取系统模块
func (h *ProcessHandler) GetSystemModules(c *gin.Context) {
	modules, err := h.processService.GetSystemModules()
	if err != nil {
		h.logger.Errorf("获取系统模块失败: %v", err)
		utils.InternalServerErrorResponse(c, "获取系统模块失败")
		return
	}

	// 简化响应处理
	if modules.Code == 200 {
		utils.SuccessResponse(c, "获取系统模块成功", modules.Data)
	} else {
		utils.InternalServerErrorResponse(c, modules.Message)
	}
}

// GetModuleInfo 获取模块详细信息
func (h *ProcessHandler) GetModuleInfo(c *gin.Context) {
	module := c.Param("module")
	if module == "" {
		utils.BadRequestResponse(c, "模块名称不能为空")
		return
	}

	info, err := h.processService.GetModuleInfo(module)
	if err != nil {
		h.logger.Errorf("获取模块信息失败: %v", err)
		utils.InternalServerErrorResponse(c, "获取模块信息失败")
		return
	}

	utils.InfoSuccessResponse(c, info)
}

// SuspendModule 挂起系统模块
func (h *ProcessHandler) SuspendModule(c *gin.Context) {
	module := c.Param("module")
	if module == "" {
		utils.BadRequestResponse(c, "模块名称不能为空")
		return
	}

	result, err := h.processService.SuspendModule(module)
	if err != nil {
		h.logger.Errorf("挂起系统模块失败: %v", err)
		utils.OperationFailedResponse(c, "挂起系统模块", err)
		return
	}

	if result.Code == 200 {
		utils.OperationSuccessResponse(c, "挂起系统模块")
	} else {
		utils.OperationFailedResponse(c, "挂起系统模块", fmt.Errorf(result.Message))
	}
}

// ResumeModule 恢复系统模块
func (h *ProcessHandler) ResumeModule(c *gin.Context) {
	module := c.Param("module")
	if module == "" {
		utils.BadRequestResponse(c, "模块名称不能为空")
		return
	}

	result, err := h.processService.ResumeModule(module)
	if err != nil {
		h.logger.Errorf("恢复系统模块失败: %v", err)
		utils.OperationFailedResponse(c, "恢复系统模块", err)
		return
	}

	if result.Code == 200 {
		utils.OperationSuccessResponse(c, "恢复系统模块")
	} else {
		utils.OperationFailedResponse(c, "恢复系统模块", fmt.Errorf(result.Message))
	}
}

// KillModule 结束系统模块
func (h *ProcessHandler) KillModule(c *gin.Context) {
	module := c.Param("module")
	if module == "" {
		utils.BadRequestResponse(c, "模块名称不能为空")
		return
	}

	result, err := h.processService.KillModule(module)
	if err != nil {
		h.logger.Errorf("结束系统模块失败: %v", err)
		utils.OperationFailedResponse(c, "结束系统模块", err)
		return
	}

	if result.Code == 200 {
		utils.OperationSuccessResponse(c, "结束系统模块")
	} else {
		utils.OperationFailedResponse(c, "结束系统模块", fmt.Errorf(result.Message))
	}
}

// EnableProcessMonitoring 启用进程监控
func (h *ProcessHandler) EnableProcessMonitoring(c *gin.Context) {
	result, err := h.processService.EnableProcessMonitoring()
	if err != nil {
		h.logger.Errorf("启用进程监控失败: %v", err)
		utils.OperationFailedResponse(c, "启用进程监控", err)
		return
	}

	if result.Code == 200 {
		utils.OperationSuccessResponse(c, "进程监控已启用")
	} else {
		utils.OperationFailedResponse(c, "启用进程监控", fmt.Errorf(result.Message))
	}
}

// DisableProcessMonitoring 禁用进程监控
func (h *ProcessHandler) DisableProcessMonitoring(c *gin.Context) {
	result, err := h.processService.DisableProcessMonitoring()
	if err != nil {
		h.logger.Errorf("禁用进程监控失败: %v", err)
		utils.OperationFailedResponse(c, "禁用进程监控", err)
		return
	}

	if result.Code == 200 {
		utils.OperationSuccessResponse(c, "进程监控已禁用")
	} else {
		utils.OperationFailedResponse(c, "禁用进程监控", fmt.Errorf(result.Message))
	}
}

// GetMonitoredProcesses 获取监控的进程列表
func (h *ProcessHandler) GetMonitoredProcesses(c *gin.Context) {
	result, err := h.processService.GetMonitoredProcesses()
	if err != nil {
		h.logger.Errorf("获取监控进程列表失败: %v", err)
		utils.InternalServerErrorResponse(c, "获取监控进程列表失败")
		return
	}

	if result.Code == 200 {
		utils.SuccessResponse(c, "获取监控进程列表成功", result.Data)
	} else {
		utils.InternalServerErrorResponse(c, result.Message)
	}
}

// GetProcessStatistics 获取进程统计信息
func (h *ProcessHandler) GetProcessStatistics(c *gin.Context) {
	stats, err := h.processService.GetProcessStatistics()
	if err != nil {
		h.logger.Errorf("获取进程统计信息失败: %v", err)
		utils.InternalServerErrorResponse(c, "获取进程统计信息失败")
		return
	}

	utils.InfoSuccessResponse(c, stats)
}
