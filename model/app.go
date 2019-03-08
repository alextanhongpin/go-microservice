package model

import (
	"database/sql"
	"log"

	"github.com/alextanhongpin/go-microservice/config"
	"github.com/alextanhongpin/go-microservice/pkg/signer"
	"go.uber.org/zap"
	validator "gopkg.in/go-playground/validator.v9"
)

type App struct {
	Config    *config.Config
	Validator *validator.Validate
	Database  *sql.DB
	Logger    *zap.Logger
	Signer    signer.Signer
}

func (a *App) Shutdown() {
	err := a.Database.Close()
	a.Logger.Error("close db failed", zap.Error(err))
	err = a.Logger.Sync()
	if err != nil {
		log.Println("sync logger failed", err)
	}
}
