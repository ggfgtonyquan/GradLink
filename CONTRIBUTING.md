# 贡献指南

感谢你关注 GradLink 并愿意参与改进！无论是修复问题、完善文档、改进测试还是提出新想法，都欢迎通过清晰、可复现的方式参与。开始前请先阅读 [README.md](README.md) 和本地运行说明 [docs/run-local.md](docs/run-local.md)。

## 开始之前

- 请先搜索现有 Issue，避免重复报告或重复实现。
- 较大的功能请先创建 Issue，说明目标、范围和拟采用的方案，待讨论后再开始实现。
- 安全问题请按照 [SECURITY.md](SECURITY.md) 私下报告，不要在公开 Issue、Discussion 或 Pull Request 中披露。

## 分支与提交

请从最新的 `main` 创建分支，例如 `fix/login-error` 或 `feat/post-search`。不要直接向 `main` 推送功能代码，所有改动都应通过 Pull Request 进入主分支。

请沿用项目现有约定：

- Go 代码使用 `gofmt`；后端按领域组织模块，并保持 handler/service/repository 分层。
- 前端使用 Vue 单文件组件（SFC）与 TypeScript；组件名使用 PascalCase，变量和函数使用 camelCase。
- MVP 面向用户的文案保持中文。
- 修复 bug 或新增行为时补充回归测试。

提交信息使用 Conventional Commits，例如：

- `feat: add post search`
- `fix: clamp reaction counters`
- `docs: clarify local setup`

## 提交前验证

后端：

```powershell
cd backend
go test ./... -count=1
go run golang.org/x/vuln/cmd/govulncheck@latest ./...
```

前端：

```powershell
cd frontend
npm test -- --run
npm run build
npm audit --audit-level=high
```

## Pull Request 要求

Pull Request 请包含：

- 改动目的和内容摘要；
- 实际运行过的验证命令及结果；
- 关联 Issue 或背景说明；
- 涉及 UI 改动时的截图或录屏。

所有贡献均以 MIT License 授权，并须遵守 [CODE_OF_CONDUCT.md](CODE_OF_CONDUCT.md)。
