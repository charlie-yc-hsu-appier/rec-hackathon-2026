package controller

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"rec-vendor-api/internal/vendor"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"go.uber.org/mock/gomock"
)

type RecommenderTestSuite struct {
	suite.Suite
	mockClient *vendor.MockClient
}

func (ts *RecommenderTestSuite) SetupTest() {
	ts.mockClient = vendor.NewMockClient(gomock.NewController(ts.T()))
}

func (ts *RecommenderTestSuite) TestRecommend() {
	tt := []struct {
		name           string
		requestURL     string
		vendorRegistry map[string]vendor.Client
		setupMock      func(mockClient *vendor.MockClient)
		wantCode       int
		wantBody       string
	}{
		{
			name:           "GIVEN a valid request THEN expect a successful response",
			requestURL:     "/r?vendor_key=test_vendor&user_id=123&w=100&h=200",
			vendorRegistry: map[string]vendor.Client{"test_vendor": ts.mockClient},
			setupMock: func(mc *vendor.MockClient) {
				mockResp := vendor.Response{ProductIDs: []string{"1"}, ProductPatch: map[string]vendor.ProductPatch{"1": {Url: "url", Image: "img"}}}
				mc.EXPECT().GetUserRecommendationItems(gomock.Any(), gomock.Any()).Return(mockResp, nil)
			},
			wantCode: http.StatusOK,
			wantBody: `{"product_ids":["1"],"product_patch":{"1":{"url":"url","image":"img"}}}`,
		},
		{
			name:           "GIVEN an invalid vendor key THEN expect a bad request response",
			requestURL:     "/r?vendor_key=bad_vendor&user_id=123&w=100&h=200",
			vendorRegistry: map[string]vendor.Client{},
			setupMock:      func(mc *vendor.MockClient) {},
			wantCode:       http.StatusBadRequest,
			wantBody:       `{"detail":"Vendor key 'bad_vendor' not supported", "status":400}`,
		},
		{
			name:           "GIVEN a missing user ID THEN expect a bad request response",
			requestURL:     "/r?vendor_key=test_vendor&w=100&h=200",
			vendorRegistry: map[string]vendor.Client{"test_vendor": ts.mockClient},
			setupMock:      func(mc *vendor.MockClient) {},
			wantCode:       http.StatusBadRequest,
			wantBody:       `{"detail":"Key: 'RecommendQuery.UserID' Error:Field validation for 'UserID' failed on the 'required' tag", "status":400}`,
		},
		{
			name:           "GIVEN an internal error THEN expect an internal server error response",
			requestURL:     "/r?vendor_key=test_vendor&user_id=123&w=100&h=200",
			vendorRegistry: map[string]vendor.Client{"test_vendor": ts.mockClient},
			setupMock: func(mc *vendor.MockClient) {
				mc.EXPECT().GetUserRecommendationItems(gomock.Any(), gomock.Any()).Return(vendor.Response{}, errors.New("fail"))
			},
			wantCode: http.StatusInternalServerError,
			wantBody: `{"detail":"fail", "status":500}`,
		},
	}

	for _, tc := range tt {
		ts.T().Run(tc.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			c.Request = httptest.NewRequest(http.MethodGet, tc.requestURL, nil)

			tc.setupMock(ts.mockClient)

			vc := NewVendorController(tc.vendorRegistry)
			vc.Recommend(c)

			require.Equal(t, tc.wantCode, w.Code)
			require.JSONEq(t, tc.wantBody, w.Body.String())
		})
	}
}

func TestRecommenderTestSuite(t *testing.T) {
	t.Parallel()
	suite.Run(t, &RecommenderTestSuite{})
}
