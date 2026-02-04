package grpc_request_info

import (
	"context"
	"rec-vendor-api/internal/telemetry"
	"testing"

	schema "github.com/plaxieappier/rec-schema/go/vendorapi"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

type ContextTestSuite struct {
	suite.Suite
}

func (ts *ContextTestSuite) TestUnaryServerInterceptor() {
	tt := []struct {
		name        string
		setupCtx    func() context.Context
		setupReq    func() any
		setupMock   func() (grpc.UnaryHandler, *context.Context)
		wantErr     bool
		expectedErr error
		validate    func(t *testing.T, capturedCtx *context.Context)
	}{
		{
			name: "GIVEN valid request THEN context is enriched with request info",
			setupCtx: func() context.Context {
				md := metadata.New(map[string]string{
					HeaderSiteID:   "site123",
					HeaderOID:      "oid456",
					HeaderBidObjID: "bidobj789",
					HeaderReqID:    "req-id-001",
				})
				return metadata.NewIncomingContext(context.Background(), md)
			},
			setupReq: func() any {
				return &schema.GetRecommendationsRequest{
					VendorKey: "test-vendor",
					Subid:     "test-subid",
				}
			},
			setupMock: func() (grpc.UnaryHandler, *context.Context) {
				var capturedCtx context.Context
				handler := func(ctx context.Context, req any) (any, error) {
					capturedCtx = ctx
					return &schema.GetRecommendationsResponse{}, nil
				}
				return handler, &capturedCtx
			},
			wantErr: false,
			validate: func(t *testing.T, capturedCtx *context.Context) {
				requestInfo := telemetry.RequestInfoFromContext(*capturedCtx)
				assert.Equal(t, "site123", requestInfo.SiteID)
				assert.Equal(t, "oid456", requestInfo.OID)
				assert.Equal(t, "bidobj789", requestInfo.BidObjID)
				assert.Equal(t, "req-id-001", requestInfo.ReqID)
				assert.Equal(t, "test-vendor", requestInfo.VendorKey)
				assert.Equal(t, "test-subid", requestInfo.SubID)
				assert.Equal(t, "GetRecommendations", requestInfo.MethodName)
			},
		},
		{
			name: "GIVEN handler returns error THEN error is propagated",
			setupCtx: func() context.Context {
				return context.Background()
			},
			setupReq: func() any {
				return &schema.GetRecommendationsRequest{}
			},
			setupMock: func() (grpc.UnaryHandler, *context.Context) {
				handler := func(ctx context.Context, req any) (any, error) {
					return nil, assert.AnError
				}
				return handler, nil
			},
			wantErr:     true,
			expectedErr: assert.AnError,
			validate:    func(t *testing.T, capturedCtx *context.Context) {},
		},
	}

	for _, tc := range tt {
		ts.T().Run(tc.name, func(t *testing.T) {
			t.Parallel()

			ctx := tc.setupCtx()
			req := tc.setupReq()
			handler, capturedCtx := tc.setupMock()

			interceptor := UnaryServerInterceptor()
			info := &grpc.UnaryServerInfo{
				FullMethod: "/vendorapi.VendorAPI/GetRecommendations",
			}

			_, err := interceptor(ctx, req, info, handler)

			if tc.wantErr {
				require.Error(t, err)
				assert.Equal(t, tc.expectedErr, err)
			} else {
				require.NoError(t, err)
				tc.validate(t, capturedCtx)
			}
		})
	}
}

func (ts *ContextTestSuite) TestExtractMetadataHeaders() {
	tt := []struct {
		name         string
		metadata     map[string]string
		hasMetadata  bool
		expectedInfo telemetry.RequestInfo
	}{
		{
			name: "GIVEN all headers present THEN all values extracted",
			metadata: map[string]string{
				HeaderSiteID:   "site123",
				HeaderOID:      "oid456",
				HeaderBidObjID: "bidobj789",
				HeaderReqID:    "req-id-001",
			},
			hasMetadata: true,
			expectedInfo: telemetry.RequestInfo{
				SiteID:   "site123",
				OID:      "oid456",
				BidObjID: "bidobj789",
				ReqID:    "req-id-001",
			},
		},
		{
			name:        "GIVEN no metadata THEN empty values",
			metadata:    map[string]string{},
			hasMetadata: true,
			expectedInfo: telemetry.RequestInfo{
				SiteID:   "",
				OID:      "",
				BidObjID: "",
				ReqID:    "",
			},
		},
		{
			name: "GIVEN partial headers THEN only present values extracted",
			metadata: map[string]string{
				HeaderSiteID: "site123",
				HeaderReqID:  "req-id-001",
			},
			hasMetadata: true,
			expectedInfo: telemetry.RequestInfo{
				SiteID:   "site123",
				OID:      "",
				BidObjID: "",
				ReqID:    "req-id-001",
			},
		},
		{
			name: "GIVEN empty header values THEN empty strings extracted",
			metadata: map[string]string{
				HeaderSiteID: "",
				HeaderOID:    "",
			},
			hasMetadata: true,
			expectedInfo: telemetry.RequestInfo{
				SiteID:   "",
				OID:      "",
				BidObjID: "",
				ReqID:    "",
			},
		},
		{
			name:        "GIVEN context without metadata THEN no panic and empty values",
			metadata:    nil,
			hasMetadata: false,
			expectedInfo: telemetry.RequestInfo{
				SiteID:   "",
				OID:      "",
				BidObjID: "",
				ReqID:    "",
			},
		},
	}

	for _, tc := range tt {
		ts.T().Run(tc.name, func(t *testing.T) {
			t.Parallel()

			var ctx context.Context
			if tc.hasMetadata {
				md := metadata.New(tc.metadata)
				ctx = metadata.NewIncomingContext(context.Background(), md)
			} else {
				ctx = context.Background()
			}
			info := &telemetry.RequestInfo{}

			extractMetadataHeaders(ctx, info)

			assert.Equal(t, tc.expectedInfo.SiteID, info.SiteID)
			assert.Equal(t, tc.expectedInfo.OID, info.OID)
			assert.Equal(t, tc.expectedInfo.BidObjID, info.BidObjID)
			assert.Equal(t, tc.expectedInfo.ReqID, info.ReqID)
		})
	}
}

func (ts *ContextTestSuite) TestSetDefaultMetadataValues() {
	tt := []struct {
		name         string
		inputInfo    telemetry.RequestInfo
		expectedInfo telemetry.RequestInfo
	}{
		{
			name: "GIVEN empty SiteID and OID THEN defaults set to unknown",
			inputInfo: telemetry.RequestInfo{
				SiteID: "",
				OID:    "",
			},
			expectedInfo: telemetry.RequestInfo{
				SiteID: "unknown",
				OID:    "unknown",
			},
		},
		{
			name: "GIVEN SiteID present but OID empty THEN only OID defaults",
			inputInfo: telemetry.RequestInfo{
				SiteID: "site123",
				OID:    "",
			},
			expectedInfo: telemetry.RequestInfo{
				SiteID: "site123",
				OID:    "unknown",
			},
		},
		{
			name: "GIVEN OID present but SiteID empty THEN only SiteID defaults",
			inputInfo: telemetry.RequestInfo{
				SiteID: "",
				OID:    "oid456",
			},
			expectedInfo: telemetry.RequestInfo{
				SiteID: "unknown",
				OID:    "oid456",
			},
		},
		{
			name: "GIVEN both SiteID and OID present THEN no defaults applied",
			inputInfo: telemetry.RequestInfo{
				SiteID: "site123",
				OID:    "oid456",
			},
			expectedInfo: telemetry.RequestInfo{
				SiteID: "site123",
				OID:    "oid456",
			},
		},
	}

	for _, tc := range tt {
		ts.T().Run(tc.name, func(t *testing.T) {
			t.Parallel()

			info := tc.inputInfo

			setDefaultMetadataValues(&info)

			assert.Equal(t, tc.expectedInfo.SiteID, info.SiteID)
			assert.Equal(t, tc.expectedInfo.OID, info.OID)
		})
	}
}

func (ts *ContextTestSuite) TestExtractRequestParams() {
	tt := []struct {
		name         string
		request      any
		expectedInfo telemetry.RequestInfo
	}{
		{
			name: "GIVEN GetRecommendationsRequest THEN VendorKey and SubID extracted",
			request: &schema.GetRecommendationsRequest{
				VendorKey: "test-vendor",
				Subid:     "test-subid",
			},
			expectedInfo: telemetry.RequestInfo{
				VendorKey: "test-vendor",
				SubID:     "test-subid",
			},
		},
		{
			name: "GIVEN GetRecommendationsRequest with empty values THEN empty strings extracted",
			request: &schema.GetRecommendationsRequest{
				VendorKey: "",
				Subid:     "",
			},
			expectedInfo: telemetry.RequestInfo{
				VendorKey: "",
				SubID:     "",
			},
		},
		{
			name:    "GIVEN different request type THEN no extraction",
			request: &struct{ SomeField string }{SomeField: "value"},
			expectedInfo: telemetry.RequestInfo{
				VendorKey: "",
				SubID:     "",
			},
		},
		{
			name:    "GIVEN nil request THEN no panic and empty values",
			request: nil,
			expectedInfo: telemetry.RequestInfo{
				VendorKey: "",
				SubID:     "",
			},
		},
	}

	for _, tc := range tt {
		ts.T().Run(tc.name, func(t *testing.T) {
			t.Parallel()

			info := &telemetry.RequestInfo{}

			extractRequestParams(tc.request, info)

			assert.Equal(t, tc.expectedInfo.VendorKey, info.VendorKey)
			assert.Equal(t, tc.expectedInfo.SubID, info.SubID)
		})
	}
}

func (ts *ContextTestSuite) TestGetMethodName() {
	tt := []struct {
		name               string
		fullMethod         string
		expectedMethodName string
	}{
		{
			name:               "GIVEN standard gRPC method THEN method name extracted",
			fullMethod:         "/vendorapi.VendorAPI/GetRecommendations",
			expectedMethodName: "GetRecommendations",
		},
		{
			name:               "GIVEN method without leading slash THEN method name extracted",
			fullMethod:         "vendorapi.VendorAPI/GetRecommendations",
			expectedMethodName: "GetRecommendations",
		},
		{
			name:               "GIVEN method without slash separator THEN returns unknown",
			fullMethod:         "/vendorapi.VendorAPI",
			expectedMethodName: "unknown",
		},
		{
			name:               "GIVEN empty string THEN returns unknown",
			fullMethod:         "",
			expectedMethodName: "unknown",
		},
	}

	for _, tc := range tt {
		ts.T().Run(tc.name, func(t *testing.T) {
			t.Parallel()

			methodName := getMethodName(tc.fullMethod)

			assert.Equal(t, tc.expectedMethodName, methodName)
		})
	}
}

func (ts *ContextTestSuite) TestExtractTraceID() {
	tt := []struct {
		name           string
		setupCtx       func() context.Context
		expectTraceID  bool
		expectedResult func(traceID trace.TraceID) string
	}{
		{
			name: "GIVEN context with valid sampled trace THEN traceID extracted",
			setupCtx: func() context.Context {
				traceID := trace.TraceID{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16}
				spanID := trace.SpanID{1, 2, 3, 4, 5, 6, 7, 8}
				spanContext := trace.NewSpanContext(trace.SpanContextConfig{
					TraceID:    traceID,
					SpanID:     spanID,
					TraceFlags: trace.FlagsSampled,
				})
				return trace.ContextWithSpanContext(context.Background(), spanContext)
			},
			expectTraceID: true,
			expectedResult: func(traceID trace.TraceID) string {
				return traceID.String()
			},
		},
		{
			name: "GIVEN context with non-sampled trace THEN traceID not extracted",
			setupCtx: func() context.Context {
				traceID := trace.TraceID{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16}
				spanID := trace.SpanID{1, 2, 3, 4, 5, 6, 7, 8}
				spanContext := trace.NewSpanContext(trace.SpanContextConfig{
					TraceID:    traceID,
					SpanID:     spanID,
					TraceFlags: 0, // Not sampled
				})
				return trace.ContextWithSpanContext(context.Background(), spanContext)
			},
			expectTraceID: false,
		},
		{
			name: "GIVEN context without trace THEN no panic and empty traceID",
			setupCtx: func() context.Context {
				return context.Background()
			},
			expectTraceID: false,
		},
	}

	for _, tc := range tt {
		ts.T().Run(tc.name, func(t *testing.T) {
			t.Parallel()

			ctx := tc.setupCtx()
			info := &telemetry.RequestInfo{}

			extractTraceID(ctx, info)

			if tc.expectTraceID {
				assert.NotEmpty(t, info.TraceID)
				spanCtx := trace.SpanContextFromContext(ctx)
				assert.Equal(t, spanCtx.TraceID().String(), info.TraceID)
			} else {
				assert.Empty(t, info.TraceID)
			}
		})
	}
}

func (ts *ContextTestSuite) TestSetRequestInfo() {
	tt := []struct {
		name         string
		setupCtx     func() context.Context
		request      any
		expectedInfo telemetry.RequestInfo
	}{
		{
			name: "GIVEN complete request with all data THEN full request info created",
			setupCtx: func() context.Context {
				md := metadata.New(map[string]string{
					HeaderSiteID:   "site123",
					HeaderOID:      "oid456",
					HeaderBidObjID: "bidobj789",
					HeaderReqID:    "req-id-001",
				})
				return metadata.NewIncomingContext(context.Background(), md)
			},
			request: &schema.GetRecommendationsRequest{
				VendorKey: "test-vendor",
				Subid:     "test-subid",
			},
			expectedInfo: telemetry.RequestInfo{
				MethodName: "GetRecommendations",
				SiteID:     "site123",
				OID:        "oid456",
				BidObjID:   "bidobj789",
				ReqID:      "req-id-001",
				VendorKey:  "test-vendor",
				SubID:      "test-subid",
			},
		},
		{
			name: "GIVEN request without metadata THEN defaults applied",
			setupCtx: func() context.Context {
				return context.Background()
			},
			request: &schema.GetRecommendationsRequest{
				VendorKey: "test-vendor",
				Subid:     "test-subid",
			},
			expectedInfo: telemetry.RequestInfo{
				MethodName: "GetRecommendations",
				SiteID:     "unknown",
				OID:        "unknown",
				VendorKey:  "test-vendor",
				SubID:      "test-subid",
			},
		},
	}

	for _, tc := range tt {
		ts.T().Run(tc.name, func(t *testing.T) {
			t.Parallel()

			ctx := tc.setupCtx()

			newCtx := setRequestInfo(ctx, "/vendorapi.VendorAPI/GetRecommendations", tc.request)

			requestInfo := telemetry.RequestInfoFromContext(newCtx)
			assert.Equal(t, tc.expectedInfo.MethodName, requestInfo.MethodName)
			assert.Equal(t, tc.expectedInfo.SiteID, requestInfo.SiteID)
			assert.Equal(t, tc.expectedInfo.OID, requestInfo.OID)
			assert.Equal(t, tc.expectedInfo.BidObjID, requestInfo.BidObjID)
			assert.Equal(t, tc.expectedInfo.ReqID, requestInfo.ReqID)
			assert.Equal(t, tc.expectedInfo.VendorKey, requestInfo.VendorKey)
			assert.Equal(t, tc.expectedInfo.SubID, requestInfo.SubID)
		})
	}
}

func TestContextTestSuite(t *testing.T) {
	t.Parallel()
	suite.Run(t, &ContextTestSuite{})
}
