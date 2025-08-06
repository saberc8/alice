# Alice - Go后端项目

Alice是一个基于Go和Gin框架的后端项目，实现了用户认证系统。项目采用DDD（领域驱动设计）架构模式，参考了Coze Studio的架构设计。

## 技术栈

- **框架**: Gin
- **数据库**: PostgreSQL
- **ORM**: GORM
- **认证**: JWT
- **密码加密**: bcrypt
- **架构模式**: DDD (领域驱动设计)

## 项目结构

```
backend-go/
├── main.go                     # 应用入口
├── go.mod                      # Go模块定义
├── go.sum                      # 依赖版本锁定
├── application/                # 应用层
│   └── application.go         # 应用初始化
├── api/                        # API层
│   ├── handler/               # 处理器
│   ├── middleware/            # 中间件
│   ├── model/                 # 请求/响应模型
│   └── router/                # 路由
├── domain/                     # 领域层
│   └── user/                  # 用户领域
│       ├── entity/            # 实体
│       ├── repository/        # 仓储接口
│       └── service/           # 领域服务
├── infra/                      # 基础设施层
│   ├── config/                # 配置
│   ├── database/              # 数据库
│   └── repository/            # 仓储实现
└── pkg/                        # 通用包
    └── logger/                # 日志
```

## API接口

### 用户认证

- `POST /api/v1/users/register` - 用户注册
- `POST /api/v1/users/login` - 用户登录
- `GET /api/v1/users/profile` - 获取用户资料（需要认证）
- `PUT /api/v1/users/profile` - 更新用户资料（需要认证）

### 健康检查

- `GET /health` - 健康检查

## 环境配置

创建 `config.yaml` 文件：

```yaml
server:
  port: ":8080"
  
database:
  host: "localhost"
  port: 5432
  username: "postgres"
  password: "password"
  dbname: "alice"
  sslmode: "disable"

jwt:
  secret_key: "your-secret-key-here"
  expires_in: 24  # 小时

log:
  level: "info"
```

## 运行项目

1. 安装依赖：
```bash
go mod tidy
```

2. 启动PostgreSQL数据库

3. 配置数据库连接信息

4. 运行项目：
```bash
go run main.go
```

服务将在 `http://localhost:8080` 启动。

## 数据库

项目会自动创建用户表（users），包含以下字段：
- id: 主键
- username: 用户名（唯一）
- password_hash: 密码哈希
- email: 邮箱（唯一）
- status: 用户状态
- created_at: 创建时间
- updated_at: 更新时间

## API示例

### 用户注册
```bash
curl -X POST http://localhost:8080/api/v1/users/register \
  -H "Content-Type: application/json" \
  -d '{
    "username": "testuser",
    "password": "password123",
    "email": "test@example.com"
  }'
```

### 用户登录
```bash
curl -X POST http://localhost:8080/api/v1/users/login \
  -H "Content-Type: application/json" \
  -d '{
    "username": "testuser",
    "password": "password123"
  }'
```

### 获取用户资料
```bash
curl -X GET http://localhost:8080/api/v1/users/profile \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

## 架构特点

1. **DDD架构**: 清晰的层次分离和领域建模
2. **依赖注入**: 松耦合的组件设计
3. **接口抽象**: 便于测试和扩展
4. **配置管理**: 统一的配置管理
5. **错误处理**: 统一的错误处理机制
6. **日志记录**: 结构化日志记录

## 许可证

Apache License 2.0
