package application

import (
	"database/sql"
	"html/template"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"github.com/alextanhongpin/pkg/gojwt"
	"github.com/alextanhongpin/pkg/grace"

	"github.com/alextanhongpin/go-microservice/implementation/authnimpl"
	"github.com/alextanhongpin/go-microservice/implementation/tokenimpl"
	"github.com/alextanhongpin/go-microservice/infrastructure"
	"github.com/alextanhongpin/go-microservice/infrastructure/repository"
	"github.com/alextanhongpin/go-microservice/usecase"
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

func (m *Manager) NewAuthnUseCase() usecase.Authn {
	users := repository.NewUser(m.Database())
	userService := authnimpl.NewService()

	tokens := repository.NewToken(m.Database())
	tokenService := tokenimpl.NewService(m.Signer())
	return authnimpl.New(users, userService, tokens, tokenService)
}

// func (m *Manager) NewViews(tpl *template.Template) {
//         authn.NewView(tpl)
// }
