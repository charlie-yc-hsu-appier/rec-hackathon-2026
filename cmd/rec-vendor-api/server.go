package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"net"

	"io"
	"net/http"
	"os"
	"os/signal"
	"runtime/debug"
	"syscall"
	"time"

	"rec-vendor-api/internal/config"
	"rec-vendor-api/internal/controller"
	logFormat "rec-vendor-api/internal/logformat"
	"rec-vendor-api/internal/middleware"
	"rec-vendor-api/internal/telemetry"
	"rec-vendor-api/internal/vendor"

	"github.com/gin-gonic/gin"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"github.com/plaxieappier/rec-go-kit/logkit"
	"github.com/plaxieappier/rec-go-kit/tracekit"
	log "github.com/sirupsen/logrus"

	vendor_grpc "rec-vendor-api/internal/grpc"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/reflection"

	schema "github.com/plaxieappier/rec-schema/go/vendorapi"
)

// @title Vendor API service
// @version 1.0
// @description Vendor API service for recommendation ecosystem
// @contact.email ai-rec-sys@appier.com
// @basePath /
// @schemes https
//
//go:generate swag init -d ../../ -g cmd/rec-vendor-api/server.go -o ../../docs --parseInternal --parseDependency

var headerMatcher = map[string]struct{}{
	"X-Requester":    {},
	"X-Rec-Siteid":   {},
	"X-Rec-Bidobjid": {},
	"X-Rec-Oid":      {},
	"X-Request-Id":   {},
	"X-Request-Ts":   {},
	"Traceparent":    {},
}

func main() {
	cfg, err := loadConfig()
	if err != nil {
		log.Fatalf("Failed to load config, err: %v", err)
	}
	// Init logging
	logkit.InitLogging(cfg.Logging, &logFormat.LogFormat{})

	// Init tracer
	shutdownFunc := initTracer(cfg.Tracing)
	defer func() {
		if err := shutdownFunc(context.Background()); err != nil {
			log.Errorf("Fail to shutdown tracer provider, err: %v.", err)
		}
	}()

	r := gin.New()
	// MUST be set to true for getting value from context
	r.ContextWithFallback = true

	r.Use(middleware.RequestInfo())

	if cfg.EnableGinLogger {
		r.Use(gin.Logger())
	}

	if cfg.Logging.Format == "json" {
		r.Use(gin.RecoveryWithWriter(io.Discard, jsonRecoveryHandler))
	} else {
		r.Use(gin.Recovery())
	}

	vendorRegistry, err := vendor.BuildRegistry(cfg.VendorConfig)
	if err != nil {
		log.Fatalf("Failed to build vendor registry, err: %v", err)
	}
	recommender := controller.NewRecommender(vendorRegistry)
	vendorManager := controller.NewVendorManager(cfg.VendorConfig)

	r.GET("/r/:vendor_key", recommender.Recommend)
	r.GET("/vendors", vendorManager.GetVendors)
	r.GET("/healthz", controller.HealthCheck)
	r.GET("/metrics", telemetry.PromHandler())

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

	// Start a gRPC server
	grpcServer := initGRPCServer(vendorRegistry, cfg.VendorConfig)
	grpcAddr := "0.0.0.0:10000"
	go func() {
		lis, err := net.Listen("tcp", grpcAddr)
		if err != nil {
			log.Fatalf("Failed to listen grpc server on %v, err: %v", grpcAddr, err)
		}
		log.Infof("Serving gRPC server on %v", grpcAddr)
		if err := grpcServer.Serve(lis); err != nil {
			log.Fatalf("Failed to serve gRPC server on %v, err: %v", grpcAddr, err)
		}
	}()
	defer grpcServer.GracefulStop()

	// gateway server
	gatewayAddr := "0.0.0.0:10001"
	gatewayServer := initGatewayServer(grpcAddr, gatewayAddr)
	go func() {
		log.Infof("Serving gateway server on %v", gatewayAddr)
		if err := gatewayServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("Failed to listen and serve gateway server on %v, err: %v", gatewayAddr, err)
		}
	}()
	defer func() {
		shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer shutdownCancel()
		if err := gatewayServer.Shutdown(shutdownCtx); err != nil {
			log.Fatalf("Failed to shutdown server, err: %v", err)
		}
	}()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit
	log.Info("Shutting down server ...")
}

func initGRPCServer(vendorRegistry map[string]vendor.Client, vendorConfig config.VendorConfig) *grpc.Server {
	service := vendor.NewService(vendorRegistry, vendorConfig)
	apiServer := vendor_grpc.NewAPIServer(service)

	grpcServer := grpc.NewServer()
	schema.RegisterVendorAPIServer(grpcServer, apiServer)

	// Register standard gRPC health service for Kubernetes probes
	healthServer := health.NewServer()
	grpc_health_v1.RegisterHealthServer(grpcServer, healthServer)
	// Set the service to serving status
	healthServer.SetServingStatus("", grpc_health_v1.HealthCheckResponse_SERVING)

	reflection.Register(grpcServer)
	return grpcServer
}

func initGatewayServer(grpcAddr string, gatewayAddr string) *http.Server {
	gatewayMux := runtime.NewServeMux(runtime.WithIncomingHeaderMatcher(func(key string) (string, bool) {
		if _, ok := headerMatcher[key]; ok {
			return key, true
		}
		return runtime.DefaultHeaderMatcher(key)
	}))
	gatewayOpts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}
	if err := schema.RegisterVendorAPIHandlerFromEndpoint(context.Background(), gatewayMux, grpcAddr, gatewayOpts); err != nil {
		log.Fatalf("Failed to register gRPC gateway, err: %v", err)
	}
	mux := http.NewServeMux()
	mux.Handle("/", gatewayMux)
	gatewayServer := &http.Server{
		Addr:    gatewayAddr,
		Handler: mux,
	}
	return gatewayServer
}

func loadConfig() (*config.Config, error) {
	var cf = flag.String("c", "", "config file")
	flag.Parse()

	cfg := &config.Config{}
	if err := config.Load(*cf, cfg); err != nil {
		return nil, err
	}
	return cfg, nil
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
