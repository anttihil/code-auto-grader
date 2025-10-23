package logger

import (
	"context"
	"log/slog"
	"os"
	"time"
)

var Logger *slog.Logger

// Initialize sets up the global logger based on environment
func Initialize(env string) {
	var handler slog.Handler

	switch env {
	case "production":
		// JSON format for production (machine-readable)
		handler = slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
			Level:       slog.LevelInfo,
			AddSource:   true, // Include file:line for errors
			ReplaceAttr: replaceTimestamp,
		})
	case "development":
		// Human-readable format for development
		handler = slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
			Level:     slog.LevelDebug,
			AddSource: false,
		})
	default:
		handler = slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
			Level: slog.LevelInfo,
		})
	}

	Logger = slog.New(handler)
	slog.SetDefault(Logger) // Set as default for standard log package
}

// replaceTimestamp formats timestamps as ISO8601 for log aggregation
func replaceTimestamp(groups []string, a slog.Attr) slog.Attr {
	if a.Key == slog.TimeKey {
		return slog.String("timestamp", a.Value.Time().Format(time.RFC3339Nano))
	}
	return a
}

// WithRequestContext creates a logger with request-scoped fields
func WithRequestContext(ctx context.Context, requestID, method, path string) *slog.Logger {
	return Logger.With(
		slog.String("request_id", requestID),
		slog.String("method", method),
		slog.String("path", path),
	)
}

// WithExecutionContext creates a logger for code execution operations
func WithExecutionContext(ctx context.Context, requestID, language, containerID string) *slog.Logger {
	attrs := []any{
		slog.String("request_id", requestID),
		slog.String("language", language),
		slog.String("component", "executor"),
	}

	if containerID != "" && len(containerID) > 12 {
		attrs = append(attrs, slog.String("container_id", containerID[:12]))
	} else if containerID != "" {
		attrs = append(attrs, slog.String("container_id", containerID))
	}

	return Logger.With(attrs...)
}

// WithAIContext creates a logger for AI service operations
func WithAIContext(ctx context.Context, requestID, operation string) *slog.Logger {
	return Logger.With(
		slog.String("request_id", requestID),
		slog.String("operation", operation),
		slog.String("component", "gemini_service"),
	)
}

// GetRequestID extracts request ID from context
func GetRequestID(ctx context.Context) string {
	if id, ok := ctx.Value("request_id").(string); ok {
		return id
	}
	return "unknown"
}
