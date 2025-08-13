package main

import (
	"context"
	"errors"
	"flag"
	"net/http"
	"os"
	"os/signal"
	"rec-vendor-api/internal/config"
	"rec-vendor-api/internal/controller"
	"syscall"
	"time"

	log "github.com/sirupsen/logrus"

	"bitbucket.org/plaxieappier/rec-go-kit/logkit"

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

	r := gin.Default()
	r.GET("/healthz", controller.HealthCheck)

	s := &http.Server{
		Addr:    "0.0.0.0:8080",
		Handler: r,
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
			log.Fatalf("Failed to listen and serve http server, err: %v", err)
		}
	}()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit
	log.Info("Shutting down server ...")
}
