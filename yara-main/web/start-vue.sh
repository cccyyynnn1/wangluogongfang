#!/bin/bash

# Yara安全服务Vue界面启动脚本

# 颜色定义
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

echo -e "${BLUE}🚀 启动Yara安全服务Vue界面...${NC}"

# 检查Node.js是否安装
if ! command -v node &> /dev/null; then
    echo -e "${RED}❌ Node.js未安装，请先安装Node.js${NC}"
    exit 1
fi

# 检查npm是否安装
if ! command -v npm &> /dev/null; then
    echo -e "${RED}❌ npm未安装，请先安装npm${NC}"
    exit 1
fi

# 检查是否在正确的目录
if [ ! -f "package.json" ]; then
    echo -e "${RED}❌ 请在web目录下运行此脚本${NC}"
    exit 1
fi

# 检查node_modules是否存在
if [ ! -d "node_modules" ]; then
    echo -e "${YELLOW}📦 安装依赖...${NC}"
    npm install
    if [ $? -ne 0 ]; then
        echo -e "${RED}❌ 依赖安装失败${NC}"
        exit 1
    fi
fi

# 检查后端服务是否运行
echo -e "${YELLOW}🔍 检查后端服务状态...${NC}"
if curl -s http://localhost:8081/api/health > /dev/null 2>&1; then
    echo -e "${GREEN}✅ 后端服务运行正常${NC}"
else
    echo -e "${YELLOW}⚠️  后端服务未运行，请先启动后端服务${NC}"
    echo -e "${BLUE}💡 提示: 在项目根目录运行 ./scripts/start.sh${NC}"
fi

# 启动开发服务器
echo -e "${BLUE}🌐 启动Vue开发服务器...${NC}"
echo -e "${GREEN}📊 服务将在 http://localhost:3001 启动${NC}"
echo -e "${GREEN}🔗 API代理到 http://localhost:8081${NC}"
echo ""
echo -e "${YELLOW}按 Ctrl+C 停止服务${NC}"
echo ""

npm run dev 