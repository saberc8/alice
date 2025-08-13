import { useMemo } from "react";
import { Navigate, Outlet } from "react-router";
import { useBackendDashboardRoutes } from "../backend";
import { GLOBAL_CONFIG } from "@/global-config";

/**
 * 动态路由出口组件
 * 这个组件会根据当前用户的菜单数据动态渲染匹配的路由
 */
export const DynamicRoutesOutlet = () => {
	const dynamicRoutes = useBackendDashboardRoutes();
	
	// 创建一个包含动态路由的路由配置
	const routeConfig = useMemo(() => [
		// 默认重定向到工作台
		{ index: true, element: <Navigate to={GLOBAL_CONFIG.defaultRoute} replace /> },
		// 动态生成的路由
		...dynamicRoutes,
	], [dynamicRoutes]);

	// 由于我们在 DashboardLayout 中，直接渲染 Outlet
	// React Router 会自动匹配子路由
	return <Outlet />;
};
