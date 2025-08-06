# Alice RBAC 权限管理系统 - 项目总结

## 项目概述

基于Go语言实现的完整RBAC(基于角色的访问控制)权限管理系统，支持用户、角色、权限和菜单的精确控制，可以精确到页面按钮级别的权限。严格遵循DDD(领域驱动设计)架构模式。

## 🎯 核心特性

### ✅ 完整的RBAC模型
- **用户(User)** - 系统使用者
- **角色(Role)** - 权限的集合
- **权限(Permission)** - 具体的操作许可(resource:action格式)
- **菜单(Menu)** - 支持分组、目录、菜单、按钮四种类型

### ✅ 灵活的权限控制
- **资源级权限**：基于resource:action格式的细粒度权限控制
- **菜单级权限**：支持页面访问控制
- **按钮级权限**：支持页面内功能按钮的权限控制
- **层级菜单**：支持无限层级的菜单树结构

### ✅ 完整的API接口
- **RESTful设计**：符合REST规范的API设计
- **统一响应格式**：标准化的API响应结构
- **错误处理**：完善的错误处理和响应机制
- **参数验证**：请求参数的验证和类型检查

### ✅ 中间件支持
- **权限中间件**：基于权限的API访问控制
- **认证中间件**：JWT认证支持
- **日志中间件**：请求日志记录
- **CORS中间件**：跨域资源共享支持

## 🏗️ 架构设计

### 分层架构
```
┌─────────────────────────────────────┐
│              API Layer              │  ← Handler + Router + Middleware
├─────────────────────────────────────┤
│           Application Layer         │  ← Service Registration
├─────────────────────────────────────┤
│            Domain Layer             │  ← Business Logic
│  ┌─────────────┐ ┌─────────────────┐ │
│  │    User     │ │      RBAC       │ │
│  │   Domain    │ │     Domain      │ │
│  └─────────────┘ └─────────────────┘ │
├─────────────────────────────────────┤
│         Infrastructure Layer        │  ← Repository Implementation
├─────────────────────────────────────┤
│            Database Layer           │  ← GORM + PostgreSQL
└─────────────────────────────────────┘
```

### 目录结构
```
backend/
├── domain/                    # 领域层
│   ├── user/                 # 用户领域
│   └── rbac/                 # RBAC领域
│       ├── entity/           # 实体
│       ├── repository/       # 仓储接口
│       └── service/          # 领域服务
├── infra/                    # 基础设施层
│   ├── config/              # 配置
│   ├── database/            # 数据库
│   └── repository/          # 仓储实现
├── api/                     # API层
│   ├── handler/            # 处理器
│   ├── middleware/         # 中间件
│   ├── model/             # API模型
│   └── router/            # 路由
├── application/            # 应用层
├── pkg/                   # 工具包
├── cmd/                   # 命令行工具
│   └── init/             # 数据初始化
└── docs/                 # 文档
```

## 📊 数据库设计

### 核心表结构

#### 用户表 (users)
- id, username, password_hash, email, status
- created_at, updated_at

#### 角色表 (roles)
- id, name, code, description, status
- created_at, updated_at

#### 权限表 (permissions)
- id, name, code, resource, action, description, status
- created_at, updated_at

#### 菜单表 (menus)
- id, parent_id, name, code, path, type, order, status, meta
- created_at, updated_at

#### 关联表
- user_roles (用户角色关联)
- role_permissions (角色权限关联)
- role_menus (角色菜单关联)

## 🔧 API接口

### 角色管理
- `POST /api/v1/roles` - 创建角色
- `GET /api/v1/roles/:id` - 获取角色
- `GET /api/v1/roles` - 角色列表
- `PUT /api/v1/roles/:id` - 更新角色
- `DELETE /api/v1/roles/:id` - 删除角色

### 权限管理
- `POST /api/v1/permissions` - 创建权限
- `GET /api/v1/permissions/:id` - 获取权限
- `GET /api/v1/permissions` - 权限列表
- `PUT /api/v1/permissions/:id` - 更新权限
- `DELETE /api/v1/permissions/:id` - 删除权限

### 菜单管理
- `POST /api/v1/menus` - 创建菜单
- `GET /api/v1/menus/:id` - 获取菜单
- `GET /api/v1/menus` - 菜单列表
- `GET /api/v1/menus/tree` - 菜单树
- `PUT /api/v1/menus/:id` - 更新菜单
- `DELETE /api/v1/menus/:id` - 删除菜单

### 权限分配
- `POST /api/v1/users/:user_id/roles` - 分配用户角色
- `POST /api/v1/roles/:role_id/permissions` - 分配角色权限
- `POST /api/v1/roles/:role_id/menus` - 分配角色菜单

### 权限查询
- `GET /api/v1/users/:user_id/roles` - 获取用户角色
- `GET /api/v1/users/:user_id/permissions` - 获取用户权限
- `GET /api/v1/users/:user_id/menus/tree` - 获取用户菜单树
- `GET /api/v1/users/:user_id/permissions/check` - 检查用户权限

## 🚀 快速开始

### 1. 环境准备
```bash
# 安装Go 1.19+
# 安装PostgreSQL
# 克隆项目
git clone <repository>
cd alice/backend
```

### 2. 配置数据库
```bash
# 启动PostgreSQL数据库
make db-setup

# 或者手动配置config.yaml中的数据库连接
```

### 3. 构建和运行
```bash
# 安装依赖
make deps

# 构建项目
make build

# 初始化RBAC数据
make init-data

# 运行应用
make run
```

### 4. 测试API
```bash
# 创建角色
curl -X POST http://localhost:8080/api/v1/roles \
  -H "Content-Type: application/json" \
  -d '{"name": "管理员", "code": "admin"}'

# 获取角色列表
curl http://localhost:8080/api/v1/roles
```

## 📝 使用示例

### 权限中间件使用
```go
// 在需要权限控制的路由上使用
router.GET("/api/v1/users", 
    middleware.RequirePermission(permissionService, "user", "read"),
    userHandler.ListUsers)

router.POST("/api/v1/users", 
    middleware.RequirePermission(permissionService, "user", "create"),
    userHandler.CreateUser)
```

### 前端权限控制
```javascript
// 获取用户菜单树
const menuTree = await fetch('/api/v1/users/current/menus/tree');

// 检查按钮权限
const hasPermission = await fetch('/api/v1/users/current/permissions/check?resource=user&action=delete');

// 根据权限显示/隐藏按钮
if (hasPermission.data.has_permission) {
    showDeleteButton();
}
```

## 🔒 安全特性

### 1. 认证授权
- JWT Token认证
- 权限中间件拦截
- 最小权限原则

### 2. 数据安全
- 密码哈希存储
- SQL注入防护(GORM)
- 参数验证

### 3. API安全
- CORS配置
- 请求日志记录
- 错误信息过滤

## 🎯 最佳实践

### 1. 权限设计
- 使用resource:action格式定义权限
- 遵循最小权限原则
- 定期审查权限分配

### 2. 菜单设计
- 合理设计菜单层级
- 使用语义化的菜单代码
- 按钮权限与菜单权限分离

### 3. 角色设计
- 基于业务职能设计角色
- 避免权限重复分配
- 定期清理无用角色

## 🚧 扩展性

### 1. 多租户支持
可扩展支持多租户架构，每个租户有独立的权限体系。

### 2. 权限缓存
可集成Redis缓存权限信息，提高查询性能。

### 3. 审计日志
可扩展审计日志功能，记录所有权限变更操作。

### 4. 工作流集成
可与工作流系统集成，支持权限申请审批流程。

## 📈 性能优化

### 1. 数据库优化
- 在关联表外键上建立索引
- 合理使用数据库连接池
- 查询优化和慢查询监控

### 2. 缓存策略
- 菜单树结构缓存
- 用户权限信息缓存
- 角色权限映射缓存

### 3. API优化
- 分页查询支持
- 批量操作接口
- 响应数据压缩

## 🧪 测试

### 单元测试
```bash
make test
```

### API测试
```bash
# 使用提供的API示例进行测试
# 参考docs/API_EXAMPLES.md
```

### 集成测试
```bash
make test-coverage
```

## 📚 文档

- [RBAC系统设计文档](docs/RBAC.md)
- [API使用示例](docs/API_EXAMPLES.md)
- [架构设计文档](docs/architecture.md)

## 🤝 贡献

1. Fork 项目
2. 创建特性分支
3. 提交变更
4. 推送到分支
5. 创建 Pull Request

## 📄 许可证

Apache License 2.0

## 💡 总结

这个RBAC权限管理系统提供了：

1. **完整的权限模型**：支持用户、角色、权限、菜单的完整关联
2. **精确的权限控制**：可以精确到页面按钮级别的权限
3. **清晰的架构设计**：严格遵循DDD分层架构
4. **丰富的API接口**：提供完整的RESTful API
5. **良好的扩展性**：便于后续功能扩展和优化
6. **完善的文档**：提供详细的使用文档和示例

该系统可以直接用于企业级应用的权限管理，也可以作为RBAC系统的参考实现。
