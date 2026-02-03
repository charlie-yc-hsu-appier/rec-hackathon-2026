package telemetry

import (
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

const (
	systemName = "vendor-api"
)

var (
	histogramBucket = []float64{0.01, 0.03, 0.05, 0.07, 0.09, 0.1, 0.2, 0.3, 0.4, 0.5, 0.75, 1}

	Metrics = newExternalPromMetrics()
)

type externalPromMetrics struct {
	RestApiDurationSeconds *prometheus.HistogramVec
	RestApiErrorTotal      *prometheus.CounterVec
	RestApiAnomalyTotal    *prometheus.CounterVec
}

func newExternalPromMetrics() externalPromMetrics {
	m := externalPromMetrics{}
	m.RestApiDurationSeconds = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Subsystem: systemName,
			Name:      "rest_api_duration_seconds",
			Help:      "Time spent calling Rest API",
			Buckets:   histogramBucket,
		}, []string{"vendor", "site", "oid"},
	)
	m.RestApiErrorTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Subsystem: systemName,
			Name:      "rest_api_error_total",
			Help:      "Error count when calling Rest API",
		}, []string{"vendor", "site", "oid"},
	)
	m.RestApiAnomalyTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Subsystem: systemName,
			Name:      "rest_api_anomaly_total",
			Help:      "Anomaly count when calling Rest API",
		}, []string{"vendor", "site", "oid", "reason"},
	)
	return m
}

func PromHandler() gin.HandlerFunc {
	h := promhttp.Handler()

	return func(c *gin.Context) {
		h.ServeHTTP(c.Writer, c.Request)
	}
}

var (
	grpcHistogramBucket = []float64{0.01, 0.05, 0.1, 0.2, 0.3, 0.4, 0.5, 0.6, 0.7, 0.8, 0.9, 1.0, 2.0}

	GrpcMetrics = newGrpcPromMetrics()
)

type GrpcPromMetrics struct {
	ServerHandledHistogram *prometheus.HistogramVec
}

func newGrpcPromMetrics() GrpcPromMetrics {
	m := GrpcPromMetrics{}
	m.ServerHandledHistogram = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Subsystem: systemName,
			Name:      "grpc_server_handled_duration_seconds",
			Help:      "request latency by status code",
			Buckets:   grpcHistogramBucket,
		}, []string{"grpc_method", "grpc_code", "site", "oid"},
	)
	return m
}
