import { AuthGuard } from "@/components/auth/auth-guard";
import { LineLoading } from "@/components/loading";
import { GLOBAL_CONFIG } from "@/global-config";
import Page403 from "@/pages/sys/error/Page403";
import { useSettings } from "@/store/settingStore";
import { cn } from "@/utils";
import { flattenTrees } from "@/utils/tree";
import { clone, concat } from "ramda";
import { Suspense } from "react";
import { Outlet, ScrollRestoration, useLocation } from "react-router";
import { useBackendNavData } from "./nav/nav-data/nav-data-backend";

const Main = () => {
	const backendNavData = useBackendNavData();
	const navData = backendNavData;
	const allItems = navData.reduce((acc: any[], group) => {
		const flattenedItems = flattenTrees(group.items);
		return concat(acc, flattenedItems);
	}, []);
	const { themeStretch } = useSettings();

	const { pathname } = useLocation();

	/**
	 * find auth by path
	 * @param path
	 * @returns
	 */
	const findAuthByPath = (path: string): string[] => {
		const foundItem = allItems.find((item) => item.path === path);
		return foundItem?.auth || [];
	};

	const currentNavAuth = findAuthByPath(pathname);

	return (
		<AuthGuard checkAny={currentNavAuth} fallback={<Page403 />}>
			<main
				data-slot="slash-layout-main"
				className={cn(
					"flex-auto w-full flex flex-col",
					"transition-[max-width] duration-300 ease-in-out",
					"px-4 sm:px-6 py-4 sm:py-6 md:px-8 mx-auto",
					{
						"max-w-full": themeStretch,
						"xl:max-w-screen-xl": !themeStretch,
					},
				)}
				style={{
					willChange: "max-width",
				}}
			>
				<Suspense fallback={<LineLoading />}>
					<Outlet />
					<ScrollRestoration />
				</Suspense>
			</main>
		</AuthGuard>
	);
};

export default Main;
