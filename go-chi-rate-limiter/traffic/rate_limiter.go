package traffic

import (
	"log"
	"sync"
	"time"

	"golang.org/x/time/rate"
)

// IPRateLimiter ...
type IPRateLimiter struct {
	limiter  *rate.Limiter
	lastSeen time.Time
}

// NewIPRateLimiter ...
func NewIPRateLimiter(r rate.Limit, b int) *IPRateLimiter {
	l := &IPRateLimiter{
		limiter:  rate.NewLimiter(r, b),
		lastSeen: time.Now(),
	}
	return l
}

// Allow tells if to allow the event (request to be processed in this case) to happen.
func (l *IPRateLimiter) Allow() bool {
	return l.limiter.Allow()
}

func (l *IPRateLimiter) RefreshLastSeen() {
	l.lastSeen = time.Now()
}

// IPRateControl ...
type IPRateControl struct {
	limiters map[string]*IPRateLimiter // rate limiters per IP
	mu       *sync.RWMutex
	r        rate.Limit
	b        int
}

// NewIPRateControl creates an instance of `IPRateControl` that allows
// up to r requests per second with an initial burst of `b` number of requests.
func NewIPRateControl(r rate.Limit, b int) *IPRateControl {
	i := &IPRateControl{
		limiters: make(map[string]*IPRateLimiter),
		mu:       &sync.RWMutex{},
		r:        r,
		b:        b,
	}
	return i
}

// addIPRateLimiter creates a new IPRateLimiter
// and keeps it internally for later use.
func (i *IPRateControl) addIPRateLimiter(ip string) *IPRateLimiter {
	limiter := NewIPRateLimiter(i.r, i.b)
	i.limiters[ip] = limiter
	return limiter
}

// GetLimiter returns the rate limiter for the provided IP address, if it exists.
// Otherwise, it calls adds the IP address to the map.
func (i *IPRateControl) GetLimiter(ip string) *IPRateLimiter {
	i.mu.Lock()
	defer i.mu.Unlock()
	limiter, exists := i.limiters[ip]
	if !exists {
		return i.addIPRateLimiter(ip)
	}
	limiter.RefreshLastSeen()
	return limiter
}

// PeriodicCleanup removes old entries of IPRateLimiter
// to keep the in-memory storage somehow in control.
func (i *IPRateControl) PeriodicCleanup() {
	for {
		time.Sleep(2 * time.Minute)
		i.mu.Lock()

		cs := 0
		for ip, limiter := range i.limiters {
			if time.Since(limiter.lastSeen) > 2*time.Minute {
				delete(i.limiters, ip)
				cs++
			}
		}
		if cs > 0 {
			log.Println(">>> Entries cleaned up:", cs)
		}
		i.mu.Unlock()
	}
}
