// Attempted the following name for package:
// - authenticator: this sounds more like a verb
// - authentication: too long
// - userlogin: is too specific, since user can also register
// - loginUser: breaks the convention, since package name is preferable a noun.
// - authz and authn is better.

package authnsvc

import (
	"context"

	"github.com/alextanhongpin/pkg/gojwt"
)

type (
	loginUseCase interface {
		Login(ctx context.Context, req LoginRequest) (*LoginResponse, error)
	}
	registerUseCase interface {
		Register(ctx context.Context, req RegisterRequest) (*RegisterResponse, error)
	}
	repository interface {
		// Reader.
		WithEmail(email string) (User, error)

		// Writer.
		Create(username, password string) (User, error)
	}
	Service struct {
		loginUseCase
		registerUseCase
	}
)

func NewService(repo repository, signer gojwt.Signer) *Service {
	createAccessToken := NewCreateAccessTokenUseCase(signer)
	return &Service{
		loginUseCase:    NewLoginUseCase(repo, createAccessToken),
		registerUseCase: NewRegisterUseCase(repo, createAccessToken),
	}
}
