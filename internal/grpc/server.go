package grpc

import (
	"context"

	"rec-vendor-api/internal/grpcutils"
	"rec-vendor-api/internal/vendor"

	"google.golang.org/protobuf/types/known/emptypb"

<<<<<<< HEAD
	// TODO: not yet generated at rec-schema yet
	schema "github.com/plaxieappier/rec-schema/go/vendorapi"
=======
	schema "github.com/plaxieappier/rec-schema/go/vendor"
>>>>>>> 6ffcbbe (start up grpc server)
)

type APIServerImpl struct {
	schema.UnimplementedVendorAPIServer
	vendorService vendor.Service
}

func NewAPIServer(service vendor.Service) *APIServerImpl {
	return &APIServerImpl{
		vendorService: service,
	}
}

func (s *APIServerImpl) GetRecommendations(ctx context.Context, req *schema.GetRecommendationsRequest) (*schema.GetRecommendationsResponse, error) {
	// Convert protobuf request to internal Request struct
	internalReq := convertToInternalRequest(ctx, req)

	recommendations, err := s.vendorService.GetRecommendations(ctx, req.VendorKey, internalReq)
	if err != nil {
		return nil, err
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

func convertToInternalRequest(ctx context.Context, req *schema.GetRecommendationsRequest) vendor.Request {
	clientIP := grpcutils.GetClientIPFromContext(ctx)

	osStr := ""
	if req.Os == schema.OperationSystem_ANDROID {
		osStr = "android"
	} else if req.Os == schema.OperationSystem_IOS {
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

func (s *APIServerImpl) GetVendors(ctx context.Context, _ *emptypb.Empty) (*schema.GetVendorsResponse, error) {
	vendors, err := s.vendorService.GetVendors(ctx)
	if err != nil {
		return nil, err
	}

	protoVendors := make([]*schema.VendorInfo, 0, len(vendors))
	for _, v := range vendors {
		protoVendors = append(protoVendors, &schema.VendorInfo{
			VendorKey:   v.VendorKey,
			RequestHost: v.RequestHost,
		})
	}
	return &schema.GetVendorsResponse{
<<<<<<< HEAD
		Vendors: []*schema.VendorInfo{
			{
				VendorKey:   vendors[0].VendorKey,
				RequestHost: vendors[0].RequestHost,
			},
		},
=======
		Vendors: protoVendors,
>>>>>>> 6ffcbbe (start up grpc server)
	}, nil
}

func (s *APIServerImpl) CheckHealthCheck(ctx context.Context, req *emptypb.Empty) (*schema.HealthcheckResponse, error) {

	return &schema.HealthcheckResponse{
		Status: "ok",
	}, nil
}
