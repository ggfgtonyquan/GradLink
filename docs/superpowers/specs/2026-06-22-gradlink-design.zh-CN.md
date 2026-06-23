# GradLink 设计方案

日期：2026-06-22

## 概要

GradLink 是一个面向国内应届生的移动端优先求职交流网站。第一版聚焦一个小而完整的社区闭环：浏览帖子、注册登录、完善应届生资料、发帖、评论、回复、点赞、收藏、举报和内容审核。

产品形态是响应式网站，不是原生 App。手机端要像社区 App 一样顺手，桌面端保持可用。后续可以考虑 PWA、小程序或原生 App，但 MVP 阶段不处理应用商店上架、App 审核和移动应用合规带来的额外复杂度。

## 目标

- 帮助国内应届生以半匿名方式交流求职经验。
- 让发帖和讨论足够顺畅，同时要求用户在发言前完善基本应届生资料。
- 控制个人项目的开发、运维和服务器成本，适合全 AI 辅助开发。
- 从第一版开始内置基础审核能力：关键词检查、待审核队列、举报、隐藏内容和封禁用户。
- 保持系统边界足够小，先做出一个真实可用的小社区。

## 暂不做的内容

- 岗位聚合。
- 投递进度管理。
- 公司主页。
- 私信。
- 实时通知。
- 推荐算法。
- 图片上传。
- 原生 App 或应用商店发布。
- 强学生认证、学生证 OCR、学信网认证、强制校园邮箱认证。

## 用户与权限

GradLink 使用半匿名身份模型。

用户状态：

- 游客：可以浏览帖子和评论。
- 注册用户：可以登录、点赞、收藏，并完善应届生资料。
- 已完善资料用户：可以发帖、评论、回复和举报内容。
- 管理员：可以审核待审内容、处理举报、隐藏帖子或评论、封禁用户、维护审核关键词。

发言门槛是“完善应届生资料”，不是“邮箱已验证”。

应届生资料字段：

- 昵称。
- 学校。
- 学历。
- 专业或方向。
- 毕业年份。
- 是否公开学校。

前台展示要保持半匿名。例如根据用户的公开设置展示为“2026 届硕士 · 计算机方向”或“上海高校学生 · 2026 届”，避免默认暴露完整身份。

## 账号与认证

MVP 使用邮箱 + 密码账号体系。

- 邮箱只用于注册和登录。
- 系统只检查邮箱格式和唯一性。
- MVP 不发送邮箱验证邮件。
- MVP 不要求双向邮箱认证。
- 发帖和评论要求用户已登录，并且已完善应届生资料。
- 密码必须哈希后存储。
- 第一版可以不做自助找回密码。早期测试阶段如有需要，可由管理员手动重置账号。

MVP 用轻量反滥用机制替代邮箱验证：

- 注册、登录、发帖、评论的基础频率限制。
- 关键词审核。
- 用户举报。
- 管理员封禁。
- 邮箱唯一约束。

如果网站进入公开使用阶段，再补充邮箱验证、密码找回、手机号验证或校园邮箱认证。

## 产品页面

### 首页信息流

首页是主入口，展示帖子列表。每条帖子包含标题、摘要、标签、作者半匿名信息、评论数、点赞数、收藏状态和发布时间。

信息流支持：

- 最新排序。
- 热门排序。
- 标签筛选。
- 标题和正文的关键词搜索。

初始标签：

- 面试经验。
- Offer 比较。
- 秋招进度。
- 春招进度。
- 简历求助。
- 笔试讨论。
- 签约三方。
- 吐槽避坑。
- 综合提问。

### 帖子详情页

帖子详情页展示：

- 标题。
- 正文。
- 标签。
- 作者半匿名资料。
- 内容状态，只有有权限的用户需要看到待审或隐藏状态。
- 点赞入口。
- 收藏入口。
- 举报入口。
- 评论列表。
- 评论回复入口。

评论支持一层回复，通过父评论 id 关联。MVP 不做无限层级楼中楼。

### 发帖页

发帖页包含：

- 标题输入。
- 正文编辑。
- 标签选择。
- 提交按钮。

提交后，后端执行关键词检查：

- 无风险内容直接发布。
- 命中风险词的内容进入待审核，并提示用户“内容等待审核”。
- 命中高风险词的内容可以直接拒绝，也可以强制进入待审核，取决于关键词配置。

### 我的页面

我的页面包含：

- 个人资料摘要。
- 编辑应届生资料。
- 我的帖子。
- 我的收藏。
- 账号设置。
- 退出登录。

如果用户资料未完善，页面应引导用户先完善资料，再发帖或评论。

### 管理后台

管理后台位于 `/admin`，复用同一个前端项目。

后台页面：

- 待审核队列。
- 举报列表。
- 用户管理。
- 关键词管理。

后台操作：

- 通过待审核内容。
- 拒绝待审核内容。
- 隐藏已发布帖子或评论。
- 恢复被隐藏内容。
- 封禁或解封用户。
- 新增、编辑、删除审核关键词。
- 将举报标记为已处理或驳回。

## 审核模型

内容状态：

- `published`：已发布，对普通用户可见。
- `pending`：待管理员审核。
- `hidden`：被管理员隐藏。
- `rejected`：审核拒绝。
- `deleted`：作者或系统软删除。

需要审核的内容类型：

- 帖子。
- 评论。

关键词规则：

- 每个关键词有风险等级。
- 每个关键词有处理动作，例如进入待审核或直接拒绝。
- 关键词命中记录需要进入后台，方便管理员判断。

举报规则：

- 已完善资料的用户可以举报帖子和评论。
- 举报包含目标类型、目标 id、举报原因、补充说明、举报人 id 和处理状态。
- 同一内容被多次举报时，应在后台中更醒目。

MVP 不做机器学习审核。

## 数据模型

核心数据表：

- `users`：账号、密码哈希、昵称、角色、状态、时间戳。
- `student_profiles`：用户 id、学校、学历、专业、毕业年份、学校是否公开、时间戳。
- `posts`：作者 id、标题、正文、状态、计数器、时间戳。
- `tags`：标签名、slug、展示顺序。
- `post_tags`：帖子和标签的关联。
- `comments`：帖子 id、作者 id、父评论 id、正文、状态、时间戳。
- `likes`：用户 id、目标类型、目标 id、时间戳。
- `favorites`：用户 id、帖子 id、时间戳。
- `reports`：举报人 id、目标类型、目标 id、原因、说明、状态、时间戳。
- `moderation_keywords`：关键词、风险等级、处理动作、是否启用、时间戳。
- `moderation_queue`：目标类型、目标 id、命中的关键词、状态、审核人 id、审核结论、时间戳。
- `sessions`：用户 id、会话 token 哈希、过期时间、时间戳。

需要重点支持的索引：

- 信息流按创建时间和热度排序。
- 按状态和作者查询帖子。
- 按帖子 id 和父评论 id 查询评论。
- 点赞和收藏需要保证同一用户对同一目标唯一。
- 标签筛选。
- 按状态查询举报。
- 按 token 哈希查询会话。

## 技术架构

项目采用单仓库，前后端分目录。

```text
GradLink/
  backend/
  frontend/
  docs/
  scripts/
```

前端：

- Vue 3。
- Vite。
- TypeScript。
- 移动端优先的响应式布局。
- 手机端底部导航。
- 桌面端顶部导航或更宽的信息流布局。

后端：

- Go。
- Gin 作为 HTTP 框架。
- GORM 负责数据库访问和迁移。
- SQLite 作为 MVP 数据库。
- REST API。
- 基于 Cookie 的服务端 Session。

代码层尽量避免不必要的 SQLite 专属 SQL，方便后续迁移到 PostgreSQL。

本地运行形态：

```text
浏览器 -> Vue 前端
Vue 前端 -> /api REST 请求
Go API -> SQLite
```

未来单服务器部署形态：

```text
Nginx -> 前端 dist 静态文件
Nginx /api -> Go API 服务
Go API -> SQLite 数据库文件
```

Node.js 只用于前端开发和构建，不作为生产环境 Web 服务。生产环境可以由 Nginx 直接托管前端构建后的静态文件。

## 后端模块边界

建议后端目录结构：

```text
backend/
  cmd/server/
  internal/auth/
  internal/users/
  internal/posts/
  internal/comments/
  internal/reactions/
  internal/moderation/
  internal/admin/
  internal/db/
```

模块职责：

- `auth`：注册、登录、退出、Session、密码哈希、认证中间件。
- `users`：用户记录和应届生资料。
- `posts`：信息流、帖子详情、发帖、编辑、删除、帖子状态。
- `comments`：评论和一层回复。
- `reactions`：点赞和收藏。
- `moderation`：关键词检查、举报、审核队列。
- `admin`：管理员接口和权限校验。
- `db`：数据库连接、迁移、种子数据。

## API 边界

代表性 API：

- `POST /api/auth/register`
- `POST /api/auth/login`
- `POST /api/auth/logout`
- `GET /api/me`
- `PUT /api/me/profile`
- `GET /api/posts`
- `POST /api/posts`
- `GET /api/posts/:id`
- `POST /api/posts/:id/comments`
- `POST /api/comments/:id/replies`
- `POST /api/posts/:id/like`
- `DELETE /api/posts/:id/like`
- `POST /api/posts/:id/favorite`
- `DELETE /api/posts/:id/favorite`
- `GET /api/tags`
- `POST /api/reports`
- `GET /api/admin/moderation`
- `POST /api/admin/moderation/:id/approve`
- `POST /api/admin/moderation/:id/reject`
- `GET /api/admin/reports`
- `POST /api/admin/reports/:id/resolve`
- `GET /api/admin/users`
- `POST /api/admin/users/:id/ban`
- `POST /api/admin/users/:id/unban`
- `GET /api/admin/keywords`
- `POST /api/admin/keywords`
- `PUT /api/admin/keywords/:id`
- `DELETE /api/admin/keywords/:id`

## 错误处理

API 返回统一 JSON 错误格式：

```json
{
  "error": {
    "code": "PROFILE_REQUIRED",
    "message": "Please complete your student profile before posting."
  }
}
```

重要错误码：

- `UNAUTHENTICATED`
- `FORBIDDEN`
- `PROFILE_REQUIRED`
- `VALIDATION_FAILED`
- `CONTENT_PENDING_REVIEW`
- `CONTENT_REJECTED`
- `NOT_FOUND`
- `RATE_LIMITED`

前端需要展示简洁、友好的错误提示。内容提交进入待审核时，应保留表单信息并明确告诉用户当前状态。

## 测试策略

后端测试重点：

- 注册和登录。
- 密码哈希。
- Session 认证中间件。
- 资料完善门槛。
- 发帖和关键词审核。
- 评论和回复。
- 点赞和收藏的幂等性。
- 举报创建。
- 管理员审核操作。

前端测试可以先轻量：

- API 客户端行为。
- 登录状态处理。
- 资料门槛 UI。
- 发帖表单校验。

MVP 必须进行手动端到端检查：

1. 注册并登录。
2. 未完善资料时尝试发帖。
3. 完善应届生资料。
4. 发布一篇正常帖子。
5. 发布一篇命中关键词的帖子，确认进入待审核。
6. 评论和回复。
7. 点赞和收藏。
8. 举报内容。
9. 以管理员身份处理审核和举报。

## 本地开发

默认端口：

- 后端：`localhost:8080`。
- 前端：`localhost:5173`。

开发流程：

- 本地运行 Go API。
- 本地运行 Vite 前端。
- Vite 将 `/api` 代理到 Go 后端。
- SQLite 数据文件放在后端本地数据目录中，并加入 git 忽略。

第一版实现时应提供种子数据，包括初始标签、管理员账号、审核关键词和几篇示例帖子，方便本地测试。

## 实现默认选择

这些默认选择用于减少后续实现计划中的摇摆，同时保持 MVP 足够轻。

- 密码哈希：使用 `golang.org/x/crypto/bcrypt`。
- Session：自定义数据库 Session 表；Cookie 保存随机 token，SQLite 保存 token 哈希；Cookie 设置为 HTTP-only。
- Gin 中间件：请求日志、异常恢复、认证、管理员权限校验、简单限流。
- 数据库迁移：MVP 使用 GORM auto-migration，并加入标签、关键词、管理员账号和示例帖子的种子逻辑。
- 前端 UI：使用轻量自定义 Vue 组件，不引入大型 UI 框架。
- 视觉方向：内容优先、移动端优先、冷静实用，像一个专注的学生求职社区，而不是营销官网。
