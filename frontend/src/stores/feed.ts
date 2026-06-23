import { computed, ref } from 'vue'
import { defineStore } from 'pinia'
import { api } from '../api/http'
import type { Post, Tag } from '../api/types'

export const seedTags: Tag[] = [
  { id: 1, name: '面试经验', slug: 'interview-experience', displayOrder: 10 },
  { id: 2, name: 'Offer 比较', slug: 'offer-comparison', displayOrder: 20 },
  { id: 3, name: '秋招进度', slug: 'autumn-recruiting', displayOrder: 30 },
  { id: 5, name: '简历求助', slug: 'resume-help', displayOrder: 50 },
  { id: 7, name: '签约三方', slug: 'tripartite-agreement', displayOrder: 70 },
  { id: 9, name: '综合提问', slug: 'general-questions', displayOrder: 90 },
]

export const seedPosts: Post[] = [
  {
    id: 1,
    title: '秋招面试被问到项目细节，怎么准备更稳？',
    excerpt: '最近投后端开发，面试官一直追问项目里的并发和数据库。我能讲功能，但讲到取舍就容易卡住。',
    body: '最近投后端开发，发现面试官很喜欢从项目细节开始追问。我能讲清楚功能，但问到为什么这样设计、数据量上来怎么办、有没有替代方案，就会回答得比较散。',
    status: 'published',
    tags: [seedTags[0], seedTags[2]],
    author: { nickname: '准备秋招的同学', school: null, educationLevel: '硕士', graduationYear: 2026 },
    counts: { comments: 12, likes: 38, favorites: 9 },
    createdAt: '2小时前',
  },
  {
    id: 2,
    title: '两个 offer 怎么选：互联网测试开发 vs 银行科技岗',
    excerpt: 'base 差距不大，一个成长快但强度高，一个稳定但技术栈偏传统。想听听大家怎么看。',
    body: '我比较看重前三年的成长，也担心强度太高影响状态，想听听大家怎么比较。',
    status: 'published',
    tags: [seedTags[1], seedTags[4]],
    author: { nickname: 'Offer 选择困难', school: null, educationLevel: '本科', graduationYear: 2026 },
    counts: { comments: 26, likes: 51, favorites: 18 },
    createdAt: '4小时前',
  },
  {
    id: 3,
    title: '简历上实习经历比较少，要不要把课程项目写详细？',
    excerpt: '目前只有一段校内项目和比赛经历，准备投产品和运营，担心简历看起来太空。',
    body: '希望大家帮我看看课程项目应该怎么写，哪些细节值得展开。',
    status: 'published',
    tags: [seedTags[3], seedTags[5]],
    author: { nickname: '改简历中', school: null, educationLevel: '本科', graduationYear: 2026 },
    counts: { comments: 8, likes: 17, favorites: 6 },
    createdAt: '今天',
  },
]

export const useFeedStore = defineStore('feed', () => {
  const posts = ref<Post[]>(seedPosts)
  const tags = ref<Tag[]>(seedTags)
  const sort = ref<'latest' | 'hot'>('latest')
  const activeTagId = ref<number | null>(null)
  const loading = ref(false)

  const visiblePosts = computed(() => {
    const filtered =
      activeTagId.value === null
        ? posts.value
        : posts.value.filter((post) => post.tags.some((tag) => tag.id === activeTagId.value))
    return [...filtered].sort((a, b) =>
      sort.value === 'hot' ? b.counts.likes + b.counts.comments - (a.counts.likes + a.counts.comments) : a.id - b.id,
    )
  })

  async function load() {
    loading.value = true
    try {
      const [tagResult, postResult] = await Promise.all([api.tags(), api.posts(`?sort=${sort.value}`)])
      tags.value = tagResult.tags
      posts.value = postResult.posts
    } catch {
      posts.value = seedPosts
      tags.value = seedTags
    } finally {
      loading.value = false
    }
  }

  return { posts, tags, sort, activeTagId, visiblePosts, loading, load }
})
