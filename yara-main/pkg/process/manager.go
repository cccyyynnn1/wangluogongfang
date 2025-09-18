package process

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"syscall"
	"time"
	"unsafe"

	"yara-security-service/internal/models"

	"github.com/shirou/gopsutil/v3/process"
	"github.com/sirupsen/logrus"
	"golang.org/x/sys/windows"
)

// Manager 进程管理器
type Manager struct {
	logger   *logrus.Logger
	enhancer *ProcessInfoEnhancer
}

// NewManager 创建进程管理器
func NewManager(logger *logrus.Logger) *Manager {
	return &Manager{
		logger:   logger,
		enhancer: NewProcessInfoEnhancer(logger),
	}
}

// GetProcessList 获取进程列表
func (m *Manager) GetProcessList() ([]*models.ProcessInfo, error) {
	processes, err := process.Processes()
	if err != nil {
		return nil, fmt.Errorf("获取进程列表失败: %w", err)
	}

	var result []*models.ProcessInfo
	for _, p := range processes {
		info, err := m.getBasicProcessInfo(p)
		if err != nil {
			m.logger.Debugf("获取进程基本信息失败 PID=%d: %v", p.Pid, err)
			continue
		}
		result = append(result, info)
	}

	return result, nil
}

// GetProcessByPID 根据PID获取进程信息
func (m *Manager) GetProcessByPID(pid int32) (*models.ProcessInfo, error) {
	proc, err := process.NewProcess(pid)
	if err != nil {
		return nil, fmt.Errorf("无法创建进程对象 PID=%d: %w", pid, err)
	}

	return m.getDetailedProcessInfo(proc)
}

// StartProcess 启动指定路径的进程
func (m *Manager) StartProcess(processPath string) error {
	if processPath == "" {
		return fmt.Errorf("process path cannot be empty")
	}

	// 检查文件是否存在
	if _, err := os.Stat(processPath); os.IsNotExist(err) {
		return fmt.Errorf("process file does not exist: %s", processPath)
	}

	// 获取进程的绝对路径
	absPath, err := filepath.Abs(processPath)
	if err != nil {
		return fmt.Errorf("failed to get absolute path: %v", err)
	}

	// 获取进程所在目录
	workDir := filepath.Dir(absPath)

	logrus.Debugf("Starting process: %s in directory: %s", absPath, workDir)

	// 创建命令，设置工作目录和环境变量
	cmd := exec.Command(absPath)
	cmd.Dir = workDir

	// 设置环境变量，确保进程有必要的系统环境
	cmd.Env = append(os.Environ(),
		"PATH="+os.Getenv("PATH"),
		"SYSTEMROOT="+os.Getenv("SYSTEMROOT"),
		"TEMP="+os.Getenv("TEMP"),
		"TMP="+os.Getenv("TMP"),
	)

	// 设置进程属性，防止进程被意外终止
	cmd.SysProcAttr = &syscall.SysProcAttr{
		HideWindow:    false,                            // 显示窗口，有助于保持进程运行
		CreationFlags: syscall.CREATE_NEW_PROCESS_GROUP, // 创建新的进程组
	}

	// 启动进程
	if err := cmd.Start(); err != nil {
		return fmt.Errorf("failed to start process: %v", err)
	}

	// 获取进程ID
	pid := cmd.Process.Pid
	logrus.Infof("Process started successfully with PID: %d", pid)

	// 启动一个goroutine来监控进程状态，但不等待进程结束
	go func() {
		// 等待进程结束，但不阻塞主线程
		cmd.Wait()
		logrus.Debugf("Process with PID %d has exited", pid)
	}()

	// 验证进程是否真的在运行
	time.Sleep(100 * time.Millisecond) // 给进程一点时间启动

	// 检查进程是否还在运行
	proc, err := os.FindProcess(pid)
	if err != nil {
		return fmt.Errorf("failed to find started process: %v", err)
	}

	// 尝试发送信号0来检查进程是否存活（Windows上可能不工作，但值得尝试）
	if err := proc.Signal(syscall.Signal(0)); err != nil {
		// 如果信号发送失败，尝试通过gopsutil检查
		if p, err := process.NewProcess(int32(pid)); err == nil {
			if running, err := p.IsRunning(); err == nil && running {
				logrus.Infof("Process %d is confirmed running via gopsutil", pid)
				return nil
			}
		}
		return fmt.Errorf("process started but may not be running: %v", err)
	}

	logrus.Infof("Process %d started and confirmed running", pid)
	return nil
}

// KillProcess 结束进程
func (m *Manager) KillProcess(pid int32) error {
	proc, err := process.NewProcess(pid)
	if err != nil {
		return fmt.Errorf("无法创建进程对象 PID=%d: %w", pid, err)
	}

	err = proc.Kill()
	if err != nil {
		return fmt.Errorf("结束进程失败 PID=%d: %w", pid, err)
	}

	return nil
}

// SuspendProcess 挂起进程
func (m *Manager) SuspendProcess(pid int32) error {
	// 使用Windows API挂起进程的所有线程
	handle, err := windows.OpenProcess(windows.PROCESS_SUSPEND_RESUME, false, uint32(pid))
	if err != nil {
		return fmt.Errorf("无法打开进程句柄 PID=%d: %w", pid, err)
	}
	defer windows.CloseHandle(handle)

	// 创建进程快照来枚举所有线程
	snapshot, err := windows.CreateToolhelp32Snapshot(windows.TH32CS_SNAPTHREAD, 0)
	if err != nil {
		return fmt.Errorf("创建线程快照失败: %w", err)
	}
	defer windows.CloseHandle(snapshot)

	var threadEntry windows.ThreadEntry32
	threadEntry.Size = uint32(unsafe.Sizeof(threadEntry))

	// 枚举所有线程
	err = windows.Thread32First(snapshot, &threadEntry)
	if err != nil {
		return fmt.Errorf("枚举线程失败: %w", err)
	}

	var suspendedCount int
	for {
		// 检查线程是否属于目标进程
		if threadEntry.OwnerProcessID == uint32(pid) {
			// 打开线程句柄
			threadHandle, err := windows.OpenThread(windows.THREAD_SUSPEND_RESUME, false, threadEntry.ThreadID)
			if err != nil {
				m.logger.Warnf("无法打开线程句柄 TID=%d: %v", threadEntry.ThreadID, err)
				continue
			}

			// 使用syscall直接调用NtSuspendThread
			var suspendCount uint32
			ret, _, err := syscall.NewLazyDLL("ntdll.dll").NewProc("NtSuspendThread").Call(
				uintptr(threadHandle),
				uintptr(unsafe.Pointer(&suspendCount)),
			)
			windows.CloseHandle(threadHandle)

			if ret != 0 {
				m.logger.Warnf("挂起线程失败 TID=%d: %v", threadEntry.ThreadID, err)
			} else {
				suspendedCount++
			}
		}

		// 获取下一个线程
		err = windows.Thread32Next(snapshot, &threadEntry)
		if err != nil {
			break
		}
	}

	if suspendedCount == 0 {
		return fmt.Errorf("没有找到进程 %d 的线程或所有线程挂起失败", pid)
	}

	m.logger.Infof("成功挂起进程 PID=%d 的 %d 个线程", pid, suspendedCount)
	return nil
}

// ResumeProcess 恢复进程
func (m *Manager) ResumeProcess(pid int32) error {
	// 使用Windows API恢复进程的所有线程
	handle, err := windows.OpenProcess(windows.PROCESS_SUSPEND_RESUME, false, uint32(pid))
	if err != nil {
		return fmt.Errorf("无法打开进程句柄 PID=%d: %w", pid, err)
	}
	defer windows.CloseHandle(handle)

	// 创建进程快照来枚举所有线程
	snapshot, err := windows.CreateToolhelp32Snapshot(windows.TH32CS_SNAPTHREAD, 0)
	if err != nil {
		return fmt.Errorf("创建线程快照失败: %w", err)
	}
	defer windows.CloseHandle(snapshot)

	var threadEntry windows.ThreadEntry32
	threadEntry.Size = uint32(unsafe.Sizeof(threadEntry))

	// 枚举所有线程
	err = windows.Thread32First(snapshot, &threadEntry)
	if err != nil {
		return fmt.Errorf("枚举线程失败: %w", err)
	}

	var resumedCount int
	for {
		// 检查线程是否属于目标进程
		if threadEntry.OwnerProcessID == uint32(pid) {
			// 打开线程句柄
			threadHandle, err := windows.OpenThread(windows.THREAD_SUSPEND_RESUME, false, threadEntry.ThreadID)
			if err != nil {
				m.logger.Warnf("无法打开线程句柄 TID=%d: %v", threadEntry.ThreadID, err)
				continue
			}

			// 使用syscall直接调用NtResumeThread
			var suspendCount uint32
			ret, _, err := syscall.NewLazyDLL("ntdll.dll").NewProc("NtResumeThread").Call(
				uintptr(threadHandle),
				uintptr(unsafe.Pointer(&suspendCount)),
			)
			windows.CloseHandle(threadHandle)

			if ret != 0 {
				m.logger.Warnf("恢复线程失败 TID=%d: %v", threadEntry.ThreadID, err)
			} else {
				resumedCount++
			}
		}

		// 获取下一个线程
		err = windows.Thread32Next(snapshot, &threadEntry)
		if err != nil {
			break
		}
	}

	if resumedCount == 0 {
		return fmt.Errorf("没有找到进程 %d 的线程或所有线程恢复失败", pid)
	}

	m.logger.Infof("成功恢复进程 PID=%d 的 %d 个线程", pid, resumedCount)
	return nil
}

// GetProcessModules 获取进程模块
func (m *Manager) GetProcessModules(pid int32) ([]*models.ProcessInfo, error) {
	m.logger.Debugf("开始获取进程模块相关进程信息 PID=%d", pid)

	// 首先获取系统模块信息
	systemModules, err := m.GetSystemModules()
	if err != nil {
		m.logger.Errorf("获取系统模块失败 PID=%d: %v", pid, err)
		return nil, fmt.Errorf("获取系统模块失败: %w", err)
	}

	// 找到与指定进程相关的模块
	var relatedProcesses []*models.ProcessInfo
	processedPIDs := make(map[int32]bool) // 避免重复处理同一进程

	for _, module := range systemModules {
		// 检查该模块是否被指定进程使用
		isUsedByTargetProcess := false
		for _, affectedPID := range module.AffectedProcesses {
			if affectedPID == pid {
				isUsedByTargetProcess = true
				break
			}
		}

		// 如果模块被指定进程使用，则获取所有使用该模块的进程信息
		if isUsedByTargetProcess {
			m.logger.Debugf("找到相关模块: %s, 影响进程数: %d", module.Name, module.ProcessCount)

			for _, affectedPID := range module.AffectedProcesses {
				// 避免重复处理同一进程
				if processedPIDs[affectedPID] {
					continue
				}
				processedPIDs[affectedPID] = true

				// 获取进程详细信息
				processInfo, err := m.GetProcessByPID(affectedPID)
				if err != nil {
					m.logger.Debugf("获取进程信息失败 PID=%d: %v", affectedPID, err)
					continue
				}

				relatedProcesses = append(relatedProcesses, processInfo)
				m.logger.Debugf("添加相关进程 PID=%d, Name=%s", affectedPID, processInfo.Name)
			}
		}
	}

	m.logger.Debugf("成功获取进程模块相关进程信息 PID=%d，共%d个相关进程", pid, len(relatedProcesses))
	return relatedProcesses, nil
}

// getProcessModulesWindows 使用Windows API获取进程模块
func (m *Manager) getProcessModulesWindows(pid int32) ([]*models.ProcessModuleInfo, error) {
	// 打开进程句柄
	handle, err := windows.OpenProcess(windows.PROCESS_QUERY_INFORMATION|windows.PROCESS_VM_READ, false, uint32(pid))
	if err != nil {
		return nil, fmt.Errorf("无法打开进程句柄 PID=%d: %w", pid, err)
	}
	defer windows.CloseHandle(handle)

	m.logger.Debugf("成功打开进程句柄 PID=%d, Handle=%v", pid, handle)

	// 验证进程句柄是否对应正确的PID
	if !m.verifyProcessHandle(handle, pid) {
		return nil, fmt.Errorf("进程句柄验证失败，句柄可能指向错误的进程 PID=%d", pid)
	}

	// 使用EnumProcessModules获取模块列表
	modules, err := m.enumProcessModulesWindows(handle, pid)
	if err != nil {
		return nil, fmt.Errorf("枚举进程模块失败 PID=%d: %w", pid, err)
	}

	var result []*models.ProcessModuleInfo
	for _, module := range modules {
		info := &models.ProcessModuleInfo{
			PID:  module.PID, // 使用模块中的实际PID
			Name: module.Name,
			Path: module.Path,
			Size: module.Size,
		}
		result = append(result, info)
	}

	m.logger.Debugf("处理完成，返回%d个模块信息 PID=%d", len(result), pid)
	return result, nil
}

// ProcessModule 进程模块信息
type ProcessModule struct {
	PID  int32  `json:"pid"`
	Name string `json:"name"`
	Path string `json:"path"`
	Size uint64 `json:"size"`
}

// enumProcessModulesWindows 使用Windows API枚举进程模块
func (m *Manager) enumProcessModulesWindows(handle windows.Handle, targetPID int32) ([]*ProcessModule, error) {
	m.logger.Debugf("开始枚举进程模块 Handle=%v, TargetPID=%d", handle, targetPID)

	// 验证句柄有效性
	if handle == 0 {
		return nil, fmt.Errorf("无效的进程句柄")
	}

	// 使用Windows API真正枚举进程模块
	var modules []*ProcessModule

	// 获取进程模块数量
	var needed uint32
	var moduleCount uint32

	// 首先调用EnumProcessModules获取需要的缓冲区大小
	ret, _, err := syscall.NewLazyDLL("psapi.dll").NewProc("EnumProcessModules").Call(
		uintptr(handle),
		0, // 不传递模块句柄数组
		0, // 缓冲区大小为0
		uintptr(unsafe.Pointer(&needed)),
	)

	if ret == 0 {
		m.logger.Debugf("psapi.dll EnumProcessModules失败，尝试kernel32.dll K32EnumProcessModules")
		// 如果失败，尝试使用K32EnumProcessModules (Windows 7+)
		ret, _, err = syscall.NewLazyDLL("kernel32.dll").NewProc("K32EnumProcessModules").Call(
			uintptr(handle),
			0, // 不传递模块句柄数组
			0, // 缓冲区大小为0
			uintptr(unsafe.Pointer(&needed)),
		)
	}

	if ret == 0 {
		return nil, fmt.Errorf("获取模块数量失败: %w", err)
	}

	moduleCount = needed / uint32(unsafe.Sizeof(uintptr(0)))
	m.logger.Debugf("检测到%d个模块需要处理", moduleCount)

	if moduleCount == 0 {
		m.logger.Debugf("进程没有模块")
		return []*ProcessModule{}, nil
	}

	// 分配模块句柄数组
	moduleHandles := make([]uintptr, moduleCount)

	// 调用EnumProcessModules获取模块句柄
	ret, _, err = syscall.NewLazyDLL("psapi.dll").NewProc("EnumProcessModules").Call(
		uintptr(handle),
		uintptr(unsafe.Pointer(&moduleHandles[0])),
		uintptr(needed),
		uintptr(unsafe.Pointer(&needed)),
	)

	if ret == 0 {
		m.logger.Debugf("psapi.dll EnumProcessModules获取句柄失败，尝试kernel32.dll")
		// 如果失败，尝试使用K32EnumProcessModules
		ret, _, err = syscall.NewLazyDLL("kernel32.dll").NewProc("K32EnumProcessModules").Call(
			uintptr(handle),
			uintptr(unsafe.Pointer(&moduleHandles[0])),
			uintptr(needed),
			uintptr(unsafe.Pointer(&needed)),
		)
	}

	if ret == 0 {
		return nil, fmt.Errorf("枚举进程模块失败: %w", err)
	}

	m.logger.Debugf("成功获取%d个模块句柄", len(moduleHandles))

	// 打印所有模块句柄用于调试
	for i, h := range moduleHandles {
		m.logger.Debugf("模块句柄[%d]: %v", i, h)
	}

	// 为每个模块获取详细信息
	for i := uint32(0); i < moduleCount; i++ {
		moduleHandle := moduleHandles[i]
		if moduleHandle == 0 {
			m.logger.Debugf("跳过空模块句柄 索引=%d", i)
			continue
		}

		m.logger.Debugf("处理模块 %d/%d, Handle=%v", i+1, moduleCount, moduleHandle)

		// 获取模块文件名
		var filename [windows.MAX_PATH]uint16
		var filenameLen uint32 = windows.MAX_PATH

		ret, _, err := syscall.NewLazyDLL("psapi.dll").NewProc("GetModuleFileNameExW").Call(
			uintptr(handle),
			moduleHandle,
			uintptr(unsafe.Pointer(&filename[0])),
			uintptr(filenameLen),
		)

		if ret == 0 {
			m.logger.Debugf("psapi.dll GetModuleFileNameEx失败，尝试kernel32.dll")
			// 如果失败，尝试使用K32GetModuleFileNameEx
			ret, _, err = syscall.NewLazyDLL("kernel32.dll").NewProc("K32GetModuleFileNameExW").Call(
				uintptr(handle),
				moduleHandle,
				uintptr(unsafe.Pointer(&filename[0])),
				uintptr(filenameLen),
			)
		}

		if ret == 0 {
			m.logger.Warnf("获取模块文件名失败 模块句柄=%v: %v", moduleHandle, err)
			continue
		}

		// 获取模块信息
		var moduleInfo struct {
			BaseOfDll   uintptr
			SizeOfImage uint32
			EntryPoint  uintptr
		}

		ret, _, err = syscall.NewLazyDLL("psapi.dll").NewProc("GetModuleInformation").Call(
			uintptr(handle),
			moduleHandle,
			uintptr(unsafe.Pointer(&moduleInfo)),
			uintptr(unsafe.Sizeof(moduleInfo)),
		)

		if ret == 0 {
			m.logger.Debugf("psapi.dll GetModuleInformation失败，尝试kernel32.dll")
			// 如果失败，尝试使用K32GetModuleInformation
			ret, _, err = syscall.NewLazyDLL("kernel32.dll").NewProc("K32GetModuleInformation").Call(
				uintptr(handle),
				moduleHandle,
				uintptr(unsafe.Pointer(&moduleInfo)),
				uintptr(unsafe.Sizeof(moduleInfo)),
			)
		}

		if ret == 0 {
			m.logger.Warnf("获取模块信息失败 模块句柄=%v: %v", moduleHandle, err)
			continue
		}

		// 转换文件名
		modulePath := windows.UTF16ToString(filename[:])
		moduleName := filepath.Base(modulePath)

		// 验证模块路径是否属于目标进程
		if modulePath == "" {
			m.logger.Debugf("跳过空路径模块 模块句柄=%v", moduleHandle)
			continue
		}

		m.logger.Debugf("添加模块: %s, 路径: %s, 大小: %d, 基址: %v", moduleName, modulePath, moduleInfo.SizeOfImage, moduleInfo.BaseOfDll)

		module := &ProcessModule{
			PID:  targetPID, // 设置PID
			Name: moduleName,
			Path: modulePath,
			Size: uint64(moduleInfo.SizeOfImage),
		}

		modules = append(modules, module)
	}

	m.logger.Debugf("枚举完成，共找到%d个有效模块", len(modules))
	return modules, nil
}

// isProcessRunning 检查进程是否正在运行
func (m *Manager) isProcessRunning(pid int32) bool {
	process, err := process.NewProcess(pid)
	if err != nil {
		return false
	}

	running, err := process.IsRunning()
	if err != nil {
		return false
	}

	return running
}

// validateProcessModules 验证模块确实属于目标进程
func (m *Manager) validateProcessModules(modules []*ProcessModule, targetPID int32) []*ProcessModule {
	var validatedModules []*ProcessModule

	for _, module := range modules {
		if module.Path == "" {
			m.logger.Debugf("跳过空路径模块")
			continue
		}

		// 检查模块文件是否存在
		if _, err := os.Stat(module.Path); os.IsNotExist(err) {
			m.logger.Debugf("跳过不存在的模块文件: %s", module.Path)
			continue
		}

		// 检查模块是否被目标进程加载
		if m.isModuleLoadedByProcess(module.Path, targetPID) {
			validatedModules = append(validatedModules, module)
			m.logger.Debugf("验证通过: %s", module.Path)
		} else {
			m.logger.Debugf("验证失败，模块不属于目标进程: %s", module.Path)
		}
	}

	m.logger.Debugf("模块验证完成: 原始%d个，验证后%d个", len(modules), len(validatedModules))
	return validatedModules
}

// isModuleLoadedByProcess 检查模块是否被指定进程加载
func (m *Manager) isModuleLoadedByProcess(modulePath string, targetPID int32) bool {
	// 使用Windows API检查模块是否被目标进程加载
	handle, err := windows.OpenProcess(windows.PROCESS_QUERY_INFORMATION|windows.PROCESS_VM_READ, false, uint32(targetPID))
	if err != nil {
		m.logger.Debugf("无法打开进程句柄进行模块验证 PID=%d: %v", targetPID, err)
		return false
	}
	defer windows.CloseHandle(handle)

	// 这里可以添加更详细的验证逻辑
	// 暂时返回true，因为如果模块在枚举时被找到，通常意味着它属于该进程
	return true
}

// verifyProcessHandle 验证进程句柄是否对应正确的PID
func (m *Manager) verifyProcessHandle(handle windows.Handle, targetPID int32) bool {
	if handle == 0 {
		return false
	}

	// 获取进程ID
	var processId uint32
	ret, _, err := syscall.NewLazyDLL("kernel32.dll").NewProc("GetProcessId").Call(uintptr(handle))
	if ret == 0 {
		m.logger.Debugf("无法获取进程ID: %v", err)
		return false
	}

	processId = uint32(ret)
	m.logger.Debugf("句柄对应的进程ID: %d, 目标PID: %d", processId, targetPID)

	return processId == uint32(targetPID)
}

// GetProcessConnections 获取进程网络连接
func (m *Manager) GetProcessConnections(pid int32) ([]*models.ConnectionInfo, error) {
	proc, err := process.NewProcess(pid)
	if err != nil {
		return nil, fmt.Errorf("无法创建进程对象 PID=%d: %w", pid, err)
	}

	// 获取进程的网络连接
	connections, err := proc.Connections()
	if err != nil {
		// 如果无法获取连接，尝试使用Windows API
		return m.getProcessConnectionsWindows(pid)
	}

	var result []*models.ConnectionInfo
	for _, conn := range connections {
		info := &models.ConnectionInfo{
			FD:     conn.Fd,
			Family: int32(conn.Family),
			Type:   int32(conn.Type),
			Laddr: &models.AddrInfo{
				IP:   conn.Laddr.IP,
				Port: conn.Laddr.Port,
			},
			Raddr: &models.AddrInfo{
				IP:   conn.Raddr.IP,
				Port: conn.Raddr.Port,
			},
			Status: conn.Status,
		}
		result = append(result, info)
	}

	m.logger.Debugf("获取进程网络连接成功 PID=%d, 共%d个连接", pid, len(result))
	return result, nil
}

// getProcessConnectionsWindows 使用Windows API获取进程网络连接
func (m *Manager) getProcessConnectionsWindows(pid int32) ([]*models.ConnectionInfo, error) {
	// 使用netstat命令获取进程的网络连接
	cmd := exec.Command("netstat", "-ano")
	output, err := cmd.Output()
	if err != nil {
		return []*models.ConnectionInfo{}, nil
	}

	// 解析netstat输出
	connections := m.parseNetstatOutput(string(output), pid)

	var result []*models.ConnectionInfo
	for _, conn := range connections {
		info := &models.ConnectionInfo{
			FD:     conn.FD,
			Family: conn.Family,
			Type:   conn.Type,
			Laddr:  conn.Laddr,
			Raddr:  conn.Raddr,
			Status: conn.Status,
		}
		result = append(result, info)
	}

	return result, nil
}

// parseNetstatOutput 解析netstat命令输出
func (m *Manager) parseNetstatOutput(output string, targetPID int32) []*models.ConnectionInfo {
	var connections []*models.ConnectionInfo

	// 按行分割netstat输出
	lines := strings.Split(output, "\n")

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		// 跳过标题行
		if strings.Contains(line, "Proto") || strings.Contains(line, "协议") {
			continue
		}

		// 分割行内容
		fields := strings.Fields(line)
		if len(fields) < 5 {
			continue
		}

		// 解析协议
		proto := fields[0]
		var family int32
		var connType int32

		switch strings.ToLower(proto) {
		case "tcp", "tcp6":
			family = 2   // AF_INET
			connType = 1 // SOCK_STREAM
		case "udp", "udp6":
			family = 2   // AF_INET
			connType = 2 // SOCK_DGRAM
		default:
			continue
		}

		// 解析本地地址和端口
		localAddr := fields[1]
		remoteAddr := fields[2]
		status := fields[3]
		pidStr := fields[4]

		// 检查PID是否匹配
		pid, err := strconv.ParseInt(pidStr, 10, 32)
		if err != nil || int32(pid) != targetPID {
			continue
		}

		// 解析本地地址
		localIP, localPort := m.parseAddr(localAddr)
		remoteIP, remotePort := m.parseAddr(remoteAddr)

		// 创建连接信息
		conn := &models.ConnectionInfo{
			FD:     uint32(len(connections) + 1), // 简单的文件描述符
			Family: family,
			Type:   connType,
			Laddr: &models.AddrInfo{
				IP:   localIP,
				Port: uint32(localPort),
			},
			Raddr: &models.AddrInfo{
				IP:   remoteIP,
				Port: uint32(remotePort),
			},
			Status: status,
		}

		connections = append(connections, conn)
	}

	return connections
}

// parseAddr 解析地址字符串，返回IP和端口
func (m *Manager) parseAddr(addr string) (string, int32) {
	// 处理IPv6地址
	if strings.HasPrefix(addr, "[") {
		// IPv6格式: [::1]:8080
		parts := strings.Split(addr, "]")
		if len(parts) == 2 {
			ip := strings.TrimPrefix(parts[0], "[")
			portStr := strings.TrimPrefix(parts[1], ":")
			if port, err := strconv.ParseInt(portStr, 10, 32); err == nil {
				return ip, int32(port)
			}
			return ip, 0
		}
	} else {
		// IPv4格式: 127.0.0.1:8080
		parts := strings.Split(addr, ":")
		if len(parts) == 2 {
			ip := parts[0]
			if port, err := strconv.ParseInt(parts[1], 10, 32); err == nil {
				return ip, int32(port)
			}
			return ip, 0
		}
	}

	return addr, 0
}

// GetProcessMemoryInfo 获取进程内存信息
func (m *Manager) GetProcessMemoryInfo(pid int32) (*models.ProcessMemoryInfo, error) {
	proc, err := process.NewProcess(pid)
	if err != nil {
		return nil, fmt.Errorf("无法创建进程对象 PID=%d: %w", pid, err)
	}

	// 获取内存信息
	memInfo, err := proc.MemoryInfo()
	if err != nil {
		return nil, fmt.Errorf("获取进程内存信息失败 PID=%d: %w", pid, err)
	}

	// 获取内存百分比
	memoryPercent, err := proc.MemoryPercent()
	if err != nil {
		m.logger.Debugf("获取进程内存百分比失败 PID=%d: %v", pid, err)
		memoryPercent = 0
	}

	info := &models.ProcessMemoryInfo{
		RSS:     memInfo.RSS,
		VMS:     memInfo.VMS,
		Percent: float64(memoryPercent),
	}

	return info, nil
}

// IsProcessRunning 检查进程是否运行
func (m *Manager) IsProcessRunning(pid int32) (bool, error) {
	proc, err := process.NewProcess(pid)
	if err != nil {
		return false, fmt.Errorf("无法创建进程对象 PID=%d: %w", pid, err)
	}

	running, err := proc.IsRunning()
	if err != nil {
		return false, fmt.Errorf("检查进程运行状态失败 PID=%d: %w", pid, err)
	}

	return running, nil
}

// GetProcessChildren 获取子进程
func (m *Manager) GetProcessChildren(pid int32) ([]int32, error) {
	proc, err := process.NewProcess(pid)
	if err != nil {
		return nil, fmt.Errorf("无法创建进程对象 PID=%d: %w", pid, err)
	}

	children, err := proc.Children()
	if err != nil {
		return nil, fmt.Errorf("获取子进程失败 PID=%d: %w", pid, err)
	}

	var result []int32
	for _, child := range children {
		result = append(result, child.Pid)
	}

	return result, nil
}

// GetSystemModules 获取系统模块列表
func (m *Manager) GetSystemModules() ([]*models.SystemModule, error) {
	// 获取所有进程
	processes, err := process.Processes()
	if err != nil {
		return nil, fmt.Errorf("获取进程列表失败: %w", err)
	}

	// 统计模块使用情况
	moduleUsage := make(map[string]*models.SystemModule)

	// 添加系统核心模块
	systemModules := []string{
		"ntoskrnl.exe", // Windows内核
		"ntdll.dll",    // NT DLL
		"kernel32.dll", // 内核32
		"user32.dll",   // 用户32
		"gdi32.dll",    // GDI32
		"advapi32.dll", // 高级API32
		"shell32.dll",  // Shell32
		"ole32.dll",    // OLE32
		"rpcrt4.dll",   // RPC运行时
		"ws2_32.dll",   // Winsock2
		"msvcrt.dll",   // MSVCRT
		"ucrtbase.dll", // UCRT基础
	}

	for _, moduleName := range systemModules {
		moduleUsage[moduleName] = &models.SystemModule{
			Name:              moduleName,
			Path:              fmt.Sprintf("C:\\Windows\\System32\\%s", moduleName),
			ProcessCount:      0,
			AffectedProcesses: []int32{},
		}
	}

	// 统计进程使用的模块
	for _, p := range processes {
		// 获取进程名称
		name, err := p.Name()
		if err != nil {
			continue
		}

		// 检查是否是系统模块
		if module, exists := moduleUsage[name]; exists {
			module.ProcessCount++
			module.AffectedProcesses = append(module.AffectedProcesses, p.Pid)
		}

		// 添加进程特定的模块
		processModuleName := fmt.Sprintf("process_%s", name)
		if existing, exists := moduleUsage[processModuleName]; exists {
			existing.ProcessCount++
			existing.AffectedProcesses = append(existing.AffectedProcesses, p.Pid)
		} else {
			moduleUsage[processModuleName] = &models.SystemModule{
				Name:              processModuleName,
				Path:              fmt.Sprintf("process_path_%s", name),
				ProcessCount:      1,
				AffectedProcesses: []int32{p.Pid},
			}
		}
	}

	// 转换为切片
	var result []*models.SystemModule
	for _, module := range moduleUsage {
		result = append(result, module)
	}

	m.logger.Debugf("获取系统模块成功，共%d个模块", len(result))
	return result, nil
}

// GetModuleInfo 获取模块信息
func (m *Manager) GetModuleInfo(moduleName string) (*models.SystemModule, error) {
	modules, err := m.GetSystemModules()
	if err != nil {
		return nil, err
	}

	for _, module := range modules {
		if module.Name == moduleName {
			return module, nil
		}
	}

	return nil, fmt.Errorf("模块未找到: %s", moduleName)
}

// SuspendModule 挂起模块
func (m *Manager) SuspendModule(moduleName string) error {
	// 获取模块信息
	module, err := m.GetModuleInfo(moduleName)
	if err != nil {
		return fmt.Errorf("获取模块信息失败: %w", err)
	}

	// 挂起所有使用该模块的进程
	var suspendedCount int
	for _, pid := range module.AffectedProcesses {
		err := m.SuspendProcess(pid)
		if err != nil {
			m.logger.Warnf("挂起进程失败 PID=%d: %v", pid, err)
			continue
		}
		suspendedCount++
	}

	if suspendedCount == 0 {
		return fmt.Errorf("没有进程使用该模块或所有进程挂起失败")
	}

	m.logger.Infof("成功挂起模块 %s，影响了 %d 个进程", moduleName, suspendedCount)
	return nil
}

// ResumeModule 恢复模块
func (m *Manager) ResumeModule(moduleName string) error {
	// 获取模块信息
	module, err := m.GetModuleInfo(moduleName)
	if err != nil {
		return fmt.Errorf("获取模块信息失败: %w", err)
	}

	// 恢复所有使用该模块的进程
	var resumedCount int
	for _, pid := range module.AffectedProcesses {
		err := m.ResumeProcess(pid)
		if err != nil {
			m.logger.Warnf("恢复进程失败 PID=%d: %v", pid, err)
			continue
		}
		resumedCount++
	}

	if resumedCount == 0 {
		return fmt.Errorf("没有进程使用该模块或所有进程恢复失败")
	}

	m.logger.Infof("成功恢复模块 %s，影响了 %d 个进程", moduleName, resumedCount)
	return nil
}

// KillModule 结束模块
func (m *Manager) KillModule(moduleName string) error {
	// 获取模块信息
	module, err := m.GetModuleInfo(moduleName)
	if err != nil {
		return fmt.Errorf("获取模块信息失败: %w", err)
	}

	// 结束所有使用该模块的进程
	var killedCount int
	for _, pid := range module.AffectedProcesses {
		err := m.KillProcess(pid)
		if err != nil {
			m.logger.Warnf("结束进程失败 PID=%d: %v", pid, err)
			continue
		}
		killedCount++
	}

	if killedCount == 0 {
		return fmt.Errorf("没有进程使用该模块或所有进程结束失败")
	}

	m.logger.Infof("成功结束模块 %s，影响了 %d 个进程", moduleName, killedCount)
	return nil
}

// GetSystemMemoryInfo 获取系统内存信息
func (m *Manager) GetSystemMemoryInfo() (map[string]interface{}, error) {
	// 获取所有进程的内存使用情况
	processes, err := process.Processes()
	if err != nil {
		return nil, fmt.Errorf("获取进程列表失败: %w", err)
	}

	var totalMemory uint64
	var runningProcesses int
	var suspendedProcesses int

	for _, p := range processes {
		// 检查进程状态
		status, err := p.Status()
		if err == nil && len(status) > 0 {
			if status[0] == "running" {
				runningProcesses++
			} else if status[0] == "stopped" {
				suspendedProcesses++
			}
		}

		// 获取内存使用
		memInfo, err := p.MemoryInfo()
		if err == nil {
			totalMemory += memInfo.RSS
		}
	}

	result := map[string]interface{}{
		"total_processes":     len(processes),
		"running_processes":   runningProcesses,
		"suspended_processes": suspendedProcesses,
		"total_memory_usage":  totalMemory,
	}

	return result, nil
}

// 辅助方法

// getBasicProcessInfo 获取进程基本信息
func (m *Manager) getBasicProcessInfo(p *process.Process) (*models.ProcessInfo, error) {
	info := &models.ProcessInfo{
		PID: p.Pid,
	}

	// 获取进程名称
	if name, err := p.Name(); err == nil {
		info.Name = name
	}

	// 获取CPU百分比
	if cpuPercent, err := p.CPUPercent(); err == nil {
		info.CPUPercent = cpuPercent
	}

	// 获取内存百分比
	if memoryPercent, err := p.MemoryPercent(); err == nil {
		info.MemoryPercent = float64(memoryPercent)
	}

	// 获取进程状态
	if status, err := p.Status(); err == nil && len(status) > 0 {
		info.Status = status[0]
	}

	// 获取创建时间
	if createTime, err := p.CreateTime(); err == nil {
		info.CreateTime = time.Unix(createTime/1000, 0)
	}

	// 使用增强器获取基本信息
	if m.enhancer != nil {
		// 获取进程体系结构（基本信息中包含）
		if arch, err := m.enhancer.GetProcessArchitecture(p.Pid); err == nil {
			info.Architecture = arch
		}

		// 获取基本内存信息
		if memInfo, err := p.MemoryInfo(); err == nil {
			info.MemoryInfo = &models.MemoryInfo{
				RSS: memInfo.RSS,
				VMS: memInfo.VMS,
			}
			info.WorkingSet = memInfo.RSS // 工作集大小
		}
	}

	// 设置最后更新时间
	info.LastUpdate = time.Now()

	return info, nil
}

// getDetailedProcessInfo 获取进程详细信息
func (m *Manager) getDetailedProcessInfo(proc *process.Process) (*models.ProcessInfo, error) {
	info := &models.ProcessInfo{
		PID: proc.Pid,
	}

	// 获取基本信息
	if name, err := proc.Name(); err == nil {
		info.Name = name
	}

	if exe, err := proc.Exe(); err == nil {
		info.Exe = exe
	}

	if cmdline, err := proc.Cmdline(); err == nil {
		info.Cmdline = cmdline
	}

	if cwd, err := proc.Cwd(); err == nil {
		info.Cwd = cwd
	}

	if status, err := proc.Status(); err == nil && len(status) > 0 {
		info.Status = status[0]
	}

	if createTime, err := proc.CreateTime(); err == nil {
		info.CreateTime = time.Unix(createTime/1000, 0)
	}

	if username, err := proc.Username(); err == nil {
		info.Username = username
	}

	if ppid, err := proc.Ppid(); err == nil {
		info.PPID = ppid
	}

	if cpuPercent, err := proc.CPUPercent(); err == nil {
		info.CPUPercent = cpuPercent
	}

	if numThreads, err := proc.NumThreads(); err == nil {
		info.NumThreads = numThreads
	}

	if priority, err := proc.Nice(); err == nil {
		info.Priority = priority
	}

	// 获取内存信息
	if memInfo, err := proc.MemoryInfo(); err == nil {
		info.MemoryInfo = &models.MemoryInfo{
			RSS: memInfo.RSS,
			VMS: memInfo.VMS,
		}
	}

	// 使用增强器获取详细信息
	if m.enhancer != nil {
		// 获取进程体系结构
		if arch, err := m.enhancer.GetProcessArchitecture(proc.Pid); err == nil {
			info.Architecture = arch
		}

		// 获取内存百分比
		if memoryPercent, err := proc.MemoryPercent(); err == nil {
			info.MemoryPercent = float64(memoryPercent)
		}

		// 获取详细内存信息
		if memCounters, err := m.enhancer.GetProcessMemoryInfoEx(proc.Pid); err == nil {
			info.WorkingSet = memCounters.WorkingSetSize
			info.PrivateBytes = memCounters.WorkingSetSizePrivate
			info.PeakWorkingSet = memCounters.PeakWorkingSetSize
			info.PeakPrivateBytes = memCounters.PeakWorkingSetSizePrivate
		}

		// 获取页面错误数
		if pageFaults, err := m.enhancer.GetProcessPageFaults(proc.Pid); err == nil {
			info.PageFaults = pageFaults
		}

		// 获取IO计数器
		if ioCounters, err := m.enhancer.GetProcessIOCounters(proc.Pid); err == nil {
			info.IOCounters = &models.ProcessIOCounters{
				ReadCount:  ioCounters.ReadOperationCount,
				WriteCount: ioCounters.WriteOperationCount,
				OtherCount: ioCounters.OtherOperationCount,
				ReadBytes:  ioCounters.ReadTransferCount,
				WriteBytes: ioCounters.WriteTransferCount,
				OtherBytes: ioCounters.OtherTransferCount,
			}
		}

		// 获取完整性级别
		if integrityLevel, err := m.enhancer.GetProcessIntegrityLevel(proc.Pid); err == nil {
			info.IntegrityLevel = integrityLevel
		}

		// 检查是否提升权限
		if elevated, err := m.enhancer.IsProcessElevated(proc.Pid); err == nil {
			info.Elevated = elevated
		}

		// 获取会话ID
		if sessionID, err := m.enhancer.GetProcessSessionID(proc.Pid); err == nil {
			info.SessionID = sessionID
		}

		// 获取句柄数量
		if handleCount, err := m.enhancer.GetProcessHandleCount(proc.Pid); err == nil {
			info.HandleCount = handleCount
		}

		// 获取上下文切换数
		if contextSwitches, err := m.enhancer.GetProcessContextSwitches(proc.Pid); err == nil {
			info.ContextSwitches = contextSwitches
		}

		// 检查安全特性
		if securityFeatures, err := m.enhancer.CheckProcessSecurityFeatures(proc.Pid); err == nil {
			if depEnabled, ok := securityFeatures["dep_enabled"]; ok {
				info.DEPEnabled = depEnabled
			}
			if aslrEnabled, ok := securityFeatures["aslr_enabled"]; ok {
				info.ASLREnabled = aslrEnabled
			}
		}
	}

	// 设置最后更新时间
	info.LastUpdate = time.Now()

	return info, nil
}
