package main

import (
	"fmt"
	"os"
	"time"

	"golang.org/x/sys/windows"
)

// TestWindowsAPIs 测试Windows API的基本功能
func TestWindowsAPIs() {
	fmt.Println("=== Windows API 用户状态信息获取方法可行性测试 ===")
	fmt.Println("开始时间:", time.Now().Format("2006-01-02 15:04:05"))

	// 测试1: 获取当前用户名
	fmt.Println("\n1. 测试获取当前用户名...")
	if username, err := os.UserHomeDir(); err == nil {
		fmt.Printf("   成功: 用户目录 = %s\n", username)
	} else {
		fmt.Printf("   失败: %v\n", err)
	}

	// 测试2: 尝试加载netapi32.dll
	fmt.Println("\n2. 测试加载netapi32.dll...")
	if netapi32, err := windows.LoadDLL("netapi32.dll"); err == nil {
		fmt.Println("   成功: netapi32.dll加载成功")

		// 尝试查找NetUserGetInfo函数
		if _, err := netapi32.FindProc("NetUserGetInfo"); err == nil {
			fmt.Println("   成功: NetUserGetInfo函数找到")
		} else {
			fmt.Printf("   失败: NetUserGetInfo函数未找到: %v\n", err)
		}
	} else {
		fmt.Printf("   失败: netapi32.dll加载失败: %v\n", err)
	}

	// 测试3: 尝试加载wtsapi32.dll
	fmt.Println("\n3. 测试加载wtsapi32.dll...")
	if wtsapi32, err := windows.LoadDLL("wtsapi32.dll"); err == nil {
		fmt.Println("   成功: wtsapi32.dll加载成功")

		// 尝试查找WTSQueryUserToken函数
		if _, err := wtsapi32.FindProc("WTSQueryUserToken"); err == nil {
			fmt.Println("   成功: WTSQueryUserToken函数找到")
		} else {
			fmt.Printf("   失败: WTSQueryUserToken函数未找到: %v\n", err)
		}
	} else {
		fmt.Printf("   失败: wtsapi32.dll加载失败: %v\n", err)
	}

	// 测试4: 尝试加载secur32.dll
	fmt.Println("\n4. 测试加载secur32.dll...")
	if secur32, err := windows.LoadDLL("secur32.dll"); err == nil {
		fmt.Println("   成功: secur32.dll加载成功")

		// 尝试查找GetUserNameExW函数
		if _, err := secur32.FindProc("GetUserNameExW"); err == nil {
			fmt.Println("   成功: GetUserNameExW函数找到")
		} else {
			fmt.Printf("   失败: GetUserNameExW函数未找到: %v\n", err)
		}
	} else {
		fmt.Printf("   失败: secur32.dll加载失败: %v\n", err)
	}

	// 测试5: 尝试加载advapi32.dll
	fmt.Println("\n5. 测试加载advapi32.dll...")
	if advapi32, err := windows.LoadDLL("advapi32.dll"); err == nil {
		fmt.Println("   成功: advapi32.dll加载成功")

		// 尝试查找LookupAccountSidW函数
		if _, err := advapi32.FindProc("LookupAccountSidW"); err == nil {
			fmt.Println("   成功: LookupAccountSidW函数找到")
		} else {
			fmt.Printf("   失败: LookupAccountSidW函数未找到: %v\n", err)
		}
	} else {
		fmt.Printf("   失败: advapi32.dll加载失败: %v\n", err)
	}

	// 测试6: 尝试获取当前进程令牌
	fmt.Println("\n6. 测试获取当前进程令牌...")
	if token, err := windows.OpenCurrentProcessToken(); err == nil {
		fmt.Println("   成功: 当前进程令牌获取成功")
		token.Close()
	} else {
		fmt.Printf("   失败: 当前进程令牌获取失败: %v\n", err)
	}

	fmt.Println("\n测试完成时间:", time.Now().Format("2006-01-02 15:04:05"))
}

// func main() {
// 	// 运行基本Windows API测试
// 	TestWindowsAPIs()

// 	// 创建用户状态测试实例
// 	test := NewUserStatusTest()

// 	// 运行所有测试
// 	test.RunAllTests()
// }
