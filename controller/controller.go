package controller

import (
	"github.com/alextanhongpin/go-microservice/model"
	"github.com/alextanhongpin/go-microservice/service"
)

type Controller struct {
	Authz  *Authz
	Health *Health
}

func New(app *model.App, services *service.Service) *Controller {
	return &Controller{
		Authz:  NewAuthz(services.Auth, app.Signer),
		Health: NewHealth(app.Config),
	}
}
