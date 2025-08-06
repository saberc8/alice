import React from 'react';
import { usePermissionCheck, useRoleCheck } from '@/store/userStore';

interface PermissionWrapperProps {
  children: React.ReactNode;
  // 权限控制方式（三选一）
  resource?: string;
  action?: string;
  permissionCode?: string;
  roleCode?: string;
  // 无权限时显示的内容
  fallback?: React.ReactNode;
  // 是否隐藏（默认false，无权限时显示fallback）
  hide?: boolean;
}

/**
 * 权限控制组件
 * 用于页面级和按钮级权限控制
 * 
 * @example
 * // 资源+操作方式
 * <PermissionWrapper resource="user" action="delete">
 *   <Button>删除用户</Button>
 * </PermissionWrapper>
 * 
 * // 权限代码方式
 * <PermissionWrapper permissionCode="user:delete">
 *   <Button>删除用户</Button>
 * </PermissionWrapper>
 * 
 * // 角色方式
 * <PermissionWrapper roleCode="admin">
 *   <AdminPanel />
 * </PermissionWrapper>
 * 
 * // 无权限时显示提示
 * <PermissionWrapper resource="user" action="create" fallback={<div>无权限</div>}>
 *   <Button>创建用户</Button>
 * </PermissionWrapper>
 * 
 * // 无权限时隐藏
 * <PermissionWrapper resource="user" action="delete" hide>
 *   <Button>删除用户</Button>
 * </PermissionWrapper>
 */
export const PermissionWrapper: React.FC<PermissionWrapperProps> = ({
  children,
  resource,
  action,
  permissionCode,
  roleCode,
  fallback = null,
  hide = false,
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
  if (hide) {
    return null;
  }

  return <>{fallback}</>;
};

// Hook方式的权限检查
export const useHasPermission = (resource: string, action: string) => {
  const { hasPermission } = usePermissionCheck();
  return hasPermission(resource, action);
};

export const useHasPermissionByCode = (code: string) => {
  const { hasPermissionByCode } = usePermissionCheck();
  return hasPermissionByCode(code);
};

export const useHasRole = (roleCode: string) => {
  const { hasRole } = useRoleCheck();
  return hasRole(roleCode);
};

export default PermissionWrapper;
