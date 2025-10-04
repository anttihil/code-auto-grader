package middleware

import (
	"net/http"
	"strings"
	"sync"
	"time"

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

func RateLimit() echo.MiddlewareFunc {
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
				return c.JSON(http.StatusTooManyRequests, map[string]string{
					"error": "Too many requests from this IP, please try again later.",
				})
			}

			limiter.mu.Unlock()
			return next(c)
		}
	}
}
