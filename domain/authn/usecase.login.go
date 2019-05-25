package authn

import (
	"context"

	"github.com/pkg/errors"

	"github.com/alextanhongpin/go-microservice/pkg/govalidator"
	"github.com/alextanhongpin/passwd"
)

// Request/response.
type (
	// LoginRequest ... (means self-explanatory)
	LoginRequest struct {
		Username string `json:"username" validate:"required,email"`
		Password string `json:"password" validate:"required,min=8"`
	}
	// LoginResponse ...
	LoginResponse struct {
		Data User `json:"data"`
	}
)

type (
	// interfaces are lowercase - clients have to implement them
	// themselves.
	loginRepository interface {
		UserWithEmail(email string) (User, error)
	}
	loginUseCase interface {
		Login(ctx context.Context, req LoginRequest) (*LoginResponse, error)
	}
)

type LoginUseCase struct {
	repo loginRepository
}

// NewLoginUseCase returns a new use case for login.
func NewLoginUseCase(repo loginRepository) *LoginUseCase {
	return &LoginUseCase{
		repo: repo,
	}
}

// Login checks if the user is authenticated.
func (l *LoginUseCase) Login(ctx context.Context, req LoginRequest) (*LoginResponse, error) {
	if err := govalidator.Validate.Struct(req); err != nil {
		return nil, errors.Wrap(err, "invalid request")
	}
	user, err := l.repo.UserWithEmail(req.Username)
	if err != nil {
		return nil, ErrInvalidUsernameOrPassword
	}
	if err := passwd.Verify(req.Password, user.HashedPassword); err != nil {
		return nil, ErrInvalidUsernameOrPassword
	}
	return &LoginResponse{user}, err
}
