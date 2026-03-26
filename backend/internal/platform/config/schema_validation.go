package config

import (
	"fmt"
	"slices"
	"strings"
)

func validateKnownConfigPaths(settings map[string]any) error {
	if len(settings) == 0 || len(generatedKnownConfigPaths) == 0 {
		return nil
	}

	seen := make(map[string]struct{})
	collectConfigPaths("", settings, seen)

	var unknown []string
	for path := range seen {
		if _, ok := generatedKnownConfigPaths[path]; ok {
			continue
		}

		unknown = append(unknown, path)
	}

	if len(unknown) == 0 {
		return nil
	}

	slices.Sort(unknown)

	return fmt.Errorf("unknown config paths: %s", strings.Join(unknown, ", "))
}

func collectConfigPaths(prefix string, value any, seen map[string]struct{}) {
	switch typed := value.(type) {
	case map[string]any:
		if prefix != "" {
			seen[prefix] = struct{}{}
		}

		for key, child := range typed {
			collectConfigPaths(joinConfigPath(prefix, key), child, seen)
		}
	case map[any]any:
		normalized := make(map[string]any, len(typed))
		for rawKey, child := range typed {
			key, ok := rawKey.(string)
			if !ok {
				continue
			}
			normalized[key] = child
		}
		collectConfigPaths(prefix, normalized, seen)
	case []any:
		if prefix != "" {
			seen[prefix] = struct{}{}
		}
		for _, child := range typed {
			collectConfigPaths(prefix, child, seen)
		}
	default:
		if prefix != "" {
			seen[prefix] = struct{}{}
		}
	}
}

func joinConfigPath(prefix, key string) string {
	if prefix == "" {
		return key
	}

	return prefix + "." + key
}
