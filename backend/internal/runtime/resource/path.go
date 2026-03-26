package resource

import (
	"os"
	"path/filepath"
	"runtime"
)

func ResolvePath(rel string) string {
	rel = filepath.Clean(rel)
	candidates := buildPathCandidates(rel)

	for _, candidate := range candidates {
		if _, err := os.Stat(candidate); err == nil {
			return candidate
		}
	}

	return candidates[0]
}

func buildPathCandidates(rel string) []string {
	candidates := make([]string, 0, 4)

	if root := projectRoot(); root != "" {
		candidates = append(candidates,
			filepath.Join(root, rel),
			filepath.Join(root, "v2", rel),
		)
	}

	candidates = append(candidates,
		rel,
		filepath.Join("v2", rel),
	)

	return uniquePaths(candidates)
}

func projectRoot() string {
	_, file, _, ok := runtime.Caller(0)
	if !ok {
		return ""
	}

	return filepath.Dir(filepath.Dir(filepath.Dir(filepath.Dir(file))))
}

func uniquePaths(paths []string) []string {
	seen := make(map[string]struct{}, len(paths))
	result := make([]string, 0, len(paths))

	for _, path := range paths {
		path = filepath.Clean(path)
		if _, ok := seen[path]; ok {
			continue
		}

		seen[path] = struct{}{}
		result = append(result, path)
	}

	return result
}
