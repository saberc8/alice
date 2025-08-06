import { Suspense } from "react";
import { Outlet } from "react-router";
import Loading from "@/components/loading";

// 用户管理主页面
function UserManagement() {
	return (
		<div className="p-6">
			<div className="mb-6">
				<h1 className="text-2xl font-bold text-foreground">用户管理</h1>
				<p className="text-muted-foreground">管理系统用户，包括创建、编辑、删除用户，以及为用户分配角色</p>
			</div>
			
			<Suspense fallback={<Loading />}>
				<Outlet />
			</Suspense>
		</div>
	);
}

// 导出懒加载组件
export const Component = UserManagement;
export default UserManagement;
