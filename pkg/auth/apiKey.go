package auth

import (
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"github.com/korvised/go-ecommerce/config"
	"time"
)

type apikey struct {
	*auth
}

func (a apikey) SignToken() string {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, a.mapClaims)
	ss, _ := token.SignedString(a.cfg.ApiKey())
	return ss
}

func newApiKey(cfg config.IJwtConfig) IAuth {
	return &apikey{
		auth: &auth{
			cfg: cfg,
			mapClaims: &MapClaims{
				Claims: nil,
				RegisteredClaims: jwt.RegisteredClaims{
					Issuer:    "ecommerce-api",
					Subject:   "api-key",
					Audience:  []string{"admin", "customer"},
					ExpiresAt: jwt.NewNumericDate(time.Now().AddDate(1, 0, 0)), // 1 year
					NotBefore: jwt.NewNumericDate(time.Now()),
					IssuedAt:  jwt.NewNumericDate(time.Now()),
				},
			},
		},
	}
}

func ParseApiKey(cfg config.IJwtConfig, tokenString string) (*MapClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &MapClaims{}, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("signing method is invalid")
		}

		return cfg.ApiKey(), nil
	})
	if err != nil {
		if errors.Is(err, jwt.ErrTokenMalformed) {
			return nil, fmt.Errorf("token format is invalid")
		} else if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, fmt.Errorf("token had expired")
		} else {
			return nil, fmt.Errorf("parse token failed: %v", err)
		}
	}

	if claims, ok := token.Claims.(*MapClaims); ok {
		return claims, nil
	} else {
		return nil, fmt.Errorf("claims type is invalid")
	}
}
