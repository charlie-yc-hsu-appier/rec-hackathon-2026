package grpc

import (
	"context"

	"rec-vendor-api/internal/config"
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

func NewAPIServer(vendorRegistry map[string]vendor.Client, vendorConfig config.VendorConfig) *APIServerImpl {
	return &APIServerImpl{
		vendorService: vendor.NewService(vendorRegistry, vendorConfig),
	}
}

func (s *APIServerImpl) GetRecommendations(ctx context.Context, req *schema.GetRecommendationsRequest) (*schema.GetRecommendationsResponse, error) {
	// TODO: Implement logic to convert protobuf request to internal Request
	// and call vendorService.GetRecommendations
	return nil, nil
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
	// TODO: Implement health check logic
	return &schema.HealthcheckResponse{
		Status: "ok",
	}, nil
}
