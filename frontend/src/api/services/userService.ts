import apiClient from "../apiClient";

import type { UserInfo, UserToken } from "#/entity";

export interface SignInReq {
	username: string;
	password: string;
}

export interface SignUpReq extends SignInReq {
	email: string;
}

// 标准API响应格式
export interface APIResponse<T = any> {
	code: number;
	message: string;
	data?: T;
}

// 登录响应数据
export interface LoginData {
	token: string;
}

// 注册响应数据
export interface RegisterData {
	user: {
		id: number;
		username: string;
		email: string;
	};
	token: string;
}

// 用户资料数据
export interface ProfileData {
	id: number;
	username: string;
	email: string;
}

// 后端登录响应格式
export type SignInRes = APIResponse<LoginData>;

// 后端注册响应格式  
export type SignUpRes = APIResponse<RegisterData>;

// 后端用户资料响应格式
export type ProfileRes = APIResponse<ProfileData>;

export enum UserApi {
	SignIn = "/api/v1/users/login",
	SignUp = "/api/v1/users/register", 
	Profile = "/api/v1/users/profile",
	Logout = "/auth/logout",
	Refresh = "/auth/refresh",
	User = "/user",
}

const signin = (data: SignInReq) => apiClient.post<SignInRes>({ url: UserApi.SignIn, data });
const signup = (data: SignUpReq) => apiClient.post<SignUpRes>({ url: UserApi.SignUp, data });
const getProfile = () => apiClient.get<ProfileRes>({ url: UserApi.Profile });
const logout = () => apiClient.get({ url: UserApi.Logout });
const findById = (id: string) => apiClient.get<UserInfo[]>({ url: `${UserApi.User}/${id}` });

export default {
	signin,
	signup,
	getProfile,
	findById,
	logout,
};
