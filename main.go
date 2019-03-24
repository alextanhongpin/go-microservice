package main

import (
	"context"
	"database/sql"
	"errors"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/rs/xid"
	"go.uber.org/zap"

	"github.com/alextanhongpin/go-microservice/api"
	"github.com/alextanhongpin/go-microservice/api/middleware"
	"github.com/alextanhongpin/go-microservice/config"
	"github.com/alextanhongpin/go-microservice/database"
	"github.com/alextanhongpin/go-microservice/pkg/logger"
	"github.com/alextanhongpin/go-microservice/service/authnsvc"
	"github.com/alextanhongpin/go-microservice/service/health"
	"github.com/alextanhongpin/go-microservice/service/usersvc"
	"github.com/alextanhongpin/pkg/gojwt"
	"github.com/alextanhongpin/pkg/grace"
	"github.com/alextanhongpin/pkg/ratelimiter"
	"github.com/alextanhongpin/pkg/requestid"
)

func main() {
	var shutdowns grace.Shutdowns
	cfg := config.New()

	// Create a namespace for the service running.
	log := logger.New(cfg.Env,
		zap.String("app", cfg.Name),
		zap.String("host", cfg.Hostname))

	defer log.Sync()

	// We are replacing the global logger here. Since logging happens at
	// all level, it will be a little pointless to pass down the logger
	// through dependency injection to all levels. You may still do that if
	// that is your preferred way of working.
	zap.ReplaceGlobals(log)

	var db *sql.DB
	{
		// Type conversion.
		db = database.NewProduction(database.Option(cfg.Database))
		defer db.Close()

		db.SetMaxOpenConns(10)
		db.SetMaxIdleConns(5)
		db.SetConnMaxLifetime(time.Hour)
	}
	var signer gojwt.Signer
	{
		var (
			audience     = cfg.Audience
			issuer       = cfg.Issuer
			semver       = cfg.Semver
			secret       = cfg.Secret
			expiresAfter = 10080 * time.Minute // 1 Week.
			scope        = api.ScopeDefault.String()
			role         = api.RoleGuest.String()
		)
		opt := gojwt.Option{
			Secret:       []byte(secret),
			ExpiresAfter: expiresAfter,
			DefaultClaims: &gojwt.Claims{
				Semver: semver,
				Scope:  scope,
				Role:   role,
				StandardClaims: jwt.StandardClaims{
					Audience: audience,
					Issuer:   issuer,
				},
			},
			Validator: func(c *gojwt.Claims) error {
				if c.Semver != semver ||
					c.Issuer != issuer ||
					c.Audience != audience {
					return errors.New("invalid token")
				}
				return nil
			},
		}
		signer = gojwt.New(opt)
	}

	r := gin.New()

	// Setup middlewares.
	r.Use(gin.Recovery())
	r.Use(cors.Default())

	// Custom middlewares.
	{
		provider := requestid.RequestID(func() (string, error) {
			return xid.New().String(), nil
		})
		r.Use(middleware.RequestIDProvider(provider))
	}
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

		repo := authnsvc.NewRepository(db)
		createAccessTokenUseCase := authnsvc.NewCreateAccessTokenUseCase(signer)

		ctl := authnsvc.NewController(authnsvc.UseCase{
			Login:    authnsvc.NewLoginUseCase(repo, createAccessTokenUseCase),
			Register: authnsvc.NewRegisterUseCase(repo, createAccessTokenUseCase),
		})

		// Endpoint throttled.
		var (
			interval     = ratelimiter.Per(time.Minute, 12) // 1 req every 5 seconds.
			burst        = 1
			limiter      = ratelimiter.New(interval, burst)
			every        = 1 * time.Minute
			expiresAfter = 1 * time.Minute
		)
		shutdown := limiter.CleanupVisitor(every, expiresAfter)
		shutdowns.Append(shutdown)

		throttled := r.Group("/", middleware.RateLimiter(limiter))
		throttled.POST("/login", ctl.PostLogin)
		throttled.POST("/register", ctl.PostRegister)

	}
	{
		repo := usersvc.NewRepository(db)
		ctl := usersvc.NewController(usersvc.UseCase{
			UserInfo: usersvc.NewUserInfoUseCase(repo),
		})
		r.POST("/userinfo", bearerAuthorizer, ctl.PostUserInfo)
		// r.GET("/users/:userID", basicAuthorizer.ctl.GetUsers)
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
	shutdowns.Append(shutdown)

	<-grace.Signal()
	// // Listen to the os signal for CTLR + C.
	// quit := make(chan os.Signal)
	// signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	// <-quit

	// Create a global context cancellation to orchestrate graceful
	// shutdown for different services.
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	shutdowns.Close(ctx)
}
