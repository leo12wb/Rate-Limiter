package limiter

import (
	"testing"
	"time"

	"github.com/matryer/is"
)

func TestCounter(t *testing.T) {
	t.Run("Increment should use default limit on non specified IP", func(t *testing.T) {
		is := is.New(t)
		ip := "127.0.0.1"
		defaultLimit := int64(101)
		sut := newLimitCounter(
			WithDefaultLimit(defaultLimit, time.Second),
		)

		info, err := sut.Increment("", ip)
		is.NoErr(err)
		is.Equal(info.Limit, defaultLimit)
		is.Equal(info.Remaining, defaultLimit-1)
		is.Equal(info.Reset, time.Now().UTC().Add(time.Second).Truncate(time.Second))

		info, err = sut.Increment("", ip)
		is.NoErr(err)
		is.Equal(info.Limit, defaultLimit)
		is.Equal(info.Remaining, defaultLimit-2)
		is.Equal(info.Reset, time.Now().UTC().Add(time.Second).Truncate(time.Second))

		info, err = sut.Increment("", ip)
		is.NoErr(err)
		is.Equal(info.Limit, defaultLimit)
		is.Equal(info.Remaining, defaultLimit-3)
		is.Equal(info.Reset, time.Now().UTC().Add(time.Second).Truncate(time.Second))
	})

	t.Run("Increment should use provided limit on matching IP", func(t *testing.T) {
		is := is.New(t)
		ip := "192.160.1.123"
		requestsPerSecond := int64(42)
		sut := newLimitCounter(
			WithDefaultLimit(100, time.Minute),
			WithIPLimits(map[string]Limit{
				ip: {RequestsPerSecond: requestsPerSecond, TimeBlocked: 3 * time.Second},
			}),
		)

		info, err := sut.Increment("", ip)
		is.NoErr(err)
		is.Equal(info.Limit, requestsPerSecond)
		is.Equal(info.Remaining, requestsPerSecond-1)
		is.Equal(info.Reset, time.Now().UTC().Add(time.Second).Truncate(time.Second))

		info, err = sut.Increment("", ip)
		is.NoErr(err)
		is.Equal(info.Limit, requestsPerSecond)
		is.Equal(info.Remaining, requestsPerSecond-2)
		is.Equal(info.Reset, time.Now().UTC().Add(time.Second).Truncate(time.Second))

		info, err = sut.Increment("", ip)
		is.NoErr(err)
		is.Equal(info.Limit, requestsPerSecond)
		is.Equal(info.Remaining, requestsPerSecond-3)
		is.Equal(info.Reset, time.Now().UTC().Add(time.Second).Truncate(time.Second))
	})

	t.Run("Increment should use provided limit on matching Token", func(t *testing.T) {
		is := is.New(t)
		ip := "192.160.1.124"
		token := "SOME_TOKEN"
		defaultLimit := int64(123)
		sut := newLimitCounter(
			WithDefaultLimit(10, time.Minute),
			WithTokenLimits(map[string]Limit{
				token: {RequestsPerSecond: defaultLimit, TimeBlocked: 3 * time.Second},
			}),
		)

		info, err := sut.Increment(token, ip)
		is.NoErr(err)
		is.Equal(info.Limit, defaultLimit)
		is.Equal(info.Remaining, defaultLimit-1)
		is.Equal(info.Reset, time.Now().UTC().Add(time.Second).Truncate(time.Second))

		info, err = sut.Increment(token, ip)
		is.NoErr(err)
		is.Equal(info.Limit, defaultLimit)
		is.Equal(info.Remaining, defaultLimit-2)
		is.Equal(info.Reset, time.Now().UTC().Add(time.Second).Truncate(time.Second))

		info, err = sut.Increment(token, ip)
		is.NoErr(err)
		is.Equal(info.Limit, defaultLimit)
		is.Equal(info.Remaining, defaultLimit-3)
		is.Equal(info.Reset, time.Now().UTC().Add(time.Second).Truncate(time.Second))
	})
}
