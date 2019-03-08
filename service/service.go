package service

import (
	"github.com/alextanhongpin/go-microservice/model"
	"github.com/alextanhongpin/go-microservice/service/authsvc"
)

type Service struct {
	Auth authsvc.Service
}

func New(app *model.App) *Service {
	var service Service

	{
		svc := authsvc.New(authsvc.Option{
			Repo:      authsvc.NewRepository(app.Database),
			Validator: app.Validator,
		})
		service.Auth = svc
	}

	return &service
}
