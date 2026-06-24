import { fireEvent, render, screen, waitFor } from '@testing-library/vue'
import { createPinia, setActivePinia } from 'pinia'
import ModerationView from '../ModerationView.vue'
import { useAuthStore } from '../../stores/auth'

function renderModeration(fetchMock: ReturnType<typeof vi.fn>) {
  vi.stubGlobal('fetch', fetchMock)
  const pinia = createPinia()
  setActivePinia(pinia)
  const auth = useAuthStore()
  auth.user = {
    id: 1,
    email: 'admin@example.com',
    nickname: '管理员',
    role: 'admin',
    status: 'active',
    profileComplete: true,
  }

  return render(ModerationView, {
    global: {
      plugins: [pinia],
    },
  })
}

afterEach(() => {
  vi.unstubAllGlobals()
})

test('loads moderation queue from the admin API instead of showing demo rows', async () => {
  renderModeration(
    vi.fn(async (input: RequestInfo | URL) => {
      const url = String(input)
      if (url === '/api/admin/reports') {
        return new Response(JSON.stringify({ reports: [] }), { status: 200 })
      }
      return new Response(JSON.stringify({ items: [] }), { status: 200 })
    }),
  )

  await waitFor(() => expect(screen.getByText('暂无待审核内容。')).toBeTruthy())
  expect(screen.queryByText('命中关键词：加微信')).toBeNull()
  expect(screen.queryByText('举报原因：广告引流')).toBeNull()
})

test('approves a pending moderation item through the admin API', async () => {
  let approved = false
  const fetchMock = vi.fn(async (input: RequestInfo | URL) => {
    const url = String(input)
    if (url === '/api/admin/reports') {
      return new Response(JSON.stringify({ reports: [] }), { status: 200 })
    }
    if (url === '/api/admin/moderation/7/approve') {
      approved = true
      return new Response(JSON.stringify({ id: 7, status: 'published' }), { status: 200 })
    }
    return new Response(
      JSON.stringify({
        items: approved
          ? []
          : [
              {
                id: 7,
                targetType: 'post',
                targetId: 42,
                matchedKeywords: '内推',
                status: 'pending',
                content: {
                  title: '求内推建议',
                  body: '想问问大家怎么判断内推是否靠谱。',
                  author: '应届生',
                },
              },
            ],
      }),
      { status: 200 },
    )
  })

  renderModeration(fetchMock)

  expect(await screen.findByText('命中关键词：内推')).toBeTruthy()
  await fireEvent.click(screen.getByRole('button', { name: '通过' }))

  await waitFor(() => {
    expect(fetchMock).toHaveBeenCalledWith(
      '/api/admin/moderation/7/approve',
      expect.objectContaining({ method: 'POST', credentials: 'include' }),
    )
  })
})
