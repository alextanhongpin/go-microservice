package authsvc

import "github.com/alextanhongpin/go-microservice/model"

type (
	Repository interface {
		GetUser(username string) (*model.User, error)
		CreateUser(username, password string) (*model.User, error)
	}
	// RepositoryImpl struct{}
)
