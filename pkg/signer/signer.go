package signer

import (
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/pkg/errors"
)

type (
	Signer interface {
		Sign(Claims) (string, error)
		Verify(tokenString string) (*Claims, error)
		NewClaims(user, role, scope string) Claims
	}
	Option struct {
		Secret            []byte
		DurationInMinutes time.Duration

		// The software organization that issues this token, e.g. Alex
		// Inc.
		Issuer string

		// The identity of the intended recipient of the token, e.g.
		// paymentsvc.
		Audience string

		// The current service API version. Note that we can easily
		// expire all the old tokens simply by changing the version.
		Semver string
	}
	SignerImpl struct {
		opt Option
	}
	Claims struct {
		// The actor of the system.
		Role string `json:"role"`

		// The resources the actor can access.
		Scope string `json:"scope"`

		// A specific version to expire all the old tokens.
		Semver string `json:"version"`
		jwt.StandardClaims
	}
)

func New(opt Option) *SignerImpl {
	return &SignerImpl{opt}
}

func (s *SignerImpl) Sign(claims Claims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	ss, err := token.SignedString(s.opt.Secret)
	return ss, errors.Wrap(err, "signing token failed")
}

func (s *SignerImpl) Verify(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return s.opt.Secret, nil
	})
	// Apparently this is possible by sending Authorization: Bearer
	// undefined.
	if token == nil {
		return nil, errors.New("invalid authorization header")
	}
	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		// The most basic validation - checking if this is the exact
		// issuer. Note that you can create multiple signer for
		// different tokens.
		if claims.Semver != s.opt.Semver ||
			claims.Issuer != s.opt.Issuer ||
			claims.Audience != s.opt.Audience {
			return nil, errors.New("invalid token")
		}
		return claims, nil
	}
	return nil, err
}

func (s *SignerImpl) NewClaims(user, role, scope string) Claims {
	now := time.Now()
	return Claims{
		Role:   role,
		Scope:  scope,
		Semver: s.opt.Semver,
		StandardClaims: jwt.StandardClaims{
			Audience:  s.opt.Audience,
			Issuer:    s.opt.Issuer,
			ExpiresAt: now.Add(s.opt.DurationInMinutes).Unix(),
			IssuedAt:  now.Unix(),
			Subject:   user,
		},
	}
}
