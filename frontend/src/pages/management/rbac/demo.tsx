import React from 'react';
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/ui/card';
import { Badge } from '@/ui/badge';
import { Button } from '@/ui/button';
import { Alert, AlertDescription } from '@/ui/alert';
import { 
  User, 
  Shield, 
  Menu, 
  Key, 
  Users, 
  UserCheck, 
  Settings,
  CheckCircle,
  AlertCircle,
  Info,
} from 'lucide-react';
import { 
  useUserInfo, 
  useUserRoles, 
  useUserPermissions, 
  useUserMenuTree,
  usePermissionCheck,
  useRoleCheck 
} from '@/store/userStore';
import { PermissionWrapper } from '@/components/auth/PermissionWrapper';
import { DynamicMenu, MenuBreadcrumb } from '@/components/nav/DynamicMenu';

const RBACDemo: React.FC = () => {
  const userInfo = useUserInfo();
  const userRoles = useUserRoles();
  const userPermissions = useUserPermissions();
  const userMenuTree = useUserMenuTree();
  const { hasPermission, hasPermissionByCode } = usePermissionCheck();
  const { hasRole, hasAnyRole } = useRoleCheck();

  return (
    <div className="space-y-6">
      {/* 页面标题和面包屑 */}
      <div className="space-y-2">
        <h1 className="text-3xl font-bold">RBAC 权限系统演示</h1>
        <MenuBreadcrumb />
        <p className="text-gray-600">
          展示基于角色的访问控制系统的各种功能和权限检查
        </p>
      </div>

      {/* 用户信息 */}
      <Card>
        <CardHeader>
          <CardTitle className="flex items-center gap-2">
            <User className="w-5 h-5" />
            当前用户信息
          </CardTitle>
        </CardHeader>
        <CardContent className="space-y-4">
          <div className="grid md:grid-cols-2 gap-4">
            <div>
              <h4 className="font-medium mb-2">基本信息</h4>
              <div className="space-y-2 text-sm">
                <div className="flex justify-between">
                  <span className="text-gray-600">用户ID:</span>
                  <span>{userInfo.id || '-'}</span>
                </div>
                <div className="flex justify-between">
                  <span className="text-gray-600">用户名:</span>
                  <span>{userInfo.username || '-'}</span>
                </div>
                <div className="flex justify-between">
                  <span className="text-gray-600">邮箱:</span>
                  <span>{userInfo.email || '-'}</span>
                </div>
                <div className="flex justify-between">
                  <span className="text-gray-600">状态:</span>
                  <Badge variant={userInfo.status === 'active' ? 'default' : 'secondary'}>
                    {userInfo.status === 'active' ? '激活' : '禁用'}
                  </Badge>
                </div>
              </div>
            </div>
            
            <div>
              <h4 className="font-medium mb-2">系统信息</h4>
              <div className="space-y-2 text-sm">
                <div className="flex justify-between">
                  <span className="text-gray-600">角色数量:</span>
                  <span>{userRoles.length}</span>
                </div>
                <div className="flex justify-between">
                  <span className="text-gray-600">权限数量:</span>
                  <span>{userPermissions.length}</span>
                </div>
                <div className="flex justify-between">
                  <span className="text-gray-600">菜单数量:</span>
                  <span>{userMenuTree.length}</span>
                </div>
                <div className="flex justify-between">
                  <span className="text-gray-600">登录时间:</span>
                  <span>{new Date().toLocaleString()}</span>
                </div>
              </div>
            </div>
          </div>
        </CardContent>
      </Card>

      {/* 用户角色 */}
      <Card>
        <CardHeader>
          <CardTitle className="flex items-center gap-2">
            <UserCheck className="w-5 h-5" />
            用户角色 ({userRoles.length})
          </CardTitle>
        </CardHeader>
        <CardContent>
          {userRoles.length > 0 ? (
            <div className="space-y-3">
              {userRoles.map((role) => (
                <div key={role.id} className="border rounded-lg p-3">
                  <div className="flex items-center justify-between">
                    <div>
                      <div className="font-medium">{role.name}</div>
                      <div className="text-sm text-gray-600">
                        代码: <code className="bg-gray-100 px-1 rounded">{role.code}</code>
                      </div>
                      {role.description && (
                        <div className="text-sm text-gray-500 mt-1">{role.description}</div>
                      )}
                    </div>
                    <Badge variant={role.status === 'active' ? 'default' : 'secondary'}>
                      {role.status === 'active' ? '激活' : '禁用'}
                    </Badge>
                  </div>
                </div>
              ))}
            </div>
          ) : (
            <div className="text-center text-gray-500 py-8">
              <UserCheck className="w-12 h-12 mx-auto mb-4 text-gray-300" />
              <p>暂无分配角色</p>
            </div>
          )}
        </CardContent>
      </Card>

      {/* 用户权限 */}
      <Card>
        <CardHeader>
          <CardTitle className="flex items-center gap-2">
            <Shield className="w-5 h-5" />
            用户权限 ({userPermissions.length})
          </CardTitle>
        </CardHeader>
        <CardContent>
          {userPermissions.length > 0 ? (
            <div className="grid md:grid-cols-2 gap-3">
              {userPermissions.map((permission) => (
                <div key={permission.id} className="border rounded-lg p-3">
                  <div className="flex items-start justify-between">
                    <div className="flex-1">
                      <div className="font-medium">{permission.name}</div>
                      <div className="text-sm text-gray-600 mt-1">
                        <code className="bg-gray-100 px-1 rounded">{permission.code}</code>
                      </div>
                      <div className="flex gap-2 mt-2">
                        <Badge variant="outline" className="text-xs">
                          {permission.resource}
                        </Badge>
                        <Badge variant="secondary" className="text-xs">
                          {permission.action}
                        </Badge>
                      </div>
                      {permission.description && (
                        <div className="text-sm text-gray-500 mt-1">{permission.description}</div>
                      )}
                    </div>
                    <Badge variant={permission.status === 'active' ? 'default' : 'secondary'} className="text-xs">
                      {permission.status === 'active' ? '激活' : '禁用'}
                    </Badge>
                  </div>
                </div>
              ))}
            </div>
          ) : (
            <div className="text-center text-gray-500 py-8">
              <Shield className="w-12 h-12 mx-auto mb-4 text-gray-300" />
              <p>暂无权限</p>
            </div>
          )}
        </CardContent>
      </Card>

      {/* 权限检查演示 */}
      <Card>
        <CardHeader>
          <CardTitle className="flex items-center gap-2">
            <Key className="w-5 h-5" />
            权限检查演示
          </CardTitle>
          <CardDescription>
            演示不同的权限检查方法和权限控制组件
          </CardDescription>
        </CardHeader>
        <CardContent className="space-y-4">
          {/* Hook方式权限检查 */}
          <div>
            <h4 className="font-medium mb-3">Hook方式权限检查</h4>
            <div className="grid md:grid-cols-2 gap-3">
              {[
                { resource: 'user', action: 'read', name: '查看用户' },
                { resource: 'user', action: 'create', name: '创建用户' },
                { resource: 'role', action: 'read', name: '查看角色' },
                { resource: 'permission', action: 'read', name: '查看权限' },
              ].map((item) => (
                <div key={`${item.resource}:${item.action}`} className="flex items-center justify-between border rounded-lg p-3">
                  <span className="text-sm">{item.name}</span>
                  <div className="flex items-center gap-2">
                    <code className="text-xs bg-gray-100 px-1 rounded">
                      {item.resource}:{item.action}
                    </code>
                    {hasPermission(item.resource, item.action) ? (
                      <CheckCircle className="w-4 h-4 text-green-500" />
                    ) : (
                      <AlertCircle className="w-4 h-4 text-red-500" />
                    )}
                  </div>
                </div>
              ))}
            </div>
          </div>

          {/* 组件方式权限控制 */}
          <div>
            <h4 className="font-medium mb-3">组件方式权限控制</h4>
            <div className="space-y-3">
              <PermissionWrapper 
                resource="user" 
                action="create"
                fallback={
                  <Alert>
                    <AlertCircle className="h-4 w-4" />
                    <AlertDescription>
                      您没有创建用户的权限，此按钮被隐藏
                    </AlertDescription>
                  </Alert>
                }
              >
                <Button className="w-full">
                  <Users className="w-4 h-4 mr-2" />
                  创建用户 (需要 user:create 权限)
                </Button>
              </PermissionWrapper>

              <PermissionWrapper 
                resource="role" 
                action="delete"
                fallback={
                  <Alert>
                    <AlertCircle className="h-4 w-4" />
                    <AlertDescription>
                      您没有删除角色的权限，此按钮被隐藏
                    </AlertDescription>
                  </Alert>
                }
              >
                <Button variant="destructive" className="w-full">
                  <UserCheck className="w-4 h-4 mr-2" />
                  删除角色 (需要 role:delete 权限)
                </Button>
              </PermissionWrapper>

              <PermissionWrapper 
                roleCode="admin"
                fallback={
                  <Alert>
                    <AlertCircle className="h-4 w-4" />
                    <AlertDescription>
                      您不是管理员，此功能被隐藏
                    </AlertDescription>
                  </Alert>
                }
              >
                <Button variant="outline" className="w-full">
                  <Settings className="w-4 h-4 mr-2" />
                  系统设置 (需要 admin 角色)
                </Button>
              </PermissionWrapper>
            </div>
          </div>
        </CardContent>
      </Card>

      {/* 用户菜单树 */}
      <div className="grid md:grid-cols-2 gap-6">
        <Card>
          <CardHeader>
            <CardTitle className="flex items-center gap-2">
              <Menu className="w-5 h-5" />
              用户菜单
            </CardTitle>
            <CardDescription>
              根据用户权限动态生成的菜单
            </CardDescription>
          </CardHeader>
          <CardContent>
            <DynamicMenu />
          </CardContent>
        </Card>

        {/* 角色检查演示 */}
        <Card>
          <CardHeader>
            <CardTitle className="flex items-center gap-2">
              <Info className="w-5 h-5" />
              角色检查演示
            </CardTitle>
          </CardHeader>
          <CardContent className="space-y-3">
            {[
              'admin',
              'manager', 
              'user',
              'guest',
            ].map((roleCode) => (
              <div key={roleCode} className="flex items-center justify-between border rounded-lg p-3">
                <span className="text-sm">角色: {roleCode}</span>
                <div className="flex items-center gap-2">
                  {hasRole(roleCode) ? (
                    <CheckCircle className="w-4 h-4 text-green-500" />
                  ) : (
                    <AlertCircle className="w-4 h-4 text-red-500" />
                  )}
                </div>
              </div>
            ))}
          </CardContent>
        </Card>
      </div>

      {/* 系统说明 */}
      <Alert>
        <Info className="h-4 w-4" />
        <AlertDescription>
          <div className="space-y-2">
            <h4 className="font-medium">RBAC 系统说明</h4>
            <p className="text-sm">
              此页面展示了完整的RBAC权限管理系统功能，包括：
            </p>
            <ul className="text-sm space-y-1 list-disc list-inside ml-4">
              <li>用户基本信息和系统状态</li>
              <li>用户角色分配和管理</li>
              <li>用户权限列表和状态</li>
              <li>多种权限检查方法的演示</li>
              <li>基于权限的组件显示控制</li>
              <li>动态菜单生成和权限控制</li>
              <li>角色检查和多角色支持</li>
            </ul>
            <p className="text-sm mt-2">
              系统支持页面级和按钮级的精确权限控制，确保用户只能访问被授权的功能。
            </p>
          </div>
        </AlertDescription>
      </Alert>
    </div>
  );
};

export default RBACDemo;
