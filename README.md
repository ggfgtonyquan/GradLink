# GradLink

[![CI](https://github.com/ggfgtonyquan/GradLink/actions/workflows/ci.yml/badge.svg)](https://github.com/ggfgtonyquan/GradLink/actions/workflows/ci.yml) [![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](LICENSE)

GradLink 是一个面向国内应届生的移动端优先求职交流社区 MVP，提供登录、应届生资料门槛、内容互动和轻量管理能力。

> **重要：** 本项目适合本地开发、学习和功能验证。公开部署前必须配置 HTTPS、强管理员密码、Secure Cookie、可靠备份、限流，以及完善的内容和隐私治理。

## 功能

- 邮箱密码认证与会话管理；完善应届生资料后才能发布内容。
- 帖子、评论、回复、点赞和收藏。
- 标签、个人资料与活动记录。
- 举报、关键词审核和管理员处理能力。
- Vue 3 移动端优先界面。

## 技术栈

- 后端：Go、Gin、GORM、SQLite。
- 前端：Vue 3、Vite、TypeScript、Pinia、Vue Router。
- 测试：Go testing、Vitest、Testing Library。

## 项目结构

- `backend/`：Go REST API、领域模块和服务入口。
- `frontend/`：Vue 3 前端应用。
- `docs/`：产品、实现和本地运行文档。

## 环境要求

- Go 1.26.6 或更高的兼容版本。
- Node.js 20 或更高版本。
- npm。

## 快速启动

后端（PowerShell）：

```powershell
cd backend
$env:GRADLINK_ADMIN_PASSWORD="<set-a-strong-admin-password>"
go run ./cmd/server
```

前端（新开一个终端）：

```powershell
cd frontend
npm install
npm run dev
```

默认地址：

- 前端：<http://localhost:5173>
- 后端：<http://localhost:8080>
- 健康检查：<http://localhost:8080/api/health>

管理员账号为 `admin@gradlink.local`，仅在设置 `GRADLINK_ADMIN_PASSWORD` 后启用。更完整的启动、重置和验证说明见 [docs/run-local.md](docs/run-local.md)。

## 配置

| 变量 | 默认值 | 说明 |
| --- | --- | --- |
| `GRADLINK_ADDR` | `:8080` | 后端监听地址 |
| `GRADLINK_DB_PATH` | `data/gradlink.db` | SQLite 数据库路径 |
| `GRADLINK_ADMIN_PASSWORD` | 未设置（管理员禁用） | 管理员密码；请使用强密码 |
| `GRADLINK_COOKIE_SECURE` | `false` | 本地 HTTP 开发可保持 `false`；生产 HTTPS 部署必须设为 `true` |

## 验证

```powershell
cd backend
go test ./... -count=1
go run golang.org/x/vuln/cmd/govulncheck@latest ./...

cd ..\frontend
npm test -- --run
npm run build
npm audit --audit-level=high
```

## 贡献与安全

欢迎通过 [CONTRIBUTING.md](CONTRIBUTING.md) 参与贡献，并遵守 [CODE_OF_CONDUCT.md](CODE_OF_CONDUCT.md)。安全问题请按照 [SECURITY.md](SECURITY.md) 私下报告，不要通过公开 Issue 披露。

## 许可证

本项目采用 [MIT License](LICENSE)。
