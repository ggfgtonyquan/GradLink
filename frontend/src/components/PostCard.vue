<template>
  <RouterLink :to="`/posts/${post.id}`" class="post-card surface">
    <h2 class="post-title">{{ post.title }}</h2>
    <p class="post-excerpt">{{ post.excerpt }}</p>
    <div class="tag-row">
      <span class="tag">{{ tagNames }}</span>
    </div>
    <div class="meta-row">
      <span>{{ authorLine }}</span>
      <span class="metric-row">
        <span v-if="post.favorited" class="favorite-marker" aria-label="已收藏" title="已收藏">
          <Bookmark aria-hidden="true" fill="currentColor" />
        </span>
        <span>评论 {{ post.counts.comments }}</span>
        <span>赞 {{ post.counts.likes }}</span>
      </span>
    </div>
  </RouterLink>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { Bookmark } from '@lucide/vue'
import type { Post } from '../api/types'

const props = defineProps<{ post: Post }>()

const authorLine = computed(() => {
  const school = props.post.author.school ? props.post.author.school : `${props.post.author.graduationYear}届`
  return `${school} · ${props.post.author.educationLevel}`
})
const tagNames = computed(() => props.post.tags.map((tag) => tag.name).join(' / '))
</script>

<style scoped>
.post-card {
  display: block;
  padding: 14px;
}

.post-title {
  margin: 0 0 7px;
  font-size: 16px;
  line-height: 1.35;
  font-weight: 760;
}

.post-excerpt {
  margin: 0 0 10px;
  color: #42504b;
  font-size: 14px;
  line-height: 1.55;
}

.tag-row {
  display: flex;
  flex-wrap: wrap;
  gap: 6px;
  margin-bottom: 10px;
}

.tag {
  border-radius: 999px;
  background: var(--brand-soft);
  color: var(--brand);
  padding: 4px 8px;
  font-size: 12px;
  font-weight: 650;
}

.meta-row {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 8px;
  color: var(--subtle);
  font-size: 12px;
}

.metric-row {
  display: flex;
  gap: 10px;
  align-items: center;
  flex: 0 0 auto;
}

.favorite-marker {
  display: inline-grid;
  place-items: center;
  color: var(--brand);
}

.favorite-marker svg {
  width: 14px;
  height: 14px;
}
</style>
