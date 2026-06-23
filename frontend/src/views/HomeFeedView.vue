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
      :tags="feed.tags"
      :selected-ids="feed.activeTagId ? [feed.activeTagId] : []"
      :sort="feed.sort"
      show-sort
      @update:sort="setSort"
      @toggle-tag="toggleTag"
    />

    <div class="feed-list">
      <PostCard v-for="post in feed.visiblePosts" :key="post.id" :post="post" />
      <p v-if="!feed.visiblePosts.length" class="empty-state">暂时没有帖子。</p>
    </div>
  </section>
</template>

<script setup lang="ts">
import { onMounted } from 'vue'
import PostCard from '../components/PostCard.vue'
import TagChips from '../components/TagChips.vue'
import { useFeedStore } from '../stores/feed'

const feed = useFeedStore()

onMounted(() => {
  feed.load()
})

function toggleTag(tagId: number) {
  feed.activeTagId = feed.activeTagId === tagId ? null : tagId
  feed.load()
}

function setSort(sort: 'latest' | 'hot') {
  feed.sort = sort
  feed.load()
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

.empty-state {
  margin: 0;
  border: 1px dashed var(--line);
  border-radius: 8px;
  padding: 16px;
  color: var(--muted);
  text-align: center;
  font-size: 14px;
}
</style>
