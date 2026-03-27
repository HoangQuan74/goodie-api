package auth

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
	apperrors "github.com/kainguyen/goodie-api/pkg/errors"
)

type Claims struct {
	UserID string `json:"user_id"`
	Role   string `json:"role"`
	jwt.RegisteredClaims
}

type TokenPair struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresAt    int64  `json:"expires_at"`
}

type JWTConfig struct {
	Secret     string
	AccessTTL  time.Duration
	RefreshTTL time.Duration
	Issuer     string
}

type JWTManager struct {
	config JWTConfig
}

func NewJWTManager(cfg JWTConfig) *JWTManager {
	return &JWTManager{config: cfg}
}

func (m *JWTManager) GenerateTokenPair(userID, role string) (*TokenPair, error) {
	now := time.Now()

	accessClaims := Claims{
		UserID: userID,
		Role:   role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(now.Add(m.config.AccessTTL)),
			IssuedAt:  jwt.NewNumericDate(now),
			Issuer:    m.config.Issuer,
		},
	}
	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims)
	accessStr, err := accessToken.SignedString([]byte(m.config.Secret))
	if err != nil {
		return nil, apperrors.InternalWrap("failed to sign access token", err)
	}

	refreshClaims := Claims{
		UserID: userID,
		Role:   role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(now.Add(m.config.RefreshTTL)),
			IssuedAt:  jwt.NewNumericDate(now),
			Issuer:    m.config.Issuer,
		},
	}
	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims)
	refreshStr, err := refreshToken.SignedString([]byte(m.config.Secret))
	if err != nil {
		return nil, apperrors.InternalWrap("failed to sign refresh token", err)
	}

	return &TokenPair{
		AccessToken:  accessStr,
		RefreshToken: refreshStr,
		ExpiresAt:    accessClaims.ExpiresAt.Unix(),
	}, nil
}

func (m *JWTManager) ValidateToken(tokenStr string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, apperrors.Unauthorized("invalid signing method")
		}
		return []byte(m.config.Secret), nil
	})
	if err != nil {
		return nil, apperrors.Unauthorized("invalid or expired token")
	}

	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, apperrors.Unauthorized("invalid token claims")
	}

	return claims, nil
}
