#!/bin/bash

echo "正在构建 Yara 安全服务..."

# 清理之前的构建
if [ -f "bin/yara-security-service" ]; then
    rm bin/yara-security-service
fi

# 构建服务器
echo "构建服务器..."
go build -o bin/yara-security-service cmd/server/main.go

if [ $? -eq 0 ]; then
    echo "构建成功！服务器文件位置: bin/yara-security-service"
    echo "文件大小:"
    ls -lh bin/yara-security-service
else
    echo "构建失败！"
    exit 1
fi

echo ""
echo "构建完成！"
echo "运行服务器: ./bin/yara-security-service"
