package security

import (
	"fmt"
	"os"
	"runtime"
	"time"
	"unsafe"

	"golang.org/x/sys/windows"
)

// Windows API constants
const (
	FILE_ATTRIBUTE_NORMAL      = 0x00000080
	FILE_SHARE_READ            = 0x00000001
	FILE_SHARE_WRITE           = 0x00000002
	FILE_SHARE_DELETE          = 0x00000004
	OPEN_EXISTING              = 3
	GENERIC_READ               = 0x80000000
	FILE_FLAG_BACKUP_SEMANTICS = 0x02000000
	FILE_FLAG_POSIX_SEMANTICS  = 0x01000000
)

// Windows API structures
type FILETIME struct {
	DwLowDateTime  uint32
	DwHighDateTime uint32
}

type BY_HANDLE_FILE_INFORMATION struct {
	DwFileAttributes     uint32
	FtCreationTime       FILETIME
	FtLastAccessTime     FILETIME
	FtLastWriteTime      FILETIME
	DwVolumeSerialNumber uint32
	NFileSizeHigh        uint32
	NFileSizeLow         uint32
	NNumberOfLinks       uint32
	NFileIndexHigh       uint32
	NFileIndexLow        uint32
}

// Lazy DLL loading
var (
	kernel32 = windows.NewLazySystemDLL("kernel32.dll")
)

// LazyProc declarations
var (
	createFileW                = kernel32.NewProc("CreateFileW")
	getFileInformationByHandle = kernel32.NewProc("GetFileInformationByHandle")
	closeHandle                = kernel32.NewProc("CloseHandle")
	setFileTime                = kernel32.NewProc("SetFileTime")
)

// FileTimeInfo 文件时间信息
type FileTimeInfo struct {
	CreationTime   time.Time
	LastAccessTime time.Time
	LastWriteTime  time.Time
}

// GetFileTimeInfo 获取文件的真实时间信息
func GetFileTimeInfo(filePath string) (*FileTimeInfo, error) {
	if runtime.GOOS != "windows" {
		// 在非Windows系统上，使用标准方法
		return getFileTimeInfoStandard(filePath)
	}

	// 在Windows上使用Windows API
	return getFileTimeInfoWindows(filePath)
}

// getFileTimeInfoWindows 使用Windows API获取文件时间信息
func getFileTimeInfoWindows(filePath string) (*FileTimeInfo, error) {
	// 转换文件路径为UTF16
	pathPtr, err := windows.UTF16PtrFromString(filePath)
	if err != nil {
		return nil, fmt.Errorf("转换文件路径失败: %w", err)
	}

	// 打开文件句柄，使用更高的权限和标志
	handle, _, err := createFileW.Call(
		uintptr(unsafe.Pointer(pathPtr)),
		uintptr(GENERIC_READ),
		uintptr(FILE_SHARE_READ|FILE_SHARE_WRITE|FILE_SHARE_DELETE),
		0,
		uintptr(OPEN_EXISTING),
		uintptr(FILE_ATTRIBUTE_NORMAL|FILE_FLAG_BACKUP_SEMANTICS),
		0,
	)

	if handle == uintptr(windows.InvalidHandle) {
		return nil, fmt.Errorf("无法打开文件: %w", err)
	}
	defer closeHandle.Call(handle)

	// 获取文件信息
	var fileInfo BY_HANDLE_FILE_INFORMATION
	ret, _, err := getFileInformationByHandle.Call(
		handle,
		uintptr(unsafe.Pointer(&fileInfo)),
	)

	if ret == 0 {
		return nil, fmt.Errorf("获取文件信息失败: %w", err)
	}

	// 转换时间，保持纳秒精度
	creationTime := fileTimeToTimePrecise(fileInfo.FtCreationTime)
	lastAccessTime := fileTimeToTimePrecise(fileInfo.FtLastAccessTime)
	lastWriteTime := fileTimeToTimePrecise(fileInfo.FtLastWriteTime)

	return &FileTimeInfo{
		CreationTime:   creationTime,
		LastAccessTime: lastAccessTime,
		LastWriteTime:  lastWriteTime,
	}, nil
}

// getFileTimeInfoStandard 使用标准方法获取文件时间信息
func getFileTimeInfoStandard(filePath string) (*FileTimeInfo, error) {
	fileInfo, err := os.Stat(filePath)
	if err != nil {
		return nil, fmt.Errorf("获取文件信息失败: %w", err)
	}

	// 在Windows上，标准方法通常只能获取修改时间
	// 创建时间和访问时间通常不可用
	return &FileTimeInfo{
		CreationTime:   fileInfo.ModTime(), // 在大多数系统上，创建时间不可用
		LastAccessTime: fileInfo.ModTime(), // 在大多数系统上，访问时间不可用
		LastWriteTime:  fileInfo.ModTime(),
	}, nil
}

// fileTimeToTime 将Windows FILETIME转换为Go time.Time（保持向后兼容）
func fileTimeToTime(ft FILETIME) time.Time {
	return fileTimeToTimePrecise(ft)
}

// fileTimeToTimePrecise 将Windows FILETIME转换为Go time.Time，保持纳秒精度
func fileTimeToTimePrecise(ft FILETIME) time.Time {
	if ft.DwHighDateTime == 0 && ft.DwLowDateTime == 0 {
		return time.Time{}
	}

	// Windows时间是从1601年1月1日开始的100纳秒间隔
	// 转换为Unix时间戳，保持纳秒精度
	windowsTime := int64(ft.DwHighDateTime)<<32 | int64(ft.DwLowDateTime)
	unixTime := (windowsTime - 116444736000000000) / 10000000
	nanoseconds := (windowsTime - 116444736000000000) % 10000000 * 100

	return time.Unix(unixTime, nanoseconds)
}

// GetFileCreationTime 获取文件创建时间
func GetFileCreationTime(filePath string) (time.Time, error) {
	timeInfo, err := GetFileTimeInfo(filePath)
	if err != nil {
		return time.Time{}, err
	}
	return timeInfo.CreationTime, nil
}

// GetFileAccessTime 获取文件访问时间
func GetFileAccessTime(filePath string) (time.Time, error) {
	timeInfo, err := GetFileTimeInfo(filePath)
	if err != nil {
		return time.Time{}, err
	}
	return timeInfo.LastAccessTime, nil
}

// GetFileAccessTimeWithUpdate 获取文件访问时间，并强制更新访问时间
func GetFileAccessTimeWithUpdate(filePath string) (time.Time, error) {
	if runtime.GOOS != "windows" {
		return GetFileAccessTime(filePath)
	}

	// 在Windows上，先尝试获取当前访问时间
	currentTime, err := GetFileAccessTime(filePath)
	if err != nil {
		return time.Time{}, err
	}

	// 强制更新访问时间到当前时间
	if err := updateFileAccessTime(filePath); err != nil {
		// 即使更新失败，也返回当前获取到的时间
		return currentTime, nil
	}

	// 重新获取更新后的访问时间
	updatedTime, err := GetFileAccessTime(filePath)
	if err != nil {
		return currentTime, nil
	}

	return updatedTime, nil
}

// updateFileAccessTime 更新文件的访问时间到当前时间
func updateFileAccessTime(filePath string) error {
	// 转换文件路径为UTF16
	pathPtr, err := windows.UTF16PtrFromString(filePath)
	if err != nil {
		return fmt.Errorf("转换文件路径失败: %w", err)
	}

	// 打开文件句柄，需要写入权限来更新时间
	handle, _, err := createFileW.Call(
		uintptr(unsafe.Pointer(pathPtr)),
		uintptr(GENERIC_READ|0x40000000), // GENERIC_READ | GENERIC_WRITE
		uintptr(FILE_SHARE_READ|FILE_SHARE_WRITE|FILE_SHARE_DELETE),
		0,
		uintptr(OPEN_EXISTING),
		uintptr(FILE_ATTRIBUTE_NORMAL),
		0,
	)

	if handle == uintptr(windows.InvalidHandle) {
		return fmt.Errorf("无法打开文件进行时间更新: %w", err)
	}
	defer closeHandle.Call(handle)

	// 获取当前系统时间
	now := time.Now()

	// 转换为Windows FILETIME格式
	fileTime := timeToFileTime(now)

	// 更新访问时间（保持创建时间和修改时间不变）
	ret, _, err := setFileTime.Call(
		handle,
		0,                                  // 不修改创建时间
		uintptr(unsafe.Pointer(&fileTime)), // 更新访问时间
		0,                                  // 不修改修改时间
	)

	if ret == 0 {
		return fmt.Errorf("更新文件访问时间失败: %w", err)
	}

	return nil
}

// timeToFileTime 将Go time.Time转换为Windows FILETIME
func timeToFileTime(t time.Time) FILETIME {
	if t.IsZero() {
		return FILETIME{0, 0}
	}

	// 转换为Windows时间（从1601年1月1日开始的100纳秒间隔）
	unixTime := t.Unix()
	nanoseconds := int64(t.Nanosecond())

	// Windows时间 = (Unix时间 + 11644473600) * 10000000 + 纳秒/100
	windowsTime := (unixTime+11644473600)*10000000 + nanoseconds/100

	return FILETIME{
		DwLowDateTime:  uint32(windowsTime & 0xFFFFFFFF),
		DwHighDateTime: uint32(windowsTime >> 32),
	}
}

// GetFileModTime 获取文件修改时间
func GetFileModTime(filePath string) (time.Time, error) {
	timeInfo, err := GetFileTimeInfo(filePath)
	if err != nil {
		return time.Time{}, err
	}
	return timeInfo.LastWriteTime, nil
}

// IsFileTimeAvailable 检查文件时间是否可用
func IsFileTimeAvailable(filePath string) (bool, error) {
	timeInfo, err := GetFileTimeInfo(filePath)
	if err != nil {
		return false, err
	}

	// 检查时间是否为零值（表示不可用）
	zeroTime := time.Time{}
	return !timeInfo.CreationTime.Equal(zeroTime) ||
		!timeInfo.LastAccessTime.Equal(zeroTime) ||
		!timeInfo.LastWriteTime.Equal(zeroTime), nil
}
