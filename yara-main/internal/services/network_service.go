package services

import (
	"fmt"
	"strings"
	"sync"
	"time"

	"yara-security-service/internal/models"
	"yara-security-service/pkg/network"

	"github.com/sirupsen/logrus"
)

// NetworkService 网络服务
type NetworkService struct {
	manager *network.Manager
	logger  *logrus.Logger
	mu      sync.RWMutex

	// 监控状态
	monitoringEnabled    bool
	monitoredConnections map[string]*models.NetworkConnection
	connectionHistory    []*models.NetworkConnection
	stopMonitor          chan bool
}

// NewNetworkService 创建网络服务
func NewNetworkService(manager *network.Manager, logger *logrus.Logger) *NetworkService {
	return &NetworkService{
		manager:              manager,
		logger:               logger,
		monitoringEnabled:    false,
		monitoredConnections: make(map[string]*models.NetworkConnection),
		connectionHistory:    make([]*models.NetworkConnection, 0),
		stopMonitor:          make(chan bool),
	}
}

// GetNetworkConnections 获取网络连接
func (s *NetworkService) GetNetworkConnections() ([]*models.NetworkConnection, error) {
	connections, err := s.manager.GetNetworkConnections()
	if err != nil {
		return nil, fmt.Errorf("获取网络连接失败: %w", err)
	}

	s.logger.Debugf("获取网络连接成功: %d 个连接", len(connections))
	return connections, nil
}

// GetTCPConnections 获取TCP连接
func (s *NetworkService) GetTCPConnections() ([]*models.NetworkConnection, error) {
	connections, err := s.manager.GetTCPConnections()
	if err != nil {
		return nil, fmt.Errorf("获取TCP连接失败: %w", err)
	}

	s.logger.Debugf("获取TCP连接成功: %d 个连接", len(connections))
	return connections, nil
}

// GetUDPConnections 获取UDP连接
func (s *NetworkService) GetUDPConnections() ([]*models.NetworkConnection, error) {
	connections, err := s.manager.GetUDPConnections()
	if err != nil {
		return nil, fmt.Errorf("获取UDP连接失败: %w", err)
	}

	s.logger.Debugf("获取UDP连接成功: %d 个连接", len(connections))
	return connections, nil
}

// GetConnectionsByPID 根据PID获取连接
func (s *NetworkService) GetConnectionsByPID(pid int32) ([]*models.NetworkConnection, error) {
	connections, err := s.manager.GetConnectionsByPID(pid)
	if err != nil {
		return nil, fmt.Errorf("根据PID获取连接失败: %w", err)
	}

	s.logger.Debugf("根据PID获取连接成功: PID=%d, 连接数=%d", pid, len(connections))
	return connections, nil
}

// CloseConnection 关闭连接
func (s *NetworkService) CloseConnection(connectionID string) error {
	err := s.manager.CloseConnection(connectionID)
	if err != nil {
		return fmt.Errorf("关闭连接失败: %w", err)
	}

	s.logger.Infof("关闭连接成功: %s", connectionID)
	return nil
}

// CloseConnectionByProcessName 根据进程名称关闭连接
func (s *NetworkService) CloseConnectionByProcessName(processName string) error {
	// 获取所有网络连接
	connections, err := s.manager.GetNetworkConnections()
	if err != nil {
		return fmt.Errorf("获取网络连接失败: %w", err)
	}

	// 查找匹配进程名称的连接
	var closedCount int
	for _, conn := range connections {
		if strings.EqualFold(conn.ProcessName, processName) {
			// 关闭匹配的连接
			err := s.manager.CloseConnection(conn.ID)
			if err != nil {
				s.logger.Warnf("关闭连接失败: %s, 错误: %v", conn.ID, err)
				continue
			}
			closedCount++
		}
	}

	if closedCount == 0 {
		return fmt.Errorf("未找到进程名称为 %s 的连接", processName)
	}

	s.logger.Infof("成功关闭进程 %s 的 %d 个连接", processName, closedCount)
	return nil
}

// CloseAllConnectionsByProcess 关闭指定进程的所有网络连接（强化版）
func (s *NetworkService) CloseAllConnectionsByProcess(processName string) error {
	s.logger.Infof("开始关闭进程 %s 的所有网络连接", processName)

	err := s.manager.CloseAllConnectionsByProcess(processName)
	if err != nil {
		return fmt.Errorf("关闭进程 %s 的所有网络连接失败: %w", processName, err)
	}

	s.logger.Infof("成功关闭进程 %s 的所有网络连接", processName)
	return nil
}

// CloseAllConnectionsByPID 关闭指定PID的所有网络连接
func (s *NetworkService) CloseAllConnectionsByPID(pid int32) error {
	s.logger.Infof("开始关闭进程 %d 的所有网络连接", pid)

	err := s.manager.CloseAllConnectionsByPID(pid)
	if err != nil {
		return fmt.Errorf("关闭进程 %d 的所有网络连接失败: %w", pid, err)
	}

	s.logger.Infof("成功关闭进程 %d 的所有网络连接", pid)
	return nil
}

// ForceCloseProcessConnections 强制关闭进程的网络连接
func (s *NetworkService) ForceCloseProcessConnections(processName string) error {
	s.logger.Infof("开始强制关闭进程 %s 的网络连接", processName)

	// 获取所有网络连接
	connections, err := s.manager.GetNetworkConnections()
	if err != nil {
		return fmt.Errorf("获取网络连接失败: %w", err)
	}

	// 查找匹配进程名称的连接
	var processConnections []*models.NetworkConnection
	for _, conn := range connections {
		if strings.EqualFold(conn.ProcessName, processName) {
			processConnections = append(processConnections, conn)
		}
	}

	if len(processConnections) == 0 {
		return fmt.Errorf("未找到进程名称为 %s 的网络连接", processName)
	}

	s.logger.Infof("找到进程 %s 的 %d 个网络连接，开始强制关闭", processName, len(processConnections))

	// 获取进程信息
	processes, err := s.manager.GetProcessesByName(processName)
	if err != nil {
		return fmt.Errorf("获取进程信息失败: %w", err)
	}

	// 尝试强制关闭每个进程的网络连接
	var successCount int
	for _, proc := range processes {
		if err := s.manager.ForceCloseProcessConnections(proc.Pid); err == nil {
			successCount++
			s.logger.Infof("强制关闭进程 %d 的网络连接成功", proc.Pid)
		} else {
			s.logger.Warnf("强制关闭进程 %d 的网络连接失败: %v", proc.Pid, err)
		}
	}

	if successCount == 0 {
		return fmt.Errorf("所有进程的网络连接强制关闭都失败")
	}

	s.logger.Infof("成功强制关闭进程 %s 的网络连接: %d/%d 个进程", processName, successCount, len(processes))
	return nil
}

// GetNetworkInterfaces 获取网络接口
func (s *NetworkService) GetNetworkInterfaces() ([]*models.NetworkInterface, error) {
	interfaces, err := s.manager.GetNetworkInterfaces()
	if err != nil {
		return nil, fmt.Errorf("获取网络接口失败: %w", err)
	}

	// 转换为具体的结构体类型
	var result []*models.NetworkInterface
	for _, iface := range interfaces {
		// 安全地提取地址信息
		var addresses []string
		if addrs, ok := iface["addresses"].([]string); ok {
			addresses = addrs
		}

		// 安全地提取MTU信息
		var mtu int
		if mtuVal, ok := iface["mtu"]; ok {
			if mtuInt, ok := mtuVal.(int); ok {
				mtu = mtuInt
			}
		}

		// 安全地提取接口状态
		var isUp bool
		if upVal, ok := iface["is_up"]; ok {
			if upBool, ok := upVal.(bool); ok {
				isUp = upBool
			}
		}

		var isLoopback bool
		if loopbackVal, ok := iface["is_loopback"]; ok {
			if loopbackBool, ok := loopbackVal.(bool); ok {
				isLoopback = loopbackBool
			}
		}

		// 创建网络接口对象
		networkInterface := &models.NetworkInterface{
			Name:       iface["name"].(string),
			Addresses:  addresses,
			IsUp:       isUp,
			IsLoopback: isLoopback,
			MTU:        mtu,
			Speed:      0, // 简化实现
		}

		// 填充扩展字段
		if indexVal, ok := iface["index"]; ok {
			if indexInt, ok := indexVal.(int); ok {
				networkInterface.Index = indexInt
			}
		}

		if ipv4Addrs, ok := iface["ipv4_addresses"]; ok {
			if ipv4List, ok := ipv4Addrs.([]string); ok {
				networkInterface.IPv4Addresses = ipv4List
			}
		}

		if ipv6Addrs, ok := iface["ipv6_addresses"]; ok {
			if ipv6List, ok := ipv6Addrs.([]string); ok {
				networkInterface.IPv6Addresses = ipv6List
			}
		}

		if macAddr, ok := iface["mac_address"]; ok {
			if macStr, ok := macAddr.(string); ok {
				networkInterface.MACAddress = macStr
			}
		}

		if isMulticast, ok := iface["is_multicast"]; ok {
			if multicastBool, ok := isMulticast.(bool); ok {
				networkInterface.IsMulticast = multicastBool
			}
		}

		if isBroadcast, ok := iface["is_broadcast"]; ok {
			if broadcastBool, ok := isBroadcast.(bool); ok {
				networkInterface.IsBroadcast = broadcastBool
			}
		}

		if isPointToPoint, ok := iface["is_pointtopoint"]; ok {
			if p2pBool, ok := isPointToPoint.(bool); ok {
				networkInterface.IsPointToPoint = p2pBool
			}
		}

		if flags, ok := iface["flags"]; ok {
			if flagsStr, ok := flags.(string); ok {
				networkInterface.Flags = flagsStr
			}
		}

		if hardwareAddr, ok := iface["hardware_addr"]; ok {
			if hwAddrStr, ok := hardwareAddr.(string); ok {
				networkInterface.HardwareAddr = hwAddrStr
			}
		}

		// 添加扩展信息到元数据字段（如果模型支持）
		if extendedInfo, ok := iface["windows_info"]; ok {
			if windowsInfo, ok := extendedInfo.(map[string]interface{}); ok {
				networkInterface.WindowsInfo = windowsInfo
				s.logger.Debugf("网络接口 %s 包含Windows特定信息: %+v", networkInterface.Name, windowsInfo)
			}
		}

		// 添加错误信息（如果有）
		if errorMsg, ok := iface["error"]; ok {
			if errorStr, ok := errorMsg.(string); ok {
				networkInterface.Error = errorStr
			}
		}

		result = append(result, networkInterface)
	}

	s.logger.Debugf("获取网络接口成功: %d 个接口", len(result))
	return result, nil
}

// GetNetworkStats 获取网络统计
func (s *NetworkService) GetNetworkStats() (*models.NetworkStats, error) {
	stats, err := s.manager.GetNetworkStats()
	if err != nil {
		return nil, fmt.Errorf("获取网络统计失败: %w", err)
	}

	// 转换为具体的结构体类型
	networkStats := &models.NetworkStats{
		Timestamp:    time.Now(),
		Monitored:    true,                        // 简化实现
		Interfaces:   []models.NetworkInterface{}, // 简化实现
		TotalBytes:   int64(stats["bytes_sent"].(uint64) + stats["bytes_recv"].(uint64)),
		TotalPackets: int64(stats["packets_sent"].(uint64) + stats["packets_recv"].(uint64)),
		ThreatCount:  0, // 简化实现
	}

	s.logger.Debug("获取网络统计成功")
	return networkStats, nil
}

// GetConnectionsByPort 根据端口获取连接
func (s *NetworkService) GetConnectionsByPort(port int) ([]*models.NetworkConnection, error) {
	connections, err := s.manager.GetConnectionsByPort(port)
	if err != nil {
		return nil, fmt.Errorf("根据端口获取连接失败: %w", err)
	}

	s.logger.Debugf("根据端口获取连接成功: 端口=%d, 连接数=%d", port, len(connections))
	return connections, nil
}

// GetConnectionsByIP 根据IP地址获取连接
func (s *NetworkService) GetConnectionsByIP(ip string) ([]*models.NetworkConnection, error) {
	connections, err := s.manager.GetConnectionsByIP(ip)
	if err != nil {
		return nil, fmt.Errorf("根据IP地址获取连接失败: %w", err)
	}

	s.logger.Debugf("根据IP地址获取连接成功: IP=%s, 连接数=%d", ip, len(connections))
	return connections, nil
}

// IsPortInUse 检查端口是否被使用
func (s *NetworkService) IsPortInUse(port int) (bool, error) {
	inUse, err := s.manager.IsPortInUse(port)
	if err != nil {
		return false, fmt.Errorf("检查端口使用状态失败: %w", err)
	}

	s.logger.Debugf("检查端口使用状态: 端口=%d, 使用中=%v", port, inUse)
	return inUse, nil
}

// GetListeningPorts 获取监听端口
func (s *NetworkService) GetListeningPorts() ([]int, error) {
	ports, err := s.manager.GetListeningPorts()
	if err != nil {
		return nil, fmt.Errorf("获取监听端口失败: %w", err)
	}

	s.logger.Debugf("获取监听端口成功: %d 个端口", len(ports))
	return ports, nil
}

// GetPortsInUse 获取正在使用的端口
func (s *NetworkService) GetPortsInUse() ([]int, error) {
	// 获取所有网络连接
	connections, err := s.manager.GetNetworkConnections()
	if err != nil {
		return nil, fmt.Errorf("获取网络连接失败: %w", err)
	}

	// 收集所有使用的端口
	portMap := make(map[int]bool)
	var portsInUse []int

	for _, conn := range connections {
		// 添加本地端口
		if conn.LocalPort > 0 {
			port := conn.LocalPort
			if !portMap[port] {
				portMap[port] = true
				portsInUse = append(portsInUse, port)
			}
		}

		// 添加远程端口
		if conn.RemotePort > 0 {
			port := conn.RemotePort
			if !portMap[port] {
				portMap[port] = true
				portsInUse = append(portsInUse, port)
			}
		}
	}

	s.logger.Debugf("获取正在使用的端口成功: %d 个端口", len(portsInUse))
	return portsInUse, nil
}

// GetEstablishedConnections 获取已建立的连接
func (s *NetworkService) GetEstablishedConnections() ([]*models.NetworkConnection, error) {
	connections, err := s.manager.GetEstablishedConnections()
	if err != nil {
		return nil, fmt.Errorf("获取已建立连接失败: %w", err)
	}

	s.logger.Debugf("获取已建立连接成功: %d 个连接", len(connections))
	return connections, nil
}

// EnableMonitoring 启用连接监控
func (s *NetworkService) EnableMonitoring() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.monitoringEnabled {
		return fmt.Errorf("监控已启用")
	}

	s.monitoringEnabled = true
	s.logger.Info("网络连接监控已启用")

	// 启动后台监控进程
	go s.monitorConnections()

	return nil
}

// DisableMonitoring 禁用连接监控
func (s *NetworkService) DisableMonitoring() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if !s.monitoringEnabled {
		return fmt.Errorf("监控未启用")
	}

	s.monitoringEnabled = false
	s.stopMonitor <- true
	s.logger.Info("网络连接监控已禁用")
	return nil
}

// monitorConnections 后台监控连接
func (s *NetworkService) monitorConnections() {
	ticker := time.NewTicker(10 * time.Second) // 每10秒检查一次
	defer ticker.Stop()

	for {
		select {
		case <-s.stopMonitor:
			return
		case <-ticker.C:
			s.updateMonitoredConnections()
		}
	}
}

// updateMonitoredConnections 更新监控的连接
func (s *NetworkService) updateMonitoredConnections() {
	connections, err := s.GetNetworkConnections()
	if err != nil {
		s.logger.Warnf("获取网络连接失败: %v", err)
		return
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	// 清空当前监控的连接
	s.monitoredConnections = make(map[string]*models.NetworkConnection)

	// 添加所有当前连接到监控列表
	for _, conn := range connections {
		key := fmt.Sprintf("%s:%d-%s:%d", conn.LocalAddr, conn.LocalPort, conn.RemoteAddr, conn.RemotePort)
		s.monitoredConnections[key] = conn
		s.connectionHistory = append(s.connectionHistory, conn)
	}

	// 保持历史记录在合理范围内
	if len(s.connectionHistory) > 1000 {
		s.connectionHistory = s.connectionHistory[len(s.connectionHistory)-500:]
	}

	s.logger.Debugf("更新监控连接: %d 个连接", len(connections))
}

// AddMonitoredConnection 添加监控连接
func (s *NetworkService) AddMonitoredConnection(connection *models.NetworkConnection) {
	s.mu.Lock()
	defer s.mu.Unlock()

	key := fmt.Sprintf("%s:%d-%s:%d", connection.LocalAddr, connection.LocalPort, connection.RemoteAddr, connection.RemotePort)
	s.monitoredConnections[key] = connection
	s.connectionHistory = append(s.connectionHistory, connection)

	// 保持历史记录在合理范围内
	if len(s.connectionHistory) > 1000 {
		s.connectionHistory = s.connectionHistory[len(s.connectionHistory)-500:]
	}

	s.logger.Debugf("添加监控连接: %s", key)
}

// GetMonitoredConnections 获取监控的连接
func (s *NetworkService) GetMonitoredConnections() []*models.NetworkConnection {
	s.mu.RLock()
	defer s.mu.RUnlock()

	connections := make([]*models.NetworkConnection, 0, len(s.monitoredConnections))
	for _, conn := range s.monitoredConnections {
		connections = append(connections, conn)
	}

	return connections
}

// GetConnectionHistory 获取连接历史
func (s *NetworkService) GetConnectionHistory() []*models.NetworkConnection {
	s.mu.RLock()
	defer s.mu.RUnlock()

	history := make([]*models.NetworkConnection, len(s.connectionHistory))
	copy(history, s.connectionHistory)
	return history
}

// GetNetworkStatistics 获取网络统计信息
func (s *NetworkService) GetNetworkStatistics() map[string]interface{} {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return map[string]interface{}{
		"monitoring_enabled":    s.monitoringEnabled,
		"monitored_connections": len(s.monitoredConnections),
		"connection_history":    len(s.connectionHistory),
		"last_updated":          time.Now(),
	}
}

// ClearConnectionHistory 清空连接历史
func (s *NetworkService) ClearConnectionHistory() {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.connectionHistory = make([]*models.NetworkConnection, 0)
	s.logger.Info("连接历史已清空")
}

// GetConnectionByID 根据ID获取连接
func (s *NetworkService) GetConnectionByID(connectionID string) ([]*models.NetworkConnection, error) {
	connections, err := s.GetNetworkConnections()
	if err != nil {
		return nil, err
	}

	var matchingConnections []*models.NetworkConnection
	for _, conn := range connections {
		if conn.ID == connectionID {
			matchingConnections = append(matchingConnections, conn)
		}
	}

	if len(matchingConnections) == 0 {
		return nil, fmt.Errorf("未找到连接: %s", connectionID)
	}

	return matchingConnections, nil
}

// BlockConnection 阻止连接
func (s *NetworkService) BlockConnection(connectionID string) error {
	// 这里可以实现阻止连接的逻辑
	// 例如通过防火墙规则阻止特定连接
	s.logger.Infof("阻止连接: %s", connectionID)
	return nil
}

// AllowConnection 允许连接
func (s *NetworkService) AllowConnection(connectionID string) error {
	// 这里可以实现允许连接的逻辑
	s.logger.Infof("允许连接: %s", connectionID)
	return nil
}

// isInterfaceUp 检查接口是否启用
func (s *NetworkService) isInterfaceUp(iface map[string]interface{}) bool {
	flags, ok := iface["flags"].(string)
	if !ok {
		return false
	}
	return strings.Contains(flags, "up")
}

// isLoopbackInterface 检查是否为回环接口
func (s *NetworkService) isLoopbackInterface(iface map[string]interface{}) bool {
	name, ok := iface["name"].(string)
	if !ok {
		return false
	}
	return strings.Contains(strings.ToLower(name), "loopback") ||
		strings.Contains(strings.ToLower(name), "lo")
}

// isMulticastInterface 检查是否为多播接口
func (s *NetworkService) isMulticastInterface(iface map[string]interface{}) bool {
	flags, ok := iface["flags"].(string)
	if !ok {
		return false
	}
	return strings.Contains(flags, "multicast")
}

// isBroadcastInterface 检查是否为广播接口
func (s *NetworkService) isBroadcastInterface(iface map[string]interface{}) bool {
	flags, ok := iface["flags"].(string)
	if !ok {
		return false
	}
	return strings.Contains(flags, "broadcast")
}
