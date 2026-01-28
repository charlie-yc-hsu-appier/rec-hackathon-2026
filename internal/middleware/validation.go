package middleware

import (
	"context"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func ValidationUnaryInterceptor(ctx context.Context, req any, _ *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
	if v, ok := req.(interface{ ValidateAll() error }); ok {
		if err := v.ValidateAll(); err != nil {
			return nil, status.Errorf(codes.InvalidArgument, "Validation failed: %v", err)
		}
	} else if v, ok := req.(interface{ Validate() error }); ok {
		if err := v.Validate(); err != nil {
			return nil, status.Errorf(codes.InvalidArgument, "Validation failed: %v", err)
		}
	}

	return handler(ctx, req)
}
