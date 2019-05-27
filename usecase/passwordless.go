package usecase

import "context"

type Passwordless interface {
	Start(ctx context.Context, req PasswordlessStartRequest) *PasswordlessStartResponse
	Authorize(ctx context.Context, req PasswordlessAuthorizeRequest) (*PasswordlessAuthorizeResponse, error)
}

type (
	PasswordlessStartRequest struct {
		Connection string `json:"connection" validate:"required,oneof:email sms"`
		Send       string `json:"send" validate:"required,oneof=link code"`
		Email      string `json:"email" validate:"required,email"`
	}
	PasswordlessStartResponse struct {
		Link string `json:"link"`
		Code string `json:"code"`
	}
)

type (
	PasswordlessAuthorizeRequest struct {
		Code string `json:"code" validate:"required"`
	}
	PasswordlessAuthorizeResponse struct {
		AccessToken string `json:"access_token" validate:"required"`
	}
)
