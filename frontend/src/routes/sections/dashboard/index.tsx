import { GLOBAL_CONFIG } from "@/global-config";
import DashboardLayout from "@/layouts/dashboard";
import LoginAuthGuard from "@/routes/components/login-auth-guard";
import { Navigate, type RouteObject } from "react-router";
import { useBackendDashboardRoutes } from "./backend";
import { frontendDashboardRoutes } from "./frontend";
import { useMemo } from "react";

const useDashboardRoutes = (): RouteObject[] => {
	const backendRoutes = useBackendDashboardRoutes();
	
	return useMemo(() => {
		if (GLOBAL_CONFIG.routerMode === "frontend") {
			return frontendDashboardRoutes;
		}
		return backendRoutes;
	}, [backendRoutes]);
// 由于路由需要静态导出，我们需要创建一个动态路由包装器
// 但是为了简化，我们暂时保持静态路由结构
// 在实际应用中，动态路由应该在应用启动时根据用户权限生成
}
export const dashboardRoutes: RouteObject[] = [
	{
		element: (
			<LoginAuthGuard>
				<DashboardLayout />
			</LoginAuthGuard>
		),
		children: [
			{ index: true, element: <Navigate to={GLOBAL_CONFIG.defaultRoute} replace /> }, 
			...frontendDashboardRoutes // 使用前端路由作为基础，实际菜单权限通过组件级别控制
		],
	},
];
