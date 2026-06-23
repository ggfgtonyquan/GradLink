<template>
  <div class="app-shell">
    <header class="top-bar">
      <div>
        <p class="eyebrow">GradLink</p>
        <h1>应届生交流</h1>
      </div>
      <button class="profile-button" type="button" aria-label="打开个人中心">我</button>
    </header>

    <main class="content">
      <slot>
        <section class="welcome-panel" aria-labelledby="welcome-title">
          <p class="section-label">今日讨论</p>
          <h2 id="welcome-title">找到同届同学的求职经验</h2>
          <p>分享面试、Offer、签约和避坑信息，让毕业季少一点信息差。</p>
        </section>
      </slot>
    </main>

    <nav class="bottom-nav" aria-label="主导航">
      <a href="#" class="nav-item active">首页</a>
      <a href="#" class="nav-item">发帖</a>
      <a href="#" class="nav-item">收藏</a>
      <a href="#" class="nav-item">我的</a>
    </nav>
  </div>
</template>

<style scoped>
:global(*) {
  box-sizing: border-box;
}

:global(body) {
  margin: 0;
  background: #eef3f0;
  color: #18201c;
  font-family:
    Inter, "PingFang SC", "Microsoft YaHei", system-ui, -apple-system, sans-serif;
}

:global(a) {
  color: inherit;
  text-decoration: none;
}

.app-shell {
  min-height: 100vh;
  max-width: 430px;
  margin: 0 auto;
  background: #f8faf7;
  display: flex;
  flex-direction: column;
  box-shadow: 0 0 32px rgb(21 43 33 / 10%);
}

.top-bar {
  position: sticky;
  top: 0;
  z-index: 1;
  display: flex;
  align-items: center;
  justify-content: space-between;
  min-height: 72px;
  padding: 14px 18px;
  background: rgb(248 250 247 / 94%);
  border-bottom: 1px solid #dde6df;
  backdrop-filter: blur(16px);
}

.eyebrow,
.section-label {
  margin: 0;
  color: #527264;
  font-size: 12px;
  font-weight: 700;
}

h1,
h2,
p {
  letter-spacing: 0;
}

h1 {
  margin: 3px 0 0;
  font-size: 21px;
  line-height: 1.2;
}

.profile-button {
  width: 40px;
  height: 40px;
  border: 1px solid #d2ded7;
  border-radius: 50%;
  background: #ffffff;
  color: #315849;
  font-size: 15px;
  font-weight: 700;
}

.content {
  flex: 1;
  padding: 16px 16px 88px;
}

.welcome-panel {
  border: 1px solid #dfe8e2;
  border-radius: 8px;
  background: #ffffff;
  padding: 18px;
}

.welcome-panel h2 {
  margin: 8px 0;
  font-size: 19px;
  line-height: 1.35;
}

.welcome-panel p:last-child {
  margin: 0;
  color: #5f6f68;
  font-size: 14px;
  line-height: 1.7;
}

.bottom-nav {
  position: fixed;
  left: 50%;
  bottom: 0;
  width: min(100%, 430px);
  transform: translateX(-50%);
  display: grid;
  grid-template-columns: repeat(4, 1fr);
  min-height: 64px;
  padding: 8px 8px max(8px, env(safe-area-inset-bottom));
  border-top: 1px solid #d9e4dc;
  background: rgb(255 255 255 / 96%);
  backdrop-filter: blur(18px);
}

.nav-item {
  display: flex;
  align-items: center;
  justify-content: center;
  min-width: 0;
  min-height: 44px;
  border-radius: 8px;
  color: #68776f;
  font-size: 13px;
  font-weight: 700;
}

.nav-item.active {
  background: #e4f3eb;
  color: #226348;
}
</style>
