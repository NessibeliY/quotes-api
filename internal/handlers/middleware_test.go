package handlers

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestRateLimiter_AllowRequestsWithinLimit(t *testing.T) {
	rl := NewRateLimiter(2, time.Minute)

	var called int
	testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		called++
		w.WriteHeader(http.StatusOK)
	})

	wrapped := rl.LimitMiddleware(testHandler)

	req := httptest.NewRequest("GET", "/", nil)
	req.RemoteAddr = "127.0.0.1:1234"

	for i := 0; i < 2; i++ {
		w := httptest.NewRecorder()
		wrapped.ServeHTTP(w, req)
		if w.Code != http.StatusOK {
			t.Errorf("expected status code 200, got %d", w.Code)
		}
	}

	w := httptest.NewRecorder()
	wrapped.ServeHTTP(w, req)
	if w.Code != http.StatusTooManyRequests {
		t.Errorf("expected status code 429, got %d", w.Code)
	}
	if called != 2 {
		t.Errorf("expected 2 calls, got %d", called)
	}
}

func TestRateLimiter_ResetsTokensAfterWindow(t *testing.T) {
	rl := NewRateLimiter(1, time.Millisecond*10)

	req := httptest.NewRequest("GET", "/", nil)
	req.RemoteAddr = "127.0.0.1:1234"
	testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})
	wrapped := rl.LimitMiddleware(testHandler)

	w := httptest.NewRecorder()
	wrapped.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Errorf("expected status code 200, got %d", w.Code)
	}

	w = httptest.NewRecorder()
	wrapped.ServeHTTP(w, req)
	if w.Code != http.StatusTooManyRequests {
		t.Errorf("expected status code 429, got %d", w.Code)
	}

	time.Sleep(time.Millisecond * 20)

	w = httptest.NewRecorder()
	wrapped.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Errorf("expected status code 200 after window, got %d", w.Code)
	}
}

func TestRateLimiter_CleanupVisitorsRemovesOldEntries(t *testing.T) {
	rl := NewRateLimiter(1, time.Millisecond*10)

	// добавляем старого посетителя
	rl.visitors["127.0.0.1:1234"] = &Visitor{
		lastSeen: time.Now().Add(-20 * time.Millisecond),
		tokens:   0,
	}

	// добавляем нового посетителя
	rl.visitors["127.0.0.1:5678"] = &Visitor{
		lastSeen: time.Now(),
		tokens:   1,
	}

	rl.cleanupOnce()

	if len(rl.visitors) != 1 {
		t.Fatalf("expected 1 visitor left, got %d", len(rl.visitors))
	}
	if _, ok := rl.visitors["127.0.0.1:5678"]; !ok {
		t.Fatalf("expected visitor with ip 127.0.0.1:5678, got %v", rl.visitors)
	}
}
