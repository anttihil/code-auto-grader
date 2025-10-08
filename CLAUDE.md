# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

Code Auto-Grader is an LLM-powered educational platform that generates coding exercises and provides AI-driven feedback on student submissions. It's a full-stack application with a Go backend and SolidJS TypeScript frontend.

## Architecture

**Backend (`/app` directory):**
- Go HTTP server using standard library
- Google Gemini AI integration for exercise generation and code evaluation
- Sandboxed code execution for JavaScript and Python
- Security middleware (rate limiting, CORS, security headers)

**Frontend (`/web` directory):**
- SolidJS with TypeScript and Vite
- Real-time code editor with execution and AI feedback
- Responsive CSS design with gradients

## Development Commands

**Backend (from `/app` directory):**
```bash
go mod download          # Install dependencies
go run main.go          # Run development server on :3000
go build -o code-auto-grader  # Build production binary
```

**Frontend (from `/web` directory):**
```bash
npm install             # Install dependencies
npm run dev            # Development server on :5173 (proxies API to :3000)
npm run build          # Production build (outputs to dist/)
npm run typecheck      # TypeScript type checking
```

**Full Development Setup:**
1. Set up backend: `cd app && cp .env.example .env` (add GEMINI_API_KEY)
2. Build Docker images: `cd app && docker build -f docker/javascript.Dockerfile -t code-grader-javascript . && docker build -f docker/python.Dockerfile -t code-grader-python .`
3. Run backend: `cd app && go run main.go` (builds images automatically if not present)
4. Run frontend: `cd web && npm run dev` (in separate terminal)

## Key API Endpoints

- `GET /api/exercise/sample` - Get sample exercise
- `POST /api/exercise/generate` - Generate custom exercise (requires difficulty, topic, language)
- `POST /api/submission/run` - Execute code without evaluation
- `POST /api/submission/submit` - Submit code for AI evaluation

## Code Execution System

The backend securely executes JavaScript (Node.js) and Python code in Docker containers:
- Isolated containers with no network access
- Read-only filesystem
- 5-second timeout limits
- Memory limit: 128MB
- CPU limit: 0.5 cores
- Process limit: 50 PIDs
- No capabilities, privilege escalation prevented

## File Structure

- `app/main.go` - Server setup and routing
- `app/handlers/` - HTTP request handlers
- `app/services/` - Business logic (AI, code execution)
- `app/middleware/` - Security, logging, rate limiting
- `web/src/App.tsx` - Main SolidJS component with all UI logic
- `web/vite.config.ts` - Build config with API proxy to backend

## Environment Setup

**Required:**
- Go 1.20+ (project uses 1.24.7)
- Node.js and npm
- Gemini API key (add to `app/.env`)
- **Docker** (required for code execution)

**Docker Setup:**
```bash
# Images are built automatically on first run
# Or build manually:
cd app
docker build -f docker/javascript.Dockerfile -t code-grader-javascript .
docker build -f docker/python.Dockerfile -t code-grader-python .
```

**Security Features:**
- Rate limiting (100 requests per 15 minutes)
- Docker container isolation (network disabled, read-only filesystem, resource limits)
- Security headers and CORS protection

## Production Deployment

1. Build frontend: `cd web && npm run build`
2. Build backend: `cd app && go build -o code-auto-grader`
3. Run: `./app/code-auto-grader` (serves both frontend and API)

## TypeScript Configuration

The frontend uses strict TypeScript with:
- SolidJS JSX support
- Modern ES2020+ features
- Vite bundler with HMR
- Type checking via `npm run typecheck`