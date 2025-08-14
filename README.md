# Alice 全栈项目

<div align="center">

![Alice Logo](https://via.placeholder.com/200x80/4A90E2/FFFFFF?text=Alice)

*一个现代化的全栈 Web 应用*

[![Go](https://img.shields.io/badge/Go-1.19+-00ADD8?style=flat&logo=go)](https://golang.org/)
[![Vue](https://img.shields.io/badge/Vue-19+-61DAFB?style=flat&logo=Vue)](https://Vuejs.org/)
[![TypeScript](https://img.shields.io/badge/TypeScript-5.6+-3178C6?style=flat&logo=typescript)](https://www.typescriptlang.org/)
[![License](https://img.shields.io/badge/License-MIT-green.svg)](LICENSE)

</div>

## 📖 项目简介

Alice 是一个基于现代技术栈的全栈 Web 应用，采用前后端分离架构：

- **前端**: Vue + TypeScript + Vite
- **后端**: Go + Gin + GORM + PostgreSQL
- **架构**: DDD (领域驱动设计)

## 🚀 快速开始

### 一键启动
```bash
# Linux/macOS
./start.sh

# Windows
start.bat

# 或使用 npm 脚本
npm run dev
```

> 📋 **详细启动说明**: [STARTUP.md](./STARTUP.md) | [快速启动](./README_STARTUP.md)

### 访问地址
- 🌐 **前端应用**: http://localhost:8091
- 🔌 **后端 API**: http://localhost:8090

## 🏗️ 项目结构

```
alice/
├── 📁 backend/          # Go 后端服务
│   ├── main.go         # 程序入口
│   ├── api/            # API 层
│   ├── application/    # 应用层
│   ├── domain/         # 领域层
│   └── infra/          # 基础设施层
├── 📁 admin/        # Vue 前端应用
│   ├── src/            # 源代码
│   ├── public/         # 静态资源
│   └── package.json    # 依赖配置
├── 📁 docs/            # 项目文档
├── 🚀 start.sh         # 启动脚本 (Linux/macOS)
├── 🚀 start.bat        # 启动脚本 (Windows)
└── 📄 STARTUP.md       # 详细启动说明
```

## ✨ 主要特性

### 后端特性
- 🏛️ **DDD 架构**: 领域驱动设计，清晰的分层架构
- 🔐 **用户认证**: JWT Token 认证机制
- 🗄️ **数据库**: PostgreSQL + GORM ORM
- 🔒 **安全加密**: bcrypt 密码加密
- 📝 **RESTful API**: 标准的 REST API 设计
- ⚡ **高性能**: Gin 框架，高并发处理能力

### 前端特性
- ⚛️ **Vue3**: 最新的 Vue 特性
- 🎯 **TypeScript**: 类型安全的 JavaScript
- ⚡ **Vite**: 快速的构建工具
## 🔧 环境要求

- **Node.js**: 20.x
- **Go**: 1.19+
- **pnpm**: 10.8.0+
- **PostgreSQL**: 13+ (可选，可使用 SQLite 开发)

## 📚 文档导航

- 📋 [详细启动说明](./STARTUP.md)
- ⚡ [快速启动指南](./README_STARTUP.md)
- 🏗️ [后端文档](./backend/README.md)
- 🎨 [前端文档](./frontend/README.md)
- 📐 [架构设计](./docs/architecture.md)

## 🛠️ 开发指令

## 🤝 贡献指南

1. Fork 本仓库
2. 创建功能分支 (`git checkout -b feature/AmazingFeature`)
3. 提交更改 (`git commit -m 'Add some AmazingFeature'`)
4. 推送到分支 (`git push origin feature/AmazingFeature`)
5. 开启 Pull Request

请遵循项目的代码规范：
- **后端**: 遵循 Go 官方代码规范
- **前端**: 使用 ESLint + Prettier + Biome

## 📄 许可证

本项目基于 MIT 许可证开源 - 查看 [LICENSE](LICENSE) 文件了解详情

## 👥 团队

- **Backend**: Go + DDD 架构
- **Frontend**: Vue + TypeScript
- **DevOps**: Docker + 自动化部署

## 🔗 相关链接


<div align="center">

**Alice Project** - 构建现代化的全栈应用 🚀

</div>

模型：agent+gpt5或者agent+claude4
第一步：先让ai设计go的后端基础架构
第二步：让ai设计一个rbac管理的api以及表设计，生成md格式的文档，在刚刚的ddd架构基础之上开始开发
第三步：实现后台管理系统的前端的登录快速启动模块
第四步：实现前端的动态菜单渲染
第五步：实现第一个user的table的增删改查
第六步：批量实现菜单，角色，权限管理
