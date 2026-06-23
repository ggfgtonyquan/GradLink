<template>
  <section class="page">
    <div>
      <h1 class="screen-title">发布讨论</h1>
      <p class="hint">写清楚背景、问题和你希望大家帮你判断的点。</p>
    </div>

    <form class="surface compose-form" @submit.prevent="submit">
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
        {{ submitting ? '发布中' : '发布' }}
      </button>
      <p v-if="message" class="hint">{{ message }}</p>
      <p class="hint">如果内容命中审核关键词，会先进入待审核队列。</p>
    </form>
  </section>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { Send } from '@lucide/vue'
import { api, ApiError } from '../api/http'
import TagChips from '../components/TagChips.vue'
import { seedTags } from '../stores/feed'

const router = useRouter()
const title = ref('两个 offer 怎么选：互联网测试开发 vs 银行科技岗')
const body = ref('base 差距不大，一个成长快但强度高，一个稳定但技术栈偏传统。我比较看重前三年的成长，也担心强度太高影响状态。')
const selectedTagIds = ref<number[]>([2, 7])
const submitting = ref(false)
const message = ref('')

function toggleTag(tagId: number) {
  if (selectedTagIds.value.includes(tagId)) {
    selectedTagIds.value = selectedTagIds.value.filter((id) => id !== tagId)
    return
  }
  if (selectedTagIds.value.length < 3) {
    selectedTagIds.value = [...selectedTagIds.value, tagId]
  }
}

async function submit() {
  submitting.value = true
  message.value = ''
  try {
    const post = await api.createPost({
      title: title.value,
      body: body.value,
      tagIds: selectedTagIds.value,
    })
    message.value = post.status === 'pending' ? '内容已进入待审核队列。' : '内容已发布。'
    await router.push(`/posts/${post.id}`)
  } catch (error) {
    message.value = error instanceof ApiError ? error.message : '发布失败'
  } finally {
    submitting.value = false
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
