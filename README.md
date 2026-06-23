# GradLink

GradLink is currently at the Task 1 scaffold stage. This repository contains a minimal Go backend health check and a minimal Vue mobile shell.

## Scaffold

- `backend/`: Go module with a Gin server entry point and `GET /api/health`.
- `frontend/`: Vue 3, Vite, TypeScript, and Vitest app with a light mobile shell, top bar, content area, and bottom navigation.

## Task 1 Commands

Backend:

```bash
cd backend
go test ./internal/http -run TestHealthEndpoint -v
go run ./cmd/server
```

Frontend:

```bash
cd frontend
npm install
npm test -- --run src/components/__tests__/AppShell.test.ts
npm run build
npm run dev
```

The backend listens on `http://localhost:8080` when started. The Vite dev server prints the frontend URL after `npm run dev`.
