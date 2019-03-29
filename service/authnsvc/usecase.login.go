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
	LoginRequest struct {
		Username string `json:"username" validate:"required,email"`
		Password string `json:"password" validate:"required,min=8"`
	}
	LoginResponse struct {
		AccessToken string `json:"access_token"`
	}
	// interfaces are lowercase - clients have to implement them
	// themselves.
	loginRepository interface {
		WithEmail(email string) (User, error)
	}
	LoginUseCase struct {
		users loginRepository
		createAccessTokenUseCase
	}
)

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
	token, err := l.CreateAccessToken(user.ID)
	return &LoginResponse{token}, errors.Wrap(err, "create access token failed")
}

func NewLoginUseCase(users loginRepository, createAccessToken createAccessTokenUseCase) *LoginUseCase {
	return &LoginUseCase{
		users:                    users,
		createAccessTokenUseCase: createAccessToken,
	}
}
