import type { ResultStatus } from "./enum";

export interface Result<T = unknown> {
	status: ResultStatus;
	message: string;
	data: T;
}

// 后端标准API响应格式
export interface APIResponse<T = any> {
	code: number;
	message: string;
	data?: T;
}

// 分页请求参数
export interface PaginationReq {
	page?: number;
	page_size?: number;
}

// 分页响应数据
export interface PaginationData<T> {
	list: T[];
	total: number;
	page: number;
	page_size: number;
}
