import { render, screen } from '@testing-library/vue'
import { createPinia } from 'pinia'
import { createMemoryHistory, createRouter } from 'vue-router'
import LoginView from '../LoginView.vue'

test('starts with empty login credentials', async () => {
  const router = createRouter({
    history: createMemoryHistory(),
    routes: [
      { path: '/', component: LoginView },
      { path: '/register', component: { template: '<div />' } },
    ],
  })

  render(LoginView, {
    global: {
      plugins: [createPinia(), router],
    },
  })

  await router.push('/')
  await router.isReady()

  expect(screen.getByLabelText('邮箱')).toHaveProperty('value', '')
  expect(screen.getByLabelText('密码')).toHaveProperty('value', '')
})
