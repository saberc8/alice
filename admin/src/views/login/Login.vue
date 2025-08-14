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
		<a-card class="login-card" :bordered="false">
			<h2 class="title">Alice Admin 登录</h2>
			<a-form :model="form" layout="vertical" @submit.prevent="onSubmit">
				<a-form-item field="username" label="用户名">
					<a-input v-model="form.username" placeholder="请输入用户名" />
				</a-form-item>
				<a-form-item field="password" label="密码">
					<a-input-password v-model="form.password" placeholder="请输入密码" />
				</a-form-item>
				<a-space fill>
					<a-button type="primary" long :loading="loading" @click="onSubmit">登录</a-button>
				</a-space>
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
	background: #f5f7fa;
}
.login-card {
	width: 360px;
	box-shadow: 0 6px 24px rgba(0,0,0,0.08);
}
.title {
	text-align: center;
	margin-bottom: 12px;
}
</style>
