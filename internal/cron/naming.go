package cron

import (
	"crypto/md5" //nolint:gosec // G501: md5 used only for cron artifact name fingerprint, not security
	"encoding/hex"
	"fmt"
	"path/filepath"
)

// CanonicalConfigPath returns the absolute config path with symlinks resolved.
func CanonicalConfigPath(path string) (string, error) {
	abs, err := filepath.Abs(path)
	if err != nil {
		return "", fmt.Errorf("resolve config path: %w", err)
	}

	clean, err := filepath.EvalSymlinks(abs)
	if err != nil {
		return "", fmt.Errorf("resolve config symlinks: %w", err)
	}

	return clean, nil
}

// LockFileName returns the cron lock filename for a config.
// cfgPath is the canonical config path; cfgBasename is the config filename without extension.
func LockFileName(cfgPath, cfgBasename string) string {
	hash := configHash(cfgPath)
	return fmt.Sprintf("crestic-cron-%s-%s.lock", cfgBasename, hash)
}

// StateFileName returns the cron state filename for a config.
// cfgPath is the canonical config path; cfgBasename is the config filename without extension.
func StateFileName(cfgPath, cfgBasename string) string {
	hash := configHash(cfgPath)
	return fmt.Sprintf("crestic-cron-state-%s-%s.json", cfgBasename, hash)
}

func configHash(cfgPath string) string {
	//nolint:gosec // G401: md5 is used only for cron artifact name fingerprint, not security
	sum := md5.Sum([]byte(cfgPath))
	return hex.EncodeToString(sum[:])
}
