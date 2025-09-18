package models

import (
	"time"
)

// Response 通用响应结构
type Response struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
	Time    time.Time   `json:"time"`
}

// FileInfo 文件信息
type FileInfo struct {
	Path         string    `json:"path"`
	Name         string    `json:"name"`
	Size         int64     `json:"size"`
	IsDir        bool      `json:"is_dir"`
	ModTime      time.Time `json:"mod_time"`
	CreateTime   time.Time `json:"create_time"`
	AccessTime   time.Time `json:"access_time"`
	Permissions  string    `json:"permissions"`
	Owner        string    `json:"owner"`
	Group        string    `json:"group"`
	MD5          string    `json:"md5,omitempty"`
	SHA256       string    `json:"sha256,omitempty"`
	IsSuspicious bool      `json:"is_suspicious,omitempty"`
	ThreatLevel  string    `json:"threat_level,omitempty"`
}

// ScanResult 扫描结果
type ScanResult struct {
	FilePath     string            `json:"file_path"`
	IsInfected   bool              `json:"is_infected"`
	Threats      []ThreatInfo      `json:"threats"`
	ScanTime     time.Time         `json:"scan_time"`
	ScanDuration time.Duration     `json:"scan_duration"`
	FileInfo     *FileInfo         `json:"file_info"`
	Metadata     map[string]string `json:"metadata,omitempty"`
}

// ThreatInfo 威胁信息
type ThreatInfo struct {
	RuleName    string `json:"rule_name"`
	Description string `json:"description"`
	Severity    string `json:"severity"`
	Category    string `json:"category"`
	Tags        string `json:"tags"`
}

// ProcessInfo 进程信息
type ProcessInfo struct {
	PID         int32             `json:"pid"`
	Name        string            `json:"name"`
	Exe         string            `json:"exe"`
	Cmdline     string            `json:"cmdline"`
	Cwd         string            `json:"cwd"`
	Status      string            `json:"status"`
	CPUPercent  float64           `json:"cpu_percent"`
	MemoryInfo  *MemoryInfo       `json:"memory_info"`
	CreateTime  time.Time         `json:"create_time"`
	Username    string            `json:"username"`
	PPID        int32             `json:"ppid"`
	NumThreads  int32             `json:"num_threads"`
	NumFiles    int32             `json:"num_files"`
	Priority    int32             `json:"priority"`
	Modules     []string          `json:"modules"`
	Children    []int32           `json:"children"`
	Connections []*ConnectionInfo `json:"connections"`

	// 新增字段
	Architecture     string             `json:"architecture"`       // 进程体系结构 (x64, x86, ARM64等)
	MemoryPercent    float64            `json:"memory_percent"`     // 内存使用百分比
	WorkingSet       uint64             `json:"working_set"`        // 工作集大小
	PrivateBytes     uint64             `json:"private_bytes"`      // 私有字节数
	PeakWorkingSet   uint64             `json:"peak_working_set"`   // 峰值工作集
	PeakPrivateBytes uint64             `json:"peak_private_bytes"` // 峰值私有字节
	PageFaults       uint64             `json:"page_faults"`        // 页面错误数
	IOCounters       *ProcessIOCounters `json:"io_counters"`        // IO计数器
	ContextSwitches  uint64             `json:"context_switches"`   // 上下文切换数
	HandleCount      uint32             `json:"handle_count"`       // 句柄数量
	SessionID        uint32             `json:"session_id"`         // 会话ID
	IntegrityLevel   string             `json:"integrity_level"`    // 完整性级别
	Elevated         bool               `json:"elevated"`           // 是否提升权限
	DEPEnabled       bool               `json:"dep_enabled"`        // 数据执行保护是否启用
	ASLREnabled      bool               `json:"aslr_enabled"`       // 地址空间布局随机化是否启用
	LastUpdate       time.Time          `json:"last_update"`        // 最后更新时间
}

// MemoryInfo 内存信息
type MemoryInfo struct {
	RSS    uint64 `json:"rss"`
	VMS    uint64 `json:"vms"`
	HWM    uint64 `json:"hwm"`
	Data   uint64 `json:"data"`
	Stack  uint64 `json:"stack"`
	Locked uint64 `json:"locked"`
	Swap   uint64 `json:"swap"`
}

// ConnectionInfo 连接信息
type ConnectionInfo struct {
	FD     uint32    `json:"fd"`
	Family int32     `json:"family"`
	Type   int32     `json:"type"`
	Laddr  *AddrInfo `json:"laddr"`
	Raddr  *AddrInfo `json:"raddr"`
	Status string    `json:"status"`
}

// AddrInfo 地址信息
type AddrInfo struct {
	IP   string `json:"ip"`
	Port uint32 `json:"port"`
}

// ProcessIOCounters 进程IO计数器
type ProcessIOCounters struct {
	ReadCount  uint64 `json:"read_count"`  // 读取操作次数
	WriteCount uint64 `json:"write_count"` // 写入操作次数
	ReadBytes  uint64 `json:"read_bytes"`  // 读取字节数
	WriteBytes uint64 `json:"write_bytes"` // 写入字节数
	OtherCount uint64 `json:"other_count"` // 其他操作次数
	OtherBytes uint64 `json:"other_bytes"` // 其他操作字节数
}

// SystemMemoryInfo 系统内存信息
type SystemMemoryInfo struct {
	Total     uint64  `json:"total"`
	Available uint64  `json:"available"`
	Used      uint64  `json:"used"`
	Free      uint64  `json:"free"`
	Percent   float64 `json:"percent"`
}

// RegistryKey 注册表键
type RegistryKey struct {
	Path         string            `json:"path"`
	Name         string            `json:"name"`
	Type         string            `json:"type"`
	Value        interface{}       `json:"value"`
	SubKeys      []string          `json:"sub_keys,omitempty"`
	Values       map[string]string `json:"values,omitempty"`
	LastModified time.Time         `json:"last_modified"`
}

// SecurityStatus 安全状态
type SecurityStatus struct {
	IsProtected        bool                   `json:"is_protected"`
	LastScanTime       time.Time              `json:"last_scan_time"`
	ThreatCount        int                    `json:"threat_count"`
	QuarantineCount    int                    `json:"quarantine_count"`
	RealTimeProtection bool                   `json:"real_time_protection"`
	DefinitionsVersion string                 `json:"definitions_version"`
	LastUpdateTime     time.Time              `json:"last_update_time"`
	ScanStatistics     *ScanStatistics        `json:"scan_statistics,omitempty"`
	RulesInfo          map[string]interface{} `json:"rules_info,omitempty"`
}

// DirectoryScanRequest 目录扫描请求
type DirectoryScanRequest struct {
	Path            string   `json:"path"`
	Recursive       bool     `json:"recursive"`
	MaxDepth        int      `json:"max_depth"`
	IncludePatterns []string `json:"include_patterns,omitempty"`
	ExcludePatterns []string `json:"exclude_patterns,omitempty"`
}

// ProcessOperationRequest 进程操作请求
type ProcessOperationRequest struct {
	PID       int32  `json:"pid"`
	Operation string `json:"operation"`
	Timeout   int    `json:"timeout,omitempty"`
}

// RegistryOperationRequest 注册表操作请求
type RegistryOperationRequest struct {
	Path  string      `json:"path"`
	Name  string      `json:"name,omitempty"`
	Type  string      `json:"type,omitempty"`
	Value interface{} `json:"value,omitempty"`
}

// NetworkOperationRequest 网络操作请求
type NetworkOperationRequest struct {
	ConnectionID string `json:"connection_id"`
	Operation    string `json:"operation"`
}

// UserInfo 用户信息
type UserInfo struct {
	UID         string    `json:"uid"`
	GID         string    `json:"gid"`
	Username    string    `json:"username"`
	Name        string    `json:"name"`
	HomeDir     string    `json:"home_dir"`
	Shell       string    `json:"shell"`
	LastLogin   time.Time `json:"last_login"`
	IsActive    bool      `json:"is_active"`
	IsAdmin     bool      `json:"is_admin"`
	PasswordAge int       `json:"password_age"`
	Groups      []string  `json:"groups,omitempty"`
}

// ErrorResponse 错误响应
type ErrorResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Details string `json:"details,omitempty"`
}

// UserPermissions 用户权限
type UserPermissions struct {
	UserID      string          `json:"user_id"`
	Username    string          `json:"username"`
	IsAdmin     bool            `json:"is_admin"`
	Permissions map[string]bool `json:"permissions"`
	CheckTime   time.Time       `json:"check_time"`
}

// PasswordPolicy 密码策略
type PasswordPolicy struct {
	MinLength           int  `json:"min_length"`
	RequireUppercase    bool `json:"require_uppercase"`
	RequireLowercase    bool `json:"require_lowercase"`
	RequireNumbers      bool `json:"require_numbers"`
	RequireSpecialChars bool `json:"require_special_chars"`
	MaxAge              int  `json:"max_age"`
	HistoryCount        int  `json:"history_count"`
	LockoutThreshold    int  `json:"lockout_threshold"`
	LockoutDuration     int  `json:"lockout_duration"`
}

// AccountStatus 账户状态
type AccountStatus struct {
	Username          string    `json:"username"`
	IsActive          bool      `json:"is_active"`
	IsLocked          bool      `json:"is_locked"`
	IsPasswordExpired bool      `json:"is_password_expired"`
	LastLogin         time.Time `json:"last_login"`
	CheckTime         time.Time `json:"check_time"`
}

// UserSession 用户会话
type UserSession struct {
	SessionID  string    `json:"session_id"`
	Username   string    `json:"username"`
	LoginTime  time.Time `json:"login_time"`
	LastActive time.Time `json:"last_active"`
	IPAddress  string    `json:"ip_address"`
	UserAgent  string    `json:"user_agent"`
	IsActive   bool      `json:"is_active"`
}

// GroupInfo 组信息
type GroupInfo struct {
	GID      string `json:"gid"`
	Name     string `json:"name"`
	Username string `json:"username"`
}

// LoginRecord 登录记录
type LoginRecord struct {
	Username   string    `json:"username"`
	LoginTime  time.Time `json:"login_time"`
	LogoutTime time.Time `json:"logout_time,omitempty"`
	IPAddress  string    `json:"ip_address"`
	UserAgent  string    `json:"user_agent"`
	Status     string    `json:"status"`
}

// QuarantineItem 隔离项目
type QuarantineItem struct {
	FilePath       string    `json:"file_path"`
	OriginalPath   string    `json:"original_path"`
	QuarantineTime time.Time `json:"quarantine_time"`
	Reason         string    `json:"reason"`
	Status         string    `json:"status"`
}

// ScanStatistics 扫描统计
type ScanStatistics struct {
	TotalScans    int64         `json:"total_scans"`
	InfectedFiles int64         `json:"infected_files"`
	CleanFiles    int64         `json:"clean_files"`
	FailedScans   int64         `json:"failed_scans"`
	TotalScanTime time.Duration `json:"total_scan_time"`
}

// ProcessStatistics 进程统计
type ProcessStatistics struct {
	TotalProcesses   int64   `json:"total_processes"`
	RunningProcesses int64   `json:"running_processes"`
	TotalMemoryUsage int64   `json:"total_memory_usage"`
	TotalCPUUsage    float64 `json:"total_cpu_usage"`
	MonitoredCount   int64   `json:"monitored_count"`
}

// ProcessMemoryInfo 进程内存信息
type ProcessMemoryInfo struct {
	PID       int32   `json:"pid"`
	RSS       uint64  `json:"rss"`
	VMS       uint64  `json:"vms"`
	Percent   float64 `json:"percent"`
	Available uint64  `json:"available"`
	Used      uint64  `json:"used"`
	Free      uint64  `json:"free"`
	Total     uint64  `json:"total"`
}

// PermissionResult 权限检查结果
type PermissionResult struct {
	Username    string            `json:"username"`
	Permissions map[string]bool   `json:"permissions"`
	IsAdmin     bool              `json:"is_admin"`
	CheckTime   time.Time         `json:"check_time"`
	Details     map[string]string `json:"details,omitempty"`
}

// BufferScanRequest 缓冲区扫描请求
type BufferScanRequest struct {
	Data       string `json:"data"`
	Identifier string `json:"identifier"`
}

// FileListResponse 文件列表响应
type FileListResponse struct {
	Directory string     `json:"directory"`
	Files     []FileInfo `json:"files"`
	Count     int        `json:"count"`
	Error     string     `json:"error,omitempty"`
}

// ProcessListResponse 进程列表响应
type ProcessListResponse struct {
	Processes []ProcessInfo `json:"processes"`
	Count     int           `json:"count"`
	Error     string        `json:"error,omitempty"`
}

// NetworkConnectionsResponse 网络连接响应
type NetworkConnectionsResponse struct {
	Connections []NetworkConnection `json:"connections"`
	Count       int                 `json:"count"`
	Error       string              `json:"error,omitempty"`
}

// RegistrySearchRequest 注册表搜索请求
type RegistrySearchRequest struct {
	Path      string `json:"path"`
	Pattern   string `json:"pattern"`
	Recursive bool   `json:"recursive"`
	MaxDepth  int    `json:"max_depth"`
}

// RegistrySearchResponse 注册表搜索响应
type RegistrySearchResponse struct {
	Results []RegistryKey `json:"results"`
	Count   int           `json:"count"`
	Error   string        `json:"error,omitempty"`
}

// PasswordValidationRequest 密码验证请求
type PasswordValidationRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// PasswordValidationResponse 密码验证响应
type PasswordValidationResponse struct {
	IsValid     bool     `json:"is_valid"`
	Errors      []string `json:"errors,omitempty"`
	Suggestions []string `json:"suggestions,omitempty"`
}

// PasswordChangeRequest 密码更改请求
type PasswordChangeRequest struct {
	Username    string `json:"username"`
	OldPassword string `json:"old_password"`
	NewPassword string `json:"new_password"`
}

// UserLockRequest 用户锁定请求
type UserLockRequest struct {
	Username string `json:"username"`
	Reason   string `json:"reason,omitempty"`
}

// SessionKillRequest 会话终止请求
type SessionKillRequest struct {
	SessionID string `json:"session_id"`
	Reason    string `json:"reason,omitempty"`
}

// CacheStats 缓存统计
type CacheStats struct {
	HitCount      int64     `json:"hit_count"`
	MissCount     int64     `json:"miss_count"`
	HitRate       float64   `json:"hit_rate"`
	Size          int64     `json:"size"`
	MaxSize       int64     `json:"max_size"`
	EvictionCount int64     `json:"eviction_count"`
	LastCleanup   time.Time `json:"last_cleanup"`
}

// ScanHistoryItem 扫描历史项
type ScanHistoryItem struct {
	ID           string        `json:"id"`
	FilePath     string        `json:"file_path"`
	ScanTime     time.Time     `json:"scan_time"`
	IsInfected   bool          `json:"is_infected"`
	ThreatCount  int           `json:"threat_count"`
	ScanDuration time.Duration `json:"scan_duration"`
	Status       string        `json:"status"`
}

// 请求模型定义
type ScanFileRequest struct {
	Path string `json:"path" binding:"required"`
}

type StartProcessRequest struct {
	Command    string   `json:"command" binding:"required"`
	Args       []string `json:"args"`
	WorkingDir string   `json:"working_dir"`
}

type CloseConnectionRequest struct {
	ConnectionID string `json:"connection_id"`
	ProcessName  string `json:"process_name"`
}

type CloseAllConnectionsRequest struct {
	ProcessName string `json:"process_name"`
}

type ForceCloseConnectionsRequest struct {
	ProcessName string `json:"process_name"`
}

type CreateRegistryKeyRequest struct {
	Path string `json:"path" binding:"required"`
}

type SetRegistryValueRequest struct {
	Path  string      `json:"path" binding:"required"`
	Name  string      `json:"name" binding:"required"`
	Value interface{} `json:"value" binding:"required"`
	Type  string      `json:"type" binding:"required"`
}

type SearchRegistryRequest struct {
	Root         string `json:"root" binding:"required"`
	Pattern      string `json:"pattern"`
	ValuePattern string `json:"value_pattern"`
}

type QuarantineFileRequest struct {
	FilePath string `json:"file_path" binding:"required"`
	Reason   string `json:"reason"`
}

type RestoreFileRequest struct {
	FilePath string `json:"file_path" binding:"required"`
}

type CheckPermissionsRequest struct {
	Username    string   `json:"username" binding:"required"`
	Permissions []string `json:"permissions" binding:"required"`
}

type ValidatePasswordRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type KillSessionRequest struct {
	Username  string `json:"username" binding:"required"`
	SessionID string `json:"session_id" binding:"required"`
}

type ChangePasswordRequest struct {
	Username    string `json:"username" binding:"required"`
	OldPassword string `json:"old_password" binding:"required"`
	NewPassword string `json:"new_password" binding:"required"`
}
