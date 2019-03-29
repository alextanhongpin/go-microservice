// Attempted the following name for package:
// - authenticator: this sounds more like a verb
// - authentication: too long
// - userlogin: is too specific, since user can also register
// - loginUser: breaks the convention, since package name is preferable a noun.
// - authz and authn is better.

package authnsvc

import (
	"context"

	"github.com/pkg/errors"

	"github.com/alextanhongpin/go-microservice/pkg/govalidator"
	"github.com/alextanhongpin/passwd"
)

type (
	// LoginRequest ... (means self-explanatory)
	LoginRequest struct {
		Username string `json:"username" validate:"required,email"`
		Password string `json:"password" validate:"required,min=8"`
	}
	// LoginResponse ...
	LoginResponse struct {
		User User `json:"user"`
	}
	// interfaces are lowercase - clients have to implement them
	// themselves.
	loginRepository interface {
		WithEmail(email string) (User, error)
	}
	// LoginUseCase ...
	LoginUseCase struct {
		users loginRepository
	}
)

// NewLoginUseCase returns a new use case for login.
func NewLoginUseCase(users loginRepository) *LoginUseCase {
	return &LoginUseCase{
		users: users,
	}
}

// Login checks if the user is authenticated.
func (l *LoginUseCase) Login(ctx context.Context, req LoginRequest) (*LoginResponse, error) {
	if err := govalidator.Validate.Struct(req); err != nil {
		return nil, errors.Wrap(err, "validate login request failed")
	}
	user, err := l.users.WithEmail(req.Username)
	if err != nil {
		return nil, errors.Wrap(err, "get user failed")
	}
	if err := passwd.Verify(req.Password, user.HashedPassword); err != nil {
		return nil, errors.Wrap(err, "verify password failed")
	}
	return &LoginResponse{user}, errors.Wrap(err, "verify password failed")
}
