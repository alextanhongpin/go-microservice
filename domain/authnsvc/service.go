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
	"github.com/pkg/errors"
)

type (
	repository interface {
		// Reader.
		WithEmail(email string) (User, error)

		// Writer.
		Create(username, password string) (User, error)
	}

	service interface {
		// Base UseCase.
		Login(ctx context.Context, req LoginRequest) (*LoginResponse, error)
		Register(ctx context.Context, req RegisterRequest) (*RegisterResponse, error)
		CreateAccessToken(userID string) (string, error)

		// UseCases with include UseCase
		LoginWithAccessToken(context.Context, LoginRequest) (string, error)
		RegisterWithAccessToken(context.Context, RegisterRequest) (string, error)
	}

	// Service represents the authentication service.
	Service struct {
		service
		login             *LoginUseCase
		register          *RegisterUseCase
		createAccessToken *CreateAccessTokenUseCase
	}
)

// NewService returns the individual usecases + compound use cases (use cases
// which includes other usecases).
func NewService(repo repository, signer gojwt.Signer) *Service {
	return &Service{
		login:             NewLoginUseCase(repo),
		register:          NewRegisterUseCase(repo),
		createAccessToken: NewCreateAccessTokenUseCase(signer),
	}
}

// LoginWithAccessToken logins an existing user, and generates an access token
// if the login succeeds.
func (s *Service) LoginWithAccessToken(ctx context.Context, req LoginRequest) (string, error) {
	res, err := s.Login(ctx, req)
	if err != nil {
		return "", errors.Wrap(err, "login failed")
	}
	accessToken, err := s.CreateAccessToken(res.User.ID)
	return accessToken, errors.Wrap(err, "login with access token failed")
}

// RegisterWithAccessToken registers a new user, and generates an access token
// if the registration succeeds.
func (s *Service) RegisterWithAccessToken(ctx context.Context, req RegisterRequest) (string, error) {
	res, err := s.Register(ctx, req)
	if err != nil {
		return "", errors.Wrap(err, "register failed")
	}
	accessToken, err := s.CreateAccessToken(res.User.ID)
	return accessToken, errors.Wrap(err, "register with access token failed")
}
