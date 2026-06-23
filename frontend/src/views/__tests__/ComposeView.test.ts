import { render, screen } from '@testing-library/vue'
import ComposeView from '../ComposeView.vue'

test('shows compose form and moderation hint', () => {
  render(ComposeView)

  expect(screen.getByLabelText('标题')).toBeTruthy()
  expect(screen.getByLabelText('正文')).toBeTruthy()
  expect(screen.getByText('Offer 比较')).toBeTruthy()
  expect(screen.getByText('如果内容命中审核关键词，会先进入待审核队列。')).toBeTruthy()
})
