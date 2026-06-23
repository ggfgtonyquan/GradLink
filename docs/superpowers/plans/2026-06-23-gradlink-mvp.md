# GradLink MVP Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Build a mobile-first web community for domestic new graduates with email login, a student-profile gate, posts, comments, likes, favorites, moderation, and a small admin console.

**Architecture:** Single repo, split into `backend` and `frontend`. The backend is a Go + Gin API with SQLite, GORM, bcrypt, and cookie sessions. The frontend is a Vue 3 + Vite + TypeScript app with custom UI components and Vue Router. Implement the product as vertical slices so each milestone can be tested end-to-end before the next one starts.

**Tech Stack:** Go 1.24+, Gin, GORM, SQLite, bcrypt, Vue 3, Vite, TypeScript, Vue Router, Pinia, Vitest, Vue Testing Library, Playwright/manual browser checks.

---

## File Structure Before Coding

Plan the repo around small, focused files.

- Root: `.gitignore`, `README.md`
- Backend: `backend/go.mod`, `backend/cmd/server/main.go`, `backend/internal/http/router.go`, `backend/internal/db/*`, `backend/internal/auth/*`, `backend/internal/users/*`, `backend/internal/posts/*`, `backend/internal/comments/*`, `backend/internal/reactions/*`, `backend/internal/moderation/*`, `backend/internal/admin/*`, `backend/internal/testutil/*`
- Frontend: `frontend/package.json`, `frontend/index.html`, `frontend/vite.config.ts`, `frontend/tsconfig.json`, `frontend/src/main.ts`, `frontend/src/App.vue`, `frontend/src/router.ts`, `frontend/src/api/*`, `frontend/src/stores/*`, `frontend/src/views/*`, `frontend/src/views/__tests__/*`, `frontend/src/components/*`, `frontend/src/components/__tests__/*`
- Docs: `docs/run-local.md`

Keep each file responsible for one slice of behavior. Avoid giant catch-all files.

Backend integration tests should share helpers from `backend/internal/testutil/api.go`:

- `NewAPIServer(t *testing.T) *APIServer`
- `(*APIServer).NewAgent() *Agent`
- `(*Agent).PostJSON(t *testing.T, path string, body any) *httptest.ResponseRecorder`
- `(*Agent).Get(t *testing.T, path string) *httptest.ResponseRecorder`
- `RegisterAndLogin(t *testing.T, s *APIServer, email string) *Agent`
- `CompleteProfile(t *testing.T, a *Agent)`
- `CreatePublishedPost(t *testing.T, a *Agent) int64`
- `AssertJSONError(t *testing.T, rec *httptest.ResponseRecorder, status int, code string)`

---

### Task 1: Scaffold the repo and a visible health check

**Files:**
- Create: `.gitignore`
- Create: `README.md`
- Create: `backend/go.mod`
- Create: `backend/cmd/server/main.go`
- Create: `backend/internal/http/router.go`
- Create: `backend/internal/http/health_test.go`
- Create: `frontend/package.json`
- Create: `frontend/index.html`
- Create: `frontend/vite.config.ts`
- Create: `frontend/tsconfig.json`
- Create: `frontend/src/main.ts`
- Create: `frontend/src/App.vue`
- Create: `frontend/src/components/AppShell.vue`
- Create: `frontend/src/components/__tests__/AppShell.test.ts`

- [ ] **Step 1: Write the failing backend health test**

```go
package http_test

import (
	stdhttp "net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	apphttp "gradlink/internal/http"
)

func TestHealthEndpoint(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	apphttp.RegisterRoutes(r, nil)

	req := httptest.NewRequest(stdhttp.MethodGet, "/api/health", nil)
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)

	if rec.Code != stdhttp.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}
	if rec.Body.String() != `{"ok":true}` {
		t.Fatalf("unexpected body: %s", rec.Body.String())
	}
}
```

- [ ] **Step 2: Run the backend test and confirm it fails**

Run:

```bash
cd backend
go test ./internal/http -run TestHealthEndpoint -v
```

Expected: FAIL because the HTTP router and health handler do not exist yet.

- [ ] **Step 3: Write the minimal backend implementation**

```go
package http

import "github.com/gin-gonic/gin"

func RegisterRoutes(r *gin.Engine, deps any) {
	r.GET("/api/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"ok": true})
	})
}
```

Create `backend/cmd/server/main.go` so the server starts with Gin, loads the router, and listens on `:8080`.

- [ ] **Step 4: Write the failing frontend shell test**

```ts
import { render, screen } from '@testing-library/vue'
import AppShell from '../AppShell.vue'

test('renders bottom navigation', () => {
  render(AppShell)
  expect(screen.getByText('首页')).toBeTruthy()
  expect(screen.getByText('发帖')).toBeTruthy()
  expect(screen.getByText('收藏')).toBeTruthy()
  expect(screen.getByText('我的')).toBeTruthy()
})
```

- [ ] **Step 5: Run the frontend test and confirm it fails**

Run:

```bash
cd frontend
npm test -- --run src/components/__tests__/AppShell.test.ts
```

Expected: FAIL because the Vue app shell and test setup do not exist yet.

- [ ] **Step 6: Write the minimal frontend shell**

Implement `AppShell.vue` with a top bar, content slot, and bottom navigation. Render it from `App.vue` so the first screen already looks like the approved mobile reference.

- [ ] **Step 7: Run both checks again**

Run:

```bash
cd backend
go test ./internal/http -run TestHealthEndpoint -v

cd ../frontend
npm test -- --run src/components/__tests__/AppShell.test.ts
npm run build
```

Expected: all pass.

- [ ] **Step 8: Commit**

```bash
git add .
git commit -m "chore: scaffold gradlink shell"
```

---

### Task 2: Add SQLite, models, migrations, and seed data

**Files:**
- Create: `backend/internal/db/db.go`
- Create: `backend/internal/db/models.go`
- Create: `backend/internal/db/migrate.go`
- Create: `backend/internal/db/seed.go`
- Create: `backend/internal/db/migrate_test.go`
- Create: `backend/internal/db/seed_test.go`
- Modify: `backend/cmd/server/main.go`
- Modify: `backend/internal/http/router.go`

- [ ] **Step 1: Write the failing schema test**

```go
func TestRequiredTablesExist(t *testing.T) {
	db := testDB(t)
	want := []string{
		"users",
		"student_profiles",
		"posts",
		"tags",
		"post_tags",
		"comments",
		"likes",
		"favorites",
		"reports",
		"moderation_keywords",
		"moderation_queue",
		"sessions",
	}
	for _, table := range want {
		if !tableExists(t, db, table) {
			t.Fatalf("missing table %s", table)
		}
	}
}
```

- [ ] **Step 2: Run the schema test and confirm failure**

Run:

```bash
cd backend
go test ./internal/db -run TestRequiredTablesExist -v
```

Expected: FAIL because migrations and the DB helper are missing.

- [ ] **Step 3: Implement the database layer**

Create:

- `Open(path string) (*gorm.DB, error)` in `db.go`
- `Migrate(db *gorm.DB) error` in `migrate.go`
- `Seed(db *gorm.DB) error` in `seed.go`

Define the core models in `models.go` with the exact fields from the approved spec. Keep `published`, `pending`, `hidden`, `rejected`, and `deleted` as string statuses.

- [ ] **Step 4: Add a seed test for tags and admin account**

```go
func TestSeedCreatesTagsAndAdmin(t *testing.T) {
	db := testDB(t)
	if countRows(t, db, "tags") < 8 {
		t.Fatal("expected seed tags")
	}
	if !adminUserExists(t, db) {
		t.Fatal("expected admin user seed")
	}
}
```

- [ ] **Step 5: Run database tests**

Run:

```bash
cd backend
go test ./internal/db -v
```

Expected: pass.

- [ ] **Step 6: Wire the DB into server startup**

`main.go` should open `data/gradlink.db`, run migrations, run seed data, and inject the DB handle into route dependencies.

- [ ] **Step 7: Commit**

```bash
git add .
git commit -m "feat: add sqlite schema and seeds"
```

---

### Task 3: Build email/password auth, sessions, and student profile gating

**Files:**
- Create: `backend/internal/auth/handler.go`
- Create: `backend/internal/auth/service.go`
- Create: `backend/internal/auth/password.go`
- Create: `backend/internal/auth/session_store.go`
- Create: `backend/internal/auth/middleware.go`
- Create: `backend/internal/auth/auth_test.go`
- Create: `backend/internal/users/profile_handler.go`
- Create: `backend/internal/users/profile_test.go`
- Create: `backend/internal/testutil/api.go`
- Modify: `backend/internal/http/router.go`
- Modify: `backend/internal/db/models.go`

- [ ] **Step 1: Write the failing auth flow test**

```go
func TestRegisterLoginAndProfileGate(t *testing.T) {
	s := testutil.NewAPIServer(t)
	a := testutil.RegisterAndLogin(t, s, "student@example.com")

	rec := a.PostJSON(t, "/api/posts", map[string]any{
		"title":  "未完善资料时不能发帖",
		"body":   "这条请求应该被资料门槛拦截。",
		"tagIds": []int64{1},
	})
	testutil.AssertJSONError(t, rec, 403, "PROFILE_REQUIRED")

	testutil.CompleteProfile(t, a)

	me := a.Get(t, "/api/me")
	if me.Code != 200 {
		t.Fatalf("expected /api/me 200, got %d", me.Code)
	}
	if !strings.Contains(me.Body.String(), `"profileComplete":true`) {
		t.Fatalf("expected completed profile response, got %s", me.Body.String())
	}
}
```

- [ ] **Step 2: Run the auth test and confirm it fails**

Run:

```bash
cd backend
go test ./internal/auth -run TestRegisterLoginAndProfileGate -v
```

Expected: FAIL because the auth service, session handling, and profile gate are not implemented yet.

- [ ] **Step 3: Implement registration, login, logout, and `GET /api/me`**

Use `bcrypt` for password hashing, store sessions in SQLite, and set the session token in an HTTP-only cookie. The cookie should identify the session; the database should store only the token hash.

- [ ] **Step 4: Implement profile completion and the gating middleware**

`PUT /api/me/profile` should save:

- nickname
- school
- education level
- major or direction
- graduation year
- school visibility

The posting/commenting middleware should reject requests from logged-in users whose profile is incomplete with:

```json
{"error":{"code":"PROFILE_REQUIRED","message":"Please complete your student profile before posting."}}
```

- [ ] **Step 5: Run auth and profile tests**

Run:

```bash
cd backend
go test ./internal/auth ./internal/users -v
```

Expected: pass.

- [ ] **Step 6: Commit**

```bash
git add .
git commit -m "feat: add auth and profile gate"
```

---

### Task 4: Implement posts, tags, feed sorting, and detail pages on the API

**Files:**
- Create: `backend/internal/posts/handler.go`
- Create: `backend/internal/posts/service.go`
- Create: `backend/internal/posts/repository.go`
- Create: `backend/internal/posts/posts_test.go`
- Create: `backend/internal/posts/feed_test.go`
- Create: `backend/internal/posts/detail_test.go`
- Create: `backend/internal/tags/seed.go`
- Modify: `backend/internal/http/router.go`
- Modify: `backend/internal/db/models.go`

- [ ] **Step 1: Write the failing posts API test**

```go
func TestCreateAndListPosts(t *testing.T) {
	s := testutil.NewAPIServer(t)
	a := testutil.RegisterAndLogin(t, s, "poster@example.com")
	testutil.CompleteProfile(t, a)

	create := a.PostJSON(t, "/api/posts", map[string]any{
		"title":  "秋招面试被问到项目细节，怎么准备更稳？",
		"body":   "最近投后端开发，面试官一直追问项目里的并发和数据库。",
		"tagIds": []int64{1, 3},
	})
	if create.Code != 201 {
		t.Fatalf("expected create 201, got %d: %s", create.Code, create.Body.String())
	}

	feed := a.Get(t, "/api/posts?sort=latest")
	if !strings.Contains(feed.Body.String(), "秋招面试被问到项目细节") {
		t.Fatalf("expected created post in feed, got %s", feed.Body.String())
	}

	detail := a.Get(t, "/api/posts/1")
	if detail.Code != 200 {
		t.Fatalf("expected detail 200, got %d", detail.Code)
	}
}
```

- [ ] **Step 2: Run the posts test and confirm it fails**

Run:

```bash
cd backend
go test ./internal/posts -run TestCreateAndListPosts -v
```

Expected: FAIL because the posts handlers and feed logic do not exist yet.

- [ ] **Step 3: Implement post creation, feed listing, and detail lookup**

Support:

- title
- body
- up to 3 tags
- latest sort
- hot sort based on a simple score derived from comments and likes
- post status

Keep the initial search simple: title/body filter in SQL, no separate search engine.

- [ ] **Step 4: Add tag list support**

Expose `GET /api/tags` and seed the initial tags from the approved spec:

- 面试经验
- Offer 比较
- 秋招进度
- 春招进度
- 简历求助
- 笔试讨论
- 签约三方
- 吐槽避坑
- 综合提问

- [ ] **Step 5: Run posts tests**

Run:

```bash
cd backend
go test ./internal/posts ./internal/tags -v
```

Expected: pass.

- [ ] **Step 6: Commit**

```bash
git add .
git commit -m "feat: add post feed and tags"
```

---

### Task 5: Add comments, replies, likes, favorites, and reports

**Files:**
- Create: `backend/internal/comments/handler.go`
- Create: `backend/internal/comments/service.go`
- Create: `backend/internal/comments/comments_test.go`
- Create: `backend/internal/reactions/handler.go`
- Create: `backend/internal/reactions/reactions_test.go`
- Create: `backend/internal/moderation/report_handler.go`
- Create: `backend/internal/moderation/report_test.go`
- Modify: `backend/internal/http/router.go`
- Modify: `backend/internal/db/models.go`

- [ ] **Step 1: Write the failing interaction test**

```go
func TestCommentLikeFavoriteAndReport(t *testing.T) {
	s := testutil.NewAPIServer(t)
	a := testutil.RegisterAndLogin(t, s, "reader@example.com")
	testutil.CompleteProfile(t, a)
	postID := testutil.CreatePublishedPost(t, a)

	comment := a.PostJSON(t, fmt.Sprintf("/api/posts/%d/comments", postID), map[string]any{
		"body": "可以按背景、难点、方案、结果四段准备。",
	})
	if comment.Code != 201 {
		t.Fatalf("expected comment 201, got %d", comment.Code)
	}

	reply := a.PostJSON(t, "/api/comments/1/replies", map[string]any{
		"body": "这个结构很清楚，我试一下。",
	})
	if reply.Code != 201 {
		t.Fatalf("expected reply 201, got %d", reply.Code)
	}

	if rec := a.PostJSON(t, fmt.Sprintf("/api/posts/%d/like", postID), nil); rec.Code != 200 {
		t.Fatalf("expected like 200, got %d", rec.Code)
	}
	if rec := a.PostJSON(t, fmt.Sprintf("/api/posts/%d/favorite", postID), nil); rec.Code != 200 {
		t.Fatalf("expected favorite 200, got %d", rec.Code)
	}
	report := a.PostJSON(t, "/api/reports", map[string]any{
		"targetType": "post",
		"targetId":   postID,
		"reason":     "广告引流",
		"detail":     "这条内容疑似引导加微信。",
	})
	if report.Code != 201 {
		t.Fatalf("expected report 201, got %d", report.Code)
	}
}
```

- [ ] **Step 2: Run the interaction tests and confirm failure**

Run:

```bash
cd backend
go test ./internal/comments ./internal/reactions ./internal/moderation -run TestCommentLikeFavoriteAndReport -v
```

Expected: FAIL because the handlers are missing.

- [ ] **Step 3: Implement one-level reply comments**

Store `parent_comment_id` and keep the nesting depth to one reply level in the API.

- [ ] **Step 4: Implement like and favorite idempotency**

The same user should not create duplicate likes or favorites for the same target. Repeated requests should return success without duplicating rows.

- [ ] **Step 5: Implement the report endpoint**

Allow reporting both posts and comments with:

- target type
- target id
- reason
- detail

Reports should be visible to admins in the moderation queue.

- [ ] **Step 6: Run interaction tests**

Run:

```bash
cd backend
go test ./internal/comments ./internal/reactions ./internal/moderation -v
```

Expected: pass.

- [ ] **Step 7: Commit**

```bash
git add .
git commit -m "feat: add comments reactions and reports"
```

---

### Task 6: Build keyword moderation and the admin console API

**Files:**
- Create: `backend/internal/moderation/keywords_handler.go`
- Create: `backend/internal/moderation/queue.go`
- Create: `backend/internal/moderation/moderation_test.go`
- Create: `backend/internal/admin/handler.go`
- Create: `backend/internal/admin/admin_test.go`
- Modify: `backend/internal/http/router.go`
- Modify: `backend/internal/db/models.go`

- [ ] **Step 1: Write the failing moderation test**

```go
func TestKeywordReviewAndAdminActions(t *testing.T) {
	s := testutil.NewAPIServer(t)
	user := testutil.RegisterAndLogin(t, s, "risk@example.com")
	testutil.CompleteProfile(t, user)

	pending := user.PostJSON(t, "/api/posts", map[string]any{
		"title":  "有人说可以加微信内推，靠谱吗？",
		"body":   "对方要求加微信并购买培训服务，我不确定是不是广告。",
		"tagIds": []int64{9},
	})
	if pending.Code != 202 {
		t.Fatalf("expected pending review 202, got %d: %s", pending.Code, pending.Body.String())
	}

	admin := testutil.LoginAsSeedAdmin(t, s)
	queue := admin.Get(t, "/api/admin/moderation")
	if !strings.Contains(queue.Body.String(), "加微信") {
		t.Fatalf("expected keyword hit in queue, got %s", queue.Body.String())
	}

	approve := admin.PostJSON(t, "/api/admin/moderation/1/approve", nil)
	if approve.Code != 200 {
		t.Fatalf("expected approve 200, got %d", approve.Code)
	}
}
```

- [ ] **Step 2: Run moderation tests and confirm failure**

Run:

```bash
cd backend
go test ./internal/moderation ./internal/admin -run TestKeywordReviewAndAdminActions -v
```

Expected: FAIL because keyword review and admin actions do not exist yet.

- [ ] **Step 3: Implement keyword checks and moderation queue writes**

Use the seed keyword table to determine whether content is:

- published immediately
- sent to pending review
- rejected

Persist keyword hits to the moderation queue so admins can inspect why a post was flagged.

- [ ] **Step 4: Implement admin-only endpoints**

Admin endpoints should support:

- listing pending content
- approving or rejecting content
- listing reports
- resolving or rejecting reports
- listing users
- banning or unbanning users
- listing keywords
- creating, updating, and deleting keywords

- [ ] **Step 5: Run moderation/admin tests**

Run:

```bash
cd backend
go test ./internal/moderation ./internal/admin -v
```

Expected: pass.

- [ ] **Step 6: Commit**

```bash
git add .
git commit -m "feat: add moderation and admin api"
```

---

### Task 7: Build the Vue app shell and the main mobile views

**Files:**
- Create: `frontend/src/router.ts`
- Create: `frontend/src/api/http.ts`
- Create: `frontend/src/api/types.ts`
- Create: `frontend/src/stores/auth.ts`
- Create: `frontend/src/stores/feed.ts`
- Create: `frontend/src/views/HomeFeedView.vue`
- Create: `frontend/src/views/PostDetailView.vue`
- Create: `frontend/src/views/ComposeView.vue`
- Create: `frontend/src/views/ProfileView.vue`
- Create: `frontend/src/views/LoginView.vue`
- Create: `frontend/src/views/RegisterView.vue`
- Create: `frontend/src/views/ModerationView.vue`
- Create: `frontend/src/components/PostCard.vue`
- Create: `frontend/src/components/TagChips.vue`
- Create: `frontend/src/components/BottomNav.vue`
- Create: `frontend/src/views/__tests__/HomeFeedView.test.ts`
- Create: `frontend/src/views/__tests__/ComposeView.test.ts`
- Modify: `frontend/src/App.vue`
- Modify: `frontend/src/main.ts`

- [ ] **Step 1: Write the failing front-end view test**

```ts
import { render, screen } from '@testing-library/vue'
import HomeFeedView from '../HomeFeedView.vue'

test('shows mobile feed cards and filter chips', () => {
  render(HomeFeedView)
  expect(screen.getByText('最新')).toBeTruthy()
  expect(screen.getByText('热门')).toBeTruthy()
  expect(screen.getByText('面试经验')).toBeTruthy()
  expect(screen.getByText('秋招面试被问到项目细节，怎么准备更稳？')).toBeTruthy()
})
```

- [ ] **Step 2: Run the front-end test and confirm failure**

Run:

```bash
cd frontend
npm test -- --run src/views/__tests__/HomeFeedView.test.ts
```

Expected: FAIL because the routed views and feed components do not exist yet.

- [ ] **Step 3: Implement the app shell and router**

Create a mobile-first shell with:

- top bar
- main content
- bottom navigation

Wire routes for:

- `/`
- `/posts/:id`
- `/compose`
- `/login`
- `/register`
- `/profile`
- `/moderation`

- [ ] **Step 4: Implement the feed, detail, compose, and profile views**

Match the approved low-fidelity layout:

- feed first on home
- detail page centered on reading and reply
- compose page with title/body/tags
- profile page with student-info completion

- [ ] **Step 5: Add the auth store and API wrapper**

The API client should:

- send cookies
- normalize JSON errors
- expose typed methods for auth, posts, comments, reactions, reports, and moderation

The auth store should track:

- current user
- profile completion state
- role

- [ ] **Step 6: Run front-end tests and a production build**

Run:

```bash
cd frontend
npm test -- --run
npm run build
```

Expected: pass.

- [ ] **Step 7: Commit**

```bash
git add .
git commit -m "feat: add frontend shell and views"
```

---

### Task 8: Wire the end-to-end flow, local docs, and manual verification

**Files:**
- Create: `docs/run-local.md`
- Modify: `README.md`
- Modify: `backend/cmd/server/main.go`
- Modify: `frontend/vite.config.ts`

- [ ] **Step 1: Write the end-to-end verification checklist**

```md
1. Start backend and frontend locally.
2. Register with email and password.
3. Log in.
4. Try to post before completing the student profile.
5. Complete the student profile.
6. Create a clean post.
7. Create a keyword-hit post and confirm pending review.
8. Comment, reply, like, favorite, and report.
9. Log in as admin and process moderation items.
```

- [ ] **Step 2: Implement the exact dev commands in docs**

`docs/run-local.md` should document:

- `cd backend && go run ./cmd/server`
- `cd frontend && npm install && npm run dev`
- the local URLs for both apps
- where the SQLite file lives
- how to reset local data safely

- [ ] **Step 3: Verify the full journey manually**

Use the browser to confirm:

- the home feed renders posts immediately
- the profile gate blocks posting until completed
- posts can move into pending review
- replies and reactions persist
- admin actions change visible state

- [ ] **Step 4: Clean up rough edges**

Fix any wording, spacing, overflow, or response-code mismatch discovered during manual testing.

- [ ] **Step 5: Commit**

```bash
git add .
git commit -m "feat: complete gradlink mvp flow"
```

---

## Coverage Check

This plan covers the approved spec and UI notes:

- Mobile-first website instead of a native app: Tasks 1 and 7.
- Email/password login with no email verification: Task 3.
- Student-profile gate before posting and commenting: Task 3.
- Posts, comments, replies, likes, favorites, and reports: Tasks 4 and 5.
- Keyword moderation, pending review, and admin console: Task 6.
- SQLite and one-server deployment path: Tasks 2 and 8.
- Approved UI structure and content-first layout: Tasks 1 and 7.

## Self-Review Notes

- No unresolved filler terms remain.
- The naming is consistent across backend and frontend tasks: `auth`, `posts`, `comments`, `reactions`, `moderation`, and `admin`.
- The plan starts with a visible vertical slice, not a giant all-at-once build.
- Each task is independently testable and ends with a commit checkpoint.
