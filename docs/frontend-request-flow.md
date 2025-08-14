# 前端请求流程说明（基于后端现有实现）

- Base URL: /api/v1
- 认证方式: Bearer Token（放在请求头 Authorization: Bearer <token>）
- 通用响应包装: { code: number, message: string, data?: any }

## 1) 登录获取 Token

- 接口: POST /auth/login
- 入参:
  - username: string
  - password: string
- 成功返回(data):
  - token: string

示例请求体:
{
  "username": "admin",
  "password": "123456"
}

示例响应:
{
  "code": 200,
  "message": "login successful",
  "data": {
    "token": "<JWT_TOKEN>"
  }
}

可选：也可以先注册再返回 token
- 接口: POST /auth/register
- 成功返回(data): { user: { id, username, email }, token }

## 2) 获取当前用户资料（含角色）

- 接口: GET /auth/profile
- 头部: Authorization: Bearer <token>
- 成功返回(data):
  - id: number
  - username: string
  - email: string
  - roles?: Array<{ id: string; name: string; code: string }>

示例响应:
{
  "code": 200,
  "message": "success",
  "data": {
    "id": 1,
    "username": "admin",
    "email": "admin@example.com",
    "roles": [
      { "id": "e9f1...", "name": "管理员", "code": "admin" },
      { "id": "a12b...", "name": "编辑",   "code": "editor" }
    ]
  }
}

说明:
- 角色信息来自后端聚合，字段为精简版(id/name/code)，用于前端后续选择或直接使用。

## 3) 获取菜单

菜单有两种常用获取路径，按需求二选一：

### 3.1 按“选中的单个角色”获取菜单

- 列表（非树）:
  - 接口: GET /roles/{roleId}/menus
- 树形:
  - 接口: GET /roles/{roleId}/menus/tree

请求头都需要 Authorization: Bearer <token>

返回数据结构（示例为树形，字段与 docs/tree.json 类似，省略部分字段）：
[
  {
    "id": "3c12c347-...",
    "parent_id": null,
    "name": "仪表板",
    "code": "dashboard",
    "type": 0,
    "order": 1,
    "status": "active",
    "meta": {},
    "children": [
      {
        "id": "5c7f5c2e-...",
        "parent_id": "3c12c347-...",
        "name": "工作台",
        "code": "workbench",
        "path": "/workbench",
        "type": 2,
        "order": 1,
        "status": "active",
        "meta": {
          "icon": "local:ic-workbench",
          "component": "/pages/dashboard/workbench",
          "perms": ["system:menu:list", "system:menu:get", ...]
        }
      }
    ]
  }
]

适用场景:
- 前端允许“角色切换”来查看不同角色的菜单/权限视图。

### 3.2 按“当前用户综合权限”获取菜单

- 列表（非树）:
  - 接口: GET /users/{userId}/menus
- 树形:
  - 接口: GET /users/{userId}/menus/tree

说明:
- {userId} 为 profile 中返回的 id（数值），需要转成字符串拼接在路径中。
- 该方式由后端合并多角色菜单并去重/排序，更适合运行态按“用户整体权限”生成路由。

## 4) 其他接口（可选）

- 获取用户的角色列表（如需单独拉取/校验）
  - 接口: GET /users/{userId}/roles
- 获取完整菜单树（不做权限过滤，多用于管理端）
  - 接口: GET /menus/tree

## 调用顺序建议

1. POST /auth/login -> 获得 token
2. GET /auth/profile -> 获得用户基本信息与 roles
3. 决策：
   - 若按“角色视角”渲染：选定一个 roleId -> GET /roles/{roleId}/menus 或 /roles/{roleId}/menus/tree
   - 若按“用户综合权限”渲染：使用用户 id -> GET /users/{userId}/menus/tree
4. 根据返回的菜单结构，动态生成前端路由与侧边栏；使用 meta.perms 控制按钮级权限

## 错误与重试

- 未携带或 token 过期: 接口返回 401（code=401, message="unauthorized"），应跳转登录或刷新 token。
- 其余错误按服务端 message 处理，可统一 toast/重试。

## 头部与示例

请求头（示例）:
Authorization: Bearer <JWT_TOKEN>
Content-Type: application/json

