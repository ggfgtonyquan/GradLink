import { render, screen } from '@testing-library/vue'
import AppShell from '../AppShell.vue'

test('renders bottom navigation', () => {
  render(AppShell)
  expect(screen.getByText('首页')).toBeTruthy()
  expect(screen.getByText('发帖')).toBeTruthy()
  expect(screen.getByText('收藏')).toBeTruthy()
  expect(screen.getByText('我的')).toBeTruthy()
})
