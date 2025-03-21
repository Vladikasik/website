package middleware

import (
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
)

// SetupMiddleware configures all middleware for the application
func SetupMiddleware(handler http.Handler, trustedOrigins []string, logger *log.Logger) http.Handler {
	return securityHeaders(
		cors.Handler(corsConfig(trustedOrigins))(
			rateLimiter(
				middleware.RequestID(
					middleware.RealIP(
						middleware.Logger(
							middleware.Recoverer(handler),
						),
					),
				),
			),
		),
	)
}

// securityHeaders adds security headers to all responses
func securityHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.Header().Set("X-Frame-Options", "deny")
		w.Header().Set("X-XSS-Protection", "1; mode=block")
		w.Header().Set("Referrer-Policy", "same-origin")
		w.Header().Set("Content-Security-Policy", "default-src 'self'; frame-ancestors 'none'")

		next.ServeHTTP(w, r)
	})
}

// rateLimiter provides basic IP-based rate limiting
func rateLimiter(next http.Handler) http.Handler {
	// Simple in-memory store for rate limiting
	// In production, use a distributed cache like Redis
	type visitor struct {
		count      int
		lastAccess time.Time
	}
	
	visitors := make(map[string]*visitor)
	
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ip := r.Header.Get("X-Real-IP")
		if ip == "" {
			ip = r.RemoteAddr
		}
		
		// Clean up IP if needed (remove port)
		if idx := strings.LastIndex(ip, ":"); idx != -1 {
			ip = ip[:idx]
		}
		
		// Get current visitor or create new one
		v, exists := visitors[ip]
		now := time.Now()
		
		if !exists {
			visitors[ip] = &visitor{
				count:      1,
				lastAccess: now,
			}
		} else {
			// Reset if last access was more than a minute ago
			if now.Sub(v.lastAccess) > time.Minute {
				v.count = 0
				v.lastAccess = now
			}
			
			// Increment counter
			v.count++
			
			// If too many requests (more than 60 per minute)
			if v.count > 60 {
				w.Header().Set("Retry-After", "60")
				http.Error(w, "Rate limit exceeded", http.StatusTooManyRequests)
				return
			}
		}
		
		next.ServeHTTP(w, r)
	})
}

// corsConfig creates a CORS configuration with the trusted origins
func corsConfig(trustedOrigins []string) cors.Options {
	// If no trusted origins are provided, add defaults including localhost for dev
	if len(trustedOrigins) == 0 {
		trustedOrigins = []string{
			"https://aynshteyn.dev",
			"http://localhost:3000",
			"http://127.0.0.1:3000",
		}
	}
	
	return cors.Options{
		AllowedOrigins:   trustedOrigins,
		AllowedMethods:   []string{"GET", "POST", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300, // Maximum value not refused by any of the browsers
	}
} 