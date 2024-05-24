package web_session

import (
	"crypto/sha256"
	"errors"
	"fmt"

	"github.com/leo12wb/Rate-Limiter/internal/value_objects"
)

type WebSession struct {
	IP                  string
	ApiToken            string
	maxRequestPerSecond uint
	ExpireSeconds       uint
}

const CounterSuffix = "_counter"
const TimerSuffix = "_lastest_request"

func NewWebSession(IP string, ApiToken string, requestLimits value_objects.RequestLimits) (WebSession, error) {
	res := WebSession{
		IP:            IP,
		ApiToken:      ApiToken,
		ExpireSeconds: requestLimits.ExpireSeconds,
	}
	err := res.IsValid()
	if err != nil {
		return WebSession{}, err
	}
	if res.ApiToken != "" {
		res.maxRequestPerSecond = requestLimits.IPLimit
	} else {
		res.maxRequestPerSecond = requestLimits.APILimit
	}
	return res, nil
}

func (h *WebSession) IsValid() error {
	if h.IP == "" {
		return errors.New("invalid IP address")
	}
	return nil
}

func (h *WebSession) GetSessionId() string {
	//The API Token precedes the IP address.
	if h.ApiToken != "" {
		return fmt.Sprintf("%x", sha256.Sum256([]byte(h.ApiToken)))
	}
	return fmt.Sprintf("%x", sha256.Sum256([]byte(h.IP)))

}

func (h *WebSession) GetRequestCounterId() string {
	return h.GetSessionId() + CounterSuffix
}
func (h *WebSession) GetRequestTimerId() string {
	return h.GetSessionId() + TimerSuffix
}

func (h *WebSession) GetRequestsLimitInSeconds() int64 {
	return int64(h.maxRequestPerSecond)
}
func (h *WebSession) GetExpireInSeconds() int64 {
	return int64(h.ExpireSeconds)
}
