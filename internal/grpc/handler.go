package grpc

import (
	"context"
	"errors"
	"net/url"

	"rec-vendor-api/internal/config"
	"rec-vendor-api/internal/vendor"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"

	controller_errors "rec-vendor-api/internal/controller/errors"

	grpc_realip "github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/realip"
	schema "github.com/plaxieappier/rec-schema/go/vendorapi"
	log "github.com/sirupsen/logrus"
)

type Handler interface {
	GetRecommendations(context.Context, *schema.GetRecommendationsRequest) (*schema.GetRecommendationsResponse, error)
	GetVendors(context.Context, *emptypb.Empty) (*schema.GetVendorsResponse, error)
	CheckHealthCheck(context.Context, *emptypb.Empty) (*schema.HealthcheckResponse, error)
}

type HandlerImpl struct {
	schema.UnimplementedVendorAPIServer
	vendorRegistry map[string]vendor.Client
	vendorConfig   config.VendorConfig
}

func NewHandler(vendorRegistry map[string]vendor.Client, vendorConfig config.VendorConfig) *HandlerImpl {
	return &HandlerImpl{
		vendorRegistry: vendorRegistry,
		vendorConfig:   vendorConfig,
	}
}

func (s *HandlerImpl) GetRecommendations(ctx context.Context, req *schema.GetRecommendationsRequest) (*schema.GetRecommendationsResponse, error) {
	internalReq := convertToInternalRequest(ctx, req)

	vendorKey := req.VendorKey
	vendorClient := s.vendorRegistry[vendorKey]
	if vendorClient == nil {
		log.WithContext(ctx).Errorf("Invalid vendor key: %s", vendorKey)
		return nil, status.Errorf(codes.InvalidArgument, "Vendor key '%s' not supported", vendorKey)
	}
	products, err := vendorClient.GetUserRecommendationItems(ctx, internalReq)
	if err != nil {
		var badRequestErr *controller_errors.BadRequestError
		if errors.As(err, &badRequestErr) {
			log.WithContext(ctx).Errorf("VendorClient returned BadRequestError. err: %v", err)
			return nil, status.Errorf(codes.InvalidArgument, "VendorClient returned BadRequestError. err: %v", err)
		}
		log.WithContext(ctx).Errorf("Fail to recommend any products. err: %v", err)
		return nil, status.Errorf(codes.Internal, "Fail to recommend any products. err: %v", err)
	}
	protoProducts := make([]*schema.ProductInfo, 0, len(products))
	for _, product := range products {
		protoProducts = append(protoProducts, &schema.ProductInfo{
			ProductId: product.ProductID,
			Url:       product.Url,
			Image:     product.Image,
			Price:     product.Price,
			SalePrice: product.SalePrice,
			Currency:  product.Currency,
		})
	}
	return &schema.GetRecommendationsResponse{
		Products: protoProducts,
	}, nil
}

func (s *HandlerImpl) GetVendors(ctx context.Context, _ *emptypb.Empty) (*schema.GetVendorsResponse, error) {
	vendors := make([]*schema.VendorInfo, 0, len(s.vendorConfig.Vendors))
	for _, v := range s.vendorConfig.Vendors {
		requestHost := ""
		if parsedURL, err := url.Parse(v.Request.URL); err == nil {
			requestHost = parsedURL.Host
		} else {
			return nil, status.Errorf(codes.InvalidArgument, "failed to parse request URL for vendor %s: %s", v.Name, err.Error())
		}
		vendors = append(vendors, &schema.VendorInfo{
			VendorKey:   v.Name,
			RequestHost: requestHost,
		})
	}
	return &schema.GetVendorsResponse{
		Vendors: vendors,
	}, nil
}

func (s *HandlerImpl) CheckHealthCheck(ctx context.Context, req *emptypb.Empty) (*schema.HealthcheckResponse, error) {
	return &schema.HealthcheckResponse{
		Status: "ok",
	}, nil
}

func convertToInternalRequest(ctx context.Context, req *schema.GetRecommendationsRequest) vendor.Request {
	clientIP, exists := grpc_realip.FromContext(ctx)
	clientIPStr := ""
	if exists {
		clientIPStr = clientIP.String()
	}

	osStr := ""
	switch req.Os {
	case schema.OperationSystem_ANDROID:
		osStr = "android"
	case schema.OperationSystem_IOS:
		osStr = "ios"
	}

	return vendor.Request{
		UserID:          req.UserId,
		ClickID:         req.ClickId,
		ImgWidth:        int(req.W),
		ImgHeight:       int(req.H),
		WebHost:         req.WebHost,
		BundleID:        req.BundleId,
		AdType:          int(req.Adtype),
		PartnerID:       req.PartnerId,
		OS:              osStr,
		SubID:           req.Subid,
		KeetaCampaignID: req.KCampaignId,
		Latitude:        req.Lat,
		Longitude:       req.Lon,
		ClientIP:        clientIPStr,
	}
}
