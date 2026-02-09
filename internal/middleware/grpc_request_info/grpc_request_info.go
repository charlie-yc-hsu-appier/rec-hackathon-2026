// TODO: to be renamed to request_info after gin/nginx retirement, and also remove request_info.go and request_info_test.go
package grpc_request_info

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
		requestInfo := buildRequestInfo(ctx, info.FullMethod, req)
		ctx = telemetry.RequestInfoToContext(ctx, requestInfo)
		return handler(ctx, req)
	}
}

func buildRequestInfo(ctx context.Context, fullMethod string, req any) telemetry.RequestInfo {
	traceID := extractTraceID(ctx)
	siteID, oid, bidObjID, reqID := extractMetadataHeaders(ctx)
	siteID = applyDefault(siteID, "unknown")
	oid = applyDefault(oid, "unknown")
	vendorKey, subID := extractRequestParams(req)

	return telemetry.RequestInfo{
		MethodName: getMethodName(fullMethod),
		TraceID:    traceID,
		SiteID:     siteID,
		OID:        oid,
		BidObjID:   bidObjID,
		ReqID:      reqID,
		VendorKey:  vendorKey,
		SubID:      subID,
	}
}

func extractTraceID(ctx context.Context) string {
	spanCtx := trace.SpanContextFromContext(ctx)
	if spanCtx.HasTraceID() && spanCtx.IsSampled() {
		return spanCtx.TraceID().String()
	}
	return ""
}

func extractMetadataHeaders(ctx context.Context) (siteID, oid, bidObjID, reqID string) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return "", "", "", ""
	}

	getFirst := func(key string) string {
		if values := md.Get(key); len(values) > 0 {
			return values[0]
		}
		return ""
	}

	siteID = getFirst(HeaderSiteID)
	oid = getFirst(HeaderOID)
	bidObjID = getFirst(HeaderBidObjID)
	reqID = getFirst(HeaderReqID)
	return siteID, oid, bidObjID, reqID
}

func applyDefault(value, defaultValue string) string {
	if value == "" {
		return defaultValue
	}
	return value
}

func extractRequestParams(req any) (vendorKey, subID string) {
	if getRecommendationsReq, ok := req.(*schema.GetRecommendationsRequest); ok {
		return getRecommendationsReq.VendorKey, getRecommendationsReq.Subid
	}
	return "", ""
}

func getMethodName(fullMethod string) string {
	fullMethod = strings.TrimPrefix(fullMethod, "/")
	if _, after, found := strings.Cut(fullMethod, "/"); found {
		return after
	}

	return "unknown"
}
