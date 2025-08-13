import { GLOBAL_CONFIG } from "@/global-config";
import DashboardLayout from "@/layouts/dashboard";
import LoginAuthGuard from "@/routes/components/login-auth-guard";
import { Navigate, type RouteObject } from "react-router";
import { useBackendDashboardRoutes } from "./backend";
import { useMemo } from "react";
import { DynamicRouteResolver } from "./components/DynamicRouteResolver";

/**
 * 动态仪表板路由生成器 Hook
 * 基于后端菜单数据动态生成路由，所有组件都使用懒加载
 * 
 * 使用示例：
 * ```tsx
 * const MyRouterComponent = () => {
 *   const routes = useDashboardRoutes();
 *   return <Routes>{routes}</Routes>;
 * };
 * ```
 */
export const useDashboardRoutes = (): RouteObject[] => {
	// 获取基于后端菜单数据的动态路由
	const backendRoutes = useBackendDashboardRoutes();

	return useMemo(() => [
		{
			path: "/",
			element: (
				<LoginAuthGuard>
					<DashboardLayout />
				</LoginAuthGuard>
			),
			children: [
				// 默认重定向到工作台
				{ index: true, element: <Navigate to={GLOBAL_CONFIG.defaultRoute} replace /> },
				// 动态生成的路由（基于后端菜单数据）
				...backendRoutes,
			],
		},
	], [backendRoutes]);
};

/**
 * 静态路由配置（用于兼容现有的路由系统）
 * 
 * 使用通配符路由来捕获所有路径，然后通过 DynamicRouteResolver 
 * 根据用户菜单数据动态解析和渲染对应的组件
 */
export const dashboardRoutes: RouteObject[] = [
	{
		path: "/",
		element: (
			<LoginAuthGuard>
				<DashboardLayout />
			</LoginAuthGuard>
		),
		children: [
			// 默认重定向到工作台
			{ index: true, element: <Navigate to={GLOBAL_CONFIG.defaultRoute} replace /> },
			// 通配符路由：捕获所有未匹配的路径，动态解析路由
			{ path: "*", element: <DynamicRouteResolver /> },
		],
	},
];
