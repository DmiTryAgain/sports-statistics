package tg

import (
	"context"
	"sync"
	"time"
)

const (
	minRequestInterval  = 200 * time.Millisecond
	rateLimitCleanupTTL = 5 * time.Minute
)

type rateLimiter struct {
	mu       *sync.RWMutex
	lastSeen map[int64]time.Time
}

func newRateLimiter(ctx context.Context) *rateLimiter {
	rl := &rateLimiter{
		lastSeen: make(map[int64]time.Time),
		mu:       &sync.RWMutex{},
	}
	go rl.cleanupLoop(ctx)
	return rl
}

func (rl *rateLimiter) Allow(userID int64) bool {
	if rl.allow(userID) {
		return true
	}

	rl.mu.Lock()
	defer rl.mu.Unlock()

	rl.lastSeen[userID] = time.Now()
	return false
}

func (rl *rateLimiter) allow(userID int64) bool {
	rl.mu.RLock()
	defer rl.mu.RUnlock()
	last, ok := rl.lastSeen[userID]
	return !ok || time.Since(last) >= rateLimitCleanupTTL
}

func (rl *rateLimiter) cleanupLoop(ctx context.Context) {
	ticker := time.NewTicker(1 * time.Minute)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			rl.mu.Lock()
			now := time.Now()
			for k, v := range rl.lastSeen {
				if now.Sub(v) > rateLimitCleanupTTL {
					delete(rl.lastSeen, k)
				}
			}
			rl.mu.Unlock()
		}
	}
}
