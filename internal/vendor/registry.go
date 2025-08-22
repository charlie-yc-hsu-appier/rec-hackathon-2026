package vendor

import (
	"log"
	"rec-vendor-api/internal/config"
	"rec-vendor-api/internal/strategy"

	"bitbucket.org/plaxieappier/rec-go-kit/httpkit"
)

func BuildRegistry(config config.VendorConfig) map[string]Client {
	registry := map[string]Client{}
	for _, v := range config.Vendors {
		var httpClient httpkit.Client
		if v.WithProxy {
			httpClient = buildHTTPClient(httpkit.WithProxy(config.ProxyURL))
		} else {
			httpClient = buildHTTPClient()
		}

		client := NewClient(
			v,
			httpClient,
			config.Timeout,
			strategy.BuildHeader(v.Name),
			strategy.BuildRequester(v.Name),
			strategy.BuildUnmarshaler(v.Name),
			strategy.BuildTracker(v.Name))

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
