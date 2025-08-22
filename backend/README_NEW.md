# Alice - 企业级Go后端项目

Alice 是一个基于 Go 和 Gin 框架的企业级后端项目，采用 DDD（领域驱动设计）架构模式。项目提供了完整的用户认证系统，支持 JWT 认证、角色权限管理，并遵循现代软件工程的最佳实践。

## 🚀 特性

- **DDD 架构**: 清晰的分层架构，易于维护和扩展
- **用户认证**: 完整的注册、登录、JWT 认证系统
- **角色权限**: 基于 RBAC 的权限管理系统
- **数据库**: PostgreSQL + GORM ORM
- **安全性**: bcrypt 密码加密，JWT 令牌验证
- **配置管理**: YAML 配置文件支持
- **容器化**: Docker 和 Docker Compose 支持
- **API 文档**: Swagger/OpenAPI 文档
- **测试**: 单元测试和集成测试
- **CI/CD**: GitHub Actions 工作流

## 🏃 快速开始

### 环境要求

- Go 1.21+
- PostgreSQL 13+
- Make

## 📁 项目结构

```
alice/
├── README.md                   # 项目说明文档
├── LICENSE                     # 开源许可证
├── Dockerfile                  # Docker 构建文件
├── docker-compose.yml          # Docker Compose 配置
├── Makefile                    # 构建和管理脚本
├── .gitignore                  # Git 忽略文件
├── .github/                    # GitHub 配置
│   ├── workflows/              # CI/CD 工作流
│   │   ├── ci.yml             # 持续集成
│   │   └── release.yml        # 发布流程
│   ├── ISSUE_TEMPLATE/         # Issue 模板
│   ├── PULL_REQUEST_TEMPLATE.md # PR 模板
│   └── CODEOWNERS             # 代码所有者
├── docs/                       # 项目文档
│   ├── api/                   # API 文档
│   ├── architecture/          # 架构文档
│   └── development/           # 开发文档
├── scripts/                    # 脚本文件
├── configs/                    # 配置文件模板
├── go.mod                      # Go 模块定义
├── go.sum                      # 依赖版本锁定
├── main.go                     # 应用入口
├── application/                # 应用层
│   └── application.go         # 应用初始化和依赖注入
├── api/                        # API 层（接口层）
│   ├── handler/               # HTTP 处理器
│   ├── middleware/            # 中间件
│   ├── model/                 # 请求/响应模型
│   └── router/                # 路由配置
├── domain/                     # 领域层（核心业务逻辑）
│   ├── user/                  # 用户领域
│   │   ├── entity/            # 用户实体
│   │   ├── repository/        # 用户仓储接口
│   │   └── service/           # 用户领域服务
│   └── role/                  # 角色领域
│       ├── entity/            # 角色实体
│       ├── repository/        # 角色仓储接口
│       └── service/           # 角色领域服务
├── infra/                      # 基础设施层
│   ├── config/                # 配置管理
│   ├── database/              # 数据库连接和迁移
│   ├── repository/            # 仓储实现
│   └── cache/                 # 缓存实现
├── pkg/                        # 通用工具包
│   ├── logger/                # 日志工具
│   ├── validator/             # 数据验证
│   ├── errors/                # 错误处理
│   └── utils/                 # 工具函数
└── test/                       # 测试文件
    ├── integration/           # 集成测试
    ├── fixtures/              # 测试数据
    └── mocks/                 # 模拟对象
```

## 🏗 架构设计

### DDD 分层架构

Alice 项目采用 DDD（领域驱动设计）的分层架构：

#### 1. API 层 (api/)
- **职责**: 处理 HTTP 请求，数据验证，响应格式化
- **组件**: Handler, Middleware, Router, Model
- **原则**: 不包含业务逻辑，只负责协议转换

#### 2. 应用层 (application/)
- **职责**: 用例编排，事务管理，领域服务协调
- **组件**: Application Service, DTO, Use Cases
- **原则**: 不包含业务规则，只负责流程控制

#### 3. 领域层 (domain/)
- **职责**: 核心业务逻辑，业务规则，领域模型
- **组件**: Entity, Value Object, Domain Service, Repository Interface
- **原则**: 独立于技术实现，包含核心业务逻辑

#### 4. 基础设施层 (infra/)
- **职责**: 技术实现，外部系统集成，数据持久化
- **组件**: Repository Implementation, Database, Cache, Config
- **原则**: 实现领域层定义的接口

### 依赖关系

```
API Layer     →     Application Layer
     ↓                      ↓
Infrastructure ←  Domain Layer
```

## 📝 开发规范

### 新增功能规范

以新增 Role 接口为例，需要遵循以下步骤：

#### 1. 领域层开发 (domain/role/)

**实体定义** (`entity/role.go`):
```go
package entity

type Role struct {
    ID          uint   `json:"id" gorm:"primaryKey"`
    Name        string `json:"name" gorm:"uniqueIndex;not null"`
    Description string `json:"description"`
    Permissions []Permission `json:"permissions" gorm:"many2many:role_permissions;"`
    CreatedAt   time.Time `json:"created_at"`
    UpdatedAt   time.Time `json:"updated_at"`
}
```

**仓储接口** (`repository/role_repository.go`):
```go
package repository

type RoleRepository interface {
    Create(role *entity.Role) error
    GetByID(id uint) (*entity.Role, error)
    GetByName(name string) (*entity.Role, error)
    Update(role *entity.Role) error
    Delete(id uint) error
    List(offset, limit int) ([]*entity.Role, int64, error)
}
```

**领域服务** (`service/role_service.go`):
```go
package service

type RoleService interface {
    CreateRole(name, description string) (*entity.Role, error)
    GetRoleByID(id uint) (*entity.Role, error)
    UpdateRole(id uint, updates map[string]interface{}) error
    DeleteRole(id uint) error
    AssignPermissions(roleID uint, permissionIDs []uint) error
}
```

#### 2. 基础设施层实现 (infra/repository/)

**仓储实现** (`role_repository_impl.go`):
```go
package repository

type roleRepositoryImpl struct {
    db *gorm.DB
}

func NewRoleRepository(db *gorm.DB) repository.RoleRepository {
    return &roleRepositoryImpl{db: db}
}
```

#### 3. API 层实现 (api/)

**请求模型** (`model/role_request.go`):
```go
package model

type CreateRoleRequest struct {
    Name        string `json:"name" binding:"required,min=2,max=50"`
    Description string `json:"description" binding:"max=200"`
}
```

**处理器** (`handler/role_handler.go`):
```go
package handler

type RoleHandler struct {
    roleService service.RoleService
}

func (h *RoleHandler) CreateRole(c *gin.Context) {
    // 实现创建角色逻辑
}
```

**路由注册** (`router/router.go`):
```go
func (r *Router) SetupRoutes() *gin.Engine {
    // 添加角色相关路由
    roleGroup := v1.Group("/roles")
    roleGroup.Use(middleware.JWTAuth())
    {
        roleGroup.POST("", r.roleHandler.CreateRole)
        roleGroup.GET("/:id", r.roleHandler.GetRole)
        // ...
    }
}
```

#### 4. 应用层集成 (application/)

**依赖注入** (`application.go`):
```go
var (
    RoleSvc service.RoleService
)

func Init(ctx context.Context, cfg *config.Config) error {
    // 初始化角色相关服务
    roleRepo := repository.NewRoleRepository(db)
    RoleSvc = service.NewRoleService(roleRepo)
}
```

### 编码规范

#### 命名规范
- **包名**: 小写，简短，有意义
- **文件名**: 蛇形命名 `user_service.go`
- **接口**: 大写开头，以接口功能命名 `UserService`
- **结构体**: 大写开头，驼峰命名 `UserHandler`
- **方法**: 大写开头（公开），小写开头（私有）
- **常量**: 全大写，下划线分隔 `USER_STATUS_ACTIVE`

#### 错误处理
```go
var (
    ErrRoleNotFound      = errors.New("role not found")
    ErrRoleAlreadyExists = errors.New("role already exists")
    ErrInvalidPermission = errors.New("invalid permission")
)
```

#### 日志记录
```go
logger.Infof("Creating role: %s", roleName)
logger.Errorf("Failed to create role: %v", err)
```

## 📚 API 文档

API 文档使用 Swagger/OpenAPI 3.0 规范，访问地址：
- 开发环境: http://localhost:8090/swagger/index.html
- 生产环境: (依据部署域名) https://<your-domain>/swagger/index.html

### 生成 / 更新文档

项目使用 [swag](https://github.com/swaggo/swag) 从代码注释生成文档，步骤如下：

```bash
cd backend
make swagger  # 等价于: go install github.com/swaggo/swag/cmd/swag@latest && swag init -g main.go -o docs
```

生成后会在 `backend/docs` 目录出现 `docs.go swagger.json swagger.yaml` 文件。

### 注释示例

如下是 `Register` 接口的注释样例：

```go
// @Summary 用户注册
// @Description 注册新用户
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body model.RegisterRequest true "注册请求"
// @Success 200 {object} model.APIResponse{data=model.RegisterResponse}
// @Failure 400 {object} model.APIResponse
// @Router /auth/register [post]
```

更多注释可参考 `api/handler/user_handler.go`。

### 访问控制

当前示例项目默认所有环境均暴露 `/swagger/*any`。生产环境可在路由中增加开关，例如通过环境变量：

```go
if cfg.Server.EnableSwagger {
    router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
}
```

并在配置中添加 `enableSwagger` 开关。

### 主要 API 端点

#### 认证相关
- `POST /api/v1/auth/register` - 用户注册
- `POST /api/v1/auth/login` - 用户登录
- `POST /api/v1/auth/refresh` - 刷新令牌

#### 用户管理
- `GET /api/v1/users/profile` - 获取用户资料
- `PUT /api/v1/users/profile` - 更新用户资料
- `GET /api/v1/users` - 用户列表（管理员）

#### 角色管理
- `POST /api/v1/roles` - 创建角色
- `GET /api/v1/roles` - 角色列表
- `GET /api/v1/roles/:id` - 获取角色详情
- `PUT /api/v1/roles/:id` - 更新角色
- `DELETE /api/v1/roles/:id` - 删除角色


### 环境变量

| 变量名 | 描述 | 默认值 |
|--------|------|--------|
| `SERVER_PORT` | 服务端口 | `:8090` |
| `DB_HOST` | 数据库主机 | `localhost` |
| `DB_PORT` | 数据库端口 | `5432` |
| `DB_USERNAME` | 数据库用户名 | `postgres` |
| `DB_PASSWORD` | 数据库密码 | - |
| `DB_NAME` | 数据库名 | `alice` |
| `JWT_SECRET` | JWT 密钥 | - |
| `ENABLE_SWAGGER` | 是否启用 swagger 文档 | `true` |
