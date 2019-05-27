package authnimpl

import (
	"context"

	"github.com/pkg/errors"

	"github.com/alextanhongpin/go-microservice/domain/token"
	"github.com/alextanhongpin/go-microservice/domain/user"
	"github.com/alextanhongpin/go-microservice/pkg/database"
	"github.com/alextanhongpin/go-microservice/pkg/govalidator"
	"github.com/alextanhongpin/go-microservice/usecase"
)

type Authn struct {
	users        user.Repository
	userService  user.Service
	tokens       token.Repository
	tokenService token.Service
}

// token.Service - the interface that belongs to the token domain.
// service.Token - the implementation.
func New(
	users user.Repository,
	userService user.Service,
	tokens token.Repository,
	tokenService token.Service,
) *Authn {
	return &Authn{
		users:        users,
		userService:  userService,
		tokens:       tokens,
		tokenService: tokenService,
	}
}

func (a *Authn) Login(ctx context.Context, req usecase.LoginRequest) (*usecase.LoginResponse, error) {
	if err := govalidator.Validate.Struct(req); err != nil {
		return nil, ErrInvalidRequest
	}
	usr, err := a.users.WithEmail(req.Email)
	if err != nil {
		return nil, ErrInvalidUsernameOrPassword
	}
	if ok := a.userService.ComparePassword(usr, req.Password); !ok {
		return nil, ErrInvalidUsernameOrPassword
	}
	accessToken, err := a.tokenService.CreateAccessToken(usr.ID)
	return &usecase.LoginResponse{
		Data:        usr,
		AccessToken: accessToken,
	}, err
}

func (a *Authn) Register(ctx context.Context, req usecase.RegisterRequest) (*usecase.RegisterResponse, error) {
	if err := govalidator.Validate.Struct(req); err != nil {
		return nil, errors.Wrap(err, "invalid register request")
	}
	hashedPassword, err := a.userService.HashPassword(req.Password)
	if err != nil {
		return nil, err
	}
	usr, err := a.users.Create(req.Email, hashedPassword)
	if err != nil {
		if database.IsDuplicateEntry(err) {
			return nil, ErrDuplicateUser
		}
		return nil, err
	}
	accessToken, err := a.tokenService.CreateAccessToken(usr.ID)
	return &usecase.RegisterResponse{
		Data:        usr,
		AccessToken: accessToken,
	}, err
}

func (a *Authn) RecoverPassword(ctx context.Context, req usecase.RecoverPasswordRequest) (*usecase.RecoverPasswordResponse, error) {
	if err := govalidator.Validate.Struct(req); err != nil {
		return nil, err
	}
	usr, err := a.users.WithEmail(req.Email)
	if err != nil {
		return nil, ErrInvalidRequest
	}
	tkn := token.New(usr.ID)
	// Store the hashed token, and send the user the unhashed. When
	// getting, hash the token provided by the user and compare with the
	// one in the database. This way, if the tokens in the database has
	// been leaked out, the token on the client side cannot be compromised.
	success, err := a.tokens.Create(tkn.UserID, tkn.Hashed())
	return &usecase.RecoverPasswordResponse{
		Token:   tkn.Token,
		Success: success,
	}, err
}

func (a *Authn) ResetPassword(ctx context.Context, req usecase.ResetPasswordRequest) (*usecase.ResetPasswordResponse, error) {
	if err := govalidator.Validate.Struct(req); err != nil {
		return nil, err
	}
	if req.Password != req.ConfirmPassword {
		return nil, ErrInvalidPassword
	}

	tkn, err := a.tokens.WithValue(token.HashToken(req.Token))
	if err != nil {
		return nil, err
	}
	if tkn.HasExpired(token.TTL) {
		// Delete the token from the table.
		_, _ = a.tokens.Delete(token.HashToken(req.Token))
		return nil, ErrTokenExpired
	}

	// Hash the password before storing.
	hashedPassword, err := a.userService.HashPassword(req.Password)
	if err != nil {
		return nil, err
	}
	success, err := a.users.ChangePassword(tkn.UserID, hashedPassword)
	if err != nil {
		return nil, err
	}
	// Password should not be the same as the old password.
	if !success {
		return nil, ErrInvalidPassword
	}
	// Remove the one-time used only token.
	success, err = a.tokens.Delete(token.HashToken(req.Token))
	return &usecase.ResetPasswordResponse{
		Success: success,
	}, err
}

func (a *Authn) ChangePassword(ctx context.Context, req usecase.ChangePasswordRequest) (*usecase.ChangePasswordResponse, error) {
	if err := govalidator.Validate.Struct(req); err != nil {
		return nil, err
	}
	usr, err := a.users.WithID(req.ContextUserID)
	if err != nil {
		return nil, err
	}
	// Check if the provided password is correct.
	if ok := a.userService.ComparePassword(usr, req.OldPassword); !ok {
		return nil, ErrInvalidRequest
	}

	// Hashed the new password.
	hashedPassword, err := a.userService.HashPassword(req.NewPassword)
	if err != nil {
		return nil, err
	}
	success, err := a.users.ChangePassword(req.ContextUserID, hashedPassword)
	return &usecase.ChangePasswordResponse{
		Success: success,
	}, err
}
