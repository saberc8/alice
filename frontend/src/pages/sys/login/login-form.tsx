import type { SignInReq } from "@/api/services/userService";
import { Icon } from "@/components/icon";
import { GLOBAL_CONFIG } from "@/global-config";
import { useSignIn } from "@/store/userStore";
import { Button } from "@/ui/button";
import { Checkbox } from "@/ui/checkbox";
import { Form, FormControl, FormField, FormItem, FormLabel, FormMessage } from "@/ui/form";
import { Input } from "@/ui/input";
import { cn } from "@/utils";
import { Loader2 } from "lucide-react";
import { useState } from "react";
import { useForm } from "react-hook-form";
import { useNavigate } from "react-router";
import { toast } from "sonner";
import { LoginStateEnum, useLoginStateContext } from "./providers/login-provider";

export function LoginForm({ className, ...props }: React.ComponentPropsWithoutRef<"form">) {
	const [loading, setLoading] = useState(false);
	const [remember, setRemember] = useState(true);
	const navigatge = useNavigate();

	const { loginState, setLoginState } = useLoginStateContext();
	const signIn = useSignIn();

	const form = useForm<SignInReq>({
		defaultValues: {
			username: '',
			password: '',
		},
	});

	if (loginState !== LoginStateEnum.LOGIN) return null;

	const handleFinish = async (values: SignInReq) => {
		setLoading(true);
		try {
			await signIn(values);
			navigatge(GLOBAL_CONFIG.defaultRoute, { replace: true });
			toast.success("登录成功", {
				closeButton: true,
			});
		} finally {
			setLoading(false);
		}
	};

	return (
		<div className={cn("flex flex-col gap-6", className)}>
			<Form {...form} {...props}>
				<form onSubmit={form.handleSubmit(handleFinish)} className="space-y-4">
					<div className="flex flex-col items-center gap-2 text-center">
						<h1 className="text-2xl font-bold">登录到您的账户</h1>
						<p className="text-balance text-sm text-muted-foreground">输入您的账号登录到您的账户</p>
					</div>

					<FormField
						control={form.control}
						name="username"
						rules={{ required: "请输入账号" }}
						render={({ field }) => (
							<FormItem>
								<FormLabel>账号</FormLabel>
								<FormControl>
									<Input placeholder="请输入账号" {...field} />
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
								<FormLabel>密码</FormLabel>
								<FormControl>
									<Input type="password" placeholder="请输入密码" {...field} suppressHydrationWarning />
								</FormControl>
								<FormMessage />
							</FormItem>
						)}
					/>

					{/* 记住我/忘记密码 */}
					<div className="flex flex-row justify-between">
						<div className="flex items-center space-x-2">
							<Checkbox
								id="remember"
								checked={remember}
								onCheckedChange={(checked) => setRemember(checked === "indeterminate" ? false : checked)}
							/>
							<label
								htmlFor="remember"
								className="text-sm font-medium leading-none peer-disabled:cursor-not-allowed peer-disabled:opacity-70"
							>
								记住我
							</label>
						</div>
						<Button variant="link" onClick={() => setLoginState(LoginStateEnum.RESET_PASSWORD)} size="sm">
							忘记密码?
						</Button>
					</div>

					{/* 登录按钮 */}
					<Button type="submit" className="w-full">
						{loading && <Loader2 className="animate-spin mr-2" />}
						登录
					</Button>

					{/* 手机登录/二维码登录 */}
					<div className="grid gap-4 sm:grid-cols-2">
						<Button variant="outline" className="w-full" onClick={() => setLoginState(LoginStateEnum.MOBILE)}>
							<Icon icon="uil:mobile-android" size={20} />
							手机登录
						</Button>
						<Button variant="outline" className="w-full" onClick={() => setLoginState(LoginStateEnum.QR_CODE)}>
							<Icon icon="uil:qrcode-scan" size={20} />
							二维码登录
						</Button>
					</div>

					{/* 其他登录方式 */}
					<div className="relative text-center text-sm after:absolute after:inset-0 after:top-1/2 after:z-0 after:flex after:items-center after:border-t after:border-border">
						<span className="relative z-10 bg-background px-2 text-muted-foreground">其他登录方式</span>
					</div>
					<div className="flex cursor-pointer justify-around text-2xl">
						<Button variant="ghost" size="icon">
							<Icon icon="mdi:github" size={24} />
						</Button>
						<Button variant="ghost" size="icon">
							<Icon icon="mdi:wechat" size={24} />
						</Button>
						<Button variant="ghost" size="icon">
							<Icon icon="ant-design:google-circle-filled" size={24} />
						</Button>
					</div>

					{/* 注册 */}
					<div className="text-center text-sm">
						没有账号？
						<Button variant="link" className="px-1" onClick={() => setLoginState(LoginStateEnum.REGISTER)}>
							注册
						</Button>
					</div>
				</form>
			</Form>
		</div>
	);
}

export default LoginForm;
