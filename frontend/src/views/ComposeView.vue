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
        <TagChips :tags="tags" :selected-ids="selectedTagIds" @toggle-tag="toggleTag" />
      </div>

      <div v-if="gateMessage" class="gate-box">
        <p>{{ gateMessage }}</p>
        <RouterLink :to="gateLink" class="primary-button">{{ gateAction }}</RouterLink>
      </div>

      <button v-else class="primary-button full" type="submit" :disabled="submitting">
        <Send aria-hidden="true" />
        {{ submitting ? '发布中' : '发布' }}
      </button>
      <p v-if="message" class="hint">{{ message }}</p>
      <p class="hint">如果内容命中审核关键词，会先进入待审核队列。</p>
    </form>
  </section>
</template>

<script setup lang="ts">
import { computed, onMounted, ref } from 'vue'
import { useRouter } from 'vue-router'
import { Send } from '@lucide/vue'
import { api, ApiError } from '../api/http'
import type { Tag } from '../api/types'
import TagChips from '../components/TagChips.vue'
import { seedTags } from '../stores/feed'
import { useAuthStore } from '../stores/auth'

const router = useRouter()
const auth = useAuthStore()
const title = ref('')
const body = ref('')
const tags = ref<Tag[]>(seedTags)
const selectedTagIds = ref<number[]>([])
const submitting = ref(false)
const message = ref('')

const gateMessage = computed(() => {
  if (!auth.isLoggedIn) return '登录后才能发布讨论。'
  if (!auth.profileComplete) return '完善应届生资料后才能发帖。'
  return ''
})
const gateLink = computed(() => (auth.isLoggedIn ? '/profile' : '/login'))
const gateAction = computed(() => (auth.isLoggedIn ? '去完善资料' : '去登录'))

onMounted(async () => {
  await auth.loadMe()
  try {
    tags.value = (await api.tags()).tags
  } catch {
    tags.value = seedTags
  }
})

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
  if (gateMessage.value) return
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

.gate-box {
  display: grid;
  gap: 10px;
  border-radius: 8px;
  background: var(--surface-soft);
  padding: 12px;
}

.gate-box p {
  margin: 0;
  color: var(--muted);
  font-size: 13px;
}
</style>
