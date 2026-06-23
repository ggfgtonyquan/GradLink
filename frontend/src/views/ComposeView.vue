<template>
  <section class="page">
    <div>
      <h1 class="screen-title">发布讨论</h1>
      <p class="hint">写清楚背景、问题和你希望大家帮你判断的点。</p>
    </div>

    <form class="surface compose-form" @submit.prevent>
      <div class="field">
        <label for="post-title">标题</label>
        <input id="post-title" v-model="title" aria-label="标题" />
      </div>

      <div class="field">
        <label for="post-body">正文</label>
        <textarea id="post-body" v-model="body" aria-label="正文" />
      </div>

      <div class="field">
        <label>选择标签</label>
        <TagChips :tags="seedTags" :selected-ids="selectedTagIds" @toggle-tag="toggleTag" />
      </div>

      <button class="primary-button full" type="submit">
        <Send aria-hidden="true" />
        发布
      </button>
      <p class="hint">如果内容命中审核关键词，会先进入待审核队列。</p>
    </form>
  </section>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { Send } from '@lucide/vue'
import TagChips from '../components/TagChips.vue'
import { seedTags } from '../stores/feed'

const title = ref('两个 offer 怎么选：互联网测试开发 vs 银行科技岗')
const body = ref('base 差距不大，一个成长快但强度高，一个稳定但技术栈偏传统。我比较看重前三年的成长，也担心强度太高影响状态。')
const selectedTagIds = ref<number[]>([2, 7])

function toggleTag(tagId: number) {
  if (selectedTagIds.value.includes(tagId)) {
    selectedTagIds.value = selectedTagIds.value.filter((id) => id !== tagId)
    return
  }
  if (selectedTagIds.value.length < 3) {
    selectedTagIds.value = [...selectedTagIds.value, tagId]
  }
}
</script>

<style scoped>
.compose-form {
  padding: 16px;
}

.full {
  width: 100%;
}

.primary-button svg {
  width: 17px;
  height: 17px;
}
</style>
