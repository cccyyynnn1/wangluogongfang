#!/bin/bash

# Yara安全服务启动脚本

echo "🚀 启动Yara安全服务..."

# 颜色定义
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# 检查Go是否安装
if ! command -v go &> /dev/null; then
    echo -e "${RED}错误: Go未安装${NC}"
    echo "请访问 https://golang.org/dl/ 安装Go"
    exit 1
fi

# 检查Go版本
GO_VERSION=$(go version | awk '{print $3}' | sed 's/go//')
echo -e "${BLUE}Go版本: $GO_VERSION${NC}"

# 检查配置文件
if [ ! -f "configs/config.yaml" ]; then
    echo -e "${RED}错误: 配置文件不存在 configs/config.yaml${NC}"
    exit 1
fi

# 检查规则目录
if [ ! -d "rules" ]; then
    echo -e "${YELLOW}创建规则目录...${NC}"
    mkdir -p rules
fi

# 检查隔离目录
if [ ! -d "quarantine" ]; then
    echo -e "${YELLOW}创建隔离目录...${NC}"
    mkdir -p quarantine
fi

# 检查日志目录
if [ ! -d "logs" ]; then
    echo -e "${YELLOW}创建日志目录...${NC}"
    mkdir -p logs
fi

# 检查临时目录
if [ ! -d "temp" ]; then
    echo -e "${YELLOW}创建临时目录...${NC}"
    mkdir -p temp
fi

# 设置环境变量
export GO_ENV=production
export CGO_ENABLED=1
export GIN_MODE=release

# 检查依赖
echo -e "${BLUE}检查依赖...${NC}"
go mod tidy

# 检查依赖是否成功
if [ $? -ne 0 ]; then
    echo -e "${RED}错误: 依赖检查失败${NC}"
    exit 1
fi

# 清理旧的构建文件
echo -e "${BLUE}清理旧的构建文件...${NC}"
rm -f bin/yara-security-service

# 构建项目
echo -e "${BLUE}构建项目...${NC}"
go build -ldflags="-s -w" -o bin/yara-security-service cmd/server/main.go

# 检查构建是否成功
if [ $? -ne 0 ]; then
    echo -e "${RED}错误: 构建失败${NC}"
    exit 1
fi

# 检查构建文件是否存在
if [ ! -f "bin/yara-security-service" ]; then
    echo -e "${RED}错误: 构建文件不存在${NC}"
    exit 1
fi

# 设置可执行权限
chmod +x bin/yara-security-service

# 检查端口是否被占用
PORT=8081
if lsof -Pi :$PORT -sTCP:LISTEN -t >/dev/null ; then
    echo -e "${YELLOW}警告: 端口 $PORT 已被占用${NC}"
    echo -e "${YELLOW}请确保没有其他实例在运行${NC}"
fi

# 启动服务
echo -e "${BLUE}启动服务...${NC}"
echo -e "${GREEN}服务将在 http://localhost:$PORT 启动${NC}"
echo -e "${GREEN}API文档: http://localhost:$PORT/api/health${NC}"

# 启动后端服务
echo -e "${BLUE}启动后端服务...${NC}"
./bin/yara-security-service > logs/yara-service-$(date +%Y%m%d-%H%M%S).log 2>&1 &
BACKEND_PID=$!

# 等待后端服务启动
echo -e "${YELLOW}等待后端服务启动...${NC}"
sleep 5

# 检查后端服务是否启动成功
RETRY_COUNT=0
MAX_RETRIES=10
while [ $RETRY_COUNT -lt $MAX_RETRIES ]; do
    if curl -s http://localhost:8081/api/health > /dev/null 2>&1; then
        echo -e "${GREEN}✅ 后端服务启动成功${NC}"
        break
    else
        RETRY_COUNT=$((RETRY_COUNT + 1))
        echo -e "${YELLOW}等待后端服务启动... (尝试 $RETRY_COUNT/$MAX_RETRIES)${NC}"
        sleep 2
    fi
done

if [ $RETRY_COUNT -eq $MAX_RETRIES ]; then
    echo -e "${RED}❌ 后端服务启动失败${NC}"
    kill $BACKEND_PID 2>/dev/null
    echo -e "${YELLOW}请检查日志文件 logs/yara-service-*.log${NC}"
    exit 1
fi

# 检查Web服务器文件
if [ -f "web/server.go" ]; then
    echo -e "${BLUE}启动前端服务...${NC}"
    cd web
    go run server.go > ../logs/web-server-$(date +%Y%m%d-%H%M%S).log 2>&1 &
    FRONTEND_PID=$!
    cd ..
    
    # 等待前端服务启动
    sleep 3
    
    # 检查前端服务是否启动成功
    if curl -s http://localhost:3001 > /dev/null 2>&1; then
        echo -e "${GREEN}✅ 前端服务启动成功${NC}"
    else
        echo -e "${YELLOW}⚠️  前端服务可能未完全启动，请稍等...${NC}"
    fi
    
    echo ""
    echo -e "${GREEN}🎉 所有服务启动完成！${NC}"
    echo ""
    echo -e "${BLUE}📊 后端服务: http://localhost:8081${NC}"
    echo -e "${BLUE}🌐 前端服务: http://localhost:3001${NC}"
    echo -e "${BLUE}📋 API文档: http://localhost:8081/api/health${NC}"
    echo ""
    echo -e "${YELLOW}按 Ctrl+C 停止所有服务${NC}"
    
    # 等待中断信号
    trap 'echo ""; echo -e "${YELLOW}正在停止服务...${NC}"; kill $BACKEND_PID $FRONTEND_PID 2>/dev/null; echo -e "${GREEN}✅ 服务已停止${NC}"; exit 0' INT
    
    # 保持脚本运行
    wait
else
    echo -e "${YELLOW}前端服务文件不存在，仅启动后端服务${NC}"
    echo ""
    echo -e "${GREEN}🎉 后端服务启动完成！${NC}"
    echo ""
    echo -e "${BLUE}📊 后端服务: http://localhost:8081${NC}"
    echo -e "${BLUE}📋 API文档: http://localhost:8081/api/health${NC}"
    echo ""
    echo -e "${YELLOW}按 Ctrl+C 停止服务${NC}"
    
    # 等待中断信号
    trap 'echo ""; echo -e "${YELLOW}正在停止服务...${NC}"; kill $BACKEND_PID 2>/dev/null; echo -e "${GREEN}✅ 服务已停止${NC}"; exit 0' INT
    
    # 保持脚本运行
    wait
fi 