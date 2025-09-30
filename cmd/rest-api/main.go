package main

import (
	"context"
	"fmt"
	"github.com/anpotashev/mpdgo/internal/api/middleware"
	v1 "github.com/anpotashev/mpdgo/internal/api/v1"
	"github.com/anpotashev/mpdgo/pkg/mpdapi"
	"github.com/gorilla/mux"
	"log/slog"
	"net/http"
	"os"
	"time"
)

func main() {
	ctx := context.Background()
	//logHandler := slog.NewJSONHandler(
	//	os.Stdout,
	//	&slog.HandlerOptions{
	//		Level: slog.LevelDebug,
	//	})
	logHandler := &contextHandler{slog.NewJSONHandler(
		os.Stdout,
		&slog.HandlerOptions{
			Level: slog.LevelDebug,
		})}
	//tintLogHandler := tint.NewHandler(os.Stdout, &tint.Options{Level: slog.LevelDebug, TimeFormat: time.Kitchen})
	//logHandler := tint.NewHandler(os.Stdout, &tint.Options{Level: slog.LevelDebug, TimeFormat: time.Kitchen})
	log := slog.New(logHandler)
	slog.SetDefault(log)
	mpdapi.SetLogger(log)
	//liblog := slog.New(tint.NewHandler(os.Stdout, &tint.Options{Level: slog.LevelDebug, TimeFormat: time.Kitchen}))
	//mpdapi.SetLogger(liblog)
	api, err := mpdapi.NewMpdApi(ctx, "192.168.0.110", 6600, "12345678", true, 100, 3, time.Millisecond*200, time.Second*10)

	if err != nil {
		panic(err)
	}
	router := mux.NewRouter()
	router.Use(middleware.LoggerContextMiddleware)
	v1.New(router.PathPrefix("/v1").Subrouter(), api)
	srv := http.Server{
		Addr:    fmt.Sprintf(":%d", 8080),
		Handler: router,
	}
	panic(srv.ListenAndServe())
}

type contextHandler struct {
	slog.Handler
}

func (contextHandler *contextHandler) Handle(ctx context.Context, r slog.Record) error {
	if rId, ok := ctx.Value(middleware.RequestIdContextAttributeName).(string); ok && rId != "" {
		r.AddAttrs(slog.String("request_id", rId))
	}
	return contextHandler.Handler.Handle(ctx, r)
}
