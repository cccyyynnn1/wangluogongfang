#!/bin/bash

# 运行测试脚本
echo "🧪 开始运行Yara安全服务测试..."

# 颜色定义
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# 设置环境变量
export GO_ENV=test
export CGO_ENABLED=1
export GIN_MODE=test

# 检查Go是否安装
if ! command -v go &> /dev/null; then
    echo -e "${RED}错误: Go未安装${NC}"
    exit 1
fi

# 检查Go版本
GO_VERSION=$(go version | awk '{print $3}' | sed 's/go//')
echo -e "${BLUE}Go版本: $GO_VERSION${NC}"

# 检查依赖
echo -e "${BLUE}检查依赖...${NC}"
go mod tidy

# 创建测试目录
mkdir -p test-results
mkdir -p coverage

# 运行单元测试
echo -e "${BLUE}运行单元测试...${NC}"
go test -v -race -coverprofile=coverage/unit.out ./pkg/... 2>&1 | tee test-results/unit-tests.log
UNIT_TEST_RESULT=$?

# 运行集成测试
echo -e "${BLUE}运行集成测试...${NC}"
go test -v -race -coverprofile=coverage/integration.out ./internal/... 2>&1 | tee test-results/integration-tests.log
INTEGRATION_TEST_RESULT=$?

# 运行API测试
echo -e "${BLUE}运行API测试...${NC}"
go test -v -race -coverprofile=coverage/api.out ./cmd/... 2>&1 | tee test-results/api-tests.log
API_TEST_RESULT=$?

# 运行所有测试
echo -e "${BLUE}运行所有测试...${NC}"
go test -v -race -coverprofile=coverage/all.out ./... 2>&1 | tee test-results/all-tests.log
ALL_TEST_RESULT=$?

# 生成测试覆盖率报告
echo -e "${BLUE}生成测试覆盖率报告...${NC}"
go tool cover -html=coverage/all.out -o coverage/coverage.html
go tool cover -func=coverage/all.out > coverage/coverage.txt

# 显示覆盖率统计
echo -e "${BLUE}测试覆盖率统计:${NC}"
cat coverage/coverage.txt | tail -1

# 检查测试结果
TOTAL_FAILURES=0
if [ $UNIT_TEST_RESULT -ne 0 ]; then
    echo -e "${RED}❌ 单元测试失败${NC}"
    TOTAL_FAILURES=$((TOTAL_FAILURES + 1))
fi

if [ $INTEGRATION_TEST_RESULT -ne 0 ]; then
    echo -e "${RED}❌ 集成测试失败${NC}"
    TOTAL_FAILURES=$((TOTAL_FAILURES + 1))
fi

if [ $API_TEST_RESULT -ne 0 ]; then
    echo -e "${RED}❌ API测试失败${NC}"
    TOTAL_FAILURES=$((TOTAL_FAILURES + 1))
fi

if [ $ALL_TEST_RESULT -ne 0 ]; then
    echo -e "${RED}❌ 部分测试失败${NC}"
    TOTAL_FAILURES=$((TOTAL_FAILURES + 1))
fi

# 生成测试报告摘要
echo -e "${BLUE}生成测试报告摘要...${NC}"
cat > test-results/summary.txt << EOF
Yara安全服务测试报告
====================
测试时间: $(date)
Go版本: $GO_VERSION

测试结果:
- 单元测试: $([ $UNIT_TEST_RESULT -eq 0 ] && echo "通过" || echo "失败")
- 集成测试: $([ $INTEGRATION_TEST_RESULT -eq 0 ] && echo "通过" || echo "失败")
- API测试: $([ $API_TEST_RESULT -eq 0 ] && echo "通过" || echo "失败")
- 总体测试: $([ $ALL_TEST_RESULT -eq 0 ] && echo "通过" || echo "失败")

失败测试数: $TOTAL_FAILURES

测试文件位置:
- 单元测试日志: test-results/unit-tests.log
- 集成测试日志: test-results/integration-tests.log
- API测试日志: test-results/api-tests.log
- 所有测试日志: test-results/all-tests.log
- 覆盖率报告: coverage/coverage.html
- 覆盖率统计: coverage/coverage.txt
EOF

# 显示测试报告摘要
echo ""
echo -e "${BLUE}==================== 测试报告摘要 ====================${NC}"
cat test-results/summary.txt

# 最终结果
if [ $TOTAL_FAILURES -eq 0 ]; then
    echo -e "${GREEN}🎉 所有测试通过！${NC}"
    echo -e "${GREEN}📊 覆盖率报告: coverage/coverage.html${NC}"
    exit 0
else
    echo -e "${RED}❌ 有 $TOTAL_FAILURES 个测试类别失败${NC}"
    echo -e "${YELLOW}请检查测试日志文件获取详细信息${NC}"
    exit 1
fi 