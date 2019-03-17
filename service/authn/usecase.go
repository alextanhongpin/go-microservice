// Attempted the following name for package:
// - authenticator: this sounds more like a verb
// - authentication: too long
// - userlogin: is too specific, since user can also register
// - loginUser: breaks the convention, since package name is preferable a noun.
// - authz and authn is better.

package authn

import (
	"github.com/pkg/errors"

	"github.com/alextanhongpin/go-microservice/api"
	"github.com/alextanhongpin/go-microservice/pkg/passport"
	"github.com/alextanhongpin/go-microservice/service"
	"github.com/alextanhongpin/passwd"
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
		GetUser(email string) (User, error)
	}
)

// NewLoginUseCase returns a new LoginUseCase that includes the access token
// creation use case.
func NewLoginUseCase(
	repo LoginUseCaseRepository,
	createAccessToken CreateAccessTokenUseCase,
) LoginUseCase {
	return func(req LoginRequest) (*LoginResponse, error) {
		if err := service.Validate.Struct(req); err != nil {
			return nil, errors.Wrap(err, "validate login request failed")
		}
		user, err := repo.GetUser(req.Username)
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
		CreateUser(username, password string) (User, error)
	}
)

func NewRegisterUseCase(
	repo RegisterUseCaseRepository,
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
		user, err := repo.CreateUser(req.Username, hashedPassword)
		if err != nil {
			return nil, errors.Wrap(err, "create user failed")
		}
		token, err := createAccessToken(user.ID)
		return &RegisterResponse{token}, errors.Wrap(err, "create access token failed")
	}
}

type CreateAccessTokenUseCase func(user string) (string, error)

func NewCreateAccessTokenUseCase(signer passport.Signer) CreateAccessTokenUseCase {
	return func(user string) (string, error) {
		if len(user) == 0 {
			return "", errors.New("user is required")
		}
		role := api.RoleUser.String()
		scope := api.Scopes(api.ScopeProfile, api.ScopeOpenID)
		claims := signer.NewClaims(user, role, scope)
		accessToken, err := signer.Sign(claims)
		return accessToken, errors.Wrap(err, "sign token failed")
	}
}
