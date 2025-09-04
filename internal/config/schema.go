package config

import (
	"time"

	"github.com/plaxieappier/rec-go-kit/logkit"
	"github.com/plaxieappier/rec-go-kit/tracekit"
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
	Name        string            `mapstructure:"name"`
	RequestURL  string            `mapstructure:"request_url"`
	TrackingURL string            `mapstructure:"tracking_url"`
	WithProxy   bool              `mapstructure:"with_proxy"`
	AccessKey   string            `mapstructure:"access_key"`
	SecretKey   string            `mapstructure:"secret_key"`
	SizeCodeMap map[string]string `mapstructure:"size_code_map"`
	UserAgent   string            `mapstructure:"user_agent"`
}
