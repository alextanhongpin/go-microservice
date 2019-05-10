package infrastructure

import (
	"errors"
	"time"

	jwt "github.com/dgrijalva/jwt-go"

	"github.com/alextanhongpin/go-microservice/api"
	"github.com/alextanhongpin/pkg/gojwt"
)

func NewSigner(cfg *Config) *gojwt.JwtSigner {
	var (
		audience     = cfg.Audience
		issuer       = cfg.Issuer
		semver       = cfg.Semver
		secret       = cfg.Secret
		expiresAfter = 10080 * time.Minute // 1 Week.
		scope        = api.ScopeDefault.String()
		role         = api.RoleGuest.String()
	)
	opt := gojwt.Option{
		Secret:       []byte(secret),
		ExpiresAfter: expiresAfter,
		DefaultClaims: &gojwt.Claims{
			Semver: semver,
			Scope:  scope,
			Role:   role,
			StandardClaims: jwt.StandardClaims{
				Audience: audience,
				Issuer:   issuer,
			},
		},
		Validator: func(c *gojwt.Claims) error {
			if c.Semver != semver ||
				c.Issuer != issuer ||
				c.Audience != audience {
				return errors.New("invalid token")
			}
			return nil
		},
	}
	return gojwt.New(opt)
}
