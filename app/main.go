package main

import (
	"log"
	"os"

	"github.com/anttihil/code-auto-grader/handlers"
	customMiddleware "github.com/anttihil/code-auto-grader/middleware"
	"github.com/anttihil/code-auto-grader/services"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using system environment variables")
	}

	// Initialize services
	log.Println("Initializing services...")

	codeExecutor, err := services.NewCodeExecutionService()
	if err != nil {
		log.Fatalf("Failed to initialize code execution service: %v", err)
	}
	defer codeExecutor.Cleanup() // Clean up container pools on shutdown
	log.Println("Code execution service initialized successfully")

	geminiService := services.NewGeminiService()
	log.Println("Gemini service initialized")

	// Create handlers with injected dependencies
	submissionHandler := handlers.NewSubmissionHandler(codeExecutor)
	exerciseHandler := handlers.NewExerciseHandler(geminiService)
	log.Println("Handlers initialized with dependencies")

	// Get port from environment or use default
	port := os.Getenv("PORT")
	if port == "" {
		port = "5000"
	}

	// Create Echo instance
	e := echo.New()

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORS())
	e.Use(middleware.Secure())
	e.Use(customMiddleware.RateLimit())

	// Serve static files from web/dist directory (Vite build output)
	e.Static("/", "./web/dist")

	// API routes
	e.GET("/api/health", handlers.HealthCheck)
	e.GET("/api/exercise/sample", exerciseHandler.GetSampleExercise)
	e.POST("/api/exercise/generate", exerciseHandler.GenerateExercise)
	e.POST("/api/submission/run", submissionHandler.RunCode)
	e.POST("/api/submission/submit", submissionHandler.SubmitCode)
	e.GET("/api/pool/status", submissionHandler.GetPoolStatus)

	// Start server
	log.Printf("Server is running on port %s\n", port)
	log.Printf("Access the app at http://localhost:%s\n", port)
	e.Logger.Fatal(e.Start(":" + port))
}
