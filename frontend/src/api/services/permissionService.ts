import apiClient from "../apiClient";
import type { APIResponse, PaginationReq, PaginationData } from "@/types/api";
import type { Permission } from "@/types/entity";

// 权限创建请求
export interface CreatePermissionReq {
  name: string;
  code: string;
  resource: string;
  action: string;
  description?: string;
  status?: "active" | "inactive";
}

// 权限更新请求
export interface UpdatePermissionReq extends CreatePermissionReq {}

// 权限列表查询参数
export interface PermissionListReq extends PaginationReq {
  name?: string;
  code?: string;
  resource?: string;
  action?: string;
  status?: "active" | "inactive";
}

// 权限检查请求
export interface CheckPermissionReq {
  resource: string;
  action: string;
}

// 权限检查响应
export interface CheckPermissionRes {
  has_permission: boolean;
}

export enum PermissionApi {
  Permissions = "/api/v1/permissions",
  PermissionDetail = "/api/v1/permissions/:id",
  CheckUserPermission = "/api/v1/users/:user_id/permissions/check",
  GetUserPermissions = "/api/v1/users/:user_id/permissions",
}

// 创建权限
const createPermission = (data: CreatePermissionReq) => 
  apiClient.post<APIResponse<Permission>>({ 
    url: PermissionApi.Permissions, 
    data 
  });

// 获取权限列表
const getPermissions = (params?: PermissionListReq) => 
  apiClient.get<APIResponse<PaginationData<Permission>>>({ 
    url: PermissionApi.Permissions, 
    params 
  });

// 获取权限详情
const getPermissionById = (id: string) => 
  apiClient.get<APIResponse<Permission>>({ 
    url: PermissionApi.PermissionDetail.replace(':id', id) 
  });

// 更新权限
const updatePermission = (id: string, data: UpdatePermissionReq) => 
  apiClient.put<APIResponse<Permission>>({ 
    url: PermissionApi.PermissionDetail.replace(':id', id), 
    data 
  });

// 删除权限
const deletePermission = (id: string) => 
  apiClient.delete<APIResponse<void>>({ 
    url: PermissionApi.PermissionDetail.replace(':id', id) 
  });

// 获取用户权限
const getUserPermissions = (userId: string) => 
  apiClient.get<APIResponse<Permission[]>>({ 
    url: PermissionApi.GetUserPermissions.replace(':user_id', userId) 
  });

// 检查用户权限
const checkUserPermission = (userId: string, params: CheckPermissionReq) => 
  apiClient.get<APIResponse<CheckPermissionRes>>({ 
    url: PermissionApi.CheckUserPermission.replace(':user_id', userId),
    params 
  });

export default {
  createPermission,
  getPermissions,
  getPermissionById,
  updatePermission,
  deletePermission,
  getUserPermissions,
  checkUserPermission,
};
