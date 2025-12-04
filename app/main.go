package main

import (
	"context"
	"log/slog"
	"os"

	"net/http"
	"time"

	codeexecution "github.com/anttihil/code-auto-grader/features/code-execution"
	"github.com/anttihil/code-auto-grader/features/exercises"
	learningmaterials "github.com/anttihil/code-auto-grader/features/learning-materials"
	"github.com/anttihil/code-auto-grader/features/progress"
	"github.com/anttihil/code-auto-grader/features/submissions"
	techtrees "github.com/anttihil/code-auto-grader/features/tech-trees"
	"github.com/anttihil/code-auto-grader/internal/logger"
	"github.com/anttihil/code-auto-grader/repository"
	"github.com/anttihil/code-auto-grader/shared/config"
	"github.com/anttihil/code-auto-grader/shared/database"
	"github.com/anttihil/code-auto-grader/shared/gemini"
	customMiddleware "github.com/anttihil/code-auto-grader/shared/middleware"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	// Load environment variables from .env file
	if err := godotenv.Load(); err != nil {
		// Use standard output before logger is initialized
		slog.Info("No .env file found, using system environment variables")
	}

	// Load and validate configuration (single source of truth)
	cfg, err := config.Load()
	if err != nil {
		slog.Error("configuration error", slog.String("error", err.Error()))
		os.Exit(1)
	}

	// Initialize logger based on environment
	logger.Initialize(cfg.Environment)

	logger.Logger.Info("application starting",
		slog.String("environment", cfg.Environment),
		slog.String("port", cfg.Port),
	)

	// Initialize database
	logger.Logger.Info("initializing database connection")
	ctx := context.Background()
	dbConfig := database.DefaultConfig()
	dbConfig.URL = cfg.DatabaseURL
	dbConfig.AutoMigrate = cfg.DBAutoMigrate
	dbConfig.MigrationsPath = cfg.DBMigrationsPath
	db, err := database.New(ctx, dbConfig)
	if err != nil {
		logger.Logger.Error("failed to initialize database",
			slog.String("error", err.Error()),
		)
		os.Exit(1)
	}
	defer func() {
		logger.Logger.Info("closing database connection")
		db.Close()
	}()
	logger.Logger.Info("database connection established successfully")

	// Initialize repositories
	logger.Logger.Info("initializing repositories")
	techTreeRepo := repository.NewTechTreeRepo(db.Pool)
	nodeRepo := repository.NewNodeRepo(db.Pool)
	exerciseRepo := repository.NewExerciseRepo(db.Pool)
	submissionRepo := repository.NewSubmissionRepo(db.Pool)
	learningMaterialRepo := repository.NewLearningMaterialRepo(db.Pool)
	// userRepo := repository.NewUserRepo(db.Pool) // Not used yet in dogfooding MVP
	logger.Logger.Info("repositories initialized successfully")

	// Initialize feature services
	logger.Logger.Info("initializing feature services")
	progressRepo := progress.NewRepository(db.Pool)
	progressService := progress.NewService(progressRepo)
	logger.Logger.Info("progress service initialized")

	// Initialize services
	logger.Logger.Info("initializing services")

	codeExecutor, err := codeexecution.NewCodeExecutionServiceWithConfig(cfg)
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

	geminiService := gemini.NewGeminiServiceWithConfig(cfg)
	logger.Logger.Info("gemini service initialized")

	// Create adapters to bridge services to domain interfaces
	logger.Logger.Info("creating service adapters")
	executorAdapter := submissions.NewExecutorAdapter(codeExecutor)
	geminiGradingAdapter := submissions.NewGeminiAdapter(geminiService)
	geminiContentAdapter := exercises.NewGeminiAdapter(geminiService)

	// Create domain services (business logic layer)
	logger.Logger.Info("initializing domain services")
	gradingService := submissions.NewGradingService(
		executorAdapter,
		geminiGradingAdapter,
		submissionRepo,
		exerciseRepo,
		progressService,
	)
	learningService := techtrees.NewLearningService(
		techTreeRepo,
		nodeRepo,
		exerciseRepo,
		progressService,
		submissionRepo,
	)
	contentService := exercises.NewContentService(
		geminiContentAdapter,
		exerciseRepo,
		submissionRepo,
	)
	learningMaterialService := learningmaterials.NewService(learningMaterialRepo)
	logger.Logger.Info("domain services initialized successfully")

	// Create API handlers (HTTP adapters)
	logger.Logger.Info("initializing API handlers")
	submissionHandler := submissions.NewSubmissionHandler(gradingService)
	techTreeHandler := techtrees.NewTechTreeHandler(learningService)
	nodeHandler := techtrees.NewNodeHandler(learningService)
	exerciseHandler := exercises.NewExerciseHandler(contentService)
	learningMaterialHandler := learningmaterials.NewHandler(learningMaterialService)
	logger.Logger.Info("API handlers initialized with domain services")

	// Create Echo instance
	e := echo.New()

	// Middleware (order matters!)
	e.Use(customMiddleware.RequestID())                     // Generate request ID first
	e.Use(customMiddleware.StructuredLogger(logger.Logger)) // Log with request ID
	e.Use(middleware.Recover())
	e.Use(middleware.CORS())
	e.Use(middleware.Secure())
	e.Use(customMiddleware.RateLimit(logger.Logger)) // Rate limiting with logging

	// Serve static files from web/dist directory (Vite build output)
	e.Static("/", "./web/dist")

	// API routes
	e.GET("/api/health", func(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]interface{}{
			"status":    "ok",
			"timestamp": time.Now().Format(time.RFC3339),
		})
	})

	// Exercise routes (legacy)
	e.GET("/api/exercise/sample", exerciseHandler.GetSampleExercise)
	e.POST("/api/exercise/generate", exerciseHandler.GenerateExercise)

	// Tech tree navigation routes
	e.GET("/api/tech-tree/:treeId", techTreeHandler.GetTechTree)
	e.GET("/api/node/:nodeId", nodeHandler.GetNode)
	e.GET("/api/exercise/:exerciseId", exerciseHandler.GetExercise)

	// Submission routes
	e.POST("/api/code/run/:exerciseId", submissionHandler.RunCode)
	e.POST("/api/submission/:exerciseId", submissionHandler.SubmitToExercise)
	e.GET("/api/pool/status", submissionHandler.GetPoolStatus)

	// Learning materials routes
	e.GET("/api/nodes/:nodeId/materials", learningMaterialHandler.GetByNode)
	e.POST("/api/nodes/:nodeId/materials", learningMaterialHandler.Create)
	e.GET("/api/materials/:id", learningMaterialHandler.GetByID)
	e.PUT("/api/materials/:id", learningMaterialHandler.Update)
	e.DELETE("/api/materials/:id", learningMaterialHandler.Delete)

	// Start server
	logger.Logger.Info("application ready",
		slog.String("port", cfg.Port),
		slog.String("address", "http://localhost:"+cfg.Port),
	)
	e.Logger.Fatal(e.Start(":" + cfg.Port))
}
