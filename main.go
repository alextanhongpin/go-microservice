package main

import (
	"time"

	"github.com/alextanhongpin/go-microservice/api"
	"github.com/alextanhongpin/go-microservice/api/middleware"
	"github.com/alextanhongpin/go-microservice/domain/authnsvc"
	"github.com/alextanhongpin/go-microservice/domain/health"
	"github.com/alextanhongpin/go-microservice/domain/usersvc"
	"github.com/alextanhongpin/go-microservice/infrastructure"
	"github.com/alextanhongpin/pkg/grace"
	"github.com/alextanhongpin/pkg/ratelimiter"
)

func main() {
	infra := infrastructure.New()
	// Gracefully terminate all dependencies, together with the server.
	defer infra.Shutdown()

	var (
		signer = infra.Signer()
		cfg    = infra.Config()
		r      = infra.Router()
		db     = infra.Database()
	)

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
		svc := authnsvc.NewService(repo, signer)
		ctl := authnsvc.NewController(svc)

		// Endpoint throttled.
		var (
			interval     = ratelimiter.Per(time.Minute, 12) // 1 req every 5 seconds.
			burst        = 1
			limiter      = ratelimiter.New(interval, burst)
			every        = 1 * time.Minute
			expiresAfter = 1 * time.Minute
		)
		shutdown := limiter.CleanupVisitor(every, expiresAfter)
		infra.OnShutdown(shutdown)

		throttled := r.Group("/", middleware.RateLimiter(limiter))
		throttled.POST("/login", ctl.PostLogin)
		throttled.POST("/register", ctl.PostRegister)

	}
	{
		repo := usersvc.NewRepository(db)
		svc := usersvc.NewService(repo)
		ctl := usersvc.NewController(svc)
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

	// Starts a new server with the given port. Returns a shutdown method
	// for the server.
	shutdown := grace.New(r, cfg.Port)

	// Coordinate server shutdown with the infrastructure dependencies.
	infra.OnShutdown(shutdown)

	// Listen to the os signal for CTLR + C.
	<-grace.Signal()
}
