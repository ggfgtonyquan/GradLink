import { render, screen } from '@testing-library/vue'
import { createPinia, setActivePinia } from 'pinia'
import { createRouter, createWebHistory } from 'vue-router'
import { useAuthStore } from '../../stores/auth'
import AppShell from '../AppShell.vue'

function renderShell(userRole?: 'user' | 'admin') {
  const pinia = createPinia()
  setActivePinia(pinia)
  const auth = useAuthStore()
  auth.user = userRole
    ? {
        id: 1,
        email: `${userRole}@example.com`,
        nickname: userRole,
        role: userRole,
        status: 'active',
        profileComplete: true,
      }
    : null
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
      plugins: [pinia, router],
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

test('shows moderation entry only to admins', () => {
  const { unmount } = renderShell()
  expect(screen.queryByLabelText('审核后台')).toBeNull()
  unmount()

  renderShell('admin')
  expect(screen.getByLabelText('审核后台').getAttribute('href')).toBe('/moderation')
})
