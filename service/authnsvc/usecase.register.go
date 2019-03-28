// Attempted the following name for package:
// - authenticator: this sounds more like a verb
// - authentication: too long
// - userlogin: is too specific, since user can also register
// - loginUser: breaks the convention, since package name is preferable a noun.
// - authz and authn is better.

package authnsvc

import (
	"github.com/pkg/errors"

	"github.com/alextanhongpin/go-microservice/database"
	"github.com/alextanhongpin/go-microservice/service"
	"github.com/alextanhongpin/passwd"
)

type (
	RegisterRequest struct {
		Username string `json:"username" validate:"required,email"`
		Password string `json:"password" validate:"required,min=8"`
	}
	RegisterResponse struct {
		AccessToken string `json:"access_token"`
	}
	RegisterUseCase    func(req RegisterRequest) (*RegisterResponse, error)
	RegisterRepository interface {
		Create(username, password string) (User, error)
	}
)

func NewRegisterUseCase(
	users RegisterRepository,
	createAccessToken CreateAccessTokenUseCase,
) RegisterUseCase {
	return func(req RegisterRequest) (*RegisterResponse, error) {
		if err := service.Validate.Struct(req); err != nil {
			return nil, errors.Wrap(err, "validate register request failed")
		}
		// NOTE: There's no checking if the user exists, because there should
		// be a constraint in the database that the username/email is unique.
		hashedPassword, err := passwd.Hash(req.Password)
		if err != nil {
			return nil, errors.Wrap(err, "hash password failed")
		}
		user, err := users.Create(req.Username, hashedPassword)
		if err != nil {
			if database.IsDuplicateEntry(err) {
				return nil, errors.New("user already exists")
			}
			return nil, errors.Wrap(err, "create user failed")
		}
		token, err := createAccessToken(user.ID)
		return &RegisterResponse{token}, errors.Wrap(err, "create access token failed")
	}
}
