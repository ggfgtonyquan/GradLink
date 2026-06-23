<template>
  <section class="page">
    <div class="home-head">
      <div>
        <p class="eyebrow">今日讨论</p>
        <h1 class="screen-title">应届生求职交流</h1>
      </div>
      <RouterLink to="/compose" class="primary-button">发帖</RouterLink>
    </div>

    <TagChips
      :tags="seedTags"
      :selected-ids="selectedTagId ? [selectedTagId] : []"
      :sort="sort"
      show-sort
      @update:sort="sort = $event"
      @toggle-tag="toggleTag"
    />

    <div class="feed-list">
      <PostCard v-for="post in filteredPosts" :key="post.id" :post="post" />
    </div>
  </section>
</template>

<script setup lang="ts">
import { computed, ref } from 'vue'
import PostCard from '../components/PostCard.vue'
import TagChips from '../components/TagChips.vue'
import { seedPosts, seedTags } from '../stores/feed'

const sort = ref<'latest' | 'hot'>('latest')
const selectedTagId = ref<number | null>(null)

const filteredPosts = computed(() => {
  const filtered =
    selectedTagId.value === null
      ? seedPosts
      : seedPosts.filter((post) => post.tags.some((tag) => tag.id === selectedTagId.value))

  return [...filtered].sort((a, b) => {
    if (sort.value === 'hot') {
      return b.counts.likes + b.counts.comments - (a.counts.likes + a.counts.comments)
    }
    return a.id - b.id
  })
})

function toggleTag(tagId: number) {
  selectedTagId.value = selectedTagId.value === tagId ? null : tagId
}
</script>

<style scoped>
.home-head {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 12px;
}

.eyebrow {
  margin: 0 0 4px;
  color: var(--subtle);
  font-size: 12px;
  font-weight: 700;
}

.feed-list {
  display: grid;
  gap: 10px;
}
</style>
