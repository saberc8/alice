## Alice - 全栈示例：Go 后端 + Vue3 管理端 + Flutter 客户端

一个采用 DDD 分层的企业级示例项目，包含：
- 后端：Go + Gin + GORM + PostgreSQL，内置用户/角色/权限（RBAC）、菜单、文件存储（MinIO，可选）、WebSocket 即时通信与 Swagger 文档
- 管理端：Vue 3 + TypeScript + Vite 的后台管理界面
- 客户端：Flutter 移动/Web 多端示例（登录、好友、聊天等）


### 技术栈一览
- 后端：Go 1.23+（module: `alice`）、Gin、GORM(PostgreSQL)、JWT、Swagger(swag)、WebSocket、MinIO SDK、YAML 配置
- 管理端：Vue3、Vite、TypeScript、Pinia、Vue Router
- 客户端：Flutter 3.x（iOS/Android/Web/macOS/Windows）


## 目录结构
```
alice/
├─ backend/             # Go 后端（DDD：api / application / domain / infra / pkg）
│  ├─ main.go           # 应用入口（加载配置、初始化依赖、注册路由/Swagger）
│  ├─ config.yaml       # 后端配置（服务、数据库、JWT、MinIO 等）
│  ├─ Makefile          # 构建与开发脚本（build/run/test/swagger/init-data 等）
│  ├─ docs/             # Swagger 产物与后端文档
│  ├─ api/              # Handler / Middleware / Model / Router
│  ├─ application/      # 应用编排与依赖注入
│  ├─ domain/           # 领域模型：rbac、user、chat、moment 等
│  ├─ infra/            # 配置、数据库、仓储实现、对象存储等
│  └─ pkg/              # 日志等通用包
├─ admin/               # Vue3 + Vite 管理端
├─ client_flutter/      # Flutter 客户端
└─ docs/                # 跨端使用说明（前端请求流、WS、样例数据等）
```


## 快速开始

> macOS + zsh 环境；命令均在项目根或对应子目录执行。

### 1) 后端（Go）
前置依赖：
- Go 1.23+（建议 1.23/1.24），Make，Git
- PostgreSQL 13+（可用 Docker 快速启动）
- 可选：本地 MinIO（对象存储）

步骤：
1. 启动数据库（二选一）
	- 使用 Make（Docker 运行 PostgreSQL 13）：
	  ```bash
	  cd backend
	  make db-setup   # 拉起 postgres:13（默认端口 5432，用户 postgres，密码 password，库 alice）
	  ```
	- 或使用自有 PostgreSQL，并在 `backend/config.yaml` 中填写连接信息。
2. 配置后端
	- 编辑 `backend/config.yaml`，至少确认以下关键项：
	  - server.port: 8090
	  - database: host/port/username/password/dbname/sslmode
	  - jwt.secret_key、jwt.expires_in（小时）
	  - minio（可选）：endpoint、access-key、secret-key、base-url
	- 注意：不要将真实秘钥提交到版本库，生产环境请改用环境变量或安全的配置托管。
3. 启动服务
	```bash
	cd backend
	make deps
	make run          # 首次会先 build，再运行，监听 0.0.0.0:8090
	```
4. 初始化 RBAC 菜单/权限（可选，但推荐）
	```bash
	cd backend
	make init-data
	```
5. Swagger 文档
	- 如在路由中启用 Swagger，访问：http://localhost:8090/swagger/index.html（BasePath: /api/v1）
	- 若需重新生成：`cd backend && make swagger`

常见接口分组（示例，具体以路由为准）：
- 认证：`POST /api/v1/auth/register`、`POST /api/v1/auth/login`、`POST /api/v1/auth/refresh`
- 用户：`GET/PUT /api/v1/users/profile`、`GET /api/v1/users`
- 角色/权限/菜单（RBAC）：`/api/v1/roles`、`/api/v1/permissions`、`/api/v1/menus`、`/api/v1/menus/tree`
- 即时通信（App）：WebSocket `GET /api/v1/app/chat/ws`（Header 或 query 携带 Bearer token），历史 `GET /api/v1/app/chat/history/{peer_id}`


### 2) 管理端（Vue3 + Vite）
前置依赖：Node.js 18+，建议使用 pnpm。

```bash
cd admin
pnpm install
pnpm dev
```

配置后端地址：新建 `.env.local`（示例）
```
VITE_API_BASE=http://localhost:8090/api/v1
```

如出现 CORS 问题，请在后端开启/放宽跨域或通过本地代理解决。


### 3) 客户端（Flutter）
前置依赖：Flutter 3.x SDK 与对应平台工具链。

```bash
cd client_flutter
flutter pub get              # 或使用 VS Code 任务 “flutter_pub_get”
flutter run -d <device>      # 选择 iOS/Android/Web/macos/windows 其一
```

将后端 API 基地址配置到 Flutter 工程（通常在 `lib/core` 的配置/常量文件中，按项目实际路径修改为 `http://localhost:8090/api/v1`）。


## WebSocket 聊天（简要）
- 连接：`ws://localhost:8090/api/v1/app/chat/ws?token=<JWT>` 或在 Header 使用 `Authorization: Bearer <JWT>`
- 发送示例：
  ```json
  { "type": "text", "to": 1024, "content": "hello" }
  ```
- 历史：`GET /api/v1/app/chat/history/{peer_id}?page=1&page_size=20`
更多细节见 `docs/ws.md`。


## 常用命令（后端 Makefile 摘要）
- 基础：`make deps`、`make build`、`make run`、`make dev`、`make clean`
- 文档：`make swagger`（自动生成至 `backend/docs/`）
- 测试：`make test`、`make test-coverage`（生成 `coverage.html`）
- RBAC 初始化：`make init-data`
- Docker：`make docker-build`、`make docker-run`
- 数据库（Docker）：`make db-setup`、`make db-start`、`make db-stop`、`make db-remove`


## 开发规范与架构
- 架构：DDD 分层（api / application / domain / infra），领域模型与仓储接口解耦，基础设施实现落在 infra
- 日志：`pkg/logger` 统一输出
- 配置：YAML +（可扩展环境变量覆盖）
参考：`backend/docs/architecture.md`、`backend/README_NEW.md`、`docs/frontend-request-flow.md`。


## 常见问题（FAQ）
- 端口占用：修改 `backend/config.yaml` 的 `server.port`，或关闭占用进程
- 数据库连不通：检查 Docker 容器是否启动、密码是否与配置一致；本地装有 PostgreSQL 时避免端口冲突
- Swagger 404：需确认路由中已注册 Swagger；并用 `make swagger` 重新生成文档
- CORS 跨域：本地开发建议使用代理或在后端中开启 CORS 中间件
- MinIO 访问：确保 `minio.base-url` 正确，并设置可访问的 endpoint/凭据（生产勿使用明文秘钥）


## 参考与文档
- 后端文档与 Swagger：`backend/docs/`
- 前端请求流与接口示例：`docs/frontend-request-flow.md`、`docs/login.json`
- WebSocket 使用说明：`docs/ws.md`


## 许可
若未提供许可证文件，请在公司/团队内部遵循合规策略；对外开源前，请补充 LICENSE。
