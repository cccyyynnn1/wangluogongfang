package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
	"time"
	"unsafe"

	"golang.org/x/sys/windows"
)

// UserStatusTest 用户状态测试结构
type UserStatusTest struct {
	logger *log.Logger
}

// NewUserStatusTest 创建新的用户状态测试实例
func NewUserStatusTest() *UserStatusTest {
	return &UserStatusTest{
		logger: log.New(os.Stdout, "[UserStatusTest] ", log.LstdFlags),
	}
}

// TestResult 测试结果结构
type TestResult struct {
	Method        string        `json:"method"`
	Success       bool          `json:"success"`
	Error         string        `json:"error,omitempty"`
	Data          interface{}   `json:"data,omitempty"`
	ExecutionTime time.Duration `json:"execution_time"`
}

// UserStatusInfo 用户状态信息结构
type UserStatusInfo struct {
	Username        string    `json:"username"`
	FullName        string    `json:"full_name"`
	IsOnline        bool      `json:"is_online"`
	IsLocked        bool      `json:"is_locked"`
	LastLoginTime   time.Time `json:"last_login_time"`
	AccountDisabled bool      `json:"account_disabled"`
	PasswordExpired bool      `json:"password_expired"`
	HomeDirectory   string    `json:"home_directory"`
	ProfilePath     string    `json:"profile_path"`
	LogonScript     string    `json:"logon_script"`
	Description     string    `json:"description"`
	Workstations    string    `json:"workstations"`
	LogonHours      string    `json:"logon_hours"`
	UnitsPerWeek    uint32    `json:"units_per_week"`
	MaxStorage      uint32    `json:"max_storage"`
	LogonCount      uint32    `json:"logon_count"`
	CountryCode     uint32    `json:"country_code"`
	CodePage        uint32    `json:"code_page"`
	UserComment     string    `json:"user_comment"`
	Parms           string    `json:"parms"`
	AuthFlags       uint32    `json:"auth_flags"`
	LastLogoff      time.Time `json:"last_logoff"`
	LastLogon       time.Time `json:"last_logon"`
	BadPwCount      uint16    `json:"bad_pw_count"`
	NumLogons       uint16    `json:"num_logons"`
	LogonServer     string    `json:"logon_server"`
}

// Windows API常量
const (
	USER_PRIV_GUEST = 0
	USER_PRIV_USER  = 1
	USER_PRIV_ADMIN = 2

	UF_SCRIPT                                 = 0x0001
	UF_ACCOUNTDISABLE                         = 0x0002
	UF_HOMEDIR_REQUIRED                       = 0x0008
	UF_LOCKOUT                                = 0x0010
	UF_PASSWD_NOTREQD                         = 0x0020
	UF_PASSWD_CANT_CHANGE                     = 0x0040
	UF_ENCRYPTED_TEXT_PASSWORD_ALLOWED        = 0x0080
	UF_TEMP_DUPLICATE_ACCOUNT                 = 0x0100
	UF_NORMAL_ACCOUNT                         = 0x0200
	UF_INTERDOMAIN_TRUST_ACCOUNT              = 0x0400
	UF_WORKSTATION_TRUST_ACCOUNT              = 0x0800
	UF_SERVER_TRUST_ACCOUNT                   = 0x1000
	UF_DONT_EXPIRE_PASSWD                     = 0x20000
	UF_MNS_LOGON_ACCOUNT                      = 0x40000
	UF_SMARTCARD_REQUIRED                     = 0x40000
	UF_TRUSTED_FOR_DELEGATION                 = 0x80000
	UF_NOT_DELEGATED                          = 0x100000
	UF_USE_DES_KEY_ONLY                       = 0x200000
	UF_DONT_REQUIRE_PREAUTH                   = 0x400000
	UF_PASSWORD_EXPIRED                       = 0x800000
	UF_TRUSTED_TO_AUTHENTICATE_FOR_DELEGATION = 0x1000000
)

// Windows API结构体
type USER_INFO_1 struct {
	Name        *uint16
	Password    *uint16
	PasswordAge uint32
	Priv        uint32
	HomeDir     *uint16
	Comment     *uint16
	Flags       uint32
	ScriptPath  *uint16
}

type USER_INFO_3 struct {
	Name         *uint16
	Password     *uint16
	PasswordAge  uint32
	Priv         uint32
	HomeDir      *uint16
	Comment      *uint16
	Flags        uint32
	ScriptPath   *uint16
	AuthFlags    uint32
	FullName     *uint16
	UsrComment   *uint16
	Parms        *uint16
	Workstations *uint16
	LastLogon    uint32
	LastLogoff   uint32
	AcctExpires  uint32
	MaxStorage   uint32
	UnitsPerWeek uint32
	LogonHours   *uint16
	BadPwCount   uint16
	NumLogons    uint16
	LogonServer  *uint16
	CountryCode  uint16
	CodePage     uint16
}

// 测试所有方法
func (t *UserStatusTest) TestAllMethods() []TestResult {
	t.logger.Println("开始测试所有用户状态获取方法...")

	var results []TestResult

	// 测试Windows API方法
	results = append(results, t.testNetUserGetInfo())
	results = append(results, t.testNetUserEnum())
	results = append(results, t.testWTSQueryUserToken())
	results = append(results, t.testGetUserNameEx())
	results = append(results, t.testLookupAccountSid())

	// 测试命令行方法
	results = append(results, t.testNetUserCommand())
	results = append(results, t.testWMICUserAccount())
	results = append(results, t.testQUserCommand())
	results = append(results, t.testWhoamiCommand())

	// 测试注册表方法
	results = append(results, t.testRegistryUserInfo())

	return results
}

// 测试NetUserGetInfo方法
func (t *UserStatusTest) testNetUserGetInfo() TestResult {
	start := time.Now()
	result := TestResult{Method: "NetUserGetInfo"}

	defer func() {
		result.ExecutionTime = time.Since(start)
	}()

	// 加载netapi32.dll
	netapi32, err := windows.LoadDLL("netapi32.dll")
	if err != nil {
		result.Success = false
		result.Error = fmt.Sprintf("加载netapi32.dll失败: %v", err)
		return result
	}

	// 获取NetUserGetInfo函数
	netUserGetInfo, err := netapi32.FindProc("NetUserGetInfo")
	if err != nil {
		result.Success = false
		result.Error = fmt.Sprintf("查找NetUserGetInfo函数失败: %v", err)
		return result
	}

	// 获取当前用户名
	currentUser, err := os.UserHomeDir()
	if err != nil {
		result.Success = false
		result.Error = fmt.Sprintf("获取当前用户目录失败: %v", err)
		return result
	}
	// 简化用户名获取
	currentUser = "Administrator" // 使用默认用户名进行测试

	// 转换用户名为UTF16
	username, err := windows.UTF16PtrFromString(currentUser)
	if err != nil {
		result.Success = false
		result.Error = fmt.Sprintf("转换用户名为UTF16失败: %v", err)
		return result
	}

	// 调用NetUserGetInfo
	var buffer *byte
	var entriesRead uint32
	var totalEntries uint32
	var resumeHandle uint32

	ret, _, err := netUserGetInfo.Call(
		uintptr(0), // 本地服务器
		uintptr(unsafe.Pointer(username)),
		1, // 信息级别
		uintptr(unsafe.Pointer(&buffer)),
		uintptr(unsafe.Pointer(&entriesRead)),
		uintptr(unsafe.Pointer(&totalEntries)),
		uintptr(unsafe.Pointer(&resumeHandle)),
	)

	if ret != 0 {
		result.Success = false
		result.Error = fmt.Sprintf("NetUserGetInfo调用失败: ret=%d, err=%v", ret, err)
		return result
	}

	// 解析返回的数据
	if buffer != nil && entriesRead > 0 {
		userInfo := (*USER_INFO_1)(unsafe.Pointer(buffer))

		userStatus := &UserStatusInfo{
			Username:        currentUser,
			AccountDisabled: (userInfo.Flags & UF_ACCOUNTDISABLE) != 0,
			PasswordExpired: (userInfo.Flags & UF_PASSWORD_EXPIRED) != 0,
		}

		result.Success = true
		result.Data = userStatus

		// 释放内存
		if freeProc, err := netapi32.FindProc("NetApiBufferFree"); err == nil {
			freeProc.Call(uintptr(unsafe.Pointer(buffer)))
		}
	} else {
		result.Success = false
		result.Error = "NetUserGetInfo返回空数据"
	}

	return result
}

// 测试NetUserEnum方法
func (t *UserStatusTest) testNetUserEnum() TestResult {
	start := time.Now()
	result := TestResult{Method: "NetUserEnum"}

	defer func() {
		result.ExecutionTime = time.Since(start)
	}()

	// 加载netapi32.dll
	netapi32, err := windows.LoadDLL("netapi32.dll")
	if err != nil {
		result.Success = false
		result.Error = fmt.Sprintf("加载netapi32.dll失败: %v", err)
		return result
	}

	// 获取NetUserEnum函数
	netUserEnum, err := netapi32.FindProc("NetUserEnum")
	if err != nil {
		result.Success = false
		result.Error = fmt.Sprintf("查找NetUserEnum函数失败: %v", err)
		return result
	}

	// 调用NetUserEnum
	var buffer *byte
	var entriesRead uint32
	var totalEntries uint32
	var resumeHandle uint32

	ret, _, err := netUserEnum.Call(
		uintptr(0), // 本地服务器
		1,          // 信息级别
		0,          // 过滤器
		uintptr(unsafe.Pointer(&buffer)),
		0xFFFFFFFF, // 最大偏好
		uintptr(unsafe.Pointer(&entriesRead)),
		uintptr(unsafe.Pointer(&totalEntries)),
		uintptr(unsafe.Pointer(&resumeHandle)),
	)

	if ret != 0 {
		result.Success = false
		result.Error = fmt.Sprintf("NetUserEnum调用失败: ret=%d, err=%v", ret, err)
		return result
	}

	// 解析返回的数据
	if buffer != nil && entriesRead > 0 {
		var users []string
		ptr := buffer

		for i := uint32(0); i < entriesRead; i++ {
			userInfo := (*USER_INFO_1)(unsafe.Pointer(ptr))
			if userInfo.Name != nil {
				username := windows.UTF16ToString([]uint16{*userInfo.Name})
				users = append(users, username)
			}
			ptr = (*byte)(unsafe.Pointer(uintptr(unsafe.Pointer(ptr)) + unsafe.Sizeof(*userInfo)))
		}

		result.Success = true
		result.Data = map[string]interface{}{
			"total_users": totalEntries,
			"users":       users,
		}

		// 释放内存
		if freeProc, err := netapi32.FindProc("NetApiBufferFree"); err == nil {
			freeProc.Call(uintptr(unsafe.Pointer(buffer)))
		}
	} else {
		result.Success = false
		result.Error = "NetUserEnum返回空数据"
	}

	return result
}

// 测试WTSQueryUserToken方法
func (t *UserStatusTest) testWTSQueryUserToken() TestResult {
	start := time.Now()
	result := TestResult{Method: "WTSQueryUserToken"}

	defer func() {
		result.ExecutionTime = time.Since(start)
	}()

	// 加载wtsapi32.dll
	wtsapi32, err := windows.LoadDLL("wtsapi32.dll")
	if err != nil {
		result.Success = false
		result.Error = fmt.Sprintf("加载wtsapi32.dll失败: %v", err)
		return result
	}

	// 获取WTSQueryUserToken函数
	wtsQueryUserToken, err := wtsapi32.FindProc("WTSQueryUserToken")
	if err != nil {
		result.Success = false
		result.Error = fmt.Sprintf("查找WTSQueryUserToken函数失败: %v", err)
		return result
	}

	// 获取当前会话ID
	wtsGetActiveConsoleSessionId, err := wtsapi32.FindProc("WTSGetActiveConsoleSessionId")
	if err != nil {
		result.Success = false
		result.Error = fmt.Sprintf("查找WTSGetActiveConsoleSessionId函数失败: %v", err)
		return result
	}

	sessionId, _, err := wtsGetActiveConsoleSessionId.Call()
	if sessionId == 0xFFFFFFFF {
		result.Success = false
		result.Error = "获取活动会话ID失败"
		return result
	}

	// 查询用户令牌
	var userToken windows.Handle
	ret, _, err := wtsQueryUserToken.Call(
		sessionId,
		uintptr(unsafe.Pointer(&userToken)),
	)

	if ret == 0 {
		result.Success = false
		result.Error = fmt.Sprintf("WTSQueryUserToken调用失败: err=%v", err)
		return result
	}

	// 获取令牌信息
	tokenInfo := &UserStatusInfo{
		IsOnline: true,
	}

	result.Success = true
	result.Data = tokenInfo

	// 关闭令牌句柄
	windows.CloseHandle(userToken)

	return result
}

// 测试GetUserNameEx方法
func (t *UserStatusTest) testGetUserNameEx() TestResult {
	start := time.Now()
	result := TestResult{Method: "GetUserNameEx"}

	defer func() {
		result.ExecutionTime = time.Since(start)
	}()

	// 加载secur32.dll
	secur32, err := windows.LoadDLL("secur32.dll")
	if err != nil {
		result.Success = false
		result.Error = fmt.Sprintf("加载secur32.dll失败: %v", err)
		return result
	}

	// 获取GetUserNameExW函数
	getUserNameEx, err := secur32.FindProc("GetUserNameExW")
	if err != nil {
		result.Success = false
		result.Error = fmt.Sprintf("查找GetUserNameExW函数失败: %v", err)
		return result
	}

	// 调用GetUserNameExW
	var size uint32 = 256
	buffer := make([]uint16, size)

	ret, _, err := getUserNameEx.Call(
		2, // NameDisplay
		uintptr(unsafe.Pointer(&buffer[0])),
		uintptr(unsafe.Pointer(&size)),
	)

	if ret == 0 {
		result.Success = false
		result.Error = fmt.Sprintf("GetUserNameExW调用失败: err=%v", err)
		return result
	}

	displayName := windows.UTF16ToString(buffer[:size])

	result.Success = true
	result.Data = map[string]interface{}{
		"display_name": displayName,
	}

	return result
}

// 测试LookupAccountSid方法
func (t *UserStatusTest) testLookupAccountSid() TestResult {
	start := time.Now()
	result := TestResult{Method: "LookupAccountSid"}

	defer func() {
		result.ExecutionTime = time.Since(start)
	}()

	// 加载advapi32.dll
	advapi32, err := windows.LoadDLL("advapi32.dll")
	if err != nil {
		result.Success = false
		result.Error = fmt.Sprintf("加载advapi32.dll失败: %v", err)
		return result
	}

	// 获取LookupAccountSidW函数
	lookupAccountSid, err := advapi32.FindProc("LookupAccountSidW")
	if err != nil {
		result.Success = false
		result.Error = fmt.Sprintf("查找LookupAccountSidW函数失败: %v", err)
		return result
	}

	// 获取当前进程令牌
	processToken, err := windows.OpenCurrentProcessToken()
	if err != nil {
		result.Success = false
		result.Error = fmt.Sprintf("打开当前进程令牌失败: %v", err)
		return result
	}
	defer processToken.Close()

	// 获取令牌用户信息
	tokenUser, err := processToken.GetTokenUser()
	if err != nil {
		result.Success = false
		result.Error = fmt.Sprintf("获取令牌用户信息失败: %v", err)
		return result
	}

	// 调用LookupAccountSidW
	var nameSize uint32 = 256
	var domainSize uint32 = 256
	var sidType uint32

	nameBuffer := make([]uint16, nameSize)
	domainBuffer := make([]uint16, domainSize)

	// 获取SID的字节数组
	// 手动计算SID长度，因为tokenUser.User.SidLength可能不可用
	sidBytes := (*[1 << 30]byte)(unsafe.Pointer(tokenUser.User.Sid))[:8] // 使用固定长度8作为最小SID长度

	ret, _, err := lookupAccountSid.Call(
		uintptr(0), // 本地系统
		uintptr(unsafe.Pointer(&sidBytes[0])),
		uintptr(unsafe.Pointer(&nameBuffer[0])),
		uintptr(unsafe.Pointer(&nameSize)),
		uintptr(unsafe.Pointer(&domainBuffer[0])),
		uintptr(unsafe.Pointer(&domainSize)),
		uintptr(unsafe.Pointer(&sidType)),
	)

	if ret == 0 {
		result.Success = false
		result.Error = fmt.Sprintf("LookupAccountSidW调用失败: err=%v", err)
		return result
	}

	username := windows.UTF16ToString(nameBuffer[:nameSize])
	domain := windows.UTF16ToString(domainBuffer[:domainSize])

	result.Success = true
	result.Data = map[string]interface{}{
		"username": username,
		"domain":   domain,
		"sid_type": sidType,
	}

	return result
}

// 测试net user命令
func (t *UserStatusTest) testNetUserCommand() TestResult {
	start := time.Now()
	result := TestResult{Method: "net user command"}

	defer func() {
		result.ExecutionTime = time.Since(start)
	}()

	// 执行net user命令
	cmd := exec.Command("net", "user")
	output, err := cmd.Output()
	if err != nil {
		result.Success = false
		result.Error = fmt.Sprintf("执行net user命令失败: %v", err)
		return result
	}

	// 解析输出
	lines := strings.Split(string(output), "\n")
	var users []string

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line != "" && !strings.Contains(line, "命令成功完成") && !strings.Contains(line, "The command completed successfully") {
			users = append(users, line)
		}
	}

	result.Success = true
	result.Data = map[string]interface{}{
		"users": users,
	}

	return result
}

// 测试wmic useraccount命令
func (t *UserStatusTest) testWMICUserAccount() TestResult {
	start := time.Now()
	result := TestResult{Method: "wmic useraccount command"}

	defer func() {
		result.ExecutionTime = time.Since(start)
	}()

	// 执行wmic useraccount命令
	cmd := exec.Command("wmic", "useraccount", "get", "name,disabled,lockout,description", "/format:csv")
	output, err := cmd.Output()
	if err != nil {
		result.Success = false
		result.Error = fmt.Sprintf("执行wmic useraccount命令失败: %v", err)
		return result
	}

	// 解析CSV输出
	lines := strings.Split(string(output), "\n")
	var userAccounts []map[string]string

	for i, line := range lines {
		if i == 0 || strings.TrimSpace(line) == "" {
			continue
		}

		fields := strings.Split(line, ",")
		if len(fields) >= 4 {
			account := map[string]string{
				"name":        strings.TrimSpace(fields[0]),
				"disabled":    strings.TrimSpace(fields[1]),
				"lockout":     strings.TrimSpace(fields[2]),
				"description": strings.TrimSpace(fields[3]),
			}
			userAccounts = append(userAccounts, account)
		}
	}

	result.Success = true
	result.Data = map[string]interface{}{
		"user_accounts": userAccounts,
	}

	return result
}

// 测试quser命令
func (t *UserStatusTest) testQUserCommand() TestResult {
	start := time.Now()
	result := TestResult{Method: "quser command"}

	defer func() {
		result.ExecutionTime = time.Since(start)
	}()

	// 执行quser命令
	cmd := exec.Command("quser")
	output, err := cmd.Output()
	if err != nil {
		result.Success = false
		result.Error = fmt.Sprintf("执行quser命令失败: %v", err)
		return result
	}

	// 解析输出
	lines := strings.Split(string(output), "\n")
	var sessions []map[string]string

	for i, line := range lines {
		if i == 0 || strings.TrimSpace(line) == "" {
			continue
		}

		fields := strings.Fields(line)
		if len(fields) >= 4 {
			session := map[string]string{
				"username":   fields[0],
				"session":    fields[1],
				"idle_time":  fields[2],
				"logon_time": fields[3],
			}
			sessions = append(sessions, session)
		}
	}

	result.Success = true
	result.Data = map[string]interface{}{
		"sessions": sessions,
	}

	return result
}

// 测试whoami命令
func (t *UserStatusTest) testWhoamiCommand() TestResult {
	start := time.Now()
	result := TestResult{Method: "whoami command"}

	defer func() {
		result.ExecutionTime = time.Since(start)
	}()

	// 执行whoami命令
	cmd := exec.Command("whoami")
	output, err := cmd.Output()
	if err != nil {
		result.Success = false
		result.Error = fmt.Sprintf("执行whoami命令失败: %v", err)
		return result
	}

	username := strings.TrimSpace(string(output))

	// 执行whoami /all命令获取详细信息
	cmdAll := exec.Command("whoami", "/all")
	outputAll, err := cmdAll.Output()
	if err != nil {
		result.Success = false
		result.Error = fmt.Sprintf("执行whoami /all命令失败: %v", err)
		return result
	}

	result.Success = true
	result.Data = map[string]interface{}{
		"username": username,
		"details":  string(outputAll),
	}

	return result
}

// 测试注册表用户信息
func (t *UserStatusTest) testRegistryUserInfo() TestResult {
	start := time.Now()
	result := TestResult{Method: "registry user info"}

	defer func() {
		result.ExecutionTime = time.Since(start)
	}()

	// 尝试读取注册表用户信息
	// 注意：这需要管理员权限
	var key windows.Handle
	err := windows.RegOpenKeyEx(
		windows.HKEY_LOCAL_MACHINE,
		windows.StringToUTF16Ptr(`SOFTWARE\Microsoft\Windows NT\CurrentVersion\Winlogon`),
		0,
		windows.KEY_READ,
		&key,
	)
	if err != nil {
		result.Success = false
		result.Error = fmt.Sprintf("打开注册表键失败: %v", err)
		return result
	}
	defer windows.RegCloseKey(key)

	// 读取默认用户名
	var defaultUserNameSize uint32 = 256
	defaultUserName := make([]uint16, defaultUserNameSize)

	err = windows.RegQueryValueEx(
		key,
		windows.StringToUTF16Ptr("DefaultUserName"),
		nil,
		nil,
		(*byte)(unsafe.Pointer(&defaultUserName[0])),
		&defaultUserNameSize,
	)

	if err != nil {
		result.Success = false
		result.Error = fmt.Sprintf("读取DefaultUserName失败: %v", err)
		return result
	}

	username := windows.UTF16ToString(defaultUserName[:defaultUserNameSize/defaultUserNameSize])

	result.Success = true
	result.Data = map[string]interface{}{
		"default_username": username,
	}

	return result
}

// 运行所有测试并生成报告
func (t *UserStatusTest) RunAllTests() {
	t.logger.Println("=== 用户状态信息获取方法可行性测试 ===")
	t.logger.Println("开始时间:", time.Now().Format("2006-01-02 15:04:05"))

	results := t.TestAllMethods()

	// 统计结果
	successCount := 0
	totalCount := len(results)

	t.logger.Println("\n=== 测试结果汇总 ===")
	for _, result := range results {
		status := "失败"
		if result.Success {
			status = "成功"
			successCount++
		}

		t.logger.Printf("方法: %-25s | 状态: %s | 耗时: %v",
			result.Method, status, result.ExecutionTime)

		if !result.Success {
			t.logger.Printf("  错误: %s", result.Error)
		}
	}

	t.logger.Printf("\n成功率: %d/%d (%.1f%%)",
		successCount, totalCount, float64(successCount)/float64(totalCount)*100)

	// 推荐最佳方法
	t.logger.Println("\n=== 推荐方法 ===")
	for _, result := range results {
		if result.Success {
			t.logger.Printf("✓ %s - 推荐使用", result.Method)
		}
	}

	t.logger.Println("\n测试完成时间:", time.Now().Format("2006-01-02 15:04:05"))
}
