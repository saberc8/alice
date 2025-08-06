# Alice 项目架构设计文档

## 项目概述

Alice 是一个基于 DDD（领域驱动设计）架构的企业级 Go 后端项目，专注于用户认证和权限管理系统。项目采用分层架构，确保代码的可维护性、可扩展性和可测试性。

## 架构原则

### 1. 依赖倒置原则
- 高层模块不应该依赖低层模块，两者都应该依赖抽象
- 领域层定义接口，基础设施层实现接口

### 2. 关注点分离
- 每一层只关注自己的职责
- 业务逻辑与技术实现分离

### 3. 接口隔离原则
- 客户端不应该被迫依赖它不使用的接口
- 定义小而专注的接口

## 分层架构详解

### API 层 (api/)

#### 职责
- HTTP 请求处理
- 数据验证和转换
- 响应格式化
- 协议适配

#### 组件说明

**Handler (api/handler/)**
```go
// 处理器负责处理 HTTP 请求
type UserHandler struct {
    userService service.UserService
}

func (h *UserHandler) Register(c *gin.Context) {
    var req model.RegisterRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(400, gin.H{"error": err.Error()})
        return
    }
    
    user, err := h.userService.RegisterUser(req.Username, req.Email, req.Password)
    if err != nil {
        c.JSON(500, gin.H{"error": err.Error()})
        return
    }
    
    c.JSON(201, model.NewUserResponse(user))
}
```

**Middleware (api/middleware/)**
```go
// 中间件处理横切关注点
func JWTAuth() gin.HandlerFunc {
    return gin.HandlerFunc(func(c *gin.Context) {
        token := extractToken(c)
        if !validateToken(token) {
            c.AbortWithStatusJSON(401, gin.H{"error": "unauthorized"})
            return
        }
        c.Next()
    })
}
```

**Model (api/model/)**
```go
// 请求和响应模型
type RegisterRequest struct {
    Username string `json:"username" binding:"required,min=3,max=20"`
    Email    string `json:"email" binding:"required,email"`
    Password string `json:"password" binding:"required,min=6"`
}

type UserResponse struct {
    ID       uint   `json:"id"`
    Username string `json:"username"`
    Email    string `json:"email"`
    Status   string `json:"status"`
}
```

**Router (api/router/)**
```go
// 路由配置
func (r *Router) SetupRoutes() *gin.Engine {
    engine := gin.New()
    
    v1 := engine.Group("/api/v1")
    {
        auth := v1.Group("/auth")
        {
            auth.POST("/register", r.userHandler.Register)
            auth.POST("/login", r.userHandler.Login)
        }
        
        users := v1.Group("/users")
        users.Use(middleware.JWTAuth())
        {
            users.GET("/profile", r.userHandler.GetProfile)
            users.PUT("/profile", r.userHandler.UpdateProfile)
        }
    }
    
    return engine
}
```

### 应用层 (application/)

#### 职责
- 用例编排
- 事务管理
- 领域服务协调
- 依赖注入配置

#### 组件说明

**Application Service**
```go
// 应用服务负责用例编排
type Application struct {
    userService service.UserService
    roleService service.RoleService
}

func (app *Application) RegisterUserWithRole(username, email, password, roleName string) (*entity.User, error) {
    // 开始事务
    tx := app.db.Begin()
    defer func() {
        if r := recover(); r != nil {
            tx.Rollback()
        }
    }()
    
    // 创建用户
    user, err := app.userService.RegisterUser(username, email, password)
    if err != nil {
        tx.Rollback()
        return nil, err
    }
    
    // 分配角色
    role, err := app.roleService.GetRoleByName(roleName)
    if err != nil {
        tx.Rollback()
        return nil, err
    }
    
    err = app.userService.AssignRole(user.ID, role.ID)
    if err != nil {
        tx.Rollback()
        return nil, err
    }
    
    tx.Commit()
    return user, nil
}
```

### 领域层 (domain/)

#### 职责
- 核心业务逻辑
- 业务规则
- 领域模型
- 业务概念表达

#### 组件说明

**Entity (domain/{module}/entity/)**
```go
// 实体包含身份标识和业务逻辑
type User struct {
    ID           uint      `json:"id" gorm:"primaryKey"`
    Username     string    `json:"username" gorm:"uniqueIndex;not null"`
    Email        string    `json:"email" gorm:"uniqueIndex;not null"`
    PasswordHash string    `json:"-" gorm:"not null"`
    Status       UserStatus `json:"status" gorm:"not null;default:1"`
    CreatedAt    time.Time `json:"created_at"`
    UpdatedAt    time.Time `json:"updated_at"`
}

// 业务方法
func (u *User) IsActive() bool {
    return u.Status == UserStatusActive
}

func (u *User) CanLogin() bool {
    return u.IsActive() && u.PasswordHash != ""
}

func (u *User) ChangePassword(newPassword string) error {
    if len(newPassword) < 6 {
        return ErrInvalidPassword
    }
    
    hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
    if err != nil {
        return err
    }
    
    u.PasswordHash = string(hashedPassword)
    return nil
}
```

**Repository Interface (domain/{module}/repository/)**
```go
// 仓储接口定义数据访问抽象
type UserRepository interface {
    Create(user *entity.User) error
    GetByID(id uint) (*entity.User, error)
    GetByUsername(username string) (*entity.User, error)
    GetByEmail(email string) (*entity.User, error)
    Update(user *entity.User) error
    Delete(id uint) error
    List(offset, limit int) ([]*entity.User, int64, error)
    ExistsByUsername(username string) (bool, error)
    ExistsByEmail(email string) (bool, error)
}
```

**Domain Service (domain/{module}/service/)**
```go
// 领域服务包含业务逻辑
type UserService interface {
    RegisterUser(username, email, password string) (*entity.User, error)
    AuthenticateUser(username, password string) (*entity.User, error)
    UpdateUser(id uint, updates map[string]interface{}) error
    DeactivateUser(id uint) error
    ChangePassword(id uint, oldPassword, newPassword string) error
}

type userServiceImpl struct {
    userRepo repository.UserRepository
}

func (s *userServiceImpl) RegisterUser(username, email, password string) (*entity.User, error) {
    // 业务规则验证
    if len(username) < 3 {
        return nil, ErrInvalidUsername
    }
    
    if len(password) < 6 {
        return nil, ErrInvalidPassword
    }
    
    // 检查用户名是否已存在
    exists, err := s.userRepo.ExistsByUsername(username)
    if err != nil {
        return nil, err
    }
    if exists {
        return nil, ErrUserAlreadyExists
    }
    
    // 检查邮箱是否已存在
    exists, err = s.userRepo.ExistsByEmail(email)
    if err != nil {
        return nil, err
    }
    if exists {
        return nil, ErrEmailAlreadyExists
    }
    
    // 创建用户
    user := &entity.User{
        Username: username,
        Email:    email,
        Status:   entity.UserStatusActive,
    }
    
    err = user.ChangePassword(password)
    if err != nil {
        return nil, err
    }
    
    err = s.userRepo.Create(user)
    if err != nil {
        return nil, err
    }
    
    return user, nil
}
```

### 基础设施层 (infra/)

#### 职责
- 技术实现
- 外部系统集成
- 数据持久化
- 配置管理

#### 组件说明

**Repository Implementation (infra/repository/)**
```go
// 仓储实现
type userRepositoryImpl struct {
    db *gorm.DB
}

func NewUserRepository(db *gorm.DB) repository.UserRepository {
    return &userRepositoryImpl{db: db}
}

func (r *userRepositoryImpl) Create(user *entity.User) error {
    return r.db.Create(user).Error
}

func (r *userRepositoryImpl) GetByID(id uint) (*entity.User, error) {
    var user entity.User
    err := r.db.First(&user, id).Error
    if err != nil {
        if errors.Is(err, gorm.ErrRecordNotFound) {
            return nil, domain_user.ErrUserNotFound
        }
        return nil, err
    }
    return &user, nil
}

func (r *userRepositoryImpl) GetByUsername(username string) (*entity.User, error) {
    var user entity.User
    err := r.db.Where("username = ?", username).First(&user).Error
    if err != nil {
        if errors.Is(err, gorm.ErrRecordNotFound) {
            return nil, domain_user.ErrUserNotFound
        }
        return nil, err
    }
    return &user, nil
}
```

**Database (infra/database/)**
```go
// 数据库连接和配置
type Database struct {
    DB *gorm.DB
}

func NewDatabase(cfg *config.Config) (*Database, error) {
    dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
        cfg.Database.Host,
        cfg.Database.Port,
        cfg.Database.Username,
        cfg.Database.Password,
        cfg.Database.Name,
    )
    
    db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
        Logger: logger.Default.LogMode(logger.Info),
    })
    if err != nil {
        return nil, err
    }
    
    // 自动迁移
    err = db.AutoMigrate(&entity.User{}, &entity.Role{})
    if err != nil {
        return nil, err
    }
    
    return &Database{DB: db}, nil
}
```

**Configuration (infra/config/)**
```go
// 配置管理
type Config struct {
    Server   ServerConfig   `yaml:"server"`
    Database DatabaseConfig `yaml:"database"`
    JWT      JWTConfig      `yaml:"jwt"`
}

type ServerConfig struct {
    Port string `yaml:"port"`
}

type DatabaseConfig struct {
    Host     string `yaml:"host"`
    Port     int    `yaml:"port"`
    Username string `yaml:"username"`
    Password string `yaml:"password"`
    Name     string `yaml:"name"`
}

type JWTConfig struct {
    Secret string `yaml:"secret"`
    Expire int    `yaml:"expire"`
}

func LoadConfig() (*Config, error) {
    var cfg Config
    
    // 从文件加载
    data, err := ioutil.ReadFile("config.yaml")
    if err == nil {
        err = yaml.Unmarshal(data, &cfg)
        if err != nil {
            return nil, err
        }
    }
    
    // 环境变量覆盖
    if port := os.Getenv("SERVER_PORT"); port != "" {
        cfg.Server.Port = port
    }
    
    if host := os.Getenv("DB_HOST"); host != "" {
        cfg.Database.Host = host
    }
    
    // 设置默认值
    if cfg.Server.Port == "" {
        cfg.Server.Port = ":8090"
    }
    
    return &cfg, nil
}
```

## 新功能实现指南

### 角色权限系统实现示例

以实现角色权限系统为例，展示如何按照 DDD 架构添加新功能：

#### 1. 领域层设计

**角色实体 (domain/role/entity/role.go)**
```go
package entity

type Role struct {
    ID          uint         `json:"id" gorm:"primaryKey"`
    Name        string       `json:"name" gorm:"uniqueIndex;not null"`
    Description string       `json:"description"`
    Permissions []Permission `json:"permissions" gorm:"many2many:role_permissions;"`
    CreatedAt   time.Time    `json:"created_at"`
    UpdatedAt   time.Time    `json:"updated_at"`
}

type Permission struct {
    ID       uint   `json:"id" gorm:"primaryKey"`
    Code     string `json:"code" gorm:"uniqueIndex;not null"`
    Name     string `json:"name" gorm:"not null"`
    Resource string `json:"resource" gorm:"not null"`
    Action   string `json:"action" gorm:"not null"`
}

// 业务方法
func (r *Role) HasPermission(permissionCode string) bool {
    for _, permission := range r.Permissions {
        if permission.Code == permissionCode {
            return true
        }
    }
    return false
}

func (r *Role) AddPermission(permission Permission) {
    if !r.HasPermission(permission.Code) {
        r.Permissions = append(r.Permissions, permission)
    }
}
```

**角色仓储接口 (domain/role/repository/role_repository.go)**
```go
package repository

type RoleRepository interface {
    Create(role *entity.Role) error
    GetByID(id uint) (*entity.Role, error)
    GetByName(name string) (*entity.Role, error)
    Update(role *entity.Role) error
    Delete(id uint) error
    List(offset, limit int) ([]*entity.Role, int64, error)
    GetRoleWithPermissions(id uint) (*entity.Role, error)
}

type PermissionRepository interface {
    Create(permission *entity.Permission) error
    GetByID(id uint) (*entity.Permission, error)
    GetByCode(code string) (*entity.Permission, error)
    List() ([]*entity.Permission, error)
    GetByResource(resource string) ([]*entity.Permission, error)
}
```

**角色领域服务 (domain/role/service/role_service.go)**
```go
package service

type RoleService interface {
    CreateRole(name, description string) (*entity.Role, error)
    GetRoleByID(id uint) (*entity.Role, error)
    UpdateRole(id uint, updates map[string]interface{}) error
    DeleteRole(id uint) error
    AssignPermissions(roleID uint, permissionIDs []uint) error
    RemovePermissions(roleID uint, permissionIDs []uint) error
    GetRolePermissions(roleID uint) ([]*entity.Permission, error)
}

type roleServiceImpl struct {
    roleRepo       repository.RoleRepository
    permissionRepo repository.PermissionRepository
}

func (s *roleServiceImpl) CreateRole(name, description string) (*entity.Role, error) {
    // 检查角色名是否已存在
    existing, err := s.roleRepo.GetByName(name)
    if err == nil && existing != nil {
        return nil, ErrRoleAlreadyExists
    }
    
    role := &entity.Role{
        Name:        name,
        Description: description,
    }
    
    err = s.roleRepo.Create(role)
    if err != nil {
        return nil, err
    }
    
    return role, nil
}

func (s *roleServiceImpl) AssignPermissions(roleID uint, permissionIDs []uint) error {
    role, err := s.roleRepo.GetRoleWithPermissions(roleID)
    if err != nil {
        return err
    }
    
    for _, permissionID := range permissionIDs {
        permission, err := s.permissionRepo.GetByID(permissionID)
        if err != nil {
            return err
        }
        role.AddPermission(*permission)
    }
    
    return s.roleRepo.Update(role)
}
```

#### 2. 基础设施层实现

**角色仓储实现 (infra/repository/role_repository_impl.go)**
```go
package repository

type roleRepositoryImpl struct {
    db *gorm.DB
}

func NewRoleRepository(db *gorm.DB) domain_repository.RoleRepository {
    return &roleRepositoryImpl{db: db}
}

func (r *roleRepositoryImpl) Create(role *entity.Role) error {
    return r.db.Create(role).Error
}

func (r *roleRepositoryImpl) GetByID(id uint) (*entity.Role, error) {
    var role entity.Role
    err := r.db.First(&role, id).Error
    if err != nil {
        if errors.Is(err, gorm.ErrRecordNotFound) {
            return nil, domain_role.ErrRoleNotFound
        }
        return nil, err
    }
    return &role, nil
}

func (r *roleRepositoryImpl) GetRoleWithPermissions(id uint) (*entity.Role, error) {
    var role entity.Role
    err := r.db.Preload("Permissions").First(&role, id).Error
    if err != nil {
        if errors.Is(err, gorm.ErrRecordNotFound) {
            return nil, domain_role.ErrRoleNotFound
        }
        return nil, err
    }
    return &role, nil
}
```

#### 3. API 层实现

**角色请求模型 (api/model/role_request.go)**
```go
package model

type CreateRoleRequest struct {
    Name        string `json:"name" binding:"required,min=2,max=50"`
    Description string `json:"description" binding:"max=200"`
}

type UpdateRoleRequest struct {
    Name        string `json:"name" binding:"omitempty,min=2,max=50"`
    Description string `json:"description" binding:"max=200"`
}

type AssignPermissionsRequest struct {
    PermissionIDs []uint `json:"permission_ids" binding:"required"`
}
```

**角色响应模型 (api/model/role_response.go)**
```go
package model

type RoleResponse struct {
    ID          uint                 `json:"id"`
    Name        string               `json:"name"`
    Description string               `json:"description"`
    Permissions []PermissionResponse `json:"permissions,omitempty"`
    CreatedAt   time.Time            `json:"created_at"`
    UpdatedAt   time.Time            `json:"updated_at"`
}

type PermissionResponse struct {
    ID       uint   `json:"id"`
    Code     string `json:"code"`
    Name     string `json:"name"`
    Resource string `json:"resource"`
    Action   string `json:"action"`
}

func NewRoleResponse(role *entity.Role) *RoleResponse {
    resp := &RoleResponse{
        ID:          role.ID,
        Name:        role.Name,
        Description: role.Description,
        CreatedAt:   role.CreatedAt,
        UpdatedAt:   role.UpdatedAt,
    }
    
    for _, permission := range role.Permissions {
        resp.Permissions = append(resp.Permissions, PermissionResponse{
            ID:       permission.ID,
            Code:     permission.Code,
            Name:     permission.Name,
            Resource: permission.Resource,
            Action:   permission.Action,
        })
    }
    
    return resp
}
```

**角色处理器 (api/handler/role_handler.go)**
```go
package handler

type RoleHandler struct {
    roleService service.RoleService
}

func NewRoleHandler(roleService service.RoleService) *RoleHandler {
    return &RoleHandler{roleService: roleService}
}

func (h *RoleHandler) CreateRole(c *gin.Context) {
    var req model.CreateRoleRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(400, gin.H{"error": err.Error()})
        return
    }
    
    role, err := h.roleService.CreateRole(req.Name, req.Description)
    if err != nil {
        if errors.Is(err, domain_role.ErrRoleAlreadyExists) {
            c.JSON(409, gin.H{"error": "Role already exists"})
            return
        }
        c.JSON(500, gin.H{"error": "Internal server error"})
        return
    }
    
    c.JSON(201, model.NewRoleResponse(role))
}

func (h *RoleHandler) GetRole(c *gin.Context) {
    idStr := c.Param("id")
    id, err := strconv.ParseUint(idStr, 10, 32)
    if err != nil {
        c.JSON(400, gin.H{"error": "Invalid role ID"})
        return
    }
    
    role, err := h.roleService.GetRoleByID(uint(id))
    if err != nil {
        if errors.Is(err, domain_role.ErrRoleNotFound) {
            c.JSON(404, gin.H{"error": "Role not found"})
            return
        }
        c.JSON(500, gin.H{"error": "Internal server error"})
        return
    }
    
    c.JSON(200, model.NewRoleResponse(role))
}

func (h *RoleHandler) AssignPermissions(c *gin.Context) {
    idStr := c.Param("id")
    id, err := strconv.ParseUint(idStr, 10, 32)
    if err != nil {
        c.JSON(400, gin.H{"error": "Invalid role ID"})
        return
    }
    
    var req model.AssignPermissionsRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(400, gin.H{"error": err.Error()})
        return
    }
    
    err = h.roleService.AssignPermissions(uint(id), req.PermissionIDs)
    if err != nil {
        if errors.Is(err, domain_role.ErrRoleNotFound) {
            c.JSON(404, gin.H{"error": "Role not found"})
            return
        }
        c.JSON(500, gin.H{"error": "Internal server error"})
        return
    }
    
    c.JSON(200, gin.H{"message": "Permissions assigned successfully"})
}
```

#### 4. 路由注册

**更新路由 (api/router/router.go)**
```go
func (r *Router) SetupRoutes() *gin.Engine {
    engine := gin.New()
    
    v1 := engine.Group("/api/v1")
    {
        // 现有路由...
        
        // 角色管理路由
        roles := v1.Group("/roles")
        roles.Use(middleware.JWTAuth())
        roles.Use(middleware.RequirePermission("role:read"))
        {
            roles.GET("", r.roleHandler.ListRoles)
            roles.GET("/:id", r.roleHandler.GetRole)
        }
        
        roles.Use(middleware.RequirePermission("role:write"))
        {
            roles.POST("", r.roleHandler.CreateRole)
            roles.PUT("/:id", r.roleHandler.UpdateRole)
            roles.DELETE("/:id", r.roleHandler.DeleteRole)
            roles.POST("/:id/permissions", r.roleHandler.AssignPermissions)
            roles.DELETE("/:id/permissions", r.roleHandler.RemovePermissions)
        }
    }
    
    return engine
}
```

#### 5. 权限中间件

**权限检查中间件 (api/middleware/permission.go)**
```go
package middleware

func RequirePermission(permissionCode string) gin.HandlerFunc {
    return gin.HandlerFunc(func(c *gin.Context) {
        userID, exists := c.Get("user_id")
        if !exists {
            c.AbortWithStatusJSON(401, gin.H{"error": "Unauthorized"})
            return
        }
        
        // 获取用户角色和权限
        hasPermission, err := checkUserPermission(userID.(uint), permissionCode)
        if err != nil {
            c.AbortWithStatusJSON(500, gin.H{"error": "Internal server error"})
            return
        }
        
        if !hasPermission {
            c.AbortWithStatusJSON(403, gin.H{"error": "Insufficient permissions"})
            return
        }
        
        c.Next()
    })
}

func checkUserPermission(userID uint, permissionCode string) (bool, error) {
    // 实现权限检查逻辑
    // 1. 获取用户角色
    // 2. 检查角色是否有对应权限
    return true, nil // 简化实现
}
```

#### 6. 应用层集成

**更新依赖注入 (application/application.go)**
```go
var (
    UserSvc service.UserService
    RoleSvc service.RoleService
    PermissionSvc service.PermissionService
)

func Init(ctx context.Context, cfg *config.Config) error {
    // 数据库初始化
    db, err := database.NewDatabase(cfg)
    if err != nil {
        return err
    }
    
    // 仓储初始化
    userRepo := repository.NewUserRepository(db.DB)
    roleRepo := repository.NewRoleRepository(db.DB)
    permissionRepo := repository.NewPermissionRepository(db.DB)
    
    // 服务初始化
    UserSvc = service.NewUserService(userRepo)
    RoleSvc = service.NewRoleService(roleRepo, permissionRepo)
    PermissionSvc = service.NewPermissionService(permissionRepo)
    
    return nil
}
```

## 数据库设计

### 用户表 (users)
```sql
CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    username VARCHAR(50) NOT NULL UNIQUE,
    email VARCHAR(100) NOT NULL UNIQUE,
    password_hash VARCHAR(255) NOT NULL,
    status INTEGER NOT NULL DEFAULT 1,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

CREATE INDEX idx_users_username ON users(username);
CREATE INDEX idx_users_email ON users(email);
CREATE INDEX idx_users_status ON users(status);
```

### 角色表 (roles)
```sql
CREATE TABLE roles (
    id SERIAL PRIMARY KEY,
    name VARCHAR(50) NOT NULL UNIQUE,
    description TEXT,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

CREATE INDEX idx_roles_name ON roles(name);
```

### 权限表 (permissions)
```sql
CREATE TABLE permissions (
    id SERIAL PRIMARY KEY,
    code VARCHAR(100) NOT NULL UNIQUE,
    name VARCHAR(100) NOT NULL,
    resource VARCHAR(50) NOT NULL,
    action VARCHAR(50) NOT NULL,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

CREATE INDEX idx_permissions_code ON permissions(code);
CREATE INDEX idx_permissions_resource ON permissions(resource);
```

### 用户角色关联表 (user_roles)
```sql
CREATE TABLE user_roles (
    user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    role_id INTEGER NOT NULL REFERENCES roles(id) ON DELETE CASCADE,
    assigned_at TIMESTAMP DEFAULT NOW(),
    PRIMARY KEY (user_id, role_id)
);
```

### 角色权限关联表 (role_permissions)
```sql
CREATE TABLE role_permissions (
    role_id INTEGER NOT NULL REFERENCES roles(id) ON DELETE CASCADE,
    permission_id INTEGER NOT NULL REFERENCES permissions(id) ON DELETE CASCADE,
    PRIMARY KEY (role_id, permission_id)
);
```

## 错误处理策略

### 错误分类

#### 1. 领域错误
```go
// domain/user/errors.go
var (
    ErrUserNotFound      = errors.New("user not found")
    ErrUserAlreadyExists = errors.New("user already exists")
    ErrEmailAlreadyExists = errors.New("email already exists")
    ErrInvalidPassword   = errors.New("invalid password")
    ErrInvalidUsername   = errors.New("invalid username")
)

// domain/role/errors.go
var (
    ErrRoleNotFound      = errors.New("role not found")
    ErrRoleAlreadyExists = errors.New("role already exists")
    ErrInvalidPermission = errors.New("invalid permission")
)
```

#### 2. 基础设施错误
```go
// pkg/errors/infrastructure.go
var (
    ErrDatabaseConnection = errors.New("database connection failed")
    ErrCacheConnection    = errors.New("cache connection failed")
    ErrConfigLoad         = errors.New("config load failed")
)
```

#### 3. API 错误
```go
// pkg/errors/api.go
type APIError struct {
    Code    int    `json:"code"`
    Message string `json:"message"`
    Details string `json:"details,omitempty"`
}

func (e APIError) Error() string {
    return e.Message
}

var (
    ErrInvalidRequest   = APIError{Code: 400, Message: "Invalid request"}
    ErrUnauthorized     = APIError{Code: 401, Message: "Unauthorized"}
    ErrForbidden        = APIError{Code: 403, Message: "Forbidden"}
    ErrNotFound         = APIError{Code: 404, Message: "Resource not found"}
    ErrConflict         = APIError{Code: 409, Message: "Resource conflict"}
    ErrInternalServer   = APIError{Code: 500, Message: "Internal server error"}
)
```

### 错误处理中间件
```go
// api/middleware/error.go
func ErrorHandler() gin.HandlerFunc {
    return gin.HandlerFunc(func(c *gin.Context) {
        c.Next()
        
        if len(c.Errors) > 0 {
            err := c.Errors.Last().Err
            
            var apiErr APIError
            if errors.As(err, &apiErr) {
                c.JSON(apiErr.Code, apiErr)
                return
            }
            
            // 处理领域错误
            switch {
            case errors.Is(err, domain_user.ErrUserNotFound):
                c.JSON(404, APIError{Code: 404, Message: "User not found"})
            case errors.Is(err, domain_user.ErrUserAlreadyExists):
                c.JSON(409, APIError{Code: 409, Message: "User already exists"})
            default:
                c.JSON(500, APIError{Code: 500, Message: "Internal server error"})
            }
        }
    })
}
```

## 测试策略

### 单元测试

#### 领域层测试
```go
// domain/user/service/user_service_test.go
func TestUserService_RegisterUser(t *testing.T) {
    tests := []struct {
        name     string
        username string
        email    string
        password string
        want     *entity.User
        wantErr  error
    }{
        {
            name:     "successful registration",
            username: "testuser",
            email:    "test@example.com",
            password: "password123",
            want: &entity.User{
                Username: "testuser",
                Email:    "test@example.com",
                Status:   entity.UserStatusActive,
            },
            wantErr: nil,
        },
        {
            name:     "username too short",
            username: "ab",
            email:    "test@example.com",
            password: "password123",
            want:     nil,
            wantErr:  ErrInvalidUsername,
        },
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            mockRepo := mocks.NewUserRepository(t)
            service := NewUserService(mockRepo)
            
            if tt.wantErr == nil {
                mockRepo.On("ExistsByUsername", tt.username).Return(false, nil)
                mockRepo.On("ExistsByEmail", tt.email).Return(false, nil)
                mockRepo.On("Create", mock.AnythingOfType("*entity.User")).Return(nil)
            }
            
            got, err := service.RegisterUser(tt.username, tt.email, tt.password)
            
            if tt.wantErr != nil {
                assert.Error(t, err)
                assert.ErrorIs(t, err, tt.wantErr)
                assert.Nil(t, got)
            } else {
                assert.NoError(t, err)
                assert.Equal(t, tt.want.Username, got.Username)
                assert.Equal(t, tt.want.Email, got.Email)
                assert.Equal(t, tt.want.Status, got.Status)
            }
        })
    }
}
```

#### API 层测试
```go
// api/handler/user_handler_test.go
func TestUserHandler_Register(t *testing.T) {
    gin.SetMode(gin.TestMode)
    
    tests := []struct {
        name           string
        request        model.RegisterRequest
        mockSetup      func(*mocks.UserService)
        expectedStatus int
        expectedBody   string
    }{
        {
            name: "successful registration",
            request: model.RegisterRequest{
                Username: "testuser",
                Email:    "test@example.com",
                Password: "password123",
            },
            mockSetup: func(mockService *mocks.UserService) {
                user := &entity.User{
                    ID:       1,
                    Username: "testuser",
                    Email:    "test@example.com",
                    Status:   entity.UserStatusActive,
                }
                mockService.On("RegisterUser", "testuser", "test@example.com", "password123").Return(user, nil)
            },
            expectedStatus: 201,
        },
        {
            name: "user already exists",
            request: model.RegisterRequest{
                Username: "existinguser",
                Email:    "existing@example.com",
                Password: "password123",
            },
            mockSetup: func(mockService *mocks.UserService) {
                mockService.On("RegisterUser", "existinguser", "existing@example.com", "password123").Return(nil, domain_user.ErrUserAlreadyExists)
            },
            expectedStatus: 409,
        },
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            mockService := mocks.NewUserService(t)
            tt.mockSetup(mockService)
            
            handler := NewUserHandler(mockService)
            
            router := gin.New()
            router.POST("/register", handler.Register)
            
            body, _ := json.Marshal(tt.request)
            req, _ := http.NewRequest("POST", "/register", bytes.NewBuffer(body))
            req.Header.Set("Content-Type", "application/json")
            
            w := httptest.NewRecorder()
            router.ServeHTTP(w, req)
            
            assert.Equal(t, tt.expectedStatus, w.Code)
        })
    }
}
```

### 集成测试
```go
// test/integration/user_test.go
func TestUserIntegration(t *testing.T) {
    // 设置测试数据库
    db := setupTestDB(t)
    defer cleanupTestDB(t, db)
    
    // 初始化应用
    app := setupTestApp(t, db)
    
    t.Run("User Registration Flow", func(t *testing.T) {
        // 注册用户
        registerReq := map[string]string{
            "username": "testuser",
            "email":    "test@example.com",
            "password": "password123",
        }
        
        body, _ := json.Marshal(registerReq)
        req, _ := http.NewRequest("POST", "/api/v1/auth/register", bytes.NewBuffer(body))
        req.Header.Set("Content-Type", "application/json")
        
        w := httptest.NewRecorder()
        app.ServeHTTP(w, req)
        
        assert.Equal(t, 201, w.Code)
        
        var response map[string]interface{}
        json.Unmarshal(w.Body.Bytes(), &response)
        
        assert.Equal(t, "testuser", response["username"])
        assert.Equal(t, "test@example.com", response["email"])
    })
}
```

## 性能优化策略

### 1. 数据库优化
- 合理使用索引
- 查询优化
- 连接池配置
- 读写分离

### 2. 缓存策略
```go
// 用户信息缓存
type CachedUserRepository struct {
    repo  repository.UserRepository
    cache cache.Cache
}

func (r *CachedUserRepository) GetByID(id uint) (*entity.User, error) {
    key := fmt.Sprintf("user:%d", id)
    
    // 先从缓存获取
    var user entity.User
    if err := r.cache.Get(key, &user); err == nil {
        return &user, nil
    }
    
    // 缓存不存在，从数据库获取
    user, err := r.repo.GetByID(id)
    if err != nil {
        return nil, err
    }
    
    // 写入缓存
    r.cache.Set(key, user, 10*time.Minute)
    
    return user, nil
}
```

### 3. 并发控制
```go
// 使用 Context 控制超时
func (s *userServiceImpl) RegisterUser(ctx context.Context, username, email, password string) (*entity.User, error) {
    // 检查 context 是否已取消
    select {
    case <-ctx.Done():
        return nil, ctx.Err()
    default:
    }
    
    // 业务逻辑...
    
    return user, nil
}
```

## 监控和日志

### 结构化日志
```go
// pkg/logger/logger.go
type Logger struct {
    *logrus.Logger
}

func NewLogger() *Logger {
    log := logrus.New()
    log.SetFormatter(&logrus.JSONFormatter{})
    return &Logger{log}
}

func (l *Logger) WithContext(ctx context.Context) *logrus.Entry {
    entry := l.WithFields(logrus.Fields{})
    
    if traceID := ctx.Value("trace_id"); traceID != nil {
        entry = entry.WithField("trace_id", traceID)
    }
    
    if userID := ctx.Value("user_id"); userID != nil {
        entry = entry.WithField("user_id", userID)
    }
    
    return entry
}
```

### 指标收集
```go
// pkg/metrics/metrics.go
var (
    HTTPRequestDuration = prometheus.NewHistogramVec(
        prometheus.HistogramOpts{
            Name: "http_request_duration_seconds",
            Help: "HTTP request duration in seconds",
        },
        []string{"method", "path", "status"},
    )
    
    DatabaseQueryDuration = prometheus.NewHistogramVec(
        prometheus.HistogramOpts{
            Name: "database_query_duration_seconds",
            Help: "Database query duration in seconds",
        },
        []string{"operation", "table"},
    )
)

func init() {
    prometheus.MustRegister(HTTPRequestDuration)
    prometheus.MustRegister(DatabaseQueryDuration)
}
```

## 部署和运维

### Docker 配置
```dockerfile
# Dockerfile
FROM golang:1.21-alpine AS builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o alice .

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /app/alice .
COPY --from=builder /app/config.yaml .

EXPOSE 8090
CMD ["./alice"]
```

### Kubernetes 部署
```yaml
# k8s/deployment.yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: alice
spec:
  replicas: 3
  selector:
    matchLabels:
      app: alice
  template:
    metadata:
      labels:
        app: alice
    spec:
      containers:
      - name: alice
        image: alice:latest
        ports:
        - containerPort: 8090
        env:
        - name: DB_HOST
          value: "postgresql"
        - name: DB_PASSWORD
          valueFrom:
            secretKeyRef:
              name: db-secret
              key: password
        resources:
          requests:
            memory: "128Mi"
            cpu: "100m"
          limits:
            memory: "256Mi"
            cpu: "200m"
        livenessProbe:
          httpGet:
            path: /health
            port: 8090
          initialDelaySeconds: 30
          periodSeconds: 10
        readinessProbe:
          httpGet:
            path: /ready
            port: 8090
          initialDelaySeconds: 5
          periodSeconds: 5
```

## 总结

Alice 项目展示了如何使用 DDD 架构构建一个企业级的 Go 后端应用。通过清晰的分层设计、依赖倒置原则和关注点分离，项目具有良好的可维护性、可扩展性和可测试性。

### 架构优势

1. **清晰的职责分离**: 每一层都有明确的职责，便于维护和测试
2. **高内聚低耦合**: 模块之间的依赖关系清晰，易于修改和扩展
3. **技术无关性**: 业务逻辑不依赖具体的技术实现
4. **可测试性**: 通过依赖注入和接口抽象，易于编写单元测试

### 最佳实践

1. **遵循 SOLID 原则**: 单一职责、开放封闭、里氏替换、接口隔离、依赖倒置
2. **使用依赖注入**: 降低模块间的耦合度
3. **错误处理**: 定义明确的错误类型和处理策略
4. **日志记录**: 结构化日志，便于调试和监控
5. **测试驱动**: 编写全面的单元测试和集成测试

这个架构设计为项目的长期发展奠定了坚实的基础，支持团队协作开发和系统的持续演进。
