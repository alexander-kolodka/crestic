package healthchecks

import (
	"context"
)

// Dummy is a no-op Service used when healthchecks are disabled.
type Dummy struct{}

func (s *Dummy) Start(_ context.Context, _, _ string) error {
	return nil
}

func (s *Dummy) Success(_ context.Context, _, _ string) error {
	return nil
}

func (s *Dummy) Fail(_ context.Context, _, _ string) error {
	return nil
}
