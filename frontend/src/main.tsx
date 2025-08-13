import "./global.css";
import "./theme/theme.css";
import ReactDOM from "react-dom/client";
import { Outlet, RouterProvider, createBrowserRouter } from "react-router";
import App from "./App";
import menuService from "./api/services/menuService";
import { registerLocalIcons } from "./components/icon";
import { GLOBAL_CONFIG } from "./global-config";
import ErrorBoundary from "./routes/components/error-boundary";
import { routesSection } from "./routes/sections";
import { urlJoin } from "./utils";
import useUserStore from "./store/userStore";

await registerLocalIcons();

// 如果是后端路由模式，预加载菜单数据到 store
if (GLOBAL_CONFIG.routerMode === "backend") {
	try {
		console.log("后端路由模式：预加载菜单数据...");
		const menuResponse = await menuService.getMenuList();
		
		if (menuResponse.code === 200 && menuResponse.data) {
			// 将菜单数据保存到 store
			const { setUserMenuTree } = useUserStore.getState().actions;
			setUserMenuTree(menuResponse.data);
			console.log("菜单数据预加载完成，路由数量:", menuResponse.data.length);
		} else {
			console.warn("菜单数据加载失败:", menuResponse);
		}
	} catch (error) {
		console.error("预加载菜单数据失败:", error);
		// 即使失败也继续启动应用，菜单数据会在用户登录后加载
	}
}

const router = createBrowserRouter(
	[
		{
			Component: () => (
				<App>
					<Outlet />
				</App>
			),
			errorElement: <ErrorBoundary />,
			children: routesSection,
		},
	],
	{
		basename: GLOBAL_CONFIG.publicPath,
	},
);

const root = ReactDOM.createRoot(document.getElementById("root") as HTMLElement);
root.render(<RouterProvider router={router} />);
