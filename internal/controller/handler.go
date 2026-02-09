package controller

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
	HealthCheck(context.Context, *emptypb.Empty) (*schema.HealthcheckResponse, error)
}

type HandlerImpl struct {
	schema.UnimplementedVendorAPIServer
	vendorRegistry map[string]vendor.Client
	vendorInfo     []*schema.VendorInfo
}

func NewHandler(vendorRegistry map[string]vendor.Client, vendorConfig config.VendorConfig) (*HandlerImpl, error) {
	return &HandlerImpl{
		vendorRegistry: vendorRegistry,
		vendorInfo:     initVendorInfo(vendorConfig),
	}, nil
}

func (s *HandlerImpl) GetRecommendations(ctx context.Context, req *schema.GetRecommendationsRequest) (*schema.GetRecommendationsResponse, error) {
	vendorKey := req.VendorKey
	vendorClient := s.vendorRegistry[vendorKey]
	if vendorClient == nil {
		log.WithContext(ctx).Errorf("Invalid vendor key: %s", vendorKey)
		return nil, status.Errorf(codes.InvalidArgument, "Vendor key '%s' not supported", vendorKey)
	}

	vendorReq := toVendorRequest(ctx, req)
	products, err := vendorClient.GetUserRecommendationItems(ctx, vendorReq)
	if err != nil {
		var badRequestErr *controller_errors.BadRequestError
		if errors.As(err, &badRequestErr) {
			log.WithContext(ctx).Errorf("VendorClient returned BadRequestError. err: %v", err)
			return nil, status.Errorf(codes.InvalidArgument, "VendorClient returned BadRequestError. err: %v", err)
		}
		log.WithContext(ctx).Errorf("Fail to recommend any products. err: %v", err)
		return nil, status.Errorf(codes.Internal, "Fail to recommend any products. err: %v", err)
	}

	return toProto(products)
}

func (s *HandlerImpl) GetVendors(_ context.Context, _ *emptypb.Empty) (*schema.GetVendorsResponse, error) {
	return &schema.GetVendorsResponse{
		Vendors: s.vendorInfo,
	}, nil
}

func (s *HandlerImpl) HealthCheck(_ context.Context, _ *emptypb.Empty) (*schema.HealthcheckResponse, error) {
	return &schema.HealthcheckResponse{
		Status: "ok",
	}, nil
}

func toVendorRequest(ctx context.Context, req *schema.GetRecommendationsRequest) vendor.Request {
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
		ClientIP:        getClientIP(ctx),
	}
}

func initVendorInfo(vendorConfig config.VendorConfig) []*schema.VendorInfo {
	vendors := make([]*schema.VendorInfo, 0, len(vendorConfig.Vendors))
	for _, v := range vendorConfig.Vendors {
		requestHost := ""
		if parsedURL, err := url.Parse(v.Request.URL); err == nil {
			requestHost = parsedURL.Host
		}
		vendors = append(vendors, &schema.VendorInfo{
			VendorKey:   v.Name,
			RequestHost: requestHost,
		})
	}
	return vendors
}

func toProto(products []vendor.ProductInfo) (*schema.GetRecommendationsResponse, error) {
	protoProducts := make([]*schema.ProductInfo, len(products))
	for i, product := range products {
		protoProducts[i] = &schema.ProductInfo{
			ProductId: product.ProductID,
			Url:       product.Url,
			Image:     product.Image,
			Price:     product.Price,
			SalePrice: product.SalePrice,
			Currency:  product.Currency,
		}
	}
	return &schema.GetRecommendationsResponse{
		Products: protoProducts,
	}, nil
}

func getClientIP(ctx context.Context) string {
	if realIP, ok := grpc_realip.FromContext(ctx); ok && realIP.IsValid() {
		return realIP.String()
	}
	return ""
}
