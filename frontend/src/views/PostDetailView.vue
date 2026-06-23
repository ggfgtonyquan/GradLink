<template>
  <article class="page">
    <RouterLink to="/" class="back-link">返回首页</RouterLink>

    <section class="surface detail">
      <h1 class="screen-title">{{ post.title }}</h1>
      <p class="hint">{{ authorLine }} · {{ post.createdAt }}</p>
      <div class="tag-row">
        <span v-for="tag in post.tags" :key="tag.id" class="tag">{{ tag.name }}</span>
      </div>
      <p class="article">{{ post.body }}</p>

      <div class="action-bar">
        <button class="secondary-button" type="button">
          <ThumbsUp aria-hidden="true" />
          赞 {{ post.counts.likes }}
        </button>
        <button class="secondary-button" type="button">
          <Bookmark aria-hidden="true" />
          收藏
        </button>
        <button class="secondary-button" type="button">
          <Flag aria-hidden="true" />
          举报
        </button>
      </div>
    </section>

    <section class="comments">
      <h2>评论 {{ post.counts.comments }}</h2>
      <article class="comment surface">
        <p class="hint">2026届本科 · 软件工程</p>
        <p>可以按背景、难点、方案、结果、复盘准备，每段都配一个可以量化的点。</p>
        <div class="reply">回复：这个结构很清楚，我试一下。</div>
      </article>
      <article class="comment surface">
        <p class="hint">2025届硕士 · 后端方向</p>
        <p>准备一版 2 分钟讲法和一版 8 分钟讲法，先讲主线，再根据追问展开。</p>
      </article>
    </section>

    <form class="comment-composer" @submit.prevent>
      <input aria-label="写评论" placeholder="写评论..." />
      <button class="primary-button" type="submit">发送</button>
    </form>
  </article>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { useRoute } from 'vue-router'
import { Bookmark, Flag, ThumbsUp } from '@lucide/vue'
import { seedPosts } from '../stores/feed'

const route = useRoute()
const postId = Number(route.params.id || 1)
const post = computed(() => seedPosts.find((item) => item.id === postId) ?? seedPosts[0])
const authorLine = computed(() => `${post.value.author.graduationYear}届${post.value.author.educationLevel}`)
</script>

<style scoped>
.back-link {
  color: var(--brand);
  font-size: 13px;
  font-weight: 700;
}

.detail,
.comment {
  padding: 14px;
}

.tag-row {
  display: flex;
  flex-wrap: wrap;
  gap: 6px;
  margin: 12px 0;
}

.tag {
  border-radius: 999px;
  background: var(--brand-soft);
  color: var(--brand);
  padding: 4px 8px;
  font-size: 12px;
  font-weight: 650;
}

.article {
  margin: 14px 0;
  color: #27332f;
  font-size: 15px;
  line-height: 1.75;
}

.action-bar {
  display: grid;
  grid-template-columns: repeat(3, 1fr);
  gap: 8px;
}

.action-bar svg {
  width: 16px;
  height: 16px;
}

.comments {
  display: grid;
  gap: 10px;
}

.comments h2 {
  margin: 4px 0 0;
  font-size: 15px;
}

.comment p {
  margin: 6px 0;
  color: #33413c;
  font-size: 14px;
  line-height: 1.55;
}

.reply {
  margin-top: 8px;
  border-radius: 8px;
  background: var(--surface-soft);
  color: var(--muted);
  padding: 10px;
  font-size: 13px;
  line-height: 1.5;
}

.comment-composer {
  position: sticky;
  bottom: 72px;
  display: flex;
  gap: 8px;
  border: 1px solid var(--line);
  border-radius: 8px;
  background: rgb(255 255 255 / 96%);
  padding: 10px;
}

.comment-composer input {
  min-width: 0;
  flex: 1;
  border: 1px solid var(--line);
  border-radius: 999px;
  padding: 0 14px;
}
</style>
