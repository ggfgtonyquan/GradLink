<template>
  <section class="page auth-view">
    <h1 class="screen-title">登录 GradLink</h1>
    <form class="surface auth-form" @submit.prevent>
      <div class="field">
        <label for="email">邮箱</label>
        <input id="email" type="email" autocomplete="email" value="student@example.com" />
      </div>
      <div class="field">
        <label for="password">密码</label>
        <input id="password" type="password" autocomplete="current-password" value="password" />
      </div>
      <button class="primary-button" type="submit">登录</button>
      <RouterLink to="/register" class="hint">还没有账号？去注册</RouterLink>
    </form>
  </section>
</template>

<style scoped>
.auth-form {
  display: grid;
  padding: 16px;
}
</style>
