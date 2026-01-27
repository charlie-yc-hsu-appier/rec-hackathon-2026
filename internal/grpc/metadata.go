package grpc

import (
	"context"
	"net"
	"strings"

	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/peer"
)

// GetClientIPFromContext extracts the client IP from gRPC context metadata.
// It follows a similar algorithm to Gin's ClientIP():
// 1. Checks X-Real-IP header first (most reliable when behind a single proxy)
// 2. Checks X-Forwarded-For header and parses the first valid IP from comma-separated list
// 3. Falls back to peer address from gRPC context if headers are unavailable
func GetClientIPFromContext(ctx context.Context) string {
	md, ok := metadata.FromIncomingContext(ctx)
	if ok {
		// Check X-Real-IP first (most reliable, set by nginx as $remote_addr)
		if realIPs := md.Get("x-real-ip"); len(realIPs) > 0 {
			if ip := parseIP(realIPs[0]); ip != "" {
				return ip
			}
		}

		if forwardedFor := md.Get("x-forwarded-for"); len(forwardedFor) > 0 {
			for _, value := range forwardedFor {
				ips := strings.Split(value, ",")
				for _, ipStr := range ips {
					if ip := parseIP(ipStr); ip != "" {
						return ip
					}
				}
			}
		}
	}

	if p, ok := peer.FromContext(ctx); ok && p.Addr != nil {
		if addr, ok := p.Addr.(*net.TCPAddr); ok {
			return addr.IP.String()
		}
		// Handle other address types (e.g., Unix socket)
		addrStr := p.Addr.String()
		if host, _, err := net.SplitHostPort(addrStr); err == nil {
			return host
		}
		return addrStr
	}

	return ""
}

func parseIP(ipStr string) string {
	ipStr = strings.TrimSpace(ipStr)
	if ipStr == "" {
		return ""
	}
	if ip := net.ParseIP(ipStr); ip != nil {
		return ip.String()
	}
	return ""
}
