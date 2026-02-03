package context

import (
	"context"
	"rec-vendor-api/internal/telemetry"
	"strings"

	schema "github.com/plaxieappier/rec-schema/go/vendorapi"
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
		ctx = setRequestInfo(ctx, info.FullMethod, req)
		return handler(ctx, req)
	}
}

func setRequestInfo(ctx context.Context, fullMethod string, req any) context.Context {
	requestInfo := telemetry.RequestInfo{
		MethodName: getMethodName(fullMethod),
	}

	extractTraceID(ctx, &requestInfo)
	extractMetadataHeaders(ctx, &requestInfo)
	setDefaultMetadataValues(&requestInfo)
	extractRequestParams(req, &requestInfo)

	return telemetry.RequestInfoToContext(ctx, requestInfo)
}

func extractTraceID(ctx context.Context, info *telemetry.RequestInfo) {
	spanCtx := trace.SpanContextFromContext(ctx)
	if spanCtx.HasTraceID() && spanCtx.IsSampled() {
		info.TraceID = spanCtx.TraceID().String()
	}
}

func extractMetadataHeaders(ctx context.Context, info *telemetry.RequestInfo) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return
	}

	info.SiteID = strings.Join(md.Get(HeaderSiteID), "")
	info.OID = strings.Join(md.Get(HeaderOID), "")
	info.BidObjID = strings.Join(md.Get(HeaderBidObjID), "")
	info.ReqID = strings.Join(md.Get(HeaderReqID), "")
}

func setDefaultMetadataValues(info *telemetry.RequestInfo) {
	if info.SiteID == "" {
		info.SiteID = "unknown"
	}
	if info.OID == "" {
		info.OID = "unknown"
	}
}

func extractRequestParams(req any, info *telemetry.RequestInfo) {
	if getRecommendationsReq, ok := req.(*schema.GetRecommendationsRequest); ok {
		info.VendorKey = getRecommendationsReq.VendorKey
		info.SubID = getRecommendationsReq.Subid
	}
}

func getMethodName(fullMethod string) string {
	fullMethod = strings.TrimPrefix(fullMethod, "/")
	if before, _, found := strings.Cut(fullMethod, "/"); found {
		return before
	}

	return "unknown"
}
