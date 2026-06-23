import { render, screen } from '@testing-library/vue'
import HomeFeedView from '../HomeFeedView.vue'

test('shows mobile feed cards and filter chips', () => {
  render(HomeFeedView)

  expect(screen.getByText('最新')).toBeTruthy()
  expect(screen.getByText('热门')).toBeTruthy()
  expect(screen.getByText('面试经验')).toBeTruthy()
  expect(screen.getByText('秋招面试被问到项目细节，怎么准备更稳？')).toBeTruthy()
})
