export type UserRole = 'user' | 'admin'
export type ContentStatus = 'published' | 'pending' | 'hidden' | 'rejected' | 'deleted'

export interface Tag {
  id: number
  name: string
  slug: string
  displayOrder: number
}

export interface Author {
  nickname: string
  school: string | null
  educationLevel: string
  graduationYear: number
}

export interface Counts {
  comments: number
  likes: number
  favorites: number
}

export interface Post {
  id: number
  title: string
  body?: string
  excerpt: string
  status: ContentStatus
  tags: Tag[]
  author: Author
  counts: Counts
  createdAt: string
}

export interface Comment {
  id: number
  postId: number
  parentCommentId: number | null
  body: string
  status: ContentStatus
  counts: Counts
}

export interface CurrentUser {
  id: number
  email: string
  nickname: string
  role: UserRole
  status: string
  profileComplete?: boolean
}

export interface ApiErrorPayload {
  error: {
    code: string
    message: string
  }
}

export interface ModerationItem {
  id: number
  targetType: 'post' | 'comment'
  targetId: number
  matchedKeywords: string
  status: string
}
