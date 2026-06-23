<template>
  <section class="page auth-view">
    <h1 class="screen-title">注册账号</h1>
    <form class="surface auth-form" @submit.prevent>
      <div class="field">
        <label for="nickname">昵称</label>
        <input id="nickname" value="准备秋招的同学" />
      </div>
      <div class="field">
        <label for="email">邮箱</label>
        <input id="email" type="email" autocomplete="email" />
      </div>
      <div class="field">
        <label for="password">密码</label>
        <input id="password" type="password" autocomplete="new-password" />
      </div>
      <button class="primary-button" type="submit">注册</button>
      <p class="hint">邮箱只用于注册和登录，MVP 阶段不要求邮箱验证。</p>
    </form>
  </section>
</template>

<style scoped>
.auth-form {
  display: grid;
  padding: 16px;
}
</style>
