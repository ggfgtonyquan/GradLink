# GradLink

[![CI](https://github.com/ggfgtonyquan/GradLink/actions/workflows/ci.yml/badge.svg)](https://github.com/ggfgtonyquan/GradLink/actions/workflows/ci.yml) [![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](LICENSE)

GradLink 是一个面向国内应届生的移动端优先求职交流网站。MVP 包含邮箱密码登录、应届生资料门槛、帖子、评论、回复、点赞、收藏、举报、关键词审核和轻量管理后台。

> **重要：** 本项目适合本地开发、学习和功能验证。公开部署前必须配置 HTTPS、强管理员密码、Secure Cookie、可靠备份、限流，以及完善的内容和隐私治理。

## 功能

- 邮箱密码注册、登录和数据库会话。
- 应届生资料完善门槛。
- 帖子、评论、回复、点赞和收藏。
- 标签筛选、用户主页和活动记录。
- 内容举报、关键词审核和管理员处理流程。
- 移动端优先的 Vue 单页界面。

## 技术栈

- 后端：Go、Gin、GORM、SQLite。
- 前端：Vue 3、Vite、TypeScript、Pinia、Vue Router。
- 测试：Go testing、Vitest、Testing Library。

## 项目结构

```text
backend/   Go REST API 和领域模块
frontend/  Vue 3 Web 前端
docs/      产品设计、实施计划和运行说明
```

## 环境要求

- Go 1.26.6 或更高的兼容版本。
- Node.js 20.19+、22.12+ 或 24+（推荐当前 LTS）。
- npm。

## 快速启动

后端（PowerShell）：

```powershell
cd backend
$env:GRADLINK_ADDR="127.0.0.1:8080"
Remove-Item Env:GRADLINK_ADMIN_PASSWORD -ErrorAction SilentlyContinue
go run ./cmd/server
```

上面的启动方式适用于新建的本地数据库：会写入一个禁用的管理员种子账号，不会启用管理员登录。

可选：启用管理员或轮换现有管理员密码（PowerShell）：

```powershell
cd backend
$env:GRADLINK_ADDR="127.0.0.1:8080"
$secureAdminPassword = Read-Host "请输入强管理员密码" -AsSecureString
$adminCredential = New-Object System.Net.NetworkCredential('', $secureAdminPassword)
$env:GRADLINK_ADMIN_PASSWORD = $adminCredential.Password
Remove-Variable secureAdminPassword, adminCredential
go run ./cmd/server
```

输入过程会被隐藏；为使服务器进程能够读取密码，环境变量在当前进程中必须保存明文。

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

管理员账号为 `admin@gradlink.local`。新建数据库且未设置 `GRADLINK_ADMIN_PASSWORD` 时，该账号保持禁用；设置变量启动服务会启用该账号或轮换其现有密码。管理员一旦启用，后续启动时取消设置该变量不会将其禁用或轮换密码；如需在本地禁用或重置管理员，请重置本地数据库。更完整的启动、重置和验证说明见 [docs/run-local.md](docs/run-local.md)。

## 配置

| 变量 | 默认值 | 说明 |
| --- | --- | --- |
| `GRADLINK_ADDR` | `:8080` | 默认监听所有网络接口；本地开发建议设为 `127.0.0.1:8080` |
| `GRADLINK_DB_PATH` | `data/gradlink.db` | SQLite 数据库路径 |
| `GRADLINK_ADMIN_PASSWORD` | 未设置 | 新建数据库时写入禁用的管理员；在后续启动时设置该变量会启用管理员或轮换密码 |
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
