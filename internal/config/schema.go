package config

import "bitbucket.org/plaxieappier/rec-go-kit/logkit"

type Config struct {
	Logging logkit.Config `mapstructure:"logging"`
}
