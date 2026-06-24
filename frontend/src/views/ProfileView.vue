<template>
  <section class="page">
    <div class="surface profile-card">
      <div class="avatar">{{ initials }}</div>
      <div>
        <h1 class="screen-title">{{ auth.user?.nickname ?? '未登录' }}</h1>
        <p class="hint">{{ profileLine }}</p>
      </div>
    </div>

    <div class="menu-list">
      <button class="menu-row" type="button" @click="activePanel = 'posts'">我的帖子 <span>{{ myPosts.length }}</span></button>
      <RouterLink class="menu-row" to="/favorites">我的收藏 <span>{{ myFavorites.length }}</span></RouterLink>
      <button class="menu-row" type="button" @click="togglePanel('settings')">账号设置 <span>{{ auth.user?.email ?? '' }}</span></button>
    </div>

    <section v-if="activePanel === 'posts'" id="posts" class="panel">
      <h2>我的帖子</h2>
      <PostCard v-for="post in myPosts" :key="post.id" :post="post" />
      <p v-if="!myPosts.length" class="empty-state">还没有发布过帖子。</p>
    </section>

    <section v-if="activePanel === 'settings'" class="surface settings-panel">
      <h2>账号设置</h2>
      <form class="profile-form compact" @submit.prevent="submit">
        <div class="field">
          <label for="nickname">昵称</label>
          <input id="nickname" v-model="form.nickname" />
        </div>
        <div class="field">
          <label for="school">学校</label>
          <input id="school" v-model="form.school" />
        </div>
        <div class="field">
          <label for="education">学历</label>
          <select id="education" v-model="form.educationLevel">
            <option>本科</option>
            <option>硕士</option>
            <option>博士</option>
          </select>
        </div>
        <div class="field">
          <label for="major">专业或方向</label>
          <input id="major" v-model="form.major" />
        </div>
        <div class="field">
          <label for="graduation-year">毕业年份</label>
          <input id="graduation-year" v-model.number="form.graduationYear" type="number" min="2020" max="2035" />
        </div>
        <label class="check-row">
          <input v-model="form.schoolVisible" type="checkbox" />
          公开学校信息
        </label>
        <div v-if="!auth.isLoggedIn" class="gate-box">
          <p>登录后可以保存资料、发帖和评论。</p>
          <RouterLink to="/login" class="primary-button">去登录</RouterLink>
        </div>
        <button v-else class="primary-button" type="submit">{{ submitting ? '保存中' : '保存资料' }}</button>
        <p v-if="message" class="hint">{{ message }}</p>
      </form>
      <p class="hint">邮箱：{{ auth.user?.email ?? '未登录' }}</p>
      <button v-if="auth.isLoggedIn" class="secondary-button" type="button" @click="logout">退出登录</button>
    </section>
  </section>
</template>

<script setup lang="ts">
import { computed, onMounted, reactive, ref } from 'vue'
import { useRouter } from 'vue-router'
import { api, ApiError } from '../api/http'
import type { Post } from '../api/types'
import PostCard from '../components/PostCard.vue'
import { useAuthStore } from '../stores/auth'

const auth = useAuthStore()
const router = useRouter()
const submitting = ref(false)
const message = ref('')
const activePanel = ref<'posts' | 'settings' | null>('posts')
const myPosts = ref<Post[]>([])
const myFavorites = ref<Post[]>([])
const form = reactive({
  nickname: '',
  school: '',
  educationLevel: '',
  major: '',
  graduationYear: new Date().getFullYear(),
  schoolVisible: true,
})

const profileLine = computed(() => {
  const profile = auth.user?.profile
  if (!profile) return auth.isLoggedIn ? '请完善应届生资料' : '登录后查看个人中心'
  return `${profile.graduationYear}届${profile.educationLevel} · ${profile.major}`
})

const initials = computed(() => (auth.user?.nickname ?? 'GL').slice(0, 2).toUpperCase())

function syncProfileForm() {
  if (!auth.user) return
  form.nickname = auth.user.nickname
  if (auth.user.profile) {
    form.school = auth.user.profile.school
    form.educationLevel = auth.user.profile.educationLevel
    form.major = auth.user.profile.major
    form.graduationYear = auth.user.profile.graduationYear
    form.schoolVisible = auth.user.profile.schoolVisible
  }
}

async function loadActivity() {
  if (!auth.isLoggedIn) return
  const [postsResult, favoritesResult] = await Promise.all([api.myPosts(), api.myFavorites()])
  myPosts.value = postsResult.posts
  myFavorites.value = favoritesResult.posts
}

function togglePanel(panel: 'posts' | 'settings') {
  activePanel.value = activePanel.value === panel ? null : panel
}

onMounted(async () => {
  await auth.loadMe()
  syncProfileForm()
  await loadActivity()
})

async function submit() {
  if (!auth.isLoggedIn) {
    await router.push('/login')
    return
  }
  submitting.value = true
  message.value = ''
  try {
    const result = await api.saveProfile({ ...form })
    await auth.loadMe()
    syncProfileForm()
    message.value = result.profileComplete ? '资料已完善。' : '资料已保存。'
  } catch (error) {
    message.value = error instanceof ApiError ? error.message : '保存失败'
  } finally {
    submitting.value = false
  }
}

async function logout() {
  await auth.logout()
  myPosts.value = []
  myFavorites.value = []
  await router.push('/login')
}
</script>

<style scoped>
.profile-card,
.profile-form {
  padding: 16px;
}

.profile-card {
  display: flex;
  align-items: center;
  gap: 12px;
}

.avatar {
  width: 48px;
  height: 48px;
  border-radius: 50%;
  display: grid;
  place-items: center;
  background: var(--brand-soft);
  color: var(--brand);
  font-weight: 800;
}

.profile-form h2 {
  margin: 0 0 12px;
  font-size: 16px;
}

.profile-form.compact {
  padding: 0;
}

.check-row {
  display: flex;
  align-items: center;
  gap: 8px;
  margin: 0 0 14px;
  color: var(--muted);
  font-size: 13px;
  font-weight: 700;
}

.gate-box,
.settings-panel {
  display: grid;
  gap: 10px;
}

.gate-box {
  border-radius: 8px;
  background: var(--surface-soft);
  padding: 12px;
  color: var(--muted);
  font-size: 13px;
}

.gate-box p {
  margin: 0;
}

.menu-list {
  display: grid;
  gap: 8px;
}

.menu-row {
  min-height: 48px;
  border: 1px solid var(--line);
  border-radius: 8px;
  background: var(--surface);
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 0 12px;
  color: var(--text);
  font-weight: 650;
}

.menu-row span {
  color: var(--subtle);
  max-width: 180px;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.panel {
  display: grid;
  gap: 10px;
}

.panel h2,
.settings-panel h2 {
  margin: 0;
  font-size: 16px;
}

.settings-panel {
  padding: 16px;
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
