import { AvatarGroup } from "@/components/avatar-group";
import { Icon } from "@/components/icon";
import { Avatar, AvatarImage } from "@/ui/avatar";
import { Badge } from "@/ui/badge";
import { Button } from "@/ui/button";
import { Card, CardContent } from "@/ui/card";
import { Text } from "@/ui/typography";
import dayjs from "dayjs";

export default function ProjectsTab() {
	const items = [
		{
			icon: <Icon icon="logos:react" size={40} />,
			name: "Admin Template",
			client: "John Doe",
			desc: "Time is our most valuable asset, that is why we want to help you save it by creating…",
			startDate: dayjs().subtract(6, "month"),
			deadline: dayjs().add(6, "month"),
			messages: 236,
			allHours: "98/135",
			allTasks: 135,
			closedTasks: 98,
		}
	];
	return (
		<div className="grid grid-cols-1 gap-4 lg:grid-cols-2">
			{items.map((item) => (
				<Card key={item.name} className="flex w-full flex-col">
					<CardContent>
						<header className="flex w-full items-center">
							{item.icon}

							<div className="flex flex-col">
								<Text variant="body1" className="ml-4">
									{item.name}
								</Text>
								<Text variant="caption" className="ml-4">
									Client: {item.client}
								</Text>
							</div>

							<div className="ml-auto flex opacity-70">
								<Button variant="ghost" size="icon">
									<Icon icon="fontisto:more-v-a" size={18} />
								</Button>
							</div>
						</header>

						<main className="mt-4 w-full">
							<div className="my-2 flex justify-between">
								<Text variant="body1">
									Start Date:
									<Text variant="caption" className="ml-2">
										{item.startDate.format("DD/MM/YYYY")}
									</Text>
								</Text>

								<Text variant="body1">
									Deadline:
									<Text variant="caption" className="ml-2">
										{item.deadline.format("DD/MM/YYYY")}
									</Text>
								</Text>
							</div>
							<span className="opacity-70">{item.desc}</span>
						</main>

						<footer className="flex w-full  flex-col items-center">
							<div className="mb-4 flex w-full justify-between">
								<span>
									<Text variant="body1">All Hours:</Text>
									<Text variant="caption" className="ml-2">
										{item.allHours}
									</Text>
								</span>

								<Badge variant="warning">{item.deadline.diff(dayjs(), "day")} days left</Badge>
							</div>
							<div className="flex w-full ">
								<div className="ml-auto flex items-center opacity-50">
									<Icon icon="solar:chat-round-line-linear" size={24} />
									<Text variant="subTitle2" className="ml-2">
										{item.messages}
									</Text>
								</div>
							</div>
						</footer>
					</CardContent>
				</Card>
			))}
		</div>
	);
}
