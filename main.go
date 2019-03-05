package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"go.uber.org/zap"

	"github.com/alextanhongpin/logging/config"
	"github.com/alextanhongpin/logging/pkg/grace"
	"github.com/alextanhongpin/logging/pkg/logger"
	"github.com/alextanhongpin/logging/router"
)

func main() {
	cfg := config.New()
	svc := fmt.Sprintf("mathsvc-%s", cfg.Hostname)

	log := logger.New(cfg.Env, svc)
	defer log.Sync()
	zap.ReplaceGlobals(log)

	r := router.New(cfg)

	shutdown := grace.New(r, cfg.Port)

	// Listen to the os signal for ctrl+c.
	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	// TODO: Close other dependencies here.
	// var wg sync.WaitGroup
	// wg.Add(1)
	// wg.Wait()
	shutdown()
}
