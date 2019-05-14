package application

import (
	"database/sql"

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
	Router() *gin.Engine
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

func (m *Manager) NewAuthnRepository() *authn.Repository {
	return authn.NewRepository(m.Database())
}

func (m *Manager) NewAuthnUseCase() *authn.UseCase {
	return authn.NewUseCase(m.NewAuthnRepository(), m.Signer())
}

func (m *Manager) NewAuthnController() *authn.Controller {
	return authn.NewController(m.NewAuthnUseCase())
}
