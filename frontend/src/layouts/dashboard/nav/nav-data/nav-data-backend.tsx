import { Icon } from "@/components/icon";
import type { NavItemDataProps, NavProps } from "@/components/nav";
import type { MenuTree } from "@/types/entity";
import { Badge } from "@/ui/badge";
import { convertFlatToTree } from "@/utils/tree";
import useUserStore from "@/store/userStore";

const convertChildren = (children?: MenuTree[]): NavItemDataProps[] => {
	if (!children?.length) return [];

	return children.map((child) => ({
		title: child.name,
		path: child.path || "",
		icon: child.icon ? typeof child.icon === "string" ? <Icon icon={child.icon} size="24" /> : child.icon : null,
		caption: child.caption,
		info: child.info ? <Badge variant="default">{child.info}</Badge> : null,
		disabled: child.disabled,
		externalLink: child.externalLink,
		auth: child.auth,
		hidden: child.hidden,
		children: convertChildren(child.children),
	}));
};

const convert = (menuTree: MenuTree[]): NavProps["data"] => {
	return menuTree.map((item) => ({
		name: item.name,
		items: convertChildren(item.children),
	}));
};

// Hook to get dynamic nav data from user store
export const useBackendNavData = (): NavProps["data"] => {
	const userMenuTree = useUserStore((state) => state.userMenuTree);
	return convert(userMenuTree);
};

// Export empty array as fallback for direct import (should use hook instead)
export const backendNavData: NavProps["data"] = [];
