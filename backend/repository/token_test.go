package repository_test

import (
	"context"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/yDog-1/wodun/backend/pkg/testing/container"
	"github.com/yDog-1/wodun/backend/repository"
)

type storeClaims struct {
	jwt.RegisteredClaims
	Jti string `json:"jti"`
}

func (c storeClaims) GetJTI() (string, error) {
	return c.Jti, nil
}

func TestTokenRepository_SaveJTIAndExistsJTI(t *testing.T) {
	ctx := context.Background()

	// Redisコンテナを起動
	client, terminate := container.NewRedisContainer(t, ctx, container.RedisContainerInput(
		container.WithRedisImage("redis:8-alpine"),
	))
	defer terminate()

	// TokenRepositoryを生成
	repo := repository.NewTokenRepository(client)

	// テストデータ
	userID := "user123"
	jti := uuid.New().String()
	expirationTime := time.Now().Add(time.Minute) // 1分後に期限切れ

	// storeClaimsを作成
	claims := storeClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   userID,
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
		Jti: jti,
	}

	// SaveJTIのテスト
	err := repo.SaveJTI(ctx, claims)
	assert.NoError(t, err, "SaveJTI should not return an error")

	// ExistsJTIのテスト (存在するJTI)
	exists, err := repo.ExistsJTI(ctx, userID, jti)
	assert.NoError(t, err, "ExistsJTI should not return an error for existing JTI")
	assert.True(t, exists, "ExistsJTI should return true for existing JTI")

	// ExistsJTIのテスト (存在しないJTI)
	nonExistentJTI := uuid.New().String()
	exists, err = repo.ExistsJTI(ctx, userID, nonExistentJTI)
	assert.NoError(t, err, "ExistsJTI should not return an error for non-existent JTI")
	assert.False(t, exists, "ExistsJTI should return false for non-existent JTI")

	// ExistsJTIのテスト (異なるユーザーID)
	otherUserID := "user456"
	exists, err = repo.ExistsJTI(ctx, otherUserID, jti)
	assert.NoError(t, err, "ExistsJTI should not return an error for different user ID")
	assert.False(t, exists, "ExistsJTI should return false for different user ID")

	// JTIが期限切れになった後のテスト (短い期限で保存し直す)
	shortExpirationTime := time.Now().Add(time.Second) // 1秒後に期限切れ
	shortClaims := storeClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   userID,
			ExpiresAt: jwt.NewNumericDate(shortExpirationTime),
		},
		Jti: jti,
	}
	err = repo.SaveJTI(ctx, shortClaims)
	assert.NoError(t, err, "SaveJTI should not return an error when saving with short expiration")

	time.Sleep(2 * time.Second) // 期限切れまで待つ

	exists, err = repo.ExistsJTI(ctx, userID, jti)
	assert.NoError(t, err, "ExistsJTI should not return an error after expiration")
	assert.False(t, exists, "ExistsJTI should return false after expiration")
}
