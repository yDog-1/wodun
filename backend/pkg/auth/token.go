package auth

import (
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

const (
	// アクセストークンの有効期限 (1時間)
	accessTokenExpire = time.Hour
	// リフレッシュトークンの有効期限 (1週間)
	refreshTokenExpire = time.Hour * 24 * 7
)

type clock interface {
	Now() time.Time
}

type Token struct {
	Exp   time.Time
	Iat   time.Time
	Iss   string
	Sub   string
	Aud   jwt.ClaimStrings
	Jti   string
	Uname string
}

type TokenService struct {
	store         TokenStore
	issuer        string
	audience      string
	accessSecret  []byte
	refreshSecret []byte
	clock         clock
}

type TokenStore interface {
	// jtiを保存する
	SaveJTI(jti string) error
	// jtiが存在するか確認する
	ExistsJTI(jti string) (bool, error)
}

// TokenServiceを生成する
func NewTokenService(store TokenStore, clock clock) (*TokenService, error) {
	if store == nil {
		return nil, errors.New("store is nil")
	}
	if clock == nil {
		return nil, errors.New("clock is nil")
	}
	iss, ok := os.LookupEnv("TOKEN_ISSUER")
	if !ok {
		return nil, errors.New("TOKEN_ISSUER is not set")
	}
	aud, ok := os.LookupEnv("TOKEN_AUDIENCE")
	if !ok {
		return nil, errors.New("TOKEN_AUDIENCE is not set")
	}
	as, ok := os.LookupEnv("TOKEN_ACCESS_SECRET")
	if !ok {
		return nil, errors.New("TOKEN_ACCESS_SECRET is not set")
	}
	rs, ok := os.LookupEnv("TOKEN_REFRESH_SECRET")
	if !ok {
		return nil, errors.New("TOKEN_REFRESH_SECRET is not set")
	}
	return &TokenService{
		store:         store,
		issuer:        iss,
		audience:      aud,
		accessSecret:  []byte(as),
		refreshSecret: []byte(rs),
		clock:         clock,
	}, nil
}

// トークンを生成する
func (ts *TokenService) GenerateToken(id, uniqueName string) (accessToken string, refreshToken string, err error) {
	at, err := ts.generateAccessToken(id, uniqueName)
	if err != nil {
		return "", "", err
	}
	rt, err := ts.generateRefreshToken(id)
	if err != nil {
		return "", "", err
	}
	return at, rt, nil
}

// アクセストークンを生成する
func (ts *TokenService) generateAccessToken(id, uniqueName string) (string, error) {
	jti := uuid.New().String()
	claims := jwt.MapClaims{
		"exp":   ts.clock.Now().Add(accessTokenExpire).Unix(),
		"iat":   ts.clock.Now().Unix(),
		"iss":   ts.issuer,
		"sub":   id,
		"aud":   ts.audience,
		"jti":   jti,
		"uname": uniqueName,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString(ts.accessSecret)
	if err != nil {
		return "", err
	}
	if err := ts.store.SaveJTI(jti); err != nil {
		return "", err
	}

	return tokenString, nil
}

// リフレッシュトークンを生成する
func (ts *TokenService) generateRefreshToken(id string) (string, error) {
	jti := uuid.New().String()
	claims := jwt.MapClaims{
		"exp": ts.clock.Now().Add(refreshTokenExpire).Unix(),
		"iss": ts.issuer,
		"sub": id,
		"jti": jti,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString(ts.refreshSecret)
	if err != nil {
		return "", err
	}
	if err := ts.store.SaveJTI(jti); err != nil {
		return "", err
	}

	return tokenString, nil
}

func (ts *TokenService) ParseAccessToken(token string) (*Token, error) {
	t, err := jwt.Parse(token, func(t *jwt.Token) (interface{}, error) {
		if t.Method.Alg() != jwt.SigningMethodHS256.Alg() {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}
		return []byte(ts.accessSecret), nil
	})

	if claims, ok := t.Claims.(jwt.MapClaims); ok && t.Valid {
		exp, err := claims.GetExpirationTime()
		if err != nil {
			return nil, err
		}
		expUTC := exp.Time.UTC()
		iat, err := claims.GetIssuedAt()
		if err != nil {
			return nil, err
		}
		iatUTC := iat.Time.UTC()
		iss, err := claims.GetIssuer()
		if err != nil {
			return nil, err
		}
		sub, err := claims.GetSubject()
		if err != nil {
			return nil, err
		}
		aud, err := claims.GetAudience()
		if err != nil {
			return nil, err
		}
		jti, ok := claims["jti"].(string)
		if !ok {
			return nil, errors.New("jti is not set")
		}
		uname, ok := claims["uname"].(string)
		if !ok {
			return nil, errors.New("uname is not set")
		}
		return &Token{
			Exp:   expUTC,
			Iat:   iatUTC,
			Iss:   iss,
			Sub:   sub,
			Aud:   aud,
			Jti:   jti,
			Uname: uname,
		}, nil
	} else {
		return nil, err
	}
}

func (ts *TokenService) ParseRefreshToken(token string) (*Token, error) {
	t, err := jwt.Parse(token, func(t *jwt.Token) (interface{}, error) {
		if t.Method.Alg() != jwt.SigningMethodHS256.Alg() {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}
		return []byte(ts.accessSecret), nil
	})

	if claims, ok := t.Claims.(jwt.MapClaims); ok && t.Valid {
		exp, err := claims.GetExpirationTime()
		if err != nil {
			return nil, err
		}
		expUTC := exp.Time.UTC()
		iss, err := claims.GetIssuer()
		if err != nil {
			return nil, err
		}
		sub, err := claims.GetSubject()
		if err != nil {
			return nil, err
		}
		jti, ok := claims["jti"].(string)
		if !ok {
			return nil, errors.New("jti is not set")
		}
		return &Token{
			Exp: expUTC,
			Iss: iss,
			Sub: sub,
			Jti: jti,
		}, nil
	} else {
		return nil, err
	}
}
