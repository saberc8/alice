import { Navigate, type RouteObject } from "react-router";
import { authRoutes } from "./auth";
import { dashboardRoutes } from "./dashboard/index";
import { mainRoutes } from "./main";

export const routesSection: RouteObject[] = [
	// Root redirect
	{ path: "/", element: <Navigate to="/auth/login" replace /> },
	// Login redirect for backward compatibility
	{ path: "/login", element: <Navigate to="/auth/login" replace /> },
	// Auth
	...authRoutes,
	// Dashboard
	...dashboardRoutes,
	// Main
	...mainRoutes,
	// No Match
	{ path: "*", element: <Navigate to="/404" replace /> },
];
