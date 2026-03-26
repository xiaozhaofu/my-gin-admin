package config

import (
	"io"

	legacyconfig "go_sleep_admin/config"
)

func EmbeddedConfig(env string) (io.Reader, error) {
	return legacyconfig.EmbeddedConfig(env), nil
}
