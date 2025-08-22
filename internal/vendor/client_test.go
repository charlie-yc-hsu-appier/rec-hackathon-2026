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

	"bitbucket.org/plaxieappier/rec-go-kit/httpkit"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"go.uber.org/mock/gomock"
)

type mockHeaderStrategy struct{}

func (m *mockHeaderStrategy) GenerateHeaders(params header.Params) map[string]string {
	return nil
}

type mockRequestURLStrategy struct{}

func (m *mockRequestURLStrategy) GenerateRequestURL(params requester.Params) string {
	return "http://test-url"
}

type mockRespUnmarshalStrategy struct {
	resp *[]unmarshaler.CoupangPartnerResp
	err  error
}

func (m *mockRespUnmarshalStrategy) UnmarshalResponse(body []byte) (*[]unmarshaler.CoupangPartnerResp, error) {
	return m.resp, m.err
}

type mockTrackingURLStrategy struct{}

func (m *mockTrackingURLStrategy) GenerateTrackingURL(params tracker.Params) string {
	return "http://tracking-url"
}

type VendorClientTestSuite struct {
	suite.Suite
	mockRestClient *httpkit.MockClient
}

func (ts *VendorClientTestSuite) SetupTest() {
	ts.mockRestClient = httpkit.NewMockClient(gomock.NewController(ts.T()))
}

func (ts *VendorClientTestSuite) TestGetUserRecommendationItems() {
	tt := []struct {
		name          string
		mockResp      *httpkit.Response
		mockRespErr   error
		mockUnmarshal *mockRespUnmarshalStrategy
		wantErr       bool
		want          Response
	}{
		{
			name:          "GIVEN valid response THEN expect success",
			mockResp:      &httpkit.Response{Body: []byte(`[{"productId":1,"productUrl":"url1","productImage":"img1"}]`)},
			mockUnmarshal: &mockRespUnmarshalStrategy{resp: &[]unmarshaler.CoupangPartnerResp{{ProductID: 1, ProductURL: "url1", ProductImage: "img1"}}},
			want:          Response{ProductIDs: []string{"1"}, ProductPatch: map[string]ProductPatch{"1": {Url: "http://tracking-url", Image: "img1"}}},
		},
		{
			name:        "GIVEN network error THEN expect error",
			mockRespErr: errors.New("network error"),
			wantErr:     true,
			want:        Response{},
		},
		{
			name:          "GIVEN unmarshal error THEN expect error",
			mockResp:      &httpkit.Response{Body: []byte("invalid json")},
			mockUnmarshal: &mockRespUnmarshalStrategy{err: fmt.Errorf("invalid format. body: %v", "invalid json")},
			wantErr:       true,
			want:          Response{},
		},
	}
	for _, tc := range tt {
		ts.T().Run(tc.name, func(t *testing.T) {
			vc := NewClient(
				config.Vendor{Name: "test-vendor"},
				ts.mockRestClient,
				1*time.Second,
				&mockHeaderStrategy{},
				&mockRequestURLStrategy{},
				tc.mockUnmarshal,
				&mockTrackingURLStrategy{},
			)
			req := httpkit.NewRequest("http://test-url")
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
