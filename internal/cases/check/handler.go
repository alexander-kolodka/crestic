package check

import (
	"context"

	"github.com/alexander-kolodka/crestic/internal/entity"
	"github.com/alexander-kolodka/crestic/internal/logger"
	"github.com/alexander-kolodka/crestic/internal/restic"
)

type Command struct {
	Repos []*entity.Repository
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
	for _, repo := range cmd.Repos {
		repoCtx := logger.WithRepoFields(ctx, repo)
		err := h.checkRepo(repoCtx, repo)
		if err != nil {
			return err
		}
	}

	return nil
}

func (h *Handler) checkRepo(ctx context.Context, r *entity.Repository) error {
	isRepoInitialized, err := h.restic.IsRepoInitialized(ctx, r)
	if err != nil {
		return err
	}

	if !isRepoInitialized {
		err = h.restic.Init(ctx, r)
		if err != nil {
			return err
		}
	}

	return h.restic.Check(ctx, r)
}
