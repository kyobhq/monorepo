package middlewares

import (
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

type LimiterConfig struct {
	MaxRequests uint
	Window      time.Duration
}

type ipLimiter struct {
	mu      sync.Mutex
	buckets map[string]*bucket
}

type bucket struct {
	count uint
	exp   time.Time
}

func RateLimiter(cfg LimiterConfig) gin.HandlerFunc {
	l := &ipLimiter{
		buckets: make(map[string]*bucket),
	}

	go func() {
		t := time.NewTicker(cfg.Window)
		for range t.C {
			l.gc()
		}
	}()

	return func(c *gin.Context) {
		key := c.RemoteIP()
		allowed := l.hit(key, cfg.MaxRequests, cfg.Window)

		if !allowed {
			c.Header("Retry-After", cfg.Window.String())
			c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{
				"error":       "rate limit exceeded",
				"limit":       cfg.MaxRequests,
				"time_window": cfg.Window.String(),
			})
			return
		}

		c.Next()
	}
}

func (l *ipLimiter) hit(key string, max uint, window time.Duration) bool {
	now := time.Now()

	l.mu.Lock()
	defer l.mu.Unlock()

	b, ok := l.buckets[key]
	if !ok || now.After(b.exp) {
		l.buckets[key] = &bucket{
			count: 1,
			exp:   now.Add(window),
		}
		return true
	}

	if b.count > max {
		return false
	}

	b.count++
	return true
}

func (l *ipLimiter) gc() {
	now := time.Now()

	l.mu.Lock()

	for k, b := range l.buckets {
		if now.After(b.exp) {
			delete(l.buckets, k)
		}
	}

	l.mu.Unlock()
}
