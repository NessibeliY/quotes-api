package handlers

import (
	"net/http"
	"sync"
	"time"
)

type RateLimiter struct {
	visitors map[string]*Visitor
	mu       sync.Mutex
	rate     int           // запросов
	window   time.Duration // за сколько времени
}

type Visitor struct {
	lastSeen time.Time
	tokens   int
}

func NewRateLimiter(rate int, window time.Duration) *RateLimiter {
	rl := &RateLimiter{
		visitors: make(map[string]*Visitor),
		rate:     rate,
		window:   window,
	}

	go rl.cleanupVisitors() // каждую минуту чистим устаревших юзеров
	return rl
}

// cleanupVisitors удаляет устаревших юзеров с мапы visitors
func (rl *RateLimiter) cleanupVisitors() {
	for {
		time.Sleep(time.Minute)
		rl.cleanupOnce()
	}
}

func (rl *RateLimiter) cleanupOnce() {
	rl.mu.Lock()
	for ip, v := range rl.visitors {
		if time.Since(v.lastSeen) > rl.window {
			delete(rl.visitors, ip)
		}
	}
	rl.mu.Unlock()
}

func (rl *RateLimiter) LimitMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ip := r.RemoteAddr
		rl.mu.Lock()
		v, ok := rl.visitors[ip]
		if !ok {
			v = &Visitor{
				lastSeen: time.Now(),
				tokens:   rl.rate - 1,
			}
			rl.visitors[ip] = v
			rl.mu.Unlock()
			next.ServeHTTP(w, r)
			return
		}

		if time.Since(v.lastSeen) > rl.window {
			v.tokens = rl.rate - 1
			v.lastSeen = time.Now()
			rl.mu.Unlock()
			next.ServeHTTP(w, r)
			return
		}

		v.lastSeen = time.Now()
		if v.tokens > 0 {
			v.tokens--
			rl.mu.Unlock()
			next.ServeHTTP(w, r)
			return
		}

		rl.mu.Unlock()
		http.Error(w, "too many requests", http.StatusTooManyRequests)
	})
}
