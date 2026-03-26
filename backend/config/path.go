package config

import "go_sleep_admin/internal/runtime/resource"

func ResolvePath(rel string) string {
	return resource.ResolvePath(rel)
}
