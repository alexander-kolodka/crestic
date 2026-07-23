package cron

import (
	"crypto/md5" //nolint:gosec // G501: md5 used only for cron artifact name fingerprint, not security
	"encoding/hex"
	"fmt"
	"path/filepath"
	"strings"
)

// LockFileName returns the cron lock filename for a config.
// cfgPath is the canonical config path.
func LockFileName(cfgPath string) string {
	hash := configHash(cfgPath)
	return fmt.Sprintf("crestic-cron-%s-%s.lock", configBasename(cfgPath), hash)
}

// StateFileName returns the cron state filename for a config.
// cfgPath is the canonical config path.
func StateFileName(cfgPath string) string {
	hash := configHash(cfgPath)
	return fmt.Sprintf("crestic-cron-state-%s-%s.json", configBasename(cfgPath), hash)
}

func configBasename(cfgPath string) string {
	return strings.TrimSuffix(filepath.Base(cfgPath), filepath.Ext(cfgPath))
}

func configHash(cfgPath string) string {
	//nolint:gosec // G401: md5 is used only for cron artifact name fingerprint, not security
	sum := md5.Sum([]byte(cfgPath))
	return hex.EncodeToString(sum[:])
}
