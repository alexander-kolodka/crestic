package hooks

import (
	"context"

	"github.com/alexander-kolodka/crestic/internal/engine/shell"
)

// WithJobEnv sets CRESTIC_JOB_NAME and optionally CRESTIC_ERROR for job hooks.
func WithJobEnv(ctx context.Context, jobName string, err error) context.Context {
	env := map[string]string{
		"CRESTIC_JOB_NAME": jobName,
	}
	if err != nil {
		env["CRESTIC_ERROR"] = err.Error()
	}
	return shell.WithEnv(ctx, env)
}

// WithPipelineEnv sets CRESTIC_PIPELINE_NAME and optionally CRESTIC_ERROR for pipeline hooks.
func WithPipelineEnv(ctx context.Context, pipelineName string, err error) context.Context {
	env := map[string]string{
		"CRESTIC_PIPELINE_NAME": pipelineName,
	}
	if err != nil {
		env["CRESTIC_ERROR"] = err.Error()
	}
	return shell.WithEnv(ctx, env)
}
