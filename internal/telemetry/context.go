package telemetry

import (
	"context"
)

type RequestInfo struct {
	SiteID     string
	OID        string
	VendorKey  string
	SubID      string
	TraceID    string
	BidObjID   string
	ReqID      string
	MethodName string
}

type reqInfoKey struct{}

func RequestInfoToContext(ctx context.Context, requestInfo RequestInfo) context.Context {
	if ctx == nil {
		ctx = context.Background()
	}

	return context.WithValue(ctx, reqInfoKey{}, requestInfo)
}

func RequestInfoFromContext(ctx context.Context) RequestInfo {
	if ctx == nil {
		return RequestInfo{}
	}
	if requestInfo, ok := ctx.Value(reqInfoKey{}).(RequestInfo); ok {
		return requestInfo
	}

	return RequestInfo{}
}
