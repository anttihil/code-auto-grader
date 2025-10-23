package main

import (
	"log/slog"
	"os"

	"github.com/anttihil/code-auto-grader/handlers"
	"github.com/anttihil/code-auto-grader/internal/logger"
	customMiddleware "github.com/anttihil/code-auto-grader/middleware"
	"github.com/anttihil/code-auto-grader/services"
	"github.com/anttihil/code-auto-grader/services/executor"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		// Use standard output before logger is initialized
		slog.Info("No .env file found, using system environment variables")
	}

	// Initialize logger based on environment
	env := os.Getenv("ENVIRONMENT")
	if env == "" {
		env = "development"
	}
	logger.Initialize(env)

	logger.Logger.Info("application starting",
		slog.String("environment", env),
	)

	// Get port from environment or use default
	port := os.Getenv("PORT")
	if port == "" {
		port = "5000"
	}

	// Initialize services
	logger.Logger.Info("initializing services")

	codeExecutor, err := executor.NewCodeExecutionService()
	if err != nil {
		logger.Logger.Error("failed to initialize code execution service",
			slog.String("error", err.Error()),
		)
		os.Exit(1)
	}
	defer func() {
		logger.Logger.Info("cleaning up resources")
		codeExecutor.Cleanup()
		logger.Logger.Info("cleanup completed")
	}()
	logger.Logger.Info("code execution service initialized successfully")

	geminiService := services.NewGeminiService()
	logger.Logger.Info("gemini service initialized")

	// Create handlers with injected dependencies
	submissionHandler := handlers.NewSubmissionHandler(codeExecutor, geminiService)
	exerciseHandler := handlers.NewExerciseHandler(geminiService)
	logger.Logger.Info("handlers initialized with dependencies")

	// Create Echo instance
	e := echo.New()

	// Middleware (order matters!)
	e.Use(customMiddleware.RequestID())                  // Generate request ID first
	e.Use(customMiddleware.StructuredLogger(logger.Logger)) // Log with request ID
	e.Use(middleware.Recover())
	e.Use(middleware.CORS())
	e.Use(middleware.Secure())
	e.Use(customMiddleware.RateLimit(logger.Logger)) // Rate limiting with logging

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
	logger.Logger.Info("application ready",
		slog.String("port", port),
		slog.String("address", "http://localhost:"+port),
	)
	e.Logger.Fatal(e.Start(":" + port))
}
