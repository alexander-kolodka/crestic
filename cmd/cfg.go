package cmd

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"

	"github.com/alexander-kolodka/crestic/internal/dto"
	"github.com/alexander-kolodka/crestic/internal/entity"
)

// findConfigFile searches for config file in priority order:
// 1. Custom file specified with configPath parameter
// 2. ./crestic.yaml (current directory)
// 3. ~/crestic.yaml (home directory)
// 4. ~/.crestic/crestic.yaml
// 5. ~/.config/crestic/crestic.yaml.
func findConfigFile(configPath string) (string, error) {
	if configPath != "" {
		_, err := os.Stat(configPath)
		if err != nil {
			return "", fmt.Errorf("config file %s: %w", configPath, err)
		}
		return configPath, nil
	}

	home, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("get home directory: %w", err)
	}

	// Priority paths to search
	paths := []string{
		"./crestic.yaml",
		filepath.Join(home, "crestic.yaml"),
		filepath.Join(home, ".crestic", "crestic.yaml"),
		filepath.Join(home, ".config", "crestic", "crestic.yaml"),
	}

	for _, path := range paths {
		_, err = os.Stat(path)
		if err == nil {
			return path, nil
		}
	}

	return "", errors.New("no config file found in default locations " +
		"(./crestic.yaml, ~/crestic.yaml, ~/.crestic/crestic.yaml, ~/.config/crestic/crestic.yaml)",
	)
}

func loadConfig(cfgPath string) (*entity.Config, error) {
	cfgPath, err := findConfigFile(cfgPath)
	if err != nil {
		return nil, err
	}

	b, err := os.ReadFile(cfgPath)
	if err != nil {
		return nil, fmt.Errorf("read config %s: %w", cfgPath, err)
	}

	var cfg dto.Config
	err = yaml.Unmarshal(b, &cfg)
	if err != nil {
		return nil, fmt.Errorf("parse yaml from %s: %w", cfgPath, err)
	}

	entityCfg, err := dto.ToEntity(cfg)
	if err != nil {
		return nil, err
	}

	return entityCfg, nil
}
