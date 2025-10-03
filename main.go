package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/anttihil/code-auto-grader/handlers"
	"github.com/anttihil/code-auto-grader/middleware"
	"github.com/joho/godotenv"
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

	// Create HTTP server with middlewares
	mux := http.NewServeMux()

	// Serve static files from public directory
	fs := http.FileServer(http.Dir("./public"))
	mux.Handle("/", fs)

	// API routes
	mux.HandleFunc("/api/health", handlers.HealthCheck)
	mux.HandleFunc("/api/exercise/sample", handlers.GetSampleExercise)
	mux.HandleFunc("/api/exercise/generate", handlers.GenerateExercise)
	mux.HandleFunc("/api/submission/run", handlers.RunCode)
	mux.HandleFunc("/api/submission/submit", handlers.SubmitCode)

	// Wrap mux with middlewares
	handler := middleware.CORS(
		middleware.RateLimit(
			middleware.Security(
				middleware.Logger(mux),
			),
		),
	)

	// Create server
	server := &http.Server{
		Addr:         ":" + port,
		Handler:      handler,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	log.Printf("Server is running on port %s\n", port)
	log.Printf("Access the app at http://localhost:%s\n", port)

	if err := server.ListenAndServe(); err != nil {
		log.Fatal(fmt.Errorf("server failed to start: %w", err))
	}
}
