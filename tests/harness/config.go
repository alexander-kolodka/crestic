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

// Hooks are lifecycle hook commands for a pipeline or job.
type Hooks struct {
	Before  []string
	Success []string
	Failure []string
}

// Empty reports whether no hook commands are configured.
func (h Hooks) Empty() bool {
	return len(h.Before) == 0 && len(h.Success) == 0 && len(h.Failure) == 0
}

type jobEntry struct {
	Type     string
	Name     string
	From     []string
	FromRepo string
	To       string
	Options  map[string]any
	Hooks    Hooks
}

type pipelineEntry struct {
	Name  string
	Cron  string
	Hooks Hooks
	Jobs  []jobEntry
}

type configData struct {
	Repos     []repoEntry
	Pipelines []*pipelineEntry
}

func CresticBin() string {
	bin, _ := filepath.Abs(filepath.Join("..", "bin", "crestic"))
	return bin
}
