package handlers

import (
	"strconv"

	"yara-security-service/internal/models"
	"yara-security-service/internal/services"
	"yara-security-service/internal/utils"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// NetworkHandler 网络处理器
type NetworkHandler struct {
	networkService *services.NetworkService
	logger         *logrus.Logger
}

// NewNetworkHandler 创建网络处理器
func NewNetworkHandler(networkService *services.NetworkService, logger *logrus.Logger) *NetworkHandler {
	return &NetworkHandler{
		networkService: networkService,
		logger:         logger,
	}
}

// GetConnections 获取网络连接列表
func (h *NetworkHandler) GetConnections(c *gin.Context) {
	connections, err := h.networkService.GetNetworkConnections()
	if err != nil {
		h.logger.Errorf("获取网络连接失败: %v", err)
		utils.InternalServerErrorResponse(c, "获取网络连接失败")
		return
	}

	utils.ListSuccessResponse(c, connections, len(connections))
}

// GetConnectionsByPID 根据PID获取网络连接
func (h *NetworkHandler) GetConnectionsByPID(c *gin.Context) {
	pidStr := c.Param("pid")
	pid, err := strconv.Atoi(pidStr)
	if err != nil {
		utils.BadRequestResponse(c, "无效的PID参数")
		return
	}

	connections, err := h.networkService.GetConnectionsByPID(int32(pid))
	if err != nil {
		h.logger.Errorf("获取进程网络连接失败: %v", err)
		utils.InternalServerErrorResponse(c, "获取进程网络连接失败")
		return
	}

	utils.ListSuccessResponse(c, connections, len(connections))
}

// GetConnectionsByPort 根据端口获取网络连接
func (h *NetworkHandler) GetConnectionsByPort(c *gin.Context) {
	portStr := c.Param("port")
	port, err := strconv.Atoi(portStr)
	if err != nil {
		utils.BadRequestResponse(c, "无效的端口参数")
		return
	}

	connections, err := h.networkService.GetConnectionsByPort(port)
	if err != nil {
		h.logger.Errorf("获取端口网络连接失败: %v", err)
		utils.InternalServerErrorResponse(c, "获取端口网络连接失败")
		return
	}

	utils.ListSuccessResponse(c, connections, len(connections))
}

// GetConnectionsByIP 根据IP获取网络连接
func (h *NetworkHandler) GetConnectionsByIP(c *gin.Context) {
	ip, _ := utils.GetDecodedPathParam(c, "ip")
	if ip == "" {
		utils.BadRequestResponse(c, "IP地址不能为空")
		return
	}

	connections, err := h.networkService.GetConnectionsByIP(ip)
	if err != nil {
		h.logger.Errorf("获取IP网络连接失败: %v", err)
		utils.InternalServerErrorResponse(c, "获取IP网络连接失败")
		return
	}

	utils.ListSuccessResponse(c, connections, len(connections))
}

// CloseConnection 关闭网络连接
func (h *NetworkHandler) CloseConnection(c *gin.Context) {
	var request models.CloseConnectionRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		utils.ValidationErrorResponse(c, err.Error())
		return
	}

	// 验证至少提供了一个参数
	if request.ConnectionID == "" && request.ProcessName == "" {
		utils.BadRequestResponse(c, "至少需要提供connection_id或process_name中的一个")
		return
	}

	var err error
	if request.ConnectionID != "" {
		// 根据连接ID关闭连接
		err = h.networkService.CloseConnection(request.ConnectionID)
	} else {
		// 根据进程名称关闭连接
		err = h.networkService.CloseConnectionByProcessName(request.ProcessName)
	}

	if err != nil {
		h.logger.Errorf("关闭网络连接失败: %v", err)
		utils.OperationFailedResponse(c, "关闭网络连接", err)
		return
	}

	utils.OperationSuccessResponse(c, "关闭网络连接")
}

// CloseAllConnectionsByProcess 关闭指定进程的所有网络连接
func (h *NetworkHandler) CloseAllConnectionsByProcess(c *gin.Context) {
	var request models.CloseAllConnectionsRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		utils.ValidationErrorResponse(c, err.Error())
		return
	}

	if request.ProcessName == "" {
		utils.BadRequestResponse(c, "进程名称不能为空")
		return
	}

	err := h.networkService.CloseAllConnectionsByProcess(request.ProcessName)
	if err != nil {
		h.logger.Errorf("关闭进程所有网络连接失败: %v", err)
		utils.OperationFailedResponse(c, "关闭进程所有网络连接", err)
		return
	}

	utils.OperationSuccessResponse(c, "关闭进程所有网络连接")
}

// CloseAllConnectionsByPID 关闭指定PID的所有网络连接
func (h *NetworkHandler) CloseAllConnectionsByPID(c *gin.Context) {
	pidStr := c.Param("pid")
	pid, err := strconv.Atoi(pidStr)
	if err != nil {
		utils.BadRequestResponse(c, "无效的PID参数")
		return
	}

	err = h.networkService.CloseAllConnectionsByPID(int32(pid))
	if err != nil {
		h.logger.Errorf("关闭PID所有网络连接失败: %v", err)
		utils.OperationFailedResponse(c, "关闭PID所有网络连接", err)
		return
	}

	utils.OperationSuccessResponse(c, "关闭PID所有网络连接")
}

// ForceCloseProcessConnections 强制关闭指定进程的网络连接
func (h *NetworkHandler) ForceCloseProcessConnections(c *gin.Context) {
	var request models.ForceCloseConnectionsRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		utils.ValidationErrorResponse(c, err.Error())
		return
	}

	if request.ProcessName == "" {
		utils.BadRequestResponse(c, "进程名称不能为空")
		return
	}

	err := h.networkService.ForceCloseProcessConnections(request.ProcessName)
	if err != nil {
		h.logger.Errorf("强制关闭进程网络连接失败: %v", err)
		utils.OperationFailedResponse(c, "强制关闭进程网络连接", err)
		return
	}

	utils.OperationSuccessResponse(c, "强制关闭进程网络连接")
}

// GetInterfaces 获取网络接口列表
func (h *NetworkHandler) GetInterfaces(c *gin.Context) {
	interfaces, err := h.networkService.GetNetworkInterfaces()
	if err != nil {
		h.logger.Errorf("获取网络接口失败: %v", err)
		utils.InternalServerErrorResponse(c, "获取网络接口失败")
		return
	}

	utils.ListSuccessResponse(c, interfaces, len(interfaces))
}

// GetNetworkStats 获取网络统计信息
func (h *NetworkHandler) GetNetworkStats(c *gin.Context) {
	stats, err := h.networkService.GetNetworkStats()
	if err != nil {
		h.logger.Errorf("获取网络统计信息失败: %v", err)
		utils.InternalServerErrorResponse(c, "获取网络统计信息失败")
		return
	}

	utils.StatusSuccessResponse(c, stats)
}

// GetPortsInUse 获取正在使用的端口
func (h *NetworkHandler) GetPortsInUse(c *gin.Context) {
	ports, err := h.networkService.GetPortsInUse()
	if err != nil {
		h.logger.Errorf("获取正在使用的端口失败: %v", err)
		utils.InternalServerErrorResponse(c, "获取正在使用的端口失败")
		return
	}

	utils.ListSuccessResponse(c, ports, len(ports))
}

// GetListeningPorts 获取监听端口
func (h *NetworkHandler) GetListeningPorts(c *gin.Context) {
	ports, err := h.networkService.GetListeningPorts()
	if err != nil {
		h.logger.Errorf("获取监听端口失败: %v", err)
		utils.InternalServerErrorResponse(c, "获取监听端口失败")
		return
	}

	utils.ListSuccessResponse(c, ports, len(ports))
}

// GetEstablishedConnections 获取已建立的连接
func (h *NetworkHandler) GetEstablishedConnections(c *gin.Context) {
	connections, err := h.networkService.GetEstablishedConnections()
	if err != nil {
		h.logger.Errorf("获取已建立连接失败: %v", err)
		utils.InternalServerErrorResponse(c, "获取已建立连接失败")
		return
	}

	utils.ListSuccessResponse(c, connections, len(connections))
}

// GetTCPConnections 获取TCP连接
func (h *NetworkHandler) GetTCPConnections(c *gin.Context) {
	connections, err := h.networkService.GetTCPConnections()
	if err != nil {
		h.logger.Errorf("获取TCP连接失败: %v", err)
		utils.InternalServerErrorResponse(c, "获取TCP连接失败")
		return
	}

	utils.ListSuccessResponse(c, connections, len(connections))
}

// GetUDPConnections 获取UDP连接
func (h *NetworkHandler) GetUDPConnections(c *gin.Context) {
	connections, err := h.networkService.GetUDPConnections()
	if err != nil {
		h.logger.Errorf("获取UDP连接失败: %v", err)
		utils.InternalServerErrorResponse(c, "获取UDP连接失败")
		return
	}

	utils.ListSuccessResponse(c, connections, len(connections))
}

// EnableNetworkMonitoring 启用网络监控
func (h *NetworkHandler) EnableNetworkMonitoring(c *gin.Context) {
	err := h.networkService.EnableMonitoring()
	if err != nil {
		h.logger.Errorf("启用网络监控失败: %v", err)
		utils.OperationFailedResponse(c, "启用网络监控", err)
		return
	}

	utils.OperationSuccessResponse(c, "网络监控已启用")
}

// DisableNetworkMonitoring 禁用网络监控
func (h *NetworkHandler) DisableNetworkMonitoring(c *gin.Context) {
	err := h.networkService.DisableMonitoring()
	if err != nil {
		h.logger.Errorf("禁用网络监控失败: %v", err)
		utils.OperationFailedResponse(c, "禁用网络监控", err)
		return
	}

	utils.OperationSuccessResponse(c, "网络监控已禁用")
}

// GetMonitoredConnections 获取监控的连接
func (h *NetworkHandler) GetMonitoredConnections(c *gin.Context) {
	connections := h.networkService.GetMonitoredConnections()
	utils.ListSuccessResponse(c, connections, len(connections))
}

// GetConnectionHistory 获取连接历史
func (h *NetworkHandler) GetConnectionHistory(c *gin.Context) {
	history := h.networkService.GetConnectionHistory()
	utils.ListSuccessResponse(c, history, len(history))
}

// GetConnectionByID 根据连接ID获取连接信息
func (h *NetworkHandler) GetConnectionByID(c *gin.Context) {
	connectionID := c.Param("id")
	if connectionID == "" {
		utils.BadRequestResponse(c, "连接ID不能为空")
		return
	}

	connections, err := h.networkService.GetConnectionByID(connectionID)
	if err != nil {
		h.logger.Errorf("获取连接信息失败: %v", err)
		utils.InternalServerErrorResponse(c, "获取连接信息失败")
		return
	}

	utils.ListSuccessResponse(c, connections, len(connections))
}

// IsPortInUse 检查指定端口是否被占用
func (h *NetworkHandler) IsPortInUse(c *gin.Context) {
	portStr := c.Param("port")
	port, err := strconv.Atoi(portStr)
	if err != nil {
		utils.BadRequestResponse(c, "无效的端口参数")
		return
	}

	inUse, err := h.networkService.IsPortInUse(port)
	if err != nil {
		h.logger.Errorf("检查端口占用状态失败: %v", err)
		utils.InternalServerErrorResponse(c, "检查端口占用状态失败")
		return
	}

	utils.SuccessResponse(c, "端口占用状态检查完成", gin.H{
		"port":   port,
		"in_use": inUse,
	})
}
