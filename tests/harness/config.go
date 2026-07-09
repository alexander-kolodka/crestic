package harness

import (
	_ "embed"
	"path/filepath"
)

//go:embed testdata/config.yaml.tmpl
var configTemplate string

const (
	defaultPasswordCMD = "echo testpass"
)

type repoEntry struct {
	Name        string
	Path        string
	PasswordCMD string
}

type jobEntry struct {
	Type     string
	Name     string
	From     []string
	FromRepo string
	To       string
	Options  map[string]any
}

type pipelineEntry struct {
	Name string
	Cron string
	Jobs []jobEntry
}

type configData struct {
	Repos     []repoEntry
	Pipelines []*pipelineEntry
}

func CresticBin() string {
	bin, _ := filepath.Abs(filepath.Join("..", "bin", "crestic"))
	return bin
}
