import { Suspense } from "react";
import { Outlet } from "react-router";
import Loading from "@/components/loading";

// 权限管理主页面
function PermissionManagement() {
	return (
		<div className="p-6">
			<div className="mb-6">
				<h1 className="text-2xl font-bold text-foreground">权限管理</h1>
				<p className="text-muted-foreground">管理系统权限，定义资源和操作的访问控制</p>
			</div>
			
			<Suspense fallback={<Loading />}>
				<Outlet />
			</Suspense>
		</div>
	);
}

// 导出懒加载组件
export const Component = PermissionManagement;
export default PermissionManagement;
