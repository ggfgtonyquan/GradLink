import type { Comment, CurrentUser, ModerationItem, Post, Tag } from './types'

const jsonHeaders = { 'Content-Type': 'application/json' }

export class ApiError extends Error {
  code: string

  constructor(code: string, message: string) {
    super(message)
    this.code = code
  }
}

async function request<T>(path: string, init: RequestInit = {}): Promise<T> {
  const response = await fetch(path, {
    credentials: 'include',
    ...init,
    headers: {
      ...(init.body ? jsonHeaders : {}),
      ...init.headers,
    },
  })

  const text = await response.text()
  const data = text ? JSON.parse(text) : {}
  if (!response.ok) {
    const error = data.error ?? { code: 'REQUEST_FAILED', message: '请求失败' }
    throw new ApiError(error.code, error.message)
  }
  return data as T
}

export const api = {
  me: () => request<CurrentUser>('/api/me'),
  login: (email: string, password: string) =>
    request<CurrentUser>('/api/auth/login', {
      method: 'POST',
      body: JSON.stringify({ email, password }),
    }),
  register: (email: string, password: string, nickname: string) =>
    request<CurrentUser>('/api/auth/register', {
      method: 'POST',
      body: JSON.stringify({ email, password, nickname }),
    }),
  logout: () => request<{ ok: boolean }>('/api/auth/logout', { method: 'POST' }),
  saveProfile: (payload: Record<string, unknown>) =>
    request<{ profileComplete: boolean }>('/api/me/profile', {
      method: 'PUT',
      body: JSON.stringify(payload),
    }),
  tags: () => request<{ tags: Tag[] }>('/api/tags'),
  posts: (query = '') => request<{ posts: Post[] }>(`/api/posts${query}`),
  post: (id: number) => request<Post>(`/api/posts/${id}`),
  createPost: (payload: { title: string; body: string; tagIds: number[] }) =>
    request<Post>('/api/posts', {
      method: 'POST',
      body: JSON.stringify(payload),
    }),
  createComment: (postId: number, body: string) =>
    request<Comment>(`/api/posts/${postId}/comments`, {
      method: 'POST',
      body: JSON.stringify({ body }),
    }),
  likePost: (postId: number) => request(`/api/posts/${postId}/like`, { method: 'POST' }),
  favoritePost: (postId: number) => request(`/api/posts/${postId}/favorite`, { method: 'POST' }),
  report: (payload: { targetType: string; targetId: number; reason: string; detail: string }) =>
    request('/api/reports', { method: 'POST', body: JSON.stringify(payload) }),
  moderation: () => request<{ items: ModerationItem[] }>('/api/admin/moderation'),
}
