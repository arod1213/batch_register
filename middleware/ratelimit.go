package middleware

import (
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

type fixedWindowCounter struct {
	windowStart time.Time
	count       int
}

// RateLimitSongs limits GET /songs requests per client IP using a fixed-window counter.
// This protects the API (and your terminal) from accidental client-side infinite polling loops.
func RateLimitSongs(maxRequests int, window time.Duration) gin.HandlerFunc {
	var mu sync.Mutex
	counters := map[string]*fixedWindowCounter{}

	return func(c *gin.Context) {
		// NOTE: In Gin, FullPath() may be empty in early middleware stages.
		// Use the raw URL path so this always applies.
		if c.Request.Method != http.MethodGet || c.Request.URL.Path != "/songs" {
			c.Next()
			return
		}

		ip := c.ClientIP()
		now := time.Now()

		mu.Lock()
		ctr, ok := counters[ip]
		if !ok {
			ctr = &fixedWindowCounter{windowStart: now, count: 0}
			counters[ip] = ctr
		}

		// reset window
		if now.Sub(ctr.windowStart) >= window {
			ctr.windowStart = now
			ctr.count = 0
		}

		ctr.count++
		count := ctr.count
		mu.Unlock()

		if count > maxRequests {
			c.Header("Retry-After", "1")
			c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{
				"error": "rate limited: too many /songs requests (check client polling loop)",
			})
			return
		}

		c.Next()
	}
}
