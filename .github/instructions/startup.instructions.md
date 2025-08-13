---
applyTo: '**'
---
# Alice 项目启动指南

## 📋 项目概述

Alice 是一个全栈应用项目，包含：
- **后端**: Go 语言开发的 REST API 服务
- **前端**: React + TypeScript + Vite 开发的现代化 Web 应用

## 🏗️ 项目结构

```
alice/
├── backend/          # Go 后端服务
│   ├── main.go      # 主程序入口
│   ├── go.mod       # Go 模块文件
│   └── ...
├── admin/         #  前端管理后台应用
│   ├── package.json # 前端依赖配置
│   ├── src/         # 源代码目录
│   └── ...
├── start.sh         # Linux/macOS 启动脚本
├── start.bat        # Windows 启动脚本
└── package.json     # 根目录项目配置
```

## 🔧 环境要求

### 必需环境
- **Node.js**: 20.x （根据 `engines` 配置）
- **Go**: 1.19+ 
- **pnpm**: 10.8.0+ （推荐使用指定版本）

### 环境检查
```bash
# 检查 Node.js 版本
node --version

# 检查 Go 版本
go version

# 检查 pnpm 版本
pnpm --version
```

## 🚀 快速启动

### 方式 1: 一键启动脚本 (推荐)

#### Linux/macOS
```bash
# 给脚本添加执行权限（首次运行）
chmod +x start.sh

# 启动项目
./start.sh
```

#### Windows
```cmd
# 双击运行或在命令行执行
start.bat
```

### 方式 2: 使用 npm 脚本
```bash
# 安装根目录依赖（首次运行）
npm install

# 同时启动前后端
npm run dev
```

## 🌐 访问地址

启动成功后，可以通过以下地址访问：

- **前端应用**: http://localhost:8091 (Vite 默认端口)
- **后端 API**: http://localhost:8090 (Go 服务默认端口)

> 注意: 具体端口可能因配置而异，请查看终端输出的实际地址

现在你可以开始开发 Alice 项目了！🚀
