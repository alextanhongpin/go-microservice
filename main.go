package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	validator "gopkg.in/go-playground/validator.v9"

	"github.com/alextanhongpin/go-microservice/config"
	"github.com/alextanhongpin/go-microservice/controller"
	"github.com/alextanhongpin/go-microservice/database"
	"github.com/alextanhongpin/go-microservice/middleware"
	"github.com/alextanhongpin/go-microservice/model"
	"github.com/alextanhongpin/go-microservice/pkg/grace"
	"github.com/alextanhongpin/go-microservice/pkg/logger"
	"github.com/alextanhongpin/go-microservice/pkg/signer"
	"github.com/alextanhongpin/go-microservice/service"
)

func main() {
	app := makeApp()
	services := service.New(app)
	controllers := controller.New(app, services)

	handler := makeHandler(app, controllers)

	// Graceful shutdown for the server.
	shutdown := grace.New(handler, app.Config.Port, app.Shutdown)

	// Listen to the os signal for ctrl+c.
	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	// TODO: Close other dependencies here.
	// var wg sync.WaitGroup
	// wg.Add(1)
	// wg.Wait()

	// Create a global context cancellation to orchestrate graceful
	// shutdown for different services.
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	shutdown(ctx)
}

func makeApp() *model.App {
	cfg := config.New()

	// Create a namespace for the service running.
	log := logger.New(cfg.Env, "your_app", cfg.Hostname)

	// We are replacing the global logger here. Since logging happens at
	// all level, it will be a little pointless to pass down the logger
	// through dependency injection to all levels. You may still do that if
	// that is your preferred way of working.
	zap.ReplaceGlobals(log)

	// This will panic if the environment variables are not set.
	db := database.NewProduction()
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(time.Hour)

	signMgr := signer.New(signer.Option{
		Secret:            []byte(cfg.Secret),
		DurationInMinutes: 10080 * time.Minute,
		Issuer:            cfg.Issuer,
		Audience:          cfg.Audience,
		Semver:            cfg.Semver,
	})

	validate := validator.New()

	return &model.App{
		Config:    cfg,
		Validator: validate,
		Database:  db,
		Logger:    log,
		Signer:    signMgr,
	}
}

func makeHandler(app *model.App, controllers *controller.Controller) http.Handler {
	r := gin.New()

	// Setup middlewares.
	r.Use(gin.Recovery())
	r.Use(cors.Default())

	// Custom middlewares.
	r.Use(middleware.RequestID())
	r.Use(middleware.Logger(zap.L(), time.RFC3339, true))
	// TODO: Include authorization signer.

	// Health endpoint.
	{
		ctl := controllers.Health
		r.GET("/health", ctl.GetHealth)
		r.GET("/protected", middleware.Authz(app.Signer), ctl.GetHealth)
	}

	// Register endpoint.
	{
		ctl := controllers.Authz
		r.POST("/login", ctl.PostLogin)
		r.POST("/register", ctl.PostRegister)
	}

	// Handle no route.
	r.NoRoute(func(c *gin.Context) {
		// TODO: Cleanup message.
		c.JSON(http.StatusNotFound, gin.H{
			"code":    "PAGE_NOT_FOUND",
			"message": "Page not found",
		})
	})
	return r
}
