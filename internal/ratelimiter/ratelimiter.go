package ratelimiter

import (
	"sync"
	"time"
)

// UserRateLimiter implements a per-user rate limiter using a token bucket algorithm
type UserRateLimiter struct {
	buckets map[string]*rateLimiter
	mutex   sync.Mutex
}

// NewUserRateLimiter creates a new UserRateLimiter
func NewUserRateLimiter() *UserRateLimiter {
	return &UserRateLimiter{
		buckets: make(map[string]*rateLimiter),
	}
}

// Allow checks if a request is allowed for a specific user
func (ul *UserRateLimiter) Allow(user string) bool {
	ul.mutex.Lock()
	defer ul.mutex.Unlock()

	// Get or create the rate limiter for the user
	limiter, exists := ul.buckets[user]
	if !exists {
		limiter = newRateLimiter(15, time.Second)
		ul.buckets[user] = limiter
	}

	// Check if the user is currently rate limited
	if limiter.isLimited() {
		return false
	}

	// Check if the request is allowed for the user
	if !limiter.allow() {
		// If the request is not allowed, mark the user as rate limited for 30 seconds
		limiter.setLimitDuration(30 * time.Second)
		return false
	}

	return true
}

// rateLimiter implements a simple token bucket rate limiter
type rateLimiter struct {
	mutex           sync.Mutex
	capacity        int           // Maximum tokens the bucket can hold
	tokens          int           // Current number of tokens in the bucket
	fillRate        time.Duration // Rate at which the bucket is refilled
	lastUpdated     time.Time     // Last time the bucket was updated
	limitExpiration time.Time     // Time when the rate limit expires
}

// newRateLimiter creates a new RateLimiter with given capacity and fill rate
func newRateLimiter(capacity int, fillRate time.Duration) *rateLimiter {
	return &rateLimiter{
		capacity:    capacity,
		tokens:      capacity,
		fillRate:    fillRate,
		lastUpdated: time.Now(),
	}
}

// allow checks if a request is allowed based on the rate limiter's settings
func (r *rateLimiter) allow() bool {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	now := time.Now()

	// Refill the bucket based on elapsed time
	elapsedTime := now.Sub(r.lastUpdated)
	refillAmount := int(elapsedTime / r.fillRate)
	r.tokens = min(r.capacity, r.tokens+refillAmount)

	// Update the last update time
	r.lastUpdated = now

	// If there are enough tokens, allow the request
	if r.tokens > 0 {
		r.tokens--
		return true
	}
	return false
}

// isLimited checks if the rate limiter is currently in a limited state
func (r *rateLimiter) isLimited() bool {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	return time.Now().Before(r.limitExpiration)
}

// setLimitDuration sets the duration for which the rate limiter should be limited
func (r *rateLimiter) setLimitDuration(duration time.Duration) {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	r.limitExpiration = time.Now().Add(duration)
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
