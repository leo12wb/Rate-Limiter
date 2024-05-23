package entity

import (
	"context"
	"encoding/json"
	"time"
)

type RateLimiterInfo struct {
	Key       string
	Requests  int
	Every     int
	Remaining int
	Reset     int64
}

func (l RateLimiterInfo) MarshalBinary() ([]byte, error) {
	return json.Marshal(l)
}

type DatabaseRepository interface {
	Create(context.Context, RateLimiterInfo, time.Duration) error
	Read(context.Context, string) (*RateLimiterInfo, error)
	CheckLimit(context.Context, string, int, time.Duration) (bool, error)
}
