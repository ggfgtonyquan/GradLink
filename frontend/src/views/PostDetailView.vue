<template>
  <article class="page">
    <RouterLink to="/" class="back-link">返回首页</RouterLink>

    <p v-if="loading" class="hint">正在加载帖子...</p>
    <p v-if="message" class="hint">{{ message }}</p>

    <section v-if="post" class="surface detail">
      <h1 class="screen-title">{{ post.title }}</h1>
      <p class="hint">{{ authorLine }} · {{ post.createdAt }}</p>
      <div class="tag-row">
        <span v-for="tag in post.tags" :key="tag.id" class="tag">{{ tag.name }}</span>
      </div>
      <p class="article">{{ post.body }}</p>

      <div class="action-bar">
        <button class="secondary-button" type="button" @click="likePost">
          <ThumbsUp aria-hidden="true" />
          赞 {{ post.counts.likes }}
        </button>
        <button class="secondary-button" type="button" @click="favoritePost">
          <Bookmark aria-hidden="true" />
          收藏
        </button>
        <button class="secondary-button" type="button">
          <Flag aria-hidden="true" />
          举报
        </button>
      </div>
    </section>

    <section v-if="post" class="comments">
      <h2>评论 {{ post.counts.comments }}</h2>
      <article v-for="comment in comments" :key="comment.id" class="comment surface">
        <p class="hint">{{ commentAuthorLine(comment) }} · {{ comment.createdAt }}</p>
        <p>{{ comment.body }}</p>
      </article>
      <p v-if="!comments.length" class="empty-state">还没有评论。</p>
    </section>

    <form v-if="post" class="comment-composer" @submit.prevent="submitComment">
      <input v-model="commentBody" aria-label="写评论" :placeholder="commentPlaceholder" />
      <button class="primary-button" type="submit" :disabled="commentSubmitting || !commentBody.trim()">发送</button>
    </form>
  </article>
</template>

<script setup lang="ts">
import { computed, onMounted, ref } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { Bookmark, Flag, ThumbsUp } from '@lucide/vue'
import { api, ApiError } from '../api/http'
import type { Comment, Post } from '../api/types'
import { useAuthStore } from '../stores/auth'

const route = useRoute()
const router = useRouter()
const auth = useAuthStore()
const postId = Number(route.params.id || 1)
const post = ref<Post | null>(null)
const comments = ref<Comment[]>([])
const commentBody = ref('')
const loading = ref(false)
const commentSubmitting = ref(false)
const message = ref('')

const authorLine = computed(() => {
  if (!post.value) return ''
  const school = post.value.author.school ? `${post.value.author.school} · ` : ''
  return `${school}${post.value.author.graduationYear}届${post.value.author.educationLevel}`
})

const commentPlaceholder = computed(() => {
  if (!auth.isLoggedIn) return '登录后评论'
  if (!auth.profileComplete) return '完善资料后评论'
  return '写评论...'
})

onMounted(async () => {
  await auth.loadMe()
  await loadDetail()
})

async function loadDetail() {
  loading.value = true
  message.value = ''
  try {
    const [postResult, commentsResult] = await Promise.all([api.post(postId), api.comments(postId)])
    post.value = postResult
    comments.value = commentsResult.comments
  } catch (error) {
    message.value = error instanceof ApiError ? error.message : '帖子加载失败'
  } finally {
    loading.value = false
  }
}

function commentAuthorLine(comment: Comment) {
  const author = comment.author
  if (!author) return '同学'
  const school = author.school ? `${author.school} · ` : ''
  return `${school}${author.graduationYear}届${author.educationLevel}`
}

async function submitComment() {
  if (!auth.isLoggedIn) {
    await router.push('/login')
    return
  }
  if (!auth.profileComplete) {
    await router.push('/profile')
    return
  }
  if (!post.value || !commentBody.value.trim()) return

  commentSubmitting.value = true
  message.value = ''
  try {
    const comment = await api.createComment(post.value.id, commentBody.value)
    commentBody.value = ''
    if (comment.status === 'pending') {
      message.value = '评论已进入待审核队列。'
    } else {
      comments.value = [...comments.value, comment]
      post.value = {
        ...post.value,
        counts: { ...post.value.counts, comments: comment.counts.comments },
      }
    }
  } catch (error) {
    message.value = error instanceof ApiError ? error.message : '评论失败'
  } finally {
    commentSubmitting.value = false
  }
}

async function likePost() {
  if (!post.value) return
  try {
    const result = (await api.likePost(post.value.id)) as { counts?: Post['counts'] }
    if (result.counts) post.value = { ...post.value, counts: result.counts }
  } catch (error) {
    message.value = error instanceof ApiError ? error.message : '点赞失败'
  }
}

async function favoritePost() {
  if (!post.value) return
  try {
    const result = (await api.favoritePost(post.value.id)) as { counts?: Post['counts'] }
    if (result.counts) post.value = { ...post.value, counts: result.counts }
  } catch (error) {
    message.value = error instanceof ApiError ? error.message : '收藏失败'
  }
}
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

.empty-state {
  margin: 0;
  border: 1px dashed var(--line);
  border-radius: 8px;
  padding: 16px;
  color: var(--muted);
  text-align: center;
  font-size: 14px;
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
