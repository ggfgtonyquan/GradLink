# Repository Guidelines

## 项目结构与模块组织

本仓库 `main` 分支包含 GradLink 的文档、后端和前端实现。

- `docs/`：产品规格、实施计划、低保真和参考原型。
- `backend/`：Go 后端，使用 Gin、GORM、SQLite。
- `backend/internal/`：后端领域模块，如 `auth`、`posts`、`comments`、`users`、`moderation`、`admin`。
- `backend/cmd/server/`：本地 API 服务入口。
- `frontend/`：Vue 3 + Vite + TypeScript 前端。
- `frontend/src/`：应用壳、路由、页面、状态、API 客户端和共享组件。
- 测试与代码就近放置：Go 使用 `*_test.go`，前端使用 `__tests__/*.test.ts`。

## 构建、测试与本地开发命令

除非特别说明，命令都在仓库根目录中运行。

- 后端开发服务：`cd backend; go run ./cmd/server`
- 后端测试：`cd backend; go test ./... -count=1`
- 安装前端依赖：`cd frontend; npm install`
- 前端开发服务：`cd frontend; npm run dev`
- 前端测试：`cd frontend; npm test -- --run`
- 前端生产构建检查：`cd frontend; npm run build`

本地服务可通过 `GRADLINK_ADDR`、`GRADLINK_DB_PATH`、`GRADLINK_ADMIN_PASSWORD` 配置端口、数据库路径和管理员密码。

## 代码风格与命名约定

Go 文件必须使用 `gofmt` 格式化；包保持小而清晰，按领域组织。已有的 handler/service/repository 分层应继续沿用。前端使用 Vue 单文件组件、TypeScript、Pinia 和 Vue Router。组件名使用 PascalCase，变量和函数使用 camelCase，路由路径使用 kebab-case。面向用户的 MVP 页面文案保持中文。

## 测试规范

修复 bug 或新增行为时应补回归测试。后端优先通过 `internal/testutil` 覆盖 HTTP handler 和 service 行为。前端使用 Vitest 与 Testing Library；涉及 store 或路由的组件测试需要挂载 Pinia/router 插件。交付前至少运行后端测试和前端测试。

## 提交与 PR 规范

当前提交历史采用 Conventional Commit 风格，例如 `feat: add frontend shell and views`、`fix: clamp reaction counters on stale aggregates`。提交应保持范围清晰、动词明确。PR 需要包含简短变更说明、已运行的验证命令、关联问题或背景；涉及 UI 的变更应附截图或录屏。

## Agent 专用说明

开始工作前先确认当前分支和工作区状态。若浏览器验证在当前环境不稳定，优先使用测试、构建和命令行 API 检查。不要覆盖或回退与当前任务无关的用户改动。
