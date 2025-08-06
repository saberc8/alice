# RBAC权限管理系统

本项目实现了一个完整的基于角色的访问控制(RBAC)权限管理系统，支持用户、角色、权限和菜单的精确控制，可以精确到页面按钮级别的权限。

## 系统架构

### 领域模型

#### 1. 用户（User）
- 用户基本信息
- 与角色多对多关联

#### 2. 角色（Role）
- 角色名称和代码
- 角色状态（激活/非激活）
- 与用户多对多关联
- 与权限多对多关联
- 与菜单多对多关联

#### 3. 权限（Permission）
- 权限名称和代码
- 资源和操作（resource:action格式）
- 权限状态

#### 4. 菜单（Menu）
- 菜单类型：分组(GROUP)、目录(CATALOGUE)、菜单(MENU)、按钮(BUTTON)
- 层级结构（父子关系）
- 菜单元数据（图标、组件路径等）
- 排序和状态

#### 5. 关联关系
- 用户角色关联（UserRole）
- 角色权限关联（RolePermission）
- 角色菜单关联（RoleMenu）

### 项目结构

```
backend/
├── domain/
│   ├── user/                   # 用户领域
│   └── rbac/                   # RBAC领域
│       ├── entity/             # 实体层
│       │   ├── role.go
│       │   ├── permission.go
│       │   ├── menu.go
│       │   └── relation.go
│       ├── repository/         # 仓储接口层
│       │   ├── role_repository.go
│       │   ├── permission_repository.go
│       │   └── menu_repository.go
│       └── service/            # 服务层
│           ├── role_service.go
│           ├── permission_service.go
│           └── menu_service.go
├── infra/
│   └── repository/             # 仓储实现层
│       ├── role_repository_impl.go
│       ├── permission_repository_impl.go
│       └── menu_repository_impl.go
├── api/
│   ├── handler/                # API处理器
│   │   ├── role_handler.go
│   │   ├── permission_handler.go
│   │   └── menu_handler.go
│   ├── middleware/
│   │   └── permission.go       # 权限中间件
│   └── router/
│       └── rbac_router.go      # RBAC路由
└── cmd/
    └── init/
        └── main.go             # 初始化数据脚本
```

## API接口

### 角色管理

| 方法 | 路径 | 说明 |
|------|------|------|
| POST | `/api/v1/roles` | 创建角色 |
| GET | `/api/v1/roles/:id` | 获取单个角色 |
| GET | `/api/v1/roles` | 获取角色列表 |
| PUT | `/api/v1/roles/:id` | 更新角色 |
| DELETE | `/api/v1/roles/:id` | 删除角色 |

### 用户角色管理

| 方法 | 路径 | 说明 |
|------|------|------|
| GET | `/api/v1/users/:user_id/roles` | 获取用户角色 |
| POST | `/api/v1/users/:user_id/roles` | 为用户分配角色 |
| DELETE | `/api/v1/users/:user_id/roles` | 移除用户角色 |

### 权限管理

| 方法 | 路径 | 说明 |
|------|------|------|
| POST | `/api/v1/permissions` | 创建权限 |
| GET | `/api/v1/permissions/:id` | 获取单个权限 |
| GET | `/api/v1/permissions` | 获取权限列表 |
| PUT | `/api/v1/permissions/:id` | 更新权限 |
| DELETE | `/api/v1/permissions/:id` | 删除权限 |

### 角色权限管理

| 方法 | 路径 | 说明 |
|------|------|------|
| GET | `/api/v1/roles/:role_id/permissions` | 获取角色权限 |
| POST | `/api/v1/roles/:role_id/permissions` | 为角色分配权限 |
| DELETE | `/api/v1/roles/:role_id/permissions` | 移除角色权限 |

### 用户权限查询

| 方法 | 路径 | 说明 |
|------|------|------|
| GET | `/api/v1/users/:user_id/permissions` | 获取用户权限 |
| GET | `/api/v1/users/:user_id/permissions/check` | 检查用户权限 |

### 菜单管理

| 方法 | 路径 | 说明 |
|------|------|------|
| POST | `/api/v1/menus` | 创建菜单 |
| GET | `/api/v1/menus/:id` | 获取单个菜单 |
| GET | `/api/v1/menus` | 获取菜单列表 |
| GET | `/api/v1/menus/tree` | 获取菜单树 |
| PUT | `/api/v1/menus/:id` | 更新菜单 |
| DELETE | `/api/v1/menus/:id` | 删除菜单 |

### 角色菜单管理

| 方法 | 路径 | 说明 |
|------|------|------|
| GET | `/api/v1/roles/:role_id/menus` | 获取角色菜单 |
| POST | `/api/v1/roles/:role_id/menus` | 为角色分配菜单 |
| DELETE | `/api/v1/roles/:role_id/menus` | 移除角色菜单 |

### 用户菜单查询

| 方法 | 路径 | 说明 |
|------|------|------|
| GET | `/api/v1/users/:user_id/menus` | 获取用户菜单 |
| GET | `/api/v1/users/:user_id/menus/tree` | 获取用户菜单树 |

## 使用示例

### 1. 创建角色

```bash
curl -X POST http://localhost:8080/api/v1/roles \
  -H "Content-Type: application/json" \
  -d '{
    "name": "管理员",
    "code": "admin",
    "description": "系统管理员",
    "status": "active"
  }'
```

### 2. 创建权限

```bash
curl -X POST http://localhost:8080/api/v1/permissions \
  -H "Content-Type: application/json" \
  -d '{
    "name": "查看用户",
    "code": "user:read",
    "resource": "user",
    "action": "read",
    "description": "查看用户信息的权限"
  }'
```

### 3. 为角色分配权限

```bash
curl -X POST http://localhost:8080/api/v1/roles/{role_id}/permissions \
  -H "Content-Type: application/json" \
  -d '{
    "permission_ids": ["permission_id_1", "permission_id_2"]
  }'
```

### 4. 为用户分配角色

```bash
curl -X POST http://localhost:8080/api/v1/users/{user_id}/roles \
  -H "Content-Type: application/json" \
  -d '{
    "role_ids": ["role_id_1", "role_id_2"]
  }'
```

### 5. 检查用户权限

```bash
curl "http://localhost:8080/api/v1/users/{user_id}/permissions/check?resource=user&action=read"
```

### 6. 创建菜单

```bash
curl -X POST http://localhost:8080/api/v1/menus \
  -H "Content-Type: application/json" \
  -d '{
    "name": "用户管理",
    "code": "user_management",
    "path": "/user",
    "type": 2,
    "order": 1,
    "meta": {
      "icon": "user",
      "component": "/pages/user/index"
    }
  }'
```

### 7. 获取用户菜单树

```bash
curl "http://localhost:8080/api/v1/users/{user_id}/menus/tree"
```

## 权限中间件使用

在需要权限控制的API上使用权限中间件：

```go
// 在路由中使用权限中间件
authenticated.GET("/users", 
    middleware.RequirePermission(permissionService, "user", "read"),
    userHandler.ListUsers)

authenticated.POST("/users", 
    middleware.RequirePermission(permissionService, "user", "create"),
    userHandler.CreateUser)
```

## 数据初始化

运行初始化脚本创建基础数据：

```bash
cd backend
go run cmd/init/main.go
```

这将创建：
- 默认角色（超级管理员、管理员、普通用户）
- 基础权限（用户、角色、权限、菜单的CRUD权限）
- 示例菜单结构

## 特性

1. **完整的RBAC模型**：支持用户、角色、权限的完整关联
2. **菜单权限控制**：支持页面级和按钮级权限控制
3. **层级菜单结构**：支持无限层级的菜单树
4. **灵活的权限检查**：支持资源:操作格式的细粒度权限控制
5. **RESTful API**：提供完整的RESTful API接口
6. **数据库自动迁移**：自动创建和维护数据库表结构
7. **初始化脚本**：提供数据初始化脚本
8. **权限中间件**：提供Gin中间件进行API权限控制

## 前端集成

前端可以通过以下方式集成RBAC系统：

1. 登录时获取用户的角色和权限信息
2. 根据用户权限动态生成菜单
3. 在页面组件中检查按钮级权限
4. 使用菜单树API构建导航结构

这个RBAC系统严格遵循了项目的架构设计，采用了清晰的分层架构，便于维护和扩展。
