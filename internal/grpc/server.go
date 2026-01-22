package grpc

import (
	"context"

	"rec-vendor-api/internal/vendor"

	"google.golang.org/protobuf/types/known/emptypb"

	// TODO: not yet generated at rec-schema yet
	schema "github.com/plaxieappier/rec-schema/go/vendorapi"
)

// APIServerImpl implements the VendorAPIServer interface from generated protobuf code
// Note: This will compile once the protobuf code is generated via: protoc --go_out=. --go-grpc_out=. api/vendor_api.proto
type APIServerImpl struct {
	schema.UnimplementedVendorAPIServer
	vendorService vendor.Service
}

// NewAPIServer creates a new gRPC API server
func NewAPIServer(vendorService vendor.Service) *APIServerImpl {
	return &APIServerImpl{
		vendorService: vendorService,
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
	return &schema.GetVendorsResponse{
		Vendors: []*schema.VendorInfo{
			{
				VendorKey:   vendors[0].VendorKey,
				RequestHost: vendors[0].RequestHost,
			},
		},
	}, nil
}

func (s *APIServerImpl) CheckHealthCheck(ctx context.Context, req *emptypb.Empty) (*schema.HealthcheckResponse, error) {
	// TODO: Implement health check logic
	return &schema.HealthcheckResponse{
		Status: "ok",
	}, nil
}
