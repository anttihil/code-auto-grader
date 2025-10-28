# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

Code Auto-Grader is an LLM-powered educational platform for generating coding exercises and providing AI-driven feedback on student submissions.

**Current Status**: MVP with code execution and AI feedback for single exercises.

**Long-term Vision**: Adaptive learning platform with tech tree navigation, mastery-based progression, and multiple exercise types (see `IMPLEMENTATION_PLAN.md`).

## Current Architecture

**Backend (`/app` directory):**
- Go HTTP server using Echo framework
- PostgreSQL database with pgx driver
- Google Gemini AI integration for exercise generation and code evaluation
- Sandboxed code execution (JavaScript, Python, Rust) via Docker containers
- Security middleware (rate limiting, CORS, security headers)
- Structured logging with request ID tracking
- Repository pattern for clean data access

**Frontend (`/web` directory):**
- SolidJS with TypeScript and Vite
- Code editor with real-time execution and AI feedback
- Single-page application (all logic in `App.tsx`)
- Responsive CSS design with gradients

**Database**: PostgreSQL with complete schema for users, tech trees, nodes, exercises, submissions, and progress tracking. See `app/database/README.md` for setup instructions.

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
├── database/                   # Database connection and utilities
│   ├── database.go            # Connection pool, migrations
│   ├── test_helpers.go        # Test database setup
│   └── README.md              # Database setup guide
├── repository/                 # Data access layer (Repository pattern)
│   ├── repository.go          # Repository interfaces
│   ├── exercise_repo.go       # Exercise CRUD operations
│   ├── submission_repo.go     # Submission CRUD operations
│   ├── tech_tree_repo.go      # Tech tree operations
│   ├── node_repo.go           # Node operations
│   ├── user_repo.go           # User operations
│   └── user_progress_repo.go  # Progress tracking
├── models/
│   └── models.go              # Database models and structs
├── migrations/                 # Database schema migrations
│   ├── 000001_init_schema.up.sql    # Create tables
│   └── 000001_init_schema.down.sql  # Drop tables
├── cmd/
│   └── seed/
│       └── main.go            # Database seeding utility
├── internal/
│   ├── logger/
│   │   └── logger.go          # Logger initialization
│   └── content/               # Content loading and parsing
│       ├── loader.go          # Load exercises from files
│       ├── parser.go          # Parse markdown to models
│       └── types.go           # Content data structures
├── middleware/
│   ├── logger.go              # Structured logging middleware
│   ├── rate_limit.go          # Rate limiting (100 req/15min)
│   └── request_id.go          # Request ID generation
├── docker/
│   ├── javascript.Dockerfile  # Node.js execution environment
│   ├── python.Dockerfile      # Python execution environment
│   └── rust.Dockerfile        # Rust execution environment
├── Makefile                    # Database and development commands
└── integration_test.go        # Integration tests
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

**Database Commands (from `/app` directory):**
```bash
make install-tools           # Install migrate CLI and other dev tools
make docker-db-up            # Start PostgreSQL in Docker
make migrate-up              # Run database migrations
make db-seed                 # Seed database with content
make db-reset                # Drop, migrate, and seed (fresh start)
make test                    # Run all tests
make help                    # Show all available commands
```

**Full Development Setup:**
1. Set up backend: `cd app && cp .env.example .env` (add GEMINI_API_KEY and DATABASE_URL)
2. Install development tools:
   ```bash
   cd app
   make install-tools  # Installs golang-migrate CLI
   ```
3. Start PostgreSQL database:
   ```bash
   make docker-db-up  # Or use your local PostgreSQL
   ```
4. Set up database:
   ```bash
   make migrate-up    # Run migrations
   make db-seed       # Seed with Rust exercises
   ```
5. Build Docker images for code execution (optional, auto-built on first run):
   ```bash
   docker build -f docker/javascript.Dockerfile -t code-grader-javascript .
   docker build -f docker/python.Dockerfile -t code-grader-python .
   docker build -f docker/rust.Dockerfile -t code-grader-rust .
   ```
6. Run backend: `go run main.go`
7. Run frontend: `cd ../web && npm run dev` (in separate terminal)
8. Open browser: http://localhost:5173

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
- PostgreSQL 14+ (or use Docker: `make docker-db-up`)
- golang-migrate CLI tool (install: `make install-tools`)
- Gemini API key (get from https://aistudio.google.com/app/apikey)

**Environment Variables (app/.env):**
```bash
# AI API
GEMINI_API_KEY=your_api_key_here

# Server Configuration
PORT=5000                           # Optional, defaults to 5000
ENVIRONMENT=development             # "development" or "production"

# Database Configuration
DATABASE_URL=postgres://code_grader:dev_password@localhost:5432/code_auto_grader?sslmode=disable
TEST_DATABASE_URL=postgres://code_grader:test_password@localhost:5433/code_auto_grader_test?sslmode=disable
DB_AUTO_MIGRATE=true                # Run migrations on startup (optional)
DB_MIGRATIONS_PATH=file://migrations  # Path to migration files (optional)

# Code Execution Pool Sizes (optional, defaults shown)
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

## Database Architecture

**Schema Overview:**

The database uses PostgreSQL with a complete schema for the learning platform:

**Core Tables:**
- **users** - User accounts (UUID, username, email, display_name)
- **tech_trees** - Learning paths (e.g., "Rust Fundamentals")
- **nodes** - Topics within tech trees with DAG structure via `node_prerequisites`
- **exercises** - Individual coding exercises with JSONB content
- **submissions** - User submission history (code, execution results, AI feedback, scores)
- **user_progress** - Progress tracking at node level (mastery scores, completion status)

**Key Features:**
- Automatic `updated_at` triggers on all tables
- JSONB storage for flexible exercise content (starter code, test cases, hints, rubrics)
- DAG (directed acyclic graph) for node prerequisites
- Mastery-based progression (80% threshold by default)
- Submission scoring: 60% test results + 40% AI evaluation

**Repository Pattern:**

Clean data access layer with interfaces and implementations:
- `ExerciseRepository` - CRUD operations for exercises
- `SubmissionRepository` - Save and retrieve submissions with history
- `TechTreeRepository` - Tech tree management
- `NodeRepository` - Node operations and prerequisite queries
- `UserRepository` - User account operations
- `UserProgressRepository` - Progress tracking and mastery calculations

**Migrations:**

Using golang-migrate for schema versioning:
- Migration files in `app/migrations/`
- Run with `make migrate-up` / `make migrate-down`
- Auto-migration on app startup (configurable via `DB_AUTO_MIGRATE`)

**Content Seeding:**

Structured content system for exercises:
- Exercise files: Markdown with JSON frontmatter
- Content loader: Parses markdown sections (starter code, tests, hints)
- Seed command: `make db-seed` or `go run cmd/seed/main.go`
- 9 complete Rust exercises included (Nodes 1-3)

See `app/database/README.md` for detailed database setup and usage instructions.

## TypeScript Configuration

The frontend uses strict TypeScript with:
- SolidJS JSX support
- Modern ES2020+ features
- Vite bundler with HMR
- Type checking via `npm run typecheck`

## Planned Features (See IMPLEMENTATION_PLAN.md)

The following are **not yet implemented** but planned for future development:

**Database Layer (✅ COMPLETED):**
- ✅ PostgreSQL schema with all tables
- ✅ Models: `TechTree`, `Node`, `Exercise`, `Submission`, `UserProgress`
- ✅ Migration system using golang-migrate
- ✅ Repository pattern for data access
- ✅ Content seeding system
- 🔄 **In Progress**: Integrate repositories with API handlers

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
