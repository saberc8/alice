import apiClient from "../apiClient";
import type { APIResponse, PaginationReq, PaginationData } from "@/types/api";
import type { Menu, MenuTree } from "@/types/entity";

// 菜单创建请求
export interface CreateMenuReq {
  parent_id?: string;
  name: string;
  code: string;
  path?: string;
  type: number; // 0:分组 1:目录 2:菜单 3:按钮
  order?: number;
  status?: "active" | "inactive";
  meta?: Record<string, any>;
}

// 菜单更新请求
export interface UpdateMenuReq extends CreateMenuReq {}

// 菜单列表查询参数
export interface MenuListReq extends PaginationReq {
  name?: string;
  code?: string;
  type?: number;
  status?: "active" | "inactive";
  parent_id?: string;
}

export enum MenuApi {
  Menus = "/api/v1/menus",
  MenuDetail = "/api/v1/menus/:id",
  MenuTree = "/api/v1/menus/tree",
  GetUserMenus = "/api/v1/users/:user_id/menus",
  GetUserMenuTree = "/api/v1/users/:user_id/menus/tree",
}

// 创建菜单
const createMenu = (data: CreateMenuReq) => 
  apiClient.post<APIResponse<Menu>>({ 
    url: MenuApi.Menus, 
    data 
  });

// 获取菜单列表
const getMenus = (params?: MenuListReq) => 
  apiClient.get<APIResponse<PaginationData<Menu>>>({ 
    url: MenuApi.Menus, 
    params 
  });

// 获取菜单详情
const getMenuById = (id: string) => 
  apiClient.get<APIResponse<Menu>>({ 
    url: MenuApi.MenuDetail.replace(':id', id) 
  });

// 更新菜单
const updateMenu = (id: string, data: UpdateMenuReq) => 
  apiClient.put<APIResponse<Menu>>({ 
    url: MenuApi.MenuDetail.replace(':id', id), 
    data 
  });

// 删除菜单
const deleteMenu = (id: string) => 
  apiClient.delete<APIResponse<void>>({ 
    url: MenuApi.MenuDetail.replace(':id', id) 
  });

// 获取菜单树
const getMenuTree = () => 
  apiClient.get<APIResponse<MenuTree[]>>({ 
    url: MenuApi.MenuTree 
  });

// 获取用户菜单
const getUserMenus = (userId: string) => 
  apiClient.get<APIResponse<Menu[]>>({ 
    url: MenuApi.GetUserMenus.replace(':user_id', userId) 
  });

// 获取用户菜单树
const getUserMenuTree = (userId: string) => 
  apiClient.get<APIResponse<MenuTree[]>>({ 
    url: MenuApi.GetUserMenuTree.replace(':user_id', userId) 
  });

// 兼容旧接口
const getMenuList = () => getMenuTree();

export default {
  createMenu,
  getMenus,
  getMenuById,
  updateMenu,
  deleteMenu,
  getMenuTree,
  getUserMenus,
  getUserMenuTree,
  getMenuList,
};
