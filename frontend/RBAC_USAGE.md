# RBAC 前端权限管理系统使用说明

## 概述

本前端权限管理系统基于 React + TypeScript 实现，严格按照后端 RBAC API 接口规范设计，提供完整的用户、角色、权限和菜单管理功能，支持页面级和按钮级权限控制。

## 系统架构

### 核心模块
```
src/
├── api/services/           # API 服务层
│   ├── userService.ts     # 用户相关接口
│   ├── roleService.ts     # 角色管理接口
│   ├── permissionService.ts # 权限管理接口
│   └── menuService.ts     # 菜单管理接口
├── store/
│   └── userStore.ts       # 用户状态管理
├── components/auth/       # 权限控制组件
│   ├── PermissionWrapper.tsx # 权限包装组件
│   └── PermissionGuard.tsx   # 路由权限守卫
├── components/nav/
│   └── DynamicMenu.tsx    # 动态菜单组件
└── pages/management/rbac/ # RBAC管理页面
    ├── index.tsx          # RBAC 主页
    ├── UserManagement.tsx # 用户管理
    ├── RoleManagement.tsx # 角色管理
    ├── PermissionManagement.tsx # 权限管理
    ├── MenuManagement.tsx # 菜单管理
    └── demo.tsx          # 功能演示页面
```

## 核心功能

### 1. 用户认证和状态管理

#### 登录流程
```typescript
import { useSignIn } from '@/store/userStore';

const { signIn, isPending } = useSignIn();

await signIn({
  username: 'admin',
  password: 'password'
});
```

#### 用户状态获取
```typescript
import { 
  useUserInfo, 
  useUserRoles, 
  useUserPermissions, 
  useUserMenuTree 
} from '@/store/userStore';

const userInfo = useUserInfo();         // 用户基本信息
const userRoles = useUserRoles();       // 用户角色列表
const userPermissions = useUserPermissions(); // 用户权限列表
const userMenuTree = useUserMenuTree(); // 用户菜单树
```

### 2. 权限检查

#### Hook 方式权限检查
```typescript
import { usePermissionCheck, useRoleCheck } from '@/store/userStore';

const { hasPermission, hasPermissionByCode } = usePermissionCheck();
const { hasRole, hasAnyRole } = useRoleCheck();

// 资源+操作方式
const canCreateUser = hasPermission('user', 'create');

// 权限代码方式  
const canDeleteUser = hasPermissionByCode('user:delete');

// 角色检查
const isAdmin = hasRole('admin');
const hasManagerRole = hasAnyRole(['admin', 'manager']);
```

#### 组件方式权限控制
```tsx
import { PermissionWrapper } from '@/components/auth/PermissionWrapper';

// 基于权限的显示控制
<PermissionWrapper resource="user" action="create">
  <Button>创建用户</Button>
</PermissionWrapper>

// 权限代码方式
<PermissionWrapper permissionCode="user:delete">
  <Button variant="destructive">删除用户</Button>
</PermissionWrapper>

// 角色方式
<PermissionWrapper roleCode="admin">
  <AdminPanel />
</PermissionWrapper>

// 无权限时显示替代内容
<PermissionWrapper 
  resource="user" 
  action="create"
  fallback={<div>无权限访问</div>}
>
  <Button>创建用户</Button>
</PermissionWrapper>

// 无权限时隐藏
<PermissionWrapper resource="user" action="delete" hide>
  <Button>删除用户</Button>
</PermissionWrapper>
```

### 3. 路由权限控制

#### 路由守卫
```tsx
import { PermissionGuard } from '@/components/auth/PermissionGuard';

// 保护整个路由
<PermissionGuard resource="user" action="read">
  <UserManagement />
</PermissionGuard>

// 角色保护
<PermissionGuard roleCode="admin" showForbidden>
  <AdminPanel />
</PermissionGuard>

// 重定向到指定页面
<PermissionGuard resource="user" action="read" redirectTo="/403">
  <UserManagement />
</PermissionGuard>
```

### 4. 动态菜单

#### 菜单组件
```tsx
import { DynamicMenu, MenuBreadcrumb } from '@/components/nav/DynamicMenu';

// 动态菜单（根据用户权限显示）
<DynamicMenu />

// 面包屑导航
<MenuBreadcrumb />
```

## API 服务

### 用户管理
```typescript
import userService from '@/api/services/userService';

// 获取用户列表
const users = await userService.getUsers({ 
  page: 1, 
  page_size: 10,
  username: 'admin' 
});

// 创建用户
const newUser = await userService.createUser({
  username: 'newuser',
  email: 'user@example.com',
  password: 'password123'
});

// 为用户分配角色
await userService.assignUserRoles(userId, {
  role_ids: ['role1', 'role2']
});

// 获取用户权限
const permissions = await userService.getUserPermissions(userId);

// 获取用户菜单树
const menuTree = await userService.getUserMenuTree(userId);
```

### 角色管理
```typescript
import roleService from '@/api/services/roleService';

// 获取角色列表
const roles = await roleService.getRoles({
  page: 1,
  page_size: 10,
  name: '管理员'
});

// 创建角色
const newRole = await roleService.createRole({
  name: '产品经理',
  code: 'product_manager',
  description: '负责产品规划和管理'
});

// 为角色分配权限
await roleService.assignRolePermissions(roleId, {
  permission_ids: ['perm1', 'perm2']
});

// 为角色分配菜单
await roleService.assignRoleMenus(roleId, {
  menu_ids: ['menu1', 'menu2']
});
```

### 权限管理
```typescript
import permissionService from '@/api/services/permissionService';

// 创建权限
const newPermission = await permissionService.createPermission({
  name: '查看用户',
  code: 'user:read',
  resource: 'user',
  action: 'read',
  description: '查看用户信息的权限'
});

// 检查用户权限
const hasPermission = await permissionService.checkUserPermission(userId, {
  resource: 'user',
  action: 'create'
});
```

### 菜单管理
```typescript
import menuService from '@/api/services/menuService';

// 创建菜单
const newMenu = await menuService.createMenu({
  parent_id: 'parent_menu_id',
  name: '用户管理',
  code: 'user_management',
  path: '/management/users',
  type: 2, // 菜单类型
  order: 1,
  meta: {
    icon: 'user',
    component: '/pages/management/users'
  }
});

// 获取菜单树
const menuTree = await menuService.getMenuTree();
```

## 使用示例

### 完整的CRUD页面实现
```tsx
import React, { useState } from 'react';
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';
import { PermissionWrapper } from '@/components/auth/PermissionWrapper';
import userService from '@/api/services/userService';

const UserManagement: React.FC = () => {
  const [searchParams, setSearchParams] = useState({ page: 1, page_size: 10 });
  const queryClient = useQueryClient();

  // 获取用户列表（需要读取权限）
  const { data: usersData, isLoading } = useQuery({
    queryKey: ['users', searchParams],
    queryFn: () => userService.getUsers(searchParams),
  });

  // 创建用户（需要创建权限）
  const createMutation = useMutation({
    mutationFn: userService.createUser,
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['users'] });
    },
  });

  // 删除用户（需要删除权限）  
  const deleteMutation = useMutation({
    mutationFn: userService.deleteUser,
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['users'] });
    },
  });

  return (
    <div className="space-y-6">
      {/* 页面标题和创建按钮 */}
      <div className="flex justify-between items-center">
        <h1 className="text-2xl font-bold">用户管理</h1>
        <PermissionWrapper resource="user" action="create">
          <Button onClick={() => setIsCreateDialogOpen(true)}>
            创建用户
          </Button>
        </PermissionWrapper>
      </div>

      {/* 用户列表 */}
      <Table>
        <TableHeader>
          <TableRow>
            <TableHead>用户名</TableHead>
            <TableHead>邮箱</TableHead>
            <TableHead>状态</TableHead>
            <TableHead className="text-right">操作</TableHead>
          </TableRow>
        </TableHeader>
        <TableBody>
          {users.map((user) => (
            <TableRow key={user.id}>
              <TableCell>{user.username}</TableCell>
              <TableCell>{user.email}</TableCell>
              <TableCell>
                <Badge variant={user.status === 'active' ? 'default' : 'secondary'}>
                  {user.status === 'active' ? '激活' : '禁用'}
                </Badge>
              </TableCell>
              <TableCell className="text-right">
                <div className="flex gap-2 justify-end">
                  <PermissionWrapper resource="user" action="update">
                    <Button variant="outline" size="sm">
                      编辑
                    </Button>
                  </PermissionWrapper>
                  
                  <PermissionWrapper resource="user" action="delete">
                    <Button 
                      variant="destructive" 
                      size="sm"
                      onClick={() => deleteMutation.mutate(user.id)}
                    >
                      删除
                    </Button>
                  </PermissionWrapper>
                </div>
              </TableCell>
            </TableRow>
          ))}
        </TableBody>
      </Table>
    </div>
  );
};
```

## 路由配置

### 添加RBAC路由
在 `src/routes/sections/dashboard/frontend.tsx` 中添加：

```tsx
{
  path: "management",
  children: [
    {
      path: "rbac",
      children: [
        { index: true, element: Component("/pages/management/rbac") },
        { path: "users", element: Component("/pages/management/rbac/UserManagement") },
        { path: "roles", element: Component("/pages/management/rbac/RoleManagement") },
        { path: "permissions", element: Component("/pages/management/rbac/PermissionManagement") },
        { path: "menus", element: Component("/pages/management/rbac/MenuManagement") },
        { path: "demo", element: Component("/pages/management/rbac/demo") },
      ],
    },
  ],
},
```

## 权限设计规范

### 权限代码格式
权限代码采用 `resource:action` 格式：
- `user:read` - 查看用户
- `user:create` - 创建用户  
- `user:update` - 更新用户
- `user:delete` - 删除用户
- `role:assign_permission` - 为角色分配权限

### 菜单类型说明
- `0` - 分组：菜单分组，不显示在导航中
- `1` - 目录：菜单目录，可包含子菜单
- `2` - 菜单：具体的页面菜单，可点击跳转
- `3` - 按钮：按钮权限，用于页面内功能控制

### 角色设计建议
- `super_admin` - 超级管理员：拥有所有权限
- `admin` - 系统管理员：管理用户、角色、权限
- `manager` - 部门经理：管理本部门相关功能
- `user` - 普通用户：基础功能使用权限

## 最佳实践

### 1. 权限检查优化
```typescript
// 推荐：使用 useMemo 缓存权限检查结果
const canManageUsers = useMemo(() => {
  return hasPermission('user', 'create') || 
         hasPermission('user', 'update') || 
         hasPermission('user', 'delete');
}, [hasPermission]);
```

### 2. 错误处理
```typescript
// API 调用时添加错误处理
try {
  const result = await userService.createUser(userData);
  toast.success('用户创建成功');
  return result;
} catch (error: any) {
  toast.error(error.message || '用户创建失败');
  throw error;
}
```

### 3. 加载状态管理
```typescript
// 使用 React Query 管理加载状态
const { data, isLoading, error } = useQuery({
  queryKey: ['users'],
  queryFn: userService.getUsers,
});

if (isLoading) return <LoadingSpinner />;
if (error) return <ErrorMessage error={error} />;
```

### 4. 权限缓存
系统会自动缓存用户权限信息到 localStorage，登录后无需重复获取。

## 故障排除

### 常见问题

1. **权限检查不生效**
   - 确认用户已登录且获取到权限数据
   - 检查权限代码格式是否正确
   - 确认后端权限数据返回正常

2. **菜单不显示**
   - 检查用户是否有对应菜单权限
   - 确认菜单类型设置正确
   - 检查菜单状态是否为激活

3. **页面访问被拒绝**
   - 确认路由配置正确
   - 检查权限守卫设置
   - 确认用户角色和权限分配

4. **API 调用失败**
   - 检查 API 接口地址配置
   - 确认认证 token 有效
   - 检查网络连接和跨域设置

### 调试技巧

1. **使用浏览器开发者工具**
   ```javascript
   // 控制台查看用户状态
   console.log('User Store:', useUserStore.getState());
   ```

2. **权限检查调试**
   ```typescript
   console.log('User Permissions:', userPermissions);
   console.log('Has Permission:', hasPermission('user', 'create'));
   ```

3. **网络请求调试**
   - 查看 Network 面板中的 API 请求
   - 检查请求头中的 Authorization 信息
   - 确认响应数据格式正确

## 总结

本 RBAC 前端权限管理系统提供了：

1. **完整的权限管理功能**：用户、角色、权限、菜单的 CRUD 操作
2. **灵活的权限控制**：支持页面级和按钮级权限控制  
3. **多种权限检查方式**：Hook、组件、路由守卫等
4. **动态菜单生成**：根据用户权限自动生成导航菜单
5. **良好的用户体验**：加载状态、错误处理、权限提示等
6. **可扩展的架构**：清晰的模块划分，便于功能扩展

系统严格按照后端 API 接口规范设计，确保前后端数据格式一致，提供了完整的权限管理解决方案。
