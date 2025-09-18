#!/bin/bash

# Yara安全服务API测试脚本 (Linux/Mac版本)
# 使用方法: ./api_test.sh [base_url]
# 默认base_url: http://localhost:8081

BASE_URL=${1:-"http://localhost:8081"}
API_BASE="$BASE_URL/api/v1"

echo "🚀 开始API测试..."
echo "📡 测试基础URL: $BASE_URL"
echo "⏰ 测试时间: $(date)"
echo ""

# 颜色定义
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# 测试计数器
TOTAL_TESTS=0
PASSED_TESTS=0
FAILED_TESTS=0

# 测试函数
test_api() {
    local test_name="$1"
    local method="$2"
    local endpoint="$3"
    local data="$4"
    local expected_status="$5"
    
    TOTAL_TESTS=$((TOTAL_TESTS + 1))
    
    echo -e "${BLUE}🔍 测试: $test_name${NC}"
    echo "   方法: $method"
    echo "   端点: $endpoint"
    
    # 构建curl命令
    local curl_cmd="curl -s -w \"%{http_code}\" -X $method"
    
    if [ "$method" = "POST" ] || [ "$method" = "PUT" ]; then
        curl_cmd="$curl_cmd -H \"Content-Type: application/json\""
    fi
    
    if [ -n "$data" ]; then
        curl_cmd="$curl_cmd -d '$data'"
    fi
    
    curl_cmd="$curl_cmd $BASE_URL$endpoint"
    
    # 执行请求
    local response=$(eval $curl_cmd)
    local status_code="${response: -3}"
    local response_body="${response%???}"
    
    echo "   状态码: $status_code"
    
    # 检查状态码
    if [ "$status_code" = "$expected_status" ]; then
        echo -e "   ${GREEN}✅ 通过${NC}"
        PASSED_TESTS=$((PASSED_TESTS + 1))
    else
        echo -e "   ${RED}❌ 失败 (期望: $expected_status, 实际: $status_code)${NC}"
        echo "   响应: $response_body"
        FAILED_TESTS=$((FAILED_TESTS + 1))
    fi
    
    echo ""
}

# ==================== 系统监控功能测试 ====================
echo "🌐 系统监控功能测试"
echo "==================="

test_api "健康检查" "GET" "/api/health" "" "200"
test_api "性能指标" "GET" "/api/metrics" "" "200"
test_api "重置指标" "POST" "/api/metrics/reset" "" "200"

# ==================== 文件管理功能测试 ====================
echo "📁 文件管理功能测试"
echo "==================="

test_api "单文件扫描" "POST" "/api/v1/file/scan" '{"path": "/usr/bin/ls"}' "200"
test_api "目录扫描" "POST" "/api/v1/file/scan-directory" '{"path": "/tmp", "recursive": true, "max_depth": 3}' "200"
test_api "缓冲区扫描" "POST" "/api/v1/file/scan-buffer" '{"identifier": "test", "data": "MZ..."}' "200"
test_api "文件信息获取" "GET" "/api/v1/file/info//usr/bin/ls" "" "200"
test_api "文件列表获取" "GET" "/api/v1/file/list?path=/tmp" "" "200"
test_api "文件复制" "POST" "/api/v1/file/copy" '{"source": "/tmp/test.txt", "dest": "/tmp/test_copy.txt"}' "200"
test_api "文件移动" "POST" "/api/v1/file/move" '{"source": "/tmp/test.txt", "dest": "/tmp/test_moved.txt"}' "200"
test_api "文件删除" "DELETE" "/api/v1/file//tmp/test.txt" "" "200"
test_api "文件哈希计算" "GET" "/api/v1/file/hash//usr/bin/ls?algorithm=sha256" "" "200"
test_api "文件哈希验证" "POST" "/api/v1/file/verify-hash" '{"file_path": "/usr/bin/ls", "algorithm": "sha256", "expected_hash": "test"}' "200"
test_api "文件多哈希获取" "GET" "/api/v1/file/hashes//usr/bin/ls" "" "200"

# ==================== 进程管理功能测试 ====================
echo "🔄 进程管理功能测试"
echo "==================="

test_api "进程列表获取" "GET" "/api/v1/process/list" "" "200"
test_api "进程信息获取" "GET" "/api/v1/process/1234" "" "200"
test_api "启动进程" "POST" "/api/v1/process/start" '{"command": "ls", "args": ["-la"], "working_dir": "/tmp"}' "200"
test_api "结束进程" "DELETE" "/api/v1/process/1234" "" "200"
test_api "暂停进程" "PUT" "/api/v1/process/1234/suspend" "" "200"
test_api "恢复进程" "PUT" "/api/v1/process/1234/resume" "" "200"
test_api "进程模块获取" "GET" "/api/v1/process/1234/modules" "" "200"
test_api "进程连接获取" "GET" "/api/v1/process/1234/connections" "" "200"
test_api "进程内存信息" "GET" "/api/v1/process/1234/memory" "" "200"
test_api "进程运行状态" "GET" "/api/v1/process/1234/running" "" "200"
test_api "进程子进程" "GET" "/api/v1/process/1234/children" "" "200"
test_api "系统模块列表" "GET" "/api/v1/process/modules" "" "200"
test_api "模块信息获取" "GET" "/api/v1/process/module/libc.so.6" "" "200"
test_api "启用进程监控" "POST" "/api/v1/process/monitoring/enable" "" "200"
test_api "禁用进程监控" "POST" "/api/v1/process/monitoring/disable" "" "200"
test_api "监控进程列表" "GET" "/api/v1/process/monitoring/list" "" "200"
test_api "进程统计信息" "GET" "/api/v1/process/statistics" "" "200"

# ==================== 注册表管理功能测试 ====================
echo "🔧 注册表管理功能测试"
echo "==================="

test_api "注册表键获取" "GET" "/api/v1/registry/key/SOFTWARE\\Microsoft\\Windows\\CurrentVersion" "" "200"
test_api "创建注册表键" "POST" "/api/v1/registry/key" '{"path": "SOFTWARE\\TestApp", "name": "TestKey"}' "200"
test_api "删除注册表键" "DELETE" "/api/v1/registry/key/SOFTWARE\\TestApp" "" "200"
test_api "设置注册表值" "PUT" "/api/v1/registry/value" '{"path": "SOFTWARE\\TestApp", "name": "TestValue", "type": "REG_SZ", "value": "test_value"}' "200"
test_api "获取注册表值" "GET" "/api/v1/registry/value/SOFTWARE\\TestApp/TestValue" "" "200"
test_api "删除注册表值" "DELETE" "/api/v1/registry/value/SOFTWARE\\TestApp/TestValue" "" "200"
test_api "注册表键列表" "GET" "/api/v1/registry/keys/SOFTWARE\\Microsoft" "" "200"
test_api "注册表值列表" "GET" "/api/v1/registry/values/SOFTWARE\\Microsoft" "" "200"
test_api "注册表搜索" "GET" "/api/v1/registry/search?path=SOFTWARE&pattern=TestApp&recursive=true&max_depth=5" "" "200"

# ==================== 网络管理功能测试 ====================
echo "🌐 网络管理功能测试"
echo "==================="

test_api "网络连接列表" "GET" "/api/v1/network/connections" "" "200"
test_api "TCP连接列表" "GET" "/api/v1/network/connections/tcp" "" "200"
test_api "UDP连接列表" "GET" "/api/v1/network/connections/udp" "" "200"
test_api "进程连接获取" "GET" "/api/v1/network/connections/pid/1234" "" "200"
test_api "端口连接获取" "GET" "/api/v1/network/connections/port/80" "" "200"
test_api "IP连接获取" "GET" "/api/v1/network/connections/ip/192.168.1.100" "" "200"
test_api "关闭连接" "DELETE" "/api/v1/network/connection/12345" "" "200"
test_api "端口使用检查" "GET" "/api/v1/network/port/80/in-use" "" "200"
test_api "监听端口列表" "GET" "/api/v1/network/listening-ports" "" "200"
test_api "已建立连接" "GET" "/api/v1/network/established-connections" "" "200"
test_api "网络接口列表" "GET" "/api/v1/network/interfaces" "" "200"
test_api "网络统计信息" "GET" "/api/v1/network/stats" "" "200"
test_api "启用网络监控" "POST" "/api/v1/network/monitoring/enable" "" "200"
test_api "禁用网络监控" "POST" "/api/v1/network/monitoring/disable" "" "200"
test_api "监控连接列表" "GET" "/api/v1/network/monitoring/connections" "" "200"
test_api "连接历史记录" "GET" "/api/v1/network/monitoring/history" "" "200"

# ==================== 安全功能测试 ====================
echo "🛡️ 安全功能测试"
echo "==================="

test_api "安全状态获取" "GET" "/api/v1/security/status" "" "200"
test_api "规则信息获取" "GET" "/api/v1/security/rules" "" "200"
test_api "重新加载规则" "POST" "/api/v1/security/reload-rules" "" "200"
test_api "清除缓存" "POST" "/api/v1/security/cache/clear" "" "200"
test_api "缓存统计信息" "GET" "/api/v1/security/cache/stats" "" "200"
test_api "隔离文件" "POST" "/api/v1/security/quarantine" '{"file_path": "/tmp/malware.exe", "reason": "Malware detected"}' "200"
test_api "恢复文件" "POST" "/api/v1/security/restore" '{"file_path": "/tmp/malware.exe"}' "200"
test_api "隔离列表获取" "GET" "/api/v1/security/quarantine/list" "" "200"
test_api "扫描历史记录" "GET" "/api/v1/security/scan-history" "" "200"

# ==================== 用户管理功能测试 ====================
echo "👥 用户管理功能测试"
echo "==================="

test_api "当前用户信息" "GET" "/api/v1/user/current" "" "200"
test_api "用户ID查询" "GET" "/api/v1/user/id/1000" "" "200"
test_api "用户名查询" "GET" "/api/v1/user/name/testuser" "" "200"
test_api "所有用户列表" "GET" "/api/v1/user/all" "" "200"
test_api "权限检查" "POST" "/api/v1/user/permissions/check" '{"username": "testuser", "permissions": ["file_scan", "process_control"]}' "200"
test_api "密码验证" "POST" "/api/v1/user/password/validate" '{"username": "testuser", "password": "TestPassword123!"}' "200"
test_api "密码策略获取" "GET" "/api/v1/user/password/policy" "" "200"
test_api "账户状态检查" "GET" "/api/v1/user/status/testuser" "" "200"
test_api "用户会话列表" "GET" "/api/v1/user/sessions/testuser" "" "200"
test_api "终止用户会话" "DELETE" "/api/v1/user/session/session_12345" "" "200"
test_api "锁定用户账户" "POST" "/api/v1/user/lock/testuser" '{"reason": "Suspicious activity"}' "200"
test_api "解锁用户账户" "POST" "/api/v1/user/unlock/testuser" "" "200"
test_api "更改用户密码" "POST" "/api/v1/user/password/change" '{"username": "testuser", "old_password": "OldPassword123!", "new_password": "NewPassword456!"}' "200"
test_api "用户组列表" "GET" "/api/v1/user/groups/testuser" "" "200"
test_api "登录历史记录" "GET" "/api/v1/user/history/testuser" "" "200"

# ==================== 错误处理测试 ====================
echo "⚠️ 错误处理测试"
echo "==================="

test_api "无效文件路径" "POST" "/api/v1/file/scan" '{"path": "/nonexistent/file"}' "404"
test_api "无效进程ID" "GET" "/api/v1/process/99999" "" "404"
test_api "无效注册表路径" "GET" "/api/v1/registry/key/INVALID\\PATH" "" "404"
test_api "无效网络连接ID" "DELETE" "/api/v1/network/connection/invalid_id" "" "404"
test_api "无效用户ID" "GET" "/api/v1/user/id/99999" "" "404"

# ==================== 性能测试 ====================
echo "⚡ 性能测试"
echo "==================="

test_api "大量文件列表" "GET" "/api/v1/file/list?path=/usr&limit=1000" "" "200"
test_api "大量进程列表" "GET" "/api/v1/process/list?limit=1000" "" "200"
test_api "大量网络连接" "GET" "/api/v1/network/connections?limit=1000" "" "200"

# ==================== 测试结果统计 ====================
echo ""
echo "==================== 测试结果统计 ===================="
echo "📊 总测试数: $TOTAL_TESTS"
echo "✅ 通过测试: $PASSED_TESTS"
echo "❌ 失败测试: $FAILED_TESTS"

if [ $TOTAL_TESTS -gt 0 ]; then
    SUCCESS_RATE=$((PASSED_TESTS * 100 / TOTAL_TESTS))
    echo "📈 成功率: ${SUCCESS_RATE}%"
fi

if [ $FAILED_TESTS -gt 0 ]; then
    echo -e "${YELLOW}⚠️  有 $FAILED_TESTS 个测试失败，请检查API服务状态${NC}"
    exit 1
else
    echo -e "${GREEN}🎉 所有测试通过！${NC}"
    exit 0
fi 