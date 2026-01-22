package vendor

import (
	"context"
	"net/url"

	"rec-vendor-api/internal/config"
)

// Service defines the business logic interface for vendor operations
type Service interface {
	GetRecommendations(ctx context.Context, req Request) ([]ProductInfo, error)
	GetVendors(ctx context.Context) ([]VendorInfo, error)
}

// VendorInfo represents vendor information
type VendorInfo struct {
	VendorKey   string
	RequestHost string
}

// ServiceImpl implements the Service interface
type ServiceImpl struct {
	vendorRegistry map[string]Client
	vendorConfig   config.VendorConfig
}

// NewService creates a new vendor service
func NewService(vendorRegistry map[string]Client, vendorConfig config.VendorConfig) Service {
	return &ServiceImpl{
		vendorRegistry: vendorRegistry,
		vendorConfig:   vendorConfig,
	}
}

// GetRecommendations will be implemented to get recommendations from a vendor
func (s *ServiceImpl) GetRecommendations(ctx context.Context, req Request) ([]ProductInfo, error) {
	// TODO: Implement logic to get recommendations
	// This will call the appropriate vendor client from the registry
	return nil, nil
}

// GetVendors returns the list of available vendors
func (s *ServiceImpl) GetVendors(ctx context.Context) ([]VendorInfo, error) {
	vendors := make([]VendorInfo, 0, len(s.vendorConfig.Vendors))
	for _, v := range s.vendorConfig.Vendors {
		requestHost := ""
		if parsedURL, err := url.Parse(v.Request.URL); err == nil {
			requestHost = parsedURL.Host
		}
		vendors = append(vendors, VendorInfo{
			VendorKey:   v.Name,
			RequestHost: requestHost,
		})
	}
	return vendors, nil
}
