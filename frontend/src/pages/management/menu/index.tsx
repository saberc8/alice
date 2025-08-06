import { Suspense } from "react";
import { Outlet } from "react-router";
import Loading from "@/components/loading";

// 菜单管理主页面
function MenuManagement() {
	return (
		<div className="p-6">
			<div className="mb-6">
				<h1 className="text-2xl font-bold text-foreground">菜单管理</h1>
				<p className="text-muted-foreground">管理系统菜单，配置菜单层级和访问权限</p>
			</div>
			
			<Suspense fallback={<Loading />}>
				<Outlet />
			</Suspense>
		</div>
	);
}

// 导出懒加载组件
export const Component = MenuManagement;
export default MenuManagement;
