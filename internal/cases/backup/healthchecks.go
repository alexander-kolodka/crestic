package backup

import (
	"context"

	"github.com/alexander-kolodka/crestic/internal/entity"
	"github.com/alexander-kolodka/crestic/internal/healthchecks"
)

type HealthChecks interface {
	Start(ctx context.Context, rid string, j *healthchecks.JobsList) error
	Success(ctx context.Context, rid string, r *entity.JobResults) error
	Fail(ctx context.Context, rid string, r *entity.JobResults) error
}
