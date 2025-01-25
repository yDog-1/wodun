package auth_test

import (
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
	"github.com/yDog-1/wodun/backend/pkg/auth"
)

type mockTokenStore struct {
}

func (m mockTokenStore) SaveJTI(jti string) error {
	return nil
}
func (m mockTokenStore) ExistsJTI(jti string) (bool, error) {
	return false, nil
}

type mockClock struct{}

// 2100-01-01 00:00:00 を返す
func (c mockClock) Now() time.Time {
	return time.Date(2100, 1, 1, 0, 0, 0, 0, time.UTC)
}

func Testトークンを生成する(t *testing.T) {
	t.Setenv("TOKEN_ISSUER", "test")
	t.Setenv("TOKEN_AUDIENCE", "test")
	t.Setenv("TOKEN_ACCESS_SECRET", "test")
	t.Setenv("TOKEN_REFRESH_SECRET", "test")
	id := "1"
	var store mockTokenStore
	ts, err := auth.NewTokenService(store, mockClock{})
	assert.Nil(t, err)
	// トークンを生成する
	accessToken, refreshToken, err := ts.GenerateToken(id, "user")
	assert.NoError(t, err)

	accessClaims, err := ts.ParseAccessToken(accessToken)
	assert.NoError(t, err)
	assert.Equal(t,
		mockClock{}.Now().Add(time.Hour),
		accessClaims.Exp,
		"expが1時間後であること",
	)
	assert.Equal(t,
		mockClock{}.Now(),
		accessClaims.Iat,
		"iatが現在時刻であること",
	)
	assert.Equal(t, "test", accessClaims.Iss)
	assert.Equal(t, id, accessClaims.Sub)
	assert.Equal(t, jwt.ClaimStrings{"test"}, accessClaims.Aud)
	assert.NotEmpty(t, accessClaims.Jti)
	assert.Equal(t, "user", accessClaims.Uname)

	refreshClaims, err := ts.ParseRefreshToken(refreshToken)
	assert.NoError(t, err)
	assert.Equal(t,
		mockClock{}.Now().Add(time.Hour*24*7),
		refreshClaims.Exp,
		"expが24時間後であること",
	)
	assert.Equal(t, "test", refreshClaims.Iss)
	assert.Equal(t, id, refreshClaims.Sub)
	assert.NotEmpty(t, refreshClaims.Jti)
}
