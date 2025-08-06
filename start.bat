@echo off
chcp 65001 >nul
setlocal EnableDelayedExpansion

:: Alice 项目启动脚本 (Windows)
:: 同时启动后端 Go 服务和前端开发服务器

echo ========================================
echo Alice 项目启动脚本
echo ========================================

:: 获取脚本所在目录
set "SCRIPT_DIR=%~dp0"
set "BACKEND_DIR=%SCRIPT_DIR%backend"
set "FRONTEND_DIR=%SCRIPT_DIR%frontend"

echo 项目路径: %SCRIPT_DIR%

:: 检查目录是否存在
if not exist "%BACKEND_DIR%" (
    echo 错误: backend 目录不存在
    pause
    exit /b 1
)

if not exist "%FRONTEND_DIR%" (
    echo 错误: frontend 目录不存在
    pause
    exit /b 1
)

:: 检查依赖
echo 检查依赖...

:: 检查 Go
where go >nul 2>&1
if errorlevel 1 (
    echo 错误: Go 未安装
    pause
    exit /b 1
)

:: 检查 pnpm
where pnpm >nul 2>&1
if errorlevel 1 (
    echo 错误: pnpm 未安装
    pause
    exit /b 1
)

echo 依赖检查完成

:: 启动后端服务
echo 启动后端服务...
cd /d "%BACKEND_DIR%"
start "Alice Backend" cmd /k "go run main.go"

:: 等待一下让后端启动
timeout /t 3 /nobreak >nul

:: 启动前端服务
echo 启动前端开发服务器...
cd /d "%FRONTEND_DIR%"
start "Alice Frontend" cmd /k "pnpm run dev"

echo.
echo 服务启动完成!
echo 后端和前端服务已在新窗口中启动
echo 关闭对应的命令行窗口即可停止服务
echo.
pause
