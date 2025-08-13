package config

import (
	"bitbucket.org/plaxieappier/rec-go-kit/logkit"
	"bitbucket.org/plaxieappier/rec-go-kit/tracekit"
)

type Config struct {
	Logging logkit.Config   `mapstructure:"logging"`
	Tracing tracekit.Config `mapstructure:"tracing"`
}
