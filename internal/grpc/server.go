package grpc

import (
	"context"

	"rec-vendor-api/internal/controller/errors"
	"rec-vendor-api/internal/vendor"

	"google.golang.org/protobuf/types/known/emptypb"

	schema "github.com/plaxieappier/rec-schema/go/vendorapi"
)

type APIServer interface {
	GetRecommendations(context.Context, *schema.GetRecommendationsRequest) (*schema.GetRecommendationsResponse, error)
	GetVendors(context.Context, *emptypb.Empty) (*schema.GetVendorsResponse, error)
	CheckHealthCheck(context.Context, *emptypb.Empty) (*schema.HealthcheckResponse, error)
}

type APIServerImpl struct {
	schema.UnimplementedVendorAPIServer
	vendorService vendor.Service
}

func NewAPIServer(service vendor.Service) APIServer {
	return &APIServerImpl{
		vendorService: service,
	}
}

func (s *APIServerImpl) GetRecommendations(ctx context.Context, req *schema.GetRecommendationsRequest) (*schema.GetRecommendationsResponse, error) {
	internalReq := convertToInternalRequest(ctx, req)

	recommendations, err := s.vendorService.GetRecommendations(ctx, req.VendorKey, internalReq)
	if err != nil {
		return nil, errors.ToGRPCStatus(err)
	}

	products := make([]*schema.ProductInfo, 0, len(recommendations))
	for _, recommendation := range recommendations {
		products = append(products, &schema.ProductInfo{
			ProductId: recommendation.ProductID,
			Url:       recommendation.Url,
			Image:     recommendation.Image,
			Price:     recommendation.Price,
			SalePrice: recommendation.SalePrice,
			Currency:  recommendation.Currency,
		})
	}
	return &schema.GetRecommendationsResponse{
		Products: products,
	}, nil
}

func (s *APIServerImpl) GetVendors(ctx context.Context, _ *emptypb.Empty) (*schema.GetVendorsResponse, error) {
	vendors, err := s.vendorService.GetVendors(ctx)
	if err != nil {
		return nil, errors.ToGRPCStatus(err)
	}

	protoVendors := make([]*schema.VendorInfo, 0, len(vendors))
	for _, v := range vendors {
		protoVendors = append(protoVendors, &schema.VendorInfo{
			VendorKey:   v.VendorKey,
			RequestHost: v.RequestHost,
		})
	}
	return &schema.GetVendorsResponse{
		Vendors: protoVendors,
	}, nil
}

func (s *APIServerImpl) CheckHealthCheck(ctx context.Context, req *emptypb.Empty) (*schema.HealthcheckResponse, error) {

	return &schema.HealthcheckResponse{
		Status: "ok",
	}, nil
}

func convertToInternalRequest(ctx context.Context, req *schema.GetRecommendationsRequest) vendor.Request {
	clientIP := GetClientIPFromContext(ctx)

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
		ClientIP:        clientIP,
	}
}
