package authn

import (
	"context"
	"time"

	"github.com/alextanhongpin/go-microservice/pkg/govalidator"
	"github.com/alextanhongpin/passwd"
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
	resetPasswordRepository interface {
		TokenWithValue(token string) (Token, error)
		UpdateUserPassword(userID, password string) (bool, error)
		DeleteToken(token string) (bool, error)
	}
	resetPasswordUseCase interface {
		ResetPassword(ctx context.Context, req ResetPasswordRequest) (*ResetPasswordResponse, error)
	}
)

type ResetPasswordUseCase struct {
	repo     resetPasswordRepository
	tokenTTL time.Duration
}

func NewResetPasswordUseCase(repo resetPasswordRepository, tokenTTL time.Duration) *ResetPasswordUseCase {
	return &ResetPasswordUseCase{
		repo:     repo,
		tokenTTL: tokenTTL,
	}
}

func (r *ResetPasswordUseCase) ResetPassword(ctx context.Context, req ResetPasswordRequest) (*ResetPasswordResponse, error) {
	if err := govalidator.Validate.Struct(req); err != nil {
		return nil, err
	}
	if req.Password != req.ConfirmPassword {
		return nil, ErrInvalidPassword
	}

	// Check if the hash token exists in the database.
	token, err := r.repo.TokenWithValue(hashToken(req.Token))
	if err != nil {
		return nil, err
	}

	// Check for the token expiry.
	if token.HasExpired(r.tokenTTL) {
		// Delete the token.
		_, _ = r.repo.DeleteToken(hashToken(req.Token))
		return nil, ErrTokenExpired
	}

	// Hash the password before storing.
	hashedPassword, err := passwd.Hash(req.Password)
	if err != nil {
		return nil, err
	}
	success, err := r.repo.UpdateUserPassword(token.UserID, hashedPassword)
	if err != nil {
		return nil, err
	}
	// Clear the previously used token to ensure it's a one-time used only
	// token.
	success, err = r.repo.DeleteToken(hashToken(req.Token))
	return &ResetPasswordResponse{
		Success: success,
	}, err
}
