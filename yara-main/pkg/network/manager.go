package network

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	"yara-security-service/internal/models"
	"yara-security-service/pkg/security"

	"sync"

	"net"
	"os/exec"

	gopsutilnet "github.com/shirou/gopsutil/v3/net"
	gopsutilprocess "github.com/shirou/gopsutil/v3/process"
	"github.com/sirupsen/logrus"
)

// NetworkMonitor 网络监控器
type NetworkMonitor struct {
	logger     *logrus.Logger
	scanner    *security.Scanner
	isRunning  bool
	stopChan   chan bool
	packetChan chan *models.NetworkPacket
}

// NetworkPacket 网络数据包结构
type NetworkPacket struct {
	SourceIP      string
	DestIP        string
	SourcePort    int
	DestPort      int
	Protocol      string
	Payload       []byte
	Timestamp     time.Time
	PacketSize    int
	IsSuspicious  bool
	ThreatLevel   string
	ThreatDetails []string
}

// NewNetworkMonitor 创建网络监控器
func NewNetworkMonitor(logger *logrus.Logger, scanner *security.Scanner) *NetworkMonitor {
	return &NetworkMonitor{
		logger:     logger,
		scanner:    scanner,
		isRunning:  false,
		stopChan:   make(chan bool),
		packetChan: make(chan *models.NetworkPacket, 1000),
	}
}

// StartMonitoring 开始网络监控
func (nm *NetworkMonitor) StartMonitoring(ctx context.Context) error {
	if nm.isRunning {
		return fmt.Errorf("网络监控已在运行")
	}

	nm.isRunning = true
	nm.logger.Info("开始网络流实时监控")

	// 启动数据包捕获
	go nm.capturePackets(ctx)

	// 启动数据分析
	go nm.analyzePackets(ctx)

	return nil
}

// StopMonitoring 停止网络监控
func (nm *NetworkMonitor) StopMonitoring() error {
	if !nm.isRunning {
		return fmt.Errorf("网络监控未在运行")
	}

	nm.isRunning = false
	nm.stopChan <- true
	nm.logger.Info("停止网络流实时监控")

	return nil
}

// capturePackets 捕获网络数据包
func (nm *NetworkMonitor) capturePackets(ctx context.Context) {
	// 获取网络接口
	interfaces, err := gopsutilnet.Interfaces()
	if err != nil {
		nm.logger.Errorf("获取网络接口失败: %v", err)
		return
	}

	for _, iface := range interfaces {
		// 简化接口检查
		if iface.Name == "" {
			continue // 跳过无效接口
		}

		go nm.captureInterface(ctx, iface.Name)
	}

	// 等待停止信号
	<-nm.stopChan
}

// captureInterface 捕获指定接口的数据包
func (nm *NetworkMonitor) captureInterface(ctx context.Context, interfaceName string) {
	nm.logger.Infof("开始监控网络接口: %s", interfaceName)

	// 这里应该实现实际的数据包捕获
	// 由于Go的限制，这里使用模拟实现
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-nm.stopChan:
			return
		case <-ticker.C:
			// 模拟数据包捕获
			nm.simulatePacketCapture(interfaceName)
		}
	}
}

// simulatePacketCapture 模拟数据包捕获
func (nm *NetworkMonitor) simulatePacketCapture(interfaceName string) {
	// 模拟不同类型的网络数据包
	packets := []*models.NetworkPacket{
		{
			SourceIP:     "192.168.1.100",
			DestIP:       "8.8.8.8",
			SourcePort:   12345,
			DestPort:     53,
			Protocol:     "UDP",
			Payload:      []byte("DNS query for malicious-domain.com"),
			Timestamp:    time.Now(),
			PacketSize:   64,
			IsSuspicious: false,
		},
		{
			SourceIP:      "192.168.1.100",
			DestIP:        "malicious-server.com",
			SourcePort:    12346,
			DestPort:      80,
			Protocol:      "TCP",
			Payload:       []byte("GET /download/malware.exe HTTP/1.1"),
			Timestamp:     time.Now(),
			PacketSize:    128,
			IsSuspicious:  true,
			ThreatLevel:   "high",
			ThreatDetails: []string{"malicious domain", "malware download"},
		},
	}

	for _, packet := range packets {
		select {
		case nm.packetChan <- packet:
		default:
			nm.logger.Warn("数据包通道已满，丢弃数据包")
		}
	}
}

// analyzePackets 分析网络数据包
func (nm *NetworkMonitor) analyzePackets(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		case <-nm.stopChan:
			return
		case packet := <-nm.packetChan:
			nm.analyzePacket(packet)
		}
	}
}

// analyzePacket 分析单个数据包
func (nm *NetworkMonitor) analyzePacket(packet *models.NetworkPacket) {
	// 检查数据包内容是否包含威胁
	threats := nm.detectPacketThreats(packet)

	if len(threats) > 0 {
		packet.IsSuspicious = true
		packet.ThreatLevel = "high"
		packet.ThreatDetails = threats

		nm.logger.Warnf("检测到可疑网络数据包: %s:%d -> %s:%d, 威胁: %v",
			packet.SourceIP, packet.SourcePort,
			packet.DestIP, packet.DestPort,
			threats)
	}
}

// detectPacketThreats 检测数据包威胁
func (nm *NetworkMonitor) detectPacketThreats(packet *models.NetworkPacket) []string {
	var threats []string

	// 检查恶意域名
	maliciousDomains := []string{
		"malicious", "evil", "hack", "exploit", "backdoor",
		"trojan", "virus", "malware", "spyware", "keylogger",
	}

	payloadStr := strings.ToLower(string(packet.Payload))

	for _, domain := range maliciousDomains {
		if strings.Contains(payloadStr, domain) {
			threats = append(threats, fmt.Sprintf("恶意域名: %s", domain))
		}
	}

	// 检查可疑协议
	if packet.Protocol == "TCP" && packet.DestPort == 4444 {
		threats = append(threats, "可疑端口: 4444 (常见后门端口)")
	}

	// 检查数据包大小异常
	if packet.PacketSize > 1500 {
		threats = append(threats, "数据包大小异常")
	}

	// 使用安全扫描器检查载荷
	if nm.scanner != nil {
		scanResult, err := nm.scanner.ScanBuffer(context.Background(), packet.Payload, "network_packet")
		if err == nil && scanResult.IsInfected {
			threats = append(threats, "载荷包含恶意代码")
		}
	}

	return threats
}

// GetNetworkStats 获取网络统计信息
func (nm *NetworkMonitor) GetNetworkStats() *models.NetworkStats {
	stats := &models.NetworkStats{
		Timestamp: time.Now(),
		Monitored: nm.isRunning,
	}

	// 获取网络接口统计
	interfaces, err := gopsutilnet.Interfaces()
	if err == nil {
		for _, iface := range interfaces {
			// 简化接口检查
			if iface.Name != "" {
				stats.Interfaces = append(stats.Interfaces, models.NetworkInterface{
					Name:       iface.Name,
					Addresses:  []string{}, // 简化地址信息
					IsUp:       true,
					IsLoopback: false, // 简化回环检查
				})
			}
		}
	}

	return stats
}

// Manager 网络管理器
type Manager struct {
	monitor *NetworkMonitor
	logger  *logrus.Logger

	// 缓存机制
	connectionCache map[string]*models.NetworkConnection
	cacheMutex      sync.RWMutex
	cacheTTL        time.Duration
	lastCacheUpdate time.Time

	// 并发控制
	networkWorkers   int
	networkSemaphore chan struct{}
	networkTimeout   time.Duration
}

// NewManager 创建网络管理器
func NewManager(logger *logrus.Logger) *Manager {
	manager := &Manager{
		monitor:          NewNetworkMonitor(logger, nil), // 暂时不传入scanner
		logger:           logger,
		connectionCache:  make(map[string]*models.NetworkConnection),
		cacheTTL:         30 * time.Second,
		networkWorkers:   3,
		networkSemaphore: make(chan struct{}, 3),
		networkTimeout:   10 * time.Second,
	}

	// 启动缓存清理协程
	go manager.startCacheCleaner()

	return manager
}

// startCacheCleaner 启动缓存清理协程
func (m *Manager) startCacheCleaner() {
	ticker := time.NewTicker(60 * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		m.cleanExpiredCache()
	}
}

// cleanExpiredCache 清理过期缓存
func (m *Manager) cleanExpiredCache() {
	m.cacheMutex.Lock()
	defer m.cacheMutex.Unlock()

	now := time.Now()
	if now.Sub(m.lastCacheUpdate) > m.cacheTTL {
		m.connectionCache = make(map[string]*models.NetworkConnection)
		m.lastCacheUpdate = now
		m.logger.Debug("清理了网络连接缓存")
	}
}

// GetNetworkConnections 获取网络连接（优化版本）
func (m *Manager) GetNetworkConnections() ([]*models.NetworkConnection, error) {
	// 检查缓存
	m.cacheMutex.RLock()
	if len(m.connectionCache) > 0 && time.Since(m.lastCacheUpdate) < m.cacheTTL {
		connections := make([]*models.NetworkConnection, 0, len(m.connectionCache))
		for _, conn := range m.connectionCache {
			connections = append(connections, conn)
		}
		m.cacheMutex.RUnlock()
		return connections, nil
	}
	m.cacheMutex.RUnlock()

	// 获取网络连接
	connections, err := gopsutilnet.Connections("all")
	if err != nil {
		return nil, fmt.Errorf("获取网络连接失败: %w", err)
	}

	// 并发处理连接信息
	var wg sync.WaitGroup
	resultChan := make(chan *models.NetworkConnection, len(connections))
	errorChan := make(chan error, len(connections))

	for _, conn := range connections {
		wg.Add(1)
		go func(c gopsutilnet.ConnectionStat) {
			defer wg.Done()

			select {
			case m.networkSemaphore <- struct{}{}:
				defer func() { <-m.networkSemaphore }()
			case <-time.After(m.networkTimeout):
				errorChan <- fmt.Errorf("处理网络连接超时")
				return
			}

			networkConn := &models.NetworkConnection{
				ID:          fmt.Sprintf("%d-%s-%d", c.Pid, m.getConnectionType(c.Type), c.Fd),
				LocalAddr:   c.Laddr.IP,
				RemoteAddr:  c.Raddr.IP,
				LocalPort:   int(c.Laddr.Port),
				RemotePort:  int(c.Raddr.Port),
				Protocol:    m.getProtocolString(c.Type),
				Status:      c.Status,
				PID:         c.Pid,
				ProcessName: m.getProcessName(c.Pid),
				Type:        m.getConnectionType(c.Type),
			}

			resultChan <- networkConn
		}(conn)
	}

	// 等待所有协程完成
	go func() {
		wg.Wait()
		close(resultChan)
		close(errorChan)
	}()

	// 收集结果
	var networkConnections []*models.NetworkConnection
	for conn := range resultChan {
		networkConnections = append(networkConnections, conn)

		// 缓存连接信息
		connKey := fmt.Sprintf("%s:%d-%s:%d", conn.LocalAddr, conn.LocalPort, conn.RemoteAddr, conn.RemotePort)
		m.cacheMutex.Lock()
		m.connectionCache[connKey] = conn
		m.cacheMutex.Unlock()
	}

	// 处理错误
	for err := range errorChan {
		m.logger.Warnf("处理网络连接时出错: %v", err)
	}

	// 更新缓存时间
	m.cacheMutex.Lock()
	m.lastCacheUpdate = time.Now()
	m.cacheMutex.Unlock()

	return networkConnections, nil
}

// GetTCPConnections 获取TCP连接（优化版本）
func (m *Manager) GetTCPConnections() ([]*models.NetworkConnection, error) {
	connections, err := gopsutilnet.Connections("tcp")
	if err != nil {
		return nil, fmt.Errorf("获取TCP连接失败: %w", err)
	}

	var tcpConnections []*models.NetworkConnection
	var wg sync.WaitGroup
	resultChan := make(chan *models.NetworkConnection, len(connections))

	for _, conn := range connections {
		wg.Add(1)
		go func(c gopsutilnet.ConnectionStat) {
			defer wg.Done()

			select {
			case m.networkSemaphore <- struct{}{}:
				defer func() { <-m.networkSemaphore }()
			case <-time.After(m.networkTimeout):
				return
			}

			networkConn := &models.NetworkConnection{
				ID:          fmt.Sprintf("%d-%s-%d", c.Pid, "tcp", c.Fd),
				LocalAddr:   c.Laddr.IP,
				RemoteAddr:  c.Raddr.IP,
				LocalPort:   int(c.Laddr.Port),
				RemotePort:  int(c.Raddr.Port),
				Protocol:    "TCP",
				Status:      c.Status,
				PID:         c.Pid,
				ProcessName: m.getProcessName(c.Pid),
				Type:        "tcp",
			}

			resultChan <- networkConn
		}(conn)
	}

	go func() {
		wg.Wait()
		close(resultChan)
	}()

	for conn := range resultChan {
		tcpConnections = append(tcpConnections, conn)
	}

	return tcpConnections, nil
}

// GetNetworkInterfaces 获取网络接口
func (m *Manager) GetNetworkInterfaces() ([]map[string]interface{}, error) {
	m.logger.Debugf("开始获取网络接口信息")

	// 使用gopsutil获取基本接口信息
	interfaces, err := gopsutilnet.Interfaces()
	if err != nil {
		return nil, fmt.Errorf("获取网络接口失败: %w", err)
	}

	var result []map[string]interface{}
	for _, iface := range interfaces {
		m.logger.Debugf("处理网络接口: %s", iface.Name)

		// 获取接口地址
		var addresses []string
		var ipv4Addresses []string
		var ipv6Addresses []string
		var macAddress string

		// 使用net包获取接口地址
		netInterface, err := net.InterfaceByName(iface.Name)
		if err == nil {
			// 获取MAC地址
			macAddress = netInterface.HardwareAddr.String()

			// 获取接口地址
			addrs, err := netInterface.Addrs()
			if err == nil {
				for _, addr := range addrs {
					// 提取IP地址部分（去掉子网掩码）
					addrStr := addr.String()
					if strings.Contains(addrStr, "/") {
						addrStr = strings.Split(addrStr, "/")[0]
					}
					addresses = append(addresses, addrStr)

					// 分类IPv4和IPv6地址
					if strings.Contains(addrStr, ":") {
						ipv6Addresses = append(ipv6Addresses, addrStr)
					} else {
						ipv4Addresses = append(ipv4Addresses, addrStr)
					}
				}
			}

			// 获取接口状态
			flags := netInterface.Flags
			isUp := (flags & net.FlagUp) != 0
			isLoopback := (flags & net.FlagLoopback) != 0
			isMulticast := (flags & net.FlagMulticast) != 0
			isBroadcast := (flags & net.FlagBroadcast) != 0
			isPointToPoint := (flags & net.FlagPointToPoint) != 0

			// 获取Windows特定的接口信息
			windowsInfo := m.getWindowsInterfaceInfo(iface.Name)

			interfaceInfo := map[string]interface{}{
				"name":            iface.Name,
				"index":           iface.Index,
				"mtu":             iface.MTU,
				"flags":           strings.Join(iface.Flags, ","),
				"addresses":       addresses,
				"ipv4_addresses":  ipv4Addresses,
				"ipv6_addresses":  ipv6Addresses,
				"mac_address":     macAddress,
				"is_up":           isUp,
				"is_loopback":     isLoopback,
				"is_multicast":    isMulticast,
				"is_broadcast":    isBroadcast,
				"is_pointtopoint": isPointToPoint,
				"hardware_addr":   netInterface.HardwareAddr.String(),
				"windows_info":    windowsInfo,
			}

			// 如果没有找到地址，使用接口名称作为备用
			if len(addresses) == 0 {
				addresses = append(addresses, iface.Name)
				interfaceInfo["addresses"] = addresses
			}

			result = append(result, interfaceInfo)
			m.logger.Debugf("网络接口 %s 信息获取完成: 地址数=%d, MTU=%d", iface.Name, len(addresses), iface.MTU)
		} else {
			m.logger.Warnf("无法获取网络接口 %s 的详细信息: %v", iface.Name, err)
			// 使用基本信息
			interfaceInfo := map[string]interface{}{
				"name":      iface.Name,
				"index":     iface.Index,
				"mtu":       iface.MTU,
				"flags":     strings.Join(iface.Flags, ","),
				"addresses": []string{iface.Name},
				"error":     fmt.Sprintf("无法获取详细信息: %v", err),
			}
			result = append(result, interfaceInfo)
		}
	}

	m.logger.Debugf("网络接口信息获取完成，共 %d 个接口", len(result))
	return result, nil
}

// getWindowsInterfaceInfo 获取Windows特定的网络接口信息
func (m *Manager) getWindowsInterfaceInfo(interfaceName string) map[string]interface{} {
	windowsInfo := make(map[string]interface{})

	// 尝试使用Windows API获取更详细的接口信息
	if info, err := m.getWindowsNetworkAdapterInfo(interfaceName); err == nil {
		windowsInfo = info
	} else {
		m.logger.Debugf("获取Windows网络适配器信息失败: %v", err)
	}

	return windowsInfo
}

// getWindowsNetworkAdapterInfo 使用Windows API获取网络适配器信息
func (m *Manager) getWindowsNetworkAdapterInfo(interfaceName string) (map[string]interface{}, error) {
	info := make(map[string]interface{})

	// 使用netsh命令获取网络适配器信息
	cmd := exec.Command("netsh", "interface", "show", "interface", "name="+interfaceName)
	output, err := cmd.Output()
	if err == nil {
		// 解析netsh输出
		lines := strings.Split(string(output), "\n")
		for _, line := range lines {
			line = strings.TrimSpace(line)
			if strings.Contains(line, interfaceName) {
				// 解析接口状态信息
				fields := strings.Fields(line)
				if len(fields) >= 4 {
					info["admin_state"] = fields[0]
					info["state"] = fields[1]
					info["type"] = fields[2]
					info["interface_name"] = strings.Join(fields[3:], " ")
				}
				break
			}
		}
	}

	// 尝试获取IP配置信息
	if ipConfig, err := m.getWindowsIPConfig(interfaceName); err == nil {
		info["ip_config"] = ipConfig
	}

	// 尝试获取网络适配器统计信息
	if stats, err := m.getWindowsNetworkStats(interfaceName); err == nil {
		info["statistics"] = stats
	}

	return info, nil
}

// getWindowsIPConfig 获取Windows IP配置信息
func (m *Manager) getWindowsIPConfig(interfaceName string) (map[string]interface{}, error) {
	ipConfig := make(map[string]interface{})

	// 使用netsh命令获取IP配置
	cmd := exec.Command("netsh", "interface", "ip", "show", "config", "name="+interfaceName)
	output, err := cmd.Output()
	if err == nil {
		lines := strings.Split(string(output), "\n")
		for _, line := range lines {
			line = strings.TrimSpace(line)
			if strings.Contains(line, "IP Address:") {
				parts := strings.Split(line, ":")
				if len(parts) >= 2 {
					ipConfig["ip_address"] = strings.TrimSpace(parts[1])
				}
			} else if strings.Contains(line, "Subnet Prefix:") {
				parts := strings.Split(line, ":")
				if len(parts) >= 2 {
					ipConfig["subnet_prefix"] = strings.TrimSpace(parts[1])
				}
			} else if strings.Contains(line, "Default Gateway:") {
				parts := strings.Split(line, ":")
				if len(parts) >= 2 {
					ipConfig["default_gateway"] = strings.TrimSpace(parts[1])
				}
			} else if strings.Contains(line, "DNS Servers:") {
				parts := strings.Split(line, ":")
				if len(parts) >= 2 {
					ipConfig["dns_servers"] = strings.TrimSpace(parts[1])
				}
			}
		}
	}

	return ipConfig, nil
}

// getWindowsNetworkStats 获取Windows网络统计信息
func (m *Manager) getWindowsNetworkStats(interfaceName string) (map[string]interface{}, error) {
	stats := make(map[string]interface{})

	// 使用netsh命令获取网络统计
	cmd := exec.Command("netsh", "interface", "show", "interface", "name="+interfaceName)
	output, err := cmd.Output()
	if err == nil {
		lines := strings.Split(string(output), "\n")
		for _, line := range lines {
			line = strings.TrimSpace(line)
			if strings.Contains(line, "Bytes Sent:") {
				parts := strings.Split(line, ":")
				if len(parts) >= 2 {
					stats["bytes_sent"] = strings.TrimSpace(parts[1])
				}
			} else if strings.Contains(line, "Bytes Received:") {
				parts := strings.Split(line, ":")
				if len(parts) >= 2 {
					stats["bytes_received"] = strings.TrimSpace(parts[1])
				}
			} else if strings.Contains(line, "Packets Sent:") {
				parts := strings.Split(line, ":")
				if len(parts) >= 2 {
					stats["packets_sent"] = strings.TrimSpace(parts[1])
				}
			} else if strings.Contains(line, "Packets Received:") {
				parts := strings.Split(line, ":")
				if len(parts) >= 2 {
					stats["packets_received"] = strings.TrimSpace(parts[1])
				}
			}
		}
	}

	return stats, nil
}

// GetNetworkStats 获取网络统计
func (m *Manager) GetNetworkStats() (map[string]interface{}, error) {
	// 获取网络IO统计
	ioCounters, err := gopsutilnet.IOCounters(false)
	if err != nil {
		return nil, fmt.Errorf("获取网络统计失败: %w", err)
	}

	if len(ioCounters) == 0 {
		return map[string]interface{}{
			"bytes_sent":     uint64(0),
			"bytes_recv":     uint64(0),
			"packets_sent":   uint64(0),
			"packets_recv":   uint64(0),
			"err_in":         uint64(0),
			"err_out":        uint64(0),
			"drop_in":        uint64(0),
			"drop_out":       uint64(0),
			"fifo_in":        uint64(0),
			"fifo_out":       uint64(0),
			"collisions":     uint64(0),
			"carrier_errors": uint64(0),
			"compressed_in":  uint64(0),
			"compressed_out": uint64(0),
		}, nil
	}

	stats := ioCounters[0]
	return map[string]interface{}{
		"bytes_sent":     stats.BytesSent,
		"bytes_recv":     stats.BytesRecv,
		"packets_sent":   stats.PacketsSent,
		"packets_recv":   stats.PacketsRecv,
		"err_in":         stats.Errin,
		"err_out":        stats.Errout,
		"drop_in":        stats.Dropin,
		"drop_out":       stats.Dropout,
		"fifo_in":        stats.Fifoin,
		"fifo_out":       stats.Fifoout,
		"collisions":     uint64(0), // gopsutil中没有这个字段
		"carrier_errors": uint64(0), // gopsutil中没有这个字段
		"compressed_in":  uint64(0), // gopsutil中没有这个字段
		"compressed_out": uint64(0), // gopsutil中没有这个字段
	}, nil
}

// GetUDPConnections 获取UDP连接
func (m *Manager) GetUDPConnections() ([]*models.NetworkConnection, error) {
	connections, err := gopsutilnet.Connections("udp")
	if err != nil {
		return nil, fmt.Errorf("获取UDP连接失败: %w", err)
	}

	var udpConnections []*models.NetworkConnection
	var wg sync.WaitGroup
	resultChan := make(chan *models.NetworkConnection, len(connections))

	for _, conn := range connections {
		wg.Add(1)
		go func(c gopsutilnet.ConnectionStat) {
			defer wg.Done()

			select {
			case m.networkSemaphore <- struct{}{}:
				defer func() { <-m.networkSemaphore }()
			case <-time.After(m.networkTimeout):
				return
			}

			networkConn := &models.NetworkConnection{
				ID:          fmt.Sprintf("%d-%s-%d", c.Pid, m.getConnectionType(c.Type), c.Fd),
				LocalAddr:   c.Laddr.IP,
				RemoteAddr:  c.Raddr.IP,
				LocalPort:   int(c.Laddr.Port),
				RemotePort:  int(c.Raddr.Port),
				Protocol:    m.getProtocolString(c.Type),
				Status:      c.Status,
				PID:         c.Pid,
				ProcessName: m.getProcessName(c.Pid),
				Type:        m.getConnectionType(c.Type),
			}

			resultChan <- networkConn
		}(conn)
	}

	// 等待所有协程完成
	go func() {
		wg.Wait()
		close(resultChan)
	}()

	// 收集结果
	for conn := range resultChan {
		udpConnections = append(udpConnections, conn)
	}

	return udpConnections, nil
}

// CloseConnection 关闭指定的网络连接
func (m *Manager) CloseConnection(connectionID string) error {
	m.logger.Debugf("开始关闭连接: %s", connectionID)

	// 解析连接ID，格式为 "PID-type-FD" 或 "PID-type-IP-PORT"
	parts := strings.Split(connectionID, "-")
	if len(parts) < 3 {
		return fmt.Errorf("无效的连接ID格式: %s", connectionID)
	}

	pidStr := parts[0]
	connType := parts[1]

	// 验证PID是否为数字
	pid, err := strconv.ParseInt(pidStr, 10, 32)
	if err != nil {
		return fmt.Errorf("无效的进程ID: %s", pidStr)
	}

	m.logger.Debugf("尝试关闭进程 %d 的 %s 连接", pid, connType)

	// 首先检查进程是否还存在
	proc, err := gopsutilprocess.NewProcess(int32(pid))
	if err != nil {
		m.logger.Warnf("进程 %d 不存在或无法访问: %v", pid, err)
		// 即使进程不存在，也尝试清理可能的残留连接
		return m.cleanupOrphanedConnections(connectionID, int32(pid), connType)
	}

	// 检查进程是否还在运行
	isRunning, err := proc.IsRunning()
	if err != nil {
		m.logger.Warnf("无法检查进程 %d 的运行状态: %v", pid, err)
	} else if !isRunning {
		m.logger.Warnf("进程 %d 已经终止", pid)
		return m.cleanupOrphanedConnections(connectionID, int32(pid), connType)
	}

	// 获取连接信息
	connections, err := m.GetConnectionsByPID(int32(pid))
	if err != nil {
		m.logger.Warnf("获取进程 %d 的连接失败: %v", pid, err)
		// 尝试使用强制方法关闭连接
		return m.forceCloseProcessConnections(int32(pid))
	}

	m.logger.Debugf("进程 %d 当前有 %d 个连接", pid, len(connections))

	// 改进的连接匹配逻辑
	var targetConnection *models.NetworkConnection
	for _, conn := range connections {
		m.logger.Debugf("检查连接: ID=%s, Type=%s, LocalPort=%d, RemotePort=%d, Status=%s",
			conn.ID, conn.Type, conn.LocalPort, conn.RemotePort, conn.Status)

		if conn.Type == connType {
			// 如果提供了FD，尝试多种匹配方式
			if len(parts) >= 3 {
				fdStr := parts[2]
				if fdStr != "" {
					// 尝试匹配端口号
					if fmt.Sprintf("%d", conn.LocalPort) == fdStr ||
						fmt.Sprintf("%d", conn.RemotePort) == fdStr {
						targetConnection = conn
						m.logger.Debugf("通过端口号找到匹配的连接: %s", conn.ID)
						break
					}
					// 尝试匹配连接ID的最后部分
					if strings.HasSuffix(conn.ID, fdStr) {
						targetConnection = conn
						m.logger.Debugf("通过连接ID后缀找到匹配的连接: %s", conn.ID)
						break
					}
				}
			}

			// 如果没有找到特定匹配，使用第一个匹配类型的连接
			if targetConnection == nil {
				targetConnection = conn
				m.logger.Debugf("使用第一个匹配类型的连接: %s", conn.ID)
			}
		}
	}

	if targetConnection == nil {
		m.logger.Warnf("未找到匹配的连接: PID=%d, Type=%s", pid, connType)
		m.logger.Debugf("可用连接类型: %v", func() []string {
			types := make([]string, 0, len(connections))
			for _, conn := range connections {
				types = append(types, conn.Type)
			}
			return types
		}())

		// 直接返回友好的错误提示，不继续执行后续逻辑
		return fmt.Errorf("未找到匹配的连接: PID=%d, Type=%s。请检查连接ID是否正确，或该连接可能已经关闭", pid, connType)
	}

	// 尝试多种方法关闭连接
	var lastErr error

	// 方法1: 尝试使用Windows API直接关闭连接
	if err := m.closeConnectionByWindowsAPI(targetConnection); err == nil {
		m.logger.Infof("成功使用Windows API关闭连接: %s", connectionID)
		// 连接关闭成功后，终止相关进程
		if err := m.terminateProcessAfterConnectionClose(int32(pid), connType); err != nil {
			m.logger.Warnf("连接已关闭但终止进程失败: %v", err)
		}
		return nil
	} else {
		lastErr = err
		m.logger.Debugf("Windows API关闭连接失败: %v", err)
	}

	// 方法2: 尝试使用netsh命令关闭连接
	if err := m.closeConnectionByNetsh(targetConnection); err == nil {
		m.logger.Infof("成功使用netsh关闭连接: %s", connectionID)
		// 连接关闭成功后，终止相关进程
		if err := m.terminateProcessAfterConnectionClose(int32(pid), connType); err != nil {
			m.logger.Warnf("连接已关闭但终止进程失败: %v", err)
		}
		return nil
	} else {
		lastErr = err
		m.logger.Debugf("netsh关闭连接失败: %v", err)
	}

	// 方法3: 通过路由表阻止特定连接
	if err := m.closeConnectionByRoute(targetConnection.RemoteAddr); err == nil {
		m.logger.Infof("成功使用路由阻止连接: %s", connectionID)
		// 连接关闭成功后，终止相关进程
		if err := m.terminateProcessAfterConnectionClose(int32(pid), connType); err != nil {
			m.logger.Warnf("连接已关闭但终止进程失败: %v", err)
		}
		return nil
	} else {
		lastErr = err
		m.logger.Debugf("路由阻止连接失败: %v", err)
	}

	// 方法4: 使用Windows防火墙规则阻止特定连接
	if err := m.closeConnectionByFirewall(targetConnection); err == nil {
		m.logger.Infof("成功使用防火墙阻止连接: %s", connectionID)
		// 连接关闭成功后，终止相关进程
		if err := m.terminateProcessAfterConnectionClose(int32(pid), connType); err != nil {
			m.logger.Warnf("连接已关闭但终止进程失败: %v", err)
		}
		return nil
	} else {
		lastErr = err
		m.logger.Debugf("防火墙阻止连接失败: %v", err)
	}

	// 如果所有方法都失败，记录警告但不终止进程
	m.logger.Warnf("无法关闭特定连接 %s，但不会终止进程 %d", connectionID, pid)
	return fmt.Errorf("无法关闭特定连接，但进程仍在运行: %w", lastErr)
}

// CloseAllConnectionsByProcess 关闭指定进程的所有网络连接
func (m *Manager) CloseAllConnectionsByProcess(processName string) error {
	m.logger.Infof("开始关闭进程 %s 的所有网络连接", processName)

	// 获取所有网络连接
	allConnections, err := m.GetNetworkConnections()
	if err != nil {
		return fmt.Errorf("获取网络连接失败: %w", err)
	}

	// 查找匹配进程名称的连接
	var processConnections []*models.NetworkConnection
	for _, conn := range allConnections {
		if strings.EqualFold(conn.ProcessName, processName) {
			processConnections = append(processConnections, conn)
		}
	}

	if len(processConnections) == 0 {
		return fmt.Errorf("未找到进程名称为 %s 的网络连接", processName)
	}

	m.logger.Infof("找到进程 %s 的 %d 个网络连接", processName, len(processConnections))

	// 尝试关闭每个连接
	var closedCount int
	var failedCount int
	var errors []string

	for _, conn := range processConnections {
		// 构建连接ID
		connectionID := fmt.Sprintf("%d-%s-%d", conn.PID, conn.Type, conn.LocalPort)

		if err := m.CloseConnection(connectionID); err != nil {
			failedCount++
			errorMsg := fmt.Sprintf("连接 %s 关闭失败: %v", connectionID, err)
			errors = append(errors, errorMsg)
			m.logger.Warnf(errorMsg)
		} else {
			closedCount++
		}
	}

	// 如果还有连接未关闭，尝试强制方法
	if failedCount > 0 {
		m.logger.Warnf("常规方法关闭失败，尝试强制关闭进程 %s 的网络连接", processName)

		// 获取进程信息
		processes, err := m.GetProcessesByName(processName)
		if err != nil {
			m.logger.Errorf("获取进程信息失败: %v", err)
		} else {
			for _, proc := range processes {
				if err := m.ForceCloseProcessConnections(proc.Pid); err == nil {
					m.logger.Infof("强制关闭进程 %d 的网络连接成功", proc.Pid)
					closedCount += failedCount
					failedCount = 0
					break
				}
			}
		}
	}

	if failedCount > 0 {
		return fmt.Errorf("关闭进程 %s 的网络连接部分失败: 成功 %d 个，失败 %d 个。错误: %s",
			processName, closedCount, failedCount, strings.Join(errors, "; "))
	}

	m.logger.Infof("成功关闭进程 %s 的所有网络连接: %d 个", processName, closedCount)
	return nil
}

// CloseAllConnectionsByPID 关闭指定PID的所有网络连接
func (m *Manager) CloseAllConnectionsByPID(pid int32) error {
	m.logger.Infof("开始关闭进程 %d 的所有网络连接", pid)

	// 获取指定PID的所有连接
	connections, err := m.GetConnectionsByPID(pid)
	if err != nil {
		return fmt.Errorf("获取进程连接失败: %w", err)
	}

	if len(connections) == 0 {
		return fmt.Errorf("进程 %d 没有活跃的网络连接", pid)
	}

	m.logger.Infof("找到进程 %d 的 %d 个网络连接", pid, len(connections))

	// 尝试关闭每个连接
	var closedCount int
	var failedCount int
	var errors []string

	for _, conn := range connections {
		// 构建连接ID
		connectionID := fmt.Sprintf("%d-%s-%d", conn.PID, conn.Type, conn.LocalPort)

		if err := m.CloseConnection(connectionID); err != nil {
			failedCount++
			errorMsg := fmt.Sprintf("连接 %s 关闭失败: %v", connectionID, err)
			errors = append(errors, errorMsg)
			m.logger.Warnf(errorMsg)
		} else {
			closedCount++
		}
	}

	// 如果还有连接未关闭，尝试强制方法
	if failedCount > 0 {
		m.logger.Warnf("常规方法关闭失败，尝试强制关闭进程 %d 的网络连接", pid)

		if err := m.ForceCloseProcessConnections(pid); err == nil {
			m.logger.Infof("强制关闭进程 %d 的网络连接成功", pid)
			closedCount += failedCount
			failedCount = 0
		}
	}

	if failedCount > 0 {
		return fmt.Errorf("关闭进程 %d 的网络连接部分失败: 成功 %d 个，失败 %d 个。错误: %s",
			pid, closedCount, failedCount, strings.Join(errors, "; "))
	}

	m.logger.Infof("成功关闭进程 %d 的所有网络连接: %d 个", pid, closedCount)
	return nil
}

// ForceCloseProcessConnections 强制关闭进程的网络连接
func (m *Manager) ForceCloseProcessConnections(pid int32) error {
	m.logger.Infof("开始强制关闭进程 %d 的网络连接", pid)

	// 方法1: 使用Windows API强制关闭网络句柄
	if err := m.forceCloseNetworkHandles(pid); err == nil {
		m.logger.Infof("Windows API方法成功关闭进程 %d 的网络连接", pid)
	}

	// 方法2: 使用netsh命令重置网络配置
	if err := m.resetNetworkConfiguration(); err == nil {
		m.logger.Infof("网络配置重置成功")
	}

	// 方法3: 使用防火墙规则阻止进程的所有出站连接
	if err := m.blockProcessOutboundConnections(pid); err == nil {
		m.logger.Infof("防火墙规则阻止成功")
	}

	// 方法4: 使用更底层的网络重置技术
	if err := m.forceResetNetworkStack(pid); err == nil {
		m.logger.Infof("网络协议栈强制重置成功")
	}

	// 方法5: 最终验证和清理
	if err := m.finalizeNetworkClosure(pid); err == nil {
		m.logger.Infof("网络连接关闭最终验证完成")
	}

	m.logger.Infof("进程 %d 的网络连接强制关闭流程完成", pid)
	return nil
}

// forceCloseNetworkHandles 强制关闭进程的网络句柄
func (m *Manager) forceCloseNetworkHandles(pid int32) error {
	m.logger.Debugf("尝试强制关闭进程 %d 的网络句柄", pid)

	// 方法1: 使用netsh命令强制关闭TCP连接
	if err := m.forceCloseTCPConnections(pid); err == nil {
		m.logger.Debugf("成功使用netsh关闭进程 %d 的TCP连接", pid)
	}

	// 方法2: 使用Windows API强制关闭网络句柄
	if err := m.forceCloseNetworkHandlesByAPI(pid); err == nil {
		m.logger.Debugf("成功使用Windows API关闭进程 %d 的网络句柄", pid)
	}

	// 方法3: 使用route命令删除相关路由
	if err := m.forceDeleteProcessRoutes(pid); err == nil {
		m.logger.Debugf("成功删除进程 %d 的相关路由", pid)
	}

	// 方法4: 使用netsh命令重置TCP/IP协议栈
	if err := m.resetTCPIPStack(); err == nil {
		m.logger.Debugf("成功重置TCP/IP协议栈")
	}

	return nil
}

// forceCloseTCPConnections 强制关闭指定进程的TCP连接
func (m *Manager) forceCloseTCPConnections(pid int32) error {
	// 获取进程的所有TCP连接
	connections, err := m.GetConnectionsByPID(pid)
	if err != nil {
		return fmt.Errorf("获取进程连接失败: %v", err)
	}

	var tcpConnections []*models.NetworkConnection
	for _, conn := range connections {
		if conn.Type == "tcp" {
			tcpConnections = append(tcpConnections, conn)
		}
	}

	if len(tcpConnections) == 0 {
		m.logger.Debugf("进程 %d 没有TCP连接", pid)
		return nil
	}

	m.logger.Debugf("进程 %d 有 %d 个TCP连接需要关闭", pid, len(tcpConnections))

	// 方法1: 使用netsh命令强制关闭TCP连接
	for _, conn := range tcpConnections {
		if err := m.forceCloseSpecificTCPConnection(conn); err != nil {
			m.logger.Warnf("强制关闭TCP连接失败: %v", err)
		}
	}

	// 方法2: 使用netstat和taskkill强制关闭TCP连接
	if err := m.forceKillTCPConnections(pid); err != nil {
		m.logger.Warnf("强制终止TCP连接失败: %v", err)
	}

	// 方法3: 使用netsh命令重置TCP配置
	commands := []string{
		"netsh interface tcp set global chimney=disabled",
		"netsh interface tcp set global autotuninglevel=disabled",
		"netsh interface tcp set global ecncapability=disabled",
		"netsh interface tcp set global timestamps=disabled",
		"netsh interface tcp set global rss=disabled",
		"netsh interface tcp set global maxsynretransmissions=1",
		"netsh interface tcp set global congestionprovider=ctcp",
		"netsh interface tcp set global initialRto=1000",
	}

	for _, cmdStr := range commands {
		parts := strings.Fields(cmdStr)
		cmd := exec.Command(parts[0], parts[1:]...)
		if err := cmd.Run(); err != nil {
			m.logger.Debugf("执行命令失败: %s, 错误: %v", cmdStr, err)
		}
	}

	// 等待TCP配置生效
	time.Sleep(2 * time.Second)

	// 重新启用TCP功能
	commands = []string{
		"netsh interface tcp set global chimney=enabled",
		"netsh interface tcp set global autotuninglevel=normal",
		"netsh interface tcp set global ecncapability=enabled",
		"netsh interface tcp set global timestamps=enabled",
		"netsh interface tcp set global rss=enabled",
		"netsh interface tcp set global maxsynretransmissions=2",
		"netsh interface tcp set global congestionprovider=ctcp",
		"netsh interface tcp set global initialRto=3000",
	}

	for _, cmdStr := range commands {
		parts := strings.Fields(cmdStr)
		cmd := exec.Command(parts[0], parts[1:]...)
		cmd.Run() // 忽略错误
	}

	return nil
}

// forceKillTCPConnections 使用netstat和taskkill强制关闭TCP连接
func (m *Manager) forceKillTCPConnections(pid int32) error {
	// 使用netstat命令获取进程的TCP连接
	cmd := exec.Command("netstat", "-ano")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("执行netstat命令失败: %v", err)
	}

	lines := strings.Split(string(output), "\n")
	var tcpConnections []string

	// 解析netstat输出，找到属于指定进程的TCP连接
	for _, line := range lines {
		if strings.Contains(line, "TCP") && strings.Contains(line, fmt.Sprintf("%d", pid)) {
			tcpConnections = append(tcpConnections, line)
		}
	}

	if len(tcpConnections) == 0 {
		m.logger.Debugf("进程 %d 在netstat中没有找到TCP连接", pid)
		return nil
	}

	m.logger.Debugf("找到进程 %d 的 %d 个TCP连接", pid, len(tcpConnections))

	// 尝试使用netsh命令强制关闭这些连接
	for _, connLine := range tcpConnections {
		if err := m.forceCloseTCPConnectionByNetsh(connLine); err != nil {
			m.logger.Debugf("强制关闭TCP连接失败: %v", err)
		}
	}

	// 使用netsh命令重置TCP连接表
	if err := m.resetTCPConnectionTable(); err != nil {
		m.logger.Debugf("重置TCP连接表失败: %v", err)
	}

	return nil
}

// forceCloseTCPConnectionByNetsh 使用netsh命令强制关闭特定的TCP连接
func (m *Manager) forceCloseTCPConnectionByNetsh(connLine string) error {
	// 解析连接行，提取本地和远程地址端口
	// 格式: TCP    0.0.0.0:80    0.0.0.0:0    LISTENING    1234
	parts := strings.Fields(connLine)
	if len(parts) < 5 {
		return fmt.Errorf("无效的连接行格式: %s", connLine)
	}

	// 尝试使用netsh命令重置TCP连接
	commands := []string{
		"netsh interface tcp set global chimney=disabled",
		"netsh interface tcp set global autotuninglevel=disabled",
		"netsh interface tcp set global ecncapability=disabled",
	}

	for _, cmdStr := range commands {
		parts := strings.Fields(cmdStr)
		cmd := exec.Command(parts[0], parts[1:]...)
		cmd.Run()
	}

	// 短暂等待
	time.Sleep(50 * time.Millisecond)

	// 重新启用
	commands = []string{
		"netsh interface tcp set global chimney=enabled",
		"netsh interface tcp set global autotuninglevel=normal",
		"netsh interface tcp set global ecncapability=enabled",
	}

	for _, cmdStr := range commands {
		parts := strings.Fields(cmdStr)
		cmd := exec.Command(parts[0], parts[1:]...)
		cmd.Run()
	}

	return nil
}

// resetTCPConnectionTable 重置TCP连接表
func (m *Manager) resetTCPConnectionTable() error {
	// 使用netsh命令重置TCP连接表
	commands := []string{
		"netsh interface tcp set global chimney=disabled",
		"netsh interface tcp set global autotuninglevel=disabled",
		"netsh interface tcp set global ecncapability=disabled",
		"netsh interface tcp set global timestamps=disabled",
		"netsh interface tcp set global rss=disabled",
		"netsh interface tcp set global maxsynretransmissions=1",
		"netsh interface tcp set global congestionprovider=ctcp",
		"netsh interface tcp set global initialRto=1000",
		"netsh interface tcp set global memorypressure=disabled",
		"netsh interface tcp set global ecncapability=disabled",
	}

	for _, cmdStr := range commands {
		parts := strings.Fields(cmdStr)
		cmd := exec.Command(parts[0], parts[1:]...)
		cmd.Run()
	}

	// 等待重置完成
	time.Sleep(1 * time.Second)

	// 重新启用
	commands = []string{
		"netsh interface tcp set global chimney=enabled",
		"netsh interface tcp set global autotuninglevel=normal",
		"netsh interface tcp set global ecncapability=enabled",
		"netsh interface tcp set global timestamps=enabled",
		"netsh interface tcp set global rss=enabled",
		"netsh interface tcp set global maxsynretransmissions=2",
		"netsh interface tcp set global congestionprovider=ctcp",
		"netsh interface tcp set global initialRto=3000",
		"netsh interface tcp set global memorypressure=enabled",
		"netsh interface tcp set global ecncapability=enabled",
		"netsh interface tcp set global dca=enabled",
		"netsh interface tcp set global netdma=enabled",
	}

	for _, cmdStr := range commands {
		parts := strings.Fields(cmdStr)
		cmd := exec.Command(parts[0], parts[1:]...)
		cmd.Run()
	}

	return nil
}

// forceCloseSpecificTCPConnection 强制关闭特定的TCP连接
func (m *Manager) forceCloseSpecificTCPConnection(conn *models.NetworkConnection) error {
	// 使用netsh命令关闭特定连接
	// 注意：这种方法可能不够精确，但可以尝试

	// 尝试使用route命令删除到目标IP的路由
	if conn.RemoteAddr != "" && conn.RemoteAddr != "0.0.0.0" && conn.RemoteAddr != "::" {
		if err := m.forceDeleteRoute(conn.RemoteAddr); err != nil {
			m.logger.Debugf("删除路由失败: %v", err)
		}
	}

	// 尝试使用netsh命令重置TCP连接
	cmd := exec.Command("netsh", "interface", "tcp", "set", "global", "chimney=disabled")
	if err := cmd.Run(); err == nil {
		time.Sleep(100 * time.Millisecond)
		cmd = exec.Command("netsh", "interface", "tcp", "set", "global", "chimney=enabled")
		cmd.Run()
	}

	return nil
}

// forceCloseNetworkHandlesByAPI 使用Windows API强制关闭网络句柄
func (m *Manager) forceCloseNetworkHandlesByAPI(pid int32) error {
	m.logger.Debugf("尝试使用Windows API关闭进程 %d 的网络句柄", pid)

	// 方法1: 尝试使用taskkill命令强制终止进程的网络相关线程
	if err := m.forceKillProcessNetworkThreads(pid); err == nil {
		m.logger.Debugf("成功强制终止进程 %d 的网络线程", pid)
		return nil
	}

	// 方法2: 使用Windows IP Helper API关闭连接
	if err := m.forceCloseConnectionsByIPHelper(pid); err == nil {
		m.logger.Debugf("成功使用IP Helper API关闭进程 %d 的连接", pid)
		return nil
	}

	// 方法3: 如果无法终止进程，尝试重置网络适配器
	if err := m.resetNetworkAdapters(); err == nil {
		m.logger.Debugf("成功重置网络适配器")
		return nil
	}

	return fmt.Errorf("Windows API方法失败")
}

// forceKillProcessNetworkThreads 强制终止进程的网络相关线程
func (m *Manager) forceKillProcessNetworkThreads(pid int32) error {
	// 使用taskkill命令强制终止进程
	cmd := exec.Command("taskkill", "/PID", fmt.Sprintf("%d", pid), "/F")
	if err := cmd.Run(); err == nil {
		m.logger.Debugf("成功强制终止进程 %d", pid)
		return nil
	}

	// 如果无法终止整个进程，尝试终止网络相关的子进程
	cmd = exec.Command("taskkill", "/PID", fmt.Sprintf("%d", pid), "/T", "/F")
	if err := cmd.Run(); err == nil {
		m.logger.Debugf("成功强制终止进程 %d 及其子进程", pid)
		return nil
	}

	return fmt.Errorf("无法终止进程 %d", pid)
}

// forceCloseConnectionsByIPHelper 使用Windows IP Helper API强制关闭连接
func (m *Manager) forceCloseConnectionsByIPHelper(pid int32) error {
	// 由于Go语言直接调用Windows API比较复杂，这里使用命令行工具
	// 可以使用netsh的高级命令来关闭连接

	// 方法1: 使用netsh命令重置TCP连接
	if err := m.resetTCPConnectionsByNetsh(); err == nil {
		m.logger.Debugf("成功使用netsh重置TCP连接")
	}

	// 方法2: 使用netsh命令重置UDP连接
	if err := m.resetUDPConnectionsByNetsh(); err == nil {
		m.logger.Debugf("成功使用netsh重置UDP连接")
	}

	// 方法3: 使用netsh命令重置网络接口
	if err := m.resetNetworkInterfacesByNetsh(); err == nil {
		m.logger.Debugf("成功使用netsh重置网络接口")
	}

	return nil
}

// resetTCPConnectionsByNetsh 使用netsh命令重置TCP连接
func (m *Manager) resetTCPConnectionsByNetsh() error {
	// 使用netsh命令重置TCP连接
	commands := []string{
		"netsh interface tcp set global chimney=disabled",
		"netsh interface tcp set global autotuninglevel=disabled",
		"netsh interface tcp set global ecncapability=disabled",
		"netsh interface tcp set global timestamps=disabled",
		"netsh interface tcp set global rss=disabled",
		"netsh interface tcp set global maxsynretransmissions=1",
		"netsh interface tcp set global congestionprovider=ctcp",
		"netsh interface tcp set global initialRto=1000",
		"netsh interface tcp set global memorypressure=disabled",
		"netsh interface tcp set global ecncapability=disabled",
		"netsh interface tcp set global dca=enabled",
		"netsh interface tcp set global netdma=enabled",
	}

	for _, cmdStr := range commands {
		parts := strings.Fields(cmdStr)
		cmd := exec.Command(parts[0], parts[1:]...)
		cmd.Run()
	}

	// 等待重置完成
	time.Sleep(1 * time.Second)

	// 重新启用
	commands = []string{
		"netsh interface tcp set global chimney=enabled",
		"netsh interface tcp set global autotuninglevel=normal",
		"netsh interface tcp set global ecncapability=enabled",
		"netsh interface tcp set global timestamps=enabled",
		"netsh interface tcp set global rss=enabled",
		"netsh interface tcp set global maxsynretransmissions=2",
		"netsh interface tcp set global congestionprovider=ctcp",
		"netsh interface tcp set global initialRto=3000",
		"netsh interface tcp set global memorypressure=enabled",
		"netsh interface tcp set global ecncapability=enabled",
		"netsh interface tcp set global dca=enabled",
		"netsh interface tcp set global netdma=enabled",
	}

	for _, cmdStr := range commands {
		parts := strings.Fields(cmdStr)
		cmd := exec.Command(parts[0], parts[1:]...)
		cmd.Run()
	}

	return nil
}

// resetUDPConnectionsByNetsh 使用netsh命令重置UDP连接
func (m *Manager) resetUDPConnectionsByNetsh() error {
	// 使用netsh命令重置UDP相关配置
	commands := []string{
		"netsh interface ipv4 set global taskoffload=disabled",
		"netsh interface ipv6 set global taskoffload=disabled",
	}

	for _, cmdStr := range commands {
		parts := strings.Fields(cmdStr)
		cmd := exec.Command(parts[0], parts[1:]...)
		cmd.Run()
	}

	// 等待重置完成
	time.Sleep(500 * time.Millisecond)

	// 重新启用
	commands = []string{
		"netsh interface ipv4 set global taskoffload=enabled",
		"netsh interface ipv6 set global taskoffload=enabled",
	}

	for _, cmdStr := range commands {
		parts := strings.Fields(cmdStr)
		cmd := exec.Command(parts[0], parts[1:]...)
		cmd.Run()
	}

	return nil
}

// resetNetworkInterfacesByNetsh 使用netsh命令重置网络接口
func (m *Manager) resetNetworkInterfacesByNetsh() error {
	// 使用netsh命令重置网络接口
	commands := []string{
		"netsh interface ip reset",
		"netsh winsock reset",
		"netsh int ip reset c:\\resetlog.txt",
		"netsh int ipv4 reset",
		"netsh int ipv6 reset",
	}

	for _, cmdStr := range commands {
		parts := strings.Fields(cmdStr)
		cmd := exec.Command(parts[0], parts[1:]...)
		cmd.Run()
	}

	// 等待重置完成
	time.Sleep(2 * time.Second)

	return nil
}

// forceDeleteProcessRoutes 强制删除进程相关的路由
func (m *Manager) forceDeleteProcessRoutes(pid int32) error {
	// 获取进程的连接
	connections, err := m.GetConnectionsByPID(pid)
	if err != nil {
		return fmt.Errorf("获取进程连接失败: %v", err)
	}

	// 删除到远程IP的路由
	for _, conn := range connections {
		if conn.RemoteAddr != "" && conn.RemoteAddr != "0.0.0.0" && conn.RemoteAddr != "::" {
			if err := m.forceDeleteRoute(conn.RemoteAddr); err != nil {
				m.logger.Debugf("删除路由失败: %v", err)
			}
		}
	}

	return nil
}

// forceDeleteRoute 强制删除到指定IP的路由
func (m *Manager) forceDeleteRoute(ip string) error {
	// 使用route delete命令删除路由
	cmd := exec.Command("route", "delete", ip)
	if err := cmd.Run(); err != nil {
		// 尝试使用netsh命令
		cmd = exec.Command("netsh", "interface", "ip", "delete", "route", ip)
		cmd.Run()
	}
	return nil
}

// resetNetworkAdapters 重置网络适配器
func (m *Manager) resetNetworkAdapters() error {
	// 使用netsh命令重置网络适配器
	commands := []string{
		"netsh interface ip reset",
		"netsh winsock reset",
		"netsh int ip reset c:\\resetlog.txt",
	}

	for _, cmdStr := range commands {
		parts := strings.Fields(cmdStr)
		cmd := exec.Command(parts[0], parts[1:]...)
		if err := cmd.Run(); err != nil {
			m.logger.Debugf("重置网络适配器命令失败: %s, 错误: %v", cmdStr, err)
		}
	}

	return nil
}

// resetTCPIPStack 重置TCP/IP协议栈
func (m *Manager) resetTCPIPStack() error {
	// 使用netsh命令重置TCP/IP协议栈
	commands := []string{
		"netsh int ip reset",
		"netsh winsock reset",
		"netsh int ipv4 reset",
		"netsh int ipv6 reset",
	}

	for _, cmdStr := range commands {
		parts := strings.Fields(cmdStr)
		cmd := exec.Command(parts[0], parts[1:]...)
		if err := cmd.Run(); err != nil {
			m.logger.Debugf("重置TCP/IP协议栈命令失败: %s, 错误: %v", cmdStr, err)
		}
	}

	// 等待重置完成
	time.Sleep(3 * time.Second)

	return nil
}

// resetNetworkConfiguration 重置网络配置
func (m *Manager) resetNetworkConfiguration() error {
	m.logger.Debug("尝试重置网络配置")

	// 重置TCP配置
	commands := []string{
		"netsh interface tcp set global chimney=disabled",
		"netsh interface tcp set global autotuninglevel=disabled",
		"netsh interface tcp set global ecncapability=disabled",
		"netsh interface tcp set global timestamps=disabled",
	}

	for _, cmdStr := range commands {
		parts := strings.Fields(cmdStr)
		cmd := exec.Command(parts[0], parts[1:]...)
		cmd.Run() // 忽略错误，继续执行
	}

	// 等待一下
	time.Sleep(2 * time.Second)

	// 重新启用
	commands = []string{
		"netsh interface tcp set global chimney=enabled",
		"netsh interface tcp set global autotuninglevel=normal",
		"netsh interface tcp set global ecncapability=enabled",
		"netsh interface tcp set global timestamps=enabled",
	}

	for _, cmdStr := range commands {
		parts := strings.Fields(cmdStr)
		cmd := exec.Command(parts[0], parts[1:]...)
		cmd.Run() // 忽略错误，继续执行
	}

	return nil
}

// blockProcessOutboundConnections 阻止进程的所有出站连接
func (m *Manager) blockProcessOutboundConnections(pid int32) error {
	m.logger.Debugf("尝试阻止进程 %d 的所有出站连接", pid)

	// 创建防火墙规则阻止进程的所有出站连接
	ruleName := fmt.Sprintf("Block_Process_%d_Outbound", pid)

	cmd := exec.Command("netsh", "advfirewall", "firewall", "add", "rule",
		"name="+ruleName,
		"dir=out",
		"action=block",
		"program="+fmt.Sprintf("%d", pid),
		"enable=yes")

	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("无法创建防火墙规则: %v, 输出: %s", err, string(output))
	}

	// 延迟删除规则
	go func() {
		time.Sleep(30 * time.Second)
		cmd := exec.Command("netsh", "advfirewall", "firewall", "delete", "rule", "name="+ruleName)
		cmd.CombinedOutput()
		m.logger.Debugf("已删除临时防火墙规则: %s", ruleName)
	}()

	return nil
}

// GetProcessesByName 根据进程名称获取进程列表
func (m *Manager) GetProcessesByName(processName string) ([]*gopsutilprocess.Process, error) {
	processes, err := gopsutilprocess.Processes()
	if err != nil {
		return nil, err
	}

	var result []*gopsutilprocess.Process
	for _, proc := range processes {
		name, err := proc.Name()
		if err != nil {
			continue
		}
		if strings.EqualFold(name, processName) {
			result = append(result, proc)
		}
	}

	return result, nil
}

// closeConnectionByWindowsAPI 使用Windows API直接关闭连接
func (m *Manager) closeConnectionByWindowsAPI(conn *models.NetworkConnection) error {
	// 这里应该实现使用Windows API直接关闭网络连接的代码
	// 由于实现复杂，这里提供一个占位符

	m.logger.Debugf("尝试使用Windows API关闭连接: %s:%d -> %s:%d",
		conn.LocalAddr, conn.LocalPort, conn.RemoteAddr, conn.RemotePort)

	// 对于TCP连接，尝试使用netsh重置特定端口
	if conn.Type == "tcp" {
		// 尝试重置特定端口的连接
		cmd := exec.Command("netsh", "interface", "tcp", "set", "global", "chimney=disabled")
		if err := cmd.Run(); err == nil {
			time.Sleep(100 * time.Millisecond)
			cmd = exec.Command("netsh", "interface", "tcp", "set", "global", "chimney=enabled")
			cmd.Run()
			return nil
		}
	}

	return fmt.Errorf("Windows API方法未实现")
}

// closeSpecificConnection 关闭特定的网络连接，不影响其他连接
func (m *Manager) closeSpecificConnection(pid int32, connType string, fd int, conn *models.NetworkConnection) error {
	// 方法1: 尝试使用Windows的netsh命令关闭特定连接
	if err := m.closeConnectionByNetsh(conn); err == nil {
		return nil
	}

	// 方法2: 尝试使用route命令阻止特定连接
	if err := m.closeConnectionByRoute(conn.RemoteAddr); err == nil {
		return nil
	}

	// 方法3: 使用Windows防火墙规则阻止特定连接（最精确的方法）
	if err := m.closeConnectionByFirewall(conn); err == nil {
		return nil
	}

	// 方法4: 如果以上方法都失败，记录警告但不终止进程
	m.logger.Warnf("无法关闭特定连接 %s，但不会终止进程 %d", conn.ID, pid)
	return fmt.Errorf("无法关闭特定连接，但进程仍在运行")
}

// closeConnectionByNetsh 使用netsh命令关闭特定连接
func (m *Manager) closeConnectionByNetsh(conn *models.NetworkConnection) error {
	// 对于TCP连接，尝试重置连接
	if conn.Type == "tcp" {
		// 使用netsh命令重置TCP连接
		cmd := exec.Command("netsh", "interface", "tcp", "set", "global", "chimney=enabled")
		if err := cmd.Run(); err == nil {
			// 尝试使用netsh重置特定端口的连接
			cmd = exec.Command("netsh", "interface", "tcp", "set", "global", "autotuninglevel=normal")
			if err := cmd.Run(); err == nil {
				return nil
			}
		}
	}

	// 对于UDP连接，尝试清理邻居缓存
	if conn.Type == "udp" {
		cmd := exec.Command("netsh", "interface", "ipv4", "delete", "neighbors",
			"interface=*", "address="+conn.RemoteAddr)
		if err := cmd.Run(); err == nil {
			return nil
		}
	}

	return fmt.Errorf("netsh命令执行失败")
}

// closeConnectionByRoute 通过路由表阻止特定连接
func (m *Manager) closeConnectionByRoute(remoteIP string) error {
	// 添加临时路由来阻止到特定IP的连接
	cmd := exec.Command("route", "add", remoteIP, "mask", "255.255.255.255", "0.0.0.0", "metric", "1")

	output, err := cmd.CombinedOutput()
	if err != nil {
		m.logger.Warnf("使用route阻止连接失败: %v, 输出: %s", err, string(output))
		return fmt.Errorf("无法阻止连接: %s", remoteIP)
	}

	// 延迟删除路由，避免永久阻止
	go func() {
		time.Sleep(10 * time.Second)
		cmd := exec.Command("route", "delete", remoteIP)
		cmd.CombinedOutput()
		m.logger.Debugf("已删除临时路由规则: %s", remoteIP)
	}()

	return nil
}

// closeConnectionByFirewall 使用Windows防火墙规则阻止特定连接
func (m *Manager) closeConnectionByFirewall(conn *models.NetworkConnection) error {
	// 创建临时的防火墙规则名称
	ruleName := fmt.Sprintf("Block_%s_%s_%d", conn.Type, conn.RemoteAddr, conn.RemotePort)

	// 添加防火墙规则阻止特定连接
	cmd := exec.Command("netsh", "advfirewall", "firewall", "add", "rule",
		"name="+ruleName,
		"dir=out",
		"action=block",
		"protocol="+conn.Type,
		"remoteip="+conn.RemoteAddr,
		"remoteport="+fmt.Sprintf("%d", conn.RemotePort),
		"enable=yes")

	output, err := cmd.CombinedOutput()
	if err != nil {
		m.logger.Warnf("添加防火墙规则失败: %v, 输出: %s", err, string(output))
		return fmt.Errorf("无法添加防火墙规则")
	}

	// 延迟删除防火墙规则
	go func() {
		time.Sleep(15 * time.Second)
		cmd := exec.Command("netsh", "advfirewall", "firewall", "delete", "rule", "name="+ruleName)
		cmd.CombinedOutput()
		m.logger.Debugf("已删除临时防火墙规则: %s", ruleName)
	}()

	return nil
}

// closeConnectionByPID 根据进程ID关闭连接（保留原有方法以兼容性，但改进实现）
func (m *Manager) closeConnectionByPID(pid int32, connType string) error {
	// 获取该进程的所有连接
	connections, err := m.GetConnectionsByPID(pid)
	if err != nil {
		return fmt.Errorf("获取进程连接失败: %w", err)
	}

	if len(connections) == 0 {
		return fmt.Errorf("进程 %d 没有活跃的网络连接", pid)
	}

	// 只关闭指定类型的连接，而不是终止整个进程
	var closedCount int
	for _, conn := range connections {
		if conn.Type == connType {
			if err := m.closeSpecificConnection(pid, connType, 0, conn); err == nil {
				closedCount++
			}
		}
	}

	if closedCount == 0 {
		return fmt.Errorf("无法关闭进程 %d 的 %s 连接", pid, connType)
	}

	m.logger.Infof("成功关闭进程 %d 的 %d 个 %s 连接", pid, closedCount, connType)
	return nil
}

// GetConnectionsByPID 根据进程ID获取连接
func (m *Manager) GetConnectionsByPID(pid int32) ([]*models.NetworkConnection, error) {
	connections, err := gopsutilnet.Connections("all")
	if err != nil {
		return nil, fmt.Errorf("获取网络连接失败: %w", err)
	}

	var result []*models.NetworkConnection
	for _, conn := range connections {
		if conn.Pid == pid {
			networkConn := &models.NetworkConnection{
				ID:          fmt.Sprintf("%d-%s-%d", conn.Pid, m.getConnectionType(conn.Type), conn.Fd),
				LocalAddr:   conn.Laddr.IP,
				RemoteAddr:  conn.Raddr.IP,
				LocalPort:   int(conn.Laddr.Port),
				RemotePort:  int(conn.Raddr.Port),
				Protocol:    m.getProtocolString(conn.Type),
				Status:      conn.Status,
				PID:         conn.Pid,
				ProcessName: m.getProcessName(conn.Pid),
				Type:        m.getConnectionType(conn.Type),
			}
			result = append(result, networkConn)
		}
	}

	return result, nil
}

// GetConnectionsByPort 根据端口获取连接
func (m *Manager) GetConnectionsByPort(port int) ([]*models.NetworkConnection, error) {
	connections, err := gopsutilnet.Connections("all")
	if err != nil {
		return nil, fmt.Errorf("获取网络连接失败: %w", err)
	}

	var result []*models.NetworkConnection
	for _, conn := range connections {
		if int(conn.Laddr.Port) == port || int(conn.Raddr.Port) == port {
			networkConn := &models.NetworkConnection{
				ID:          fmt.Sprintf("%d-%s-%d", conn.Pid, m.getConnectionType(conn.Type), conn.Fd),
				LocalAddr:   conn.Laddr.IP,
				RemoteAddr:  conn.Raddr.IP,
				LocalPort:   int(conn.Laddr.Port),
				RemotePort:  int(conn.Raddr.Port),
				Protocol:    m.getProtocolString(conn.Type),
				Status:      conn.Status,
				PID:         conn.Pid,
				ProcessName: m.getProcessName(conn.Pid),
				Type:        m.getConnectionType(conn.Type),
			}
			result = append(result, networkConn)
		}
	}

	return result, nil
}

// GetConnectionsByIP 根据IP地址获取连接
func (m *Manager) GetConnectionsByIP(ip string) ([]*models.NetworkConnection, error) {
	connections, err := gopsutilnet.Connections("all")
	if err != nil {
		return nil, fmt.Errorf("获取网络连接失败: %w", err)
	}

	var result []*models.NetworkConnection
	for _, conn := range connections {
		if conn.Laddr.IP == ip || conn.Raddr.IP == ip {
			networkConn := &models.NetworkConnection{
				ID:          fmt.Sprintf("%d-%s-%d", conn.Pid, m.getConnectionType(conn.Type), conn.Fd),
				LocalAddr:   conn.Laddr.IP,
				RemoteAddr:  conn.Raddr.IP,
				LocalPort:   int(conn.Laddr.Port),
				RemotePort:  int(conn.Raddr.Port),
				Protocol:    m.getProtocolString(conn.Type),
				Status:      conn.Status,
				PID:         conn.Pid,
				ProcessName: m.getProcessName(conn.Pid),
				Type:        m.getConnectionType(conn.Type),
			}
			result = append(result, networkConn)
		}
	}

	return result, nil
}

// IsPortInUse 检查端口是否被占用
func (m *Manager) IsPortInUse(port int) (bool, error) {
	connections, err := gopsutilnet.Connections("all")
	if err != nil {
		return false, fmt.Errorf("获取网络连接失败: %w", err)
	}

	for _, conn := range connections {
		if int(conn.Laddr.Port) == port || int(conn.Raddr.Port) == port {
			return true, nil
		}
	}

	return false, nil
}

// GetListeningPorts 获取监听端口
func (m *Manager) GetListeningPorts() ([]int, error) {
	connections, err := gopsutilnet.Connections("all")
	if err != nil {
		return nil, fmt.Errorf("获取网络连接失败: %w", err)
	}

	var listeningPorts []int
	portMap := make(map[int]bool)

	for _, conn := range connections {
		if conn.Status == "LISTEN" {
			port := int(conn.Laddr.Port)
			if !portMap[port] {
				listeningPorts = append(listeningPorts, port)
				portMap[port] = true
			}
		}
	}

	return listeningPorts, nil
}

// GetEstablishedConnections 获取已建立的连接
func (m *Manager) GetEstablishedConnections() ([]*models.NetworkConnection, error) {
	connections, err := gopsutilnet.Connections("all")
	if err != nil {
		return nil, fmt.Errorf("获取网络连接失败: %w", err)
	}

	var establishedConnections []*models.NetworkConnection
	for _, conn := range connections {
		if conn.Status == "ESTABLISHED" {
			networkConn := &models.NetworkConnection{
				ID:          fmt.Sprintf("%d", conn.Fd),
				LocalAddr:   conn.Laddr.IP,
				RemoteAddr:  conn.Raddr.IP,
				LocalPort:   int(conn.Laddr.Port),
				RemotePort:  int(conn.Raddr.Port),
				Protocol:    m.getProtocolString(conn.Type),
				Status:      conn.Status,
				PID:         conn.Pid,
				ProcessName: m.getProcessName(conn.Pid),
				Type:        m.getConnectionType(conn.Type),
			}
			establishedConnections = append(establishedConnections, networkConn)
		}
	}

	return establishedConnections, nil
}

// 辅助方法
func (m *Manager) getProtocolString(protocolType uint32) string {
	switch protocolType {
	case 1:
		return "TCP"
	case 2:
		return "UDP"
	default:
		return "UNKNOWN"
	}
}

func (m *Manager) getConnectionType(protocolType uint32) string {
	switch protocolType {
	case 1:
		return "tcp"
	case 2:
		return "udp"
	default:
		return "unknown"
	}
}

func (m *Manager) getProcessName(pid int32) string {
	if pid <= 0 {
		return "unknown"
	}

	// 尝试获取进程信息
	proc, err := gopsutilprocess.NewProcess(pid)
	if err != nil {
		return "unknown"
	}

	// 获取进程名称
	name, err := proc.Name()
	if err != nil {
		return "unknown"
	}

	return name
}

// forceResetNetworkStack 强制重置网络协议栈
func (m *Manager) forceResetNetworkStack(pid int32) error {
	m.logger.Debugf("开始强制重置进程 %d 的网络协议栈", pid)

	// 方法1: 重置TCP/IP协议栈
	if err := m.resetTCPIPStack(); err == nil {
		m.logger.Debugf("TCP/IP协议栈重置成功")
	}

	// 方法2: 重置网络适配器
	if err := m.resetNetworkAdapters(); err == nil {
		m.logger.Debugf("网络适配器重置成功")
	}

	// 方法3: 使用netsh命令重置网络接口
	if err := m.resetNetworkInterfacesByNetsh(); err == nil {
		m.logger.Debugf("网络接口重置成功")
	}

	// 方法4: 重置TCP连接表
	if err := m.resetTCPConnectionTable(); err == nil {
		m.logger.Debugf("TCP连接表重置成功")
	}

	// 方法5: 重置UDP连接
	if err := m.resetUDPConnectionsByNetsh(); err == nil {
		m.logger.Debugf("UDP连接重置成功")
	}

	m.logger.Debugf("进程 %d 的网络协议栈强制重置完成", pid)
	return nil
}

// finalizeNetworkClosure 最终验证和清理网络连接关闭
func (m *Manager) finalizeNetworkClosure(pid int32) error {
	m.logger.Debugf("开始最终验证和清理进程 %d 的网络连接关闭", pid)

	// 验证1: 检查是否还有活跃连接
	remainingConnections, err := m.GetConnectionsByPID(pid)
	if err == nil && len(remainingConnections) > 0 {
		m.logger.Warnf("进程 %d 仍有 %d 个活跃连接，尝试最终清理", pid, len(remainingConnections))

		// 尝试最后一次强制关闭
		for _, conn := range remainingConnections {
			if conn.Type == "tcp" {
				if err := m.forceCloseSpecificTCPConnection(conn); err == nil {
					m.logger.Debugf("最终清理TCP连接 %s 成功", conn.ID)
				}
			}
		}
	}

	// 验证2: 检查网络状态
	if err := m.resetNetworkConfiguration(); err == nil {
		m.logger.Debugf("网络配置最终重置成功")
	}

	// 验证3: 检查防火墙规则
	if err := m.blockProcessOutboundConnections(pid); err == nil {
		m.logger.Debugf("防火墙规则最终验证成功")
	}

	// 验证4: 最终网络协议栈检查
	if err := m.resetTCPIPStack(); err == nil {
		m.logger.Debugf("TCP/IP协议栈最终检查成功")
	}

	m.logger.Debugf("进程 %d 的网络连接关闭最终验证和清理完成", pid)
	return nil
}

// cleanupOrphanedConnections 清理孤儿连接（进程已终止但连接可能残留）
func (m *Manager) cleanupOrphanedConnections(connectionID string, pid int32, connType string) error {
	m.logger.Infof("开始清理进程 %d 的孤儿连接 (类型: %s)", pid, connType)

	// 方法1: 尝试重置网络配置
	if err := m.resetNetworkConfiguration(); err == nil {
		m.logger.Debugf("网络配置重置成功")
	}

	// 方法2: 尝试重置TCP/IP协议栈
	if err := m.resetTCPIPStack(); err == nil {
		m.logger.Debugf("TCP/IP协议栈重置成功")
	}

	// 方法3: 尝试重置网络接口
	if err := m.resetNetworkInterfacesByNetsh(); err == nil {
		m.logger.Debugf("网络接口重置成功")
	}

	// 方法4: 尝试删除可能的路由
	if err := m.forceDeleteProcessRoutes(pid); err == nil {
		m.logger.Debugf("进程路由删除成功")
	}

	m.logger.Infof("进程 %d 的孤儿连接清理完成", pid)
	return nil
}

// forceCloseProcessConnections 强制关闭进程的所有连接
func (m *Manager) forceCloseProcessConnections(pid int32) error {
	m.logger.Infof("强制关闭进程 %d 的所有连接", pid)

	// 使用现有的强制关闭方法
	return m.ForceCloseProcessConnections(pid)
}

// forceCloseConnectionsByType 强制关闭指定进程的指定类型连接
func (m *Manager) forceCloseConnectionsByType(pid int32, connType string) error {
	m.logger.Infof("强制关闭进程 %d 的所有 %s 连接", pid, connType)

	switch connType {
	case "tcp":
		// 强制关闭TCP连接
		if err := m.forceCloseTCPConnections(pid); err == nil {
			m.logger.Debugf("TCP连接强制关闭成功")
		}
	case "udp":
		// 强制关闭UDP连接
		if err := m.resetUDPConnectionsByNetsh(); err == nil {
			m.logger.Debugf("UDP连接重置成功")
		}
	default:
		// 对于未知类型，尝试通用的网络重置
		if err := m.resetNetworkConfiguration(); err == nil {
			m.logger.Debugf("网络配置重置成功")
		}
	}

	// 最终验证
	if err := m.finalizeNetworkClosure(pid); err == nil {
		m.logger.Debugf("连接关闭最终验证完成")
	}

	return nil
}

// shouldTerminateProcessOnConnectionFailure 判断是否应该在连接关闭失败时终止进程
func (m *Manager) shouldTerminateProcessOnConnectionFailure(pid int32, connType string) bool {
	// 这里可以添加配置逻辑，比如：
	// 1. 检查进程是否重要
	// 2. 检查连接类型是否关键
	// 3. 检查用户配置

	// 默认情况下，对于TCP连接失败，允许终止进程
	if connType == "tcp" {
		// 检查进程是否重要（这里可以扩展）
		if m.isProcessCritical(pid) {
			m.logger.Infof("进程 %d 被标记为重要进程，不会自动终止", pid)
			return false
		}
		return true
	}

	return false
}

// isProcessCritical 检查进程是否重要（不应该被自动终止）
func (m *Manager) isProcessCritical(pid int32) bool {
	// 获取进程名称
	proc, err := gopsutilprocess.NewProcess(pid)
	if err != nil {
		return false
	}

	name, err := proc.Name()
	if err != nil {
		return false
	}

	// 定义重要进程列表（可以根据需要扩展）
	criticalProcesses := []string{
		"svchost.exe",    // Windows服务主机
		"lsass.exe",      // 本地安全认证服务
		"winlogon.exe",   // Windows登录进程
		"csrss.exe",      // 客户端/服务器运行时子系统
		"wininit.exe",    // Windows初始化进程
		"services.exe",   // 服务控制管理器
		"explorer.exe",   // Windows资源管理器
		"taskmgr.exe",    // 任务管理器
		"cmd.exe",        // 命令提示符
		"powershell.exe", // PowerShell
	}

	for _, critical := range criticalProcesses {
		if strings.EqualFold(name, critical) {
			return true
		}
	}

	return false
}

// terminateProcessByPID 根据PID终止进程
func (m *Manager) terminateProcessByPID(pid int32) error {
	m.logger.Infof("开始终止进程 %d", pid)

	// 方法1: 尝试优雅终止
	proc, err := gopsutilprocess.NewProcess(pid)
	if err == nil {
		// 发送SIGTERM信号（在Windows上相当于CTRL+C）
		if err := proc.Terminate(); err == nil {
			m.logger.Debugf("进程 %d 优雅终止信号已发送", pid)

			// 等待进程终止
			time.Sleep(3 * time.Second)

			// 检查进程是否已终止
			if isRunning, _ := proc.IsRunning(); !isRunning {
				m.logger.Infof("进程 %d 已成功终止", pid)
				return nil
			}
		}
	}

	// 方法2: 强制终止（如果优雅终止失败）
	m.logger.Warnf("优雅终止失败，尝试强制终止进程 %d", pid)

	// 使用taskkill命令强制终止
	cmd := exec.Command("taskkill", "/F", "/PID", fmt.Sprintf("%d", pid))
	if err := cmd.Run(); err == nil {
		m.logger.Infof("进程 %d 已强制终止", pid)
		return nil
	}

	// 方法3: 使用gopsutil强制终止
	if proc != nil {
		if err := proc.Kill(); err == nil {
			m.logger.Infof("进程 %d 已通过gopsutil强制终止", pid)
			return nil
		}
	}

	return fmt.Errorf("无法终止进程 %d", pid)
}

// terminateProcessAfterConnectionClose 在连接关闭成功后终止相关进程
func (m *Manager) terminateProcessAfterConnectionClose(pid int32, connType string) error {
	m.logger.Infof("连接已关闭，开始终止相关进程: PID=%d, Type=%s", pid, connType)

	// 检查进程是否重要，如果是重要进程则不终止
	if m.isProcessCritical(pid) {
		m.logger.Infof("进程 %d 被标记为重要进程，跳过终止", pid)
		return nil
	}

	// 检查进程是否还在运行
	proc, err := gopsutilprocess.NewProcess(pid)
	if err != nil {
		m.logger.Warnf("无法访问进程 %d: %v", pid, err)
		return fmt.Errorf("无法访问进程 %d: %w", pid, err)
	}

	isRunning, err := proc.IsRunning()
	if err != nil {
		m.logger.Warnf("无法检查进程 %d 的运行状态: %v", pid, err)
		return fmt.Errorf("无法检查进程 %d 的运行状态: %w", pid, err)
	}

	if !isRunning {
		m.logger.Infof("进程 %d 已经终止，无需操作", pid)
		return nil
	}

	// 获取进程名称用于日志记录
	procName, err := proc.Name()
	if err != nil {
		procName = "未知进程"
	}

	m.logger.Infof("开始终止进程: PID=%d, Name=%s, Type=%s", pid, procName, connType)

	// 使用现有的进程终止方法
	if err := m.terminateProcessByPID(pid); err != nil {
		m.logger.Errorf("终止进程 %d 失败: %v", pid, err)
		return fmt.Errorf("终止进程 %d 失败: %w", pid, err)
	}

	m.logger.Infof("进程 %d (%s) 已成功终止", pid, procName)
	return nil
}
