package metrics

import (
	"context"
	"rec-vendor-api/internal/telemetry"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/status"
)

func UnaryServerInterceptor(metrics *telemetry.GrpcPromMetrics) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
		startTime := time.Now()
		resp, err := handler(ctx, req)
		duration := time.Since(startTime).Seconds()

		code := status.Code(err).String()
		requestInfo := telemetry.RequestInfoFromContext(ctx)

		metrics.ServerHandledHistogram.WithLabelValues(
			requestInfo.MethodName,
			code,
			requestInfo.SiteID,
			requestInfo.OID,
		).Observe(duration)

		return resp, err
	}
}
