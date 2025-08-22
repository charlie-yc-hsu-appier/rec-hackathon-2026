package vendor

import (
	"rec-vendor-api/internal/config"
	"rec-vendor-api/internal/strategy"

	"bitbucket.org/plaxieappier/rec-go-kit/httpkit"
)

func BuildRegistry(config config.VendorConfig) (map[string]Client, error) {
	registry := map[string]Client{}
	for _, v := range config.Vendors {
		var opts []httpkit.ClientOption
		if v.WithProxy {
			opts = append(opts, httpkit.WithProxy(config.ProxyURL))
		}
		httpClient, err := httpkit.NewClient(opts...)
		if err != nil {
			return nil, err
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
	return registry, nil
}
