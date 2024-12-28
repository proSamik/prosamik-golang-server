package middleware

import (
	"net/http"
	"strings"
	"sync"
	"time"
)

type visitor struct {
	count      int
	lastAccess time.Time
}

type RateLimiter struct {
	visitors map[string]*visitor
	mu       sync.Mutex
	rate     int           // Number of allowed requests
	window   time.Duration // Time window to track
}

func NewRateLimiter(rate int, window time.Duration) *RateLimiter {
	return &RateLimiter{
		visitors: make(map[string]*visitor),
		rate:     rate,
		window:   window,
	}
}

func getIP(r *http.Request) string {
	// First try X-Forwarded-For header
	forwarded := r.Header.Get("X-Forwarded-For")
	if forwarded != "" {
		// Take the first IP if there are multiple
		return strings.Split(forwarded, ",")[0]
	}

	// Fallback to RemoteAddr
	return r.RemoteAddr
}

func (rl *RateLimiter) RateLimitMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Get IP address from request
		ip := getIP(r)

		rl.mu.Lock()
		v, exists := rl.visitors[ip]
		if !exists {
			// First visit, create new visitor
			rl.visitors[ip] = &visitor{
				count:      1,
				lastAccess: time.Now(),
			}
			rl.mu.Unlock()
			next(w, r)
			return
		}

		// Check if window has expired
		if time.Since(v.lastAccess) > rl.window {
			// Reset counter
			v.count = 1
			v.lastAccess = time.Now()
		} else if v.count >= rl.rate {
			// Too many requests
			rl.mu.Unlock()
			http.Error(w, "Too many requests", http.StatusTooManyRequests)
			return
		} else {
			// Increment counter
			v.count++
		}

		v.lastAccess = time.Now()
		rl.mu.Unlock()

		next(w, r)
	}
}
