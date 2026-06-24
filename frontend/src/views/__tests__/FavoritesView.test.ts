import { render, screen } from '@testing-library/vue'
import { createPinia } from 'pinia'
import { createRouter, createWebHistory } from 'vue-router'
import FavoritesView from '../FavoritesView.vue'

afterEach(() => {
  vi.unstubAllGlobals()
})

test('shows a filled favorite marker for posts loaded from my favorites', async () => {
  vi.stubGlobal(
    'fetch',
    vi.fn(async (input: RequestInfo | URL) => {
      const url = String(input)
      if (url === '/api/me') {
        return new Response(
          JSON.stringify({
            id: 1,
            email: 'student@example.com',
            nickname: '同学',
            role: 'user',
            status: 'active',
            profileComplete: true,
          }),
          { status: 200 },
        )
      }
      if (url === '/api/me/favorites') {
        return new Response(
          JSON.stringify({
            posts: [
              {
                id: 1,
                title: '收藏的帖子',
                body: '正文',
                excerpt: '正文',
                status: 'published',
                tags: [{ id: 1, name: '面试经验', slug: 'interview-experience', displayOrder: 1 }],
                author: { nickname: '作者', school: null, educationLevel: '本科', graduationYear: 2026 },
                counts: { comments: 0, likes: 0, favorites: 1 },
                favorited: true,
                createdAt: '刚刚',
              },
            ],
          }),
          { status: 200 },
        )
      }
      return new Response(JSON.stringify({ error: { code: 'NOT_FOUND', message: 'Not found' } }), { status: 404 })
    }),
  )
  const router = createRouter({
    history: createWebHistory(),
    routes: [
      { path: '/login', component: { template: '<div />' } },
      { path: '/posts/:id', component: { template: '<div />' } },
    ],
  })

  render(FavoritesView, {
    global: {
      plugins: [createPinia(), router],
    },
  })

  expect(await screen.findByText('收藏的帖子')).toBeTruthy()
  expect(screen.getByLabelText('已收藏')).toBeTruthy()
})
