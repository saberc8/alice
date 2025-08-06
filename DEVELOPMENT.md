# 开发环境配置指南

## 📋 环境检查清单

在开始开发之前，请确保以下环境已正确安装：

### ✅ 必需环境

#### 1. Node.js 20.x
```bash
# 检查版本
node --version
# 应显示: v20.x.x

# 如果版本不匹配，建议使用 nvm 管理
# macOS/Linux
curl -o- https://raw.githubusercontent.com/nvm-sh/nvm/v0.39.0/install.sh | bash
nvm install 20
nvm use 20

# Windows
# 下载并安装 Node.js 20.x from https://nodejs.org/
```

#### 2. Go 1.19+
```bash
# 检查版本
go version
# 应显示: go version go1.19.x 或更高

# 安装 Go (如果未安装)
# macOS
brew install go

# Linux
wget https://golang.org/dl/go1.21.0.linux-amd64.tar.gz
sudo tar -C /usr/local -xzf go1.21.0.linux-amd64.tar.gz

# Windows
# 下载并安装 Go from https://golang.org/dl/
```

#### 3. pnpm 10.8.0+
```bash
# 检查版本
pnpm --version
# 应显示: 10.8.0 或更高

# 安装 pnpm
npm install -g pnpm@latest

# 或者使用官方安装脚本
curl -fsSL https://get.pnpm.io/install.sh | sh -
```

### 🛠️ 推荐工具

#### Git
```bash
# 检查版本
git --version

# 配置 Git (首次使用)
git config --global user.name "你的姓名"
git config --global user.email "你的邮箱"
```

#### 数据库 (可选)
```bash
# PostgreSQL (生产环境推荐)
# macOS
brew install postgresql
brew services start postgresql

# 或者使用 Docker
docker run --name alice-postgres -e POSTGRES_PASSWORD=password -d -p 5432:5432 postgres:13
```

## 🔧 IDE 配置

### VS Code (推荐)

#### 必需扩展
```json
{
  "recommendations": [
    "golang.go",                    // Go 语言支持
    "bradlc.vscode-tailwindcss",   // TailwindCSS 智能提示
    "esbenp.prettier-vscode",      // 代码格式化
    "ms-vscode.vscode-typescript-next", // TypeScript 支持
    "biomejs.biome"                // Biome 代码检查
  ]
}
```

#### 工作区设置
创建 `.vscode/settings.json`:
```json
{
  "go.gopath": "",
  "go.goroot": "",
  "go.toolsManagement.checkForUpdates": "local",
  "editor.formatOnSave": true,
  "editor.codeActionsOnSave": {
    "source.organizeImports": true
  },
  "typescript.preferences.useAliasesForRenames": false,
  "emmet.includeLanguages": {
    "typescript": "html",
    "typescriptreact": "html"
  }
}
```

## 🌍 环境变量配置

### 后端环境变量
创建 `backend/.env`:
```env
# 数据库配置
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=password
DB_NAME=alice
DB_SSLMODE=disable

# JWT 配置
JWT_SECRET=your-super-secret-jwt-key
JWT_EXPIRES_IN=24h

# 服务器配置
SERVER_PORT=8090
SERVER_MODE=debug

# 日志配置
LOG_LEVEL=debug
LOG_FORMAT=json
```

### 前端环境变量
创建 `frontend/.env.local`:
```env
# API 基础地址
VITE_API_BASE_URL=http://localhost:8090

# 应用配置
VITE_APP_TITLE=Alice Admin
VITE_APP_VERSION=1.0.0

# 开发模式配置
VITE_DEV_TOOLS=true
```

## 📦 依赖管理

### 安装所有依赖
```bash
# 在项目根目录执行
make install-deps

# 或者手动安装
npm install                    # 根目录依赖
cd frontend && pnpm install   # 前端依赖
cd ../backend && go mod tidy  # 后端依赖
```

### Go 模块配置
```bash
# 设置 Go 代理 (中国用户推荐)
go env -w GOPROXY=https://goproxy.cn,direct
go env -w GOSUMDB=sum.golang.google.cn

# 启用模块模式
go env -w GO111MODULE=on
```

## 🔥 热重载配置

### 前端热重载
前端已内置 Vite HMR，无需额外配置。

### 后端热重载 (可选)
安装 Air 进行 Go 热重载：
```bash
# 安装 Air
go install github.com/cosmtrek/air@latest

# 在 backend 目录创建 .air.toml
# (配置文件内容略，可参考 Air 官方文档)

# 使用 Air 启动
cd backend
air
```

## 🐳 Docker 开发环境 (可选)

### Docker Compose 配置
创建 `docker-compose.dev.yml`:
```yaml
version: '3.8'
services:
  postgres:
    image: postgres:13
    environment:
      POSTGRES_DB: alice
      POSTGRES_USER: postgres
      POSTGRES_PASSWORD: password
    ports:
      - "5432:5432"
    volumes:
      - postgres_data:/var/lib/postgresql/data

  redis:
    image: redis:7-alpine
    ports:
      - "6379:6379"

volumes:
  postgres_data:
```

启动开发环境：
```bash
docker-compose -f docker-compose.dev.yml up -d
```

## 🧪 代码质量工具

### 前端代码检查
```bash
cd frontend

# Biome 检查
pnpm run check

# 修复自动修复的问题
pnpm run check --apply

# TypeScript 类型检查
pnpm run type-check
```

### 后端代码检查
```bash
cd backend

# 格式化代码
go fmt ./...

# 代码检查
go vet ./...

# 使用 golangci-lint (需要安装)
golangci-lint run
```

## 🚨 常见问题解决

### 1. 端口冲突
```bash
# 查看端口占用
lsof -i :8090  # 后端端口
lsof -i :5173  # 前端端口

# 修改端口 (在相应的配置文件中)
```

### 2. pnpm 安装失败
```bash
# 清除缓存
pnpm store prune

# 重新安装
rm -rf node_modules pnpm-lock.yaml
pnpm install
```

### 3. Go 模块下载失败
```bash
# 清理模块缓存
go clean -modcache

# 重新下载
go mod download
```

## ✅ 验证安装

运行以下命令验证环境配置是否正确：

```bash
# 检查所有环境
./scripts/check-env.sh

# 或者手动检查
node --version && go version && pnpm --version
```

如果所有命令都正常输出版本号，说明环境配置完成！

现在可以运行 `./start.sh` 启动项目了。
