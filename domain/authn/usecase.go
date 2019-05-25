package authn

import (
	"context"
	"time"

	"github.com/pkg/errors"

	"github.com/alextanhongpin/pkg/gojwt"
)

type (
	usecase interface {
		service
		loginUseCase
		registerUseCase
		changePasswordUseCase
		recoverPasswordUseCase
		resetPasswordUseCase
		// Extended use cases.
		LoginWithAccessToken(context.Context, LoginRequest) (string, error)
		RegisterWithAccessToken(context.Context, RegisterRequest) (string, error)
	}

	// UseCase represents the authentication usecases.
	UseCase struct {
		*Service
		*LoginUseCase
		*RegisterUseCase
		*ChangePasswordUseCase
		*RecoverPasswordUseCase
		*ResetPasswordUseCase
	}
)

// NewUseCase returns the individual usecases + compound use cases (use cases
// which includes other usecases).
func NewUseCase(repo Repository, signer gojwt.Signer, tokenTTL time.Duration) (*UseCase, func()) {
	recoverPassword, shutdown := NewRecoverPasswordUseCase(repo, tokenTTL)
	return &UseCase{
			// Service.
			Service: NewService(signer),
			// UseCase.
			LoginUseCase:           NewLoginUseCase(repo),
			RegisterUseCase:        NewRegisterUseCase(repo),
			ChangePasswordUseCase:  NewChangePasswordUseCase(repo),
			ResetPasswordUseCase:   NewResetPasswordUseCase(repo, tokenTTL),
			RecoverPasswordUseCase: recoverPassword,
		}, func() {
			shutdown()
		}
}

// LoginWithAccessToken logins an existing user, and generates an access token
// if the login succeeds.
func (u *UseCase) LoginWithAccessToken(ctx context.Context, req LoginRequest) (string, error) {
	res, err := u.Login(ctx, req)
	if err != nil {
		return "", errors.Wrap(err, "login failed")
	}
	accessToken, err := u.Service.CreateAccessToken(res.Data.ID)
	return accessToken, errors.Wrap(err, "login with access token failed")
}

// RegisterWithAccessToken registers a new user, and generates an access token
// if the registration succeeds.
func (u *UseCase) RegisterWithAccessToken(ctx context.Context, req RegisterRequest) (string, error) {
	res, err := u.Register(ctx, req)
	if err != nil {
		return "", errors.Wrap(err, "register failed")
	}
	accessToken, err := u.Service.CreateAccessToken(res.User.ID)
	return accessToken, errors.Wrap(err, "register with access token failed")
}
