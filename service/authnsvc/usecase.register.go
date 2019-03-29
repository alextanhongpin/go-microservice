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

	"github.com/alextanhongpin/go-microservice/database"
	"github.com/alextanhongpin/go-microservice/pkg/govalidator"
	"github.com/alextanhongpin/passwd"
)

type (
	// RegisterRequest ...
	RegisterRequest struct {
		Username string `json:"username" validate:"required,email"`
		Password string `json:"password" validate:"required,min=8"`
	}
	// RegisterResponse ...
	RegisterResponse struct {
		User User `json:"user"`
	}
	registerRepository interface {
		Create(username, password string) (User, error)
	}
	// RegisterUseCase ...
	RegisterUseCase struct {
		users registerRepository
	}
)

// NewRegisterUseCase returns a new use case to register user.
func NewRegisterUseCase(users registerRepository) *RegisterUseCase {
	return &RegisterUseCase{
		users: users,
	}
}

// Register creates a new account for new users.
func (r *RegisterUseCase) Register(ctx context.Context, req RegisterRequest) (*RegisterResponse, error) {
	if err := govalidator.Validate.Struct(req); err != nil {
		return nil, errors.Wrap(err, "validate register request failed")
	}
	// NOTE: There's no checking if the user exists, because there should
	// be a constraint in the database that the username/email is unique.
	hashedPassword, err := passwd.Hash(req.Password)
	if err != nil {
		return nil, errors.Wrap(err, "hash password failed")
	}
	user, err := r.users.Create(req.Username, hashedPassword)
	if err != nil {
		if database.IsDuplicateEntry(err) {
			return nil, errors.New("user already exists")
		}
		return nil, errors.Wrap(err, "create user failed")
	}
	return &RegisterResponse{user}, errors.Wrap(err, "register user failed")
}
