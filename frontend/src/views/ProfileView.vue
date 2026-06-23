<template>
  <section class="page">
    <div class="surface profile-card">
      <div class="avatar">GL</div>
      <div>
        <h1 class="screen-title">准备秋招的同学</h1>
        <p class="hint">2026届硕士 · 计算机方向</p>
      </div>
    </div>

    <form class="surface profile-form" @submit.prevent>
      <h2>完善应届生资料</h2>
      <div class="field">
        <label for="nickname">昵称</label>
        <input id="nickname" value="准备秋招的同学" />
      </div>
      <div class="field">
        <label for="school">学校</label>
        <input id="school" value="上海某高校" />
      </div>
      <div class="field">
        <label for="education">学历</label>
        <select id="education">
          <option>本科</option>
          <option selected>硕士</option>
          <option>博士</option>
        </select>
      </div>
      <div class="field">
        <label for="major">专业或方向</label>
        <input id="major" value="计算机方向" />
      </div>
      <button class="primary-button" type="submit">保存并继续</button>
    </form>

    <div id="favorites" class="menu-list">
      <button class="menu-row" type="button">我的帖子 <span>3</span></button>
      <button class="menu-row" type="button">我的收藏 <span>12</span></button>
      <button class="menu-row" type="button">账号设置 <span></span></button>
    </div>
  </section>
</template>

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
}
</style>
