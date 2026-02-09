// TODO: to be removed after gin/nginx retirement
package middleware

import (
	"rec-vendor-api/internal/telemetry"

	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/otel/trace"
)

const (
	headerSiteID   = "x-rec-siteid"
	headerOID      = "x-rec-oid"
	headerBidObjID = "x-rec-bidobjid"
	headerReqID    = "x-request-id"
	paramVendorKey = "vendor_key"
	querySubID     = "subid"
)

type requestInfoMiddleware struct{}

func (m *requestInfoMiddleware) apply(c *gin.Context) {
	ctx := c.Request.Context()
	requestInfo := m.buildRequestInfo(c)
	ctx = telemetry.RequestInfoToContext(ctx, requestInfo)

	c.Request = c.Request.WithContext(ctx)

	// Continue to the next middleware or handler
	c.Next()
}

func (m *requestInfoMiddleware) buildRequestInfo(c *gin.Context) telemetry.RequestInfo {
	traceId := ""
	spanCtx := trace.SpanContextFromContext(c)
	if spanCtx.HasTraceID() && spanCtx.IsSampled() {
		traceId = spanCtx.TraceID().String()
	}

	return telemetry.RequestInfo{
		SiteID:    c.GetHeader(headerSiteID),
		OID:       c.GetHeader(headerOID),
		VendorKey: c.Param(paramVendorKey),
		SubID:     c.Query(querySubID),
		TraceID:   traceId,
		BidObjID:  c.GetHeader(headerBidObjID),
		ReqID:     c.GetHeader(headerReqID),
	}
}

func RequestInfo() gin.HandlerFunc {
	middleware := requestInfoMiddleware{}
	return func(c *gin.Context) {
		middleware.apply(c)
	}
}
