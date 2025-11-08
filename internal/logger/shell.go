package logger

import (
	"context"
	"encoding/json"
	"io"
	"strings"
)

// ShellWriter wraps zerolog for shell command output with indentation.
type ShellWriter struct {
	ctx context.Context
}

func NewShellWriter(ctx context.Context) io.Writer {
	return &ShellWriter{ctx: ctx}
}

func (w *ShellWriter) Write(p []byte) (int, error) {
	if len(p) == 0 {
		return 0, nil
	}

	msg := string(p)
	lines := strings.Split(strings.TrimRight(msg, "\n"), "\n")
	source := GetSource(w.ctx)

	if !shouldParseJSON(w.ctx, source) {
		return w.logLines(p, lines)
	}

	return w.logJSON(p, lines)
}

func (w *ShellWriter) logJSON(p []byte, lines []string) (int, error) {
	log := FromContext(w.ctx)
	for _, line := range lines {
		if strings.TrimSpace(line) == "" {
			continue
		}

		var jsonObj any
		err := json.Unmarshal([]byte(line), &jsonObj)
		if err != nil {
			log.Debug().Err(err).Str("line", line).
				Msg("error parsing json line from restic output")
			continue
		}

		log.Info().Interface("restic_output", jsonObj).Msg("restic output")
	}

	return len(p), nil
}

func (w *ShellWriter) logLines(p []byte, lines []string) (int, error) {
	log := FromContext(w.ctx)

	for _, line := range lines {
		if line != "" {
			log.Info().Msg(line)
		}
	}

	return len(p), nil
}

func shouldParseJSON(ctx context.Context, source string) bool {
	if source != "restic" {
		return false
	}

	return IsJSONMode(ctx)
}
