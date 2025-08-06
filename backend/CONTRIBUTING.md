# Contributing to Alice

We love your input! We want to make contributing to Alice as easy and transparent as possible, whether it's:

- Reporting a bug
- Discussing the current state of the code
- Submitting a fix
- Proposing new features
- Becoming a maintainer

## Development Process

We use GitHub to host code, to track issues and feature requests, as well as accept pull requests.

### Pull Requests

1. Fork the repo and create your branch from `main`.
2. If you've added code that should be tested, add tests.
3. If you've changed APIs, update the documentation.
4. Ensure the test suite passes.
5. Make sure your code lints.
6. Issue that pull request!

### Development Setup

1. **Clone the repository**
```bash
git clone https://github.com/coze-dev/alice.git
cd alice
```

2. **Install dependencies**
```bash
make deps
```

3. **Set up development environment**
```bash
cp config.yaml.example config.yaml
# Edit config.yaml with your local settings
```

4. **Start the development server**
```bash
make dev
```

### Code Style

#### Go Code Style

We follow the standard Go conventions:

- Use `gofmt` to format your code
- Use `golint` to lint your code
- Follow the [Effective Go](https://golang.org/doc/effective_go.html) guidelines
- Follow the [Go Code Review Comments](https://github.com/golang/go/wiki/CodeReviewComments)

#### DDD Architecture Conventions

##### Domain Layer Rules

1. **Entities must be in `domain/{module}/entity/`**
```go
// domain/user/entity/user.go
package entity

type User struct {
    ID       uint   `json:"id" gorm:"primaryKey"`
    Username string `json:"username" gorm:"uniqueIndex;not null"`
    // ...
}
```

2. **Repository interfaces must be in `domain/{module}/repository/`**
```go
// domain/user/repository/user_repository.go
package repository

type UserRepository interface {
    Create(user *entity.User) error
    GetByID(id uint) (*entity.User, error)
    // ...
}
```

3. **Domain services must be in `domain/{module}/service/`**
```go
// domain/user/service/user_service.go
package service

type UserService interface {
    RegisterUser(username, email, password string) (*entity.User, error)
    AuthenticateUser(username, password string) (*entity.User, error)
    // ...
}
```

##### Infrastructure Layer Rules

1. **Repository implementations must be in `infra/repository/`**
```go
// infra/repository/user_repository_impl.go
package repository

type userRepositoryImpl struct {
    db *gorm.DB
}

func NewUserRepository(db *gorm.DB) domain_repository.UserRepository {
    return &userRepositoryImpl{db: db}
}
```

##### API Layer Rules

1. **Handlers must be in `api/handler/`**
```go
// api/handler/user_handler.go
package handler

type UserHandler struct {
    userService service.UserService
}

func NewUserHandler(userService service.UserService) *UserHandler {
    return &UserHandler{userService: userService}
}
```

2. **Request/Response models must be in `api/model/`**
```go
// api/model/user_request.go
package model

type RegisterRequest struct {
    Username string `json:"username" binding:"required,min=3,max=20"`
    Email    string `json:"email" binding:"required,email"`
    Password string `json:"password" binding:"required,min=6"`
}
```

3. **Middleware must be in `api/middleware/`**
```go
// api/middleware/auth.go
package middleware

func JWTAuth() gin.HandlerFunc {
    return gin.HandlerFunc(func(c *gin.Context) {
        // JWT validation logic
    })
}
```

### New Feature Implementation Guide

When adding a new feature (e.g., Role management), follow this checklist:

#### 1. Domain Layer Implementation

**Step 1**: Create entity
```bash
mkdir -p domain/role/entity
touch domain/role/entity/role.go
```

**Step 2**: Define repository interface
```bash
mkdir -p domain/role/repository
touch domain/role/repository/role_repository.go
```

**Step 3**: Define domain service interface
```bash
mkdir -p domain/role/service
touch domain/role/service/role_service.go
```

#### 2. Infrastructure Layer Implementation

**Step 4**: Implement repository
```bash
touch infra/repository/role_repository_impl.go
```

**Step 5**: Implement domain service
```bash
touch infra/service/role_service_impl.go
```

#### 3. API Layer Implementation

**Step 6**: Create request/response models
```bash
touch api/model/role_request.go
touch api/model/role_response.go
```

**Step 7**: Create handler
```bash
touch api/handler/role_handler.go
```

**Step 8**: Register routes
- Update `api/router/router.go`

#### 4. Application Layer Integration

**Step 9**: Update dependency injection
- Update `application/application.go`

#### 5. Testing

**Step 10**: Add tests
```bash
touch domain/role/service/role_service_test.go
touch api/handler/role_handler_test.go
touch infra/repository/role_repository_test.go
```

### Testing Guidelines

#### Unit Tests
- Test files should end with `_test.go`
- Use table-driven tests when possible
- Mock external dependencies
- Aim for > 80% test coverage

#### Integration Tests
- Place in `test/integration/`
- Test complete workflows
- Use test databases

#### Example Test Structure
```go
func TestUserService_RegisterUser(t *testing.T) {
    tests := []struct {
        name    string
        input   RegisterInput
        want    *entity.User
        wantErr bool
    }{
        {
            name: "successful registration",
            input: RegisterInput{
                Username: "testuser",
                Email:    "test@example.com",
                Password: "password123",
            },
            want: &entity.User{
                Username: "testuser",
                Email:    "test@example.com",
            },
            wantErr: false,
        },
        // Add more test cases
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            // Test implementation
        })
    }
}
```

### Database Migrations

When adding new database schema changes:

1. **Create migration file**
```bash
# Use timestamp for ordering
touch migrations/20231201120000_create_roles_table.sql
```

2. **Migration content**
```sql
-- Up migration
CREATE TABLE roles (
    id SERIAL PRIMARY KEY,
    name VARCHAR(50) NOT NULL UNIQUE,
    description TEXT,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

-- Down migration (at the bottom)
-- DROP TABLE roles;
```

### Error Handling Conventions

1. **Define domain-specific errors**
```go
// domain/role/errors.go
package role

import "errors"

var (
    ErrRoleNotFound      = errors.New("role not found")
    ErrRoleAlreadyExists = errors.New("role already exists")
    ErrInvalidRoleName   = errors.New("invalid role name")
)
```

2. **Use structured error responses**
```go
// pkg/errors/api_error.go
type APIError struct {
    Code    string `json:"code"`
    Message string `json:"message"`
    Details string `json:"details,omitempty"`
}
```

### Logging Conventions

1. **Use structured logging**
```go
logger.WithFields(logger.Fields{
    "user_id": userID,
    "action":  "create_role",
}).Info("Creating new role")
```

2. **Log levels**
- `Debug`: Detailed information for debugging
- `Info`: General information about app behavior
- `Warn`: Warning messages for unusual situations
- `Error`: Error messages for failures
- `Fatal`: Critical errors that cause app termination

### Documentation Requirements

1. **API Documentation**
   - All public APIs must have Swagger/OpenAPI documentation
   - Include example requests and responses

2. **Code Documentation**
   - Public functions and types must have Go doc comments
   - Complex business logic should be well-commented

3. **Architecture Documentation**
   - Update architecture docs when adding new modules
   - Document design decisions in ADRs (Architecture Decision Records)

### Commit Message Format

We follow the [Conventional Commits](https://www.conventionalcommits.org/) specification:

```
<type>[optional scope]: <description>

[optional body]

[optional footer(s)]
```

**Types:**
- `feat`: A new feature
- `fix`: A bug fix
- `docs`: Documentation only changes
- `style`: Changes that do not affect the meaning of the code
- `refactor`: A code change that neither fixes a bug nor adds a feature
- `perf`: A code change that improves performance
- `test`: Adding missing tests or correcting existing tests
- `chore`: Changes to the build process or auxiliary tools

**Examples:**
```
feat(auth): add role-based access control

Add RBAC support for user authorization with roles and permissions.

Closes #123
```

```
fix(user): handle duplicate username error

Properly handle database constraint errors when username already exists.
```

### Issue and Pull Request Templates

When creating issues or pull requests, please use the provided templates and fill out all relevant sections.

### License

By contributing, you agree that your contributions will be licensed under the Apache 2.0 License.

## Questions?

Feel free to contact the maintainers if you have any questions about contributing!
