package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/pkg/errors"
	"go.uber.org/zap"

	"github.com/alextanhongpin/logging/pkg/logger"
	"github.com/alextanhongpin/logging/pkg/xreqid"
)

var buildDate time.Time
var uptime time.Time

func main() {
	uptime = time.Now()
	// env := "production"
	env := "development"
	hostname, err := os.Hostname()
	if err != nil {
		log.Fatal(err)
	}

	date := os.Getenv("BUILD_DATE")
	log.Println(date)
	buildDate, err = time.Parse(time.RFC3339, date)
	if err != nil {
		log.Fatal(err)
	}

	svc := fmt.Sprintf("mathsvc-%s", hostname)
	log := logger.New(env, svc)
	defer log.Sync()
	undo := zap.ReplaceGlobals(log)
	defer undo()

	http.HandleFunc("/health", withLogging(log, withRequestID(controller)))
	log.Info("listening to port *:8080")
	http.ListenAndServe(":8080", nil)
}

func withLogging(log *zap.Logger, next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Info("middleware: start")
		next.ServeHTTP(w, r)
		log.Info("middleware: end")
	}
}

func withRequestID(next http.HandlerFunc) http.HandlerFunc {
	log := zap.L()
	return func(w http.ResponseWriter, r *http.Request) {
		reqID := xreqid.FromHeader(w, r)
		ctx := xreqid.ContextWithRequestID(r.Context(), reqID)
		log.Info("withRequestID: start")
		next.ServeHTTP(w, r.WithContext(ctx))
		log.Info("withRequestID: end")
	}
}

type Health struct {
	BuildDate time.Time `json:"build_date,omitempty"`
	GitTag    string    `json:"git_tag,omitempty"`
	Uptime    string    `json:"uptime"`
}

func controller(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	// The request id must be provided.
	reqID, _ := xreqid.FromContext(ctx)
	log := zap.L().With(logger.ReqIdField(reqID))
	err := service(ctx)
	log.Error("controller", zap.Error(err))
	// Furnish the health endpoint with necessary information.
	res := Health{
		BuildDate: buildDate,
		GitTag:    os.Getenv("TAG"),
		Uptime:    time.Since(uptime).String(),
	}
	json.NewEncoder(w).Encode(res)
}

func service(ctx context.Context) error {
	reqID, _ := xreqid.FromContext(ctx)
	log := zap.L().With(logger.ReqIdField(reqID))
	log.Info("service: start")
	repository(ctx)
	log.Info("service: end")
	// Stack trace added to this line.
	return errors.Wrap(errors.New("hello"), "service")
}

func repository(ctx context.Context) {
	reqID, _ := xreqid.FromContext(ctx)
	log := zap.L().With(logger.ReqIdField(reqID))
	log.Info("repository: start")
	// Do work.
	log.Info("repository: end")
}
