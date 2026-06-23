import { fireEvent, render, screen, waitFor } from '@testing-library/vue'
import { createPinia } from 'pinia'
import { createRouter, createWebHistory } from 'vue-router'
import RegisterView from '../RegisterView.vue'

test('shows a helpful backend connection error instead of raw fetch failure', async () => {
  vi.stubGlobal(
    'fetch',
    vi.fn(async () => {
      throw new TypeError('Failed to fetch')
    }),
  )
  const router = createRouter({
    history: createWebHistory(),
    routes: [{ path: '/', component: { template: '<div />' } }],
  })

  render(RegisterView, {
    global: {
      plugins: [createPinia(), router],
    },
  })

  await fireEvent.update(screen.getByLabelText('昵称'), '新同学')
  await fireEvent.update(screen.getByLabelText('邮箱'), 'new@example.com')
  await fireEvent.update(screen.getByLabelText('密码'), 'correct horse battery staple')
  await fireEvent.click(screen.getByRole('button', { name: '注册' }))

  await waitFor(() => {
    expect(screen.getByText('无法连接后端服务，请确认后端已启动。')).toBeTruthy()
  })

  vi.unstubAllGlobals()
})
