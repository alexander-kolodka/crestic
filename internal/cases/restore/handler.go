package restore

import (
	"context"

	"github.com/alexander-kolodka/crestic/internal/entity"
	"github.com/alexander-kolodka/crestic/internal/logger"
	"github.com/alexander-kolodka/crestic/internal/restic"
)

type Command struct {
	Repo     *entity.Repository
	Target   string
	Snapshot string
}

type Handler struct {
	restic *restic.Service
}

func NewHandler(restic *restic.Service) *Handler {
	return &Handler{
		restic: restic,
	}
}

func (h *Handler) Handle(ctx context.Context, cmd *Command) error {
	ctx = logger.WithRepoFields(ctx, cmd.Repo)
	ctx = logger.FromContext(ctx).With().
		Str("snapshot", cmd.Snapshot).
		Str("target", cmd.Target).
		Logger().WithContext(ctx)

	return h.restic.Restore(ctx, cmd.Repo, cmd.Target, cmd.Snapshot)
}
