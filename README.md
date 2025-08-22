# Alice

这是 Alice 项目的仓库，一个包含后端服务、管理后台（前端）和 Flutter 客户端的全栈示例工程。

## 项目简介

Alice 旨在作为一个可运行的示例系统，演示用户管理、RBAC（基于角色的权限控制）、动态菜单、时刻（moment）功能、聊天（WebSocket）等常见后端与前端交互场景。仓库包含：
- `backend/`：使用 Go 开发的后端服务（业务层、领域层、仓储、路由与中间件）。
- `admin/`：基于 Vue 3 + Vite 的管理后台，负责管理员界面、RBAC 管理、MinIO 管理等功能。
- `client_flutter/`：移动端/跨平台客户端，基于 Flutter，包含应用 UI 和与后端交互的功能。

## 目录概览

仓库顶层结构（摘要）：
- `backend/`：Go 服务源码与可执行文件，DDD 领域驱动设计。
	- `api/`, `domain/`, `infra/`, `application/` 等分层组织。
- `admin/`：前端源码（Vite + Vue）。
- `client_flutter/`：Flutter 客户端源码。
- `docs/`：接口与架构说明。

（详细目录请参考仓库中的各子目录。）

## 核心模块说明

- 后端（`backend/`）
	- 技术栈：Go
	- 主要职责：提供 REST / WebSocket API、用户/权限/菜单管理、时刻与聊天服务、存储（MinIO）对接与数据库访问。
	- 代码组织：按领域与分层（domain / infra / application / api / router / middleware）划分，便于维护与单元测试。

- 管理后台（`admin/`）
	- 技术栈：Vue 3, Vite, TypeScript
	- 主要职责：管理员登录、RBAC 管理（菜单、角色、权限、用户）、MinIO 浏览器等。

- Flutter 客户端（`client_flutter/`）
	- 技术栈：Flutter
	- 主要职责：移动端用户界面，调用后端 API，包含基础页面与模块化特性。

## 快速开始（本地开发）

以下命令在 macOS / zsh 环境下执行。请先确保已安装 Go、Node.js（推荐 pnpm）、Flutter，以及数据库/MinIO 等依赖（若需要）。
1) 后端（开发模式）

	- 进入目录并拉取依赖：
	```bash
	cd backend
	go mod download
		```

	- 运行服务（示例）：
	```bash
	# 直接运行
	go run main.go

	# 或构建并运行二进制
	go build -o bin/alice ./...
	./bin/alice

	- 重要文件：`backend/config.yaml`、`backend/docs/swagger.yaml`（API 参考）
2) 管理后台（开发模式）

	- 进入 `admin` 并安装依赖（仓库使用 pnpm/ npm）：
	```bash
	cd admin
	# 使用 pnpm（若已安装）
	pnpm install
	# 或使用 npm
	npm install

	- 启动开发服务器：
		```bash
		# pnpm
	pnpm run dev
	# 或 npm
	npm run dev

3) Flutter 客户端

	- 获取依赖并运行：
	```bash
	cd client_flutter
	flutter pub get
		flutter run

	- 若想只安装依赖（无运行）：
	```bash
	flutter pub get
	```

## 常用开发任务

- 运行后端单元/集成测试：`cd backend && go test ./...`
- 生成或查看 API 文档：`backend/docs/swagger.yaml` 或通过项目提供的文档路由（若已启动服务）。
- 构建前端生产包：`cd admin && pnpm build` 或 `npm run build`。

## 配置与环境

- 后端配置：`backend/config.yaml`，包含数据库、MinIO、JWT 等设置。启动前请根据本地环境调整。
- 前端：`admin/.env*`（见仓库内的 `.env.development`, `.env.production` 等）。

## 贡献与维护

- 问题与功能请求请在仓库 Issue 中提交。PR 请遵循项目的代码风格和测试覆盖要求。

## 参考资料

- 仓库内部 `docs/` 中包含架构文档与 API 示例（`docs/API_EXAMPLES.md` 等）。
