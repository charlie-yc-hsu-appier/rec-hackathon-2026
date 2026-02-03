package middleware

import (
	"context"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type allValidator interface {
	ValidateAll() error
}

type validator interface {
	Validate() error
}

func ValidationUnaryInterceptor(ctx context.Context, req any, _ *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
	if v, ok := req.(allValidator); ok {
		if err := v.ValidateAll(); err != nil {
			return nil, status.Errorf(codes.InvalidArgument, "Validation failed: %v", err)
		}
	} else if v, ok := req.(validator); ok {
		if err := v.Validate(); err != nil {
			return nil, status.Errorf(codes.InvalidArgument, "Validation failed: %v", err)
		}
	}

	return handler(ctx, req)
}
