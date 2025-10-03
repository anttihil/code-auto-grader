# Code Auto-Grader 🎓

LLM-based code learning platform that uses AI to generate coding exercises and provide intelligent feedback on student submissions.

## Features

- 🤖 **AI-Powered Exercise Generation**: Generate custom coding exercises using Gemini AI
- 📝 **Interactive Code Editor**: Built-in code editor with syntax highlighting (CodeMirror)
- ⚡ **Secure Code Execution**: Run submitted code securely in an isolated environment
- 💬 **Intelligent Feedback**: Get detailed AI-powered feedback on code submissions
- 🔒 **Security Features**: Rate limiting, helmet protection, and code sanitization
- 🌐 **Multi-Language Support**: JavaScript and Python support (extensible)

## Architecture

### Backend (Go)
- **API Endpoints**:
  - `/api/exercise/generate` - Generate new exercises using Gemini AI
  - `/api/exercise/sample` - Get a sample exercise for testing
  - `/api/submission/submit` - Submit code for AI evaluation
  - `/api/submission/run` - Run code without evaluation
  - `/api/health` - Health check endpoint

- **Services**:
  - `geminiService`: Integration with Google Gemini API for exercise generation and code evaluation
  - `codeExecutionService`: Secure code execution with timeout and resource limits

### Frontend (HTML/CSS/JavaScript)
- Modern, responsive web interface
- Built-in code editor with syntax highlighting
- Real-time code execution and feedback display

### Security Measures
1. **Code Execution**:
   - Runs in isolated temporary files
   - 5-second timeout limit
   - Resource limits (128MB memory)
   - Dangerous pattern detection (blocks fs, child_process, etc.)
   - Automatic cleanup of temporary files

2. **API Security**:
   - Rate limiting (100 requests per 15 minutes)
   - Security headers for HTTP protection
   - CORS enabled
   - Input validation

## Setup Instructions

### Prerequisites
- Go 1.20 or higher
- A Gemini API key (get one at https://makersuite.google.com/app/apikey)

### Installation

1. Clone the repository:
```bash
git clone https://github.com/anttihil/code-auto-grader.git
cd code-auto-grader
```

2. Install Go dependencies:
```bash
go mod download
```

3. Create a `.env` file from the example:
```bash
cp .env.example .env
```

4. Add your Gemini API key to the `.env` file:
```
GEMINI_API_KEY=your_actual_api_key_here
PORT=3000
```

### Running the Application

Start the server:
```bash
# Build the application
go build -o code-auto-grader

# Run the server
./code-auto-grader
```

Or run directly without building:
```bash
go run main.go
```

The application will be available at `http://localhost:3000`

## Usage

1. **Load an Exercise**:
   - Click "Load Sample" to try a demo exercise
   - Or use "Generate New Exercise" to create a custom problem using AI

2. **Write Code**:
   - Use the code editor to write your solution
   - Select the appropriate programming language

3. **Test Your Code**:
   - Click "Run Code" to execute and see output
   - Check for errors or issues

4. **Submit for Grading**:
   - Click "Submit for Grading" to get AI-powered feedback
   - Review the detailed feedback including:
     - Score (0-100)
     - Strengths and areas for improvement
     - Code quality assessment
     - Specific suggestions

## API Examples

### Generate an Exercise

```bash
curl -X POST http://localhost:3000/api/exercise/generate \
  -H "Content-Type: application/json" \
  -d '{
    "difficulty": "beginner",
    "topic": "arrays",
    "language": "javascript"
  }'
```

### Submit Code for Evaluation

```bash
curl -X POST http://localhost:3000/api/submission/submit \
  -H "Content-Type: application/json" \
  -d '{
    "code": "function sum(a, b) { return a + b; }",
    "exercise": "Write a function that sums two numbers",
    "language": "javascript"
  }'
```

## Project Structure

```
code-auto-grader/
├── main.go                    # Main server file
├── handlers/
│   └── handlers.go            # HTTP handlers for all endpoints
├── middleware/
│   └── middleware.go          # Security, CORS, rate limiting middleware
├── services/
│   ├── gemini.go              # Gemini AI integration
│   └── code_execution.go      # Secure code execution
├── public/
│   ├── index.html             # Frontend UI
│   ├── styles.css             # Styling
│   └── app.js                 # Frontend JavaScript
├── .env.example               # Environment variables template
├── .gitignore
├── go.mod                     # Go module definition
├── go.sum                     # Go dependencies checksums
└── README.md
```

## Security Considerations

⚠️ **Important Security Notes**:

1. **Code Execution**: The current implementation runs code on the server with basic security measures. For production use, consider:
   - Using Docker containers for complete isolation
   - Implementing a sandboxed execution environment
   - Using cloud-based code execution services (e.g., AWS Lambda, Azure Functions)

2. **API Key**: Never commit your `.env` file with real API keys to version control

3. **Rate Limiting**: Adjust rate limits based on your expected usage

4. **Resource Limits**: Monitor and adjust memory/CPU limits for code execution

## Future Enhancements

- [ ] Docker-based code execution for better isolation
- [ ] User authentication and progress tracking
- [ ] Exercise database and categories
- [ ] More programming language support
- [ ] Test case validation
- [ ] Code plagiarism detection
- [ ] Learning paths and curriculum
- [ ] Leaderboard and gamification

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## License

ISC

## Troubleshooting

### Server won't start
- Check if port 3000 is already in use
- Verify that all dependencies are installed (`npm install`)
- Ensure `.env` file exists with valid `GEMINI_API_KEY`

### Code execution fails
- Ensure Node.js and Python3 are installed on the system
- Check file system permissions for `/tmp` directory
- Verify timeout settings are appropriate

### AI responses are slow
- This is normal for AI-powered features
- Consider implementing caching for common exercises
- Check your Gemini API quota and rate limits
