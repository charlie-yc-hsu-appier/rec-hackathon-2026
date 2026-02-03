package logging

import (
	"context"
	"rec-vendor-api/internal/telemetry"

	grpc_logging "github.com/grpc-ecosystem/go-grpc-middleware/logging"
	grpc_logrus "github.com/grpc-ecosystem/go-grpc-middleware/logging/logrus"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/peer"
)

const (
	MethodGetRecommendations = "GetRecommendations"
)

func UnaryServerInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
		requestInfo := telemetry.RequestInfoFromContext(ctx)
		if requestInfo.MethodName == MethodGetRecommendations {
			resp, err := handler(ctx, req)

			code := grpc_logging.DefaultErrorToCode(err)
			level := grpc_logrus.DefaultCodeToLevel(code)
			logAccessLog(ctx, requestInfo.MethodName, code.String(), level)

			return resp, err
		}
		return handler(ctx, req)
	}
}

func logAccessLog(ctx context.Context, method string, codeStr string, level log.Level) {
	remoteAddr := "unknown"
	if p, ok := peer.FromContext(ctx); ok {
		remoteAddr = p.Addr.String()
	}
	message := "access_log: " + method + " finished unary call with code " + codeStr
	log.WithContext(ctx).WithFields(log.Fields{
		"remote_addr": remoteAddr,
	}).Log(level, message)
}
