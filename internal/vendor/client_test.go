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
	"rec-vendor-api/internal/strategy/requester"
	"rec-vendor-api/internal/strategy/tracker"
	"rec-vendor-api/internal/strategy/unmarshaler"

	"github.com/plaxieappier/rec-go-kit/httpkit"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"go.uber.org/mock/gomock"
)

type VendorClientTestSuite struct {
	suite.Suite
	mockRestClient  *httpkit.MockClient
	mockHeader      *header.MockStrategy
	mockRequester   *requester.MockStrategy
	mockUnmarshaler *unmarshaler.MockStrategy
	mockTracker     *tracker.MockStrategy
}

func (ts *VendorClientTestSuite) SetupTest() {
	ts.mockRestClient = httpkit.NewMockClient(gomock.NewController(ts.T()))
	ctrl := gomock.NewController(ts.T())
	ts.mockHeader = header.NewMockStrategy(ctrl)
	ts.mockRequester = requester.NewMockStrategy(ctrl)
	ts.mockUnmarshaler = unmarshaler.NewMockStrategy(ctrl)
	ts.mockTracker = tracker.NewMockStrategy(ctrl)
}

func (ts *VendorClientTestSuite) TestGetUserRecommendationItems() {
	tt := []struct {
		name         string
		mockResp     *httpkit.Response
		mockRespErr  error
		mockStrategy func()
		wantErr      bool
		want         Response
	}{
		{
			name:     "GIVEN valid response THEN expect success",
			mockResp: &httpkit.Response{Body: []byte(`[{"productId":1,"productUrl":"url1","productImage":"img1"}]`)},
			mockStrategy: func() {
				ts.mockHeader.EXPECT().GenerateHeaders(gomock.Any()).Return(map[string]string{"Authorization": "Bearer test"})
				ts.mockRequester.EXPECT().GenerateRequestURL(gomock.Any()).Return("http://test-url")
				ts.mockUnmarshaler.EXPECT().UnmarshalResponse(gomock.Any()).Return(&[]unmarshaler.CoupangPartnerResp{{ProductID: 1, ProductURL: "url1", ProductImage: "img1"}}, nil)
				ts.mockTracker.EXPECT().GenerateTrackingURL(gomock.Any()).Return("http://tracking-url")

			},
			want: Response{ProductIDs: []string{"1"}, ProductPatch: map[string]ProductPatch{"1": {Url: "http://tracking-url", Image: "img1"}}},
		},
		{
			name:        "GIVEN network error THEN expect error",
			mockRespErr: errors.New("network error"),
			mockStrategy: func() {
				ts.mockHeader.EXPECT().GenerateHeaders(gomock.Any()).Return(map[string]string{"Authorization": "Bearer test"})
				ts.mockRequester.EXPECT().GenerateRequestURL(gomock.Any()).Return("http://test-url")

			},
			wantErr: true,
			want:    Response{},
		},
		{
			name:     "GIVEN unmarshal error THEN expect error",
			mockResp: &httpkit.Response{Body: []byte("invalid json")},
			mockStrategy: func() {
				ts.mockHeader.EXPECT().GenerateHeaders(gomock.Any()).Return(map[string]string{"Authorization": "Bearer test"})
				ts.mockRequester.EXPECT().GenerateRequestURL(gomock.Any()).Return("http://test-url")
				ts.mockUnmarshaler.EXPECT().UnmarshalResponse(gomock.Any()).Return(nil, fmt.Errorf("invalid format. body: %v", "invalid json"))
			},
			wantErr: true,
			want:    Response{},
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
			req := httpkit.NewRequest("http://test-url")
			req = req.PatchHeaders(map[string]string{"Authorization": "Bearer test"})
			req = req.SetMetrics(
				telemetry.Metrics.RestApiDurationSeconds.WithLabelValues("test-vendor"),
				telemetry.Metrics.RestApiErrorTotal.WithLabelValues("test-vendor"),
			)
			ts.mockRestClient.EXPECT().Get(gomock.Any(), req, 1*time.Second, []int{200}).Return(tc.mockResp, tc.mockRespErr)

			got, err := vc.GetUserRecommendationItems(context.Background(), Request{UserID: "u1"})
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
