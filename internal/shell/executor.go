package shell

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"

	"github.com/samber/lo"

	"github.com/alexander-kolodka/crestic/internal/logger"
)

// Executor runs shell commands with full stdout/stderr logging.
// All output is duplicated to console and captured in logs.
type Executor struct{}

// Result contains the outcome of a command execution.
type Result struct {
	ExitCode int
	Stdout   string
	Stderr   string
	Error    error
}

// NewExecutor creates a new Executor with the given logger.
func NewExecutor() *Executor {
	return &Executor{}
}

// Run executes a command with timeout control and handling of ignored exit codes.
// All stdout/stderr is written to console and logs. Returns Result with exit code and output.
// If context has silent output enabled, stdout/stderr are suppressed.
func (r *Executor) Run(ctx context.Context, service string, args ...string) *Result {
	cmd := exec.CommandContext(ctx, service, args...)

	cmd.Env = os.Environ()
	cmd.Env = append(cmd.Env,
		lo.MapToSlice(getEnvVars(ctx), func(key, value string) string {
			return fmt.Sprintf("%s=%s", key, value)
		})...,
	)

	if shouldPrintCommands(ctx) {
		cmdStr := formatCommand(service, args)
		log := logger.FromContext(ctx)
		log.Info().Msg("\n" + cmdStr + "\n")
	}

	var stdoutBuf, stderrBuf bytes.Buffer

	if isSilent(ctx) {
		cmd.Stdout = &stdoutBuf
		cmd.Stderr = &stderrBuf
	} else {
		shellWriter := logger.NewShellWriter(ctx)
		cmd.Stdout = io.MultiWriter(shellWriter, &stdoutBuf)
		cmd.Stderr = io.MultiWriter(shellWriter, &stderrBuf)
	}

	err := cmd.Run()

	result := &Result{
		ExitCode: 0,
		Stdout:   stdoutBuf.String(),
		Stderr:   stderrBuf.String(),
		Error:    err,
	}

	var exitError *exec.ExitError
	if err != nil && errors.As(err, &exitError) {
		result.ExitCode = exitError.ExitCode()
		result.Error = fmt.Errorf(`%s: %s`, exitError.Error(), result.Stderr)
	}

	return result
}

func formatCommand(cmd string, args []string) string {
	b := new(strings.Builder)
	b.WriteString(cmd)

	for _, a := range args {
		b.WriteByte(' ')

		// Quote if contains spaces or special characters
		if strings.ContainsAny(a, " \t\n\"'`$\\") {
			b.WriteString(`"`)
			b.WriteString(strings.ReplaceAll(a, `"`, `\"`))
			b.WriteString(`"`)
			continue
		}

		b.WriteString(a)
	}

	return b.String()
}
