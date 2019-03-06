package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"go.uber.org/zap"

	"github.com/alextanhongpin/go-microservice/config"
	"github.com/alextanhongpin/go-microservice/pkg/grace"
	"github.com/alextanhongpin/go-microservice/pkg/logger"
	"github.com/alextanhongpin/go-microservice/router"
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

	// TODO: Setup database.
	// {
	//         db, err := database.NewProduction()
	//         if err != nil {
	//                 log.Fatal(err)
	//         }
	//         defer db.Close()
	//         db.SetMaxOpenConns(10)
	//         db.SetMaxIdleConns(5)
	//         db.SetConnMaxLifetime(time.Hour)
	// }
	// Define service facade here.

	// Router takes in the service facade, and orchestrate the endpoints.
	r := router.New(cfg)

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
