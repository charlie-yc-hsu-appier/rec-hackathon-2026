package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/signal"
	"rec-vendor-api/internal/config"
	"rec-vendor-api/internal/controller"
	"runtime/debug"
	"syscall"
	"time"

	log "github.com/sirupsen/logrus"

	"bitbucket.org/plaxieappier/rec-go-kit/logkit"
	"bitbucket.org/plaxieappier/rec-go-kit/tracekit"

	"github.com/gin-gonic/gin"
)

func main() {
	var cf = flag.String("c", "", "config file")
	flag.Parse()

	cfg := &config.Config{}
	if err := config.Load(*cf, cfg); err != nil {
		log.Fatalf("Failed to load config, err: %v", err)
	}
	logkit.InitLogging(cfg.Logging, &logkit.BaseLogFormat{})

	// Init tracer
	shutdownFunc := initTracer(cfg.Tracing)
	defer func() {
		if err := shutdownFunc(context.Background()); err != nil {
			log.Errorf("Fail to shutdown tracer provider, err: %v.", err)
		}
	}()

	r := gin.New()

	if cfg.EnableGinLogger {
		r.Use(gin.Logger())
	}

	if cfg.Logging.Format == "json" {
		r.Use(gin.RecoveryWithWriter(io.Discard, jsonRecoveryHandler))
	} else {
		r.Use(gin.Recovery())
	}

	r.GET("/healthz", controller.HealthCheck)

	addr := "0.0.0.0:8080"
	s := &http.Server{
		Addr:    addr,
		Handler: tracekit.OtelHTTPHandler(r, "vendor-api", cfg.Tracing),
	}
	defer func() {
		timeoutCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		if err := s.Shutdown(timeoutCtx); err != nil {
			log.Fatalf("Failed to shutdown server, err: %v", err)
		}
	}()
	go func() {
		if err := s.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("Failed to listen and serve http server on %s, err: %v", addr, err)
		}
	}()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit
	log.Info("Shutting down server ...")
}

func initTracer(cfg tracekit.Config) func(context.Context) error {
	if !cfg.Enable {
		return func(context.Context) error { return nil }
	}

	shutdownFunc, err := tracekit.InitProvider(cfg)
	if err != nil {
		log.Errorf("Fail to initialize tracer provider, err: %v", err)
		return func(context.Context) error { return nil }
	}

	return shutdownFunc
}

func jsonRecoveryHandler(ctx *gin.Context, recovered any) {
	log.WithContext(ctx).WithField("stack", string(debug.Stack())).Error(fmt.Sprintf("%v", recovered))
	ctx.AbortWithStatus(http.StatusInternalServerError)
}
