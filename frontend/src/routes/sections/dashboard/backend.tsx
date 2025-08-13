import type { MenuMetaInfo, MenuTree } from "@/types/entity";
import { PermissionType } from "@/types/enum";
import { convertFlatToTree } from "@/utils/tree";
import type { RouteObject } from "react-router";
import { Navigate } from "react-router";
import { Component } from "./utils";
import useUserStore from "@/store/userStore";
import { useMemo } from "react";

/**
 * get route path from menu path and parent path
 * @param menuPath '/a/b/c'
 * @param parentPath '/a/b'
 * @returns '/c'
 *
 * @example
 * getRoutePath('/a/b/c', '/a/b') // '/c'
 */
const getRoutePath = (menuPath?: string, parentPath?: string) => {
	const menuPathArr = menuPath?.split("/").filter(Boolean) || [];
	const parentPathArr = parentPath?.split("/").filter(Boolean) || [];

	// remove parentPath items from menuPath
	const result = menuPathArr.slice(parentPathArr.length).join("/");
	console.log(`getRoutePath: menuPath="${menuPath}", parentPath="${parentPath}", result="${result}"`);
	return result;
};

/**
 * generate props for menu component
 * @param metaInfo
 * @returns
 */
const generateProps = (metaInfo: MenuMetaInfo) => {
	const props: any = {};
	if (metaInfo.externalLink) {
		props.src = metaInfo.externalLink?.toString() || "";
	}
	return props;
};

/**
 * convert menu to route
 * @param items
 * @param parent
 * @returns
 */
const convertToRoute = (items: MenuTree[], parent?: MenuTree): RouteObject[] => {
	const routes: RouteObject[] = [];

	const processItem = (item: MenuTree) => {
		console.log(`处理菜单项: ${item.name}, type: ${item.type}, path: ${item.path}, parent: ${parent?.name || 'null'}, 是否有children: ${(item.children || []).length > 0}, 是否有component: ${!!item.meta?.component}`);
		
		// if group, process children
		if (item.type === PermissionType.GROUP) {
			console.log(`处理分组: ${item.name}, 直接处理 ${(item.children || []).length} 个子项`);
			for (const child of item.children || []) {
				processItem(child);
			}
		}

		// if catalogue, process children
		if (item.type === PermissionType.CATALOGUE) {
			const children = item.children || [];
			if (children.length > 0) {
				// 查找第一个有组件的子菜单项
				const firstChildWithComponent = children.find(child => child.meta?.component || child.component);
				console.log(`创建目录路由: ${getRoutePath(item.path, parent?.path)} (${children.length} 个子项, 第一个有组件的: ${firstChildWithComponent?.name})`);
				
				routes.push({
					path: getRoutePath(item.path, parent?.path),
					children: [
						// 如果有第一个可访问的子页面，则重定向到它
						...(firstChildWithComponent && firstChildWithComponent.path ? [{
							index: true,
							element: <Navigate to={getRoutePath(firstChildWithComponent.path, item.path)} replace />,
						}] : []),
						// 递归处理所有子菜单
						...convertToRoute(children, item),
					],
				});
			}
		}

		if (item.type === PermissionType.MENU) {
			const componentPath = item.meta?.component || item.component;
			const hasChildren = (item.children || []).length > 0;
			if (componentPath) {
				const props = generateProps(item);
				console.log(`创建页面路由: ${getRoutePath(item.path, parent?.path)} -> ${componentPath}`);
				routes.push({
					path: getRoutePath(item.path, parent?.path),
					element: Component(componentPath, props),
				});
			}
		}
	};

	for (const item of items) {
		processItem(item);
	}
	return routes;
};

/**
 * Hook to get dynamic routes from user store with memoization
 */
export const useBackendDashboardRoutes = (): RouteObject[] => {
	const userMenuTree = useUserStore((state) => state.userMenuTree);
	
	return useMemo(() => {
		return convertToRoute(userMenuTree);
	}, [userMenuTree]);
};

/**
 * Get static routes based on current user menu tree (for compatibility)
 * Note: This will only work after user store is initialized
 */
export const getBackendDashboardRoutes = (): RouteObject[] => {
	const userMenuTree = useUserStore.getState().userMenuTree;
	
	if (!userMenuTree || userMenuTree.length === 0) {
		console.log("用户菜单数据为空，返回基础路由结构");
		// 返回一些基础路由以防止 404 错误
		// 这些路由会在用户登录后被动态路由替换
		return [
			{
				path: "*",
				element: <div style={{ padding: '20px', textAlign: 'center' }}>
					<p>正在加载菜单数据222...</p>
					<p>如果长时间显示此消息，请刷新页面或联系管理员</p>
					<p>调试信息</p>
					<p>当前路径: {window.location.pathname}</p>
					<p>动态路由数量: 0</p>
				</div>
			}
		];
	}
	
	console.log("用户菜单数据:", userMenuTree);
	const routes = convertToRoute(userMenuTree);
	console.log(`从用户菜单数据生成了 ${routes.length} 个动态路由`);
	console.log("生成的路由详细信息:");
	routes.forEach((route, index) => {
		console.log(`路由 ${index + 1}:`, {
			path: route.path,
			hasElement: !!route.element,
			hasChildren: !!(route.children && route.children.length > 0),
			childrenCount: route.children?.length || 0,
			children: route.children?.map(child => ({
				path: child.path,
				index: child.index,
				hasElement: !!child.element
			})) || []
		});
	});
	return routes;
};
