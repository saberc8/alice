import apiClient from "../apiClient";
import type { APIResponse, PaginationReq, PaginationData } from "@/types/api";
import type { Role, Permission, Menu } from "@/types/entity";

// 角色创建请求
export interface CreateRoleReq {
  name: string;
  code: string;
  description?: string;
  status?: "active" | "inactive";
}

// 角色更新请求
export interface UpdateRoleReq extends CreateRoleReq {}

// 角色权限分配请求
export interface AssignRolePermissionsReq {
  permission_ids: string[];
}

// 角色菜单分配请求
export interface AssignRoleMenusReq {
  menu_ids: string[];
}

// 角色列表查询参数
export interface RoleListReq extends PaginationReq {
  name?: string;
  code?: string;
  status?: "active" | "inactive";
}

export enum RoleApi {
  Roles = "/api/v1/roles",
  RoleDetail = "/api/v1/roles/:id",
  RolePermissions = "/api/v1/roles/:id/permissions",
  RoleMenus = "/api/v1/roles/:id/menus",
}

// 创建角色
const createRole = (data: CreateRoleReq) => 
  apiClient.post<APIResponse<Role>>({ 
    url: RoleApi.Roles, 
    data 
  });

// 获取角色列表
const getRoles = (params?: RoleListReq) => 
  apiClient.get<APIResponse<PaginationData<Role>>>({ 
    url: RoleApi.Roles, 
    params 
  });

// 获取角色详情
const getRoleById = (id: string) => 
  apiClient.get<APIResponse<Role>>({ 
    url: RoleApi.RoleDetail.replace(':id', id) 
  });

// 更新角色
const updateRole = (id: string, data: UpdateRoleReq) => 
  apiClient.put<APIResponse<Role>>({ 
    url: RoleApi.RoleDetail.replace(':id', id), 
    data 
  });

// 删除角色
const deleteRole = (id: string) => 
  apiClient.delete<APIResponse<void>>({ 
    url: RoleApi.RoleDetail.replace(':id', id) 
  });

// 获取角色权限
const getRolePermissions = (roleId: string) => 
  apiClient.get<APIResponse<Permission[]>>({ 
    url: RoleApi.RolePermissions.replace(':id', roleId) 
  });

// 分配角色权限
const assignRolePermissions = (roleId: string, data: AssignRolePermissionsReq) => 
  apiClient.post<APIResponse<void>>({ 
    url: RoleApi.RolePermissions.replace(':id', roleId), 
    data 
  });

// 移除角色权限
const removeRolePermissions = (roleId: string, data: AssignRolePermissionsReq) => 
  apiClient.delete<APIResponse<void>>({ 
    url: RoleApi.RolePermissions.replace(':id', roleId), 
    data 
  });

// 获取角色菜单
const getRoleMenus = (roleId: string) => 
  apiClient.get<APIResponse<Menu[]>>({ 
    url: RoleApi.RoleMenus.replace(':id', roleId) 
  });

// 分配角色菜单
const assignRoleMenus = (roleId: string, data: AssignRoleMenusReq) => 
  apiClient.post<APIResponse<void>>({ 
    url: RoleApi.RoleMenus.replace(':id', roleId), 
    data 
  });

// 移除角色菜单
const removeRoleMenus = (roleId: string, data: AssignRoleMenusReq) => 
  apiClient.delete<APIResponse<void>>({ 
    url: RoleApi.RoleMenus.replace(':id', roleId), 
    data 
  });

export default {
  createRole,
  getRoles,
  getRoleById,
  updateRole,
  deleteRole,
  getRolePermissions,
  assignRolePermissions,
  removeRolePermissions,
  getRoleMenus,
  assignRoleMenus,
  removeRoleMenus,
};
