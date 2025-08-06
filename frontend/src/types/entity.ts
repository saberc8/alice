import type { NavItemDataProps } from "@/components/nav/types";
import type { BasicStatus, PermissionType } from "./enum";

export interface UserToken {
	accessToken?: string;
	refreshToken?: string;
}

export interface UserInfo {
	id: string;
	email: string;
	username: string;
	password?: string;
	avatar?: string;
	roles?: Role[];
	status?: BasicStatus;
	permissions?: Permission[];
	menu?: MenuTree[];
}

export interface CommonOptions {
	status?: BasicStatus;
	desc?: string;
	createdAt?: string;
	updatedAt?: string;
}

export interface User extends CommonOptions {
	id: string; // uuid
	username: string;
	password?: string;
	email: string;
	phone?: string;
	avatar?: string;
	password_hash?: string;
	created_at?: string;
	updated_at?: string;
}

// 角色实体 - 匹配后端角色表结构
export interface Role extends CommonOptions {
	id: string; // uuid
	name: string;
	code: string;
	description?: string;
	created_at?: string;
	updated_at?: string;
}

// 权限实体 - 匹配后端权限表结构
export interface Permission extends CommonOptions {
	id: string; // uuid
	name: string;
	code: string; // resource:action  example: "user:read"
	resource: string;
	action: string;
	description?: string;
	created_at?: string;
	updated_at?: string;
}

// 菜单实体 - 匹配后端菜单表结构
export interface Menu extends CommonOptions, MenuMetaInfo {
	id: string; // uuid
	parent_id?: string;
	name: string;
	code: string;
	path?: string;
	order?: number;
	type: PermissionType; // 0:分组 1:目录 2:菜单 3:按钮
	meta?: Record<string, any>;
	created_at?: string;
	updated_at?: string;
}

export type MenuMetaInfo = Partial<Pick<NavItemDataProps, "path" | "icon" | "caption" | "info" | "disabled" | "auth" | "hidden">> & {
	externalLink?: URL;
	component?: string;
};

export type MenuTree = Menu & {
	children?: MenuTree[];
};

// RBAC关联关系实体
export interface UserRole {
	user_id: string;
	role_id: string;
	created_at?: string;
}

export interface RolePermission {
	role_id: string;
	permission_id: string;
	created_at?: string;
}

export interface RoleMenu {
	role_id: string;
	menu_id: string;
	created_at?: string;
}
