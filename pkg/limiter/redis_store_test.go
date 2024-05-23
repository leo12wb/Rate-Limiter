package limiter

import (
	"testing"
	"time"

	"github.com/leo12wb/Rate-Limiter/internal/testutils"
	"github.com/matryer/is"
	"github.com/redis/go-redis/v9"
)

func TestRedisStore(t *testing.T) {
	redisContainer := testutils.CreateRedisTestingContainer(t)
	client := redis.NewClient(&redis.Options{
		Addr: redisContainer.Endpoint,
	})

	t.Run("Increment should create entry on non existing entry", func(t *testing.T) {
		is := is.New(t)
		sut := NewRedisStore(client)
		key := "non-existing-key"

		got, err := sut.Increment(key)

		is.NoErr(err)
		is.Equal(got, int64(1))
	})

	t.Run("Increment should increment value by 1 on existing entry", func(t *testing.T) {
		is := is.New(t)
		sut := NewRedisStore(client)
		key := "some-specific-key"

		got, err := sut.Increment(key)
		is.NoErr(err)
		is.Equal(got, int64(1))

		got, err = sut.Increment(key)
		is.NoErr(err)
		is.Equal(got, int64(2))

		got, err = sut.Increment(key)
		is.NoErr(err)
		is.Equal(got, int64(3))
	})

	t.Run("Increment should restart on increments after one second", func(t *testing.T) {
		is := is.New(t)
		sut := NewRedisStore(client)
		key := "another-specific-key"

		got, err := sut.Increment(key)
		is.NoErr(err)
		is.Equal(got, int64(1))

		time.Sleep(time.Second)

		got, err = sut.Increment(key)
		is.NoErr(err)
		is.Equal(got, int64(1))
	})

	t.Run("Set should replace value on existing key", func(t *testing.T) {
		is := is.New(t)
		sut := NewRedisStore(client)
		key := "yet-another-key"

		got, err := sut.Increment(key)
		is.Equal(got, int64(1))
		is.NoErr(err)

		err = sut.Set(key, 5, time.Second*2)
		is.NoErr(err)

		got, err = sut.Increment(key)
		is.Equal(got, int64(6))
		is.NoErr(err)
	})

	t.Run("Set should now allow the value to reset on given duration", func(t *testing.T) {
		is := is.New(t)
		sut := NewRedisStore(client)
		key := "some-key"

		err := sut.Set(key, 6, time.Second*2)
		is.NoErr(err)

		got, err := sut.Increment(key)
		is.Equal(got, int64(7))
		is.NoErr(err)

		time.Sleep(1 * time.Second)
		got, err = sut.Increment(key)
		is.Equal(got, int64(8))
		is.NoErr(err)

		time.Sleep(1 * time.Second)
		got, err = sut.Increment(key)
		is.Equal(got, int64(1))
		is.NoErr(err)
	})
}
