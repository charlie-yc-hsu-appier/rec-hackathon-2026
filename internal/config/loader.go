package config

import (
	"errors"
	"fmt"

	"github.com/go-playground/validator/v10"
	log "github.com/sirupsen/logrus"

	"path"

	"github.com/kelseyhightower/envconfig"
	"github.com/spf13/viper"
)

func Load(configPath string, cfg *Config) error {
	configName := path.Base(configPath)
	ext := path.Ext(configPath)
	dir := path.Dir(configPath)

	if configPath == "" {
		return errors.New("config path is empty")
	}
	if ext != ".yaml" {
		return errors.New("only accept .yaml file")
	}

	v := viper.New()
	v.SetConfigType("yaml")
	v.SetConfigName(configName)
	v.AddConfigPath(dir)

	if err := v.ReadInConfig(); err != nil {
		return err
	}

	// Unmarshal the configuration into the struct
	if err := v.Unmarshal(cfg); err != nil {
		return err
	}

	validate := validator.New()
	if err := validate.Struct(cfg); err != nil {
		return fmt.Errorf("config validation failed: %w", err)
	}

	return nil
}

func LoadConfigFromEnv(cfg any) {
	err := envconfig.Process("", cfg)
	if err != nil {
		log.Fatalf("failed to load env config: %v", err)
	}
}
