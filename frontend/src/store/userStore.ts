import { useMutation } from "@tanstack/react-query";
import { create } from "zustand";
import { createJSONStorage, persist } from "zustand/middleware";

import userService, { type SignInReq } from "@/api/services/userService";

import { toast } from "sonner";
import type { UserInfo, UserToken } from "#/entity";
import { StorageEnum } from "#/enum";

type UserStore = {
	userInfo: Partial<UserInfo>;
	userToken: UserToken;

	actions: {
		setUserInfo: (userInfo: UserInfo) => void;
		setUserToken: (token: UserToken) => void;
		clearUserInfoAndToken: () => void;
	};
};

const useUserStore = create<UserStore>()(
	persist(
		(set) => ({
			userInfo: {},
			userToken: {},
			actions: {
				setUserInfo: (userInfo) => {
					set({ userInfo });
				},
				setUserToken: (userToken) => {
					set({ userToken });
				},
				clearUserInfoAndToken() {
					set({ userInfo: {}, userToken: {} });
				},
			},
		}),
		{
			name: "userStore", // name of the item in the storage (must be unique)
			storage: createJSONStorage(() => localStorage), // (optional) by default, 'localStorage' is used
			partialize: (state) => ({
				[StorageEnum.UserInfo]: state.userInfo,
				[StorageEnum.UserToken]: state.userToken,
			}),
		},
	),
);

export const useUserInfo = () => useUserStore((state) => state.userInfo);
export const useUserToken = () => useUserStore((state) => state.userToken);
export const useUserPermissions = () => useUserStore((state) => state.userInfo.permissions || []);
export const useUserRoles = () => useUserStore((state) => state.userInfo.roles || []);
export const useUserActions = () => useUserStore((state) => state.actions);

export const useSignIn = () => {
	const { setUserToken, setUserInfo } = useUserActions();

	const signInMutation = useMutation({
		mutationFn: userService.signin,
	});

	const signIn = async (data: SignInReq) => {
		try {
			const res = await signInMutation.mutateAsync(data);
			
			// 检查响应格式
			if (res.code !== 200) {
				throw new Error(res.message || 'Login failed');
			}
			
			const { token } = res.data!;
			// 设置token
			setUserToken({ accessToken: token });
			
			// 获取用户资料
			try {
				const profileRes = await userService.getProfile();
				if (profileRes.code === 200 && profileRes.data) {
					setUserInfo({
						id: profileRes.data.id.toString(),
						username: profileRes.data.username,
						email: profileRes.data.email,
					});
				} else {
					console.error("获取用户资料失败:", profileRes.message);
					// 即使获取资料失败，也保持登录状态，只设置基本信息
					setUserInfo({ username: data.username });
				}
			} catch (profileError) {
				console.error("获取用户资料失败:", profileError);
				// 即使获取资料失败，也保持登录状态，只设置基本信息
				setUserInfo({ username: data.username });
			}
		} catch (err: any) {
			toast.error(err.message || 'Login failed', {
				position: "top-center",
			});
			throw err;
		}
	};

	return signIn;
};

export default useUserStore;
