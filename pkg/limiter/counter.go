package limiter

import (
	"fmt"
	"time"
)

type CounterStore interface {
	Set(key string, value int64, timeout time.Duration) error
	Increment(key string) (int64, error)
}

type LimitInfo struct {
	Limit     int64
	Remaining int64
	Reset     time.Time
}

type limitCounter struct {
	config Config
}

func newLimitCounter(options ...Option) *limitCounter {
	return &limitCounter{
		config: newConfig(options...),
	}
}

func (c *limitCounter) Increment(token, ip string) (LimitInfo, error) {
	ipKey := IpKey(ip)
	ipCount, err := c.config.Store.Increment(ipKey)
	if err != nil {
		return LimitInfo{}, err
	}

	tokenKey := TokenKey(token)
	tokenCount, err := c.config.Store.Increment(tokenKey)
	if err != nil {
		return LimitInfo{}, err
	}

	ipLimit, ok := c.config.Limits[ipKey]
	if !ok {
		ipLimit = c.config.DefaultLimit
	}

	tokenLimit, ok := c.config.Limits[tokenKey]
	if !ok {
		tokenLimit = c.config.DefaultLimit
	}

	tokenInfo := LimitInfo{
		Limit:     tokenLimit.RequestsPerSecond,
		Remaining: tokenLimit.RequestsPerSecond - tokenCount,
		Reset:     time.Now().UTC().Truncate(time.Second).Add(1 * time.Second),
	}
	ipInfo := LimitInfo{
		Limit:     ipLimit.RequestsPerSecond,
		Remaining: ipLimit.RequestsPerSecond - ipCount,
		Reset:     time.Now().UTC().Truncate(time.Second).Add(1 * time.Second),
	}

	if ipInfo.Remaining < 0 {
		c.config.Store.Set(ipKey, ipLimit.RequestsPerSecond+1, ipLimit.TimeBlocked)
		ipInfo.Reset = time.Now().UTC().Truncate(time.Second).Add(ipLimit.TimeBlocked)
	}
	if tokenInfo.Remaining < 0 {
		c.config.Store.Set(tokenKey, tokenLimit.RequestsPerSecond+1, tokenLimit.TimeBlocked)
		tokenInfo.Reset = time.Now().UTC().Truncate(time.Second).Add(tokenLimit.TimeBlocked)
	}

	if token != "" && tokenInfo.Remaining > ipInfo.Remaining {
		return tokenInfo, nil
	} else {
		return ipInfo, nil
	}
}

func IpKey(ip string) string {
	return fmt.Sprintf("IP:%s", ip)
}

func TokenKey(token string) string {
	return fmt.Sprintf("TK:%s", token)
}
