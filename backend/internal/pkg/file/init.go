package file

import (
	"os"

	v2config "go_sleep_admin/internal/platform/config"
	"go_sleep_admin/internal/runtime/resource"
)

func InitWithConfig(cfg *v2config.Config) {
	if cfg == nil {
		return
	}

	if cfg.Upload.Driver != "" && cfg.Upload.Driver != "local" {
		return
	}

	_ = os.MkdirAll(resource.ResolvePath(cfg.Upload.LocalDir), 0o755)
}
