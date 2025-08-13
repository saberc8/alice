import { Icon } from "@/components/icon";
import { useSettings } from "@/store/settingStore";
import { useMenuRefresh } from "@/store/userStore";
import { Button } from "@/ui/button";
import { cn } from "@/utils";
import type { ReactNode } from "react";
import { useState } from "react";
import { toast } from "sonner";
import AccountDropdown from "../components/account-dropdown";
import BreadCrumb from "../components/bread-crumb";
import NoticeButton from "../components/notice";
import SearchBar from "../components/search-bar";
import SettingButton from "../components/setting-button";

interface HeaderProps {
	leftSlot?: ReactNode;
}

export default function Header({ leftSlot }: HeaderProps) {
	const { breadCrumb } = useSettings();
	const { refreshMenu } = useMenuRefresh();
	const [isRefreshing, setIsRefreshing] = useState(false);

	const handleRefreshMenu = async () => {
		setIsRefreshing(true);
		try {
			await refreshMenu();
			toast.success('菜单已刷新', { position: "top-center" });
			console.log('菜单刷新完成');
		} catch (error) {
			toast.error('菜单刷新失败', { position: "top-center" });
			console.error('菜单刷新失败:', error);
		} finally {
			setIsRefreshing(false);
		}
	};
	return (
		<header
			data-slot="slash-layout-header"
			className={cn(
				"sticky top-0 left-0 right-0 z-app-bar",
				"flex items-center justify-between px-2 grow-0 shrink-0",
				"bg-background/60 backdrop-blur-xl",
				"h-[var(--layout-header-height)] ",
			)}
		>
			<div className="flex items-center">
				{leftSlot}

				<div className="hidden md:block ml-4">{breadCrumb && <BreadCrumb />}</div>
			</div>

			<div className="flex items-center gap-1">
				<SearchBar />
				<Button
					variant="ghost"
					size="icon"
					className="rounded-full"
					onClick={handleRefreshMenu}
					disabled={isRefreshing}
					title="刷新菜单"
				>
					<Icon 
						icon={isRefreshing ? "line-md:loading-twotone-loop" : "material-symbols:refresh"} 
						size={18} 
					/>
				</Button>
				<Button
					variant="ghost"
					size="icon"
					className="rounded-full"
					onClick={() => window.open("https://github.com/d3george/slash-admin")}
				>
					<Icon icon="mdi:github" size={24} />
				</Button>
				<Button
					variant="ghost"
					size="icon"
					className="rounded-full"
					onClick={() => window.open("https://discord.gg/fXemAXVNDa")}
				>
					<Icon icon="carbon:logo-discord" size={24} />
				</Button>
				<NoticeButton />
				<SettingButton />
				<AccountDropdown />
			</div>
		</header>
	);
}
