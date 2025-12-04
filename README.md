# Code Auto-Grader

LLM-based code learning platform that uses AI to generate coding exercises, execute code securely in Docker containers, and provide intelligent feedback on student submissions.

## Features

- **AI-Powered Exercise Generation**: Generate custom coding exercises using Google Gemini AI
- **Interactive Code Editor**: Built-in code editor with SolidJS reactive framework
- **Secure Docker Execution**: Run submitted code in isolated Docker containers with resource limits
- **Intelligent Grading**: Weighted scoring (60% tests + 40% AI evaluation) with 80% mastery threshold
- **Learning Paths**: Tech tree navigation with prerequisite-based progression
- **Progress Tracking**: Track completion and mastery across exercises and nodes
- **Multi-Language Support**: JavaScript, Python, and Rust with language-specific container pools

## Technology Stack

**Backend:**
- Go 1.24+ with Echo web framework
- PostgreSQL 15 with pgx driver
- Google Gemini AI API
- Docker for isolated code execution

**Frontend:**
- SolidJS with TypeScript
- Vite build tool with HMR
- Marked for Markdown rendering

## Architecture

### Backend (Go + Echo)

The backend follows a **feature-based vertical slice architecture**:

```
app/
├── main.go                     # Entry point, DI, routes
├── features/
│   ├── submissions/            # Code submission & grading
│   ├── exercises/              # Exercise management
│   ├── tech-trees/             # Learning path navigation
│   ├── progress/               # User progress tracking
│   └── code-execution/         # Docker executor & pool
├── repository/                 # Data access layer
├── models/                     # Shared domain models
├── migrations/                 # SQL migrations
└── shared/
    ├── config/                 # Environment configuration
    ├── database/               # PostgreSQL connection
    ├── gemini/                 # AI client
    └── middleware/             # HTTP middleware
```

Each feature contains:
- `api.go` - HTTP handlers
- `service.go` - Business logic
- `models.go` - Domain types
- `dto.go` - Request/response types
- `*_test.go` - Tests

### Frontend (Vite + SolidJS)

```
web/
├── src/
│   ├── App.tsx                 # Main application component
│   ├── index.tsx               # Entry point
│   └── index.css               # Styles
├── vite.config.ts              # Vite configuration
└── dist/                       # Built production files
```

### API Endpoints

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/api/health` | Health check |
| GET | `/api/tech-tree/:treeId` | Get tech tree with progress |
| GET | `/api/node/:nodeId` | Get node with exercises |
| GET | `/api/exercise/:exerciseId` | Get exercise with history |
| POST | `/api/code/run/:exerciseId` | Run code without grading |
| POST | `/api/submission/:exerciseId` | Submit for full evaluation |
| GET | `/api/pool/status` | Code execution pool status |
| GET | `/api/exercise/sample` | Sample exercise (testing) |
| POST | `/api/exercise/generate` | Generate exercise via AI |

### Security Measures

1. **Code Execution**:
   - Isolated Docker containers per language
   - Non-root user (`coderunner`) in containers
   - Memory limits and execution timeouts
   - Dangerous pattern detection
   - Container pool with configurable sizes

2. **API Security**:
   - Rate limiting (100 requests/15 minutes)
   - Security headers
   - CORS enabled
   - Request ID tracing
   - Structured logging

## Setup Instructions

### Prerequisites

- Go 1.24 or higher
- Node.js 18+ and npm
- Docker and Docker Compose
- PostgreSQL 15 (or use Docker Compose)
- Gemini API key ([get one here](https://makersuite.google.com/app/apikey))

### Installation

1. Clone the repository:
```bash
git clone https://github.com/anttihil/code-auto-grader.git
cd code-auto-grader
```

2. Start the database:
```bash
cd app
docker-compose up -d
```

3. Install Go dependencies:
```bash
cd app
go mod download
```

4. Build the frontend:
```bash
cd web
npm install
npm run build
```

5. Create environment configuration:
```bash
cp .env.example .env
```

6. Configure your `.env` file:
```env
# Server
PORT=5000
ENVIRONMENT=development

# Database
DB_USER=postgres
DB_PASSWORD=postgres
DB_HOST=localhost
DB_PORT=5432
DB_NAME=code_grader
DB_SSLMODE=disable
DB_AUTO_MIGRATE=true

# External Services
GEMINI_API_KEY=your_api_key_here

# Code Execution Pool Sizes
JS_POOL_SIZE=5
PY_POOL_SIZE=5
RUST_POOL_SIZE=2
MAX_CONCURRENT_EXECUTIONS=15
```

### Running the Application

**Development Mode:**

Run backend and frontend separately for hot reloading:

```bash
# Terminal 1: Start the Go backend
cd app && go run main.go

# Terminal 2: Start the Vite dev server
cd web && npm run dev
```

- Backend: `http://localhost:5000`
- Frontend (with HMR): `http://localhost:5173`

**Production Mode:**

```bash
# Build frontend
cd web && npm run build

# Build and run backend (serves built frontend)
cd app
go build -o code-auto-grader
./code-auto-grader
```

Access at `http://localhost:5000`

## Usage

1. **Navigate Learning Paths**:
   - Browse available tech trees
   - View nodes and their prerequisites
   - Track progress through the curriculum

2. **Complete Exercises**:
   - Read exercise description and requirements
   - Write your solution in the code editor
   - Use "Run Code" to test without grading

3. **Submit for Grading**:
   - Click "Submit" for full evaluation
   - Receive weighted score (60% tests + 40% AI)
   - Review detailed AI feedback
   - Achieve 80% mastery to unlock next exercises

## Development

### Commands

```bash
# Backend
cd app && go run main.go          # Start server
cd app && make test               # Run all tests
cd app && make db-reset           # Reset database

# Frontend
cd web && npm run dev             # Dev server with HMR
cd web && npm run build           # Production build

# Database
cd app && migrate create -ext sql -dir migrations -seq name
```

### Adding a New Feature

1. Create feature directory: `app/features/[feature]/`
2. Add required files: `api.go`, `service.go`, `models.go`, `dto.go`
3. Define repository interfaces in `app/repository/`
4. Register routes in `app/main.go`
5. Write tests (100% coverage required)

## API Examples

### Submit Code for Evaluation

```bash
curl -X POST http://localhost:5000/api/submission/[exercise-id] \
  -H "Content-Type: application/json" \
  -d '{
    "code": "fn main() { println!(\"Hello\"); }",
    "language": "rust"
  }'
```

Response:
```json
{
  "submission_id": "uuid",
  "execution_output": "Hello\n",
  "test_results": { "passed": 3, "failed": 0 },
  "ai_feedback": {
    "score": 85,
    "strengths": ["Clean code", "Correct logic"],
    "improvements": ["Add error handling"]
  },
  "test_score": 100,
  "ai_score": 85,
  "total_score": 94,
  "attempt_number": 1
}
```

### Get Tech Tree

```bash
curl http://localhost:5000/api/tech-tree/[tree-id]
```

## Troubleshooting

### Database connection fails
- Verify PostgreSQL is running: `docker-compose ps`
- Check connection settings in `.env`
- Ensure database exists: `createdb code_grader`

### Code execution fails
- Verify Docker is running: `docker ps`
- Check container pool status: `GET /api/pool/status`
- Review container logs: `docker logs [container-id]`

### AI responses are slow or fail
- Verify `GEMINI_API_KEY` is set correctly
- Check Gemini API quota limits
- Consider caching for repeated exercises

## Contributing

Contributions are welcome! Please follow the architecture guidelines in `CLAUDE.md` and ensure:
- Feature-based organization
- 100% test coverage
- No business logic in handlers

## License

ISC
