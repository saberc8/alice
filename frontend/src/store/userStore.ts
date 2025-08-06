import { useMutation } from "@tanstack/react-query";
import { create } from "zustand";
import { createJSONStorage, persist } from "zustand/middleware";

import userService, { type SignInReq } from "@/api/services/userService";
import menuService from "@/api/services/menuService";
import permissionService from "@/api/services/permissionService";

import { toast } from "sonner";
import type { UserInfo, UserToken, Role, Permission, MenuTree } from "#/entity";
import { StorageEnum } from "#/enum";

type UserStore = {
	userInfo: Partial<UserInfo>;
	userToken: UserToken;
	userRoles: Role[];
	userPermissions: Permission[];
	userMenuTree: MenuTree[];

	actions: {
		setUserInfo: (userInfo: UserInfo) => void;
		setUserToken: (token: UserToken) => void;
		setUserRoles: (roles: Role[]) => void;
		setUserPermissions: (permissions: Permission[]) => void;
		setUserMenuTree: (menuTree: MenuTree[]) => void;
		clearUserInfoAndToken: () => void;
		loadUserData: (userId: string) => Promise<void>;
	};
};

const useUserStore = create<UserStore>()(
	persist(
		(set, get) => ({
			userInfo: {},
			userToken: {},
			userRoles: [],
			userPermissions: [],
			userMenuTree: [],
			actions: {
				setUserInfo: (userInfo) => {
					set({ userInfo });
				},
				setUserToken: (userToken) => {
					set({ userToken });
				},
				setUserRoles: (userRoles) => {
					set({ userRoles });
				},
				setUserPermissions: (userPermissions) => {
					set({ userPermissions });
				},
				setUserMenuTree: (userMenuTree) => {
					set({ userMenuTree });
				},
				clearUserInfoAndToken() {
					set({ 
						userInfo: {}, 
						userToken: {},
						userRoles: [],
						userPermissions: [],
						userMenuTree: [],
					});
				},
				// 加载用户完整数据（角色、权限、菜单）
				loadUserData: async (userId: string) => {
					try {
						const { setUserRoles, setUserPermissions, setUserMenuTree } = get().actions;
						
						// 并行获取用户角色、权限和菜单
						const [rolesRes, permissionsRes, menuTreeRes] = await Promise.allSettled([
							userService.getUserRoles(userId),
							userService.getUserPermissions(userId),
							userService.getUserMenuTree(userId),
						]);

						// 处理角色
						if (rolesRes.status === 'fulfilled' && rolesRes.value.code === 200) {
							setUserRoles(rolesRes.value.data || []);
						}

						// 处理权限
						if (permissionsRes.status === 'fulfilled' && permissionsRes.value.code === 200) {
							setUserPermissions(permissionsRes.value.data || []);
						}

						// 处理菜单
						if (menuTreeRes.status === 'fulfilled' && menuTreeRes.value.code === 200) {
							setUserMenuTree(menuTreeRes.value.data || []);
						}
					} catch (error) {
						console.error('加载用户数据失败:', error);
					}
				},
			},
		}),
		{
			name: "userStore", // name of the item in the storage (must be unique)
			storage: createJSONStorage(() => localStorage), // (optional) by default, 'localStorage' is used
			partialize: (state) => ({
				[StorageEnum.UserInfo]: state.userInfo,
				[StorageEnum.UserToken]: state.userToken,
				userRoles: state.userRoles,
				userPermissions: state.userPermissions,
				userMenuTree: state.userMenuTree,
			}),
		},
	),
);

export const useUserInfo = () => useUserStore((state) => state.userInfo);
export const useUserToken = () => useUserStore((state) => state.userToken);
export const useUserPermissions = () => useUserStore((state) => state.userPermissions);
export const useUserRoles = () => useUserStore((state) => state.userRoles);
export const useUserMenuTree = () => useUserStore((state) => state.userMenuTree);
export const useUserActions = () => useUserStore((state) => state.actions);

export const useSignIn = () => {
	const { setUserToken, setUserInfo, loadUserData } = useUserActions();

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
			
			const { token, user } = res.data!;
			
			// 设置token
			setUserToken({ accessToken: token });
			
			// 设置用户基本信息
			if (user) {
				setUserInfo({
					id: user.id,
					username: user.username,
					email: user.email,
					status: user.status,
				});
				
				// 加载用户角色、权限、菜单数据
				await loadUserData(user.id);
			} else {
				// 如果登录响应中没有用户信息，尝试获取用户资料
				try {
					const profileRes = await userService.getProfile();
					if (profileRes.code === 200 && profileRes.data) {
						const userData = profileRes.data;
						setUserInfo({
							id: userData.id,
							username: userData.username,
							email: userData.email,
							status: userData.status,
						});
						
						// 加载用户角色、权限、菜单数据
						await loadUserData(userData.id);
					}
				} catch (profileError) {
					console.error("获取用户资料失败:", profileError);
					// 即使获取资料失败，也保持登录状态，只设置基本信息
					setUserInfo({ username: data.username });
				}
			}
			
			toast.success('登录成功', { position: "top-center" });
		} catch (err: any) {
			toast.error(err.message || 'Login failed', {
				position: "top-center",
			});
			throw err;
		}
	};

	return { signIn, isPending: signInMutation.isPending };
};

// 权限检查 Hook
export const usePermissionCheck = () => {
	const userPermissions = useUserPermissions();
	
	const hasPermission = (resource: string, action: string) => {
		return userPermissions.some(
			permission => permission.resource === resource && permission.action === action
		);
	};
	
	const hasPermissionByCode = (code: string) => {
		return userPermissions.some(permission => permission.code === code);
	};
	
	return { hasPermission, hasPermissionByCode };
};

// 角色检查 Hook
export const useRoleCheck = () => {
	const userRoles = useUserRoles();
	
	const hasRole = (roleCode: string) => {
		return userRoles.some(role => role.code === roleCode);
	};
	
	const hasAnyRole = (roleCodes: string[]) => {
		return roleCodes.some(code => hasRole(code));
	};
	
	return { hasRole, hasAnyRole };
};

export default useUserStore;
