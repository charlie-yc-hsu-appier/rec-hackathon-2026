package vendor

import (
	"context"
	"errors"
	"fmt"
	"rec-vendor-api/internal/config"
	"rec-vendor-api/internal/telemetry"
	"testing"
	"time"

	"rec-vendor-api/internal/strategy/body"
	"rec-vendor-api/internal/strategy/header"
	"rec-vendor-api/internal/strategy/unmarshaler"
	"rec-vendor-api/internal/strategy/url"

	"github.com/plaxieappier/rec-go-kit/httpkit"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"go.uber.org/mock/gomock"
)

type VendorClientTestSuite struct {
	suite.Suite
	mockRestClient  *httpkit.MockClient
	mockHeader      *header.MockStrategy
	mockRequester   *url.MockStrategy
	mockBody        *body.MockStrategy
	mockUnmarshaler *unmarshaler.MockStrategy
	mockTracker     *url.MockStrategy
}

func (ts *VendorClientTestSuite) SetupTest() {
	ctrl := gomock.NewController(ts.T())
	ts.mockRestClient = httpkit.NewMockClient(ctrl)
	ts.mockHeader = header.NewMockStrategy(ctrl)
	ts.mockRequester = url.NewMockStrategy(ctrl)
	ts.mockBody = body.NewMockStrategy(ctrl)
	ts.mockUnmarshaler = unmarshaler.NewMockStrategy(ctrl)
	ts.mockTracker = url.NewMockStrategy(ctrl)
}

func (ts *VendorClientTestSuite) TestGetUserRecommendationItems() {
	generatedURL := "http://test-url"
	generatedHeaders := map[string]string{"Authorization": "Bearer test"}
	generatedBody := map[string]interface{}{"userId": "u1"}

	tt := []struct {
		name         string
		httpMethod   string
		mockStrategy func()
		wantErr      bool
		want         []ProductInfo
	}{
		{
			name:       "GIVEN valid GET response THEN expect success",
			httpMethod: "GET",
			mockStrategy: func() {
				ts.mockRequester.EXPECT().GenerateURL(gomock.Any(), gomock.Any()).Return(generatedURL, nil)
				ts.mockHeader.EXPECT().GenerateHeaders(gomock.Any()).Return(generatedHeaders)

				req := httpkit.NewRequest(generatedURL)
				req = req.PatchHeaders(generatedHeaders)
				req = req.SetMetrics(
					telemetry.Metrics.RestApiDurationSeconds.WithLabelValues("test-vendor", "test-site", "test-oid"),
					telemetry.Metrics.RestApiErrorTotal.WithLabelValues("test-vendor", "test-site", "test-oid"),
				)
				ts.mockRestClient.EXPECT().Get(gomock.Any(), req, 1*time.Second, []int{200}).
					Return(&httpkit.Response{Body: []byte(`[{"productId":1,"productUrl":"url1","productImage":"img1"}]`)}, nil)
				ts.mockUnmarshaler.EXPECT().UnmarshalResponse(gomock.Any(), gomock.Any()).Return([]unmarshaler.PartnerResp{{ProductID: "1", ProductURL: "url1", ProductImage: "img1"}}, nil)
				ts.mockTracker.EXPECT().GenerateURL(gomock.Any(), gomock.Any()).Return("http://tracking-url", nil)
			},
			want: []ProductInfo{{ProductID: "1", Url: "http://tracking-url", Image: "img1"}},
		},
		{
			name:       "GIVEN valid POST response THEN expect success",
			httpMethod: "POST",
			mockStrategy: func() {
				ts.mockRequester.EXPECT().GenerateURL(gomock.Any(), gomock.Any()).Return(generatedURL, nil)
				ts.mockHeader.EXPECT().GenerateHeaders(gomock.Any()).Return(generatedHeaders)
				ts.mockBody.EXPECT().GenerateBody(gomock.Any()).Return(generatedBody)

				req := httpkit.NewRequest(generatedURL)
				req = req.PatchHeaders(generatedHeaders)
				req = req.SetBody(generatedBody)
				req = req.SetMetrics(
					telemetry.Metrics.RestApiDurationSeconds.WithLabelValues("test-vendor", "test-site", "test-oid"),
					telemetry.Metrics.RestApiErrorTotal.WithLabelValues("test-vendor", "test-site", "test-oid"),
				)
				ts.mockRestClient.EXPECT().Post(gomock.Any(), req, 1*time.Second, []int{200}).
					Return(&httpkit.Response{Body: []byte(`[{"productId":2,"productUrl":"url2","productImage":"img2"}]`)}, nil)
				ts.mockUnmarshaler.EXPECT().UnmarshalResponse(gomock.Any(), gomock.Any()).Return([]unmarshaler.PartnerResp{{ProductID: "2", ProductURL: "url2", ProductImage: "img2"}}, nil)
				ts.mockTracker.EXPECT().GenerateURL(gomock.Any(), gomock.Any()).Return("http://tracking-url-post", nil)
			},
			want: []ProductInfo{{ProductID: "2", Url: "http://tracking-url-post", Image: "img2"}},
		},
		{
			name:       "GIVEN network error THEN expect error",
			httpMethod: "GET",
			mockStrategy: func() {
				ts.mockRequester.EXPECT().GenerateURL(gomock.Any(), gomock.Any()).Return(generatedURL, nil)
				ts.mockHeader.EXPECT().GenerateHeaders(gomock.Any()).Return(generatedHeaders)

				req := httpkit.NewRequest(generatedURL)
				req = req.PatchHeaders(generatedHeaders)
				req = req.SetMetrics(
					telemetry.Metrics.RestApiDurationSeconds.WithLabelValues("test-vendor", "test-site", "test-oid"),
					telemetry.Metrics.RestApiErrorTotal.WithLabelValues("test-vendor", "test-site", "test-oid"),
				)
				ts.mockRestClient.EXPECT().Get(gomock.Any(), req, 1*time.Second, []int{200}).
					Return(nil, errors.New("network error"))
			},
			wantErr: true,
		},
		{
			name:       "GIVEN unmarshal error THEN expect error",
			httpMethod: "GET",
			mockStrategy: func() {
				ts.mockRequester.EXPECT().GenerateURL(gomock.Any(), gomock.Any()).Return(generatedURL, nil)
				ts.mockHeader.EXPECT().GenerateHeaders(gomock.Any()).Return(generatedHeaders)
				ts.mockRestClient.EXPECT().Get(gomock.Any(), gomock.Any(), 1*time.Second, []int{200}).
					Return(&httpkit.Response{Body: []byte("invalid json")}, nil)
				ts.mockUnmarshaler.EXPECT().UnmarshalResponse(gomock.Any(), gomock.Any()).Return(nil, fmt.Errorf("invalid format. body: %v", "invalid json"))
			},
			wantErr: true,
		},
		{
			name:       "GIVEN request URL generation error THEN expect error",
			httpMethod: "GET",
			mockStrategy: func() {
				ts.mockRequester.EXPECT().GenerateURL(gomock.Any(), gomock.Any()).Return("", fmt.Errorf("failed to generate request URL"))
			},
			wantErr: true,
		},
	}
	for _, tc := range tt {
		ts.T().Run(tc.name, func(t *testing.T) {
			vc := NewClient(
				config.Vendor{Name: "test-vendor", HTTPMethod: tc.httpMethod},
				ts.mockRestClient,
				1*time.Second,
				ts.mockHeader,
				ts.mockRequester,
				ts.mockBody,
				ts.mockUnmarshaler,
				ts.mockTracker,
			)

			tc.mockStrategy()
			ctx := telemetry.RequestInfoToContext(context.Background(), telemetry.RequestInfo{
				SiteID: "test-site",
				OID:    "test-oid",
			})
			got, err := vc.GetUserRecommendationItems(ctx, Request{UserID: "u1"})
			require.Equal(t, tc.want, got)
			if tc.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestVendorClientTestSuite(t *testing.T) {
	t.Parallel()
	suite.Run(t, &VendorClientTestSuite{})
}

func TestCategorizeError(t *testing.T) {
	t.Parallel()

	tt := []struct {
		name     string
		restResp *httpkit.Response
		err      error
		want     string
	}{
		{
			name:     "GIVEN nil error THEN return empty string",
			restResp: nil,
			err:      nil,
			want:     "",
		},
		{
			name:     "GIVEN context deadline exceeded THEN return network timeout",
			restResp: nil,
			err:      context.DeadlineExceeded,
			want:     errNetworkTimeout,
		},
		{
			name:     "GIVEN context canceled THEN return network timeout",
			restResp: nil,
			err:      context.Canceled,
			want:     errNetworkTimeout,
		},
		{
			name:     "GIVEN io.EOF error THEN return remote connection reset",
			restResp: nil,
			err:      errors.New("Get \"xxx\": EOF"),
			want:     errRemoteConnectionReset,
		},
		{
			name:     "GIVEN connection reset error THEN return remote connection reset",
			restResp: nil,
			err:      errors.New("invalid status code, status: 502 Bad Gateway, resp body: upstream error: Post \"xxx\": read: connection reset by peer"),
			want:     errRemoteConnectionReset,
		},
		{
			name:     "GIVEN http response with 404 status THEN return invalid http status code",
			restResp: &httpkit.Response{StatusCode: 404},
			err:      errors.New("http error"),
			want:     errInvalidHTTPStatus + "404",
		},
		{
			name:     "GIVEN http response with 500 status THEN return invalid http status code",
			restResp: &httpkit.Response{StatusCode: 500},
			err:      errors.New("http error"),
			want:     errInvalidHTTPStatus + "500",
		},
		{
			name:     "GIVEN generic network error THEN return unknown network error",
			restResp: nil,
			err:      errors.New("some generic network error"),
			want:     errUnknownNetworkError,
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			got := categorizeError(tc.restResp, tc.err)
			require.Equal(t, tc.want, got)
		})
	}
}
