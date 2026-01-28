package controller

import (
	"errors"
	"net/http"
	"net/http/httptest"
	controller_errors "rec-vendor-api/internal/controller/errors"
	"testing"

	"rec-vendor-api/internal/vendor"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"go.uber.org/mock/gomock"
)

type RecommenderTestSuite struct {
	suite.Suite
	mockClient     *vendor.MockClient
	vendorRegistry map[string]vendor.Client
}

func (ts *RecommenderTestSuite) SetupTest() {
	ts.mockClient = vendor.NewMockClient(gomock.NewController(ts.T()))
	ts.vendorRegistry = map[string]vendor.Client{"test_vendor": ts.mockClient}
}

func (ts *RecommenderTestSuite) TestRecommend() {
	tt := []struct {
		name       string
		vendorKey  string
		requestURL string
		setupMock  func(mockClient *vendor.MockClient)
		wantCode   int
		wantBody   string
	}{
		{
			name:       "GIVEN a valid request THEN expect a successful response",
			vendorKey:  "test_vendor",
			requestURL: "/r/test_vendor?user_id=123&click_id=456&w=100&h=200",
			setupMock: func(mc *vendor.MockClient) {
				mockResp := []vendor.ProductInfo{{ProductID: "1", Url: "url", Image: "img"}}
				mc.EXPECT().GetUserRecommendationItems(gomock.Any(), gomock.Any()).Return(mockResp, nil)
			},
			wantCode: http.StatusOK,
			wantBody: `[{"product_id":"1","url":"url","image":"img","price":"","sale_price":"","currency":""}]`,
		},
		{
			name:       "GIVEN an invalid vendor key THEN expect a bad request response",
			vendorKey:  "bad_vendor",
			requestURL: "/r/bad_vendor?user_id=123&click_id=456&w=100&h=200",
			setupMock:  func(mc *vendor.MockClient) {},
			wantCode:   http.StatusBadRequest,
			wantBody:   `{"detail":"vendor key 'bad_vendor' not supported", "status":400}`,
		},
		{
			name:       "GIVEN a missing user ID THEN expect a bad request response",
			vendorKey:  "test_vendor",
			requestURL: "/r/test_vendor?click_id=456&w=100&h=200",
			setupMock:  func(mc *vendor.MockClient) {},
			wantCode:   http.StatusBadRequest,
			wantBody:   `{"detail":"Key: 'Request.UserID' Error:Field validation for 'UserID' failed on the 'required' tag", "status":400}`,
		},
		{
			name:       "GIVEN an BadRequestError error THEN expect an bad request error response",
			vendorKey:  "test_vendor",
			requestURL: "/r/test_vendor?user_id=123&click_id=456&w=100&h=200",
			setupMock: func(mc *vendor.MockClient) {
				mc.EXPECT().GetUserRecommendationItems(gomock.Any(), gomock.Any()).Return(nil, controller_errors.BadRequestErrorf("param missing"))
			},
			wantCode: http.StatusBadRequest,
			wantBody: `{"detail":"VendorClient returned BadRequestError. err: param missing", "status":400}`,
		},
		{
			name:       "GIVEN an internal error THEN expect an internal server error response",
			vendorKey:  "test_vendor",
			requestURL: "/r/test_vendor?user_id=123&click_id=456&w=100&h=200",
			setupMock: func(mc *vendor.MockClient) {
				mc.EXPECT().GetUserRecommendationItems(gomock.Any(), gomock.Any()).Return(nil, errors.New("fail"))
			},
			wantCode: http.StatusInternalServerError,
			wantBody: `{"detail":"fail to recommend any products for vendor test_vendor. err: fail", "status":500}`,
		},
	}

	for _, tc := range tt {
		ts.T().Run(tc.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			c.Request = httptest.NewRequest(http.MethodGet, tc.requestURL, nil)
			c.Params = []gin.Param{{Key: "vendor_key", Value: tc.vendorKey}}

			tc.setupMock(ts.mockClient)

			r := NewRecommender(ts.vendorRegistry)
			r.Recommend(c)

			require.Equal(t, tc.wantCode, w.Code)
			require.JSONEq(t, tc.wantBody, w.Body.String())
		})
	}
}

func TestRecommenderTestSuite(t *testing.T) {
	t.Parallel()
	suite.Run(t, &RecommenderTestSuite{})
}
