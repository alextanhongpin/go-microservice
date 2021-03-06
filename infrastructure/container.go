package infrastructure

import (
	"context"
	"database/sql"
	"html/template"
	"net/http"
	"sync"
	"time"

	"github.com/alextanhongpin/go-microservice/pkg/database"
	"github.com/alextanhongpin/go-microservice/pkg/logger"
	"github.com/alextanhongpin/go-microservice/presentation/middleware"
	"github.com/alextanhongpin/pkg/gojwt"
	"github.com/alextanhongpin/pkg/grace"
	"github.com/alextanhongpin/pkg/requestid"

	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/pprof"
	"github.com/gin-gonic/gin"
	"github.com/rs/xid"
	"go.uber.org/zap"
)

// Container wraps all the infrastructure components together.
type Container struct {
	// This is required to coordinate shutdown?
	// wg sync.WaitGroup
	configOnce sync.Once
	config     *Config

	dbOnce sync.Once
	db     *sql.DB
	// router *gin.Engine
	// ? Is supervisor a better naming?
	shutdownOnce sync.Once
	shutdowns    grace.Shutdowns

	signerOnce sync.Once
	signer     gojwt.Signer

	loggerOnce sync.Once
	logger     *zap.Logger

	template *template.Template

	// ! Some infrastructure should only be created once per struct. But
	// this leads to a massive "once" fields. This fields are commonly
	// called enablers, which differentiates themselves from factories,
	// where multiple instances can be created.
	// onceRouter   sync.Once
}

// NewContainer returns a new infrastructure container.
func NewContainer() *Container {
	// Our infrastructure managed all the infras shutdown.
	return &Container{
		shutdowns: make(grace.Shutdowns, 0),
	}
}

// Logger returns a new logger instance.
func (c *Container) Logger() *zap.Logger {
	c.loggerOnce.Do(func() {
		cfg := c.Config()
		log := logger.New(cfg.Env,
			zap.String("app", cfg.Name),
			zap.String("host", cfg.Hostname))
		c.OnShutdown(func(ctx context.Context) {
			log.Sync()
		})
		// We are replacing the global logger here. Since logging happens at
		// all level, it will be a little pointless to pass down the logger
		// through dependency injection to all levels. You may still do that if
		// that is your preferred way of working.
		zap.ReplaceGlobals(log)
		c.logger = log
	})
	return c.logger
}

// OnShutdown adds a new shutdown method to the supervisor.
func (c *Container) OnShutdown(fn grace.Shutdown) {
	c.shutdowns.Append(fn)
}

// Shutdown gracefully terminates all the infrastructure within the given
// context duration.
func (c *Container) Shutdown() {
	c.shutdownOnce.Do(func() {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		// Close all goroutines and the http server gracefully.
		c.shutdowns.Close(ctx)
	})
}

// Config returns a new Config that reads from the environment variables.
func (c *Container) Config() *Config {
	c.configOnce.Do(func() {
		c.config = NewConfig()
	})
	return c.config
}

// Database returns a new pointer to the database instance.
func (c *Container) Database() *sql.DB {
	c.dbOnce.Do(func() {
		cfg := c.Config()
		db := database.NewProduction(database.Option(cfg.Database))

		db.SetMaxOpenConns(5)
		db.SetMaxIdleConns(1)
		// This should be forever.
		// db.SetConnMaxLifetime(time.Hour)
		c.OnShutdown(func(ctx context.Context) {
			db.Close()
		})
		c.db = db
	})
	return c.db
}

func (c *Container) Signer() gojwt.Signer {
	c.signerOnce.Do(func() {
		c.signer = NewSigner(c.Config())
	})
	return c.signer
}

// TODO: Separate the router from the infrastructure. The router, controllers
// etc belongs to the application. This is a factory.
func (c *Container) Router(registerTemplate func(*template.Template)) *gin.Engine {
	r := gin.New()

	if registerTemplate != nil {
		tpl := template.New("")
		registerTemplate(tpl)
		r.SetHTMLTemplate(tpl)
	}

	r.Use(gin.Recovery())
	r.Use(cors.Default())

	// Custom middlewares.
	r.Use(middleware.Logger(c.Logger(), time.RFC3339, true))
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
