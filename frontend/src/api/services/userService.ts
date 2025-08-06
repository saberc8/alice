import apiClient from "../apiClient";
import type { APIResponse, PaginationReq, PaginationData } from "@/types/api";
import type { User, Role, Permission, MenuTree } from "@/types/entity";

export interface SignInReq {
	username: string;
	password: string;
}

export interface SignUpReq extends SignInReq {
	email: string;
}

// 登录响应数据
export interface LoginData {
	token: string;
	user?: User;
}

// 注册响应数据
export interface RegisterData {
	user: User;
	token: string;
}

// 用户角色分配请求
export interface AssignUserRolesReq {
	role_ids: string[];
}

// 用户列表查询参数
export interface UserListReq extends PaginationReq {
	username?: string;
	email?: string;
	status?: "active" | "inactive";
}

// 创建用户请求
export interface CreateUserReq {
	username: string;
	email: string;
	password: string;
	status?: "active" | "inactive";
}

// 更新用户请求
export interface UpdateUserReq {
	username?: string;
	email?: string;
	password?: string;
	status?: "active" | "inactive";
}

// 后端登录响应格式
export type SignInRes = APIResponse<LoginData>;

// 后端注册响应格式  
export type SignUpRes = APIResponse<RegisterData>;

// 后端用户资料响应格式
export type ProfileRes = APIResponse<User>;

export enum UserApi {
	SignIn = "/api/v1/auth/login",
	SignUp = "/api/v1/auth/register", 
	Profile = "/api/v1/auth/profile",
	Logout = "/api/v1/auth/logout",
	Refresh = "/api/v1/auth/refresh",
	Users = "/api/v1/users",
	UserDetail = "/api/v1/users/:id",
	UserRoles = "/api/v1/users/:user_id/roles",
	UserPermissions = "/api/v1/users/:user_id/permissions",
	UserMenuTree = "/api/v1/users/:user_id/menus/tree",
}

// 认证相关
const signin = (data: SignInReq) => apiClient.post<SignInRes>({ url: UserApi.SignIn, data });
const signup = (data: SignUpReq) => apiClient.post<SignUpRes>({ url: UserApi.SignUp, data });
const getProfile = () => apiClient.get<ProfileRes>({ url: UserApi.Profile });
const logout = () => apiClient.post({ url: UserApi.Logout });

// 用户管理
const getUsers = (params?: UserListReq) => 
	apiClient.get<APIResponse<PaginationData<User>>>({ url: UserApi.Users, params });

const getUserById = (id: string) => 
	apiClient.get<APIResponse<User>>({ url: UserApi.UserDetail.replace(':id', id) });

const createUser = (data: CreateUserReq) => 
	apiClient.post<APIResponse<User>>({ url: UserApi.Users, data });

const updateUser = (id: string, data: UpdateUserReq) => 
	apiClient.put<APIResponse<User>>({ url: UserApi.UserDetail.replace(':id', id), data });

const deleteUser = (id: string) => 
	apiClient.delete<APIResponse<void>>({ url: UserApi.UserDetail.replace(':id', id) });

// 用户角色管理
const getUserRoles = (userId: string) => 
	apiClient.get<APIResponse<Role[]>>({ url: UserApi.UserRoles.replace(':user_id', userId) });

const assignUserRoles = (userId: string, data: AssignUserRolesReq) => 
	apiClient.post<APIResponse<void>>({ url: UserApi.UserRoles.replace(':user_id', userId), data });

const removeUserRoles = (userId: string, data: AssignUserRolesReq) => 
	apiClient.delete<APIResponse<void>>({ url: UserApi.UserRoles.replace(':user_id', userId), data });

// 用户权限查询
const getUserPermissions = (userId: string) => 
	apiClient.get<APIResponse<Permission[]>>({ url: UserApi.UserPermissions.replace(':user_id', userId) });

const getUserMenuTree = (userId: string) => 
	apiClient.get<APIResponse<MenuTree[]>>({ url: UserApi.UserMenuTree.replace(':user_id', userId) });

// 兼容旧接口
const findById = (id: string) => getUserById(id);

export default {
	signin,
	signup,
	getProfile,
	logout,
	getUsers,
	getUserById,
	createUser,
	updateUser,
	deleteUser,
	getUserRoles,
	assignUserRoles,
	removeUserRoles,
	getUserPermissions,
	getUserMenuTree,
	findById,
};
