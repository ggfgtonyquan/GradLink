<template>
  <section class="page auth-view">
    <h1 class="screen-title">注册账号</h1>
    <form class="surface auth-form" @submit.prevent="submit">
      <div class="field">
        <label for="nickname">昵称</label>
        <input id="nickname" v-model="nickname" />
      </div>
      <div class="field">
        <label for="email">邮箱</label>
        <input id="email" v-model="email" type="email" autocomplete="email" />
      </div>
      <div class="field">
        <label for="password">密码</label>
        <input id="password" v-model="password" type="password" autocomplete="new-password" />
      </div>
      <button class="primary-button" type="submit">注册</button>
      <p class="hint">邮箱只用于注册和登录，MVP 阶段不要求邮箱验证。</p>
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
const nickname = ref('准备秋招的同学')
const email = ref('')
const password = ref('')
const errorMessage = ref('')

async function submit() {
  errorMessage.value = ''
  try {
    await auth.register(email.value, password.value, nickname.value)
    await router.push('/')
  } catch {
    errorMessage.value = auth.error || '注册失败'
  }
}
</script>

<style scoped>
.auth-form {
  display: grid;
  padding: 16px;
}
</style>
