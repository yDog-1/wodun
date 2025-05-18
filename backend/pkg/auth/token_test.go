package auth

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type mockTokenStore struct{}

func (m *mockTokenStore) SaveJTI(ctx context.Context, jti string) error {
	return nil
}

func (m *mockTokenStore) ExistsJTI(ctx context.Context, id, jti string) (bool, error) {
	return true, nil
}

type mockClock struct{}

func (c mockClock) Now() time.Time {
	return time.Date(2100, 1, 1, 0, 0, 0, 0, time.UTC)
}

const (
	testIssuer        = "test-issuer"
	testAudience      = "test-audience"
	testAccessSecret  = "test-access-secret"
	testRefreshSecret = "test-refresh-secret"
)

func TestMain(m *testing.M) {

	// テスト実行前に環境変数を設定
	_ = os.Setenv("TOKEN_ISSUER", testIssuer)
	_ = os.Setenv("TOKEN_AUDIENCE", testAudience)
	_ = os.Setenv("TOKEN_ACCESS_SECRET", testAccessSecret)
	_ = os.Setenv("TOKEN_REFRESH_SECRET", testRefreshSecret)

	// テストを実行
	code := m.Run()

	// テスト実行後に環境変数をクリーンアップ
	_ = os.Unsetenv("TOKEN_ISSUER")
	_ = os.Unsetenv("TOKEN_AUDIENCE")
	_ = os.Unsetenv("TOKEN_ACCESS_SECRET")
	_ = os.Unsetenv("TOKEN_REFRESH_SECRET")

	os.Exit(code)
}

func TestTokenService_GenerateToken(t *testing.T) {
	store := &mockTokenStore{}
	clock := mockClock{}

	ts, err := NewTokenService(store, clock)
	require.NoError(t, err)

	id := "user123"
	uniqueName := "testuser"
	accessToken, refreshToken, err := ts.GenerateToken(context.Background(), id, uniqueName)
	require.NoError(t, err)

	assert.NotEmpty(t, accessToken, "accessToken should not be empty")
	assert.NotEmpty(t, refreshToken, "refreshToken should not be empty")

	// アクセストークンの検証
	accessClaims, err := ts.ParseAccessToken(accessToken)
	require.NoError(t, err, "Failed to parse access token")

	assert.Equal(t, id, accessClaims.Sub, "Access token subject mismatch")
	assert.Equal(t, uniqueName, accessClaims.Uname, "Access token unique name mismatch")
	assert.Equal(t, testIssuer, accessClaims.Iss, "Access token issuer mismatch")
	assert.Contains(t, accessClaims.Aud, testAudience, "Access token audience mismatch")
	assert.Equal(t, clock.Now().Add(accessTokenExpire).Unix(), accessClaims.Exp.Unix(), "Access token expiration mismatch")
	assert.Equal(t, clock.Now().Unix(), accessClaims.Iat.Unix(), "Access token issued at mismatch")
	assert.NotEmpty(t, accessClaims.Jti, "Access token JTI should not be empty")

	// リフレッシュトークンの検証
	refreshClaims, err := ts.ParseRefreshToken(refreshToken)
	require.NoError(t, err, "Failed to parse refresh token")

	assert.Equal(t, id, refreshClaims.Sub, "Refresh token subject mismatch")
	assert.Equal(t, testIssuer, refreshClaims.Iss, "Refresh token issuer mismatch")
	assert.Equal(t, clock.Now().Add(refreshTokenExpire).Unix(), refreshClaims.Exp.Unix(), "Refresh token expiration mismatch")
	assert.NotEmpty(t, refreshClaims.Jti, "Refresh token JTI should not be empty")
}

func TestTokenService_ParseAccessToken(t *testing.T) {
	store := &mockTokenStore{}
	clock := mockClock{}
	ts, err := NewTokenService(store, clock)
	require.NoError(t, err)

	id := "user123"
	uniqueName := "testuser"
	accessToken, _, err := ts.GenerateToken(context.Background(), id, uniqueName)
	require.NoError(t, err)

	parsedToken, err := ts.ParseAccessToken(accessToken)
	require.NoError(t, err, "Failed to parse access token")

	assert.Equal(t, id, parsedToken.Sub, "Parsed token subject mismatch")
	assert.Equal(t, uniqueName, parsedToken.Uname, "Parsed token unique name mismatch")
	assert.Equal(t, testIssuer, parsedToken.Iss, "Parsed token issuer mismatch")
	assert.Contains(t, parsedToken.Aud, testAudience, "Parsed token audience mismatch")
	assert.Equal(t, clock.Now().Add(accessTokenExpire).Unix(), parsedToken.Exp.Unix(), "Parsed token expiration mismatch")
	assert.Equal(t, clock.Now().Unix(), parsedToken.Iat.Unix(), "Parsed token issued at mismatch")
	assert.NotEmpty(t, parsedToken.Jti, "Parsed token JTI should not be empty")

	// 無効なトークンのテスト
	invalidToken := "invalid.token.string"
	parsedToken, err = ts.ParseAccessToken(invalidToken)
	assert.Error(t, err, "ParseAccessToken should return error for invalid token")
	assert.Nil(t, parsedToken, "ParseAccessToken should return nil for invalid token")
}

func TestTokenService_ParseRefreshToken(t *testing.T) {
	store := &mockTokenStore{}
	clock := mockClock{}
	ts, err := NewTokenService(store, clock)
	require.NoError(t, err)

	id := "user123"
	_, refreshToken, err := ts.GenerateToken(context.Background(), id, "testuser") // uniqueNameはリフレッシュトークンには含まれない
	require.NoError(t, err)

	parsedToken, err := ts.ParseRefreshToken(refreshToken)
	require.NoError(t, err, "Failed to parse access token")

	assert.Equal(t, id, parsedToken.Sub, "Parsed token subject mismatch")
	assert.Equal(t, testIssuer, parsedToken.Iss, "Parsed token issuer mismatch")
	assert.Equal(t, clock.Now().Add(refreshTokenExpire).Unix(), parsedToken.Exp.Unix(), "Parsed token expiration mismatch")
	assert.NotEmpty(t, parsedToken.Jti, "Parsed token JTI should not be empty")

	// 無効なトークンのテスト
	invalidToken := "invalid.token.string"
	parsedToken, err = ts.ParseRefreshToken(invalidToken)
	assert.Error(t, err, "ParseRefreshToken should return error for invalid token")
	assert.Nil(t, parsedToken, "ParseRefreshToken should return nil for invalid token")

	// 異なるシークレットで署名されたトークンのテスト
	otherSecretTS, err := NewTokenService(store, clock)
	require.NoError(t, err)
	otherSecretTS.refreshSecret = []byte("other_refresh_secret")
	_, otherRefreshToken, err := otherSecretTS.GenerateToken(context.Background(), id, "testuser")
	require.NoError(t, err)

	parsedToken, err = ts.ParseRefreshToken(otherRefreshToken)
	assert.Error(t, err, "ParseRefreshToken should return error for token signed with different secret")
	assert.Nil(t, parsedToken, "ParseRefreshToken should return nil for token signed with different secret")
}
