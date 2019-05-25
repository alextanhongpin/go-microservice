package authn

import (
	"context"

	"github.com/pkg/errors"

	"github.com/alextanhongpin/go-microservice/pkg/govalidator"
)

// Request/response.
type (
	ChangePasswordRequest struct {
		// We could have stored the id in the ctx object, but we want
		// it to be validated.
		ContextUserID   string `json:"-" validate:"required"`
		OldPassword     string `json:"old_password" validate:"required"`
		NewPassword     string `json:"new_password" validate:"required"`
		ConfirmPassword string `json:"confirm_password" validate:"required"`
	}
	ChangePasswordResponse struct {
		Success bool `json:"success"`
	}
)

// Dependencies interface.
type (
	changePasswordRepository interface {
		UpdateUserPassword(userID, password string) (bool, error)
	}
	changePasswordUseCase interface {
		ChangePassword(ctx context.Context, req ChangePasswordRequest) (*ChangePasswordResponse, error)
	}
)

type ChangePasswordUseCase struct {
	repo changePasswordRepository
}

func NewChangePasswordUseCase(repo changePasswordRepository) *ChangePasswordUseCase {
	return &ChangePasswordUseCase{repo}
}

// ChangePassword allows an authenticated user to update the password.
func (c *ChangePasswordUseCase) ChangePassword(ctx context.Context, req ChangePasswordRequest) (*ChangePasswordResponse, error) {
	if err := govalidator.Validate.Struct(req); err != nil {
		return nil, err
	}
	if req.OldPassword == req.NewPassword {
		return nil, errors.Wrap(ErrInvalidPassword, "cannot be the same as the old password")
	}
	if req.NewPassword != req.ConfirmPassword {
		return nil, ErrInvalidPassword
	}
	ok, err := c.repo.UpdateUserPassword(req.ContextUserID, req.NewPassword)
	return &ChangePasswordResponse{
		Success: ok,
	}, err
}
