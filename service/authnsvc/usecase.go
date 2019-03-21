// Attempted the following name for package:
// - authenticator: this sounds more like a verb
// - authentication: too long
// - userlogin: is too specific, since user can also register
// - loginUser: breaks the convention, since package name is preferable a noun.
// - authz and authn is better.

package authnsvc

import (
	"github.com/pkg/errors"

	"github.com/alextanhongpin/go-microservice/api"
	"github.com/alextanhongpin/go-microservice/database"
	"github.com/alextanhongpin/go-microservice/service"
	"github.com/alextanhongpin/passwd"
	"github.com/alextanhongpin/pkg/gojwt"
)

type UseCase struct {
	Login    LoginUseCase
	Register RegisterUseCase
}

type (
	LoginRequest struct {
		Username string `json:"username" validate:"required,email"`
		Password string `json:"password" validate:"required,min=8"`
	}
	LoginResponse struct {
		AccessToken string `json:"access_token"`
	}
	LoginUseCase           func(LoginRequest) (*LoginResponse, error)
	LoginUseCaseRepository interface {
		WithEmail(email string) (User, error)
	}
)

// NewLoginUseCase returns a new LoginUseCase that includes the access token
// creation use case.
func NewLoginUseCase(
	users LoginUseCaseRepository,
	createAccessToken CreateAccessTokenUseCase,
) LoginUseCase {
	return func(req LoginRequest) (*LoginResponse, error) {
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

type (
	RegisterRequest struct {
		Username string `json:"username" validate:"required,email"`
		Password string `json:"password" validate:"required,min=8"`
	}
	RegisterResponse struct {
		AccessToken string `json:"access_token"`
	}
	RegisterUseCase           func(req RegisterRequest) (*RegisterResponse, error)
	RegisterUseCaseRepository interface {
		Create(username, password string) (User, error)
	}
)

func NewRegisterUseCase(
	users RegisterUseCaseRepository,
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

type CreateAccessTokenUseCase func(user string) (string, error)

func NewCreateAccessTokenUseCase(signer gojwt.Signer) CreateAccessTokenUseCase {
	return func(user string) (string, error) {
		if len(user) == 0 {
			return "", errors.New("user is required")
		}
		accessToken, err := signer.Sign(func(c *gojwt.Claims) error {
			c.StandardClaims.Subject = user
			c.Scope = api.Scopes(api.ScopeProfile, api.ScopeOpenID)
			// TODO: Determine role based on user role.
			c.Role = api.RoleUser.String()
			return nil
		})
		return accessToken, errors.Wrap(err, "sign token failed")
	}
}
