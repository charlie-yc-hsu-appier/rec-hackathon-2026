package vendor

import (
	"context"
	"encoding/json"
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

type mockRespUnmarshalStrategy struct{}

func (m *mockRespUnmarshalStrategy) UnmarshalResponse(body []byte) (*[]unmarshaler.CoupangPartnerResp, error) {
	var resp []unmarshaler.CoupangPartnerResp
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, fmt.Errorf("invalid format. body: %v", string(body))
	}
	return &resp, nil
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
		name        string
		mockResp    *httpkit.Response
		wantedError error
		want        Response
	}{
		{
			name:        "GIVEN valid response THEN expect success",
			mockResp:    &httpkit.Response{Body: []byte(`[{"productId":1,"productUrl":"url1","productImage":"img1"},{"productId":2,"productUrl":"url2","productImage":"img2"}]`)},
			wantedError: nil,
			want:        Response{ProductIDs: []string{"1", "2"}, ProductPatch: map[string]ProductPatch{"1": {Url: "http://tracking-url", Image: "img1"}, "2": {Url: "http://tracking-url", Image: "img2"}}},
		},
		{
			name:        "GIVEN network error THEN expect error",
			mockResp:    nil,
			wantedError: errors.New("network error"),
			want:        Response{},
		},
		{
			name:        "GIVEN unmarshal error THEN expect error",
			mockResp:    &httpkit.Response{Body: []byte("invalid json")},
			wantedError: fmt.Errorf("invalid format. body: %v", "invalid json"),
			want:        Response{},
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
				&mockRespUnmarshalStrategy{},
				&mockTrackingURLStrategy{},
			)
			req := httpkit.NewRequest("http://test-url")
			req = req.SetMetrics(
				telemetry.Metrics.RestApiDurationSeconds.WithLabelValues("test-vendor"),
				telemetry.Metrics.RestApiErrorTotal.WithLabelValues("test-vendor"),
			)
			ts.mockRestClient.EXPECT().Get(gomock.Any(), req, 1*time.Second, []int{200}).Return(tc.mockResp, tc.wantedError)
			got, err := vc.GetUserRecommendationItems(context.Background(), Request{UserID: "u1"})
			if tc.wantedError != nil {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				require.Equal(t, tc.want, got)
			}
		})
	}
}

func TestVendorClientTestSuite(t *testing.T) {
	t.Parallel()
	suite.Run(t, &VendorClientTestSuite{})
}
