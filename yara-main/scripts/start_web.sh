#!/bin/bash

# Yara安全服务Web界面启动脚本

echo "🚀 启动Yara安全服务Web界面..."

# 颜色定义
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# 检查Go是否安装
if ! command -v go &> /dev/null; then
    echo -e "${RED}错误: 未找到Go，请先安装Go${NC}"
    exit 1
fi

# 检查是否在主项目目录
if [ ! -f "go.mod" ]; then
    echo -e "${RED}错误: 请在项目根目录运行此脚本${NC}"
    exit 1
fi

# 检查配置文件
if [ ! -f "configs/config.yaml" ]; then
    echo -e "${RED}错误: 配置文件不存在 configs/config.yaml${NC}"
    exit 1
fi

# 创建必要的目录
echo -e "${BLUE}创建必要的目录...${NC}"
mkdir -p web
mkdir -p logs
mkdir -p temp
mkdir -p quarantine

# 检查依赖
echo -e "${BLUE}检查依赖...${NC}"
go mod tidy

# 设置环境变量
export GO_ENV=production
export CGO_ENABLED=1
export GIN_MODE=release

# 启动Yara安全服务（后台运行）
echo -e "${BLUE}启动Yara安全服务...${NC}"
go run cmd/server/main.go > logs/yara-service-$(date +%Y%m%d-%H%M%S).log 2>&1 &
YARA_PID=$!

# 等待服务启动
echo -e "${YELLOW}等待服务启动...${NC}"
sleep 5

# 检查服务是否启动成功
RETRY_COUNT=0
MAX_RETRIES=10
while [ $RETRY_COUNT -lt $MAX_RETRIES ]; do
    if curl -s http://localhost:8081/api/health > /dev/null 2>&1; then
        echo -e "${GREEN}✅ Yara安全服务启动成功${NC}"
        break
    else
        RETRY_COUNT=$((RETRY_COUNT + 1))
        echo -e "${YELLOW}等待服务启动... (尝试 $RETRY_COUNT/$MAX_RETRIES)${NC}"
        sleep 2
    fi
done

if [ $RETRY_COUNT -eq $MAX_RETRIES ]; then
    echo -e "${RED}❌ Yara安全服务启动失败${NC}"
    kill $YARA_PID 2>/dev/null
    echo -e "${YELLOW}请检查日志文件 logs/yara-service-*.log${NC}"
    exit 1
fi

# 检查Web服务器文件
if [ ! -f "web/server.go" ]; then
    echo -e "${RED}错误: Web服务器文件不存在 web/server.go${NC}"
    kill $YARA_PID 2>/dev/null
    exit 1
fi

# 启动Web界面服务器
echo -e "${BLUE}启动Web界面服务器...${NC}"
cd web
go run server.go > ../logs/web-server-$(date +%Y%m%d-%H%M%S).log 2>&1 &
WEB_PID=$!

# 等待Web服务器启动
sleep 3

# 检查Web服务器是否启动成功
if curl -s http://localhost:3001 > /dev/null 2>&1; then
    echo -e "${GREEN}✅ Web界面服务器启动成功${NC}"
else
    echo -e "${YELLOW}⚠️  Web界面服务器可能未完全启动，请稍等...${NC}"
fi

echo ""
echo -e "${GREEN}🎉 服务启动完成！${NC}"
echo ""
echo -e "${BLUE}📊 Yara安全服务: http://localhost:8081${NC}"
echo -e "${BLUE}🌐 Web管理界面: http://localhost:3001${NC}"
echo -e "${BLUE}📋 API文档: http://localhost:8081/api/health${NC}"
echo ""
echo -e "${YELLOW}按 Ctrl+C 停止所有服务${NC}"

# 等待中断信号
trap 'echo ""; echo -e "${YELLOW}正在停止服务...${NC}"; kill $YARA_PID $WEB_PID 2>/dev/null; echo -e "${GREEN}✅ 服务已停止${NC}"; exit 0' INT

# 保持脚本运行
wait 