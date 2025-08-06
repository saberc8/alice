import React from 'react';
import { Navigate } from 'react-router-dom';
import { usePermissionCheck, useRoleCheck } from '@/store/userStore';
import { Alert, AlertDescription } from '@/ui/alert';
import { AlertCircle } from 'lucide-react';

interface PermissionGuardProps {
  children: React.ReactNode;
  resource?: string;
  action?: string;
  permissionCode?: string;
  roleCode?: string;
  redirectTo?: string;
  showForbidden?: boolean;
}

/**
 * 路由权限守卫组件
 * 用于保护需要特定权限的路由
 * 
 * @example
 * // 资源+操作方式
 * <PermissionGuard resource="user" action="read">
 *   <UserList />
 * </PermissionGuard>
 * 
 * // 权限代码方式
 * <PermissionGuard permissionCode="user:read">
 *   <UserList />
 * </PermissionGuard>
 * 
 * // 角色方式
 * <PermissionGuard roleCode="admin">
 *   <AdminPanel />
 * </PermissionGuard>
 * 
 * // 重定向到指定页面
 * <PermissionGuard resource="user" action="read" redirectTo="/403">
 *   <UserList />
 * </PermissionGuard>
 * 
 * // 显示403页面
 * <PermissionGuard resource="user" action="read" showForbidden>
 *   <UserList />
 * </PermissionGuard>
 */
export const PermissionGuard: React.FC<PermissionGuardProps> = ({
  children,
  resource,
  action,
  permissionCode,
  roleCode,
  redirectTo,
  showForbidden = false,
}) => {
  const { hasPermission, hasPermissionByCode } = usePermissionCheck();
  const { hasRole } = useRoleCheck();

  let hasAccess = false;

  // 检查权限
  if (resource && action) {
    hasAccess = hasPermission(resource, action);
  } else if (permissionCode) {
    hasAccess = hasPermissionByCode(permissionCode);
  } else if (roleCode) {
    hasAccess = hasRole(roleCode);
  } else {
    // 如果没有指定任何权限条件，默认有权限
    hasAccess = true;
  }

  if (hasAccess) {
    return <>{children}</>;
  }

  // 无权限时的处理
  if (redirectTo) {
    return <Navigate to={redirectTo} replace />;
  }

  if (showForbidden) {
    return (
      <div className="min-h-screen flex items-center justify-center">
        <div className="max-w-md w-full">
          <Alert variant="destructive">
            <AlertCircle className="h-4 w-4" />
            <AlertDescription>
              <div className="space-y-2">
                <h3 className="font-semibold">访问被拒绝</h3>
                <p>您没有访问此页面的权限，请联系管理员。</p>
                {(resource && action) && (
                  <p className="text-sm">所需权限：{resource}:{action}</p>
                )}
                {permissionCode && (
                  <p className="text-sm">所需权限：{permissionCode}</p>
                )}
                {roleCode && (
                  <p className="text-sm">所需角色：{roleCode}</p>
                )}
              </div>
            </AlertDescription>
          </Alert>
        </div>
      </div>
    );
  }

  // 默认重定向到403页面
  return <Navigate to="/error/403" replace />;
};

export default PermissionGuard;
