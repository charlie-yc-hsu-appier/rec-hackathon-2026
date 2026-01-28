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
		request      *schema.GetRecommendationsRequest
		setupMock    func(mockClient *vendor.MockClient)
		setupCtx     func() context.Context
		wantCode     codes.Code
		wantErrMsg   string
		wantProducts []*schema.ProductInfo
	}{
		{
			name: "GIVEN a valid request THEN expect a successful response",
			request: &schema.GetRecommendationsRequest{
				VendorKey: "test_vendor",
				UserId:    "123",
				ClickId:   "456",
				W:         100,
				H:         200,
			},
			setupMock: func(mc *vendor.MockClient) {
				mockResp := []vendor.ProductInfo{
					{ProductID: "1", Url: "url", Image: "img", Price: "100", SalePrice: "80", Currency: "USD"},
				}
				mc.EXPECT().GetUserRecommendationItems(gomock.Any(), gomock.Any()).Return(mockResp, nil)
			},
			setupCtx: func() context.Context {
				return context.Background()
			},
			wantCode: codes.OK,
			wantProducts: []*schema.ProductInfo{
				{ProductId: "1", Url: "url", Image: "img", Price: "100", SalePrice: "80", Currency: "USD"},
			},
		},
		{
			name: "GIVEN an invalid vendor key THEN expect an invalid argument error",
			request: &schema.GetRecommendationsRequest{
				VendorKey: "bad_vendor",
				UserId:    "123",
				ClickId:   "456",
				W:         100,
				H:         200,
			},
			setupMock: func(mc *vendor.MockClient) {},
			setupCtx: func() context.Context {
				return context.Background()
			},
			wantCode:   codes.InvalidArgument,
			wantErrMsg: "Vendor key 'bad_vendor' not supported",
		},
		{
			name: "GIVEN a BadRequestError error THEN expect an invalid argument error response",
			request: &schema.GetRecommendationsRequest{
				VendorKey: "test_vendor",
				UserId:    "123",
				ClickId:   "456",
				W:         100,
				H:         200,
			},
			setupMock: func(mc *vendor.MockClient) {
				mc.EXPECT().GetUserRecommendationItems(gomock.Any(), gomock.Any()).Return(nil, controller_errors.BadRequestErrorf("param missing"))
			},
			setupCtx: func() context.Context {
				return context.Background()
			},
			wantCode:   codes.InvalidArgument,
			wantErrMsg: "VendorClient returned BadRequestError. err: param missing",
		},
		{
			name: "GIVEN an internal error THEN expect an internal server error response",
			request: &schema.GetRecommendationsRequest{
				VendorKey: "test_vendor",
				UserId:    "123",
				ClickId:   "456",
				W:         100,
				H:         200,
			},
			setupMock: func(mc *vendor.MockClient) {
				mc.EXPECT().GetUserRecommendationItems(gomock.Any(), gomock.Any()).Return(nil, errors.New("fail"))
			},
			setupCtx: func() context.Context {
				return context.Background()
			},
			wantCode:   codes.Internal,
			wantErrMsg: "Fail to recommend any products. err: fail",
		},
	}

	for _, tc := range tt {
		ts.T().Run(tc.name, func(t *testing.T) {
			ctx := tc.setupCtx()
			tc.setupMock(ts.mockClient)

			handler, err := NewHandler(ts.vendorRegistry, ts.vendorConfig)
			require.NoError(t, err)
			resp, err := handler.GetRecommendations(ctx, tc.request)

			if tc.wantCode == codes.OK {
				require.NoError(t, err)
				require.NotNil(t, resp)
				require.Equal(t, len(tc.wantProducts), len(resp.Products))
				for i, wantProduct := range tc.wantProducts {
					require.Equal(t, wantProduct.ProductId, resp.Products[i].ProductId)
					require.Equal(t, wantProduct.Url, resp.Products[i].Url)
					require.Equal(t, wantProduct.Image, resp.Products[i].Image)
					require.Equal(t, wantProduct.Price, resp.Products[i].Price)
					require.Equal(t, wantProduct.SalePrice, resp.Products[i].SalePrice)
					require.Equal(t, wantProduct.Currency, resp.Products[i].Currency)
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
		wantErr      bool
		wantErrMsg   string
	}{
		{
			name:         "GIVEN valid vendor config THEN expect vendors to be returned",
			vendorConfig: ts.vendorConfig,
			wantVendors: []*schema.VendorInfo{
				{VendorKey: "test_vendor", RequestHost: "example.com"},
				{VendorKey: "another_vendor", RequestHost: "another.com"},
			},
			wantErr: false,
		},
		{
			name: "GIVEN vendor config with invalid URL THEN expect an error",
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
			wantErr:    true,
			wantErrMsg: "failed to parse request URL for vendor invalid_vendor",
		},
		{
			name: "GIVEN empty vendor config THEN expect empty vendors list",
			vendorConfig: config.VendorConfig{
				Vendors: []config.Vendor{},
			},
			wantVendors: []*schema.VendorInfo{},
			wantErr:     false,
		},
	}

	for _, tc := range tt {
		ts.T().Run(tc.name, func(t *testing.T) {
			handler, err := NewHandler(ts.vendorRegistry, tc.vendorConfig)

			if tc.wantErr {
				require.Error(t, err)
				require.Nil(t, handler)
				require.Contains(t, err.Error(), tc.wantErrMsg)
			} else {
				require.NoError(t, err)
				require.NotNil(t, handler)
				resp, err := handler.GetVendors(context.Background(), &emptypb.Empty{})
				require.NoError(t, err)
				require.NotNil(t, resp)
				require.Equal(t, len(tc.wantVendors), len(resp.Vendors))
				for i, wantVendor := range tc.wantVendors {
					require.Equal(t, wantVendor.VendorKey, resp.Vendors[i].VendorKey)
					require.Equal(t, wantVendor.RequestHost, resp.Vendors[i].RequestHost)
				}
			}
		})
	}
}

func TestHandlerTestSuite(t *testing.T) {
	t.Parallel()
	suite.Run(t, &HandlerTestSuite{})
}
