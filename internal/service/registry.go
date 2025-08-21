package service

import (
	"log"
	"rec-vendor-api/internal/config"
	header "rec-vendor-api/internal/service/header_strategy"
	requester "rec-vendor-api/internal/service/request_strategy"
	trackurl "rec-vendor-api/internal/service/trackurl_strategy"
	unmarshaler "rec-vendor-api/internal/service/unmarshal_strategy"

	"bitbucket.org/plaxieappier/rec-go-kit/httpkit"
)

func InitVendors(config config.VendorConfig) map[string]Client {
	registry := map[string]Client{}
	for _, v := range config.Vendors {
		var httpClient httpkit.Client
		if v.WithProxy {
			httpClient = buildHTTPClient(httpkit.WithProxy(config.ProxyURL))
		} else {
			httpClient = buildHTTPClient()
		}

		client := NewVendorClient(
			v,
			httpClient,
			config.Timeout,
			buildHeaderStrategyByVendor(v.Name),
			buildRequestURLStrategyByVendor(v.Name),
			buildUnmarshalStrategyByVendor(v.Name),
			buildTrackingURLStrategyByVendor(v.Name))

		registry[v.Name] = client
	}
	return registry
}

func buildHTTPClient(opts ...httpkit.ClientOption) httpkit.Client {
	client, err := httpkit.NewClient(opts...)
	if err != nil {
		log.Fatalf("Fail to create http client. err: %v", err)
	}
	return client
}

func buildHeaderStrategyByVendor(name string) HeaderStrategy {
	switch name {
	default:
		return &header.NoHeaderStrategy{}
	}
}

func buildRequestURLStrategyByVendor(name string) RequestURLStrategy {
	switch name {
	default:
		return &requester.DefaultRequestURLStrategy{}
	}
}

func buildUnmarshalStrategyByVendor(name string) RespUnmarshalStrategy {
	switch name {
	default:
		return &unmarshaler.DefaultUnmarshalStrategy{}
	}
}

func buildTrackingURLStrategyByVendor(name string) TrackingURLStrategy {
	switch name {
	default:
		return &trackurl.DefaultTrackingURLStrategy{}
	}
}
