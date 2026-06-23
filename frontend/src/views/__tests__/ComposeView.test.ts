import { render, screen } from '@testing-library/vue'
import { createPinia } from 'pinia'
import { createRouter, createWebHistory } from 'vue-router'
import ComposeView from '../ComposeView.vue'

function renderCompose() {
  const router = createRouter({
    history: createWebHistory(),
    routes: [
      { path: '/', component: { template: '<div />' } },
      { path: '/login', component: { template: '<div />' } },
      { path: '/profile', component: { template: '<div />' } },
      { path: '/posts/:id', component: { template: '<div />' } },
    ],
  })

  return render(ComposeView, {
    global: {
      plugins: [createPinia(), router],
    },
  })
}

test('shows compose form and moderation hint', () => {
  renderCompose()

  expect(screen.getByLabelText('标题')).toBeTruthy()
  expect(screen.getByLabelText('正文')).toBeTruthy()
  expect(screen.getByText('Offer 比较')).toBeTruthy()
  expect(screen.getByText('如果内容命中审核关键词，会先进入待审核队列。')).toBeTruthy()
})

test('starts with empty title and body', () => {
  renderCompose()

  expect((screen.getByLabelText('标题') as HTMLInputElement).value).toBe('')
  expect((screen.getByLabelText('正文') as HTMLTextAreaElement).value).toBe('')
})
