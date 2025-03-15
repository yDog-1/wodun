package repository_test

import (
	"context"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	rediscon "github.com/testcontainers/testcontainers-go/modules/redis"
	"github.com/yDog-1/wodun/backend/repository"
)

type storeClaims struct {
	jwt.RegisteredClaims
	Jti string `json:"jti"`
}

func (c storeClaims) GetJTI() (string, error) {
	return c.Jti, nil
}

func Test_JTIを保存する(t *testing.T) {
	ctx := context.Background()
	con, err := rediscon.Run(ctx, "redis:7")
	defer func() {
		if err := testcontainers.TerminateContainer(con); err != nil {
			t.Fatal("failed to terminate container: ", err)
		}
	}()
	require.NoError(t, err)
	addr, err := con.ConnectionString(ctx)
	require.NoError(t, err)
	opts, err := redis.ParseURL(addr)
	require.NoError(t, err)
	store := redis.NewClient(opts)

	sut := repository.NewTokenRepository(store)
	sub := "test"
	exp := time.Now().Add(time.Hour)
	jti := "123"
	claims := storeClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   sub,
			ExpiresAt: jwt.NewNumericDate(exp),
		},
		Jti: jti,
	}

	err = sut.SaveJTI(ctx, claims)
	assert.NoError(t, err, "failed to save JTI")
}
