package middleware

import (
	"log/slog"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

// RateLimit implements simple rate limiting
type rateLimiter struct {
	visitors map[string]*visitor
	mu       sync.RWMutex
}

type visitor struct {
	lastSeen time.Time
	count    int
}

var limiter = &rateLimiter{
	visitors: make(map[string]*visitor),
}

func init() {
	// Cleanup old visitors every 5 minutes
	go func() {
		for {
			time.Sleep(5 * time.Minute)
			limiter.mu.Lock()
			for ip, v := range limiter.visitors {
				if time.Since(v.lastSeen) > 15*time.Minute {
					delete(limiter.visitors, ip)
				}
			}
			limiter.mu.Unlock()
		}
	}()
}

func RateLimit(logger *slog.Logger) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// Only apply rate limiting to API endpoints
			path := c.Request().URL.Path
			if !strings.HasPrefix(path, "/api/") {
				return next(c)
			}

			ip := c.RealIP()

			limiter.mu.Lock()
			v, exists := limiter.visitors[ip]
			if !exists {
				limiter.visitors[ip] = &visitor{lastSeen: time.Now(), count: 1}
				limiter.mu.Unlock()
				return next(c)
			}

			// Reset count if 15 minutes have passed
			if time.Since(v.lastSeen) > 15*time.Minute {
				v.count = 1
				v.lastSeen = time.Now()
			} else {
				v.count++
			}
			v.lastSeen = time.Now()

			if v.count > 100 {
				limiter.mu.Unlock()
				logger.Warn("rate limit exceeded",
					slog.String("ip", ip),
					slog.Int("request_count", v.count),
					slog.String("path", path),
				)
				return c.JSON(http.StatusTooManyRequests, map[string]string{
					"error": "Too many requests from this IP, please try again later.",
				})
			}

			// Log high request volume (early warning)
			if v.count > 50 && v.count%10 == 0 {
				logger.Info("high request volume",
					slog.String("ip", ip),
					slog.Int("request_count", v.count),
				)
			}

			limiter.mu.Unlock()
			return next(c)
		}
	}
}

// RequestID generates a unique request ID for each request
func RequestID() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			requestID := uuid.New().String()
			c.Set("request_id", requestID)

			// Also set in response header for client-side tracing
			c.Response().Header().Set("X-Request-ID", requestID)

			return next(c)
		}
	}
}

// StructuredLogger replaces Echo's basic logger with structured logging
func StructuredLogger(logger *slog.Logger) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			start := time.Now()

			// Get request ID (should be set by RequestID middleware)
			requestID := "unknown"
			if id, ok := c.Get("request_id").(string); ok {
				requestID = id
			}

			// Call next handler
			err := next(c)

			// Log after request completes
			req := c.Request()
			res := c.Response()
			latency := time.Since(start)

			fields := []any{
				slog.String("request_id", requestID),
				slog.String("method", req.Method),
				slog.String("path", req.URL.Path),
				slog.String("remote_ip", c.RealIP()),
				slog.Int("status", res.Status),
				slog.Int64("bytes_out", res.Size),
				slog.Duration("latency", latency),
				slog.String("user_agent", req.UserAgent()),
			}

			// Log query params for API endpoints (for debugging)
			if len(req.URL.RawQuery) > 0 {
				fields = append(fields, slog.String("query", req.URL.RawQuery))
			}

			// Log level based on status code and errors
			if err != nil {
				fields = append(fields, slog.String("error", err.Error()))
				logger.Error("HTTP request failed", fields...)
			} else if res.Status >= 500 {
				logger.Error("HTTP server error", fields...)
			} else if res.Status >= 400 {
				logger.Warn("HTTP client error", fields...)
			} else {
				logger.Info("HTTP request", fields...)
			}

			return err
		}
	}
}
