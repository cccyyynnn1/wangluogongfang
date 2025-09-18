package user

import (
	"fmt"
	"os"
	"os/exec"
	"sort"
	"strconv"
	"strings"
	"time"
	"unsafe"

	"golang.org/x/sys/windows"
)

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

// WindowsAPIEnhanced 提供Windows API增强功能的用户管理服务
type WindowsAPIEnhanced struct{}

// NewWindowsAPIEnhanced 创建新的Windows API增强实例
func NewWindowsAPIEnhanced() *WindowsAPIEnhanced {
	return &WindowsAPIEnhanced{}
}

// GetCurrentUserInfo 获取当前用户信息（推荐方法）
func (w *WindowsAPIEnhanced) GetCurrentUserInfo() (*UserStatusInfo, error) {
	// 方法1: 使用GetUserNameEx获取显示名称
	displayName, err := w.getUserNameEx()
	if err != nil {
		return nil, fmt.Errorf("GetUserNameEx失败: %v", err)
	}

	// 方法2: 使用whoami获取用户名
	username, err := w.getWhoamiUsername()
	if err != nil {
		return nil, fmt.Errorf("whoami失败: %v", err)
	}

	// 方法3: 使用LookupAccountSid获取详细信息
	userDetails, err := w.getLookupAccountSidInfo()
	if err != nil {
		// 不返回错误，因为这是补充信息
		userDetails = &UserStatusInfo{}
	}

	// 组合信息
	userInfo := &UserStatusInfo{
		Username:      username,
		FullName:      displayName,
		IsOnline:      true, // 当前用户默认在线
		HomeDirectory: w.getUserHomeDir(),
	}

	// 合并详细信息
	if userDetails != nil {
		userInfo.AccountDisabled = userDetails.AccountDisabled
		userInfo.PasswordExpired = userDetails.PasswordExpired
		userInfo.Description = userDetails.Description
	}

	return userInfo, nil
}

// GetAllUsers 获取系统中所有用户列表（推荐方法）
func (w *WindowsAPIEnhanced) GetAllUsers() ([]*UserStatusInfo, error) {
	// 主要方法: 使用NetUserEnum
	users, err := w.getNetUserEnumUsers()
	if err == nil && len(users) > 0 {
		return users, nil
	}

	// 备用方法: 使用net user命令
	users, err = w.getNetUserCommandUsers()
	if err == nil && len(users) > 0 {
		return users, nil
	}

	// 最后备用: 使用wmic useraccount
	users, err = w.getWMICUserAccountUsers()
	if err != nil {
		return nil, fmt.Errorf("所有用户获取方法都失败: %v", err)
	}

	return users, nil
}

// GetUserStatus 获取指定用户的状态信息
func (w *WindowsAPIEnhanced) GetUserStatus(username string) (*UserStatusInfo, error) {
	// 获取所有用户信息
	allUsers, err := w.GetAllUsers()
	if err != nil {
		return nil, err
	}

	// 查找指定用户
	for _, user := range allUsers {
		if user.Username == username {
			return user, nil
		}
	}

	return nil, fmt.Errorf("用户 %s 未找到", username)
}

// IsUserOnline 检查用户是否在线
func (w *WindowsAPIEnhanced) IsUserOnline(username string) (bool, error) {
	// 获取当前用户
	currentUser, err := w.getWhoamiUsername()
	if err != nil {
		return false, err
	}

	// 如果查询的是当前用户，直接返回true
	if currentUser == username {
		return true, nil
	}

	// 对于其他用户，尝试使用WTSQueryUserToken（如果可用）
	// 注意：这个方法在某些系统上可能不可用
	return w.checkUserOnlineStatus(username)
}

// GetUserStatusEnhanced 获取增强的用户状态信息
func (w *WindowsAPIEnhanced) GetUserStatusEnhanced(username string) (*UserStatusInfo, error) {
	// 获取基本用户状态
	userStatus, err := w.GetUserStatus(username)
	if err != nil {
		return nil, err
	}

	// 增强在线状态检查
	isOnline, err := w.IsUserOnline(username)
	if err == nil {
		userStatus.IsOnline = isOnline
	}

	// 获取登录历史
	loginHistory, err := w.GetUserLoginHistory(username)
	if err == nil && len(loginHistory) > 0 {
		userStatus.LastLoginTime = loginHistory[0]
		userStatus.LastLogon = loginHistory[0]
	}

	// 获取用户会话信息
	sessions, err := w.GetUserSessions(username)
	if err == nil && len(sessions) > 0 {
		userStatus.IsOnline = true
		// 可以添加更多会话相关信息
	}

	return userStatus, nil
}

// GetUserSessions 获取用户会话信息
func (w *WindowsAPIEnhanced) GetUserSessions(username string) ([]map[string]interface{}, error) {
	var sessions []map[string]interface{}

	// 尝试使用quser命令获取会话信息
	cmd := exec.Command("quser", username)
	output, err := cmd.Output()
	if err == nil {
		// 解析quser输出
		parsedSessions := w.parseQuserOutput(output, username)
		sessions = append(sessions, parsedSessions...)
	}

	// 如果quser失败，尝试使用tasklist检查进程
	if len(sessions) == 0 {
		cmd := exec.Command("tasklist", "/FI", fmt.Sprintf("USERNAME eq %s", username), "/FO", "CSV")
		output, err := cmd.Output()
		if err == nil {
			// 解析tasklist输出
			parsedSessions := w.parseTasklistOutput(output, username)
			sessions = append(sessions, parsedSessions...)
		}
	}

	return sessions, nil
}

// parseQuserOutput 解析quser命令输出
func (w *WindowsAPIEnhanced) parseQuserOutput(output []byte, username string) []map[string]interface{} {
	var sessions []map[string]interface{}

	lines := strings.Split(string(output), "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" || strings.Contains(line, "USERNAME") || strings.Contains(line, "用户名") {
			continue
		}

		fields := strings.Fields(line)
		if len(fields) >= 4 {
			session := map[string]interface{}{
				"username":   fields[0],
				"session":    fields[1],
				"id":         fields[2],
				"state":      fields[3],
				"idle_time":  "",
				"logon_time": "",
			}

			// 解析空闲时间和登录时间
			if len(fields) >= 6 {
				session["idle_time"] = fields[4]
				session["logon_time"] = fields[5]
			}

			sessions = append(sessions, session)
		}
	}

	return sessions
}

// parseTasklistOutput 解析tasklist命令输出
func (w *WindowsAPIEnhanced) parseTasklistOutput(output []byte, username string) []map[string]interface{} {
	var sessions []map[string]interface{}

	lines := strings.Split(string(output), "\n")
	for i, line := range lines {
		if i == 0 || strings.TrimSpace(line) == "" {
			continue
		}

		// 跳过标题行
		if strings.Contains(line, "Image Name") || strings.Contains(line, "映像名称") {
			continue
		}

		fields := strings.Split(line, ",")
		if len(fields) >= 5 {
			session := map[string]interface{}{
				"username": username,
				"process":  strings.Trim(fields[0], "\""),
				"pid":      strings.Trim(fields[1], "\""),
				"session":  strings.Trim(fields[2], "\""),
				"memory":   strings.Trim(fields[4], "\""),
				"type":     "process",
			}

			sessions = append(sessions, session)
		}
	}

	return sessions
}

// GetUserLoginHistory 获取用户登录历史
func (w *WindowsAPIEnhanced) GetUserLoginHistory(username string) ([]time.Time, error) {
	var loginHistory []time.Time

	// 方法1: 使用wmic useraccount获取登录信息
	userInfo, err := w.getWMICUserAccountInfo(username)
	if err == nil && userInfo.LastLogon != (time.Time{}) {
		loginHistory = append(loginHistory, userInfo.LastLogon)
	}

	// 方法2: 尝试从事件日志获取登录历史
	eventLogTimes, err := w.getLoginHistoryFromEventLog(username)
	if err == nil && len(eventLogTimes) > 0 {
		loginHistory = append(loginHistory, eventLogTimes...)
	}

	// 方法3: 尝试从net user命令获取最后登录时间
	netUserTime, err := w.getLastLoginFromNetUser(username)
	if err == nil && !netUserTime.IsZero() {
		loginHistory = append(loginHistory, netUserTime)
	}

	// 如果没有获取到任何登录时间，返回空列表
	if len(loginHistory) == 0 {
		return nil, fmt.Errorf("无法获取用户 %s 的登录历史", username)
	}

	// 按时间排序，最新的在前
	sort.Slice(loginHistory, func(i, j int) bool {
		return loginHistory[i].After(loginHistory[j])
	})

	return loginHistory, nil
}

// 私有方法实现

// getUserNameEx 使用GetUserNameEx获取用户显示名称
func (w *WindowsAPIEnhanced) getUserNameEx() (string, error) {
	secur32, err := windows.LoadDLL("secur32.dll")
	if err != nil {
		return "", fmt.Errorf("加载secur32.dll失败: %v", err)
	}

	getUserNameEx, err := secur32.FindProc("GetUserNameExW")
	if err != nil {
		return "", fmt.Errorf("查找GetUserNameExW函数失败: %v", err)
	}

	var size uint32 = 256
	buffer := make([]uint16, size)

	ret, _, err := getUserNameEx.Call(
		2, // NameDisplay
		uintptr(unsafe.Pointer(&buffer[0])),
		uintptr(unsafe.Pointer(&size)),
	)

	if ret == 0 {
		return "", fmt.Errorf("GetUserNameExW调用失败: %v", err)
	}

	displayName := windows.UTF16ToString(buffer[:size])
	return displayName, nil
}

// getWhoamiUsername 使用whoami命令获取当前用户名
func (w *WindowsAPIEnhanced) getWhoamiUsername() (string, error) {
	cmd := exec.Command("whoami")
	output, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("执行whoami命令失败: %v", err)
	}

	username := strings.TrimSpace(string(output))
	// 移除域名前缀（如果存在）
	if idx := strings.Index(username, "\\"); idx != -1 {
		username = username[idx+1:]
	}

	return username, nil
}

// getLookupAccountSidInfo 使用LookupAccountSid获取用户详细信息
func (w *WindowsAPIEnhanced) getLookupAccountSidInfo() (*UserStatusInfo, error) {
	advapi32, err := windows.LoadDLL("advapi32.dll")
	if err != nil {
		return nil, fmt.Errorf("加载advapi32.dll失败: %v", err)
	}

	lookupAccountSid, err := advapi32.FindProc("LookupAccountSidW")
	if err != nil {
		return nil, fmt.Errorf("查找LookupAccountSidW函数失败: %v", err)
	}

	// 获取当前进程令牌
	processToken, err := windows.OpenCurrentProcessToken()
	if err != nil {
		return nil, fmt.Errorf("打开当前进程令牌失败: %v", err)
	}
	defer processToken.Close()

	// 获取令牌用户信息
	tokenUser, err := processToken.GetTokenUser()
	if err != nil {
		return nil, fmt.Errorf("获取令牌用户信息失败: %v", err)
	}

	// 调用LookupAccountSidW
	var nameSize uint32 = 256
	var domainSize uint32 = 256
	var sidType uint32

	nameBuffer := make([]uint16, nameSize)
	domainBuffer := make([]uint16, domainSize)

	// 获取SID的字节数组
	sidBytes := (*[1 << 30]byte)(unsafe.Pointer(tokenUser.User.Sid))[:8]

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
		return nil, fmt.Errorf("LookupAccountSidW调用失败: %v", err)
	}

	username := windows.UTF16ToString(nameBuffer[:nameSize])
	domain := windows.UTF16ToString(domainBuffer[:domainSize])

	return &UserStatusInfo{
		Username:    username,
		LogonServer: domain,
	}, nil
}

// getUserHomeDir 获取用户主目录
func (w *WindowsAPIEnhanced) getUserHomeDir() string {
	if homeDir, err := os.UserHomeDir(); err == nil {
		return homeDir
	}
	return ""
}

// getNetUserEnumUsers 使用NetUserEnum获取用户列表
func (w *WindowsAPIEnhanced) getNetUserEnumUsers() ([]*UserStatusInfo, error) {
	netapi32, err := windows.LoadDLL("netapi32.dll")
	if err != nil {
		return nil, fmt.Errorf("加载netapi32.dll失败: %v", err)
	}

	netUserEnum, err := netapi32.FindProc("NetUserEnum")
	if err != nil {
		return nil, fmt.Errorf("查找NetUserEnum函数失败: %v", err)
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
		return nil, fmt.Errorf("NetUserEnum调用失败: ret=%d, err=%v", ret, err)
	}

	// 解析返回的数据
	if buffer == nil || entriesRead == 0 {
		return nil, fmt.Errorf("NetUserEnum返回空数据")
	}

	var users []*UserStatusInfo
	ptr := buffer

	for i := uint32(0); i < entriesRead; i++ {
		userInfo := (*USER_INFO_1)(unsafe.Pointer(ptr))
		if userInfo.Name != nil {
			username := windows.UTF16ToString([]uint16{*userInfo.Name})
			user := &UserStatusInfo{
				Username:        username,
				AccountDisabled: (userInfo.Flags & UF_ACCOUNTDISABLE) != 0,
				PasswordExpired: (userInfo.Flags & UF_PASSWORD_EXPIRED) != 0,
			}
			users = append(users, user)
		}
		ptr = (*byte)(unsafe.Pointer(uintptr(unsafe.Pointer(ptr)) + unsafe.Sizeof(*userInfo)))
	}

	// 释放内存
	if freeProc, err := netapi32.FindProc("NetApiBufferFree"); err == nil {
		freeProc.Call(uintptr(unsafe.Pointer(buffer)))
	}

	return users, nil
}

// getNetUserCommandUsers 使用net user命令获取用户列表
func (w *WindowsAPIEnhanced) getNetUserCommandUsers() ([]*UserStatusInfo, error) {
	// 首先尝试直接执行net user命令
	cmd := exec.Command("net", "user")
	output, err := cmd.Output()
	if err != nil {
		// 如果失败，尝试使用PowerShell执行
		psScript := "net user"
		cmd = exec.Command("powershell", "-Command", psScript)
		output, err = cmd.Output()
		if err != nil {
			return nil, fmt.Errorf("执行net user命令失败: %v", err)
		}
	}

	lines := strings.Split(string(output), "\n")
	var users []*UserStatusInfo

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line != "" && !strings.Contains(line, "命令成功完成") && !strings.Contains(line, "The command completed successfully") && !strings.Contains(line, "用户帐户") && !strings.Contains(line, "User accounts") && !strings.Contains(line, "---") {
			user := &UserStatusInfo{
				Username: line,
			}
			users = append(users, user)
		}
	}

	return users, nil
}

// getWMICUserAccountUsers 使用wmic useraccount命令获取用户列表
func (w *WindowsAPIEnhanced) getWMICUserAccountUsers() ([]*UserStatusInfo, error) {
	// 尝试多种wmic命令格式
	commands := []string{
		"wmic useraccount get name /format:csv",
		"wmic useraccount get name",
		"wmic useraccount list brief",
	}

	var users []*UserStatusInfo

	for _, cmdStr := range commands {
		parts := strings.Fields(cmdStr)
		cmd := exec.Command(parts[0], parts[1:]...)
		output, err := cmd.Output()
		if err == nil {
			users, err = w.parseWMICOutput(output, cmdStr)
			if err == nil && len(users) > 0 {
				return users, nil
			}
		}
	}

	return nil, fmt.Errorf("所有wmic命令都失败")
}

// parseWMICOutput 解析WMIC命令输出
func (w *WindowsAPIEnhanced) parseWMICOutput(output []byte, cmdStr string) ([]*UserStatusInfo, error) {
	var users []*UserStatusInfo
	lines := strings.Split(string(output), "\n")

	if strings.Contains(cmdStr, "/format:csv") {
		// CSV格式解析
		for i, line := range lines {
			if i == 0 || strings.TrimSpace(line) == "" {
				continue
			}
			fields := strings.Split(line, ",")
			if len(fields) >= 1 {
				username := strings.TrimSpace(fields[0])
				if username != "" && username != "Name" {
					user := &UserStatusInfo{
						Username: username,
					}
					users = append(users, user)
				}
			}
		}
	} else {
		// 标准格式解析
		for _, line := range lines {
			line = strings.TrimSpace(line)
			if line != "" && !strings.Contains(strings.ToLower(line), "name") && !strings.Contains(strings.ToLower(line), "---") {
				// 提取用户名（通常是第一个非空字段）
				fields := strings.Fields(line)
				if len(fields) > 0 {
					username := fields[0]
					if username != "" {
						user := &UserStatusInfo{
							Username: username,
						}
						users = append(users, user)
					}
				}
			}
		}
	}

	return users, nil
}

// getWMICUserAccountInfo 获取指定用户的WMI信息
func (w *WindowsAPIEnhanced) getWMICUserAccountInfo(username string) (*UserStatusInfo, error) {
	// 尝试多种wmic命令格式
	commands := []string{
		fmt.Sprintf("wmic useraccount where name='%s' get name,disabled /format:csv", username),
		fmt.Sprintf("wmic useraccount where name='%s' get name,disabled", username),
		fmt.Sprintf("wmic useraccount get name,disabled /format:csv"),
		fmt.Sprintf("wmic useraccount get name,disabled"),
	}

	for _, cmdStr := range commands {
		parts := strings.Fields(cmdStr)
		cmd := exec.Command(parts[0], parts[1:]...)
		output, err := cmd.Output()
		if err == nil {
			user, err := w.parseWMICUserInfo(output, username, cmdStr)
			if err == nil && user != nil {
				return user, nil
			}
		}
	}

	// 如果所有wmic命令都失败，尝试使用PowerShell
	return w.getUserInfoWithPowerShell(username)
}

// parseWMICUserInfo 解析WMIC用户信息输出
func (w *WindowsAPIEnhanced) parseWMICUserInfo(output []byte, targetUsername, cmdStr string) (*UserStatusInfo, error) {
	lines := strings.Split(string(output), "\n")

	if strings.Contains(cmdStr, "/format:csv") {
		// CSV格式解析
		for i, line := range lines {
			if i == 0 || strings.TrimSpace(line) == "" {
				continue
			}
			fields := strings.Split(line, ",")
			if len(fields) >= 2 {
				currentUsername := strings.TrimSpace(fields[0])
				if currentUsername == targetUsername {
					user := &UserStatusInfo{
						Username:        currentUsername,
						AccountDisabled: strings.TrimSpace(fields[1]) == "TRUE",
						IsLocked:        false, // 默认值
					}
					return user, nil
				}
			}
		}
	} else {
		// 标准格式解析
		for _, line := range lines {
			line = strings.TrimSpace(line)
			if line != "" && !strings.Contains(strings.ToLower(line), "name") && !strings.Contains(strings.ToLower(line), "---") {
				fields := strings.Fields(line)
				if len(fields) >= 2 {
					currentUsername := fields[0]
					if currentUsername == targetUsername {
						user := &UserStatusInfo{
							Username:        currentUsername,
							AccountDisabled: strings.TrimSpace(fields[1]) == "TRUE",
							IsLocked:        false, // 默认值
						}
						return user, nil
					}
				}
			}
		}
	}

	return nil, fmt.Errorf("未找到用户 %s 的信息", targetUsername)
}

// getUserInfoWithPowerShell 使用PowerShell获取用户信息
func (w *WindowsAPIEnhanced) getUserInfoWithPowerShell(username string) (*UserStatusInfo, error) {
	// 使用PowerShell的Get-LocalUser命令
	psScript := fmt.Sprintf("Get-LocalUser -Name '%s' | Select-Object Name,Enabled,Locked | ConvertTo-Csv -NoTypeInformation", username)
	cmd := exec.Command("powershell", "-Command", psScript)
	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("PowerShell获取用户信息失败: %v", err)
	}

	// 解析CSV输出
	lines := strings.Split(string(output), "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" || strings.Contains(line, "Name") {
			continue
		}
		fields := strings.Split(line, ",")
		if len(fields) >= 3 {
			currentUsername := strings.Trim(strings.TrimSpace(fields[0]), "\"")
			if currentUsername == username {
				enabled := strings.Trim(strings.TrimSpace(fields[1]), "\"")
				locked := strings.Trim(strings.TrimSpace(fields[2]), "\"")

				user := &UserStatusInfo{
					Username:        currentUsername,
					AccountDisabled: enabled == "False",
					IsLocked:        locked == "True",
				}
				return user, nil
			}
		}
	}

	return nil, fmt.Errorf("PowerShell未找到用户 %s 的信息", username)
}

// checkUserOnlineStatus 检查用户在线状态
func (w *WindowsAPIEnhanced) checkUserOnlineStatus(username string) (bool, error) {
	// 尝试使用WTSQueryUserToken检查用户会话
	// 注意：这个方法在某些系统上可能不可用
	wtsapi32, err := windows.LoadDLL("wtsapi32.dll")
	if err != nil {
		return false, fmt.Errorf("加载wtsapi32.dll失败: %v", err)
	}

	// 尝试查找WTSQueryUserToken函数
	if _, err := wtsapi32.FindProc("WTSQueryUserToken"); err != nil {
		return false, fmt.Errorf("WTSQueryUserToken函数不可用: %v", err)
	}

	// 由于WTSGetActiveConsoleSessionId在某些系统上不可用
	// 我们使用一个简化的方法：检查用户是否有活动进程
	return w.checkUserActiveProcesses(username)
}

// checkUserActiveProcesses 检查用户是否有活动进程
func (w *WindowsAPIEnhanced) checkUserActiveProcesses(username string) (bool, error) {
	// 使用tasklist命令检查用户进程
	cmd := exec.Command("tasklist", "/FI", fmt.Sprintf("USERNAME eq %s", username), "/FO", "CSV")
	output, err := cmd.Output()
	if err != nil {
		return false, fmt.Errorf("执行tasklist命令失败: %v", err)
	}

	lines := strings.Split(string(output), "\n")
	// 如果有超过1行输出（标题行+进程行），说明用户有活动进程
	return len(lines) > 1, nil
}

// getLoginHistoryFromEventLog 从Windows事件日志获取登录历史
func (w *WindowsAPIEnhanced) getLoginHistoryFromEventLog(username string) ([]time.Time, error) {
	var loginTimes []time.Time

	// 使用wevtutil命令查询安全事件日志中的登录事件
	// EventID 4624 = 成功登录
	query := fmt.Sprintf("*[System[EventID=4624] and EventData[Data[@Name='TargetUserName']='%s']]", username)
	cmd := exec.Command("wevtutil", "qe", "Security", "/q:"+query, "/c:10", "/f:json")

	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("查询事件日志失败: %v", err)
	}

	// 解析JSON输出，提取时间戳
	times, err := w.parseEventLogTimes(output)
	if err != nil {
		return nil, fmt.Errorf("解析事件日志失败: %v", err)
	}

	loginTimes = append(loginTimes, times...)
	return loginTimes, nil
}

// getLastLoginFromNetUser 从net user命令获取最后登录时间
func (w *WindowsAPIEnhanced) getLastLoginFromNetUser(username string) (time.Time, error) {
	// 尝试多种方法获取最后登录时间

	// 方法1: 使用net user命令（处理用户名包含空格的情况）
	var cmd *exec.Cmd
	if strings.Contains(username, " ") {
		// 如果用户名包含空格，使用PowerShell执行
		psScript := fmt.Sprintf("net user '%s'", username)
		cmd = exec.Command("powershell", "-Command", psScript)
	} else {
		cmd = exec.Command("net", "user", username)
	}

	output, err := cmd.Output()
	if err == nil {
		// 解析net user输出
		lastLogin, err := w.parseNetUserLastLogin(output)
		if err == nil {
			return lastLogin, nil
		}
	}

	// 方法2: 如果net user失败，尝试使用PowerShell的Get-LocalUser
	psScript := fmt.Sprintf("Get-LocalUser -Name '%s' | Select-Object LastLogon | ConvertTo-Csv -NoTypeInformation", username)
	cmd = exec.Command("powershell", "-Command", psScript)
	output, err = cmd.Output()
	if err == nil {
		lastLogin, err := w.parsePowerShellLastLogin(output)
		if err == nil {
			return lastLogin, nil
		}
	}

	// 方法3: 尝试使用wmic useraccount
	cmd = exec.Command("wmic", "useraccount", "where", fmt.Sprintf("name='%s'", username), "get", "lastlogon", "/format:csv")
	output, err = cmd.Output()
	if err == nil {
		lastLogin, err := w.parseWMICLastLogin(output)
		if err == nil {
			return lastLogin, nil
		}
	}

	// 方法4: 尝试从事件日志获取
	lastLogin, err := w.getLoginHistoryFromEventLog(username)
	if err == nil && len(lastLogin) > 0 {
		return lastLogin[0], nil
	}

	return time.Time{}, fmt.Errorf("所有方法都无法获取用户 %s 的最后登录时间", username)
}

// parseNetUserLastLogin 解析net user命令输出中的最后登录时间
func (w *WindowsAPIEnhanced) parseNetUserLastLogin(output []byte) (time.Time, error) {
	lines := strings.Split(string(output), "\n")
	for _, line := range lines {
		line = strings.ToLower(strings.TrimSpace(line))
		if strings.Contains(line, "last logon") || strings.Contains(line, "最后登录") || strings.Contains(line, "lastlogon") {
			// 提取时间信息
			parts := strings.Fields(line)
			if len(parts) >= 2 {
				timeStr := strings.Join(parts[1:], " ")
				// 尝试多种时间格式
				timeFormats := []string{
					"2006-01-02 15:04:05",
					"01/02/2006 15:04:05",
					"2006-01-02",
					"01/02/2006",
					"2006-01-02t15:04:05.0000000z07:00",
					"2006-01-02t15:04:05z",
				}

				for _, format := range timeFormats {
					if t, err := time.Parse(format, timeStr); err == nil {
						return t, nil
					}
				}
			}
		}
	}
	return time.Time{}, fmt.Errorf("无法从net user输出解析最后登录时间")
}

// parsePowerShellLastLogin 解析PowerShell Get-LocalUser输出中的最后登录时间
func (w *WindowsAPIEnhanced) parsePowerShellLastLogin(output []byte) (time.Time, error) {
	lines := strings.Split(string(output), "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" || strings.Contains(line, "LastLogon") {
			continue
		}

		// 解析CSV格式
		fields := strings.Split(line, ",")
		if len(fields) >= 2 {
			timeStr := strings.Trim(strings.TrimSpace(fields[1]), "\"")
			if timeStr != "" && timeStr != "N/A" {
				// 尝试多种时间格式
				timeFormats := []string{
					"2006-01-02 15:04:05",
					"2006-01-02T15:04:05.0000000Z07:00",
					"2006-01-02T15:04:05Z",
					"2006-01-02",
				}

				for _, format := range timeFormats {
					if t, err := time.Parse(format, timeStr); err == nil {
						return t, nil
					}
				}
			}
		}
	}
	return time.Time{}, fmt.Errorf("无法从PowerShell输出解析最后登录时间")
}

// parseWMICLastLogin 解析WMIC命令输出中的最后登录时间
func (w *WindowsAPIEnhanced) parseWMICLastLogin(output []byte) (time.Time, error) {
	lines := strings.Split(string(output), "\n")
	for i, line := range lines {
		if i == 0 || strings.TrimSpace(line) == "" {
			continue
		}

		fields := strings.Split(line, ",")
		if len(fields) >= 2 {
			timeStr := strings.TrimSpace(fields[1])
			if timeStr != "" && timeStr != "LastLogon" {
				// WMIC返回的是文件时间（从1601年1月1日开始的100纳秒间隔数）
				if timeStr != "0" && timeStr != "16010101000000.000000+000" {
					// 尝试解析文件时间格式
					if t, err := w.parseFileTime(timeStr); err == nil {
						return t, nil
					}
				}
			}
		}
	}
	return time.Time{}, fmt.Errorf("无法从WMIC输出解析最后登录时间")
}

// parseFileTime 解析Windows文件时间格式
func (w *WindowsAPIEnhanced) parseFileTime(timeStr string) (time.Time, error) {
	// 尝试解析文件时间格式
	// 文件时间是从1601年1月1日开始的100纳秒间隔数
	// 格式可能是: 133456789012345678 或 133456789012345678.000000+000

	// 移除小数部分和时区信息
	if idx := strings.Index(timeStr, "."); idx != -1 {
		timeStr = timeStr[:idx]
	}
	if idx := strings.Index(timeStr, "+"); idx != -1 {
		timeStr = timeStr[:idx]
	}
	if idx := strings.Index(timeStr, "-"); idx != -1 {
		timeStr = timeStr[:idx]
	}

	// 转换为int64
	fileTime, err := strconv.ParseInt(timeStr, 10, 64)
	if err != nil {
		return time.Time{}, fmt.Errorf("解析文件时间失败: %v", err)
	}

	// 转换为纳秒（文件时间是100纳秒间隔）
	nanoseconds := fileTime * 100

	// 转换为秒
	seconds := nanoseconds / 1000000000

	// 转换为time.Time（从1601年1月1日开始）
	// 1601年1月1日到1970年1月1日的秒数
	const fileTimeEpoch = 11644473600
	unixTime := seconds - fileTimeEpoch

	return time.Unix(unixTime, 0), nil
}

// parseEventLogTimes 解析事件日志中的时间戳
func (w *WindowsAPIEnhanced) parseEventLogTimes(output []byte) ([]time.Time, error) {
	var times []time.Time
	outputStr := string(output)

	// 查找所有TimeCreated字段
	timeCreatedIndex := 0
	for {
		timeIndex := strings.Index(outputStr[timeCreatedIndex:], "TimeCreated")
		if timeIndex == -1 {
			break
		}

		timeCreatedIndex += timeIndex

		// 查找时间值
		startIndex := strings.Index(outputStr[timeCreatedIndex:], "\"")
		if startIndex == -1 {
			continue
		}

		endIndex := strings.Index(outputStr[timeCreatedIndex+startIndex+1:], "\"")
		if endIndex == -1 {
			continue
		}

		timeStr := outputStr[timeCreatedIndex+startIndex+1 : timeCreatedIndex+startIndex+1+endIndex]

		// 尝试解析ISO 8601格式
		if t, err := time.Parse(time.RFC3339, timeStr); err == nil {
			times = append(times, t)
		} else {
			// 尝试其他格式
			if t, err := time.Parse("2006-01-02T15:04:05.0000000Z07:00", timeStr); err == nil {
				times = append(times, t)
			}
		}

		timeCreatedIndex += startIndex + endIndex + 2
	}

	return times, nil
}
