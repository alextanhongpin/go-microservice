package application

import (
	"database/sql"
	"html/template"

	"github.com/alextanhongpin/go-microservice/domain/authn"
	"github.com/alextanhongpin/go-microservice/domain/user"
	"github.com/alextanhongpin/go-microservice/infrastructure"
	"github.com/alextanhongpin/pkg/gojwt"
	"github.com/alextanhongpin/pkg/grace"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type Infrastructure interface {
	Config() *infrastructure.Config
	Database() *sql.DB
	Logger() *zap.Logger
	OnShutdown(grace.Shutdown)
	Router(func(*template.Template)) *gin.Engine
	Shutdown()
	Signer() gojwt.Signer
}

type Manager struct {
	Infrastructure
}

func NewManager() *Manager {
	return &Manager{infrastructure.NewContainer()}
}

// NewUserRepository returns a new UserRepository.
func (m *Manager) NewUserRepository() *user.Repository {
	return user.NewRepository(m.Database())
}

// NewUserService returns a new UserService.
func (m *Manager) NewUserUseCase() *user.UseCase {
	return user.NewUseCase(m.NewUserRepository())
}

func (m *Manager) NewUserController() *user.Controller {
	return user.NewController(m.NewUserUseCase())
}

func (m *Manager) NewAuthnUseCase() (*authn.UseCase, func()) {
	repo := authn.NewRepository(m.Database())
	return authn.NewUseCase(repo, m.Signer(), m.Config().PasswordTokenTTL)
}

func (m *Manager) NewAuthnController() (*authn.Controller, func()) {
	usecase, shutdown := m.NewAuthnUseCase()
	return authn.NewController(usecase), shutdown
}
