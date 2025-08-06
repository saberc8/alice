# Alice 项目启动说明

## ⚡ 快速启动

### 1. 环境准备
确保已安装以下环境：
- Node.js 20.x
- Go 1.19+
- pnpm 10.8.0+

### 2. 安装依赖
```bash
# 根目录
npm install

# 前端依赖
cd frontend && pnpm install

# 后端依赖
cd ../backend && go mod tidy
```

### 3. 启动项目

#### 🎯 推荐方式：一键启动
```bash
./start.sh          # macOS/Linux
start.bat           # Windows
```

#### 🔧 其他方式
```bash
npm run dev         # 使用 npm 脚本
```

### 4. 访问地址
- 前端：http://localhost:5173
- 后端：http://localhost:8090

## 🛑 停止服务
按 `Ctrl + C` 停止所有服务

## ❓ 遇到问题？
查看详细文档：[STARTUP.md](./STARTUP.md)

---
*Alice 项目 - 全栈开发快速启动*
