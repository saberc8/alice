import QrCode from "qrcode";
import { useEffect, useState } from "react";
import { ReturnButton } from "./components/ReturnButton";
import { LoginStateEnum, useLoginStateContext } from "./providers/login-provider";

function QrCodeFrom() {
	const { loginState, backToLogin } = useLoginStateContext();

	if (loginState !== LoginStateEnum.QR_CODE) return null;
	return (
		<>
			<div className="flex flex-col items-center gap-2 text-center">
				<h1 className="text-2xl font-bold">二维码登录</h1>
				<p className="text-balance text-sm text-muted-foreground">扫码后点击'确认'，即可完成登录</p>
			</div>

			<div className="flex w-full flex-col items-center justify-center p-4">
				<QRCodeSVG value="https://github.com/d3george/slash-admin" size={200} />
			</div>
			<ReturnButton onClick={backToLogin} />
		</>
	);
}

export default QrCodeFrom;
