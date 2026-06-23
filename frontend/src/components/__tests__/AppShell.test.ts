import { render, screen } from '@testing-library/vue'
import { createRouter, createWebHistory } from 'vue-router'
import AppShell from '../AppShell.vue'

function renderShell() {
  const router = createRouter({
    history: createWebHistory(),
    routes: [
      { path: '/', component: { template: '<div />' } },
      { path: '/compose', component: { template: '<div />' } },
      { path: '/favorites', component: { template: '<div />' } },
      { path: '/profile', component: { template: '<div />' } },
      { path: '/moderation', component: { template: '<div />' } },
    ],
  })

  return render(AppShell, {
    global: {
      plugins: [router],
    },
  })
}

test('renders bottom navigation', () => {
  renderShell()
  expect(screen.getByText('首页')).toBeTruthy()
  expect(screen.getByText('发帖')).toBeTruthy()
  expect(screen.getByText('收藏')).toBeTruthy()
  expect(screen.getByText('我的')).toBeTruthy()
})

test('links favorites to an independent page', () => {
  renderShell()

  expect(screen.getByText('收藏').closest('a')?.getAttribute('href')).toBe('/favorites')
  expect(screen.getByText('我的').closest('a')?.getAttribute('href')).toBe('/profile')
})
