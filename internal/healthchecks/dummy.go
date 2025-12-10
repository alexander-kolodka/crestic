package healthchecks

import (
	"context"

	"github.com/alexander-kolodka/crestic/internal/entity"
)

// Dummy is a no-op implementation.
type Dummy struct{}

func (s *Dummy) Start(_ context.Context, _ string, _ *JobsList) error {
	return nil
}

func (s *Dummy) Success(_ context.Context, _ string, _ *entity.JobResults) error {
	return nil
}

func (s *Dummy) Fail(_ context.Context, _ string, _ *entity.JobResults) error {
	return nil
}
