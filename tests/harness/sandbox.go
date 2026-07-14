package harness

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"text/template"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/alexander-kolodka/crestic/internal/cron"
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
	t.Parallel()

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

// WriteCronState seeds the cron state file for a pipeline with the given last run time.
func (s *Sandbox) WriteCronState(pipeline string, lastRun time.Time) {
	s.t.Helper()
	s.ensureConfig()

	canonicalPath, err := cron.CanonicalConfigPath(s.configPath)
	require.NoError(s.t, err)

	cfgBasename := strings.TrimSuffix(filepath.Base(canonicalPath), filepath.Ext(canonicalPath))
	stateFileName := cron.StateFileName(canonicalPath, cfgBasename)

	state := cron.State{
		Pipelines: map[string]cron.PipelineState{
			pipeline: {LastRun: lastRun},
		},
	}

	data, marshalErr := json.Marshal(state)
	require.NoError(s.t, marshalErr)

	statePath := filepath.Join(s.root, ".crestic", stateFileName)
	err = os.WriteFile(statePath, data, filePerm)
	require.NoError(s.t, err)
}

// Run executes crestic CLI with the sandbox config and returns captured stdout.
func (s *Sandbox) Run(ctx context.Context, args ...string) (string, error) {
	s.t.Helper()
	s.ensureConfig()

	fullArgs := append([]string{"--config", s.configPath, "--log-level", "error"}, args...)
	command := exec.CommandContext(ctx, CresticBin(), fullArgs...)
	command.Dir = s.root
	command.Env = append(os.Environ(), "HOME="+s.root)

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
