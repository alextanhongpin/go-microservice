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

	"github.com/alextanhongpin/go-microservice/service"
	"github.com/alextanhongpin/passwd"
)

type (
	LoginRequest struct {
		Username string `json:"username" validate:"required,email"`
		Password string `json:"password" validate:"required,min=8"`
	}
	LoginResponse struct {
		AccessToken string `json:"access_token"`
	}
	LoginUseCase func(context.Context, LoginRequest) (*LoginResponse, error)

	LoginRepository interface {
		WithEmail(email string) (User, error)
	}
)

// NewLoginUseCase returns a new LoginUseCase that includes the access token
// creation use case.
func NewLoginUseCase(
	users LoginRepository,
	createAccessToken CreateAccessTokenUseCase,
) LoginUseCase {
	return func(ctx context.Context, req LoginRequest) (*LoginResponse, error) {
		// TOOD: If login fails three times...
		if err := service.Validate.Struct(req); err != nil {
			return nil, errors.Wrap(err, "validate login request failed")
		}
		user, err := users.WithEmail(req.Username)
		if err != nil {
			return nil, errors.Wrap(err, "get user failed")
		}
		if err := passwd.Verify(req.Password, user.HashedPassword); err != nil {
			return nil, errors.Wrap(err, "verify password failed")
		}
		token, err := createAccessToken(user.ID)
		return &LoginResponse{token}, errors.Wrap(err, "create access token failed")
	}
}
