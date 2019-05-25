package main

import (
	"context"
	"html/template"
	"time"

	"github.com/alextanhongpin/go-microservice/api"
	"github.com/alextanhongpin/go-microservice/api/middleware"
	"github.com/alextanhongpin/go-microservice/application"
	"github.com/alextanhongpin/go-microservice/domain/health"
	"github.com/alextanhongpin/pkg/grace"
	"github.com/alextanhongpin/pkg/ratelimiter"
)

func main() {
	app := application.NewManager()
	defer app.Shutdown()

	var (
		signer = app.Signer()
		cfg    = app.Config()
		router = app.Router(func(tpl *template.Template) {
			app.NewViews(tpl)
		})
	)
	// Register all views here.

	// Middlewares/controllers are not created by the infrastructure
	// container because they are framework dependent.
	bearerAuthorizer := middleware.BearerAuthorizer(signer)
	basicAuthorizer := middleware.BasicAuthorizer(cfg.Credential)

	// Health endpoint.
	{
		ctl := health.NewController(cfg)
		router.GET("/health", ctl.GetHealth)
		router.GET("/protected", bearerAuthorizer, middleware.RoleChecker(api.RoleUser), ctl.GetHealth)
		router.GET("/basic", basicAuthorizer, ctl.GetHealth)
	}

	// Authentication endpoint.
	{
		ctl, stopBackgroundTask := app.NewAuthnController()
		app.OnShutdown(func(ctx context.Context) {
			stopBackgroundTask()
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
		app.OnShutdown(shutdown)

		throttled := router.Group("/v1", middleware.RateLimiter(limiter))
		throttled.POST("/login", ctl.PostLogin)
		throttled.POST("/register", ctl.PostRegister)
		// throttled.POST("/password/recover", ctl.PostRegister)
		// throttled.POST("/password/reset", ctl.PostRegister)
		// throttled.POST("/password/update", ctl.PostRegister)
		// HTML views will not have the version, and the names will be singular.
		router.GET("/password/reset", ctl.GetResetPasswordView)

	}

	{
		ctl := app.NewUserController()
		router.POST("/userinfo", bearerAuthorizer, ctl.PostUserInfo)
		// router.GET("/users/:userID", basicAuthorizerouter.ctl.GetUsers)
	}

	// Books endpoint with multiple roles.
	{
		// roles := api.Roles{
		//         // The scopes should be exposed per api.
		//         api.RoleAdmin: []string{"read:books", "create:books", "update:books", "delete:books"},
		//         api.RoleOwner: []string{"read:books", "create:books", "delete:books"},
		// }
		// auth := router.Group("/v1/books", middleware.BearerAuthorizer(signer))
		// auth.GET("", middleware.RoleChecker(roles.Can("read:books")...), ctl.GetBooks)
		// auth.POST("", middleware.RoleChecker(roles.Can("create:books")...), ctl.PostBooks)
		// auth.UPDATE("", middleware.RoleChecker(roles.Can("update:books")...), ctl.UpdateBooks)
		// auth.DELETE("", middleware.RoleChecker(roles.Can("delete:books")...), ctl.DeleteBooks)
		// // Endpoint with custom action.
		// auth.POST("/:id/book:approve", ctl.ApproveBooks)
	}

	// Starts a new server with the given port. Returns a shutdown method
	// for the serverouter.

	shutdown := grace.New(router, cfg.Port)

	// Coordinate server shutdown with the infrastructure dependencies.
	app.OnShutdown(shutdown)

	// Listen to the os signal for CTRL + C.
	<-grace.Signal()
}
