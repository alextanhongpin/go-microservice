package infrastructure

import (
	"context"
	"database/sql"
	"errors"
	"net/http"
	"sync"
	"time"

	"github.com/alextanhongpin/go-microservice/api"
	"github.com/alextanhongpin/go-microservice/api/middleware"
	"github.com/alextanhongpin/go-microservice/pkg/logger"
	"github.com/alextanhongpin/pkg/gojwt"
	"github.com/alextanhongpin/pkg/grace"
	"github.com/alextanhongpin/pkg/requestid"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/pprof"
	"github.com/gin-gonic/gin"
	"github.com/rs/xid"
	"go.uber.org/zap"
)

type Infrastructure struct {
	config *Config
	db     *sql.DB
	// router *gin.Engine
	// ? Is supervisor a better naming?
	shutdowns grace.Shutdowns
	signer    gojwt.Signer
	logger    *zap.Logger

	// ! Some infrastructure should only be created once per struct. But
	// this leads to a massive "once" fields. This fields are commonly
	// called enablers, which differentiates themselves from factories,
	// where multiple instances can be created.
	onceConfig sync.Once
	onceDB     sync.Once
	// onceRouter   sync.Once
	onceShutdown sync.Once
	onceSigner   sync.Once
	onceLogger   sync.Once
}

func New() *Infrastructure {
	// Our infrastructure managed all the infras
	// shutdown.
	var infra = &Infrastructure{
		shutdowns: make(grace.Shutdowns, 0),
	}
	return infra
}

func (i *Infrastructure) Logger() *zap.Logger {
	i.onceLogger.Do(func() {
		cfg := i.Config()
		log := logger.New(cfg.Env,
			zap.String("app", cfg.Name),
			zap.String("host", cfg.Hostname))
		i.OnShutdown(func(ctx context.Context) {
			log.Sync()
		})
		// We are replacing the global logger here. Since logging happens at
		// all level, it will be a little pointless to pass down the logger
		// through dependency injection to all levels. You may still do that if
		// that is your preferred way of working.
		zap.ReplaceGlobals(log)
		i.logger = log
	})
	return i.logger
}

func (i *Infrastructure) OnShutdown(fn grace.Shutdown) {
	i.shutdowns.Append(fn)
}

func (i *Infrastructure) Shutdown() {
	i.onceShutdown.Do(func() {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		// Close all goroutines and the http server gracefully.
		i.shutdowns.Close(ctx)
	})
}

func (i *Infrastructure) Config() *Config {
	i.onceConfig.Do(func() {
		i.config = NewConfig()
	})
	return i.config
}

func (i *Infrastructure) Database() *sql.DB {
	i.onceDB.Do(func() {
		cfg := i.Config()
		db := NewProductionDatabase(Option(cfg.Database))

		db.SetMaxOpenConns(10)
		db.SetMaxIdleConns(5)
		db.SetConnMaxLifetime(time.Hour)

		i.OnShutdown(func(ctx context.Context) {
			db.Close()
		})
		i.db = db
	})
	return i.db
}

func (i *Infrastructure) Signer() gojwt.Signer {
	i.onceSigner.Do(func() {
		var (
			cfg          = i.Config()
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
		i.signer = gojwt.New(opt)
	})
	return i.signer
}

// TODO: Separate the router from the infrastructure. The router, controllers
// etc belongs to the application. This is a factory.
func (i *Infrastructure) Router() *gin.Engine {
	r := gin.New()

	r.Use(gin.Recovery())
	r.Use(cors.Default())

	// Custom middlewares.
	r.Use(middleware.Logger(i.Logger(), time.RFC3339, true))
	r.Use(middleware.RequestIDProvider(
		requestid.RequestID(func() (string, error) {
			return xid.New().String(), nil
		}),
	))
	r.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusNotFound, gin.H{
			"code":    "PAGE_NOT_FOUND",
			"message": "Page not found",
		})
	})
	pprof.Register(r)
	return r
}

// Factories for repositories, use cases etc should be created here.
// func (i *Infrastructure) UserUseCase() {
// }
