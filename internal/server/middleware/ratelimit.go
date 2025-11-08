package middleware

import (
    "sync"
    "time"

    "github.com/gin-gonic/gin"
)

type clientInfo struct {
    count     int
    expiresAt time.Time
}

type rateLimiter struct {
    mu      sync.Mutex
    clients map[string]*clientInfo
    max     int
    window  time.Duration
}

// NewRateLimiter returns a middleware limiting requests per client IP.
func NewRateLimiter(max int, window time.Duration) gin.HandlerFunc {
    rl := &rateLimiter{
        clients: make(map[string]*clientInfo),
        max:     max,
        window:  window,
    }

    // background cleanup
    go func() {
        ticker := time.NewTicker(window)
        defer ticker.Stop()
        for range ticker.C {
            rl.mu.Lock()
            now := time.Now()
            for k, v := range rl.clients {
                if v.expiresAt.Before(now) {
                    delete(rl.clients, k)
                }
            }
            rl.mu.Unlock()
        }
    }()

    return func(c *gin.Context) {
        key := c.ClientIP() + ":" + c.FullPath()
        rl.mu.Lock()
        info, ok := rl.clients[key]
        if !ok || time.Now().After(info.expiresAt) {
            rl.clients[key] = &clientInfo{count: 1, expiresAt: time.Now().Add(rl.window)}
            rl.mu.Unlock()
            c.Next()
            return
        }

        if info.count >= rl.max {
            rl.mu.Unlock()
            c.AbortWithStatusJSON(429, gin.H{"error": "rate limit exceeded"})
            return
        }

        info.count++
        rl.mu.Unlock()
        c.Next()
    }
}
