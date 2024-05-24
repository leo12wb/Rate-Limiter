package rate_limit

import (
	"errors"

	"github.com/leo12wb/Rate-Limiter/internal/entity/web_session"
)

type RateLimitRepository interface {
	SetRequestCounter(session *web_session.WebSession) error
	GetLastRequestTime(session *web_session.WebSession) (int64, error)
	DecreaseTokenBucket(session *web_session.WebSession) (bool, error)
}

type ThrottledError struct {
}

func (te *ThrottledError) ThrottledError() error {
	return errors.New("you have reached the maximum number of requests or actions allowed within a certain time frame")
}
