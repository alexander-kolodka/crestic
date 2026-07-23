package paths

import (
	"fmt"
	"path/filepath"
)

// Canonical returns the absolute path with symlinks resolved.
func Canonical(path string) (string, error) {
	abs, err := filepath.Abs(path)
	if err != nil {
		return "", fmt.Errorf("resolve path: %w", err)
	}

	clean, err := filepath.EvalSymlinks(abs)
	if err != nil {
		return "", fmt.Errorf("resolve path symlinks: %w", err)
	}

	return clean, nil
}
