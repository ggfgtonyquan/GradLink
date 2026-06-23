# GradLink Design Spec

Date: 2026-06-22

## Summary

GradLink is a mobile-first web community for domestic Chinese new graduates to discuss job hunting. The first version focuses on a small but complete community loop: browse posts, register, complete a student profile, publish posts, comment, reply, like, favorite, report, and moderate content.

The product is a responsive website, not a native mobile app. It should work well on phones and remain usable on desktop. A future version may add PWA support, a mini program, or native apps, but the MVP avoids app store release and mobile app compliance overhead.

## Goals

- Help domestic new graduates exchange job-search experiences in a semi-anonymous setting.
- Make posting and discussion easy while requiring a minimal student profile before contributing.
- Keep the first version cheap to run, easy to develop with AI assistance, and simple to deploy on one server later.
- Include basic moderation from day one: keyword checks, pending review, reports, hiding content, and user bans.
- Keep the system small enough for a personal project.

## Non-Goals

- Job aggregation.
- Application tracking.
- Company profile pages.
- Private messaging.
- Real-time notifications.
- Recommendation algorithms.
- Image upload.
- Native app or app store release.
- Strong student verification, OCR, Xuexin integration, or mandatory campus email verification.

## Users and Permissions

GradLink uses a semi-anonymous identity model.

User states:

- Guest: can browse posts and comments.
- Registered user: can log in, like, favorite, and complete a student profile.
- Profile-complete user: can publish posts, comment, reply, and report content.
- Admin: can review pending content, handle reports, hide posts/comments, ban users, and manage moderation keywords.

The contribution gate is profile completion, not verified email.

Student profile fields:

- Nickname.
- School.
- Education level.
- Major or direction.
- Graduation year.
- Whether to publicly show the school.

Public identity display should remain semi-anonymous, for example "2026 master student · computer science" or "Shanghai university student · 2026", depending on the user's visibility choice.

## Authentication

The MVP uses email and password accounts.

- Email is used only for registration and login.
- The system checks email format and uniqueness.
- The MVP does not send verification emails.
- The MVP does not require bidirectional email verification.
- Posting and commenting require login plus a completed student profile.
- Passwords must be hashed before storage.
- The first version can omit self-service password reset. If needed, admins can reset accounts manually during early testing.

Anti-abuse controls replace email verification in the MVP:

- Basic rate limits for registration, login, posting, and commenting.
- Keyword-based moderation.
- User reports.
- Admin bans.
- Unique email constraint.

Email verification, password reset, phone verification, and campus email certification can be added later if the site moves toward public launch.

## Product Surface

### Home Feed

The home feed is the main screen. It lists posts with title, short excerpt, tags, author display, comment count, like count, favorite state, and publish time.

Feed controls:

- Latest sort.
- Hot sort.
- Tag filtering.
- Keyword search over title and body.

Initial tags:

- Interview experience.
- Offer comparison.
- Autumn recruitment.
- Spring recruitment.
- Resume help.
- Written test discussion.
- Signing and tripartite agreement.
- Workplace warnings.
- General question.

### Post Detail

The post detail page shows:

- Title.
- Body.
- Tags.
- Author semi-anonymous profile.
- Status if pending/hidden for users who can see it.
- Like action.
- Favorite action.
- Report action.
- Comment list.
- Reply action under comments.

Comments support one reply level through a parent comment id. The MVP does not need deeply nested threads.

### Create Post

The create post page includes:

- Title input.
- Body editor.
- Tag picker.
- Submit action.

On submit, the backend runs keyword checks:

- Clean content is published immediately.
- Risky content is saved as pending and shown as "waiting for review".
- High-risk content can be rejected or forced into pending review, depending on keyword configuration.

### My Page

The user page includes:

- Profile summary.
- Edit student profile.
- My posts.
- My favorites.
- Account settings.
- Logout.

If the student profile is incomplete, the UI should guide the user to complete it before posting or commenting.

### Admin Area

The admin area lives under `/admin` in the same frontend project.

Admin pages:

- Pending moderation queue.
- Reports.
- Users.
- Keywords.

Admin actions:

- Approve pending content.
- Reject pending content.
- Hide published posts or comments.
- Restore hidden content.
- Ban or unban users.
- Add, edit, or remove moderation keywords.
- Mark reports as resolved or rejected.

## Moderation Model

Content statuses:

- `published`: visible to normal users.
- `pending`: waiting for admin review.
- `hidden`: hidden by moderation.
- `rejected`: rejected during review.
- `deleted`: soft-deleted by the author or system.

Moderated content types:

- Posts.
- Comments.

Keyword rules:

- Each keyword has a risk level.
- Each keyword has an action, such as pending review or reject.
- Keyword hits should be recorded for admin review.

Report rules:

- Profile-complete users can report posts and comments.
- Reports include target type, target id, reason, optional detail, reporter id, and status.
- Multiple reports on the same content should make that content more visible in the admin queue.

The MVP does not need automated machine-learning moderation.

## Data Model

Core tables:

- `users`: account, password hash, nickname, role, status, timestamps.
- `student_profiles`: user id, school, education level, major, graduation year, school visibility, timestamps.
- `posts`: author id, title, body, status, counters, timestamps.
- `tags`: name, slug, display order.
- `post_tags`: post id, tag id.
- `comments`: post id, author id, parent comment id, body, status, timestamps.
- `likes`: user id, target type, target id, timestamps.
- `favorites`: user id, post id, timestamps.
- `reports`: reporter id, target type, target id, reason, detail, status, timestamps.
- `moderation_keywords`: keyword, risk level, action, enabled, timestamps.
- `moderation_queue`: target type, target id, matched keywords, status, reviewer id, decision, timestamps.
- `sessions`: user id, session token hash, expiry, timestamps.

Indexes should support:

- Feed ordering by creation time and hot score.
- Post lookup by status and author.
- Comment lookup by post id and parent id.
- Unique likes and favorites per user and target.
- Tag filtering.
- Report lookup by status.
- Session lookup by token hash.

## Technical Architecture

The repository is a single project with separate frontend and backend directories.

```text
GradLink/
  backend/
  frontend/
  docs/
  scripts/
```

Frontend:

- Vue 3.
- Vite.
- TypeScript.
- Mobile-first responsive layout.
- Bottom navigation on mobile.
- Top navigation or wider layout on desktop.

Backend:

- Go.
- Gin as the HTTP framework.
- GORM for database access and migrations.
- SQLite for the MVP database.
- REST API.
- Cookie-based server-side sessions.

The code should avoid unnecessary SQLite-specific SQL so the project can migrate to PostgreSQL later.

Runtime shape:

```text
Browser -> Vue frontend
Vue frontend -> /api REST calls
Go API -> SQLite
```

Future single-server deployment:

```text
Nginx -> frontend dist files
Nginx /api -> Go API service
Go API -> SQLite database file
```

Node.js is required for frontend development and build, but not as a production web server. Production can serve the built frontend files through Nginx.

## Backend Module Boundaries

Suggested backend structure:

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

Responsibilities:

- `auth`: registration, login, logout, sessions, password hashing, auth middleware.
- `users`: user records and student profiles.
- `posts`: feed, post detail, create/update/delete posts, post status.
- `comments`: comments and one-level replies.
- `reactions`: likes and favorites.
- `moderation`: keyword checks, reports, moderation queue.
- `admin`: admin-only endpoints and authorization.
- `db`: connection, migrations, seed data.

## API Boundaries

Representative API endpoints:

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

## Error Handling

The API should return consistent JSON errors:

```json
{
  "error": {
    "code": "PROFILE_REQUIRED",
    "message": "Please complete your student profile before posting."
  }
}
```

Important error codes:

- `UNAUTHENTICATED`
- `FORBIDDEN`
- `PROFILE_REQUIRED`
- `VALIDATION_FAILED`
- `CONTENT_PENDING_REVIEW`
- `CONTENT_REJECTED`
- `NOT_FOUND`
- `RATE_LIMITED`

The frontend should show concise, user-friendly messages and keep form data when submission moves into pending review.

## Testing Strategy

Backend tests should cover:

- Registration and login.
- Password hashing behavior.
- Session auth middleware.
- Profile completion gate.
- Post creation and keyword moderation.
- Comment and reply creation.
- Like and favorite idempotency.
- Report creation.
- Admin moderation actions.

Frontend tests can start lighter:

- API client behavior.
- Auth state handling.
- Profile gate UI.
- Post form validation.

Manual end-to-end checks are required for the MVP:

1. Register and log in.
2. Attempt to post before completing profile.
3. Complete profile.
4. Create a clean post.
5. Create a keyword-hit post and confirm it enters pending review.
6. Comment and reply.
7. Like and favorite.
8. Report content.
9. Resolve moderation and report items as admin.

## Local Development

Default ports:

- Backend: `localhost:8080`.
- Frontend: `localhost:5173`.

Development flow:

- Run Go API locally.
- Run Vite frontend locally.
- Proxy `/api` from Vite to the Go backend.
- Store SQLite data under a local backend data directory ignored by git.

The first implementation should include seed data for tags, an admin user, moderation keywords, and a few sample posts to make local testing useful.

## Implementation Defaults

These defaults remove ambiguity for the implementation plan while keeping the MVP small.

- Password hashing: use `golang.org/x/crypto/bcrypt`.
- Sessions: use a custom database-backed session table with random tokens stored in secure HTTP-only cookies and token hashes stored in SQLite.
- Gin middleware: use request logging, recovery, authentication, admin authorization, and simple rate-limiting middleware.
- Database migration: use GORM auto-migration for the MVP, plus seed logic for tags, keywords, admin account, and sample posts.
- Frontend UI: use lightweight custom Vue components rather than a large UI framework.
- Visual direction: content-first, mobile-first, calm and practical. The interface should feel like a focused student community rather than a marketing site.
