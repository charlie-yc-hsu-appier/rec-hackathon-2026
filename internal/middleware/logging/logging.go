package logging

import (
	"context"
	"fmt"
	"rec-vendor-api/internal/telemetry"

	grpc_logging "github.com/grpc-ecosystem/go-grpc-middleware/logging"
	grpc_logrus "github.com/grpc-ecosystem/go-grpc-middleware/logging/logrus"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

func UnaryServerInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
		resp, err := handler(ctx, req)
		code := grpc_logging.DefaultErrorToCode(err)
		level := grpc_logrus.DefaultCodeToLevel(code)
		requestInfo := telemetry.RequestInfoFromContext(ctx)
		message := fmt.Sprintf("%s finished unary call with code %s", requestInfo.MethodName, code.String())
		messageProducer(ctx, message, level)
		return resp, err
	}
}

func messageProducer(ctx context.Context, message string, level log.Level) {
	switch level {
	case log.DebugLevel:
		log.WithContext(ctx).Debug(message)
	case log.InfoLevel:
		log.WithContext(ctx).Info(message)
	case log.WarnLevel:
		log.WithContext(ctx).Warning(message)
	case log.ErrorLevel:
		log.WithContext(ctx).Error(message)
	case log.FatalLevel:
		log.WithContext(ctx).Fatal(message)
	case log.PanicLevel:
		log.WithContext(ctx).Panic(message)
	}
}
