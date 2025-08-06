import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/ui/card";

export default function PermissionPageTest() {
	return (
		<div className="grid grid-cols-1 gap-4 md:grid-cols-2">
			{Array.from({ length: 10 }).map((_, index) => (
				<Card key={index}>
					<CardHeader>
						<CardTitle>Card {index + 1}</CardTitle>
						<CardDescription></CardDescription>
					</CardHeader>
					<CardContent></CardContent>
				</Card>
			))}
		</div>
	);
}
