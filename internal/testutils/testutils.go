package testutils

import (
	"context"
	"testing"

	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

type redisContainer struct {
	Container testcontainers.Container
	Endpoint  string
}

func CreateRedisTestingContainer(t testing.TB) redisContainer {
	ctx := context.Background()
	req := testcontainers.ContainerRequest{
		Image:        "redis:7.2.4-alpine",
		ExposedPorts: []string{"6379/tcp"},
		WaitingFor:   wait.ForLog("Ready to accept connections"),
	}
	container, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		t.Fatalf("Could not start redis: %s", err)
	}

	endpoint, err := container.Endpoint(ctx, "")
	if err != nil {
		t.Fatal(err)
	}

	t.Cleanup(func() {
		if err := container.Terminate(ctx); err != nil {
			t.Fatalf("failed to terminate redis container: %s", err)
		}
	})

	return redisContainer{
		Container: container,
		Endpoint:  endpoint,
	}
}
