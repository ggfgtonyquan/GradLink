<template>
  <section class="page">
    <div class="moderation-head">
      <div>
        <h1 class="screen-title">审核后台</h1>
        <p class="hint">处理待审核内容、举报、用户和关键词。</p>
      </div>
      <span class="status-pill pending">待处理 2</span>
    </div>

    <div class="tabs">
      <button class="tab active" type="button">待审核</button>
      <button class="tab" type="button">举报</button>
      <button class="tab" type="button">用户</button>
      <button class="tab" type="button">关键词</button>
    </div>

    <article class="surface admin-row">
      <div class="row-top">
        <h2>命中关键词：加微信</h2>
        <span class="status-pill pending">待审核</span>
      </div>
      <p>有人说可以加微信内推，靠谱吗？对方要求加微信并购买培训服务。</p>
      <div class="actions">
        <button class="primary-button" type="button">通过</button>
        <button class="danger-button" type="button">拒绝</button>
      </div>
    </article>

    <article class="surface admin-row">
      <div class="row-top">
        <h2>举报原因：广告引流</h2>
        <span class="status-pill risk">举报</span>
      </div>
      <p>评论区重复发布课程链接，影响正常讨论。</p>
      <div class="actions">
        <button class="danger-button" type="button">隐藏内容</button>
        <button class="secondary-button" type="button">驳回举报</button>
      </div>
    </article>
  </section>
</template>

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

.admin-row {
  display: grid;
  gap: 9px;
  padding: 12px;
}

.admin-row h2 {
  margin: 0;
  font-size: 14px;
}

.admin-row p {
  margin: 0;
  color: #42504b;
  font-size: 14px;
  line-height: 1.55;
}
</style>
