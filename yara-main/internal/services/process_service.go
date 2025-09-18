package services

import (
	"context"
	"fmt"
	"path/filepath"
	"time"

	"yara-security-service/internal/models"
	"yara-security-service/pkg/process"

	"github.com/sirupsen/logrus"
)

// ProcessService 进程服务
type ProcessService struct {
	manager *process.Manager
	logger  *logrus.Logger
}

// NewProcessService 创建新的进程服务
func NewProcessService(logger *logrus.Logger) *ProcessService {
	return &ProcessService{
		manager: process.NewManager(logger),
		logger:  logger,
	}
}

// GetProcessList 获取进程列表
func (s *ProcessService) GetProcessList() (*models.Response, error) {
	startTime := time.Now()

	processes, err := s.manager.GetProcessList()
	if err != nil {
		s.logger.Errorf("获取进程列表失败: %v", err)
		return &models.Response{
			Code:    500,
			Message: fmt.Sprintf("获取进程列表失败: %v", err),
			Data:    nil,
		}, nil
	}

	// 构建响应数据
	data := map[string]interface{}{
		"processes": processes,
		"count":     len(processes),
	}

	response := &models.Response{
		Code:    200,
		Message: "获取成功",
		Data:    data,
	}

	s.logger.Debugf("获取进程列表成功，共%d个进程，耗时=%v", len(processes), time.Since(startTime))
	return response, nil
}

// GetProcessByPID 根据PID获取进程信息
func (s *ProcessService) GetProcessByPID(pid int32) (*models.Response, error) {
	startTime := time.Now()

	processInfo, err := s.manager.GetProcessByPID(pid)
	if err != nil {
		s.logger.Errorf("获取进程信息失败 PID=%d: %v", pid, err)
		return &models.Response{
			Code:    404,
			Message: fmt.Sprintf("进程不存在或无法访问: %v", err),
			Data:    nil,
		}, nil
	}

	response := &models.Response{
		Code:    200,
		Message: "获取成功",
		Data:    processInfo,
	}

	s.logger.Debugf("获取进程信息成功 PID=%d，耗时=%v", pid, time.Since(startTime))
	return response, nil
}

// StartProcess 启动新进程
func (s *ProcessService) StartProcess(ctx context.Context, command string, args []string, workingDir string) (*models.Response, error) {
	startTime := time.Now()

	// 构建完整的进程路径
	var processPath string
	if len(args) > 0 {
		// 如果有参数，使用第一个参数作为进程路径
		processPath = args[0]
	} else {
		// 否则使用command作为进程路径
		processPath = command
	}

	// 如果workingDir不为空，构建相对路径
	if workingDir != "" && !filepath.IsAbs(processPath) {
		processPath = filepath.Join(workingDir, processPath)
	}

	// 调用新的StartProcess方法
	err := s.manager.StartProcess(processPath)
	if err != nil {
		s.logger.Errorf("启动进程失败: %v", err)
		return &models.Response{
			Code:    500,
			Message: fmt.Sprintf("启动进程失败: %v", err),
			Data:    nil,
		}, nil
	}

	// 构建响应数据
	data := map[string]interface{}{
		"command":      command,
		"process_path": processPath,
		"working_dir":  workingDir,
		"start_time":   time.Now(),
		"status":       "started",
		"message":      "进程启动成功，正在运行中",
	}

	response := &models.Response{
		Code:    200,
		Message: "进程启动成功",
		Data:    data,
	}

	s.logger.Debugf("启动进程成功，路径=%s，耗时=%v", processPath, time.Since(startTime))
	return response, nil
}

// KillProcess 结束进程
func (s *ProcessService) KillProcess(pid int32) (*models.Response, error) {
	startTime := time.Now()

	err := s.manager.KillProcess(pid)
	if err != nil {
		s.logger.Errorf("结束进程失败 PID=%d: %v", pid, err)
		return &models.Response{
			Code:    500,
			Message: fmt.Sprintf("结束进程失败: %v", err),
			Data:    nil,
		}, nil
	}

	response := &models.Response{
		Code:    200,
		Message: "进程已结束",
		Data:    nil,
	}

	s.logger.Infof("进程结束成功 PID=%d，耗时=%v", pid, time.Since(startTime))
	return response, nil
}

// SuspendProcess 挂起进程
func (s *ProcessService) SuspendProcess(pid int32) (*models.Response, error) {
	startTime := time.Now()

	err := s.manager.SuspendProcess(pid)
	if err != nil {
		s.logger.Errorf("挂起进程失败 PID=%d: %v", pid, err)
		return &models.Response{
			Code:    500,
			Message: fmt.Sprintf("挂起进程失败: %v", err),
			Data:    nil,
		}, nil
	}

	response := &models.Response{
		Code:    200,
		Message: "进程已挂起",
		Data:    nil,
	}

	s.logger.Infof("进程挂起成功 PID=%d，耗时=%v", pid, time.Since(startTime))
	return response, nil
}

// ResumeProcess 恢复进程
func (s *ProcessService) ResumeProcess(pid int32) (*models.Response, error) {
	startTime := time.Now()

	err := s.manager.ResumeProcess(pid)
	if err != nil {
		s.logger.Errorf("恢复进程失败 PID=%d: %v", pid, err)
		return &models.Response{
			Code:    500,
			Message: fmt.Sprintf("恢复进程失败: %v", err),
			Data:    nil,
		}, nil
	}

	response := &models.Response{
		Code:    200,
		Message: "进程已恢复",
		Data:    nil,
	}

	s.logger.Infof("进程恢复成功 PID=%d，耗时=%v", pid, time.Since(startTime))
	return response, nil
}

// GetProcessModules 获取进程模块相关进程信息
func (s *ProcessService) GetProcessModules(pid int32) (*models.Response, error) {
	startTime := time.Now()

	relatedProcesses, err := s.manager.GetProcessModules(pid)
	if err != nil {
		s.logger.Errorf("获取进程模块相关进程信息失败 PID=%d: %v", pid, err)
		return &models.Response{
			Code:    500,
			Message: fmt.Sprintf("获取进程模块相关进程信息失败: %v", err),
			Data:    nil,
		}, nil
	}

	// 构建响应数据
	data := map[string]interface{}{
		"target_pid":        pid,
		"related_processes": relatedProcesses,
		"process_count":     len(relatedProcesses),
		"description":       fmt.Sprintf("与进程PID=%d共享模块的所有相关进程信息", pid),
		"timestamp":         time.Now(),
	}

	response := &models.Response{
		Code:    200,
		Message: "获取进程模块相关进程信息成功",
		Data:    data,
	}

	s.logger.Debugf("获取进程模块相关进程信息成功 PID=%d，共%d个相关进程，耗时=%v", pid, len(relatedProcesses), time.Since(startTime))
	return response, nil
}

// GetProcessConnections 获取进程网络连接
func (s *ProcessService) GetProcessConnections(pid int32) (*models.Response, error) {
	startTime := time.Now()

	connections, err := s.manager.GetProcessConnections(pid)
	if err != nil {
		s.logger.Errorf("获取进程网络连接失败 PID=%d: %v", pid, err)
		return &models.Response{
			Code:    500,
			Message: fmt.Sprintf("获取进程网络连接失败: %v", err),
			Data:    nil,
		}, nil
	}

	// 构建响应数据
	data := map[string]interface{}{
		"connections": connections,
	}

	response := &models.Response{
		Code:    200,
		Message: "获取成功",
		Data:    data,
	}

	s.logger.Debugf("获取进程网络连接成功 PID=%d，共%d个连接，耗时=%v", pid, len(connections), time.Since(startTime))
	return response, nil
}

// GetProcessMemoryInfo 获取进程内存信息
func (s *ProcessService) GetProcessMemoryInfo(pid int32) (*models.Response, error) {
	startTime := time.Now()

	memoryInfo, err := s.manager.GetProcessMemoryInfo(pid)
	if err != nil {
		s.logger.Errorf("获取进程内存信息失败 PID=%d: %v", pid, err)
		return &models.Response{
			Code:    500,
			Message: fmt.Sprintf("获取进程内存信息失败: %v", err),
			Data:    nil,
		}, nil
	}

	response := &models.Response{
		Code:    200,
		Message: "获取成功",
		Data:    memoryInfo,
	}

	s.logger.Debugf("获取进程内存信息成功 PID=%d，耗时=%v", pid, time.Since(startTime))
	return response, nil
}

// IsProcessRunning 检查进程是否运行
func (s *ProcessService) IsProcessRunning(pid int32) (*models.Response, error) {
	startTime := time.Now()

	running, err := s.manager.IsProcessRunning(pid)
	if err != nil {
		s.logger.Errorf("检查进程运行状态失败 PID=%d: %v", pid, err)
		return &models.Response{
			Code:    500,
			Message: fmt.Sprintf("检查进程运行状态失败: %v", err),
			Data:    nil,
		}, nil
	}

	// 构建响应数据
	data := map[string]interface{}{
		"running": running,
	}

	response := &models.Response{
		Code:    200,
		Message: "获取成功",
		Data:    data,
	}

	s.logger.Debugf("检查进程运行状态成功 PID=%d，运行状态=%v，耗时=%v", pid, running, time.Since(startTime))
	return response, nil
}

// GetProcessChildren 获取子进程
func (s *ProcessService) GetProcessChildren(pid int32) (*models.Response, error) {
	startTime := time.Now()

	children, err := s.manager.GetProcessChildren(pid)
	if err != nil {
		s.logger.Errorf("获取子进程失败 PID=%d: %v", pid, err)
		return &models.Response{
			Code:    500,
			Message: fmt.Sprintf("获取子进程失败: %v", err),
			Data:    nil,
		}, nil
	}

	// 构建响应数据
	data := map[string]interface{}{
		"children": children,
	}

	response := &models.Response{
		Code:    200,
		Message: "获取成功",
		Data:    data,
	}

	s.logger.Debugf("获取子进程成功 PID=%d，共%d个子进程，耗时=%v", pid, len(children), time.Since(startTime))
	return response, nil
}

// GetSystemModules 获取系统模块列表
func (s *ProcessService) GetSystemModules() (*models.Response, error) {
	startTime := time.Now()

	modules, err := s.manager.GetSystemModules()
	if err != nil {
		s.logger.Errorf("获取系统模块列表失败: %v", err)
		return &models.Response{
			Code:    500,
			Message: fmt.Sprintf("获取系统模块列表失败: %v", err),
			Data:    nil,
		}, nil
	}

	// 构建响应数据
	data := map[string]interface{}{
		"data":  modules,
		"count": len(modules),
	}

	response := &models.Response{
		Code:    200,
		Message: "获取列表成功",
		Time:    time.Now(),
		Data:    data,
	}

	s.logger.Debugf("获取系统模块列表成功，共%d个模块，耗时=%v", len(modules), time.Since(startTime))
	return response, nil
}

// GetModuleInfo 获取模块信息
func (s *ProcessService) GetModuleInfo(moduleName string) (*models.Response, error) {
	startTime := time.Now()

	module, err := s.manager.GetModuleInfo(moduleName)
	if err != nil {
		s.logger.Errorf("获取模块信息失败 %s: %v", moduleName, err)
		return &models.Response{
			Code:    404,
			Message: fmt.Sprintf("模块未找到: %v", err),
			Data:    nil,
		}, nil
	}

	response := &models.Response{
		Code:    200,
		Message: "获取信息成功",
		Time:    time.Now(),
		Data:    module,
	}

	s.logger.Debugf("获取模块信息成功 %s，耗时=%v", moduleName, time.Since(startTime))
	return response, nil
}

// SuspendModule 挂起模块
func (s *ProcessService) SuspendModule(moduleName string) (*models.Response, error) {
	startTime := time.Now()

	err := s.manager.SuspendModule(moduleName)
	if err != nil {
		s.logger.Errorf("挂起模块失败 %s: %v", moduleName, err)
		return &models.Response{
			Code:    500,
			Message: fmt.Sprintf("挂起模块失败: %v", err),
			Data:    nil,
		}, nil
	}

	response := &models.Response{
		Code:    200,
		Message: "挂起系统模块成功",
		Time:    time.Now(),
		Data:    map[string]interface{}{},
	}

	s.logger.Infof("挂起模块成功 %s，耗时=%v", moduleName, time.Since(startTime))
	return response, nil
}

// ResumeModule 恢复模块
func (s *ProcessService) ResumeModule(moduleName string) (*models.Response, error) {
	startTime := time.Now()

	err := s.manager.ResumeModule(moduleName)
	if err != nil {
		s.logger.Errorf("恢复模块失败 %s: %v", moduleName, err)
		return &models.Response{
			Code:    500,
			Message: fmt.Sprintf("恢复模块失败: %v", err),
			Data:    nil,
		}, nil
	}

	response := &models.Response{
		Code:    200,
		Message: "恢复系统模块成功",
		Time:    time.Now(),
		Data:    map[string]interface{}{},
	}

	s.logger.Infof("恢复模块成功 %s，耗时=%v", moduleName, time.Since(startTime))
	return response, nil
}

// KillModule 结束模块
func (s *ProcessService) KillModule(moduleName string) (*models.Response, error) {
	startTime := time.Now()

	err := s.manager.KillModule(moduleName)
	if err != nil {
		s.logger.Errorf("结束模块失败 %s: %v", moduleName, err)
		return &models.Response{
			Code:    500,
			Message: fmt.Sprintf("结束模块失败: %v", err),
			Data:    nil,
		}, nil
	}

	response := &models.Response{
		Code:    200,
		Message: "结束系统模块成功",
		Time:    time.Now(),
		Data:    map[string]interface{}{},
	}

	s.logger.Infof("结束模块成功 %s，耗时=%v", moduleName, time.Since(startTime))
	return response, nil
}

// GetProcessStatistics 获取进程统计信息
func (s *ProcessService) GetProcessStatistics() (*models.Response, error) {
	startTime := time.Now()

	stats, err := s.manager.GetSystemMemoryInfo()
	if err != nil {
		s.logger.Errorf("获取进程统计信息失败: %v", err)
		return &models.Response{
			Code:    500,
			Message: fmt.Sprintf("获取进程统计信息失败: %v", err),
			Data:    nil,
		}, nil
	}

	response := &models.Response{
		Code:    200,
		Message: "获取成功",
		Data:    stats,
	}

	s.logger.Debugf("获取进程统计信息成功，耗时=%v", time.Since(startTime))
	return response, nil
}

// EnableProcessMonitoring 启用进程监控
func (s *ProcessService) EnableProcessMonitoring() (*models.Response, error) {
	startTime := time.Now()

	// 这里可以实现真正的进程监控逻辑
	// 例如：启动监控线程、设置钩子等

	response := &models.Response{
		Code:    200,
		Message: "进程监控已启用",
		Data: map[string]interface{}{
			"monitoring_enabled": true,
			"start_time":         startTime,
		},
	}

	s.logger.Infof("进程监控已启用，耗时=%v", time.Since(startTime))
	return response, nil
}

// DisableProcessMonitoring 禁用进程监控
func (s *ProcessService) DisableProcessMonitoring() (*models.Response, error) {
	startTime := time.Now()

	// 这里可以实现真正的进程监控停止逻辑
	// 例如：停止监控线程、移除钩子等

	response := &models.Response{
		Code:    200,
		Message: "进程监控已禁用",
		Data: map[string]interface{}{
			"monitoring_enabled": false,
			"stop_time":          startTime,
		},
	}

	s.logger.Infof("进程监控已禁用，耗时=%v", time.Since(startTime))
	return response, nil
}

// GetMonitoredProcesses 获取监控的进程列表
func (s *ProcessService) GetMonitoredProcesses() (*models.Response, error) {
	startTime := time.Now()

	// 获取所有进程
	processes, err := s.manager.GetProcessList()
	if err != nil {
		s.logger.Errorf("获取进程列表失败: %v", err)
		return &models.Response{
			Code:    500,
			Message: fmt.Sprintf("获取进程列表失败: %v", err),
			Data:    nil,
		}, nil
	}

	// 筛选需要监控的进程（这里可以根据实际需求设置筛选条件）
	var monitoredProcesses []*models.ProcessInfo
	for _, proc := range processes {
		// 示例：监控所有系统进程和用户进程
		if proc.PID > 0 {
			monitoredProcesses = append(monitoredProcesses, proc)
		}
	}

	// 构建响应数据
	data := map[string]interface{}{
		"monitored_processes": monitoredProcesses,
		"count":               len(monitoredProcesses),
		"monitoring_enabled":  true,
	}

	response := &models.Response{
		Code:    200,
		Message: "获取监控进程列表成功",
		Data:    data,
	}

	s.logger.Debugf("获取监控进程列表成功，共%d个进程，耗时=%v", len(monitoredProcesses), time.Since(startTime))
	return response, nil
}
