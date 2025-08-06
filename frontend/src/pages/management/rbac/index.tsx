import React from 'react';
import { Link } from 'react-router-dom';
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/ui/card';
import { Badge } from '@/ui/badge';
import { 
  Users, 
  Shield, 
  Lock, 
  Menu,
  UserCheck,
  Settings,
  AlertCircle,
} from 'lucide-react';
import { PermissionWrapper } from '@/components/auth/PermissionWrapper';

const RBACDashboard: React.FC = () => {
  const modules = [
    {
      title: '用户管理',
      description: '管理系统用户，包括创建、编辑、删除用户，以及为用户分配角色',
      icon: <Users className="w-8 h-8 text-blue-500" />,
      link: '/management/rbac/users',
      resource: 'user',
      action: 'read',
      features: [
        '用户CRUD操作',
        '用户状态管理',
        '角色分配',
        '批量操作',
      ],
    },
    {
      title: '角色管理',
      description: '管理系统角色，为角色分配权限和菜单',
      icon: <UserCheck className="w-8 h-8 text-green-500" />,
      link: '/management/rbac/roles',
      resource: 'role',
      action: 'read',
      features: [
        '角色CRUD操作',
        '权限分配',
        '菜单分配',
        '角色状态管理',
      ],
    },
    {
      title: '权限管理',
      description: '管理系统权限，定义资源和操作的访问控制',
      icon: <Shield className="w-8 h-8 text-purple-500" />,
      link: '/management/rbac/permissions',
      resource: 'permission',
      action: 'read',
      features: [
        '权限CRUD操作',
        '资源定义',
        '操作定义',
        '权限代码管理',
      ],
    },
    {
      title: '菜单管理',
      description: '管理系统菜单，支持多级菜单和按钮权限',
      icon: <Menu className="w-8 h-8 text-orange-500" />,
      link: '/management/rbac/menus',
      resource: 'menu',
      action: 'read',
      features: [
        '菜单CRUD操作',
        '菜单树结构',
        '按钮权限',
        '菜单排序',
      ],
    },
  ];

  const systemInfo = [
    {
      label: '权限模型',
      value: 'RBAC (基于角色的访问控制)',
      type: 'info' as const,
    },
    {
      label: '权限级别',
      value: '页面级 + 按钮级',
      type: 'success' as const,
    },
    {
      label: '菜单层级',
      value: '支持无限层级',
      type: 'info' as const,
    },
    {
      label: '权限格式',
      value: 'resource:action',
      type: 'warning' as const,
    },
  ];

  return (
    <div className="space-y-6">
      {/* 页面标题 */}
      <div className="space-y-2">
        <h1 className="text-3xl font-bold">RBAC 权限管理系统</h1>
        <p className="text-gray-600">
          基于角色的访问控制(Role-Based Access Control)，支持用户、角色、权限和菜单的完整管理
        </p>
      </div>

      {/* 系统信息 */}
      <Card>
        <CardHeader>
          <CardTitle className="flex items-center gap-2">
            <Settings className="w-5 h-5" />
            系统信息
          </CardTitle>
        </CardHeader>
        <CardContent>
          <div className="grid grid-cols-2 md:grid-cols-4 gap-4">
            {systemInfo.map((info, index) => (
              <div key={index} className="space-y-1">
                <div className="text-sm text-gray-500">{info.label}</div>
                <Badge variant={info.type}>
                  {info.value}
                </Badge>
              </div>
            ))}
          </div>
        </CardContent>
      </Card>

      {/* 功能模块 */}
      <div className="grid md:grid-cols-2 gap-6">
        {modules.map((module) => (
          <PermissionWrapper
            key={module.title}
            resource={module.resource}
            action={module.action}
            fallback={
              <Card className="opacity-50">
                <CardHeader>
                  <div className="flex items-center gap-3">
                    {module.icon}
                    <div>
                      <CardTitle>{module.title}</CardTitle>
                      <CardDescription>
                        {module.description}
                      </CardDescription>
                    </div>
                  </div>
                </CardHeader>
                <CardContent>
                  <div className="flex items-center gap-2 text-amber-600">
                    <AlertCircle className="w-4 h-4" />
                    <span className="text-sm">您没有访问此模块的权限</span>
                  </div>
                </CardContent>
              </Card>
            }
          >
            <Link to={module.link}>
              <Card className="h-full hover:shadow-lg transition-shadow cursor-pointer">
                <CardHeader>
                  <div className="flex items-center gap-3">
                    {module.icon}
                    <div>
                      <CardTitle>{module.title}</CardTitle>
                      <CardDescription>
                        {module.description}
                      </CardDescription>
                    </div>
                  </div>
                </CardHeader>
                <CardContent>
                  <div className="space-y-2">
                    <div className="text-sm font-medium text-gray-700">主要功能：</div>
                    <ul className="space-y-1">
                      {module.features.map((feature, index) => (
                        <li key={index} className="text-sm text-gray-600 flex items-center gap-2">
                          <div className="w-1.5 h-1.5 bg-gray-400 rounded-full" />
                          {feature}
                        </li>
                      ))}
                    </ul>
                  </div>
                </CardContent>
              </Card>
            </Link>
          </PermissionWrapper>
        ))}
      </div>

      {/* 使用说明 */}
      <Card>
        <CardHeader>
          <CardTitle className="flex items-center gap-2">
            <AlertCircle className="w-5 h-5" />
            使用说明
          </CardTitle>
        </CardHeader>
        <CardContent className="space-y-4">
          <div>
            <h4 className="font-medium mb-2">权限控制流程：</h4>
            <ol className="space-y-2 text-sm text-gray-600">
              <li className="flex items-start gap-2">
                <Badge variant="outline" className="mt-0.5">1</Badge>
                <span>创建权限：定义系统中的资源和操作（如：user:read, user:create）</span>
              </li>
              <li className="flex items-start gap-2">
                <Badge variant="outline" className="mt-0.5">2</Badge>
                <span>创建角色：定义系统中的角色（如：管理员、普通用户）</span>
              </li>
              <li className="flex items-start gap-2">
                <Badge variant="outline" className="mt-0.5">3</Badge>
                <span>分配权限：为角色分配相应的权限</span>
              </li>
              <li className="flex items-start gap-2">
                <Badge variant="outline" className="mt-0.5">4</Badge>
                <span>创建菜单：定义系统菜单结构，支持页面和按钮级权限</span>
              </li>
              <li className="flex items-start gap-2">
                <Badge variant="outline" className="mt-0.5">5</Badge>
                <span>分配菜单：为角色分配可访问的菜单</span>
              </li>
              <li className="flex items-start gap-2">
                <Badge variant="outline" className="mt-0.5">6</Badge>
                <span>分配角色：为用户分配角色，用户将继承角色的所有权限和菜单</span>
              </li>
            </ol>
          </div>
          
          <div>
            <h4 className="font-medium mb-2">权限代码格式：</h4>
            <div className="bg-gray-50 p-3 rounded-lg text-sm">
              <code>resource:action</code>
              <div className="mt-2 text-gray-600">
                示例：
                <ul className="mt-1 space-y-1">
                  <li>• <code>user:read</code> - 查看用户</li>
                  <li>• <code>user:create</code> - 创建用户</li>
                  <li>• <code>role:assign_permission</code> - 为角色分配权限</li>
                </ul>
              </div>
            </div>
          </div>

          <div>
            <h4 className="font-medium mb-2">菜单类型说明：</h4>
            <div className="grid md:grid-cols-4 gap-4 text-sm">
              <div className="flex items-center gap-2">
                <Badge variant="outline">0</Badge>
                <span>分组 - 菜单分组</span>
              </div>
              <div className="flex items-center gap-2">
                <Badge variant="outline">1</Badge>
                <span>目录 - 菜单目录</span>
              </div>
              <div className="flex items-center gap-2">
                <Badge variant="outline">2</Badge>
                <span>菜单 - 页面菜单</span>
              </div>
              <div className="flex items-center gap-2">
                <Badge variant="outline">3</Badge>
                <span>按钮 - 按钮权限</span>
              </div>
            </div>
          </div>
        </CardContent>
      </Card>
    </div>
  );
};

export default RBACDashboard;
