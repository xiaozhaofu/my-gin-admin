package config

import (
	"fmt"
	"io"
	"os"

	"go_sleep_admin/internal/runtime/resource"

	"github.com/spf13/viper"
)

const defaultConfigType = "yml"

// Load loads typed config from disk or embedded fallback.
func Load(envFile string) (*Config, error) {
	if envFile == "" {
		return nil, fmt.Errorf("env file is required")
	}

	configPath := resource.ResolvePath("config/env/" + envFile + ".yml")
	if _, err := os.Stat(configPath); err == nil {
		v := viper.New()
		v.SetConfigFile(configPath)
		v.SetConfigType(defaultConfigType)
		v.AutomaticEnv()

		if err := v.ReadInConfig(); err != nil {
			return nil, fmt.Errorf("read config file %s: %w", configPath, err)
		}

		return decode(envFile, v)
	}

	reader, err := EmbeddedConfig(envFile)
	if err != nil {
		return nil, fmt.Errorf("load embedded config: %w", err)
	}

	return loadFromReader(envFile, reader)
}

func loadFromReader(envFile string, r io.Reader) (*Config, error) {
	v := viper.New()
	v.SetConfigType(defaultConfigType)

	if err := v.ReadConfig(r); err != nil {
		return nil, fmt.Errorf("read embedded config: %w", err)
	}

	return decode(envFile, v)
}

func decode(envFile string, v *viper.Viper) (*Config, error) {
	if err := validateKnownConfigPaths(v.AllSettings()); err != nil {
		return nil, fmt.Errorf("validate config schema: %w", err)
	}

	var cfg Config
	if err := v.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("decode config: %w", err)
	}

	cfg.Env = envFile
	if err := cfg.Validate(); err != nil {
		return nil, fmt.Errorf("validate config: %w", err)
	}

	return &cfg, nil
}
