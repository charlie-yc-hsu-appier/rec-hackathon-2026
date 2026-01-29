package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"net"

	"io"
	"net/http"
	"net/netip"
	"os"
	"os/signal"
	"runtime/debug"
	"syscall"
	"time"

	"rec-vendor-api/internal/config"
	"rec-vendor-api/internal/controller"
	logFormat "rec-vendor-api/internal/logformat"
	"rec-vendor-api/internal/middleware"
	grpc_context "rec-vendor-api/internal/middleware/context"
	grpc_logging "rec-vendor-api/internal/middleware/logging"
	"rec-vendor-api/internal/telemetry"
	"rec-vendor-api/internal/vendor"

	"github.com/gin-gonic/gin"
	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_prometheus "github.com/grpc-ecosystem/go-grpc-middleware/providers/prometheus"
	grpc_recovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	grpc_realip "github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/realip"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"github.com/plaxieappier/rec-go-kit/logkit"
	"github.com/plaxieappier/rec-go-kit/tracekit"
	"github.com/prometheus/client_golang/prometheus"
	log "github.com/sirupsen/logrus"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/keepalive"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"

	schema "github.com/plaxieappier/rec-schema/go/vendorapi"
)

const (
	systemName = "vendor-api"
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
	var cf = flag.String("c", "", "config file")
	var appType = flag.String("t", "", "app type: gin, grpc, all, default: all")
	flag.Parse()

	cfg := &config.Config{}
	if err := config.Load(*cf, cfg); err != nil {
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

	vendorRegistry, err := vendor.BuildRegistry(cfg.VendorConfig)
	if err != nil {
		log.Fatalf("Failed to build vendor registry, err: %v", err)
	}

	var ginServer *http.Server
	var grpcServer *grpc.Server
	var gatewayServer *http.Server

	// Load port configuration from environment variables
	portConfig := &config.PortConfig{}
	config.LoadConfigFromEnv(portConfig)

	grpcAddr := "0.0.0.0:" + portConfig.GrpcPort
	gatewayAddr := "0.0.0.0:" + portConfig.GatewayPort
	ginAddr := "0.0.0.0:" + portConfig.GinPort
	appTypeStr := *appType
	switch appTypeStr {
	case "gin":
		ginServer = initGinServer(cfg, vendorRegistry, ginAddr)
	case "grpc":
		grpcServer = initGRPCServer(cfg, vendorRegistry, grpcAddr)
		gatewayServer = initGatewayServer(grpcAddr, gatewayAddr)
	default:
		ginServer = initGinServer(cfg, vendorRegistry, ginAddr)
		grpcServer = initGRPCServer(cfg, vendorRegistry, grpcAddr)
		gatewayServer = initGatewayServer(grpcAddr, gatewayAddr)
	}

	// Setup graceful shutdown for all started servers
	if ginServer != nil {
		defer func() {
			timeoutCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
			defer cancel()
			if err := ginServer.Shutdown(timeoutCtx); err != nil {
				log.Errorf("Failed to shutdown gin server, err: %v", err)
			}
		}()
	}
	if grpcServer != nil {
		defer grpcServer.GracefulStop()
	}
	if gatewayServer != nil {
		defer func() {
			shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 10*time.Second)
			defer shutdownCancel()
			if err := gatewayServer.Shutdown(shutdownCtx); err != nil {
				log.Errorf("Failed to shutdown gateway server, err: %v", err)
			}
		}()
	}

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit
	log.Info("Shutting down server ...")
}

func initGinServer(cfg *config.Config, vendorRegistry map[string]vendor.Client, addr string) *http.Server {
	log.Infof("Starting gin server on %s", addr)
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

	recommender := controller.NewRecommender(vendorRegistry)
	vendorManager := controller.NewVendorManager(cfg.VendorConfig)

	r.GET("/r/:vendor_key", recommender.Recommend)
	r.GET("/vendors", vendorManager.GetVendors)
	r.GET("/healthz", controller.HealthCheck)
	r.GET("/metrics", telemetry.PromHandler())

	s := &http.Server{
		Addr:    addr,
		Handler: tracekit.OtelHTTPHandler(r, "vendor-api", cfg.Tracing),
	}
	go func() {
		if err := s.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("Failed to listen and serve http server on %s, err: %v", addr, err)
		}
	}()
	return s
}

func initGRPCServer(cfg *config.Config, vendorRegistry map[string]vendor.Client, grpcAddr string) *grpc.Server {
	log.Infof("Starting grpc server on %s", grpcAddr)
	handler, err := controller.NewHandler(vendorRegistry, cfg.VendorConfig)
	if err != nil {
		log.Fatalf("Failed to initialize grpc handler: %v", err)
	}

	recoveryFunc := func(p any) (err error) {
		return status.Errorf(codes.Internal, "panic triggered: %v", p)
	}
	recoveryOpts := []grpc_recovery.Option{
		grpc_recovery.WithRecoveryHandler(recoveryFunc),
	}

	// Trust all proxies to use X-Forwarded-For header
	// since we do not know the client's IP address, we trust all proxies.
	trustedPeers := []netip.Prefix{
		netip.MustParsePrefix("0.0.0.0/0"),
		netip.MustParsePrefix("::/0"),
	}

	grpcMetrics := initGRPCMetrics()
	grpcServer := grpc.NewServer(
		grpc.StatsHandler(otelgrpc.NewServerHandler()),
		grpc.UnaryInterceptor(grpc_middleware.ChainUnaryServer(
			grpc_recovery.UnaryServerInterceptor(recoveryOpts...),
			grpc_realip.UnaryServerInterceptor(trustedPeers, []string{grpc_realip.XForwardedFor}),
			grpc_context.UnaryServerInterceptor(),
			grpc_logging.UnaryServerInterceptor(),
			middleware.ValidationUnaryInterceptor,
			grpcMetrics.UnaryServerInterceptor(
				getMetricInterceptorOpt()...,
			),
		)),
		grpc.KeepaliveParams(keepalive.ServerParameters{
			MaxConnectionAge: cfg.Grpc.MaxConnectionAge,
		}),
		grpc.WriteBufferSize(cfg.Grpc.WriteBufferSizeKb*1024),
		grpc.ReadBufferSize(cfg.Grpc.ReadBufferSizeKb*1024),
	)
	grpcMetrics.InitializeMetrics(grpcServer)
	schema.RegisterVendorAPIServer(grpcServer, handler)

	reflection.Register(grpcServer)

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

	return grpcServer
}

func initGatewayServer(grpcAddr string, gatewayAddr string) *http.Server {
	log.Infof("Starting gateway server on %s", gatewayAddr)
	gatewayMux := runtime.NewServeMux(runtime.WithIncomingHeaderMatcher(func(key string) (string, bool) {
		if _, ok := headerMatcher[http.CanonicalHeaderKey(key)]; ok {
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

	go func() {
		log.Infof("Serving gateway server on %v", gatewayAddr)
		if err := gatewayServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("Failed to listen and serve gateway server on %v, err: %v", gatewayAddr, err)
		}
	}()

	return gatewayServer
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

func initGRPCMetrics() *grpc_prometheus.ServerMetrics {
	grpcMetrics := grpc_prometheus.NewServerMetrics(
		grpc_prometheus.WithServerHandlingTimeHistogram(
			grpc_prometheus.WithHistogramSubsystem(systemName),
			grpc_prometheus.WithHistogramBuckets([]float64{
				0.01, 0.05, 0.1, 0.2, 0.3, 0.4, 0.5, 0.6, 0.7, 0.8, 0.9, 1.0, 2.0,
			}),
		),
		grpc_prometheus.WithServerCounterOptions(
			grpc_prometheus.WithSubsystem(systemName),
		),
		grpc_prometheus.WithContextLabels("site", "oid"),
	)
	prometheus.MustRegister(grpcMetrics)
	return grpcMetrics
}

func getMetricInterceptorOpt() []grpc_prometheus.Option {
	return []grpc_prometheus.Option{
		grpc_prometheus.WithLabelsFromContext(func(ctx context.Context) prometheus.Labels {
			requestInfo := telemetry.RequestInfoFromContext(ctx)
			site := requestInfo.SiteID
			oid := requestInfo.OID
			if site == "" {
				site = "unknown"
			}
			if oid == "" {
				oid = "unknown"
			}
			return prometheus.Labels{
				"site": site,
				"oid":  oid,
			}
		}),
	}
}
