<template>
  <nav class="bottom-nav" aria-label="底部导航">
    <RouterLink v-for="item in items" :key="item.to" :to="item.to" class="tab-item">
      <component :is="item.icon" aria-hidden="true" />
      <span>{{ item.label }}</span>
    </RouterLink>
  </nav>
</template>

<script setup lang="ts">
import { Bookmark, Home, PlusCircle, UserRound } from '@lucide/vue'

const items = [
  { to: '/', label: '首页', icon: Home },
  { to: '/compose', label: '发帖', icon: PlusCircle },
  { to: '/favorites', label: '收藏', icon: Bookmark },
  { to: '/profile', label: '我的', icon: UserRound },
]
</script>

<style scoped>
.bottom-nav {
  position: fixed;
  left: 50%;
  bottom: 0;
  width: min(100%, 430px);
  transform: translateX(-50%);
  display: grid;
  grid-template-columns: repeat(4, 1fr);
  min-height: 64px;
  padding: 8px 8px max(8px, env(safe-area-inset-bottom));
  border-top: 1px solid var(--line);
  background: rgb(255 255 255 / 96%);
  backdrop-filter: blur(18px);
}

.tab-item {
  display: grid;
  min-width: 0;
  min-height: 46px;
  align-content: center;
  justify-items: center;
  gap: 3px;
  border-radius: 8px;
  color: var(--muted);
  font-size: 12px;
  font-weight: 700;
}

.tab-item svg {
  width: 19px;
  height: 19px;
  stroke-width: 2;
}

.tab-item.router-link-active {
  color: var(--brand);
  background: var(--brand-soft);
}
</style>
