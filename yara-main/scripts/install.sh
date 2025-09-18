#!/bin/bash

# Yara安全服务安装脚本

set -e

# 颜色定义
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# 日志函数
log_info() {
    echo -e "${GREEN}[INFO]${NC} $1"
}

log_warn() {
    echo -e "${YELLOW}[WARN]${NC} $1"
}

log_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

log_step() {
    echo -e "${BLUE}[STEP]${NC} $1"
}

# 检查系统要求
check_requirements() {
    log_step "检查系统要求..."
    
    # 检查Go版本
    if ! command -v go &> /dev/null; then
        log_error "Go未安装，请先安装Go 1.21或更高版本"
        echo "请访问 https://golang.org/dl/ 安装Go"
        exit 1
    fi
    
    GO_VERSION=$(go version | awk '{print $3}' | sed 's/go//')
    REQUIRED_VERSION="1.21"
    
    if [ "$(printf '%s\n' "$REQUIRED_VERSION" "$GO_VERSION" | sort -V | head -n1)" != "$REQUIRED_VERSION" ]; then
        log_error "Go版本过低，需要1.21或更高版本，当前版本: $GO_VERSION"
        exit 1
    fi
    
    log_info "Go版本检查通过: $GO_VERSION"
    
    # 检查curl
    if ! command -v curl &> /dev/null; then
        log_error "curl未安装，请先安装curl"
        exit 1
    fi
    
    log_info "系统要求检查通过"
}

# 创建目录结构
create_directories() {
    log_step "创建目录结构..."
    
    mkdir -p bin
    mkdir -p logs
    mkdir -p rules
    mkdir -p configs
    mkdir -p docs
    mkdir -p scripts
    mkdir -p temp
    mkdir -p quarantine
    mkdir -p test-results
    mkdir -p coverage
    
    log_info "目录结构创建完成"
}

# 下载依赖
download_dependencies() {
    log_step "下载Go依赖..."
    
    go mod tidy
    if [ $? -ne 0 ]; then
        log_error "依赖下载失败"
        exit 1
    fi
    
    log_info "依赖下载完成"
}

# 构建项目
build_project() {
    log_step "构建项目..."
    
    # 清理旧的构建文件
    rm -f bin/yara-security-service
    
    # 构建项目
    go build -ldflags="-s -w" -o bin/yara-security-service cmd/server/main.go
    
    if [ $? -ne 0 ]; then
        log_error "构建失败"
        exit 1
    fi
    
    # 设置可执行权限
    chmod +x bin/yara-security-service
    
    log_info "项目构建完成"
}

# 创建配置文件
create_config() {
    log_step "创建配置文件..."
    
    if [ ! -f "configs/config.yaml" ]; then
        cat > configs/config.yaml << EOF
server:
  port: 8081
  host: "0.0.0.0"
  read_timeout: 30s
  write_timeout: 30s

logging:
  level: "info"
  format: "json"
  output: "stdout"

security:
  yara_rules_path: "./rules"
  max_file_size: 100MB
  scan_timeout: 60s
  enable_real_time: true
  api_key_required: false
  api_key: ""
  allowed_ips: []

file:
  max_scan_depth: 10
  exclude_patterns:
    - "*.tmp"
    - "*.log"
    - "*.cache"
  include_extensions:
    - ".exe"
    - ".dll"
    - ".sys"
    - ".bat"
    - ".cmd"
    - ".ps1"
    - ".vbs"
    - ".js"

process:
  max_process_count: 1000
  refresh_interval: 5s

registry:
  max_key_length: 256
  max_value_size: 1MB

network:
  connection_timeout: 30s
  max_connections: 1000

cors:
  allowed_origins:
    - "*"
  allowed_methods:
    - "GET"
    - "POST"
    - "PUT"
    - "DELETE"
  allowed_headers:
    - "Content-Type"
    - "Authorization"

rate_limit:
  requests_per_minute: 60
  burst_size: 10
  window_size: 60s
EOF
        log_info "配置文件创建完成"
    else
        log_info "配置文件已存在"
    fi
}

# 创建示例规则文件
create_sample_rules() {
    log_step "创建示例规则文件..."
    
    if [ ! -f "rules/malware.yar" ]; then
        cat > rules/malware.yar << EOF
/*
   Yara规则示例文件
   包含常见的恶意软件检测规则
*/

rule Malware_Generic {
    meta:
        description = "通用恶意软件检测规则"
        author = "Yara Security Service"
        date = "2024-01-01"
        version = "1.0"
        severity = "high"
        category = "malware"
    
    strings:
        $mz_header = "MZ"
        $pe_header = "PE"
        $suspicious_string1 = "malware" nocase
        $suspicious_string2 = "trojan" nocase
        $suspicious_string3 = "virus" nocase
        $suspicious_string4 = "backdoor" nocase
    
    condition:
        $mz_header at 0 and
        $pe_header and
        any of ($suspicious_string*)
}

rule Trojan_Generic {
    meta:
        description = "特洛伊木马检测规则"
        author = "Yara Security Service"
        date = "2024-01-01"
        version = "1.0"
        severity = "critical"
        category = "trojan"
    
    strings:
        $mz_header = "MZ"
        $pe_header = "PE"
        $trojan_string1 = "cmd.exe" nocase
        $trojan_string2 = "powershell" nocase
        $trojan_string3 = "netcat" nocase
        $trojan_string4 = "nc.exe" nocase
    
    condition:
        $mz_header at 0 and
        $pe_header and
        any of ($trojan_string*)
}

rule Shellcode_Generic {
    meta:
        description = "Shellcode检测规则"
        author = "Yara Security Service"
        date = "2024-01-01"
        version = "1.0"
        severity = "critical"
        category = "shellcode"
    
    strings:
        $shellcode_pattern1 = { 90 90 90 90 90 90 90 90 }
        $shellcode_pattern2 = { 68 ?? ?? ?? ?? 68 ?? ?? ?? ?? 68 ?? ?? ?? ?? }
        $shellcode_pattern3 = { 31 C0 50 68 ?? ?? ?? ?? 68 ?? ?? ?? ?? 68 ?? ?? ?? ?? }
    
    condition:
        any of ($shellcode_pattern*)
}
EOF
        log_info "示例规则文件创建完成"
    else
        log_info "规则文件已存在"
    fi
}

# 创建服务文件
create_service_file() {
    log_step "创建系统服务文件..."
    
    if [[ "$OSTYPE" == "linux-gnu"* ]]; then
        if command -v systemctl &> /dev/null; then
            sudo tee /etc/systemd/system/yara-security.service > /dev/null << EOF
[Unit]
Description=Yara Security Service
After=network.target

[Service]
Type=simple
User=$USER
WorkingDirectory=$(pwd)
ExecStart=$(pwd)/bin/yara-security-service
Restart=always
RestartSec=10
Environment=GO_ENV=production
Environment=CGO_ENABLED=1

[Install]
WantedBy=multi-user.target
EOF
            log_info "系统服务文件创建完成"
            log_info "使用以下命令管理服务:"
            log_info "  sudo systemctl enable yara-security"
            log_info "  sudo systemctl start yara-security"
            log_info "  sudo systemctl status yara-security"
        else
            log_warn "未检测到systemctl，跳过服务文件创建"
        fi
    else
        log_warn "非Linux系统，跳过服务文件创建"
    fi
}

# 运行测试
run_tests() {
    log_step "运行测试..."
    
    if [ -f "scripts/run_tests.sh" ]; then
        chmod +x scripts/run_tests.sh
        ./scripts/run_tests.sh
        if [ $? -ne 0 ]; then
            log_warn "部分测试失败，但安装继续进行"
        else
            log_info "所有测试通过"
        fi
    else
        log_warn "测试脚本不存在，跳过测试"
    fi
}

# 显示安装完成信息
show_completion_info() {
    log_step "安装完成！"
    
    echo ""
    echo -e "${GREEN}🎉 Yara安全服务安装完成！${NC}"
    echo ""
    echo -e "${BLUE}📁 项目结构:${NC}"
    echo "  ├── bin/yara-security-service    # 可执行文件"
    echo "  ├── configs/config.yaml         # 配置文件"
    echo "  ├── rules/malware.yar           # 规则文件"
    echo "  ├── logs/                       # 日志目录"
    echo "  ├── temp/                       # 临时文件目录"
    echo "  └── quarantine/                 # 隔离文件目录"
    echo ""
    echo -e "${BLUE}🚀 启动服务:${NC}"
    echo "  ./scripts/start.sh              # 启动服务"
    echo "  ./scripts/start_web.sh          # 启动Web界面"
    echo "  ./scripts/api_test.sh           # 运行API测试"
    echo ""
    echo -e "${BLUE}📊 服务地址:${NC}"
    echo "  API服务: http://localhost:8081"
    echo "  Web界面: http://localhost:3001"
    echo "  健康检查: http://localhost:8081/api/health"
    echo ""
    echo -e "${BLUE}📋 管理命令:${NC}"
    echo "  ./bin/yara-security-service     # 直接运行"
    echo "  ./scripts/run_tests.sh          # 运行测试"
    echo "  curl http://localhost:8081/api/health  # 健康检查"
    echo ""
    echo -e "${YELLOW}⚠️  注意事项:${NC}"
    echo "  - 首次运行前请检查配置文件 configs/config.yaml"
    echo "  - 确保防火墙允许端口 8081 和 3001"
    echo "  - 生产环境建议启用API密钥认证"
    echo ""
}

# 主函数
main() {
    echo -e "${BLUE}🚀 Yara安全服务安装脚本${NC}"
    echo "=================================="
    
    check_requirements
    create_directories
    download_dependencies
    build_project
    create_config
    create_sample_rules
    create_service_file
    run_tests
    show_completion_info
}

# 运行主函数
main "$@" 