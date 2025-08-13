import { useBackendDashboardRoutes } from "@/routes/sections/dashboard/backend";
import useUserStore, { useMenuRefresh } from "@/store/userStore";
import { Button } from "@/ui/button";
import { useState } from "react";

export default function RouteDebugPage() {
	const userMenuTree = useUserStore((state) => state.userMenuTree);
	const dynamicRoutes = useBackendDashboardRoutes();
	const { refreshMenu } = useMenuRefresh();
	const [isRefreshing, setIsRefreshing] = useState(false);

	const handleRefreshMenu = async () => {
		setIsRefreshing(true);
		try {
			await refreshMenu();
			console.log('菜单刷新完成');
		} catch (error) {
			console.error('菜单刷新失败:', error);
		} finally {
			setIsRefreshing(false);
		}
	};

	return (
		<div className="p-8">
			<div className="flex justify-between items-center mb-6">
				<h1 className="text-2xl font-bold">路由调试信息</h1>
				<Button 
					onClick={handleRefreshMenu} 
					disabled={isRefreshing}
					variant="outline"
				>
					{isRefreshing ? '刷新中...' : '刷新菜单树'}
				</Button>
			</div>
			
			<div className="mb-8">
				<h2 className="text-lg font-semibold mb-2">
					用户菜单树 (共 {userMenuTree.length} 个根菜单)
				</h2>
				<pre className="bg-gray-100 p-4 rounded text-xs overflow-auto max-h-96">
					{JSON.stringify(userMenuTree, null, 2)}
				</pre>
			</div>

			<div className="mb-8">
				<h2 className="text-lg font-semibold mb-2">生成的动态路由</h2>
				<pre className="bg-gray-100 p-4 rounded text-xs overflow-auto max-h-96">
					{JSON.stringify(dynamicRoutes, null, 2)}
				</pre>
			</div>

			<div>
				<h2 className="text-lg font-semibold mb-2">组件路径测试</h2>
				{userMenuTree.map(group => 
					group.children?.map(item => 
						item.children?.map(subItem => {
							const componentPath = subItem.meta?.component || subItem.component;
							if (componentPath && subItem.type === 2) {
								return (
									<div key={subItem.id} className="mb-2 p-2 border rounded">
										<div><strong>菜单:</strong> {subItem.name}</div>
										<div><strong>路径:</strong> {subItem.path}</div>
										<div><strong>组件:</strong> {componentPath}</div>
									</div>
								);
							}
							return null;
						})
					)
				)}
			</div>
		</div>
	);
}
