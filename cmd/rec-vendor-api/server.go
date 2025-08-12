package main

import (
	"context"
	"errors"
	"flag"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"log"
	"net/http"
	"os"
	"os/signal"
	"path"
	"syscall"
	"time"
)

type Config struct {
	Test string `mapstructure:"test"`
}

func loadConfig(configPath string, cfg *Config) error {
	configName := path.Base(configPath)
	ext := path.Ext(configPath)
	dir := path.Dir(configPath)

	if configPath == "" {
		return errors.New("config path is empty")
	}
	if ext != ".yaml" {
		return errors.New("only accept .yaml file")
	}

	v := viper.New()
	v.SetConfigType("yaml")
	v.SetConfigName(configName)
	v.AddConfigPath(dir)

	if err := v.ReadInConfig(); err != nil {
		return err
	}

	// Unmarshal the configuration into the struct
	if err := v.Unmarshal(cfg); err != nil {
		return err
	}

	return nil
}

func main() {
	// Parse config
	var cf = flag.String("c", "", "config file")
	flag.Parse()

	cfg := &Config{}
	if err := loadConfig(*cf, cfg); err != nil {
		log.Fatalf("Failed to load config, err: %v", err)
	}
	log.Printf("== %v ==", cfg)

	// Example GIN service
	r := gin.Default()
	r.GET("/healthz", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "ok",
		})
	})

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
	log.Println("Shutting down server ...")
}
