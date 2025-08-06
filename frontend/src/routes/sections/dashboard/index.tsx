import { GLOBAL_CONFIG } from "@/global-config";
import DashboardLayout from "@/layouts/dashboard";
import LoginAuthGuard from "@/routes/components/login-auth-guard";
import { Navigate, type RouteObject } from "react-router";

export const dashboardRoutes: RouteObject[] = [
	{
		path: "/",
		element: (
			<LoginAuthGuard>
				<DashboardLayout />
			</LoginAuthGuard>
		),
		children: []
	},
];
