package authenticator

import (
	"github.com/pkg/errors"
	validator "gopkg.in/go-playground/validator.v9"

	"github.com/alextanhongpin/passwd"

	"github.com/alextanhongpin/go-microservice/pkg/signer"
)

type (
	Option struct {
		Repo      Repository
		Validator *validator.Validate
		Signer    signer.Signer
	}
	Service interface {
		Login(LoginRequest) (*LoginResponse, error)
		Register(RegisterRequest) (*RegisterResponse, error)
		CreateAccessToken(user, role, scope string) (string, error)
	}
	ServiceImpl struct {
		opt Option
	}
)

// New returns a new Authenticator service.
func New(opt Option) *ServiceImpl {
	return &ServiceImpl{opt}
}

type (
	LoginRequest struct {
		Username string `json:"username" validate:"email,required"`
		Password string `json:"password" validate:"required"`
	}
	LoginResponse struct {
		Data User `json:"data"`
	}
)

// Login fulfils the User Login Use Case.
// As a User,
// I want to login into the application.
func (s *ServiceImpl) Login(req LoginRequest) (*LoginResponse, error) {
	if err := s.opt.Validator.Struct(req); err != nil {
		return nil, errors.Wrap(err, "validate login request failed")
	}
	user, err := s.opt.Repo.GetUser(req.Username)
	if err != nil {
		return nil, errors.Wrap(err, "query user failed")
	}
	err = passwd.Verify(req.Password, user.HashedPassword)
	return &LoginResponse{user}, errors.Wrap(err, "verify password failed")
}

type (
	RegisterRequest struct {
		Username string `json:"username" validate:"email,required"`
		Password string `json:"password" validate:"required"`
	}
	RegisterResponse struct {
		Data User `json:"data"`
	}
)

// Register fulfils the Register User Use Case:
// As a User,
// I want to register a new User Account,
// In order to gain access to the application.
func (s *ServiceImpl) Register(req RegisterRequest) (*RegisterResponse, error) {
	if err := s.opt.Validator.Struct(req); err != nil {
		return nil, errors.Wrap(err, "validate register request failed")
	}
	// NOTE: There's no checking if the user exists, because there should
	// be a constraint in the database that the username/email is unique.
	hashedPassword, err := passwd.Hash(req.Password)
	if err != nil {
		return nil, errors.Wrap(err, "hash password failed")
	}
	user, err := s.opt.Repo.CreateUser(req.Username, hashedPassword)
	return &RegisterResponse{user}, errors.Wrap(err, "create user failed")
}

// CreateAccessToken fulfils the Authenticate User Use Case:
// As a User,
// I want to obtain a token,
// When I successfully login the system.
func (s *ServiceImpl) CreateAccessToken(user, role, scope string) (string, error) {
	claims := s.opt.Signer.NewClaims(user, role, scope)
	accessToken, err := s.opt.Signer.Sign(claims)
	return accessToken, errors.Wrap(err, "sign token failed")
}
