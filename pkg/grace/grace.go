package grace

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"
)

func New(handler http.Handler, port string) func() {
	srv := &http.Server{
		Addr:           fmt.Sprintf(":%s", port),
		Handler:        handler,
		WriteTimeout:   10 * time.Second,
		ReadTimeout:    10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	idleConnsClosed := make(chan struct{})
	go func() {
		log.Printf("listening to port *:%s\n", port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s", err)
		}
		close(idleConnsClosed)
	}()

	return func() {
		log.Println("shutting down")
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := srv.Shutdown(ctx); err != nil {
			log.Fatal("server shutdown:", err)
		}
		select {
		case <-idleConnsClosed:
			return
		case <-ctx.Done():
			log.Println("timeout 5 seconds")
		}
		log.Println("server exiting")
	}
}
