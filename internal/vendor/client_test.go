package vendor

import (
	"context"
	"errors"
	"fmt"
	"rec-vendor-api/internal/config"
	"rec-vendor-api/internal/telemetry"
	"testing"
	"time"

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
	mockUnmarshaler *unmarshaler.MockStrategy
	mockTracker     *url.MockStrategy
}

func (ts *VendorClientTestSuite) SetupTest() {
	ts.mockRestClient = httpkit.NewMockClient(gomock.NewController(ts.T()))
	ctrl := gomock.NewController(ts.T())
	ts.mockHeader = header.NewMockStrategy(ctrl)
	ts.mockRequester = url.NewMockStrategy(ctrl)
	ts.mockUnmarshaler = unmarshaler.NewMockStrategy(ctrl)
	ts.mockTracker = url.NewMockStrategy(ctrl)
}

func (ts *VendorClientTestSuite) TestGetUserRecommendationItems() {
	req := httpkit.NewRequest("http://test-url")
	req = req.PatchHeaders(map[string]string{"Authorization": "Bearer test"})
	req = req.SetMetrics(
		telemetry.Metrics.RestApiDurationSeconds.WithLabelValues("test-vendor", "test-site", "test-oid"),
		telemetry.Metrics.RestApiErrorTotal.WithLabelValues("test-vendor", "test-site", "test-oid"),
	)

	tt := []struct {
		name         string
		mockStrategy func()
		wantErr      bool
		want         []ProductInfo
	}{
		{
			name: "GIVEN valid response THEN expect success",
			mockStrategy: func() {
				ts.mockHeader.EXPECT().GenerateHeaders(gomock.Any()).Return(map[string]string{"Authorization": "Bearer test"})
				ts.mockRequester.EXPECT().GenerateURL(gomock.Any(), gomock.Any()).Return("http://test-url", nil)
				ts.mockRestClient.EXPECT().Get(gomock.Any(), req, 1*time.Second, []int{200}).
					Return(&httpkit.Response{Body: []byte(`[{"productId":1,"productUrl":"url1","productImage":"img1"}]`)}, nil)
				ts.mockUnmarshaler.EXPECT().UnmarshalResponse(gomock.Any(), gomock.Any()).Return([]unmarshaler.PartnerResp{{ProductID: "1", ProductURL: "url1", ProductImage: "img1"}}, nil)
				ts.mockTracker.EXPECT().GenerateURL(gomock.Any(), gomock.Any()).Return("http://tracking-url", nil)

			},
			want: []ProductInfo{{ProductID: "1", Url: "http://tracking-url", Image: "img1"}},
		},
		{
			name: "GIVEN network error THEN expect error",
			mockStrategy: func() {
				ts.mockHeader.EXPECT().GenerateHeaders(gomock.Any()).Return(map[string]string{"Authorization": "Bearer test"})
				ts.mockRequester.EXPECT().GenerateURL(gomock.Any(), gomock.Any()).Return("http://test-url", nil)
				ts.mockRestClient.EXPECT().Get(gomock.Any(), req, 1*time.Second, []int{200}).
					Return(nil, errors.New("network error"))

			},
			wantErr: true,
		},
		{
			name: "GIVEN unmarshal error THEN expect error",
			mockStrategy: func() {
				ts.mockHeader.EXPECT().GenerateHeaders(gomock.Any()).Return(map[string]string{"Authorization": "Bearer test"})
				ts.mockRequester.EXPECT().GenerateURL(gomock.Any(), gomock.Any()).Return("http://test-url", nil)
				ts.mockRestClient.EXPECT().Get(gomock.Any(), req, 1*time.Second, []int{200}).
					Return(&httpkit.Response{Body: []byte("invalid json")}, nil)
				ts.mockUnmarshaler.EXPECT().UnmarshalResponse(gomock.Any(), gomock.Any()).Return(nil, fmt.Errorf("invalid format. body: %v", "invalid json"))

			},
			wantErr: true,
		},
		{
			name: "GIVEN request URL generation error THEN expect error",
			mockStrategy: func() {
				ts.mockRequester.EXPECT().GenerateURL(gomock.Any(), gomock.Any()).Return("", fmt.Errorf("failed to generate request URL"))
			},
			wantErr: true,
		},
	}
	for _, tc := range tt {
		ts.T().Run(tc.name, func(t *testing.T) {
			vc := NewClient(
				config.Vendor{Name: "test-vendor"},
				ts.mockRestClient,
				1*time.Second,
				ts.mockHeader,
				ts.mockRequester,
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
