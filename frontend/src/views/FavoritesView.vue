<template>
  <section class="page">
    <div>
      <h1 class="screen-title">我的收藏</h1>
      <p class="hint">你收藏过的讨论会在这里集中查看。</p>
    </div>

    <div v-if="!auth.isLoggedIn" class="surface gate-box">
      <p>登录后可以查看收藏内容。</p>
      <RouterLink to="/login" class="primary-button">去登录</RouterLink>
    </div>

    <div v-else class="feed-list">
      <PostCard v-for="post in posts" :key="post.id" :post="post" />
      <p v-if="!posts.length" class="empty-state">还没有收藏内容。</p>
    </div>
  </section>
</template>

<script setup lang="ts">
import { onMounted, ref } from 'vue'
import { api } from '../api/http'
import type { Post } from '../api/types'
import PostCard from '../components/PostCard.vue'
import { useAuthStore } from '../stores/auth'

const auth = useAuthStore()
const posts = ref<Post[]>([])

onMounted(async () => {
  await auth.loadMe()
  if (auth.isLoggedIn) {
    posts.value = (await api.myFavorites()).posts
  }
})
</script>

<style scoped>
.feed-list {
  display: grid;
  gap: 10px;
}

.gate-box {
  display: grid;
  gap: 10px;
  padding: 16px;
}

.gate-box p,
.empty-state {
  margin: 0;
  color: var(--muted);
  font-size: 14px;
}

.empty-state {
  border: 1px dashed var(--line);
  border-radius: 8px;
  padding: 16px;
  text-align: center;
}
</style>
