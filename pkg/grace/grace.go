package grace

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"
)

// New returns a new shutdown function given a http.Handler function.
func New(handler http.Handler, port string, shutdown func()) func(context.Context) {
	srv := &http.Server{
		Addr:           fmt.Sprintf(":%s", port),
		Handler:        handler,
		WriteTimeout:   10 * time.Second,
		ReadTimeout:    10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	srv.RegisterOnShutdown(shutdown)
	idleConnsClosed := make(chan struct{})
	go func() {
		log.Printf("listening to port *:%s\n", port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s", err)
		}
		close(idleConnsClosed)
	}()

	return func(ctx context.Context) {
		log.Println("shutting down")
		if err := srv.Shutdown(ctx); err != nil {
			log.Fatal("server shutdown:", err)
		}
		select {
		case <-idleConnsClosed:
			log.Println("shutdown gracefully")
			return
		case <-ctx.Done():
			log.Println("shutdown abruptly after 5 seconds timeout")
			return
		}
	}
}
