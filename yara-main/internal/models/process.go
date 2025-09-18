package models

import (
	"time"

	"github.com/shirou/gopsutil/v3/process"
)

// ProcessEvent 进程事件
type ProcessEvent struct {
	PID         int32                  `json:"pid"`
	ProcessName string                 `json:"process_name"`
	EventType   string                 `json:"event_type"`
	Timestamp   time.Time              `json:"timestamp"`
	Details     map[string]interface{} `json:"details"`
}

// DeepProcessAnalysis 深度进程分析结果
type DeepProcessAnalysis struct {
	PID              int32                    `json:"pid"`
	Timestamp        time.Time                `json:"timestamp"`
	IsSuspicious     bool                     `json:"is_suspicious"`
	ThreatLevel      string                   `json:"threat_level"`
	RiskScore        int                      `json:"risk_score"`
	RiskFactors      []string                 `json:"risk_factors"`
	BasicInfo        *BasicProcessInfo        `json:"basic_info"`
	MemoryAnalysis   *MemoryAnalysis          `json:"memory_analysis"`
	NetworkAnalysis  *NetworkAnalysis         `json:"network_analysis"`
	FileAnalysis     *FileAnalysis            `json:"file_analysis"`
	BehaviorAnalysis *ProcessBehaviorAnalysis `json:"behavior_analysis"`
	ThreatAnalysis   *ProcessThreatAnalysis   `json:"threat_analysis"`
}

// BasicProcessInfo 基础进程信息
type BasicProcessInfo struct {
	Timestamp        time.Time `json:"timestamp"`
	Name             string    `json:"name"`
	CommandLine      string    `json:"command_line"`
	Executable       string    `json:"executable"`
	WorkingDirectory string    `json:"working_directory"`
	Status           string    `json:"status"`
	CreateTime       time.Time `json:"create_time"`
	Username         string    `json:"username"`
	ParentPID        int32     `json:"parent_pid"`
	CPUPercent       float64   `json:"cpu_percent"`
	MemoryPercent    float32   `json:"memory_percent"`
	NumThreads       int32     `json:"num_threads"`
	Priority         int32     `json:"priority"`
}

// MemoryAnalysis 内存分析
type MemoryAnalysis struct {
	Timestamp     time.Time                `json:"timestamp"`
	RSS           uint64                   `json:"rss"`
	VMS           uint64                   `json:"vms"`
	Swap          uint64                   `json:"swap"`
	MemoryPercent float32                  `json:"memory_percent"`
	MemoryMaps    []process.MemoryMapsStat `json:"memory_maps"`
	MemoryInfo    *process.MemoryInfoStat  `json:"memory_info"`
}

// ConnectionStat 连接状态
type ConnectionStat struct {
	Fd     uint32    `json:"fd"`
	Family uint32    `json:"family"`
	Type   uint32    `json:"type"`
	Laddr  *AddrInfo `json:"laddr"`
	Raddr  *AddrInfo `json:"raddr"`
	Status string    `json:"status"`
}

// NetworkAnalysis 网络分析
type NetworkAnalysis struct {
	Timestamp       time.Time               `json:"timestamp"`
	Connections     []ConnectionStat        `json:"connections"`
	ConnectionCount int                     `json:"connection_count"`
	NetworkIO       *process.IOCountersStat `json:"network_io"`
}

// FileAnalysis 文件分析
type FileAnalysis struct {
	Timestamp       time.Time               `json:"timestamp"`
	OpenFiles       []process.OpenFilesStat `json:"open_files"`
	OpenFileCount   int                     `json:"open_file_count"`
	FileDescriptors []FileDescriptor        `json:"file_descriptors"`
}

// FileDescriptor 文件描述符
type FileDescriptor struct {
	Path string `json:"path"`
	Type string `json:"type"`
}

// ProcessBehaviorAnalysis 进程行为分析
type ProcessBehaviorAnalysis struct {
	Timestamp             time.Time `json:"timestamp"`
	HasNetworkActivity    bool      `json:"has_network_activity"`
	HasFileOperations     bool      `json:"has_file_operations"`
	HasRegistryOperations bool      `json:"has_registry_operations"`
	HasProcessCreation    bool      `json:"has_process_creation"`
	HasAPIHooking         bool      `json:"has_api_hooking"`
	HasAntiDebug          bool      `json:"has_anti_debug"`
	HasVMDetection        bool      `json:"has_vm_detection"`
	HasCodeInjection      bool      `json:"has_code_injection"`
	HasDLLInjection       bool      `json:"has_dll_injection"`
	HasThreadInjection    bool      `json:"has_thread_injection"`
}

// ProcessThreatAnalysis 进程威胁分析
type ProcessThreatAnalysis struct {
	Timestamp     time.Time `json:"timestamp"`
	IsSuspicious  bool      `json:"is_suspicious"`
	ThreatTypes   []string  `json:"threat_types"`
	ThreatLevel   string    `json:"threat_level"`
	MalwareFamily string    `json:"malware_family"`
	DetectionRate float64   `json:"detection_rate"`
	Signatures    []string  `json:"signatures"`
}

// ProcessModuleInfo 进程模块信息
type ProcessModuleInfo struct {
	PID          int32    `json:"pid"`
	Name         string   `json:"name"`
	Path         string   `json:"path"`
	Size         uint64   `json:"size"`
	BaseAddress  uint64   `json:"base_address"`
	EntryPoint   uint64   `json:"entry_point"`
	IsSuspicious bool     `json:"is_suspicious"`
	ThreatLevel  string   `json:"threat_level"`
	ThreatTypes  []string `json:"threat_types"`
}

// ProcessThread 进程线程信息
type ProcessThread struct {
	ID           int32   `json:"id"`
	UserTime     float64 `json:"user_time"`
	SystemTime   float64 `json:"system_time"`
	Status       string  `json:"status"`
	Priority     int32   `json:"priority"`
	IsSuspicious bool    `json:"is_suspicious"`
}

// ProcessHandle 进程句柄信息
type ProcessHandle struct {
	ID           int32  `json:"id"`
	Name         string `json:"name"`
	Type         string `json:"type"`
	Access       string `json:"access"`
	IsSuspicious bool   `json:"is_suspicious"`
}

// ProcessEnvironment 进程环境信息
type ProcessEnvironment struct {
	Variables      map[string]string `json:"variables"`
	IsSuspicious   bool              `json:"is_suspicious"`
	SuspiciousVars []string          `json:"suspicious_vars"`
}

// ProcessSecurity 进程安全信息
type ProcessSecurity struct {
	IntegrityLevel string                 `json:"integrity_level"`
	Privileges     []string               `json:"privileges"`
	TokenInfo      map[string]interface{} `json:"token_info"`
	IsElevated     bool                   `json:"is_elevated"`
	IsSuspicious   bool                   `json:"is_suspicious"`
}

// ProcessPerformance 进程性能信息
type ProcessPerformance struct {
	Timestamp       time.Time `json:"timestamp"`
	CPUPercent      float64   `json:"cpu_percent"`
	MemoryPercent   float32   `json:"memory_percent"`
	IOReadBytes     uint64    `json:"io_read_bytes"`
	IOWriteBytes    uint64    `json:"io_write_bytes"`
	IOReadCount     uint64    `json:"io_read_count"`
	IOWriteCount    uint64    `json:"io_write_count"`
	ContextSwitches uint64    `json:"context_switches"`
	PageFaults      uint64    `json:"page_faults"`
}

// ProcessDependencies 进程依赖信息
type ProcessDependencies struct {
	Imports        []string `json:"imports"`
	Exports        []string `json:"exports"`
	Libraries      []string `json:"libraries"`
	Functions      []string `json:"functions"`
	SuspiciousAPIs []string `json:"suspicious_apis"`
}

// ProcessRegistry 进程注册表信息
type ProcessRegistry struct {
	KeysAccessed   []string `json:"keys_accessed"`
	ValuesModified []string `json:"values_modified"`
	IsSuspicious   bool     `json:"is_suspicious"`
	SuspiciousKeys []string `json:"suspicious_keys"`
}

// ProcessNetwork 进程网络信息
type ProcessNetwork struct {
	Connections           []ConnectionStat `json:"connections"`
	ListenPorts           []int            `json:"listen_ports"`
	RemoteHosts           []string         `json:"remote_hosts"`
	Protocols             []string         `json:"protocols"`
	IsSuspicious          bool             `json:"is_suspicious"`
	SuspiciousConnections []string         `json:"suspicious_connections"`
}

// ProcessFileSystem 进程文件系统信息
type ProcessFileSystem struct {
	FilesCreated    []string `json:"files_created"`
	FilesModified   []string `json:"files_modified"`
	FilesDeleted    []string `json:"files_deleted"`
	Directories     []string `json:"directories"`
	IsSuspicious    bool     `json:"is_suspicious"`
	SuspiciousFiles []string `json:"suspicious_files"`
}

// ProcessTimeline 进程时间线
type ProcessTimeline struct {
	Events     []ProcessTimelineEvent `json:"events"`
	StartTime  time.Time              `json:"start_time"`
	EndTime    time.Time              `json:"end_time"`
	Duration   time.Duration          `json:"duration"`
	EventCount int                    `json:"event_count"`
}

// ProcessTimelineEvent 进程时间线事件
type ProcessTimelineEvent struct {
	Timestamp    time.Time              `json:"timestamp"`
	EventType    string                 `json:"event_type"`
	Description  string                 `json:"description"`
	Details      map[string]interface{} `json:"details"`
	Severity     string                 `json:"severity"`
	IsSuspicious bool                   `json:"is_suspicious"`
}

// SystemModule 系统模块信息
type SystemModule struct {
	Name              string    `json:"name"`
	ProcessCount      int       `json:"process_count"`
	TotalMemory       uint64    `json:"total_memory"`
	Status            string    `json:"status"`
	LastAccessed      time.Time `json:"last_accessed"`
	AffectedProcesses []int32   `json:"affected_processes"`
	Path              string    `json:"path,omitempty"`
	Version           string    `json:"version,omitempty"`
	Company           string    `json:"company,omitempty"`
	Description       string    `json:"description,omitempty"`
}

// ProcessDetails 进程详细信息
type ProcessDetails struct {
	PID         int32             `json:"pid"`
	BasicInfo   *BasicProcessInfo `json:"basic_info"`
	MemoryInfo  *MemoryInfo       `json:"memory_info"`
	CPUInfo     *CPUInfo          `json:"cpu_info"`
	FileInfo    *FileInfo         `json:"file_info"`
	NetworkInfo *NetworkInfo      `json:"network_info"`
	ModuleInfo  *ModuleInfo       `json:"module_info"`
	ThreadInfo  *ThreadInfo       `json:"thread_info"`
	HandleInfo  *HandleInfo       `json:"handle_info"`
	UpdateTime  time.Time         `json:"update_time"`
}

// CPUInfo CPU信息
type CPUInfo struct {
	CPUPercent float64   `json:"cpu_percent"`
	NumThreads int32     `json:"num_threads"`
	Priority   int32     `json:"priority"`
	UpdateTime time.Time `json:"update_time"`
}

// ModuleInfo 模块信息
type ModuleInfo struct {
	Modules     []ModuleDetails `json:"modules"`
	ModuleCount int             `json:"module_count"`
	UpdateTime  time.Time       `json:"update_time"`
}

// ModuleDetails 模块详细信息
type ModuleDetails struct {
	Name        string `json:"name"`
	Path        string `json:"path"`
	Size        uint64 `json:"size"`
	BaseAddress uint64 `json:"base_address"`
}

// ThreadInfo 线程信息
type ThreadInfo struct {
	Threads     []ThreadDetails `json:"threads"`
	ThreadCount int             `json:"thread_count"`
	UpdateTime  time.Time       `json:"update_time"`
}

// ThreadDetails 线程详细信息
type ThreadDetails struct {
	ID         int32   `json:"id"`
	UserTime   float64 `json:"user_time"`
	SystemTime float64 `json:"system_time"`
	Status     string  `json:"status"`
	Priority   int32   `json:"priority"`
}

// HandleInfo 句柄信息
type HandleInfo struct {
	Handles     []HandleDetails `json:"handles"`
	HandleCount int             `json:"handle_count"`
	UpdateTime  time.Time       `json:"update_time"`
}

// HandleDetails 句柄详细信息
type HandleDetails struct {
	ID     int32  `json:"id"`
	Name   string `json:"name"`
	Type   string `json:"type"`
	Access string `json:"access"`
}

// NetworkInfo 网络信息
type NetworkInfo struct {
	Connections []ConnectionInfo `json:"connections"`
	ConnCount   int              `json:"conn_count"`
	UpdateTime  time.Time        `json:"update_time"`
}
