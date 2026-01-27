package vendor

import (
	"context"
	"errors"
	"net/url"

	"rec-vendor-api/internal/config"
	controller_errors "rec-vendor-api/internal/controller/errors"

	log "github.com/sirupsen/logrus"
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
		log.WithContext(ctx).Errorf("Invalid vendor key: %s", vendorKey)
		return nil, controller_errors.BadRequestErrorf("Vendor key '%s' not supported", vendorKey)
	}
	products, err := vendorClient.GetUserRecommendationItems(ctx, req)
	if err != nil {
		var badRequestErr *controller_errors.BadRequestError
		if errors.As(err, &badRequestErr) {
			log.WithContext(ctx).Errorf("VendorClient returned BadRequestError. err: %v", err)
			return nil, controller_errors.BadRequestErrorf("VendorClient returned BadRequestError. err: %v", err)
		}
		log.WithContext(ctx).Errorf("Fail to recommend any products. err: %v", err)
		return nil, controller_errors.InternalServerErrorf("Fail to recommend any products. err: %v", err)
	}
	return products, nil
}

func (s *ServiceImpl) GetVendors(ctx context.Context) ([]VendorInfo, error) {
	vendors := make([]VendorInfo, 0, len(s.vendorConfig.Vendors))
	for _, v := range s.vendorConfig.Vendors {
		requestHost := ""
		if parsedURL, err := url.Parse(v.Request.URL); err == nil {
			requestHost = parsedURL.Host
		} else {
			return nil, controller_errors.BadRequestErrorf("failed to parse request URL for vendor %s: %s", v.Name, err.Error())
		}
		vendors = append(vendors, VendorInfo{
			VendorKey:   v.Name,
			RequestHost: requestHost,
		})
	}
	return vendors, nil
}
