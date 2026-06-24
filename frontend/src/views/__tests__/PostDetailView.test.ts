import { fireEvent, render, screen, waitFor } from '@testing-library/vue'
import { createPinia } from 'pinia'
import { createRouter, createWebHistory } from 'vue-router'
import PostDetailView from '../PostDetailView.vue'

function renderDetail(fetchMock: ReturnType<typeof vi.fn>) {
  vi.stubGlobal('fetch', fetchMock)
  const router = createRouter({
    history: createWebHistory(),
    routes: [
      { path: '/', component: { template: '<div />' } },
      { path: '/login', component: { template: '<div />' } },
      { path: '/profile', component: { template: '<div />' } },
      { path: '/posts/:id', component: PostDetailView },
    ],
  })
  router.push('/posts/1')

  return render(PostDetailView, {
    global: {
      plugins: [createPinia(), router],
    },
  })
}

afterEach(() => {
  vi.unstubAllGlobals()
})

function postResponse(favorited: boolean) {
  return {
    id: 1,
    title: '收藏反馈测试',
    body: '正文',
    excerpt: '正文',
    status: 'published',
    tags: [{ id: 1, name: '面试经验', slug: 'interview-experience', displayOrder: 1 }],
    author: { nickname: '作者', school: null, educationLevel: '本科', graduationYear: 2026 },
    counts: { comments: 0, likes: 0, favorites: favorited ? 1 : 0 },
    favorited,
    createdAt: '刚刚',
  }
}

function currentUserResponse() {
  return {
    id: 1,
    email: 'student@example.com',
    nickname: '同学',
    role: 'user',
    status: 'active',
    profileComplete: true,
  }
}

test('marks the favorite button as selected after a successful favorite action', async () => {
  const fetchMock = vi.fn(async (input: RequestInfo | URL) => {
    const url = String(input)
    if (url === '/api/me') {
      return new Response(JSON.stringify(currentUserResponse()), { status: 200 })
    }
    if (url === '/api/posts/1/comments') {
      return new Response(JSON.stringify({ comments: [] }), { status: 200 })
    }
    if (url === '/api/posts/1/favorite') {
      return new Response(JSON.stringify({ counts: { comments: 0, likes: 0, favorites: 1 } }), { status: 200 })
    }
    return new Response(JSON.stringify(postResponse(false)), { status: 200 })
  })

  renderDetail(fetchMock)

  const favoriteButton = await screen.findByRole('button', { name: '收藏' })
  expect(favoriteButton.getAttribute('aria-pressed')).toBe('false')

  await fireEvent.click(favoriteButton)

  await waitFor(() => {
    expect(favoriteButton.getAttribute('aria-pressed')).toBe('true')
    expect(favoriteButton.classList.contains('is-favorited')).toBe(true)
  })
})

test('renders an already favorited post as selected and unfavorites on click', async () => {
  const fetchMock = vi.fn(async (input: RequestInfo | URL, init?: RequestInit) => {
    const url = String(input)
    if (url === '/api/me') {
      return new Response(JSON.stringify(currentUserResponse()), { status: 200 })
    }
    if (url === '/api/posts/1/comments') {
      return new Response(JSON.stringify({ comments: [] }), { status: 200 })
    }
    if (url === '/api/posts/1/favorite' && init?.method === 'DELETE') {
      return new Response(JSON.stringify({ counts: { comments: 0, likes: 0, favorites: 0 } }), { status: 200 })
    }
    return new Response(JSON.stringify(postResponse(true)), { status: 200 })
  })

  renderDetail(fetchMock)

  const favoriteButton = await screen.findByRole('button', { name: '已收藏' })
  expect(favoriteButton.getAttribute('aria-pressed')).toBe('true')
  expect(favoriteButton.classList.contains('is-favorited')).toBe(true)

  await fireEvent.click(favoriteButton)

  await waitFor(() => {
    expect(fetchMock).toHaveBeenCalledWith(
      '/api/posts/1/favorite',
      expect.objectContaining({ method: 'DELETE', credentials: 'include' }),
    )
    expect(favoriteButton.getAttribute('aria-pressed')).toBe('false')
    expect(favoriteButton.classList.contains('is-favorited')).toBe(false)
  })
})
