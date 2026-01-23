package vendor

import (
	"context"
	"fmt"
	"net/url"

	"rec-vendor-api/internal/config"
)

type Service interface {
	GetRecommendations(ctx context.Context, vendorKey string, req Request) ([]ProductInfo, error)
	GetVendors(ctx context.Context) ([]VendorInfo, error)
}

type VendorInfo struct {
	VendorKey   string
	RequestHost string
}

type ServiceImpl struct {
	vendorRegistry map[string]Client
	vendorConfig   config.VendorConfig
}

func NewService(vendorRegistry map[string]Client, vendorConfig config.VendorConfig) Service {
	return &ServiceImpl{
		vendorRegistry: vendorRegistry,
		vendorConfig:   vendorConfig,
	}
}

func (s *ServiceImpl) GetRecommendations(ctx context.Context, vendorKey string, req Request) ([]ProductInfo, error) {
	vendorClient := s.vendorRegistry[vendorKey]
	if vendorClient == nil {
		return nil, fmt.Errorf("vendor not found")
	}
	products, err := vendorClient.GetUserRecommendationItems(ctx, req)
	if err != nil {
		return nil, err
	}
	return products, nil
}

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
