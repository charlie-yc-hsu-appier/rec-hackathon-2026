package middleware

import (
	"rec-vendor-api/internal/telemetry"

	"github.com/gin-gonic/gin"
)

const (
	headerSiteID   = "x-rec-siteid"
	headerOID      = "x-rec-oid"
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
	return telemetry.RequestInfo{
		SiteID: c.GetHeader(headerSiteID),
		OID:    c.GetHeader(headerOID),
		Vendor: c.Param(paramVendorKey),
		SubID:  c.Query(querySubID),
	}
}

func RequestInfo() gin.HandlerFunc {
	middleware := requestInfoMiddleware{}
	return func(c *gin.Context) {
		middleware.apply(c)
	}
}
