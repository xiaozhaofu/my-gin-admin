package auth

import (
	"fmt"
	"slices"
	"time"

	gojwt "github.com/golang-jwt/jwt/v5"
	encryjwt "github.com/gtkit/encry/jwt"
	jwtclaims "github.com/gtkit/encry/jwt/claims"
)

type JWTManager struct {
	issuer          *encryjwt.JwtEd25519
	accessTokenTTL  time.Duration
	refreshTokenTTL time.Duration
}

type Claims struct {
	AdminID  int64    `json:"admin_id"`
	Username string   `json:"username"`
	Roles    []string `json:"roles"`
	TokenUse string   `json:"token_use"`
	gojwt.RegisteredClaims
}

func NewJWTManager(issuer *encryjwt.JwtEd25519, accessTokenTTL, refreshTokenTTL time.Duration) *JWTManager {
	return &JWTManager{
		issuer:          issuer,
		accessTokenTTL:  accessTokenTTL,
		refreshTokenTTL: refreshTokenTTL,
	}
}

func (m *JWTManager) GenerateAccessToken(adminID int64, username string, roles []string) (string, time.Time, error) {
	return m.generate(adminID, username, roles, "access", m.accessTokenTTL)
}

func (m *JWTManager) GenerateRefreshToken(adminID int64, username string, roles []string) (string, time.Time, error) {
	return m.generate(adminID, username, roles, "refresh", m.refreshTokenTTL)
}

func (m *JWTManager) Parse(tokenStr string) (*Claims, error) {
	if m == nil || m.issuer == nil {
		return nil, fmt.Errorf("jwt issuer is required")
	}

	tokenClaims, err := m.issuer.ParseToken(tokenStr)
	if err != nil {
		return nil, err
	}

	return &Claims{
		AdminID:          tokenClaims.UserID,
		Username:         tokenClaims.Subject,
		Roles:            slices.Clone(tokenClaims.Roles),
		TokenUse:         tokenClaims.Prv,
		RegisteredClaims: tokenClaims.RegisteredClaims,
	}, nil
}

func (m *JWTManager) generate(adminID int64, username string, roles []string, tokenUse string, ttl time.Duration) (string, time.Time, error) {
	if m == nil || m.issuer == nil {
		return "", time.Time{}, fmt.Errorf("jwt issuer is required")
	}

	expiresAt := time.Now().Add(ttl)

	token, err := m.issuer.GenerateToken(
		adminID,
		jwtclaims.WithSubject(username),
		jwtclaims.WithRoles(roles...),
		jwtclaims.WithPrv(tokenUse),
		jwtclaims.WithExpiresAt(ttl),
	)
	if err != nil {
		return "", time.Time{}, err
	}

	return token, expiresAt, nil
}
