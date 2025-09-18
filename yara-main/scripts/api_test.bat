@echo off
setlocal enabledelayedexpansion

REM Yara安全服务API测试脚本 (Windows版本)
REM 使用方法: api_test.bat [base_url]
REM 默认base_url: http://localhost:8081

set BASE_URL=%1
if "%BASE_URL%"=="" set BASE_URL=http://localhost:8081

echo 🚀 开始API测试...
echo 📡 测试基础URL: %BASE_URL%
echo ⏰ 测试时间: %date% %time%
echo.

REM 测试计数器
set TOTAL_TESTS=0
set PASSED_TESTS=0
set FAILED_TESTS=0

REM 测试函数
:test_api
set test_name=%~1
set method=%~2
set endpoint=%~3
set data=%~4
set expected_status=%~5

set /a TOTAL_TESTS+=1

echo 🔍 测试: %test_name%
echo    方法: %method%
echo    端点: %endpoint%

REM 构建curl命令
set curl_cmd=curl -s -w "%%{http_code}" -X %method%

if "%method%"=="POST" (
    set curl_cmd=%curl_cmd% -H "Content-Type: application/json"
)
if "%method%"=="PUT" (
    set curl_cmd=%curl_cmd% -H "Content-Type: application/json"
)

if not "%data%"=="" (
    set curl_cmd=%curl_cmd% -d "%data%"
)

set curl_cmd=%curl_cmd% %BASE_URL%%endpoint%

REM 执行请求
for /f "delims=" %%i in ('%curl_cmd%') do set response=%%i

REM 提取状态码和响应体
for /f "tokens=*" %%a in ("%response%") do (
    set status_code=%%a
    set status_code=!status_code:~-3!
    set response_body=%%a
    set response_body=!response_body:~0,-3!
)

echo    状态码: !status_code!

REM 检查状态码
if "!status_code!"=="%expected_status%" (
    echo    ✅ 通过
    set /a PASSED_TESTS+=1
) else (
    echo    ❌ 失败 (期望: %expected_status%, 实际: !status_code!)
    echo    响应: !response_body!
    set /a FAILED_TESTS+=1
)

echo.
goto :eof

REM ==================== 系统监控功能测试 ====================
echo 🌐 系统监控功能测试
echo ====================

call :test_api "健康检查" "GET" "/api/health" "" "200"
call :test_api "性能指标" "GET" "/api/metrics" "" "200"
call :test_api "重置指标" "POST" "/api/metrics/reset" "" "200"

REM ==================== 文件管理功能测试 ====================
echo 📁 文件管理功能测试
echo ====================

call :test_api "单文件扫描" "POST" "/api/v1/file/scan" "{\"path\": \"C:\\\\Windows\\\\System32\\\\notepad.exe\"}" "200"
call :test_api "目录扫描" "POST" "/api/v1/file/scan-directory" "{\"path\": \"C:\\\\temp\", \"recursive\": true, \"max_depth\": 3}" "200"
call :test_api "缓冲区扫描" "POST" "/api/v1/file/scan-buffer" "{\"identifier\": \"test\", \"data\": \"MZ...\"}" "200"
call :test_api "文件信息获取" "GET" "/api/v1/file/info/C:\\Windows\\System32\\notepad.exe" "" "200"
call :test_api "文件列表获取" "GET" "/api/v1/file/list?path=C:\\temp" "" "200"
call :test_api "文件复制" "POST" "/api/v1/file/copy" "{\"source\": \"C:\\\\temp\\\\test.txt\", \"dest\": \"C:\\\\temp\\\\test_copy.txt\"}" "200"
call :test_api "文件移动" "POST" "/api/v1/file/move" "{\"source\": \"C:\\\\temp\\\\test.txt\", \"dest\": \"C:\\\\temp\\\\test_moved.txt\"}" "200"
call :test_api "文件删除" "DELETE" "/api/v1/file/C:\\temp\\test.txt" "" "200"
call :test_api "文件哈希计算" "GET" "/api/v1/file/hash/C:\\Windows\\System32\\notepad.exe?algorithm=sha256" "" "200"
call :test_api "文件哈希验证" "POST" "/api/v1/file/verify-hash" "{\"file_path\": \"C:\\\\Windows\\\\System32\\\\notepad.exe\", \"algorithm\": \"sha256\", \"expected_hash\": \"test\"}" "200"
call :test_api "文件多哈希获取" "GET" "/api/v1/file/hashes/C:\\Windows\\System32\\notepad.exe" "" "200"

REM ==================== 进程管理功能测试 ====================
echo 🔄 进程管理功能测试
echo ====================

call :test_api "进程列表获取" "GET" "/api/v1/process/list" "" "200"
call :test_api "进程信息获取" "GET" "/api/v1/process/1234" "" "200"
call :test_api "启动进程" "POST" "/api/v1/process/start" "{\"command\": \"notepad.exe\", \"args\": [], \"working_dir\": \"C:\\\\temp\"}" "200"
call :test_api "结束进程" "DELETE" "/api/v1/process/1234" "" "200"
call :test_api "暂停进程" "PUT" "/api/v1/process/1234/suspend" "" "200"
call :test_api "恢复进程" "PUT" "/api/v1/process/1234/resume" "" "200"
call :test_api "进程模块获取" "GET" "/api/v1/process/1234/modules" "" "200"
call :test_api "进程连接获取" "GET" "/api/v1/process/1234/connections" "" "200"
call :test_api "进程内存信息" "GET" "/api/v1/process/1234/memory" "" "200"
call :test_api "进程运行状态" "GET" "/api/v1/process/1234/running" "" "200"
call :test_api "进程子进程" "GET" "/api/v1/process/1234/children" "" "200"
call :test_api "系统模块列表" "GET" "/api/v1/process/modules" "" "200"
call :test_api "模块信息获取" "GET" "/api/v1/process/module/notepad.exe" "" "200"
call :test_api "启用进程监控" "POST" "/api/v1/process/monitoring/enable" "" "200"
call :test_api "禁用进程监控" "POST" "/api/v1/process/monitoring/disable" "" "200"
call :test_api "监控进程列表" "GET" "/api/v1/process/monitoring/list" "" "200"
call :test_api "进程统计信息" "GET" "/api/v1/process/statistics" "" "200"

REM ==================== 注册表管理功能测试 ====================
echo 🔧 注册表管理功能测试
echo ====================

call :test_api "注册表键获取" "GET" "/api/v1/registry/key/SOFTWARE\\Microsoft\\Windows\\CurrentVersion" "" "200"
call :test_api "创建注册表键" "POST" "/api/v1/registry/key" "{\"path\": \"SOFTWARE\\\\TestApp\", \"name\": \"TestKey\"}" "200"
call :test_api "删除注册表键" "DELETE" "/api/v1/registry/key/SOFTWARE\\TestApp" "" "200"
call :test_api "设置注册表值" "PUT" "/api/v1/registry/value" "{\"path\": \"SOFTWARE\\\\TestApp\", \"name\": \"TestValue\", \"type\": \"REG_SZ\", \"value\": \"test_value\"}" "200"
call :test_api "获取注册表值" "GET" "/api/v1/registry/value/SOFTWARE\\TestApp/TestValue" "" "200"
call :test_api "删除注册表值" "DELETE" "/api/v1/registry/value/SOFTWARE\\TestApp/TestValue" "" "200"
call :test_api "注册表键列表" "GET" "/api/v1/registry/keys/SOFTWARE\\Microsoft" "" "200"
call :test_api "注册表值列表" "GET" "/api/v1/registry/values/SOFTWARE\\Microsoft" "" "200"
call :test_api "注册表搜索" "GET" "/api/v1/registry/search?path=SOFTWARE&pattern=TestApp&recursive=true&max_depth=5" "" "200"

REM ==================== 网络管理功能测试 ====================
echo 🌐 网络管理功能测试
echo ====================

call :test_api "网络连接列表" "GET" "/api/v1/network/connections" "" "200"
call :test_api "TCP连接列表" "GET" "/api/v1/network/connections/tcp" "" "200"
call :test_api "UDP连接列表" "GET" "/api/v1/network/connections/udp" "" "200"
call :test_api "进程连接获取" "GET" "/api/v1/network/connections/pid/1234" "" "200"
call :test_api "端口连接获取" "GET" "/api/v1/network/connections/port/80" "" "200"
call :test_api "IP连接获取" "GET" "/api/v1/network/connections/ip/192.168.1.100" "" "200"
call :test_api "关闭连接" "DELETE" "/api/v1/network/connection/12345" "" "200"
call :test_api "端口使用检查" "GET" "/api/v1/network/port/80/in-use" "" "200"
call :test_api "监听端口列表" "GET" "/api/v1/network/listening-ports" "" "200"
call :test_api "已建立连接" "GET" "/api/v1/network/established-connections" "" "200"
call :test_api "网络接口列表" "GET" "/api/v1/network/interfaces" "" "200"
call :test_api "网络统计信息" "GET" "/api/v1/network/stats" "" "200"
call :test_api "启用网络监控" "POST" "/api/v1/network/monitoring/enable" "" "200"
call :test_api "禁用网络监控" "POST" "/api/v1/network/monitoring/disable" "" "200"
call :test_api "监控连接列表" "GET" "/api/v1/network/monitoring/connections" "" "200"
call :test_api "连接历史记录" "GET" "/api/v1/network/monitoring/history" "" "200"

REM ==================== 安全功能测试 ====================
echo 🛡️ 安全功能测试
echo ====================

call :test_api "安全状态获取" "GET" "/api/v1/security/status" "" "200"
call :test_api "规则信息获取" "GET" "/api/v1/security/rules" "" "200"
call :test_api "重新加载规则" "POST" "/api/v1/security/reload-rules" "" "200"
call :test_api "清除缓存" "POST" "/api/v1/security/cache/clear" "" "200"
call :test_api "缓存统计信息" "GET" "/api/v1/security/cache/stats" "" "200"
call :test_api "隔离文件" "POST" "/api/v1/security/quarantine" "{\"file_path\": \"C:\\\\temp\\\\malware.exe\", \"reason\": \"Malware detected\"}" "200"
call :test_api "恢复文件" "POST" "/api/v1/security/restore" "{\"file_path\": \"C:\\\\temp\\\\malware.exe\"}" "200"
call :test_api "隔离列表获取" "GET" "/api/v1/security/quarantine/list" "" "200"
call :test_api "扫描历史记录" "GET" "/api/v1/security/scan-history" "" "200"

REM ==================== 用户管理功能测试 ====================
echo 👥 用户管理功能测试
echo ====================

call :test_api "当前用户信息" "GET" "/api/v1/user/current" "" "200"
call :test_api "用户ID查询" "GET" "/api/v1/user/id/1000" "" "200"
call :test_api "用户名查询" "GET" "/api/v1/user/name/testuser" "" "200"
call :test_api "所有用户列表" "GET" "/api/v1/user/all" "" "200"
call :test_api "权限检查" "POST" "/api/v1/user/permissions/check" "{\"username\": \"testuser\", \"permissions\": [\"file_scan\", \"process_control\"]}" "200"
call :test_api "密码验证" "POST" "/api/v1/user/password/validate" "{\"username\": \"testuser\", \"password\": \"TestPassword123!\"}" "200"
call :test_api "密码策略获取" "GET" "/api/v1/user/password/policy" "" "200"
call :test_api "账户状态检查" "GET" "/api/v1/user/status/testuser" "" "200"
call :test_api "用户会话列表" "GET" "/api/v1/user/sessions/testuser" "" "200"
call :test_api "终止用户会话" "DELETE" "/api/v1/user/session/session_12345" "" "200"
call :test_api "锁定用户账户" "POST" "/api/v1/user/lock/testuser" "{\"reason\": \"Suspicious activity\"}" "200"
call :test_api "解锁用户账户" "POST" "/api/v1/user/unlock/testuser" "" "200"
call :test_api "更改用户密码" "POST" "/api/v1/user/password/change" "{\"username\": \"testuser\", \"old_password\": \"OldPassword123!\", \"new_password\": \"NewPassword456!\"}" "200"
call :test_api "用户组列表" "GET" "/api/v1/user/groups/testuser" "" "200"
call :test_api "登录历史记录" "GET" "/api/v1/user/history/testuser" "" "200"

REM ==================== 错误处理测试 ====================
echo ⚠️ 错误处理测试
echo ====================

call :test_api "无效文件路径" "POST" "/api/v1/file/scan" "{\"path\": \"C:\\\\nonexistent\\\\file.exe\"}" "404"
call :test_api "无效进程ID" "GET" "/api/v1/process/99999" "" "404"
call :test_api "无效注册表路径" "GET" "/api/v1/registry/key/INVALID\\PATH" "" "404"
call :test_api "无效网络连接ID" "DELETE" "/api/v1/network/connection/invalid_id" "" "404"
call :test_api "无效用户ID" "GET" "/api/v1/user/id/99999" "" "404"

REM ==================== 性能测试 ====================
echo ⚡ 性能测试
echo ====================

call :test_api "大量文件列表" "GET" "/api/v1/file/list?path=C:\\Windows&limit=1000" "" "200"
call :test_api "大量进程列表" "GET" "/api/v1/process/list?limit=1000" "" "200"
call :test_api "大量网络连接" "GET" "/api/v1/network/connections?limit=1000" "" "200"

REM ==================== 测试结果统计 ====================
echo.
echo ==================== 测试结果统计 ====================
echo 📊 总测试数: %TOTAL_TESTS%
echo ✅ 通过测试: %PASSED_TESTS%
echo ❌ 失败测试: %FAILED_TESTS%

set /a SUCCESS_RATE=(%PASSED_TESTS% * 100) / %TOTAL_TESTS%
echo 📈 成功率: %SUCCESS_RATE%%%

if %FAILED_TESTS% GTR 0 (
    echo ⚠️  有 %FAILED_TESTS% 个测试失败，请检查API服务状态
    exit /b 1
) else (
    echo 🎉 所有测试通过！
    exit /b 0
) 