import { render, screen } from '@testing-library/vue'
import { createPinia } from 'pinia'
import { createRouter, createWebHistory } from 'vue-router'
import ProfileView from '../ProfileView.vue'

test('does not expand profile editor or show demo profile values by default', async () => {
  vi.stubGlobal(
    'fetch',
    vi.fn(async () => new Response(JSON.stringify({ error: { code: 'UNAUTHENTICATED', message: 'Authentication required' } }), { status: 401 })),
  )
  const router = createRouter({
    history: createWebHistory(),
    routes: [{ path: '/login', component: { template: '<div />' } }],
  })

  render(ProfileView, {
    global: {
      plugins: [createPinia(), router],
    },
  })

  expect(screen.queryByText('完善应届生资料')).toBeNull()
  expect(screen.queryByDisplayValue('准备秋招的同学')).toBeNull()
  expect(screen.queryByDisplayValue('上海某高校')).toBeNull()

  vi.unstubAllGlobals()
})
