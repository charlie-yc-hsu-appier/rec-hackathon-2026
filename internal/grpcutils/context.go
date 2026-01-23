package grpcutils

import (
	"context"

	"google.golang.org/grpc/metadata"
)

// GetClientIPFromContext extracts the client IP from gRPC metadata.
// It looks for the "x-forwarded-for" header, which is typically set by proxies/load balancers.
// Not exactly sure if x-forwarded-for is the best header to use, but it's the most common one.
func GetClientIPFromContext(ctx context.Context) string {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return ""
	}
	ips := md.Get("x-forwarded-for")
	if len(ips) == 0 {
		return ""
	}
	// x-forwarded-for can contain multiple IPs separated by commas
	// The first IP is typically the original client IP
	return ips[0]
}
