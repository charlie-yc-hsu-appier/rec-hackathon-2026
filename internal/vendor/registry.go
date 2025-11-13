package vendor

import (
	"rec-vendor-api/internal/config"
	"rec-vendor-api/internal/strategy"

	"github.com/plaxieappier/rec-go-kit/httpkit"
)

func BuildRegistry(config config.VendorConfig) (map[string]Client, error) {
	registry := map[string]Client{}

	// Initialize two http clients: one with proxy, one without
	httpProxyClient, err := httpkit.NewClient(httpkit.WithProxy(config.ProxyURL))
	if err != nil {
		return nil, err
	}
	httpClient, err := httpkit.NewClient()
	if err != nil {
		return nil, err
	}

	httpClients := map[bool]httpkit.Client{
		true:  httpProxyClient,
		false: httpClient,
	}

	for _, v := range config.Vendors {
		client := NewClient(
			v,
			httpClients[v.WithProxy],
			config.Timeout,
			strategy.BuildHeader(v),
			strategy.BuildRequest(v),
			strategy.BuildBody(v),
			strategy.BuildUnmarshaler(v),
			strategy.BuildTracking(v),
		)

		registry[v.Name] = client
	}
	return registry, nil
}
