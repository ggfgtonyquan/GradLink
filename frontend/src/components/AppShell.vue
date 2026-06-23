<template>
  <div class="app-shell">
    <header class="top-bar">
      <RouterLink to="/" class="brand">
        <span class="eyebrow">GradLink</span>
        <strong>应届生交流</strong>
      </RouterLink>
      <div class="top-actions">
        <RouterLink to="/moderation" class="icon-button" aria-label="审核后台">
          <ShieldCheck aria-hidden="true" />
        </RouterLink>
        <RouterLink to="/profile" class="profile-button" aria-label="打开个人中心">我</RouterLink>
      </div>
    </header>

    <main class="content">
      <slot>
        <RouterView />
      </slot>
    </main>

    <BottomNav />
  </div>
</template>

<script setup lang="ts">
import { ShieldCheck } from '@lucide/vue'
import BottomNav from './BottomNav.vue'
</script>

<style scoped>
.app-shell {
  min-height: 100vh;
  max-width: 430px;
  margin: 0 auto;
  background: var(--app);
  display: flex;
  flex-direction: column;
  box-shadow: 0 0 32px rgb(21 43 33 / 10%);
}

.top-bar {
  position: sticky;
  top: 0;
  z-index: 2;
  display: flex;
  align-items: center;
  justify-content: space-between;
  min-height: 72px;
  padding: 14px 18px;
  background: rgb(248 250 247 / 94%);
  border-bottom: 1px solid var(--line);
  backdrop-filter: blur(16px);
}

.brand {
  display: grid;
  gap: 3px;
}

.eyebrow {
  color: var(--brand);
  font-size: 12px;
  font-weight: 800;
}

.brand strong {
  font-size: 21px;
  line-height: 1.2;
}

.top-actions {
  display: flex;
  align-items: center;
  gap: 8px;
}

.icon-button,
.profile-button {
  width: 40px;
  height: 40px;
  border: 1px solid #d2ded7;
  border-radius: 50%;
  background: #ffffff;
  color: #315849;
  display: grid;
  place-items: center;
  font-size: 15px;
  font-weight: 700;
}

.icon-button svg {
  width: 18px;
  height: 18px;
}

.content {
  flex: 1;
  padding: 16px 16px 88px;
}
</style>
