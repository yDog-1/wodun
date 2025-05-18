package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/redis/go-redis/v9"
)

type tokenRepository struct {
	store *redis.Client
}

type storeClaims interface {
	jwt.Claims
	GetJTI() (string, error)
}

func NewTokenRepository(store *redis.Client) *tokenRepository {
	return &tokenRepository{store}
}

// JTIを保存する
// RedisのキーはユーザーID、値はJTIとする
func (r *tokenRepository) SaveJTI(ctx context.Context, claims storeClaims) error {
	id, err := claims.GetSubject()
	if err != nil {
		return err
	}
	jti, err := claims.GetJTI()
	if err != nil {
		return err
	}
	exptime, err := claims.GetExpirationTime()
	if err != nil {
		return err
	}
	// 有効期限をtime.Durationに変換
	exp := time.Until(exptime.Time)

	// RedisにユーザーIDをキー、JTIを値として保存。有効期限も設定する。
	return r.store.Set(ctx, id, jti, exp).Err()
}

// JTIが存在するか確認する
// RedisにユーザーIDをキーとして保存されているJTIと、引数で渡されたJTIが一致するかを確認する
func (r *tokenRepository) ExistsJTI(ctx context.Context, id, jti string) (bool, error) {
	// RedisからユーザーIDに対応するJTIを取得
	storedJTI, err := r.store.Get(ctx, id).Result()
	if err == redis.Nil {
		// キーが存在しない場合はJTIも存在しない
		return false, nil
	} else if err != nil {
		// その他のエラー
		return false, fmt.Errorf("failed to get JTI from redis: %w", err)
	}

	// 取得したJTIと引数のJTIが一致するか確認
	return storedJTI == jti, nil
}
