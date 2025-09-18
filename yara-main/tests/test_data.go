package main

import (
	"time"
	"yara-security-service/internal/models"
)

// TestData 测试数据
type TestData struct {
	Files       []*models.FileInfo
	Processes   []*models.ProcessInfo
	Registry    []*models.RegistryKey
	Network     []*models.NetworkConnection
	Users       []*models.UserInfo
	ScanResults []*models.ScanResult
	Security    *models.SecurityStatus
	Metrics     map[string]interface{}
}

// GetTestData 获取测试数据
func GetTestData() *TestData {
	return &TestData{
		Files:       getTestFiles(),
		Processes:   getTestProcesses(),
		Registry:    getTestRegistry(),
		Network:     getTestNetwork(),
		Users:       getTestUsers(),
		ScanResults: getTestScanResults(),
		Security:    getTestSecurityStatus(),
		Metrics:     getTestMetrics(),
	}
}

// getTestFiles 获取测试文件数据
func getTestFiles() []*models.FileInfo {
	now := time.Now()
	return []*models.FileInfo{
		{
			Path:         "/usr/bin/ls",
			Name:         "ls",
			Size:         138496,
			IsDir:        false,
			ModTime:      now.Add(-24 * time.Hour),
			CreateTime:   now.Add(-24 * time.Hour),
			AccessTime:   now,
			Permissions:  "-rwxr-xr-x",
			Owner:        "root",
			Group:        "root",
			MD5:          "d41d8cd98f00b204e9800998ecf8427e",
			SHA256:       "e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855",
			IsSuspicious: false,
			ThreatLevel:  "safe",
		},
		{
			Path:         "/tmp/suspicious.exe",
			Name:         "suspicious.exe",
			Size:         2048000,
			IsDir:        false,
			ModTime:      now.Add(-1 * time.Hour),
			CreateTime:   now.Add(-1 * time.Hour),
			AccessTime:   now,
			Permissions:  "-rw-r--r--",
			Owner:        "testuser",
			Group:        "users",
			MD5:          "a1b2c3d4e5f678901234567890123456",
			SHA256:       "1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef",
			IsSuspicious: true,
			ThreatLevel:  "high",
		},
		{
			Path:         "/usr/lib/libc.so.6",
			Name:         "libc.so.6",
			Size:         2097152,
			IsDir:        false,
			ModTime:      now.Add(-12 * time.Hour),
			CreateTime:   now.Add(-12 * time.Hour),
			AccessTime:   now,
			Permissions:  "-rwxr-xr-x",
			Owner:        "root",
			Group:        "root",
			MD5:          "f1e2d3c4b5a678901234567890123456",
			SHA256:       "abcdef1234567890abcdef1234567890abcdef1234567890abcdef1234567890",
			IsSuspicious: false,
			ThreatLevel:  "medium",
		},
		{
			Path:         "/tmp/test.txt",
			Name:         "test.txt",
			Size:         1024,
			IsDir:        false,
			ModTime:      now.Add(-30 * time.Minute),
			CreateTime:   now.Add(-30 * time.Minute),
			AccessTime:   now,
			Permissions:  "-rw-r--r--",
			Owner:        "testuser",
			Group:        "users",
			MD5:          "098f6bcd4621d373cade4e832627b4f6",
			SHA256:       "a94a8fe5ccb19ba61c4c0873d391e987982fbbd3",
			IsSuspicious: false,
			ThreatLevel:  "safe",
		},
		{
			Path:         "/tmp/malware.exe",
			Name:         "malware.exe",
			Size:         1048576,
			IsDir:        false,
			ModTime:      now.Add(-5 * time.Minute),
			CreateTime:   now.Add(-5 * time.Minute),
			AccessTime:   now,
			Permissions:  "-rwxr-xr-x",
			Owner:        "testuser",
			Group:        "users",
			MD5:          "malware_md5_hash_here",
			SHA256:       "malware_sha256_hash_here",
			IsSuspicious: true,
			ThreatLevel:  "critical",
		},
		{
			Path:         "/home/testuser/.ssh/id_rsa",
			Name:         "id_rsa",
			Size:         1679,
			IsDir:        false,
			ModTime:      now.Add(-7 * 24 * time.Hour),
			CreateTime:   now.Add(-7 * 24 * time.Hour),
			AccessTime:   now,
			Permissions:  "-rw-------",
			Owner:        "testuser",
			Group:        "testuser",
			MD5:          "ssh_key_md5_hash",
			SHA256:       "ssh_key_sha256_hash",
			IsSuspicious: false,
			ThreatLevel:  "safe",
		},
		{
			Path:         "/var/log/auth.log",
			Name:         "auth.log",
			Size:         5242880,
			IsDir:        false,
			ModTime:      now.Add(-10 * time.Minute),
			CreateTime:   now.Add(-10 * time.Minute),
			AccessTime:   now,
			Permissions:  "-rw-r-----",
			Owner:        "root",
			Group:        "adm",
			MD5:          "log_file_md5_hash",
			SHA256:       "log_file_sha256_hash",
			IsSuspicious: false,
			ThreatLevel:  "safe",
		},
	}
}

// getTestProcesses 获取测试进程数据
func getTestProcesses() []*models.ProcessInfo {
	now := time.Now()
	return []*models.ProcessInfo{
		{
			PID:        1234,
			Name:       "yara-server",
			Cmdline:    "/usr/bin/yara-server --config /etc/yara/config.yaml",
			Exe:        "/usr/bin/yara-server",
			Status:     "running",
			CPUPercent: 2.5,
			MemoryInfo: &models.MemoryInfo{
				RSS: 1024000,
				VMS: 2048000,
			},
			CreateTime: now.Add(-30 * time.Minute),
			Username:   "yara",
			PPID:       1,
			NumThreads: 5,
			NumFiles:   10,
			Priority:   20,
			Modules:    []string{"libc.so.6", "libyara.so.4"},
			Children:   []int32{},
			Connections: []*models.ConnectionInfo{
				{
					FD:     3,
					Family: 2,
					Type:   1,
					Laddr:  &models.AddrInfo{IP: "127.0.0.1", Port: 8081},
					Raddr:  &models.AddrInfo{IP: "0.0.0.0", Port: 0},
					Status: "LISTEN",
				},
			},
		},
		{
			PID:        5678,
			Name:       "bash",
			Cmdline:    "bash",
			Exe:        "/usr/bin/bash",
			Status:     "running",
			CPUPercent: 1.2,
			MemoryInfo: &models.MemoryInfo{
				RSS: 2048000,
				VMS: 4096000,
			},
			CreateTime:  now.Add(-2 * time.Hour),
			Username:    "testuser",
			PPID:        1,
			NumThreads:  8,
			NumFiles:    25,
			Priority:    19,
			Modules:     []string{"libc.so.6", "libdl.so.2", "libtinfo.so.6"},
			Children:    []int32{1234},
			Connections: []*models.ConnectionInfo{},
		},
		{
			PID:        9999,
			Name:       "suspicious_process",
			Cmdline:    "suspicious_process --malware",
			Exe:        "/tmp/suspicious_process",
			Status:     "running",
			CPUPercent: 15.7,
			MemoryInfo: &models.MemoryInfo{
				RSS: 5120000,
				VMS: 10485760,
			},
			CreateTime: now.Add(-5 * time.Minute),
			Username:   "testuser",
			PPID:       5678,
			NumThreads: 12,
			NumFiles:   50,
			Priority:   15,
			Modules:    []string{"libc.so.6", "libssl.so.1.1", "libcrypto.so.1.1"},
			Children:   []int32{},
			Connections: []*models.ConnectionInfo{
				{
					FD:     3,
					Family: 2,
					Type:   1,
					Laddr:  &models.AddrInfo{IP: "192.168.1.100", Port: 12345},
					Raddr:  &models.AddrInfo{IP: "192.168.1.1", Port: 80},
					Status: "ESTABLISHED",
				},
			},
		},
		{
			PID:        1001,
			Name:       "nginx",
			Cmdline:    "nginx: master process",
			Exe:        "/usr/sbin/nginx",
			Status:     "running",
			CPUPercent: 0.8,
			MemoryInfo: &models.MemoryInfo{
				RSS: 3072000,
				VMS: 6144000,
			},
			CreateTime: now.Add(-1 * time.Hour),
			Username:   "www-data",
			PPID:       1,
			NumThreads: 6,
			NumFiles:   30,
			Priority:   18,
			Modules:    []string{"libc.so.6", "libpcre.so.3", "libssl.so.1.1"},
			Children:   []int32{1002, 1003},
			Connections: []*models.ConnectionInfo{
				{
					FD:     4,
					Family: 2,
					Type:   1,
					Laddr:  &models.AddrInfo{IP: "0.0.0.0", Port: 80},
					Raddr:  &models.AddrInfo{IP: "0.0.0.0", Port: 0},
					Status: "LISTEN",
				},
			},
		},
		{
			PID:        1002,
			Name:       "nginx",
			Cmdline:    "nginx: worker process",
			Exe:        "/usr/sbin/nginx",
			Status:     "running",
			CPUPercent: 0.5,
			MemoryInfo: &models.MemoryInfo{
				RSS: 1536000,
				VMS: 3072000,
			},
			CreateTime:  now.Add(-1 * time.Hour),
			Username:    "www-data",
			PPID:        1001,
			NumThreads:  4,
			NumFiles:    15,
			Priority:    18,
			Modules:     []string{"libc.so.6", "libpcre.so.3"},
			Children:    []int32{},
			Connections: []*models.ConnectionInfo{},
		},
	}
}

// getTestRegistry 获取测试注册表数据
func getTestRegistry() []*models.RegistryKey {
	now := time.Now()
	return []*models.RegistryKey{
		{
			Path:         "SOFTWARE\\Microsoft\\Windows\\CurrentVersion",
			Name:         "CurrentVersion",
			Type:         "REG_KEY",
			Value:        nil,
			SubKeys:      []string{"Run", "Policies", "Explorer"},
			Values:       map[string]string{"ProgramFilesDir": "C:\\Program Files", "CommonFilesDir": "C:\\Program Files\\Common Files"},
			LastModified: now.Add(-1 * time.Hour),
		},
		{
			Path:         "SOFTWARE\\TestApp",
			Name:         "TestApp",
			Type:         "REG_KEY",
			Value:        nil,
			SubKeys:      []string{},
			Values:       map[string]string{"InstallPath": "C:\\TestApp", "Version": "1.0.0"},
			LastModified: now.Add(-30 * time.Minute),
		},
		{
			Path:         "SYSTEM\\CurrentControlSet\\Services",
			Name:         "Services",
			Type:         "REG_KEY",
			Value:        nil,
			SubKeys:      []string{"Dhcp", "Dnscache", "Tcpip"},
			Values:       map[string]string{},
			LastModified: now.Add(-2 * time.Hour),
		},
		{
			Path:         "SOFTWARE\\Microsoft\\Windows\\CurrentVersion\\Run",
			Name:         "Run",
			Type:         "REG_KEY",
			Value:        nil,
			SubKeys:      []string{},
			Values:       map[string]string{"WindowsDefender": "C:\\Program Files\\Windows Defender\\MsMpEng.exe", "OneDrive": "C:\\Users\\testuser\\AppData\\Local\\Microsoft\\OneDrive\\OneDrive.exe"},
			LastModified: now.Add(-15 * time.Minute),
		},
		{
			Path:         "SYSTEM\\CurrentControlSet\\Control\\Session Manager",
			Name:         "Session Manager",
			Type:         "REG_KEY",
			Value:        nil,
			SubKeys:      []string{"AppCompatCache", "BootExecute"},
			Values:       map[string]string{"BootExecute": "autocheck autochk *"},
			LastModified: now.Add(-3 * time.Hour),
		},
	}
}

// getTestNetwork 获取测试网络连接数据
func getTestNetwork() []*models.NetworkConnection {
	return []*models.NetworkConnection{
		{
			ID:          "1",
			LocalAddr:   "127.0.0.1",
			RemoteAddr:  "127.0.0.1",
			LocalPort:   8081,
			RemotePort:  0,
			Protocol:    "TCP",
			Status:      "LISTEN",
			PID:         1234,
			ProcessName: "yara-server",
			Type:        "tcp",
		},
		{
			ID:          "2",
			LocalAddr:   "192.168.1.100",
			RemoteAddr:  "8.8.8.8",
			LocalPort:   12345,
			RemotePort:  53,
			Protocol:    "UDP",
			Status:      "ESTABLISHED",
			PID:         5678,
			ProcessName: "systemd-resolve",
			Type:        "udp",
		},
		{
			ID:          "3",
			LocalAddr:   "192.168.1.100",
			RemoteAddr:  "192.168.1.1",
			LocalPort:   54321,
			RemotePort:  80,
			Protocol:    "TCP",
			Status:      "ESTABLISHED",
			PID:         9999,
			ProcessName: "suspicious_process",
			Type:        "tcp",
		},
		{
			ID:          "4",
			LocalAddr:   "0.0.0.0",
			RemoteAddr:  "0.0.0.0",
			LocalPort:   22,
			RemotePort:  0,
			Protocol:    "TCP",
			Status:      "LISTEN",
			PID:         1000,
			ProcessName: "sshd",
			Type:        "tcp",
		},
		{
			ID:          "5",
			LocalAddr:   "0.0.0.0",
			RemoteAddr:  "0.0.0.0",
			LocalPort:   80,
			RemotePort:  0,
			Protocol:    "TCP",
			Status:      "LISTEN",
			PID:         1001,
			ProcessName: "nginx",
			Type:        "tcp",
		},
		{
			ID:          "6",
			LocalAddr:   "192.168.1.100",
			RemoteAddr:  "52.84.123.45",
			LocalPort:   45678,
			RemotePort:  443,
			Protocol:    "TCP",
			Status:      "ESTABLISHED",
			PID:         2001,
			ProcessName: "chrome",
			Type:        "tcp",
		},
		{
			ID:          "7",
			LocalAddr:   "192.168.1.100",
			RemoteAddr:  "192.168.1.255",
			LocalPort:   137,
			RemotePort:  137,
			Protocol:    "UDP",
			Status:      "ESTABLISHED",
			PID:         2002,
			ProcessName: "systemd-resolve",
			Type:        "udp",
		},
	}
}

// getTestUsers 获取测试用户数据
func getTestUsers() []*models.UserInfo {
	now := time.Now()
	return []*models.UserInfo{
		{
			UID:         "1000",
			GID:         "1000",
			Username:    "testuser",
			Name:        "Test User",
			HomeDir:     "/home/testuser",
			Shell:       "/bin/bash",
			LastLogin:   now.Add(-1 * time.Hour),
			IsActive:    true,
			IsAdmin:     false,
			PasswordAge: 30,
			Groups:      []string{"users", "sudo"},
		},
		{
			UID:         "0",
			GID:         "0",
			Username:    "root",
			Name:        "root",
			HomeDir:     "/root",
			Shell:       "/bin/bash",
			LastLogin:   now.Add(-30 * time.Minute),
			IsActive:    true,
			IsAdmin:     true,
			PasswordAge: 90,
			Groups:      []string{"root", "sudo", "admin"},
		},
		{
			UID:         "1001",
			GID:         "1001",
			Username:    "admin",
			Name:        "Administrator",
			HomeDir:     "/home/admin",
			Shell:       "/bin/bash",
			LastLogin:   now.Add(-2 * time.Hour),
			IsActive:    true,
			IsAdmin:     true,
			PasswordAge: 45,
			Groups:      []string{"users", "sudo", "admin"},
		},
		{
			UID:         "65534",
			GID:         "65534",
			Username:    "nobody",
			Name:        "nobody",
			HomeDir:     "/nonexistent",
			Shell:       "/usr/sbin/nologin",
			LastLogin:   now.Add(-24 * time.Hour),
			IsActive:    false,
			IsAdmin:     false,
			PasswordAge: 0,
			Groups:      []string{"nogroup"},
		},
		{
			UID:         "1002",
			GID:         "1002",
			Username:    "yara",
			Name:        "Yara Service User",
			HomeDir:     "/home/yara",
			Shell:       "/bin/bash",
			LastLogin:   now.Add(-10 * time.Minute),
			IsActive:    true,
			IsAdmin:     false,
			PasswordAge: 60,
			Groups:      []string{"yara", "users"},
		},
		{
			UID:         "33",
			GID:         "33",
			Username:    "www-data",
			Name:        "www-data",
			HomeDir:     "/var/www",
			Shell:       "/usr/sbin/nologin",
			LastLogin:   now.Add(-5 * time.Hour),
			IsActive:    true,
			IsAdmin:     false,
			PasswordAge: 0,
			Groups:      []string{"www-data"},
		},
	}
}

// getTestScanResults 获取测试扫描结果数据
func getTestScanResults() []*models.ScanResult {
	now := time.Now()
	return []*models.ScanResult{
		{
			FilePath:     "/usr/bin/ls",
			IsInfected:   false,
			Threats:      []models.ThreatInfo{},
			ScanTime:     now.Add(-5 * time.Minute),
			ScanDuration: 2 * time.Second,
			FileInfo: &models.FileInfo{
				Path:         "/usr/bin/ls",
				Name:         "ls",
				Size:         138496,
				IsSuspicious: false,
				ThreatLevel:  "safe",
			},
			Metadata: map[string]string{
				"file_size":    "138496",
				"file_type":    "executable",
				"scan_engine":  "yara",
				"threat_count": "0",
			},
		},
		{
			FilePath:   "/tmp/suspicious.exe",
			IsInfected: true,
			Threats: []models.ThreatInfo{
				{
					RuleName:    "Malware_Generic",
					Description: "检测到通用恶意软件特征",
					Severity:    "High",
					Category:    "Malware",
					Tags:        "malware,generic",
				},
			},
			ScanTime:     now.Add(-2 * time.Minute),
			ScanDuration: 5 * time.Second,
			FileInfo: &models.FileInfo{
				Path:         "/tmp/suspicious.exe",
				Name:         "suspicious.exe",
				Size:         2048000,
				IsSuspicious: true,
				ThreatLevel:  "high",
			},
			Metadata: map[string]string{
				"file_size":    "2048000",
				"file_type":    "executable",
				"scan_engine":  "yara",
				"threat_count": "1",
			},
		},
		{
			FilePath:   "/tmp/malware.exe",
			IsInfected: true,
			Threats: []models.ThreatInfo{
				{
					RuleName:    "Trojan_Generic",
					Description: "检测到特洛伊木马",
					Severity:    "Critical",
					Category:    "Trojan",
					Tags:        "trojan,backdoor",
				},
				{
					RuleName:    "Shellcode_Generic",
					Description: "检测到Shellcode载荷",
					Severity:    "Critical",
					Category:    "Shellcode",
					Tags:        "shellcode,exploit",
				},
			},
			ScanTime:     now.Add(-1 * time.Minute),
			ScanDuration: 8 * time.Second,
			FileInfo: &models.FileInfo{
				Path:         "/tmp/malware.exe",
				Name:         "malware.exe",
				Size:         1048576,
				IsSuspicious: true,
				ThreatLevel:  "critical",
			},
			Metadata: map[string]string{
				"file_size":    "1048576",
				"file_type":    "executable",
				"scan_engine":  "yara",
				"threat_count": "2",
			},
		},
		{
			FilePath:     "/usr/lib/libc.so.6",
			IsInfected:   false,
			Threats:      []models.ThreatInfo{},
			ScanTime:     now.Add(-3 * time.Minute),
			ScanDuration: 3 * time.Second,
			FileInfo: &models.FileInfo{
				Path:         "/usr/lib/libc.so.6",
				Name:         "libc.so.6",
				Size:         2097152,
				IsSuspicious: false,
				ThreatLevel:  "medium",
			},
			Metadata: map[string]string{
				"file_size":    "2097152",
				"file_type":    "shared_library",
				"scan_engine":  "yara",
				"threat_count": "0",
			},
		},
	}
}

// getTestSecurityStatus 获取测试安全状态数据
func getTestSecurityStatus() *models.SecurityStatus {
	now := time.Now()
	return &models.SecurityStatus{
		IsProtected:        true,
		LastScanTime:       now.Add(-10 * time.Minute),
		ThreatCount:        3,
		QuarantineCount:    1,
		RealTimeProtection: true,
		DefinitionsVersion: "2024.1.1",
		LastUpdateTime:     now.Add(-1 * time.Hour),
		ScanStatistics: &models.ScanStatistics{
			TotalScans:    150,
			InfectedFiles: 3,
			CleanFiles:    147,
			FailedScans:   0,
			TotalScanTime: 25 * time.Minute,
		},
		RulesInfo: map[string]interface{}{
			"total_rules":     25,
			"active_rules":    23,
			"last_updated":    now.Add(-2 * time.Hour),
			"rule_categories": []string{"Malware", "Trojan", "Shellcode", "Network"},
		},
	}
}

// getTestMetrics 获取测试指标数据
func getTestMetrics() map[string]interface{} {
	now := time.Now()
	return map[string]interface{}{
		"api_requests": map[string]interface{}{
			"total":   1250,
			"success": 1200,
			"failed":  50,
			"rate":    12.5,
		},
		"scan_operations": map[string]interface{}{
			"total_files":    150,
			"infected_files": 3,
			"clean_files":    147,
			"scan_time_avg":  2.5,
		},
		"system_resources": map[string]interface{}{
			"cpu_usage":    15.2,
			"memory_usage": 45.8,
			"disk_usage":   67.3,
			"network_io":   1024.5,
		},
		"cache_performance": map[string]interface{}{
			"hit_rate":   85.5,
			"miss_rate":  14.5,
			"cache_size": 1000,
			"evictions":  25,
		},
		"process_monitoring": map[string]interface{}{
			"total_processes":      156,
			"monitored_processes":  12,
			"suspicious_processes": 2,
		},
		"network_monitoring": map[string]interface{}{
			"total_connections":       45,
			"established_connections": 23,
			"listening_ports":         8,
			"suspicious_connections":  1,
		},
		"last_updated": now,
	}
}

// GetTestScanRequest 获取测试扫描请求
func GetTestScanRequest() models.DirectoryScanRequest {
	return models.DirectoryScanRequest{
		Path:            "/tmp",
		Recursive:       true,
		MaxDepth:        3,
		IncludePatterns: []string{"*.exe", "*.dll", "*.bat", "*.ps1"},
		ExcludePatterns: []string{"*.tmp", "*.log"},
	}
}

// GetTestProcessOperationRequest 获取测试进程操作请求
func GetTestProcessOperationRequest() models.ProcessOperationRequest {
	return models.ProcessOperationRequest{
		PID:       1234,
		Operation: "suspend",
		Timeout:   30,
	}
}

// GetTestRegistryOperationRequest 获取测试注册表操作请求
func GetTestRegistryOperationRequest() models.RegistryOperationRequest {
	return models.RegistryOperationRequest{
		Path:  "SOFTWARE\\TestApp",
		Name:  "TestValue",
		Type:  "REG_SZ",
		Value: "test_value",
	}
}

// GetTestNetworkOperationRequest 获取测试网络操作请求
func GetTestNetworkOperationRequest() models.NetworkOperationRequest {
	return models.NetworkOperationRequest{
		ConnectionID: "12345",
		Operation:    "close",
	}
}

// GetTestUserPermissionsRequest 获取测试用户权限请求
func GetTestUserPermissionsRequest() map[string]interface{} {
	return map[string]interface{}{
		"username":    "testuser",
		"permissions": []string{"file_scan", "process_control", "network_monitor"},
	}
}

// GetTestPasswordValidationRequest 获取测试密码验证请求
func GetTestPasswordValidationRequest() map[string]interface{} {
	return map[string]interface{}{
		"username": "testuser",
		"password": "TestPassword123!",
	}
}

// GetTestFileOperationRequest 获取测试文件操作请求
func GetTestFileOperationRequest() map[string]interface{} {
	return map[string]interface{}{
		"source": "/tmp/test.txt",
		"dest":   "/tmp/test_copy.txt",
	}
}

// GetTestHashVerificationRequest 获取测试哈希验证请求
func GetTestHashVerificationRequest() map[string]interface{} {
	return map[string]interface{}{
		"file_path":     "/usr/bin/ls",
		"algorithm":     "sha256",
		"expected_hash": "e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855",
	}
}

// GetTestBufferScanRequest 获取测试缓冲区扫描请求
func GetTestBufferScanRequest() models.BufferScanRequest {
	return models.BufferScanRequest{
		Data:       "MZ...",
		Identifier: "test_buffer",
	}
}

// GetTestRegistrySearchRequest 获取测试注册表搜索请求
func GetTestRegistrySearchRequest() models.RegistrySearchRequest {
	return models.RegistrySearchRequest{
		Path:      "SOFTWARE\\Microsoft",
		Pattern:   "TestApp",
		Recursive: true,
		MaxDepth:  5,
	}
}

// GetTestPasswordChangeRequest 获取测试密码更改请求
func GetTestPasswordChangeRequest() models.PasswordChangeRequest {
	return models.PasswordChangeRequest{
		Username:    "testuser",
		OldPassword: "OldPassword123!",
		NewPassword: "NewPassword456!",
	}
}

// GetTestUserLockRequest 获取测试用户锁定请求
func GetTestUserLockRequest() models.UserLockRequest {
	return models.UserLockRequest{
		Username: "suspicious_user",
		Reason:   "Suspicious activity detected",
	}
}

// GetTestSessionKillRequest 获取测试会话终止请求
func GetTestSessionKillRequest() models.SessionKillRequest {
	return models.SessionKillRequest{
		SessionID: "session_12345",
		Reason:    "Security policy violation",
	}
}
