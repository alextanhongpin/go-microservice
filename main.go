package main

import (
	"context"
	"database/sql"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	validator "gopkg.in/go-playground/validator.v9"

	"github.com/alextanhongpin/go-microservice/api"
	"github.com/alextanhongpin/go-microservice/api/controller"
	"github.com/alextanhongpin/go-microservice/api/middleware"
	"github.com/alextanhongpin/go-microservice/config"
	"github.com/alextanhongpin/go-microservice/database"
	"github.com/alextanhongpin/go-microservice/pkg/grace"
	"github.com/alextanhongpin/go-microservice/pkg/logger"
	"github.com/alextanhongpin/go-microservice/pkg/passport"
	"github.com/alextanhongpin/go-microservice/service/authn"
)

func main() {
	cfg := config.New()

	// Create a namespace for the service running.
	log := logger.New(cfg.Env, "your_app", cfg.Hostname)
	defer log.Sync()

	// We are replacing the global logger here. Since logging happens at
	// all level, it will be a little pointless to pass down the logger
	// through dependency injection to all levels. You may still do that if
	// that is your preferred way of working.
	zap.ReplaceGlobals(log)

	var db *sql.DB
	{
		// This will panic if the environment variables are not set.
		db = database.NewProduction()
		defer db.Close()

		db.SetMaxOpenConns(10)
		db.SetMaxIdleConns(5)
		db.SetConnMaxLifetime(time.Hour)
	}

	signer := passport.New(passport.Option{
		Secret:            []byte(cfg.Secret),
		DurationInMinutes: 10080 * time.Minute,
		Issuer:            cfg.Issuer,
		Audience:          cfg.Audience,
		Semver:            cfg.Semver,
	})

	validate := validator.New()

	r := gin.New()

	// Setup middlewares.
	r.Use(gin.Recovery())
	r.Use(cors.Default())

	// Custom middlewares.
	r.Use(middleware.RequestIDProvider())
	r.Use(middleware.Logger(log, time.RFC3339, true))
	// TODO: Include authorization passport.

	// Health endpoint.
	{
		ctl := controller.NewHealth(cfg)
		r.GET("/health", ctl.GetHealth)
		r.GET("/protected", middleware.BearerAuthorizer(signer, api.RoleUser), ctl.GetHealth)
		r.GET("/basic", middleware.BasicAuthorizer(cfg.Credential), ctl.GetHealth)
	}

	// Authentication endpoint.
	{
		opt := authn.Option{
			Signer:    signer,
			Repo:      authn.NewRepository(db),
			Validator: validate,
		}
		svc := authn.New(opt)
		ctl := controller.NewAuthn(svc, signer)
		// TODO: Throttle the login and register endpoint.
		r.POST("/login", ctl.PostLogin)
		r.POST("/register", ctl.PostRegister)
	}

	// Books endpoint with multiple roles.
	{
		// roles := api.Roles{
		//         // The scopes should be exposed per api.
		//         api.RoleAdmin: []string{"read:books", "create:books", "update:books", "delete:books"},
		//         api.RoleOwner: []string{"read:books", "create:books", "delete:books"},
		// }
		// auth := middleware.BearerAuthorizer
		// r.GET("/books", auth(signer, roles.Can("read:books")...), ctl.GetBooks)
		// r.POST("/books", auth(signer, roles.Can("create:books")...), ctl.PostBooks)
		// r.UPDATE("/books", auth(signer, roles.Can("update:books")...), ctl.UpdateBooks)
		// r.DELETE("/books", auth(signer, roles.Can("delete:books")...), ctl.DeleteBooks)
		// // Endpoint with custom action.
		// r.POST("/books:approve", auth(signer), ctl.ApproveBooks)
	}
	// Handle no route.
	r.NoRoute(func(c *gin.Context) {
		// TODO: Cleanup message.
		c.JSON(http.StatusNotFound, gin.H{
			"code":    "PAGE_NOT_FOUND",
			"message": "Page not found",
		})
	})
	// Graceful shutdown for the server.
	shutdown := grace.New(r, cfg.Port)

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
