<template>
  <section class="page">
    <div class="moderation-head">
      <div>
        <h1 class="screen-title">审核后台</h1>
        <p class="hint">处理待审核内容和举报。</p>
      </div>
      <span v-if="auth.isAdmin" class="status-pill pending">待处理 {{ pendingTotal }}</span>
    </div>

    <div v-if="checkingAuth || loading" class="surface state-box">加载中...</div>

    <div v-else-if="!auth.isAdmin" class="surface state-box">
      <strong>无法访问审核后台</strong>
      <p>请先登录管理员账号。</p>
    </div>

    <template v-else>
      <div class="tabs">
        <button
          class="tab"
          :class="{ active: activeTab === 'moderation' }"
          type="button"
          @click="activeTab = 'moderation'"
        >
          待审核
        </button>
        <button class="tab" :class="{ active: activeTab === 'reports' }" type="button" @click="activeTab = 'reports'">
          举报
        </button>
      </div>

      <p v-if="message" class="hint">{{ message }}</p>

      <div v-if="activeTab === 'moderation'" class="admin-list">
        <article v-for="item in moderationItems" :key="item.id" class="surface admin-row">
          <div class="row-top">
            <h2>命中关键词：{{ item.matchedKeywords || '未记录' }}</h2>
            <span class="status-pill pending">待审核</span>
          </div>
          <p>{{ moderationSummary(item) }}</p>
          <p v-if="item.content?.author" class="meta-line">作者：{{ item.content.author }}</p>
          <div class="actions">
            <button class="primary-button" type="button" :disabled="actingId === item.id" @click="decide(item.id, 'approve')">
              通过
            </button>
            <button class="danger-button" type="button" :disabled="actingId === item.id" @click="decide(item.id, 'reject')">
              拒绝
            </button>
          </div>
        </article>
        <p v-if="!moderationItems.length" class="empty-state">暂无待审核内容。</p>
      </div>

      <div v-else class="admin-list">
        <article v-for="report in reports" :key="report.id" class="surface admin-row">
          <div class="row-top">
            <h2>举报原因：{{ report.reason }}</h2>
            <span class="status-pill risk">举报</span>
          </div>
          <p>{{ report.detail || `${targetLabel(report.targetType)} #${report.targetId}` }}</p>
          <div class="actions">
            <button class="danger-button" type="button" :disabled="actingReportId === report.id" @click="handleReport(report.id, 'resolve')">
              标记处理
            </button>
            <button class="secondary-button" type="button" :disabled="actingReportId === report.id" @click="handleReport(report.id, 'reject')">
              驳回举报
            </button>
          </div>
        </article>
        <p v-if="!reports.length" class="empty-state">暂无待处理举报。</p>
      </div>
    </template>
  </section>
</template>

<script setup lang="ts">
import { computed, onMounted, ref } from 'vue'
import { api, ApiError } from '../api/http'
import type { ModerationItem, ReportItem } from '../api/types'
import { useAuthStore } from '../stores/auth'

const auth = useAuthStore()
const checkingAuth = ref(false)
const loading = ref(false)
const message = ref('')
const activeTab = ref<'moderation' | 'reports'>('moderation')
const moderationItems = ref<ModerationItem[]>([])
const reports = ref<ReportItem[]>([])
const actingId = ref<number | null>(null)
const actingReportId = ref<number | null>(null)

const pendingTotal = computed(() => moderationItems.value.length + reports.value.length)

onMounted(async () => {
  checkingAuth.value = true
  if (!auth.user) {
    await auth.loadMe()
  }
  checkingAuth.value = false
  if (auth.isAdmin) {
    await loadQueues()
  }
})

async function loadQueues() {
  loading.value = true
  message.value = ''
  try {
    const [moderationResult, reportResult] = await Promise.all([api.moderation(), api.reports()])
    moderationItems.value = moderationResult.items
    reports.value = reportResult.reports.filter((report) => report.status === 'pending')
  } catch (error) {
    moderationItems.value = []
    reports.value = []
    message.value = error instanceof ApiError ? error.message : '审核数据加载失败'
  } finally {
    loading.value = false
  }
}

async function decide(id: number, action: 'approve' | 'reject') {
  actingId.value = id
  message.value = ''
  try {
    if (action === 'approve') {
      await api.approveModeration(id)
    } else {
      await api.rejectModeration(id)
    }
    await loadQueues()
  } catch (error) {
    message.value = error instanceof ApiError ? error.message : '审核操作失败'
  } finally {
    actingId.value = null
  }
}

async function handleReport(id: number, action: 'resolve' | 'reject') {
  actingReportId.value = id
  message.value = ''
  try {
    if (action === 'resolve') {
      await api.resolveReport(id)
    } else {
      await api.rejectReport(id)
    }
    await loadQueues()
  } catch (error) {
    message.value = error instanceof ApiError ? error.message : '举报处理失败'
  } finally {
    actingReportId.value = null
  }
}

function moderationSummary(item: ModerationItem) {
  const title = item.content?.title
  const body = item.content?.body
  if (title && body) return `${title}：${body}`
  return title || body || `${targetLabel(item.targetType)} #${item.targetId}`
}

function targetLabel(targetType: string) {
  return targetType === 'comment' ? '评论' : '帖子'
}
</script>

<style scoped>
.moderation-head,
.row-top,
.actions,
.tabs {
  display: flex;
  align-items: center;
  gap: 8px;
}

.moderation-head,
.row-top {
  justify-content: space-between;
}

.tabs {
  overflow-x: auto;
}

.tab {
  border-radius: 999px;
  background: var(--surface-soft);
  color: var(--muted);
  padding: 7px 11px;
  font-size: 13px;
  font-weight: 700;
}

.tab.active {
  color: var(--brand);
  background: var(--brand-soft);
}

.admin-list {
  display: grid;
  gap: 10px;
}

.admin-row {
  display: grid;
  gap: 9px;
  padding: 12px;
}

.admin-row h2 {
  margin: 0;
  font-size: 14px;
}

.admin-row p,
.state-box p {
  margin: 0;
  color: #42504b;
  font-size: 14px;
  line-height: 1.55;
}

.meta-line {
  color: var(--muted) !important;
  font-size: 12px !important;
}

.state-box,
.empty-state {
  border-radius: 8px;
  padding: 16px;
}

.state-box {
  display: grid;
  gap: 6px;
}

.empty-state {
  margin: 0;
  border: 1px dashed var(--line);
  color: var(--muted);
  text-align: center;
  font-size: 14px;
}
</style>
