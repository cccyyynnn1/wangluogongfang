@echo off
chcp 65001 >nul

echo 🚀 启动Yara安全服务Web界面...

REM 检查Go是否安装
go version >nul 2>&1
if errorlevel 1 (
    echo ❌ 错误: 未找到Go，请先安装Go
    echo 请访问 https://golang.org/dl/ 安装Go
    pause
    exit /b 1
)

REM 检查是否在主项目目录
if not exist "go.mod" (
    echo ❌ 错误: 请在项目根目录运行此脚本
    pause
    exit /b 1
)

REM 检查配置文件
if not exist "configs\config.yaml" (
    echo ❌ 错误: 配置文件不存在 configs\config.yaml
    pause
    exit /b 1
)

REM 创建必要的目录
echo 📁 创建必要的目录...
if not exist "web" mkdir web
if not exist "logs" mkdir logs
if not exist "temp" mkdir temp
if not exist "quarantine" mkdir quarantine

REM 检查依赖
echo 🔧 检查依赖...
go mod tidy
if errorlevel 1 (
    echo ❌ 错误: 依赖检查失败
    pause
    exit /b 1
)

REM 设置环境变量
set GO_ENV=production
set CGO_ENABLED=1
set GIN_MODE=release

echo 🚀 启动Yara安全服务...
start /B go run cmd/server/main.go > logs\yara-service-%date:~0,4%%date:~5,2%%date:~8,2%-%time:~0,2%%time:~3,2%%time:~6,2%.log 2>&1

REM 等待服务启动
echo ⏳ 等待服务启动...
timeout /t 5 /nobreak >nul

REM 检查服务是否启动成功
set RETRY_COUNT=0
set MAX_RETRIES=10

:check_service
curl -s http://localhost:8081/api/health >nul 2>&1
if errorlevel 1 (
    set /a RETRY_COUNT+=1
    if !RETRY_COUNT! lss !MAX_RETRIES! (
        echo ⏳ 等待服务启动... (尝试 !RETRY_COUNT!/!MAX_RETRIES!)
        timeout /t 2 /nobreak >nul
        goto check_service
    ) else (
        echo ❌ Yara安全服务启动失败
        echo 请检查日志文件 logs\yara-service-*.log
        pause
        exit /b 1
    )
) else (
    echo ✅ Yara安全服务启动成功
)

REM 检查Web服务器文件
if not exist "web\server.go" (
    echo ❌ 错误: Web服务器文件不存在 web\server.go
    pause
    exit /b 1
)

echo 🚀 启动Web界面服务器...
cd web
start /B go run server.go > ..\logs\web-server-%date:~0,4%%date:~5,2%%date:~8,2%-%time:~0,2%%time:~3,2%%time:~6,2%.log 2>&1

REM 等待Web服务器启动
timeout /t 3 /nobreak >nul

REM 检查Web服务器是否启动成功
curl -s http://localhost:3001 >nul 2>&1
if errorlevel 1 (
    echo ⚠️  Web界面服务器可能未完全启动，请稍等...
) else (
    echo ✅ Web界面服务器启动成功
)

echo.
echo 🎉 服务启动完成！
echo.
echo 📊 Yara安全服务: http://localhost:8081
echo 🌐 Web管理界面: http://localhost:3001
echo 📋 API文档: http://localhost:8081/api/health
echo.
echo 按任意键停止所有服务...

pause >nul

REM 停止服务
echo 🛑 正在停止服务...
taskkill /f /im go.exe >nul 2>&1
echo ✅ 服务已停止
pause 