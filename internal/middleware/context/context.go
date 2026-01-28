package context

import (
	"context"
	"rec-vendor-api/internal/telemetry"
	"strings"

	"go.opentelemetry.io/otel/trace"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

const (
	headerSiteID   = "x-rec-siteid"
	headerOID      = "x-rec-oid"
	headerBidObjID = "x-rec-bidobjid"
	headerReqID    = "x-request-id"
)

func UnaryServerInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
		ctx = setRequestInfo(ctx)
		return handler(ctx, req)
	}
}

func setRequestInfo(ctx context.Context) context.Context {
	requestInfo := telemetry.RequestInfo{}

	// Extract trace ID from OpenTelemetry span context
	spanCtx := trace.SpanContextFromContext(ctx)
	if spanCtx.HasTraceID() && spanCtx.IsSampled() {
		requestInfo.TraceID = spanCtx.TraceID().String()
	}

	// Extract metadata from gRPC context
	if md, ok := metadata.FromIncomingContext(ctx); ok {
		requestInfo.SiteID = strings.Join(md.Get(headerSiteID), "")
		requestInfo.OID = strings.Join(md.Get(headerOID), "")
		requestInfo.BidObjID = strings.Join(md.Get(headerBidObjID), "")
		requestInfo.ReqID = strings.Join(md.Get(headerReqID), "")
	}

	// Extract peer information (IP address) if available
	// Note: This is informational and not part of RequestInfo structure,
	// but we can log it if needed in the future

	return telemetry.RequestInfoToContext(ctx, requestInfo)
}
