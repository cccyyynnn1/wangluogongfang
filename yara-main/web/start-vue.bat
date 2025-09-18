@echo off
chcp 65001 >nul

echo 🚀 启动Yara安全服务Vue界面...

REM 检查Node.js是否安装
node --version >nul 2>&1
if %errorlevel% neq 0 (
    echo ❌ Node.js未安装，请先安装Node.js
    pause
    exit /b 1
)

REM 检查npm是否安装
npm --version >nul 2>&1
if %errorlevel% neq 0 (
    echo ❌ npm未安装，请先安装npm
    pause
    exit /b 1
)

REM 检查是否在正确的目录
if not exist "package.json" (
    echo ❌ 请在web目录下运行此脚本
    pause
    exit /b 1
)

REM 检查node_modules是否存在
if not exist "node_modules" (
    echo 📦 安装依赖...
    npm install
    if %errorlevel% neq 0 (
        echo ❌ 依赖安装失败
        pause
        exit /b 1
    )
)

REM 检查后端服务是否运行
echo 🔍 检查后端服务状态...
curl -s http://localhost:8081/api/health >nul 2>&1
if %errorlevel% equ 0 (
    echo ✅ 后端服务运行正常
) else (
    echo ⚠️  后端服务未运行，请先启动后端服务
    echo 💡 提示: 在项目根目录运行 scripts\start.bat
)

REM 启动开发服务器
echo 🌐 启动Vue开发服务器...
echo 📊 服务将在 http://localhost:3001 启动
echo 🔗 API代理到 http://localhost:8081
echo.
echo 按 Ctrl+C 停止服务
echo.

npm run dev 