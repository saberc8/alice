import { useMemo } from "react";
import { type RouteObject } from "react-router";
import { useDashboardRoutes } from "../sections/dashboard";
import { authRoutes } from "../sections/auth";
import { mainRoutes } from "../sections/main";
import { Navigate } from "react-router";

/**
 * 动态路由生成 Hook
 * 根据用户权限和菜单数据动态生成完整的路由配置
 */
export const useDynamicRoutes = (): RouteObject[] => {
	const dashboardRoutes = useDashboardRoutes();

	const routes = useMemo<RouteObject[]>(() => [
		// Root redirect
		{ path: "/", element: <Navigate to="/auth/login" replace /> },
		// Login redirect for backward compatibility
		{ path: "/login", element: <Navigate to="/auth/login" replace /> },
		// Auth
		...authRoutes,
		// Dashboard (动态生成)
		...dashboardRoutes,
		// Main
		...mainRoutes,
		// No Match
		{ path: "*", element: <Navigate to="/404" replace /> },
	], [dashboardRoutes]);

	return routes;
};
