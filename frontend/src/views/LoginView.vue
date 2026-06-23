<template>
  <section class="page auth-view">
    <h1 class="screen-title">登录 GradLink</h1>
    <form class="surface auth-form" @submit.prevent="submit">
      <div class="field">
        <label for="email">邮箱</label>
        <input id="email" v-model="email" type="email" autocomplete="email" />
      </div>
      <div class="field">
        <label for="password">密码</label>
        <input id="password" v-model="password" type="password" autocomplete="current-password" />
      </div>
      <button class="primary-button" type="submit">登录</button>
      <RouterLink to="/register" class="hint">还没有账号？去注册</RouterLink>
      <p v-if="errorMessage" class="hint">{{ errorMessage }}</p>
    </form>
  </section>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { useAuthStore } from '../stores/auth'

const auth = useAuthStore()
const router = useRouter()
const email = ref('student@example.com')
const password = ref('password')
const errorMessage = ref('')

async function submit() {
  errorMessage.value = ''
  try {
    await auth.login(email.value, password.value)
    await router.push('/')
  } catch {
    errorMessage.value = auth.error || '登录失败'
  }
}
</script>

<style scoped>
.auth-form {
  display: grid;
  padding: 16px;
}
</style>
