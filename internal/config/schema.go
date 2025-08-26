package config

import (
	"time"

	"bitbucket.org/plaxieappier/rec-go-kit/logkit"
	"bitbucket.org/plaxieappier/rec-go-kit/tracekit"
)

type Config struct {
	Logging         logkit.Config   `mapstructure:"logging"`
	EnableGinLogger bool            `mapstructure:"enable_gin_logger"`
	Tracing         tracekit.Config `mapstructure:"tracing"`
	VendorConfig    VendorConfig    `mapstructure:"vendor_config"`
}

type VendorConfig struct {
	ProxyURL string        `mapstructure:"proxy_url"`
	Timeout  time.Duration `mapstructure:"timeout"`
	Vendors  []Vendor      `mapstructure:"vendors"`
}

type Vendor struct {
	Name        string `mapstructure:"name"`
	RequestURL  string `mapstructure:"request_url"`
	TrackingURL string `mapstructure:"tracking_url"`
	WithProxy   bool   `mapstructure:"with_proxy"`
}
