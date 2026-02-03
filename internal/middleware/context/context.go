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
	HeaderSiteID      = "x-rec-siteid"
	HeaderOID         = "x-rec-oid"
	HeaderBidObjID    = "x-rec-bidobjid"
	HeaderReqID       = "x-request-id"
	HeaderRequestTs   = "x-request-ts"
	HeaderRequester   = "x-requester"
	HeaderTraceparent = "traceparent"
)

func UnaryServerInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
		ctx = setRequestInfo(ctx, info.FullMethod)
		return handler(ctx, req)
	}
}

func setRequestInfo(ctx context.Context, fullMethod string) context.Context {
	requestInfo := telemetry.RequestInfo{}

	// Extract trace ID from OpenTelemetry span context
	spanCtx := trace.SpanContextFromContext(ctx)
	if spanCtx.HasTraceID() && spanCtx.IsSampled() {
		requestInfo.TraceID = spanCtx.TraceID().String()
	}

	// Extract metadata from gRPC context
	if md, ok := metadata.FromIncomingContext(ctx); ok {
		requestInfo.SiteID = strings.Join(md.Get(HeaderSiteID), "")
		requestInfo.OID = strings.Join(md.Get(HeaderOID), "")
		requestInfo.BidObjID = strings.Join(md.Get(HeaderBidObjID), "")
		requestInfo.ReqID = strings.Join(md.Get(HeaderReqID), "")
	}

	requestInfo.MethodName = getMethodName(fullMethod)
	site := requestInfo.SiteID
	oid := requestInfo.OID
	if site == "" {
		site = "unknown"
	}
	if oid == "" {
		oid = "unknown"
	}
	requestInfo.SiteID = site
	requestInfo.OID = oid

	return telemetry.RequestInfoToContext(ctx, requestInfo)
}

func getMethodName(fullMethod string) string {
	fullMethod = strings.TrimPrefix(fullMethod, "/") // remove leading slash
	if before, _, found := strings.Cut(fullMethod, "/"); found {
		return before
	}

	return "unknown"
}
