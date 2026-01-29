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
	Grpc            GrpcConfig      `mapstructure:"grpc"`
}
type GrpcConfig struct {
	MaxConnectionAge  time.Duration `mapstructure:"max_connection_age"`
	WriteBufferSizeKb int           `mapstructure:"write_buffer_size_kb"`
	ReadBufferSizeKb  int           `mapstructure:"read_buffer_size_kb"`
}

type VendorConfig struct {
	ProxyURL string        `mapstructure:"proxy_url"`
	Timeout  time.Duration `mapstructure:"timeout"`
	Vendors  []Vendor      `mapstructure:"vendors"`
}

type Vendor struct {
	Name         string     `mapstructure:"name"`
	WithProxy    bool       `mapstructure:"with_proxy"`
	HTTPMethod   string     `mapstructure:"http_method"`
	AccessKey    string     `mapstructure:"access_key"`
	SecretKey    string     `mapstructure:"secret_key"`
	UserAgent    string     `mapstructure:"user_agent"`
	SceneType    string     `mapstructure:"scene_type"`
	Ver          string     `mapstructure:"ver"`
	ChannelToken string     `mapstructure:"channel_token"`
	SCaApp       string     `mapstructure:"s_ca_app"`
	SCaSecret    string     `mapstructure:"s_ca_secret"`
	Request      URLPattern `mapstructure:"request"`
	Tracking     URLPattern `mapstructure:"tracking"`
	ContentType  string     `mapstructure:"content_type"`
}

type URLPattern struct {
	URL     string  `mapstructure:"url"`
	Queries []Query `mapstructure:"queries,omitempty"`
}

type Query struct {
	Key   string `mapstructure:"key"`
	Value string `mapstructure:"value"`
}
