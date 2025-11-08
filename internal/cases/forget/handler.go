package forget

import (
	"context"

	"github.com/alexander-kolodka/crestic/internal/entity"
	"github.com/alexander-kolodka/crestic/internal/logger"
	"github.com/alexander-kolodka/crestic/internal/restic"
)

type Command struct {
	Repos  []*entity.Repository
	DryRun bool
	Prune  bool
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
	ctx = logger.FromContext(ctx).With().
		Bool("dry_run", cmd.DryRun).
		Bool("prune", cmd.Prune).
		Logger().WithContext(ctx)

	for _, repo := range cmd.Repos {
		repoCtx := logger.WithRepoFields(ctx, repo)
		err := h.forget(repoCtx, cmd, repo)
		if err != nil {
			return err
		}
	}

	return nil
}

func (h *Handler) forget(ctx context.Context, cmd *Command, r *entity.Repository) error {
	if cmd.DryRun {
		ctx = restic.WithDryRun(ctx)
	}

	r.ForgetOptions["prune"] = cmd.Prune
	return h.restic.Forget(ctx, r)
}
