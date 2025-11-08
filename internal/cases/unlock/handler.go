package unlock

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

		err := h.restic.Unlock(repoCtx, repo)
		if err != nil {
			return err
		}
	}

	return nil
}
