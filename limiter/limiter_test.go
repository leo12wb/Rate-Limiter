package limiter

import (
	"testing"
	"time"
)

func TestRateLimiter(t *testing.T) {
	rateLimiter, err := NewRateLimiter()
	if err != nil {
		t.Fatalf("Failed to create rate limiter: %v", err)
	}

	tests := []struct {
		identifier string
		maxRequests int
		window time.Duration
		expect bool
	}{
		{"127.0.0.1", 5, time.Second, true},
		{"127.0.0.1", 5, time.Second, false},
	}

	for _, tt := range tests {
		allowed, err := rateLimiter.Allow(tt.identifier, tt.maxRequests, tt.window)
		if err != nil {
			t.Errorf("Error in Allow method: %v", err)
		}
		if allowed != tt.expect {
			t.Errorf("Expected %v, got %v", tt.expect, allowed)
		}
	}
}