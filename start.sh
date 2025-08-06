#!/bin/bash

# Alice 项目启动脚本
# 同时启动后端 Go 服务和前端开发服务器

# 颜色定义
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# 打印带颜色的消息
print_message() {
    echo -e "${2}[$(date '+%Y-%m-%d %H:%M:%S')] $1${NC}"
}

# 检查依赖
check_dependencies() {
    print_message "检查依赖..." $YELLOW
    
    # 检查 Go
    if ! command -v go &> /dev/null; then
        print_message "错误: Go 未安装" $RED
        exit 1
    fi
    
    # 检查 pnpm
    if ! command -v pnpm &> /dev/null; then
        print_message "错误: pnpm 未安装" $RED
        exit 1
    fi
    
    print_message "依赖检查完成" $GREEN
}

# 清理函数
cleanup() {
    print_message "正在停止服务..." $YELLOW
    kill 0
    exit
}

# 设置信号处理
trap cleanup SIGINT SIGTERM

# 获取脚本所在目录
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
BACKEND_DIR="$SCRIPT_DIR/backend"
FRONTEND_DIR="$SCRIPT_DIR/frontend"

print_message "Alice 项目启动脚本" $BLUE
print_message "项目路径: $SCRIPT_DIR" $BLUE

# 检查目录是否存在
if [ ! -d "$BACKEND_DIR" ]; then
    print_message "错误: backend 目录不存在" $RED
    exit 1
fi

if [ ! -d "$FRONTEND_DIR" ]; then
    print_message "错误: frontend 目录不存在" $RED
    exit 1
fi

# 检查依赖
check_dependencies

# 启动后端服务
print_message "启动后端服务..." $GREEN
cd "$BACKEND_DIR"
go run main.go &
BACKEND_PID=$!

# 等待一下让后端启动
sleep 2

# 启动前端服务
print_message "启动前端开发服务器..." $GREEN
cd "$FRONTEND_DIR"
pnpm run dev &
FRONTEND_PID=$!

print_message "服务启动完成!" $GREEN
print_message "后端服务 PID: $BACKEND_PID" $BLUE
print_message "前端服务 PID: $FRONTEND_PID" $BLUE
print_message "按 Ctrl+C 停止所有服务" $YELLOW

# 等待进程
wait
