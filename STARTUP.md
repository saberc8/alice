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
├── frontend/         # React 前端应用
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

### 方式 3: 分别启动
```bash
# 终端 1: 启动后端
npm run dev:backend
# 或者
cd backend && go run main.go

# 终端 2: 启动前端
npm run dev:frontend
# 或者
cd frontend && pnpm run dev
```

## 📦 依赖安装

### 首次设置项目
```bash
# 1. 安装根目录依赖
npm install

# 2. 安装前端依赖
cd frontend
pnpm install

# 3. 安装后端依赖（Go 模块）
cd ../backend
go mod tidy
```

## 🌐 访问地址

启动成功后，可以通过以下地址访问：

- **前端应用**: http://localhost:5173 (Vite 默认端口)
- **后端 API**: http://localhost:8090 (Go 服务默认端口)

> 注意: 具体端口可能因配置而异，请查看终端输出的实际地址

## 🛠️ 可用脚本命令

### 根目录脚本
```bash
npm run dev              # 同时启动前后端
npm run dev:backend      # 仅启动后端
npm run dev:frontend     # 仅启动前端
npm run build           # 构建前端
npm run build:backend   # 构建后端
npm run install:frontend # 安装前端依赖
```

### 前端脚本 (frontend/)
```bash
pnpm run dev      # 启动开发服务器
pnpm run build    # 构建生产版本
pnpm run preview  # 预览构建结果
```

### 后端脚本 (backend/)
```bash
go run main.go    # 运行开发服务器
go build          # 构建可执行文件
make run          # 使用 Makefile 运行
make build        # 使用 Makefile 构建
```

## 🔍 故障排除

### 常见问题

#### 1. 端口占用
```bash
# 查看端口占用
lsof -i :5173  # 前端端口
lsof -i :8090  # 后端端口

# 杀死占用进程
kill -9 <PID>
```

#### 2. 依赖安装失败
```bash
# 清除缓存并重新安装
cd frontend
rm -rf node_modules pnpm-lock.yaml
pnpm install

# Go 模块问题
cd backend
go clean -modcache
go mod download
```

#### 3. Go 环境问题
```bash
# 检查 Go 环境变量
go env GOPATH
go env GOROOT

# 设置 Go 代理（中国用户）
go env -w GOPROXY=https://goproxy.cn,direct
```

#### 4. Node.js 版本不匹配
```bash
# 使用 nvm 管理 Node.js 版本
nvm install 20
nvm use 20
```

### 错误日志查看

#### 查看启动脚本输出
启动脚本会显示详细的状态信息，包括：
- ✅ 依赖检查结果
- ✅ 服务启动状态
- ✅ 进程 PID 信息
- ✅ 错误信息（如果有）

#### 查看服务日志
- **后端日志**: 在后端终端窗口查看
- **前端日志**: 在前端终端窗口查看
- **浏览器控制台**: F12 开发者工具查看前端错误

## 🔄 停止服务

### 使用启动脚本启动的服务
- 在脚本运行的终端按 `Ctrl + C`
- 脚本会自动清理所有相关进程

### 手动启动的服务
- 在各自的终端窗口按 `Ctrl + C`
- 或者关闭对应的终端窗口

## 📝 开发说明

### 代码热重载
- **前端**: Vite 提供热模块替换 (HMR)，代码修改后自动刷新
- **后端**: 需要手动重启，或使用 `air` 等工具实现热重载

### 开发工具推荐
- **VS Code**: 推荐的编辑器
- **Go 扩展**: Go 语言支持
- **ES7+ React 扩展**: React 开发支持
- **Thunder Client**: API 测试工具

## 📞 获取帮助

如果遇到问题，可以：
1. 检查本文档的故障排除部分
2. 查看项目的 `README.md` 文件
3. 检查 `backend/README.md` 和 `frontend/README.md`
4. 查看项目 Issues 或提交新的 Issue

---

## 🎉 快速开始示例

```bash
# 克隆项目后的完整启动流程
git clone <repository-url>
cd alice

# 安装依赖
npm install
cd frontend && pnpm install && cd ..
cd backend && go mod tidy && cd ..

# 启动项目
./start.sh

# 等待服务启动完成，然后访问
# 前端: http://localhost:5173
# 后端: http://localhost:8090
```

现在你可以开始开发 Alice 项目了！🚀
