package telemetry

import (
	"rec-vendor-api/internal/constants"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	grpcHistogramBucket = []float64{0.01, 0.05, 0.1, 0.2, 0.3, 0.4, 0.5, 0.6, 0.7, 0.8, 0.9, 1.0, 2.0}

	// TODO: to be removed after gin/nginx retirement
	GrpcCodeToStatusMapping = map[string]string{
		"OK":                 "success",
		"Canceled":           "error",
		"Unknown":            "error",
		"InvalidArgument":    "client_error",
		"DeadlineExceeded":   "timeout",
		"NotFound":           "client_error",
		"AlreadyExists":      "client_error",
		"PermissionDenied":   "client_error",
		"ResourceExhausted":  "server_error",
		"FailedPrecondition": "client_error",
		"Aborted":            "error",
		"OutOfRange":         "client_error",
		"Unimplemented":      "server_error",
		"Internal":           "server_error",
		"Unavailable":        "server_error",
		"DataLoss":           "server_error",
		"Unauthenticated":    "client_error",
	}

	// TODO: to be removed after gin/nginx retirement, currently used to match gin labels
	GrpcMethodToPathMapping = map[string]string{
		constants.MethodGetRecommendations: "/r/:vendor_key",
		constants.MethodGetVendors:         "/vendors",
		constants.MethodHealthCheck:        "/healthz",
	}

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
			Name:      "grpc_server_handled_latency",
			Help:      "request latency by status code",
			Buckets:   grpcHistogramBucket,
		}, []string{"grpc_method", "grpc_code", "site", "oid", "status", "path"}, // remove status and path after gin/nginx retirement
	)
	return m
}

// TODO: to be removed after gin/nginx retirement
func GetStatusFromCode(code string) string {
	if status, ok := GrpcCodeToStatusMapping[code]; ok {
		return status
	}
	return "unknown"
}

// TODO: to be removed after gin/nginx retirement
// GetPathFromMethodAndVendorKey constructs the full path from method name and vendor_key
// e.g., GetRecommendations + "adforus" -> "/r/adforus"
func GetPathFromMethodAndVendorKey(methodName, vendorKey string) string {
	pathTemplate, ok := GrpcMethodToPathMapping[methodName]
	if !ok {
		return "unknown"
	}

	// For GetRecommendations, replace :vendor_key with actual value
	if methodName == constants.MethodGetRecommendations && vendorKey != "" {
		return "/r/" + vendorKey
	}

	return pathTemplate
}
