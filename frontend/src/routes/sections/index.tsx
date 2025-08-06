import { Navigate, type RouteObject } from "react-router";
import { authRoutes } from "./auth";
import { mainRoutes } from "./main";

export const routesSection: RouteObject[] = [
	// Auth
	...authRoutes,
	// Main
	...mainRoutes,
	// No Match
	{ path: "*", element: <Navigate to="/404" replace /> },
];
