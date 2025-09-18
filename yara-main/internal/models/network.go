package models

import (
	"time"
)

// NetworkConnection 网络连接信息
type NetworkConnection struct {
	ID          string `json:"id"`
	LocalAddr   string `json:"local_addr"`
	RemoteAddr  string `json:"remote_addr"`
	LocalPort   int    `json:"local_port"`
	RemotePort  int    `json:"remote_port"`
	Protocol    string `json:"protocol"`
	Status      string `json:"status"`
	PID         int32  `json:"pid"`
	ProcessName string `json:"process_name"`
	Type        string `json:"type"`
}

// NetworkPacket 网络数据包
type NetworkPacket struct {
	SourceIP      string    `json:"source_ip"`
	DestIP        string    `json:"dest_ip"`
	SourcePort    int       `json:"source_port"`
	DestPort      int       `json:"dest_port"`
	Protocol      string    `json:"protocol"`
	Payload       []byte    `json:"payload"`
	Timestamp     time.Time `json:"timestamp"`
	PacketSize    int       `json:"packet_size"`
	IsSuspicious  bool      `json:"is_suspicious"`
	ThreatLevel   string    `json:"threat_level"`
	ThreatDetails []string  `json:"threat_details"`
}

// NetworkStats 网络统计信息
type NetworkStats struct {
	Timestamp    time.Time          `json:"timestamp"`
	Monitored    bool               `json:"monitored"`
	Interfaces   []NetworkInterface `json:"interfaces"`
	TotalBytes   int64              `json:"total_bytes"`
	TotalPackets int64              `json:"total_packets"`
	ThreatCount  int                `json:"threat_count"`
}

// NetworkInterface 网络接口信息
type NetworkInterface struct {
	Name           string                 `json:"name"`
	Index          int                    `json:"index"`
	Addresses      []string               `json:"addresses"`
	IPv4Addresses  []string               `json:"ipv4_addresses,omitempty"`
	IPv6Addresses  []string               `json:"ipv6_addresses,omitempty"`
	MACAddress     string                 `json:"mac_address,omitempty"`
	IsUp           bool                   `json:"is_up"`
	IsLoopback     bool                   `json:"is_loopback"`
	IsMulticast    bool                   `json:"is_multicast,omitempty"`
	IsBroadcast    bool                   `json:"is_broadcast,omitempty"`
	IsPointToPoint bool                   `json:"is_point_to_point,omitempty"`
	MTU            int                    `json:"mtu"`
	Speed          int64                  `json:"speed"`
	Flags          string                 `json:"flags,omitempty"`
	HardwareAddr   string                 `json:"hardware_addr,omitempty"`
	WindowsInfo    map[string]interface{} `json:"windows_info,omitempty"`
	Error          string                 `json:"error,omitempty"`
}

// NetworkThreat 网络威胁信息
type NetworkThreat struct {
	ID          string    `json:"id"`
	Type        string    `json:"type"`
	Severity    string    `json:"severity"`
	SourceIP    string    `json:"source_ip"`
	DestIP      string    `json:"dest_ip"`
	Protocol    string    `json:"protocol"`
	Description string    `json:"description"`
	Timestamp   time.Time `json:"timestamp"`
	Details     []string  `json:"details"`
}

// NetworkScanResult 网络扫描结果
type NetworkScanResult struct {
	URL          string            `json:"url"`
	IsSuspicious bool              `json:"is_suspicious"`
	Threats      []ThreatInfo      `json:"threats"`
	ScanTime     time.Time         `json:"scan_time"`
	ScanDuration time.Duration     `json:"scan_duration"`
	ThreatLevel  string            `json:"threat_level"`
	Category     string            `json:"category"`
	Metadata     map[string]string `json:"metadata,omitempty"`
}
