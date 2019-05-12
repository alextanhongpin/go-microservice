package authn

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

	usecase interface {
		service
		loginUseCase
		registerUseCase
		// Extended use cases.
		LoginWithAccessToken(context.Context, LoginRequest) (string, error)
		RegisterWithAccessToken(context.Context, RegisterRequest) (string, error)
	}

	// UseCase represents the authentication usecases.
	UseCase struct {
		usecase
		service  *Service
		login    *LoginUseCase
		register *RegisterUseCase
	}
)

// NewUseCase returns the individual usecases + compound use cases (use cases
// which includes other usecases).
func NewUseCase(repo repository, signer gojwt.Signer) *UseCase {
	return &UseCase{
		// Service.
		service: NewService(signer),
		// UseCase.
		login:    NewLoginUseCase(repo),
		register: NewRegisterUseCase(repo),
	}
}

// LoginWithAccessToken logins an existing user, and generates an access token
// if the login succeeds.
func (u *UseCase) LoginWithAccessToken(ctx context.Context, req LoginRequest) (string, error) {
	res, err := u.Login(ctx, req)
	if err != nil {
		return "", errors.Wrap(err, "login failed")
	}
	accessToken, err := u.CreateAccessToken(res.User.ID)
	return accessToken, errors.Wrap(err, "login with access token failed")
}

// RegisterWithAccessToken registers a new user, and generates an access token
// if the registration succeeds.
func (u *UseCase) RegisterWithAccessToken(ctx context.Context, req RegisterRequest) (string, error) {
	res, err := u.Register(ctx, req)
	if err != nil {
		return "", errors.Wrap(err, "register failed")
	}
	accessToken, err := u.CreateAccessToken(res.User.ID)
	return accessToken, errors.Wrap(err, "register with access token failed")
}
