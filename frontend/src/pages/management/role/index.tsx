import { Suspense } from "react";
import { Outlet } from "react-router";
import Loading from "@/components/loading";

// 角色管理主页面
function RoleManagement() {
	return (
		<div className="p-6">
			<div className="mb-6">
				<h1 className="text-2xl font-bold text-foreground">角色管理</h1>
				<p className="text-muted-foreground">管理系统角色，为角色分配权限和菜单</p>
			</div>
			
			<Suspense fallback={<Loading />}>
				<Outlet />
			</Suspense>
		</div>
	);
}

// 导出懒加载组件
export const Component = RoleManagement;
export default RoleManagement;
