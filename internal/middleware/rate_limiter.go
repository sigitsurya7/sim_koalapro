package middleware

import (
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
)

type ipLimiter struct {
	limiters map[string]*rate.Limiter
	mu       sync.Mutex
	r        rate.Limit
	b        int
	lastSeen map[string]time.Time
	window   time.Duration
}

func newIPLimiter(r rate.Limit, b int, window time.Duration) *ipLimiter {
	return &ipLimiter{
		limiters: make(map[string]*rate.Limiter),
		r:        r,
		b:        b,
		lastSeen: make(map[string]time.Time),
		window:   window,
	}
}

func (l *ipLimiter) getLimiter(ip string) *rate.Limiter {
	l.mu.Lock()
	defer l.mu.Unlock()

	limiter, exists := l.limiters[ip]
	if !exists {
		limiter = rate.NewLimiter(l.r, l.b)
		l.limiters[ip] = limiter
	}
	l.lastSeen[ip] = time.Now()

	return limiter
}

func (l *ipLimiter) cleanup() {
	l.mu.Lock()
	defer l.mu.Unlock()

	cutoff := time.Now().Add(-l.window)
	for ip, ts := range l.lastSeen {
		if ts.Before(cutoff) {
			delete(l.lastSeen, ip)
			delete(l.limiters, ip)
		}
	}
}

func RateLimit(perSecond float64, burst int, cleanupWindow time.Duration) gin.HandlerFunc {
	limiter := newIPLimiter(rate.Limit(perSecond), burst, cleanupWindow)

	go func() {
		ticker := time.NewTicker(cleanupWindow)
		defer ticker.Stop()
		for range ticker.C {
			limiter.cleanup()
		}
	}()

	return func(c *gin.Context) {
		ip := c.ClientIP()
		if !limiter.getLimiter(ip).Allow() {
			c.AbortWithStatusJSON(429, gin.H{"error": "rate_limited"})
			return
		}
		c.Next()
	}
}
