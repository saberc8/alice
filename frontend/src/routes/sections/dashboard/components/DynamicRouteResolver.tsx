import { useMemo } from "react";
import { Navigate, useLocation } from "react-router";
import { useBackendDashboardRoutes } from "../backend";
import { GLOBAL_CONFIG } from "@/global-config";

/**
 * 动态路由匹配器
 * 这个组件会根据当前路径和用户菜单数据来动态渲染对应的组件
 */
export const DynamicRouteResolver = () => {
	const location = useLocation();
	const dynamicRoutes = useBackendDashboardRoutes();
	
	// 添加调试信息
	console.log("DynamicRouteResolver 调试信息:");
	console.log("  当前路径:", location.pathname);
	console.log("  动态路由数量:", dynamicRoutes.length);
	console.log("  动态路由列表:", dynamicRoutes.map(r => ({ path: r.path, hasElement: !!r.element })));
	
	const currentRoute = useMemo(() => {
		const findMatchingRoute = (routes: any[], fullPath: string, currentPath: string = ""): any => {
			console.log(`查找路由匹配: fullPath="${fullPath}", currentPath="${currentPath}"`);

			// 统一处理：实际 URL 总是以 '/' 开头，这里保证构建出来的候选路径也以 '/' 开头
			const normalize = (p: string) => {
				if (!p) return '/';
				return p.startsWith('/') ? p : '/' + p;
			};

			for (const route of routes) {
				console.log(`  检查路由: path="${route.path}", index=${route.index}`);

				// 处理索引路由 - 当完整路径匹配当前路径时（都做规范化）
				if (route.index) {
					const normCurrent = normalize(currentPath);
					if (normalize(fullPath) === normCurrent) {
						console.log("  ✓ 找到索引路由匹配");
						return route;
					}
				}

				// 处理普通 path 路由
				if (route.path) {
					// 构建当前层级的完整路径（父路径 currentPath 可能已是规范化或空）
					const parentFull = currentPath ? normalize(currentPath) : '';
					const raw = currentPath ? `${parentFull.replace(/\/$/, '')}/${route.path}` : route.path;
					const routeFullPath = normalize(raw);
					console.log(`  构建的完整路径: "${routeFullPath}" (原始: "${raw}")`);

					// 精确匹配（考虑大小写一致性，直接比较规范化后的路径）
					if (routeFullPath === normalize(fullPath)) {
						console.log(`  ✓ 找到精确路径匹配: ${routeFullPath}`);
						return route;
					}

					// 前缀匹配（用于继续向子节点深入）
					if (normalize(fullPath).startsWith(routeFullPath + '/') && route.children && route.children.length > 0) {
						console.log(`  → 路径前缀匹配，递归查找子路由: ${routeFullPath}`);
						const childMatch = findMatchingRoute(route.children, fullPath, routeFullPath);
						if (childMatch) return childMatch;
					}
				}
			}
			console.log("  ✗ 未找到匹配的路由");
			return null;
		};
		
		const matchedRoute = findMatchingRoute(dynamicRoutes, location.pathname);
		console.log("最终匹配结果:", matchedRoute ? `找到路由` : "无匹配", matchedRoute);
		return matchedRoute;
	}, [location.pathname, dynamicRoutes]);
	
	// 如果找到匹配的路由，渲染对应的组件
	if (currentRoute) {
		return currentRoute.element || <Navigate to={GLOBAL_CONFIG.defaultRoute} replace />;
	}
	
	// 如果没有找到匹配的路由，显示加载状态或 404
	if (dynamicRoutes.length === 0) {
		return (
			<div style={{ padding: '20px', textAlign: 'center' }}>
				<p>正在加载菜单数据...</p>
				<p>如果长时间显示此消息，请刷新页面或联系管理员</p>
				<details style={{ marginTop: '10px', textAlign: 'left' }}>
					<summary>调试信息</summary>
					<pre style={{ fontSize: '12px', backgroundColor: '#f5f5f5', padding: '10px', borderRadius: '4px' }}>
						当前路径: {location.pathname}
						动态路由数量: {dynamicRoutes.length}
					</pre>
				</details>
			</div>
		);
	}
	
	// 有路由数据但没有匹配的路由，可能是 404
	return (
		<div style={{ padding: '20px', textAlign: 'center' }}>
			<h3>页面不存在</h3>
			<p>请检查URL是否正确，或者您是否有权限访问此页面</p>
			<p>
				<a href={GLOBAL_CONFIG.defaultRoute}>返回首页</a>
			</p>
			<details style={{ marginTop: '10px', textAlign: 'left' }}>
				<summary>调试信息</summary>
				<pre style={{ fontSize: '12px', backgroundColor: '#f5f5f5', padding: '10px', borderRadius: '4px' }}>
					当前路径: {location.pathname}
					可用路由: {JSON.stringify(dynamicRoutes.map(r => ({ path: r.path })), null, 2)}
				</pre>
			</details>
		</div>
	);
};
