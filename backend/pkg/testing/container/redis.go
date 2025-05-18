package container

import (
	"context"
	"net/url"
	"testing"

	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	tcredis "github.com/testcontainers/testcontainers-go/modules/redis"
)

type redisContainerInput struct {
	image string
}

type RedisContainerOption func(*redisContainerInput)

func RedisContainerInput(options ...RedisContainerOption) *redisContainerInput {
	input := &redisContainerInput{
		image: "redis:latest",
	}
	for _, o := range options {
		o(input)
	}
	return input
}

func WithRedisImage(image string) RedisContainerOption {
	return func(input *redisContainerInput) {
		input.image = image
	}
}

type RedisContainer struct {
	testcontainers.Container
	dsn string
}

func NewRedisContainer(t *testing.T, ctx context.Context, input *redisContainerInput) (*redis.Client, func()) {
	t.Helper()

	redisContainer, err := tcredis.Run(ctx,
		input.image,
	)
	require.NoError(t, err)

	connectionString, err := redisContainer.ConnectionString(ctx)
	require.NoError(t, err)

	// redis://localhost:6379 のような形式なので、 localhost:6379 に変換
	u, err := url.Parse(connectionString)
	require.NoError(t, err)
	addr := u.Host

	client := redis.NewClient(&redis.Options{
		Addr: addr,
	})

	return client, func() {
		if err := redisContainer.Terminate(ctx); err != nil {
			t.Fatalf("failed to terminate container: %v", err)
		}
	}
}

func (c *RedisContainer) GetDSN() string {
	return c.dsn
}

func (c *RedisContainer) Terminate(ctx context.Context) error {
	return c.Container.Terminate(ctx)
}
