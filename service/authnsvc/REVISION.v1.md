```go
package authnsvc

import (
	"context"

	"github.com/alextanhongpin/go-microservice/api"
	"github.com/alextanhongpin/go-microservice/database"
	"github.com/alextanhongpin/go-microservice/pkg/govalidator"
	"github.com/alextanhongpin/passwd"
	"github.com/alextanhongpin/pkg/gojwt"
	"github.com/pkg/errors"
)

type (
	repository interface {
		// Reader.
		WithEmail(email string) (User, error)

		// Writer.
		Create(username, password string) (User, error)
	}

	Service struct {
		users  repository
		signer gojwt.Signer
	}
)

func NewService(users repository, signer gojwt.Signer) *Service {
	return &Service{
		users,
		signer,
	}
}

type (
	LoginRequest struct {
		Username string `json:"username" validate:"required,email"`
		Password string `json:"password" validate:"required,min=8"`
	}
	LoginResponse struct {
		AccessToken string `json:"access_token"`
	}
)

func (s *Service) Login(ctx context.Context, req LoginRequest) (*LoginResponse, error) {
	if err := govalidator.Validate.Struct(req); err != nil {
		return nil, errors.Wrap(err, "validate login request failed")
	}
	user, err := s.users.WithEmail(req.Username)
	if err != nil {
		return nil, errors.Wrap(err, "get user failed")
	}
	if err := passwd.Verify(req.Password, user.HashedPassword); err != nil {
		return nil, errors.Wrap(err, "verify password failed")
	}
	token, err := s.CreateAccessToken(user.ID)
	return &LoginResponse{token}, errors.Wrap(err, "create access token failed")
}

type (
	RegisterRequest struct {
		Username string `json:"username" validate:"required,email"`
		Password string `json:"password" validate:"required,min=8"`
	}
	RegisterResponse struct {
		AccessToken string `json:"access_token"`
	}
)

func (s *Service) Register(ctx context.Context, req RegisterRequest) (*RegisterResponse, error) {
	if err := govalidator.Validate.Struct(req); err != nil {
		return nil, errors.Wrap(err, "validate register request failed")
	}
	// NOTE: There's no checking if the user exists, because there should
	// be a constraint in the database that the username/email is unique.
	hashedPassword, err := passwd.Hash(req.Password)
	if err != nil {
		return nil, errors.Wrap(err, "hash password failed")
	}
	user, err := s.users.Create(req.Username, hashedPassword)
	if err != nil {
		if database.IsDuplicateEntry(err) {
			return nil, errors.New("user already exists")
		}
		return nil, errors.Wrap(err, "create user failed")
	}
	token, err := s.CreateAccessToken(user.ID)
	return &RegisterResponse{token}, errors.Wrap(err, "create access token failed")
}

func (s *Service) CreateAccessToken(userID string) (string, error) {
	if err := govalidator.Validate.Varname(userID, "required"); err != nil {
		return nil, err
	}

	accessToken, err := c.signer.Sign(func(c *gojwt.Claims) error {
		c.StandardClaims.Subject = user
		c.Scope = api.Scopes(api.ScopeProfile, api.ScopeOpenID)
		// TODO: Determine role based on user role.
		c.Role = api.RoleUser.String()
		return nil
	})
	return accessToken, errors.Wrap(err, "sign token failed")
}
```
