package main

import (
	"log"
	"os"

	"github.com/anttihil/code-auto-grader/handlers"
	customMiddleware "github.com/anttihil/code-auto-grader/middleware"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using system environment variables")
	}

	// Get port from environment or use default
	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
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
	e.GET("/api/exercise/sample", handlers.GetSampleExercise)
	e.POST("/api/exercise/generate", handlers.GenerateExercise)
	e.POST("/api/submission/run", handlers.RunCode)
	e.POST("/api/submission/submit", handlers.SubmitCode)

	// Start server
	log.Printf("Server is running on port %s\n", port)
	log.Printf("Access the app at http://localhost:%s\n", port)
	e.Logger.Fatal(e.Start(":" + port))
}
