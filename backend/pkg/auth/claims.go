package auth

import (
	"github.com/golang-jwt/jwt/v5"
)

// アクセストークンの claims
type accessClaims struct {
	Issuer    string           `json:"iss"`
	Subject   string           `json:"sub"`
	Audience  jwt.ClaimStrings `json:"aud"`
	ExpiresAt *jwt.NumericDate `json:"exp"`
	IssuedAt  *jwt.NumericDate `json:"iat"`
	ID        string           `json:"jti"`
	UniqueName string         `json:"uname"`
}

// GetExpirationTime implements the Claims interface.
func (a accessClaims) GetExpirationTime() (*jwt.NumericDate, error) {
	return a.ExpiresAt, nil
}

// GetNotBefore implements the Claims interface.
func (a accessClaims) GetNotBefore() (*jwt.NumericDate, error) {
	return nil, nil
}

// GetIssuedAt implements the Claims interface.
func (a accessClaims) GetIssuedAt() (*jwt.NumericDate, error) {
	return a.IssuedAt, nil
}

// GetAudience implements the Claims interface.
func (a accessClaims) GetAudience() (jwt.ClaimStrings, error) {
	return a.Audience, nil
}

// GetIssuer implements the Claims interface.
func (a accessClaims) GetIssuer() (string, error) {
	return a.Issuer, nil
}

// GetSubject implements the Claims interface.
func (a accessClaims) GetSubject() (string, error) {
	return a.Subject, nil
}

// リフレッシュトークンの claims
type refreshClaims struct {
	Issuer    string           `json:"iss"`
	Subject   string           `json:"sub"`
	ExpiresAt *jwt.NumericDate `json:"exp"`
	ID        string           `json:"jti"`
}

// GetExpirationTime implements the Claims interface.
func (r refreshClaims) GetExpirationTime() (*jwt.NumericDate, error) {
	return r.ExpiresAt, nil
}

// GetNotBefore implements the Claims interface.
func (r refreshClaims) GetNotBefore() (*jwt.NumericDate, error) {
	return nil, nil
}

// GetIssuedAt implements the Claims interface.
func (r refreshClaims) GetIssuedAt() (*jwt.NumericDate, error) {
	return nil, nil
}

// GetAudience implements the Claims interface.
func (r refreshClaims) GetAudience() (jwt.ClaimStrings, error) {
	return nil, nil
}

// GetIssuer implements the Claims interface.
func (r refreshClaims) GetIssuer() (string, error) {
	return r.Issuer, nil
}

// GetSubject implements the Claims interface.
func (r refreshClaims) GetSubject() (string, error) {
	return r.Subject, nil
}
