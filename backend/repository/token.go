package repository

import (
	"context"
	"errors"
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
	exp := time.Duration(exptime.Unix())

	return r.store.Set(ctx, id, jti, exp).Err()
}

func (r *tokenRepository) ExistsJTI(jti string) error {
	// TODO
	return errors.New("not implemented this method: tokenStor.ExistsJTI")
}
