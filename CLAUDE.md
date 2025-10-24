# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

Code Auto-Grader is an LLM-powered educational platform for generating coding exercises and providing AI-driven feedback on student submissions.

**Current Status**: MVP with code execution and AI feedback for single exercises.

**Long-term Vision**: Adaptive learning platform with tech tree navigation, mastery-based progression, and multiple exercise types (see `IMPLEMENTATION_PLAN.md`).

## Current Architecture

**Backend (`/app` directory):**
- Go HTTP server using Echo framework
- Google Gemini AI integration for exercise generation and code evaluation
- Sandboxed code execution (JavaScript, Python, Rust) via Docker containers
- Security middleware (rate limiting, CORS, security headers)
- Structured logging with request ID tracking

**Frontend (`/web` directory):**
- SolidJS with TypeScript and Vite
- Code editor with real-time execution and AI feedback
- Single-page application (all logic in `App.tsx`)
- Responsive CSS design with gradients

**No Database Yet**: All data is ephemeral (sessions, no persistence). Database integration planned for MVP Phase 1.

## Code Execution System

The backend securely executes code in Docker containers for multiple languages:

**Supported Languages:**
- **JavaScript** (Node.js) - interpreted
- **Python** - interpreted
- **Rust** - compiled (uses rustc)

**Security Sandbox:**
- Isolated containers with no network access
- Read-only filesystem
- 5-second timeout limits
- Memory limit: 128MB
- CPU limit: 0.5 cores
- Process limit: 50 PIDs
- No capabilities, privilege escalation prevented

**Container Pool Management:**
- Pre-warmed container pools per language (reduces cold start latency)
- Automatic container rotation based on execution count and age
- Health checks to detect and replace unhealthy containers
- Configurable pool sizes via environment variables (`JS_POOL_SIZE`, `PYTHON_POOL_SIZE`, `RUST_POOL_SIZE`)

## File Structure

**Backend:**
```
app/
├── main.go                     # Server setup, routing, dependency injection
├── handlers/
│   ├── exercise.go            # Exercise generation endpoints
│   ├── health.go              # Health check endpoint
│   └── submission.go          # Code execution and evaluation endpoints
├── services/
│   ├── gemini.go              # Gemini API integration
│   └── executor/              # Code execution service
│       ├── executor.go        # Main execution logic
│       ├── config.go          # Language configurations
│       ├── pool.go            # Container pool management
│       ├── docker.go          # Docker container operations
│       └── maintenance.go     # Health checks and rotation
├── middleware/
│   ├── logger.go              # Structured logging middleware
│   ├── rate_limit.go          # Rate limiting (100 req/15min)
│   └── request_id.go          # Request ID generation
├── internal/
│   └── logger/
│       └── logger.go          # Logger initialization
└── docker/
    ├── javascript.Dockerfile  # Node.js execution environment
    ├── python.Dockerfile      # Python execution environment
    └── rust.Dockerfile        # Rust execution environment
```

**Frontend:**
```
web/
├── src/
│   ├── App.tsx               # Main component (all UI logic)
│   ├── index.tsx             # Entry point
│   └── styles.css            # Global styles
├── vite.config.ts            # Build config with API proxy to :5000
└── dist/                     # Production build output (gitignored)
```

## Development Commands

**Backend (from `/app` directory):**
```bash
go mod download               # Install dependencies
go run main.go               # Run development server on :5000 (or PORT env var)
go build -o code-auto-grader # Build production binary
go test ./...                # Run tests
```

**Frontend (from `/web` directory):**
```bash
npm install                  # Install dependencies
npm run dev                 # Development server on :5173 (proxies API to :5000)
npm run build               # Production build (outputs to dist/)
npm run typecheck           # TypeScript type checking
```

**Full Development Setup:**
1. Set up backend: `cd app && cp .env.example .env` (add GEMINI_API_KEY)
2. Build Docker images (optional, auto-built on first run):
   ```bash
   cd app
   docker build -f docker/javascript.Dockerfile -t code-grader-javascript .
   docker build -f docker/python.Dockerfile -t code-grader-python .
   docker build -f docker/rust.Dockerfile -t code-grader-rust .
   ```
3. Run backend: `cd app && go run main.go`
4. Run frontend: `cd web && npm run dev` (in separate terminal)
5. Open browser: http://localhost:5173

## API Endpoints

**Exercise Generation:**
- `GET /api/exercise/sample` - Get sample exercise (hardcoded)
- `POST /api/exercise/generate` - Generate custom exercise via Gemini
  - Body: `{ "difficulty": "easy|medium|hard", "topic": "string", "language": "javascript|python|rust" }`

**Code Submission:**
- `POST /api/submission/run` - Execute code without AI evaluation
  - Body: `{ "code": "string", "language": "javascript|python|rust" }`
  - Returns: `{ "output": "string", "error": "string", "executionTime": number }`
- `POST /api/submission/submit` - Execute code and get AI feedback
  - Body: `{ "code": "string", "language": "string", "exerciseDescription": "string" }`
  - Returns: `{ "output": "string", "feedback": "string", "executionTime": number }`

**Monitoring:**
- `GET /api/health` - Health check
- `GET /api/pool/status` - Container pool status (per language)

## Environment Setup

**Required:**
- Go 1.20+ (project uses 1.24.7)
- Node.js and npm
- Docker (required for code execution)
- Gemini API key (get from https://aistudio.google.com/app/apikey)

**Environment Variables (app/.env):**
```bash
GEMINI_API_KEY=your_api_key_here
PORT=5000                           # Optional, defaults to 5000
ENVIRONMENT=development             # "development" or "production"

# Optional: Override pool sizes (defaults shown)
JS_POOL_SIZE=3
PYTHON_POOL_SIZE=3
RUST_POOL_SIZE=2
```

**Security Features:**
- Rate limiting (100 requests per 15 minutes per IP)
- Docker container isolation (network disabled, read-only filesystem, resource limits)
- Security headers (HSTS, X-Frame-Options, etc.)
- CORS protection
- Structured logging with request ID tracking

## Production Deployment

1. Build frontend: `cd web && npm run build`
2. Build backend: `cd app && go build -o code-auto-grader`
3. Run: `./app/code-auto-grader` (serves both frontend and API)
4. Ensure Docker is installed on production host

**Infrastructure:**
- Backend: Single binary deployed to VPS/cloud (DigitalOcean, Render, Fly.io)
- Docker: Required on production host for code execution
- CDN: Optional for frontend assets (Cloudflare, CloudFront)

## TypeScript Configuration

The frontend uses strict TypeScript with:
- SolidJS JSX support
- Modern ES2020+ features
- Vite bundler with HMR
- Type checking via `npm run typecheck`

## Planned Features (See IMPLEMENTATION_PLAN.md)

The following are **not yet implemented** but planned for future development:

**Database Layer:**
- PostgreSQL for user progress, tech trees, exercises
- Models: `TechTree`, `Node`, `Exercise`, `Submission`, `UserProgress`
- Migration system for schema evolution

**Exercise Types:**
- Coding problems (currently supported)
- Debugging exercises (broken code to fix)
- Project-based (multi-step tasks)
- Conceptual questions (essay with diagrams)
- System design (architecture diagrams)

**User System:**
- Authentication (OAuth + email/password)
- User profiles and progress tracking
- Exercise attempt history

**AI Enhancements:**
- Contextual feedback (consider student history)
- Progression agent (recommend next topic)
- Content generator (create exercises from documents)
- Motivational coaching

**Tech Tree Navigation:**
- Visual tech tree interface
- Mastery-based progression
- Unlockable nodes based on performance

See `IMPLEMENTATION_PLAN.md` for detailed roadmap and `BUSINESS_PLAN.md` for strategic vision.

## Development Philosophy

**Current Phase**: Building MVP foundation with focus on:
1. **Reliability**: Stable code execution, no container escapes
2. **Quality feedback**: Useful AI responses, not generic
3. **Performance**: Fast execution via container pooling

**Avoiding**:
- Premature optimization (don't build features before validating need)
- Over-engineering (keep simple until scale requires complexity)
- Scope creep (MVP first, then iterate based on user feedback)
