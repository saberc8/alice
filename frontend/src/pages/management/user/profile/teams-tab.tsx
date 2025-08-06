import { Icon } from "@/components/icon";
import { Avatar, AvatarImage } from "@/ui/avatar";
import { Badge } from "@/ui/badge";
import { Button } from "@/ui/button";
import { Card, CardContent } from "@/ui/card";

export default function TeamsTab() {
	const items = [
		{
			icon: <Icon icon="logos:react" size={40} />,
			name: "React Developers",
			desc: "We don’t make assumptions about the rest of your technology stack, so you can develop new features in React.",
			tags: ["React", "AntD"],
		},
		{
			icon: <Icon icon="logos:vue" size={40} />,
			name: "Vue.js Dev Team",
			desc: "The development of Vue and its ecosystem is guided by an international team, some of whom have chosen to be featured below.",
			tags: ["Vue.js", "Developer"],
		},
	];
	return (
		<div className="grid grid-cols-1 gap-4 lg:grid-cols-2">
			{items.map((item) => (
				<Card key={item.name} className="flex w-full flex-col">
					<CardContent>
						<header className="flex w-full items-center">
							{item.icon}
							<span className="ml-4 text-xl opacity-70">{item.name}</span>

							<div className="ml-auto flex opacity-70">
								<Button variant="ghost" size="icon">
									<Icon icon="solar:star-line-duotone" size={18} />
								</Button>
								<Button variant="ghost" size="icon">
									<Icon icon="fontisto:more-v-a" size={18} />
								</Button>
							</div>
						</header>
						<main className="my-4 opacity-70">{item.desc}</main>
						<footer className="flex w-full items-center">
							<div className="ml-auto flex items-center gap-1">
								{item.tags.map((tag) => (
									<Badge key={tag} variant="info">
										{tag}
									</Badge>
								))}
							</div>
						</footer>
					</CardContent>
				</Card>
			))}
		</div>
	);
}
