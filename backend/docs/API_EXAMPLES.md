# RBAC API 使用示例

本文档提供了Alice RBAC系统的详细API使用示例。

## 前提条件

1. 确保数据库已配置并运行
2. 运行初始化脚本创建基础数据
3. 启动应用服务

```bash
# 设置数据库
make db-setup

# 构建并初始化数据
make rbac-setup

# 运行应用
make run
```

## API 示例

### 1. 角色管理

#### 创建角色

```bash
curl -X POST http://localhost:8080/api/v1/roles \
  -H "Content-Type: application/json" \
  -d '{
    "name": "产品经理",
    "code": "product_manager",
    "description": "负责产品规划和管理",
    "status": "active"
  }'
```

响应：
```json
{
  "code": 200,
  "message": "success",
  "data": {
    "id": "550e8400-e29b-41d4-a716-446655440000",
    "name": "产品经理",
    "code": "product_manager",
    "description": "负责产品规划和管理",
    "status": "active",
    "created_at": "2025-01-06T10:00:00Z",
    "updated_at": "2025-01-06T10:00:00Z"
  }
}
```

#### 获取角色列表

```bash
curl "http://localhost:8080/api/v1/roles?page=1&page_size=10"
```

#### 更新角色

```bash
curl -X PUT http://localhost:8080/api/v1/roles/550e8400-e29b-41d4-a716-446655440000 \
  -H "Content-Type: application/json" \
  -d '{
    "name": "高级产品经理",
    "code": "senior_product_manager",
    "description": "负责核心产品规划和团队管理",
    "status": "active"
  }'
```

### 2. 权限管理

#### 创建权限

```bash
# 创建产品管理权限
curl -X POST http://localhost:8080/api/v1/permissions \
  -H "Content-Type: application/json" \
  -d '{
    "name": "查看产品",
    "code": "product:read",
    "resource": "product",
    "action": "read",
    "description": "查看产品信息的权限"
  }'

curl -X POST http://localhost:8080/api/v1/permissions \
  -H "Content-Type: application/json" \
  -d '{
    "name": "创建产品",
    "code": "product:create",
    "resource": "product",
    "action": "create",
    "description": "创建新产品的权限"
  }'
```

#### 为角色分配权限

```bash
curl -X POST http://localhost:8080/api/v1/roles/550e8400-e29b-41d4-a716-446655440000/permissions \
  -H "Content-Type: application/json" \
  -d '{
    "permission_ids": [
      "permission_id_1",
      "permission_id_2"
    ]
  }'
```

### 3. 菜单管理

#### 创建菜单分组

```bash
curl -X POST http://localhost:8080/api/v1/menus \
  -H "Content-Type: application/json" \
  -d '{
    "name": "产品管理",
    "code": "product_management",
    "type": 0,
    "order": 3,
    "meta": {
      "icon": "product"
    }
  }'
```

#### 创建菜单目录

```bash
curl -X POST http://localhost:8080/api/v1/menus \
  -H "Content-Type: application/json" \
  -d '{
    "parent_id": "product_group_id",
    "name": "产品列表",
    "code": "product_list",
    "path": "/product",
    "type": 1,
    "order": 1,
    "meta": {
      "icon": "list"
    }
  }'
```

#### 创建菜单项

```bash
curl -X POST http://localhost:8080/api/v1/menus \
  -H "Content-Type: application/json" \
  -d '{
    "parent_id": "product_list_id",
    "name": "产品管理",
    "code": "product_manage",
    "path": "/product/manage",
    "type": 2,
    "order": 1,
    "meta": {
      "icon": "edit",
      "component": "/pages/product/manage"
    }
  }'
```

#### 创建按钮权限

```bash
curl -X POST http://localhost:8080/api/v1/menus \
  -H "Content-Type: application/json" \
  -d '{
    "parent_id": "product_manage_id",
    "name": "删除产品",
    "code": "product_delete_btn",
    "type": 3,
    "order": 1,
    "meta": {
      "auth": true
    }
  }'
```

### 4. 用户角色分配

#### 为用户分配角色

```bash
curl -X POST http://localhost:8080/api/v1/users/user123/roles \
  -H "Content-Type: application/json" \
  -d '{
    "role_ids": [
      "550e8400-e29b-41d4-a716-446655440000",
      "other_role_id"
    ]
  }'
```

#### 获取用户角色

```bash
curl "http://localhost:8080/api/v1/users/user123/roles"
```

### 5. 权限检查

#### 检查用户是否有特定权限

```bash
curl "http://localhost:8080/api/v1/users/user123/permissions/check?resource=product&action=create"
```

响应：
```json
{
  "code": 200,
  "message": "success",
  "data": {
    "has_permission": true
  }
}
```

#### 获取用户所有权限

```bash
curl "http://localhost:8080/api/v1/users/user123/permissions"
```

### 6. 菜单查询

#### 获取用户菜单树（用于前端导航）

```bash
curl "http://localhost:8080/api/v1/users/user123/menus/tree"
```

响应：
```json
{
  "code": 200,
  "message": "success",
  "data": [
    {
      "id": "dashboard_group",
      "name": "仪表板",
      "code": "dashboard",
      "type": 0,
      "order": 1,
      "children": [
        {
          "id": "workbench",
          "name": "工作台",
          "code": "workbench",
          "path": "/workbench",
          "type": 2,
          "order": 1,
          "meta": {
            "icon": "local:ic-workbench",
            "component": "/pages/dashboard/workbench"
          }
        }
      ]
    },
    {
      "id": "system_group",
      "name": "系统管理",
      "code": "system",
      "type": 0,
      "order": 2,
      "children": [
        {
          "id": "management",
          "name": "权限管理",
          "code": "management",
          "path": "/management",
          "type": 1,
          "order": 1,
          "children": [
            {
              "id": "user_management",
              "name": "用户管理",
              "code": "management:user",
              "path": "/management/user",
              "type": 2,
              "order": 1,
              "meta": {
                "component": "/pages/management/system/user"
              },
              "children": [
                {
                  "id": "user_delete_btn",
                  "name": "删除用户",
                  "code": "user_delete_btn",
                  "type": 3,
                  "order": 1
                }
              ]
            }
          ]
        }
      ]
    }
  ]
}
```

## 实际业务场景示例

### 场景1：电商系统权限设计

```bash
# 1. 创建角色
curl -X POST http://localhost:8080/api/v1/roles \
  -H "Content-Type: application/json" \
  -d '{"name": "运营人员", "code": "operator"}'

curl -X POST http://localhost:8080/api/v1/roles \
  -H "Content-Type: application/json" \
  -d '{"name": "财务人员", "code": "finance"}'

# 2. 创建权限
curl -X POST http://localhost:8080/api/v1/permissions \
  -H "Content-Type: application/json" \
  -d '{"name": "查看订单", "code": "order:read", "resource": "order", "action": "read"}'

curl -X POST http://localhost:8080/api/v1/permissions \
  -H "Content-Type: application/json" \
  -d '{"name": "处理退款", "code": "refund:process", "resource": "refund", "action": "process"}'

# 3. 分配权限给角色
curl -X POST http://localhost:8080/api/v1/roles/{operator_role_id}/permissions \
  -H "Content-Type: application/json" \
  -d '{"permission_ids": ["order:read", "refund:process"]}'

# 4. 分配角色给用户
curl -X POST http://localhost:8080/api/v1/users/operator001/roles \
  -H "Content-Type: application/json" \
  -d '{"role_ids": ["{operator_role_id}"]}'
```

### 场景2：多租户权限控制

```bash
# 租户A的管理员只能管理租户A的数据
curl -X POST http://localhost:8080/api/v1/permissions \
  -H "Content-Type: application/json" \
  -d '{
    "name": "管理租户A数据",
    "code": "tenant_a:manage",
    "resource": "tenant_a",
    "action": "manage"
  }'
```

### 场景3：动态菜单生成

前端可以基于用户菜单树API动态生成导航：

```javascript
// 前端示例代码
async function generateUserMenu() {
  const response = await fetch('/api/v1/users/current/menus/tree');
  const menuTree = await response.json();
  
  // 根据菜单数据生成导航组件
  return renderMenuComponent(menuTree.data);
}

// 检查按钮级权限
async function checkButtonPermission(buttonCode) {
  const userMenus = await getUserMenus();
  return hasButtonInMenus(userMenus, buttonCode);
}
```

## 错误处理

所有API都遵循统一的错误响应格式：

```json
{
  "code": 400,
  "message": "请求参数格式错误"
}
```

常见错误码：
- 400: 请求参数错误
- 401: 未认证
- 403: 权限不足
- 404: 资源不存在
- 500: 服务器内部错误

## 性能考虑

1. **权限缓存**：可以在Redis中缓存用户权限信息
2. **菜单缓存**：菜单树结构可以缓存，减少数据库查询
3. **分页查询**：列表接口支持分页，避免一次性加载大量数据
4. **索引优化**：在关联表的外键字段上建立索引

## 安全建议

1. **API认证**：所有RBAC管理接口都应该需要认证
2. **权限检查**：在业务API中使用权限中间件
3. **审计日志**：记录所有权限变更操作
4. **最小权限原则**：只分配必要的权限
