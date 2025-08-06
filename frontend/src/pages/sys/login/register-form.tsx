import userService from "@/api/services/userService";
import { Button } from "@/ui/button";
import { Form, FormControl, FormField, FormItem, FormMessage } from "@/ui/form";
import { Input } from "@/ui/input";
import { useMutation } from "@tanstack/react-query";
import { useForm } from "react-hook-form";
import { ReturnButton } from "./components/ReturnButton";
import { LoginStateEnum, useLoginStateContext } from "./providers/login-provider";

function RegisterForm() {
	const { loginState, backToLogin } = useLoginStateContext();

	const signUpMutation = useMutation({
		mutationFn: userService.signup,
	});

	const form = useForm({
		defaultValues: {
			username: "",
			email: "",
			password: "",
			confirmPassword: "",
		},
	});

	const onFinish = async (values: any) => {
		console.log("Received values of form: ", values);
		await signUpMutation.mutateAsync(values);
		backToLogin();
	};

	if (loginState !== LoginStateEnum.REGISTER) return null;

	return (
		<Form {...form}>
			<form onSubmit={form.handleSubmit(onFinish)} className="space-y-4">
				<div className="flex flex-col items-center gap-2 text-center">
					<h1 className="text-2xl font-bold">用户注册</h1>
				</div>

				<FormField
					control={form.control}
					name="username"
					rules={{ required: "请输入用户名" }}
					render={({ field }) => (
						<FormItem>
							<FormControl>
								<Input placeholder={"请输入用户名"} {...field} />
							</FormControl>
							<FormMessage />
						</FormItem>
					)}
				/>

				<FormField
					control={form.control}
					name="email"
					rules={{ required: "请输入邮箱" }}
					render={({ field }) => (
						<FormItem>
							<FormControl>
								<Input placeholder={"请输入邮箱"} {...field} />
							</FormControl>
							<FormMessage />
						</FormItem>
					)}
				/>

				<FormField
					control={form.control}
					name="password"
					rules={{ required: "请输入密码" }}
					render={({ field }) => (
						<FormItem>
							<FormControl>
								<Input type="password" placeholder={"请输入密码"} {...field} />
							</FormControl>
							<FormMessage />
						</FormItem>
					)}
				/>

				<FormField
					control={form.control}
					name="confirmPassword"
					rules={{
						required: "请确认密码",
						validate: (value) => value === form.getValues("password") || "两次输入的密码不一致",
					}}
					render={({ field }) => (
						<FormItem>
							<FormControl>
								<Input type="password" placeholder={"请确认密码"} {...field} />
							</FormControl>
							<FormMessage />
						</FormItem>
					)}
				/>

				<Button type="submit" className="w-full">
					注册
				</Button>

				<div className="mb-2 text-xs text-gray">
					<span>同意协议</span>
				</div>

				<ReturnButton onClick={backToLogin} />
			</form>
		</Form>
	);
}

export default RegisterForm;
