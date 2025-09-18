package process

import (
	"fmt"
	"runtime"
	"unsafe"

	"golang.org/x/sys/windows"
)

// Windows API constants for process information
const (
	PROCESS_QUERY_INFORMATION         = 0x0400
	PROCESS_QUERY_LIMITED_INFORMATION = 0x1000
	PROCESS_VM_READ                   = 0x0010

	// Memory information
	PROCESS_MEMORY_COUNTERS_EX_INFO = 0x00000001

	// Process information
	PROCESS_BASIC_INFORMATION_INFO          = 0
	PROCESS_EXTENDED_BASIC_INFORMATION_INFO = 0x00000001

	// Integrity levels
	SECURITY_MANDATORY_UNTRUSTED_RID         = 0x00000000
	SECURITY_MANDATORY_LOW_RID               = 0x00001000
	SECURITY_MANDATORY_MEDIUM_RID            = 0x00002000
	SECURITY_MANDATORY_MEDIUM_PLUS_RID       = 0x00002100
	SECURITY_MANDATORY_HIGH_RID              = 0x00003000
	SECURITY_MANDATORY_SYSTEM_RID            = 0x00004000
	SECURITY_MANDATORY_PROTECTED_PROCESS_RID = 0x00005000
)

// Windows API structures
type PROCESS_MEMORY_COUNTERS_EX struct {
	cb                        uint32
	PageFaultCount            uint32
	PeakWorkingSetSize        uint64
	WorkingSetSize            uint64
	PeakPagefileUsage         uint64
	PagefileUsage             uint64
	PeakWorkingSetSizePrivate uint64
	WorkingSetSizePrivate     uint64
	PeakPagefileUsagePrivate  uint64
	PagefileUsagePrivate      uint64
	PrivateUsage              uint64
}

type PROCESS_BASIC_INFORMATION struct {
	ExitStatus                   uint32
	PebBaseAddress               uintptr
	AffinityMask                 uintptr
	BasePriority                 uint32
	UniqueProcessId              uintptr
	InheritedFromUniqueProcessId uintptr
}

type PROCESS_EXTENDED_BASIC_INFORMATION struct {
	BasicInfo          PROCESS_BASIC_INFORMATION
	IsProtectedProcess uint32
	IsWow64Process     uint32
}

type IO_COUNTERS struct {
	ReadOperationCount  uint64
	WriteOperationCount uint64
	OtherOperationCount uint64
	ReadTransferCount   uint64
	WriteTransferCount  uint64
	OtherTransferCount  uint64
}

// Lazy DLL loading
var (
	ntdll    = windows.NewLazySystemDLL("ntdll.dll")
	kernel32 = windows.NewLazySystemDLL("kernel32.dll")
	psapi    = windows.NewLazySystemDLL("psapi.dll")
)

// LazyProc declarations
var (
	ntQueryInformationProcess = ntdll.NewProc("NtQueryInformationProcess")
	ntQuerySystemInformation  = ntdll.NewProc("NtQuerySystemInformation")
	getProcessMemoryInfo      = psapi.NewProc("GetProcessMemoryInfo")
	getProcessIoCounters      = kernel32.NewProc("GetProcessIoCounters")
	getProcessIdOfThread      = kernel32.NewProc("GetProcessIdOfThread")
	getProcessTimes           = kernel32.NewProc("GetProcessTimes")
)

// ProcessInfoEnhancer 进程信息增强器
type ProcessInfoEnhancer struct {
	logger interface {
		Debugf(format string, args ...interface{})
		Warnf(format string, args ...interface{})
		Errorf(format string, args ...interface{})
	}
}

// NewProcessInfoEnhancer 创建进程信息增强器
func NewProcessInfoEnhancer(logger interface {
	Debugf(format string, args ...interface{})
	Warnf(format string, args ...interface{})
	Errorf(format string, args ...interface{})
}) *ProcessInfoEnhancer {
	return &ProcessInfoEnhancer{
		logger: logger,
	}
}

// GetProcessArchitecture 获取进程体系结构
func (e *ProcessInfoEnhancer) GetProcessArchitecture(pid int32) (string, error) {
	if runtime.GOOS != "windows" {
		return "unknown", fmt.Errorf("only supported on Windows")
	}

	handle, err := windows.OpenProcess(PROCESS_QUERY_INFORMATION, false, uint32(pid))
	if err != nil {
		// 尝试使用受限权限
		handle, err = windows.OpenProcess(PROCESS_QUERY_LIMITED_INFORMATION, false, uint32(pid))
		if err != nil {
			return "unknown", fmt.Errorf("无法打开进程句柄: %w", err)
		}
	}
	defer windows.CloseHandle(handle)

	var basicInfo PROCESS_BASIC_INFORMATION
	var returnLength uint32

	ret, _, err := ntQueryInformationProcess.Call(
		uintptr(handle),
		PROCESS_BASIC_INFORMATION_INFO,
		uintptr(unsafe.Pointer(&basicInfo)),
		unsafe.Sizeof(basicInfo),
		uintptr(unsafe.Pointer(&returnLength)),
	)

	if ret != 0 {
		return "unknown", fmt.Errorf("查询进程信息失败: %w", err)
	}

	// 检查是否为WOW64进程
	var extendedInfo PROCESS_EXTENDED_BASIC_INFORMATION
	ret, _, err = ntQueryInformationProcess.Call(
		uintptr(handle),
		PROCESS_EXTENDED_BASIC_INFORMATION_INFO,
		uintptr(unsafe.Pointer(&extendedInfo)),
		unsafe.Sizeof(extendedInfo),
		uintptr(unsafe.Pointer(&returnLength)),
	)

	if ret == 0 {
		if extendedInfo.IsWow64Process != 0 {
			return "x86", nil // WOW64进程，32位
		}
	}

	// 根据系统架构判断
	if runtime.GOARCH == "amd64" {
		return "x64", nil
	} else if runtime.GOARCH == "386" {
		return "x86", nil
	} else if runtime.GOARCH == "arm64" {
		return "ARM64", nil
	}

	return "unknown", nil
}

// GetProcessMemoryInfoEx 获取进程详细内存信息
func (e *ProcessInfoEnhancer) GetProcessMemoryInfoEx(pid int32) (*PROCESS_MEMORY_COUNTERS_EX, error) {
	if runtime.GOOS != "windows" {
		return nil, fmt.Errorf("only supported on Windows")
	}

	handle, err := windows.OpenProcess(PROCESS_QUERY_INFORMATION|PROCESS_VM_READ, false, uint32(pid))
	if err != nil {
		return nil, fmt.Errorf("无法打开进程句柄: %w", err)
	}
	defer windows.CloseHandle(handle)

	var memCounters PROCESS_MEMORY_COUNTERS_EX
	var cb uint32

	ret, _, err := getProcessMemoryInfo.Call(
		uintptr(handle),
		uintptr(unsafe.Pointer(&memCounters)),
		unsafe.Sizeof(memCounters),
		uintptr(unsafe.Pointer(&cb)),
	)

	if ret == 0 {
		return nil, fmt.Errorf("获取进程内存信息失败: %w", err)
	}

	return &memCounters, nil
}

// GetProcessIOCounters 获取进程IO计数器
func (e *ProcessInfoEnhancer) GetProcessIOCounters(pid int32) (*IO_COUNTERS, error) {
	if runtime.GOOS != "windows" {
		return nil, fmt.Errorf("only supported on Windows")
	}

	handle, err := windows.OpenProcess(PROCESS_QUERY_INFORMATION, false, uint32(pid))
	if err != nil {
		return nil, fmt.Errorf("无法打开进程句柄: %w", err)
	}
	defer windows.CloseHandle(handle)

	var ioCounters IO_COUNTERS

	ret, _, err := getProcessIoCounters.Call(
		uintptr(handle),
		uintptr(unsafe.Pointer(&ioCounters)),
	)

	if ret == 0 {
		return nil, fmt.Errorf("获取进程IO计数器失败: %w", err)
	}

	return &ioCounters, nil
}

// GetProcessIntegrityLevel 获取进程完整性级别
func (e *ProcessInfoEnhancer) GetProcessIntegrityLevel(pid int32) (string, error) {
	if runtime.GOOS != "windows" {
		return "unknown", fmt.Errorf("only supported on Windows")
	}

	// 这里需要实现获取进程令牌的完整性级别
	// 由于涉及复杂的Windows安全API，暂时返回默认值
	return "Medium", nil
}

// IsProcessElevated 检查进程是否提升权限
func (e *ProcessInfoEnhancer) IsProcessElevated(pid int32) (bool, error) {
	if runtime.GOOS != "windows" {
		return false, fmt.Errorf("only supported on Windows")
	}

	// 这里需要实现检查进程是否提升权限的逻辑
	// 由于涉及复杂的Windows安全API，暂时返回默认值
	return false, nil
}

// GetProcessSessionID 获取进程会话ID
func (e *ProcessInfoEnhancer) GetProcessSessionID(pid int32) (uint32, error) {
	if runtime.GOOS != "windows" {
		return 0, fmt.Errorf("only supported on Windows")
	}

	// 这里需要实现获取进程会话ID的逻辑
	// 由于涉及复杂的Windows API，暂时返回默认值
	return 0, nil
}

// GetProcessHandleCount 获取进程句柄数量
func (e *ProcessInfoEnhancer) GetProcessHandleCount(pid int32) (uint32, error) {
	if runtime.GOOS != "windows" {
		return 0, fmt.Errorf("only supported on Windows")
	}

	// 这里需要实现获取进程句柄数量的逻辑
	// 由于涉及复杂的Windows API，暂时返回默认值
	return 0, nil
}

// GetProcessContextSwitches 获取进程上下文切换数
func (e *ProcessInfoEnhancer) GetProcessContextSwitches(pid int32) (uint64, error) {
	if runtime.GOOS != "windows" {
		return 0, fmt.Errorf("only supported on Windows")
	}

	// 这里需要实现获取进程上下文切换数的逻辑
	// 由于涉及复杂的Windows API，暂时返回默认值
	return 0, nil
}

// GetProcessPageFaults 获取进程页面错误数
func (e *ProcessInfoEnhancer) GetProcessPageFaults(pid int32) (uint64, error) {
	if runtime.GOOS != "windows" {
		return 0, fmt.Errorf("only supported on Windows")
	}

	memInfo, err := e.GetProcessMemoryInfoEx(pid)
	if err != nil {
		return 0, err
	}

	return uint64(memInfo.PageFaultCount), nil
}

// CheckProcessSecurityFeatures 检查进程安全特性
func (e *ProcessInfoEnhancer) CheckProcessSecurityFeatures(pid int32) (map[string]bool, error) {
	if runtime.GOOS != "windows" {
		return nil, fmt.Errorf("only supported on Windows")
	}

	// 这里需要实现检查进程安全特性的逻辑
	// 包括DEP、ASLR等
	features := map[string]bool{
		"dep_enabled":  true, // 数据执行保护
		"aslr_enabled": true, // 地址空间布局随机化
	}

	return features, nil
}
