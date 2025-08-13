package config

import (
	"bitbucket.org/plaxieappier/rec-go-kit/logkit"
	"bitbucket.org/plaxieappier/rec-go-kit/tracekit"
)

type Config struct {
	Logging         logkit.Config   `mapstructure:"logging"`
	EnableGinLogger bool            `mapstructure:"enable_gin_logger"`
	Tracing         tracekit.Config `mapstructure:"tracing"`
}
