package logging

import (
	"context"
	"fmt"
	"strings"
	"time"

	grpc_logging "github.com/grpc-ecosystem/go-grpc-middleware/logging"
	grpc_logrus "github.com/grpc-ecosystem/go-grpc-middleware/logging/logrus"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

func UnaryServerInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
		startTime := time.Now()
		resp, err := handler(ctx, req)
		methodName := getMethodName(info.FullMethod)
		code := grpc_logging.DefaultErrorToCode(err)
		level := grpc_logrus.DefaultCodeToLevel(code)
		message := fmt.Sprintf("%s finished unary call with code %s, duration: %v", methodName, code.String(), time.Since(startTime))
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

// getMethodName returns the method name from the grpc fullMethod.
// e.g. "/v1/GetProductFeatures" -> "GetProductFeatures"
func getMethodName(fullMethod string) string {
	fullMethod = strings.TrimPrefix(fullMethod, "/") // remove leading slash
	if i := strings.Index(fullMethod, "/"); i >= 0 {
		return fullMethod[i+1:]
	}

	return "unknown"
}
