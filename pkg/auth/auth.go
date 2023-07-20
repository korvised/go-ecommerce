package auth

import (
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"github.com/korvised/go-ecommerce/config"
	"github.com/korvised/go-ecommerce/modules/users"
	"math"
	"time"
)

type TokenType string

const (
	Access  TokenType = "access"
	Refresh TokenType = "refresh"
	Admin   TokenType = "admin"
	ApiKey  TokenType = "apiKey"
)

type IAuth interface {
	SignToken() string
}

type Auth struct {
	mapClaims *MapClaims
	cfg       config.IJwtConfig
}

type MapClaims struct {
	Claims *users.UserClaims `json:"claims"` // mean payload
	jwt.RegisteredClaims
}

func jwtTimeDurationCal(t int) *jwt.NumericDate {
	return jwt.NewNumericDate(time.Now().Add(time.Duration(int64(t) * int64(math.Pow10(9)))))
}

func jwtTimeRepeatAdapter(t int64) *jwt.NumericDate {
	return jwt.NewNumericDate(time.Unix(t, 0))
}

func NewAuth(tokenType TokenType, cfg config.IJwtConfig, claims *users.UserClaims) (IAuth, error) {
	switch tokenType {
	case Access:
		return newAccessToken(cfg, claims), nil
	case Refresh:
		return newRefreshToken(cfg, claims), nil
	case Admin:
		return newAccessToken(cfg, claims), nil
	default:
		return nil, fmt.Errorf("unknown token type")
	}
}

func newAccessToken(cfg config.IJwtConfig, claims *users.UserClaims) IAuth {
	return &Auth{
		cfg: cfg,
		mapClaims: &MapClaims{
			Claims: claims,
			RegisteredClaims: jwt.RegisteredClaims{
				Issuer:    "ecommerce-api",
				Subject:   "access-token",
				Audience:  []string{"customer", "admin"},
				ExpiresAt: jwtTimeDurationCal(cfg.AccessExpiresAt()),
				NotBefore: jwt.NewNumericDate(time.Now()),
				IssuedAt:  jwt.NewNumericDate(time.Now()),
			},
		},
	}
}

func newRefreshToken(cfg config.IJwtConfig, claims *users.UserClaims) IAuth {
	return &Auth{
		cfg: cfg,
		mapClaims: &MapClaims{
			Claims: claims,
			RegisteredClaims: jwt.RegisteredClaims{
				Issuer:    "ecommerce-api",
				Subject:   "refresh-token",
				Audience:  []string{"customer", "admin"},
				ExpiresAt: jwtTimeDurationCal(cfg.RefreshExpiresAt()),
				NotBefore: jwt.NewNumericDate(time.Now()),
				IssuedAt:  jwt.NewNumericDate(time.Now()),
			},
		},
	}
}

func (a Auth) SignToken() string {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, a.mapClaims)
	ss, _ := token.SignedString(a.cfg.SecretKey())
	return ss
}
