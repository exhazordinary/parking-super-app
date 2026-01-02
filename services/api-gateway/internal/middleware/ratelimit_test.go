package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestRateLimiter_Limit(t *testing.T) {
	// Create a rate limiter with 3 requests per minute
	limiter := NewRateLimiter(3, time.Minute)

	handler := limiter.Limit(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	// First 3 requests should pass
	for i := 0; i < 3; i++ {
		req := httptest.NewRequest(http.MethodGet, "/test", nil)
		req.RemoteAddr = "192.168.1.1:12345"
		rec := httptest.NewRecorder()

		handler.ServeHTTP(rec, req)

		if rec.Code != http.StatusOK {
			t.Errorf("request %d: expected status 200, got %d", i+1, rec.Code)
		}
	}

	// 4th request should be rate limited
	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	req.RemoteAddr = "192.168.1.1:12345"
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusTooManyRequests {
		t.Errorf("expected status 429, got %d", rec.Code)
	}
}

func TestRateLimiter_DifferentIPs(t *testing.T) {
	limiter := NewRateLimiter(2, time.Minute)

	handler := limiter.Limit(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	// Request from IP 1
	req1 := httptest.NewRequest(http.MethodGet, "/test", nil)
	req1.RemoteAddr = "10.0.0.1:12345"
	rec1 := httptest.NewRecorder()
	handler.ServeHTTP(rec1, req1)

	if rec1.Code != http.StatusOK {
		t.Errorf("IP1 request 1: expected 200, got %d", rec1.Code)
	}

	// Request from IP 2 should have its own limit
	req2 := httptest.NewRequest(http.MethodGet, "/test", nil)
	req2.RemoteAddr = "10.0.0.2:12345"
	rec2 := httptest.NewRecorder()
	handler.ServeHTTP(rec2, req2)

	if rec2.Code != http.StatusOK {
		t.Errorf("IP2 request 1: expected 200, got %d", rec2.Code)
	}

	// Another request from IP 1 should still work (within limit)
	req3 := httptest.NewRequest(http.MethodGet, "/test", nil)
	req3.RemoteAddr = "10.0.0.1:12345"
	rec3 := httptest.NewRecorder()
	handler.ServeHTTP(rec3, req3)

	if rec3.Code != http.StatusOK {
		t.Errorf("IP1 request 2: expected 200, got %d", rec3.Code)
	}

	// IP 1 should now be rate limited
	req4 := httptest.NewRequest(http.MethodGet, "/test", nil)
	req4.RemoteAddr = "10.0.0.1:12345"
	rec4 := httptest.NewRecorder()
	handler.ServeHTTP(rec4, req4)

	if rec4.Code != http.StatusTooManyRequests {
		t.Errorf("IP1 request 3: expected 429, got %d", rec4.Code)
	}
}

func TestRateLimiter_WindowReset(t *testing.T) {
	// Use a very short window for testing
	limiter := NewRateLimiter(1, 50*time.Millisecond)

	handler := limiter.Limit(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	// First request passes
	req1 := httptest.NewRequest(http.MethodGet, "/test", nil)
	req1.RemoteAddr = "192.168.1.1:12345"
	rec1 := httptest.NewRecorder()
	handler.ServeHTTP(rec1, req1)

	if rec1.Code != http.StatusOK {
		t.Errorf("request 1: expected 200, got %d", rec1.Code)
	}

	// Second request is rate limited
	req2 := httptest.NewRequest(http.MethodGet, "/test", nil)
	req2.RemoteAddr = "192.168.1.1:12345"
	rec2 := httptest.NewRecorder()
	handler.ServeHTTP(rec2, req2)

	if rec2.Code != http.StatusTooManyRequests {
		t.Errorf("request 2: expected 429, got %d", rec2.Code)
	}

	// Wait for window to reset
	time.Sleep(60 * time.Millisecond)

	// Request should pass again
	req3 := httptest.NewRequest(http.MethodGet, "/test", nil)
	req3.RemoteAddr = "192.168.1.1:12345"
	rec3 := httptest.NewRecorder()
	handler.ServeHTTP(rec3, req3)

	if rec3.Code != http.StatusOK {
		t.Errorf("request 3 after reset: expected 200, got %d", rec3.Code)
	}
}
