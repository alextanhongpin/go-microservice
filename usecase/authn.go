package usecase

import (
	"context"

	"github.com/alextanhongpin/go-microservice/domain/user"
)

type Authn interface {
	Login(ctx context.Context, req LoginRequest) (*LoginResponse, error)
	Register(ctx context.Context, req RegisterRequest) (*RegisterResponse, error)
	RecoverPassword(ctx context.Context, req RecoverPasswordRequest) (*RecoverPasswordResponse, error)
	ResetPassword(ctx context.Context, req ResetPasswordRequest) (*ResetPasswordResponse, error)
	ChangePassword(ctx context.Context, req ChangePasswordRequest) (*ChangePasswordResponse, error)
}

type (
	LoginRequest struct {
		Email    string `json:"email" validate:"required,email"`
		Password string `json:"password" validate:"required,min=8"`
	}
	LoginResponse struct {
		AccessToken string      `json:"access_token"`
		Data        user.Entity `json:"data"`
	}
)

type (
	RegisterRequest struct {
		Email    string `json:"email" validate:"required,email"`
		Password string `json:"password" validate:"required,min=8"`
	}
	RegisterResponse struct {
		Data        user.Entity `json:"data"`
		AccessToken string      `json:"access_token"`
	}
)

type (
	RecoverPasswordRequest struct {
		Email string `json:"required,email"`
	}
	RecoverPasswordResponse struct {
		Success bool `json:"success"`
		// Token will be sent as capability url to the user's email.
		Token string `json:"-"`
	}
)

type (
	ResetPasswordRequest struct {
		Token           string `json:"token" validate:"required"`
		Password        string `json:"password" validate:"required"`
		ConfirmPassword string `json:"confirm_password" validate:"required"`
	}
	ResetPasswordResponse struct {
		Success bool
	}
)

type (
	ChangePasswordRequest struct {
		ContextUserID   string `json:"-" validate:"required"`
		OldPassword     string `json:"old_password" validate:"required,neqfield=NewPassword"`
		NewPassword     string `json:"new_password" validate:"required,eqfield=ConfirmPassword"`
		ConfirmPassword string `json:"confirm_password" validate:"required"`
	}
	ChangePasswordResponse struct {
		Success bool `json:"success"`
	}
)
