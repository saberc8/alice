<script setup lang="ts">
import { ref } from 'vue'
import { Message } from '@arco-design/web-vue'
import { useAuthStore } from '@/stores/auth'
import { useRoute, useRouter } from 'vue-router'

const auth = useAuthStore()
const router = useRouter()
const route = useRoute()

const loading = ref(false)
const form = ref({ username: 'admin', password: '123456' })

async function onSubmit() {
	loading.value = true
	try {
		await auth.doLogin(form.value.username, form.value.password)
		Message.success('登录成功')
		const redirect = (route.query.redirect as string) || '/'
		router.replace(redirect)
	} catch (e: any) {
		Message.error(e?.message || '登录失败')
	} finally {
		loading.value = false
	}
}
</script>

<template>
	<div class="login-wrap">
		<div class="login-bg" aria-hidden="true">
			<span class="blob b1"></span>
			<span class="blob b2"></span>
			<span class="blob b3"></span>
			<span class="blob b4"></span>
		</div>
		<a-card class="login-card" :bordered="false">
			<h2 class="title">Alice Admin</h2>
			<a-form :model="form" layout="vertical" @submit.prevent="onSubmit">
				<a-form-item field="username" label="用户名">
					<a-input v-model="form.username" placeholder="请输入用户名" />
				</a-form-item>
				<a-form-item field="password" label="密码">
					<a-input-password v-model="form.password" placeholder="请输入密码" />
				</a-form-item>
				<div class="btn-wrap">
					<a-button
						class="login-btn"
						type="primary"
						:loading="loading"
						@click="onSubmit"
						>登录</a-button
					>
				</div>
			</a-form>
		</a-card>
	</div>
</template>

<style scoped lang="less">
.login-wrap {
	min-height: 100vh;
	display: flex;
	align-items: center;
	justify-content: center;
	background: #ffffff;
	position: relative;
	overflow: hidden;
}
/* 纯白画布 + 聚集式弥散色块 */
.login-bg {
	position: absolute;
	inset: 0;
	pointer-events: none;
}
.blob {
	position: absolute;
	border-radius: 50%;
	filter: blur(110px);
	opacity: 0.75;
	mix-blend-mode: normal;
	animation: pulse 32s ease-in-out infinite;
}
.b1 {
	width: 520px;
	height: 520px;
	background: rgba(255, 167, 136, 0.5);
	top: 38%;
	left: 45%;
	margin: -260px 0 0 -260px;
	animation-delay: 0s;
}
.b2 {
	width: 440px;
	height: 440px;
	background: rgba(255, 222, 140, 0.5);
	top: 40%;
	left: 58%;
	margin: -220px 0 0 -220px;
	animation-delay: -6s;
}
.b3 {
	width: 460px;
	height: 460px;
	background: rgba(158, 205, 255, 0.5);
	top: 57%;
	left: 50%;
	margin: -230px 0 0 -230px;
	animation-delay: -12s;
}
.b4 {
	width: 400px;
	height: 400px;
	background: rgba(196, 170, 255, 0.5);
	top: 50%;
	left: 38%;
	margin: -200px 0 0 -200px;
	animation-delay: -18s;
}
@keyframes pulse {
	0%,
	100% {
		transform: translate3d(0, 0, 0) scale(1);
	}
	25% {
		transform: translate3d(-20px, -14px, 0) scale(1.05);
	}
	50% {
		transform: translate3d(18px, 16px, 0) scale(1.03);
	}
	75% {
		transform: translate3d(-14px, 18px, 0) scale(1.06);
	}
}
@media (prefers-reduced-motion: reduce) {
	.blob {
		animation: none;
	}
}
.login-card {
	width: 380px;
	background: rgba(255, 255, 255, 0.85);
	backdrop-filter: blur(12px) saturate(150%);
	border: 1px solid rgba(0, 0, 0, 0.06);
	box-shadow: 0 4px 24px -6px rgba(0, 0, 0, 0.12),
		0 2px 6px -2px rgba(0, 0, 0, 0.06);
	color: #1a1a1a;
	position: relative;
	z-index: 1;
	border-radius: 18px;
}
.login-card :deep(.arco-card-body) {
	padding: 30px 34px 34px;
}
.login-card :deep(.arco-form-item-label-col) {
	color: #333;
	font-weight: 500;
}
.login-card :deep(.arco-input),
.login-card :deep(.arco-input-wrapper),
.login-card :deep(.arco-input-password) {
	background: rgba(255, 255, 255, 0.9);
	border-color: rgba(0, 0, 0, 0.12);
	color: #1a1a1a;
}
.login-card :deep(.arco-input:hover),
.login-card :deep(.arco-input-wrapper:hover) {
	border-color: #6366f1;
}
.login-card :deep(.arco-input:focus-within),
.login-card :deep(.arco-input-wrapper:focus-within) {
	border-color: #6366f1;
	box-shadow: 0 0 0 2px rgba(99, 102, 241, 0.25);
}
.login-card :deep(.arco-input::placeholder) {
	color: #9aa0b1;
}
.login-card :deep(.arco-btn-primary) {
	background: linear-gradient(140deg, #1d9df8 0%, #37c9e3 45%, #5adeb2 100%);
	border: none;
	box-shadow: 0 4px 14px -4px rgba(6, 182, 212, 0.55),
		0 2px 6px -2px rgba(16, 185, 129, 0.35);
}
.login-card :deep(.arco-btn-primary:hover) {
	filter: brightness(1.07);
}
.login-card :deep(.arco-btn-primary:active) {
	filter: brightness(0.93);
}
.btn-wrap {
	margin-top: 4px;
}
.login-btn {
	width: 100%;
	height: 40px;
	font-size: 15px;
}
.title {
	text-align: center;
	margin-bottom: 18px;
	font-weight: 600;
	letter-spacing: 0.5px;
	background: linear-gradient(
		120deg,
		#4f46e5,
		#6366f1 35%,
		#8b5cf6 70%,
		#ec4899
	);
	-webkit-background-clip: text;
	background-clip: text;
	color: transparent;
}
</style>
