package main

import (
	"context"
	"database/sql"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	validator "gopkg.in/go-playground/validator.v9"

	"github.com/alextanhongpin/go-microservice/api"
	"github.com/alextanhongpin/go-microservice/api/middleware"
	"github.com/alextanhongpin/go-microservice/config"
	"github.com/alextanhongpin/go-microservice/database"
	"github.com/alextanhongpin/go-microservice/pkg/grace"
	"github.com/alextanhongpin/go-microservice/pkg/logger"
	"github.com/alextanhongpin/go-microservice/pkg/passport"
	"github.com/alextanhongpin/go-microservice/pkg/ratelimit"
	"github.com/alextanhongpin/go-microservice/service/authn"
	"github.com/alextanhongpin/go-microservice/service/health"
)

type Shutdown func(ctx context.Context)

func main() {
	var shutdowns []Shutdown
	cfg := config.New()

	// Create a namespace for the service running.
	log := logger.New(cfg.Env,
		zap.String("app", cfg.App),
		zap.String("host", cfg.Hostname))

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

	bearerAuthorizer := middleware.BearerAuthorizer(signer)
	basicAuthorizer := middleware.BasicAuthorizer(cfg.Credential)

	// Health endpoint.
	{
		ctl := health.NewController(cfg)
		r.GET("/health", ctl.GetHealth)
		r.GET("/protected", bearerAuthorizer, middleware.RoleChecker(api.RoleUser), ctl.GetHealth)
		r.GET("/basic", basicAuthorizer, ctl.GetHealth)
	}

	// Authentication endpoint.
	{
		opt := authn.Option{
			Signer:    signer,
			Repo:      authn.NewRepository(db),
			Validator: validate,
		}
		svc := authn.NewService(opt)
		ctl := authn.NewController(svc, signer)

		// Endpoint throttled.
		var (
			interval     = ratelimit.Per(time.Minute, 5) // 5 req/minute
			burst        = 1
			limiter      = ratelimit.New(interval, burst)
			every        = 1 * time.Minute
			expiresAfter = 1 * time.Minute
		)
		shutdown := limiter.CleanupVisitor(every, expiresAfter)
		shutdowns = append(shutdowns, shutdown)

		throttled := r.Group("/", middleware.RateLimiter(limiter))
		throttled.POST("/login", ctl.PostLogin)
		throttled.POST("/register", ctl.PostRegister)
	}

	// Books endpoint with multiple roles.
	{
		// roles := api.Roles{
		//         // The scopes should be exposed per api.
		//         api.RoleAdmin: []string{"read:books", "create:books", "update:books", "delete:books"},
		//         api.RoleOwner: []string{"read:books", "create:books", "delete:books"},
		// }
		// auth := r.Group("/v1/books", middleware.BearerAuthorizer(signer))
		// auth.GET("", middleware.RoleChecker(roles.Can("read:books")...), ctl.GetBooks)
		// auth.POST("", middleware.RoleChecker(roles.Can("create:books")...), ctl.PostBooks)
		// auth.UPDATE("", middleware.RoleChecker(roles.Can("update:books")...), ctl.UpdateBooks)
		// auth.DELETE("", middleware.RoleChecker(roles.Can("delete:books")...), ctl.DeleteBooks)
		// // Endpoint with custom action.
		// auth.POST("/:id/book:approve", ctl.ApproveBooks)
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
	shutdowns = append(shutdowns, shutdown)

	// Listen to the os signal for ctrl+c.
	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	// Create a global context cancellation to orchestrate graceful
	// shutdown for different services.
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var wg sync.WaitGroup
	wg.Add(len(shutdowns))
	for _, shutdown := range shutdowns {
		go func(shutdown Shutdown) {
			defer wg.Done()
			shutdown(ctx)
		}(shutdown)
	}
	wg.Wait()
	log.Info("terminating")
}
