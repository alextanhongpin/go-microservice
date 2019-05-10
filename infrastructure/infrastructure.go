package infrastructure

import (
	"context"
	"database/sql"
	"net/http"
	"sync"
	"time"

	"github.com/alextanhongpin/go-microservice/api/middleware"
	"github.com/alextanhongpin/go-microservice/domain/usersvc"
	"github.com/alextanhongpin/go-microservice/pkg/logger"
	"github.com/alextanhongpin/pkg/gojwt"
	"github.com/alextanhongpin/pkg/grace"
	"github.com/alextanhongpin/pkg/requestid"

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

// New returns a new infrastructure container.
func New() *Infrastructure {
	// Our infrastructure managed all the infras shutdown.
	return &Infrastructure{
		shutdowns: make(grace.Shutdowns, 0),
	}
}

// Logger returns a new logger instance.
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

// OnShutdown adds a new shutdown method to the supervisor.
func (i *Infrastructure) OnShutdown(fn grace.Shutdown) {
	i.shutdowns.Append(fn)
}

// Shutdown gracefully terminates all the infrastructure within the given
// context duration.
func (i *Infrastructure) Shutdown() {
	i.onceShutdown.Do(func() {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		// Close all goroutines and the http server gracefully.
		i.shutdowns.Close(ctx)
	})
}

// Config returns a new Config that reads from the environment variables.
func (i *Infrastructure) Config() *Config {
	i.onceConfig.Do(func() {
		i.config = NewConfig()
	})
	return i.config
}

// Database returns a new pointer to the database instance.
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
		i.signer = NewSigner(i.Config())
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

// UserRepository returns a new UserRepository.
func (i *Infrastructure) UserRepository() *usersvc.Repository {
	return usersvc.NewRepository(i.Database())
}

// UserService returns a new UserService.
func (i *Infrastructure) UserService() *usersvc.Service {
	return usersvc.NewService(i.UserRepository())
}
