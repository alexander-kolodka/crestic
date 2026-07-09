package harness

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"testing"
	"text/template"

	"github.com/stretchr/testify/require"
)

const (
	dirPerm  = 0o750
	filePerm = 0o600
)

// Sandbox is an isolated integration test environment.
type Sandbox struct {
	t          *testing.T
	root       string
	configPath string
	config     configData
}

// New creates an empty sandbox with only a temp root and state directory.
func New(t *testing.T) *Sandbox {
	t.Helper()

	root := t.TempDir()
	stateDir := filepath.Join(root, ".crestic")
	require.NoError(t, os.MkdirAll(stateDir, dirPerm))

	return &Sandbox{
		t:      t,
		root:   root,
		config: configData{},
	}
}

// Mkdir creates a directory under the sandbox root and returns its absolute path.
func (s *Sandbox) Mkdir(dir string) string {
	s.t.Helper()

	path := filepath.Join(s.root, dir)
	require.NoError(s.t, os.MkdirAll(path, dirPerm))
	return path
}

// WriteFile creates a file inside dir with the given name and content.
func (s *Sandbox) WriteFile(dir, name, content string) {
	s.t.Helper()

	path := filepath.Join(s.root, dir, name)
	require.NoError(s.t, os.WriteFile(path, []byte(content), filePerm))
}

// AddRepo creates a repository directory and registers it in the config.
func (s *Sandbox) AddRepo(name string) *RepoDir {
	s.t.Helper()

	path := filepath.Join(s.root, "repos", name)
	require.NoError(s.t, os.MkdirAll(path, dirPerm))

	repo := &RepoDir{
		sandbox: s,
		Name:    name,
		Path:    path,
	}

	s.config.Repos = append(s.config.Repos, repoEntry{
		Name:        name,
		Path:        path,
		PasswordCMD: defaultPasswordCMD,
	})

	return repo
}

// AddPipeline creates a pipeline with the given name. Cron is optional.
func (s *Sandbox) AddPipeline(name string) *PipelineBuilder {
	s.t.Helper()

	entry := pipelineEntry{Name: name}
	s.config.Pipelines = append(s.config.Pipelines, &entry)

	return &PipelineBuilder{
		sandbox:  s,
		pipeline: &entry,
	}
}

func (s *Sandbox) ensureConfig() {
	if s.configPath != "" {
		return
	}

	tmpl, err := template.New("config").Parse(configTemplate)
	require.NoError(s.t, err)

	var buf bytes.Buffer
	require.NoError(s.t, tmpl.Execute(&buf, s.config))

	s.configPath = filepath.Join(s.root, "crestic.yaml")
	require.NoError(s.t, os.WriteFile(s.configPath, buf.Bytes(), filePerm))
}

// Run executes crestic CLI with the sandbox config and returns captured stdout.
func (s *Sandbox) Run(ctx context.Context, args ...string) (string, error) {
	s.t.Helper()
	s.ensureConfig()

	fullArgs := append([]string{"--config", s.configPath, "--log-level", "error"}, args...)
	command := exec.CommandContext(ctx, CresticBin(), fullArgs...)
	command.Dir = s.root

	var stdout bytes.Buffer
	var stderr bytes.Buffer
	command.Stdout = &stdout
	command.Stderr = &stderr

	err := command.Run()
	if err != nil {
		return stdout.String(), fmt.Errorf("%w\nstderr:\n%s", err, stderr.String())
	}

	return stdout.String(), nil
}
