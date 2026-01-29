package controller

import (
	"context"
	"errors"
	"testing"

	"rec-vendor-api/internal/config"
	controller_errors "rec-vendor-api/internal/controller/errors"
	"rec-vendor-api/internal/vendor"

	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"go.uber.org/mock/gomock"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/emptypb"

	schema "github.com/plaxieappier/rec-schema/go/vendorapi"
)

type HandlerTestSuite struct {
	suite.Suite
	mockClient     *vendor.MockClient
	vendorRegistry map[string]vendor.Client
	vendorConfig   config.VendorConfig
}

func (ts *HandlerTestSuite) SetupTest() {
	ts.mockClient = vendor.NewMockClient(gomock.NewController(ts.T()))
	ts.vendorRegistry = map[string]vendor.Client{"test_vendor": ts.mockClient}
	ts.vendorConfig = config.VendorConfig{
		Vendors: []config.Vendor{
			{
				Name: "test_vendor",
				Request: config.URLPattern{
					URL: "https://example.com/api",
				},
			},
			{
				Name: "another_vendor",
				Request: config.URLPattern{
					URL: "https://another.com/api",
				},
			},
		},
	}
}

func (ts *HandlerTestSuite) TestGetRecommendations() {
	tt := []struct {
		name         string
		vendorKey    string
		setupMock    func(mockClient *vendor.MockClient)
		wantCode     codes.Code
		wantErrMsg   string
		wantProducts []*schema.ProductInfo
	}{
		{
			name:      "GIVEN a valid request THEN expect a successful response",
			vendorKey: "test_vendor",
			setupMock: func(mc *vendor.MockClient) {
				mockResp := []vendor.ProductInfo{
					{ProductID: "1", Url: "url", Image: "img", Price: "100", SalePrice: "80", Currency: "USD"},
				}
				mc.EXPECT().GetUserRecommendationItems(gomock.Any(), gomock.Any()).Return(mockResp, nil)
			},
			wantCode: codes.OK,
			wantProducts: []*schema.ProductInfo{
				{ProductId: "1", Url: "url", Image: "img", Price: "100", SalePrice: "80", Currency: "USD"},
			},
		},
		{
			name:       "GIVEN an invalid vendor key THEN expect a bad request response",
			vendorKey:  "wrong_vendor_key",
			setupMock:  func(mc *vendor.MockClient) {},
			wantCode:   codes.InvalidArgument,
			wantErrMsg: "Vendor key 'wrong_vendor_key' not supported",
		},
		{
			name:      "GIVEN a BadRequestError error THEN expect an invalid argument error response",
			vendorKey: "test_vendor",
			setupMock: func(mc *vendor.MockClient) {
				mc.EXPECT().GetUserRecommendationItems(gomock.Any(), gomock.Any()).Return(nil, controller_errors.BadRequestErrorf("param missing"))
			},
			wantCode:   codes.InvalidArgument,
			wantErrMsg: "VendorClient returned BadRequestError. err: param missing",
		},
		{
			name:      "GIVEN an internal error THEN expect an internal server error response",
			vendorKey: "test_vendor",
			setupMock: func(mc *vendor.MockClient) {
				mc.EXPECT().GetUserRecommendationItems(gomock.Any(), gomock.Any()).Return(nil, errors.New("fail"))
			},
			wantCode:   codes.Internal,
			wantErrMsg: "Fail to recommend any products. err: fail",
		},
	}

	for _, tc := range tt {
		ts.T().Run(tc.name, func(t *testing.T) {
			tc.setupMock(ts.mockClient)

			handler, err := NewHandler(ts.vendorRegistry, ts.vendorConfig)
			require.NoError(t, err)
			request := &schema.GetRecommendationsRequest{
				VendorKey: tc.vendorKey,
				UserId:    "123",
				ClickId:   "456",
				W:         100,
				H:         200,
			}
			resp, err := handler.GetRecommendations(context.Background(), request)

			if tc.wantCode == codes.OK {
				require.NoError(t, err)
				require.NotNil(t, resp)
				require.Equal(t, len(tc.wantProducts), len(resp.Products))
				for i, wantProduct := range tc.wantProducts {
					require.True(t, proto.Equal(wantProduct, resp.Products[i]))
				}
			} else {
				require.Error(t, err)
				st, ok := status.FromError(err)
				require.True(t, ok)
				require.Equal(t, tc.wantCode, st.Code())
				require.Contains(t, st.Message(), tc.wantErrMsg)
			}
		})
	}
}

func (ts *HandlerTestSuite) TestGetVendors() {
	tt := []struct {
		name         string
		vendorConfig config.VendorConfig
		wantVendors  []*schema.VendorInfo
		wantErrMsg   string
	}{
		{
			name:         "GIVEN valid vendor config THEN expect vendors to be returned",
			vendorConfig: ts.vendorConfig,
			wantVendors: []*schema.VendorInfo{
				{VendorKey: "test_vendor", RequestHost: "example.com"},
				{VendorKey: "another_vendor", RequestHost: "another.com"},
			},
		},
		{
			name: "GIVEN empty vendor config THEN expect empty vendors list",
			vendorConfig: config.VendorConfig{
				Vendors: []config.Vendor{},
			},
			wantVendors: []*schema.VendorInfo{},
		},
		{
			name: "GIVEN vendor with invalid URL THEN expect vendor with empty request host",
			vendorConfig: config.VendorConfig{
				Vendors: []config.Vendor{
					{
						Name: "invalid_vendor",
						Request: config.URLPattern{
							URL: "://invalid-url",
						},
					},
				},
			},
			wantVendors: []*schema.VendorInfo{
				{VendorKey: "invalid_vendor", RequestHost: ""},
			},
		},
	}

	for _, tc := range tt {
		ts.T().Run(tc.name, func(t *testing.T) {
			handler, err := NewHandler(ts.vendorRegistry, tc.vendorConfig)

			require.NoError(t, err)
			require.NotNil(t, handler)
			resp, err := handler.GetVendors(context.Background(), &emptypb.Empty{})
			require.NoError(t, err)
			require.NotNil(t, resp)
			require.Equal(t, len(tc.wantVendors), len(resp.Vendors))
			for i, wantVendor := range tc.wantVendors {
				require.True(t, proto.Equal(wantVendor, resp.Vendors[i]))
			}
		})
	}
}

func TestHandlerTestSuite(t *testing.T) {
	t.Parallel()
	suite.Run(t, &HandlerTestSuite{})
}
