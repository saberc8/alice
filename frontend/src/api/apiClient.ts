import { GLOBAL_CONFIG } from "@/global-config";
import userStore from "@/store/userStore";
import axios, { type AxiosRequestConfig, type AxiosError, type AxiosResponse } from "axios";
import { toast } from "sonner";

const axiosInstance = axios.create({
	baseURL: GLOBAL_CONFIG.apiBaseUrl,
	timeout: 50000,
	headers: { "Content-Type": "application/json;charset=utf-8" },
});

axiosInstance.interceptors.request.use(
	(config) => {
		// 从userStore获取token
		const token = userStore.getState().userToken.accessToken;
		if (token) {
			config.headers.Authorization = `Bearer ${token}`;
		}
		return config;
	},
	(error) => Promise.reject(error),
);

axiosInstance.interceptors.response.use(
	(res: AxiosResponse) => {
		// 检查标准响应格式
		const responseData = res.data;
		
		// 如果是标准API响应格式 (有code字段)
		if (responseData && typeof responseData.code === 'number') {
			// 检查是否成功
			if (responseData.code === 200) {
				// 成功响应，返回完整的响应对象 (包含 code, message, data)
				return responseData;
			} else {
				// 业务错误，抛出错误
				const error = new Error(responseData.message || 'API request failed');
				return Promise.reject(error);
			}
		}
		
		// 兼容旧格式或其他格式的响应
		return responseData;
	},
	(error: AxiosError) => {
		const { response, message } = error || {};
		let errMsg = message || "Network error";
		
		// 尝试从响应中获取错误信息
		if (response?.data) {
			const data = response.data as any;
			if (data.message) {
				errMsg = data.message;
			} else if (data.error) {
				errMsg = data.error;
			}
		}
		
		toast.error(errMsg, { position: "top-center" });
		
		if (response?.status === 401) {
			userStore.getState().actions.clearUserInfoAndToken();
		}
		
		return Promise.reject(new Error(errMsg));
	},
);

class APIClient {
	get<T = unknown>(config: AxiosRequestConfig): Promise<T> {
		return this.request<T>({ ...config, method: "GET" });
	}
	post<T = unknown>(config: AxiosRequestConfig): Promise<T> {
		return this.request<T>({ ...config, method: "POST" });
	}
	put<T = unknown>(config: AxiosRequestConfig): Promise<T> {
		return this.request<T>({ ...config, method: "PUT" });
	}
	delete<T = unknown>(config: AxiosRequestConfig): Promise<T> {
		return this.request<T>({ ...config, method: "DELETE" });
	}
	request<T = unknown>(config: AxiosRequestConfig): Promise<T> {
		return axiosInstance.request<any, T>(config);
	}
}

export default new APIClient();
