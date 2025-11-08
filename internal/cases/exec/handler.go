package exec

import (
	"context"

	"github.com/alexander-kolodka/crestic/internal/entity"
	"github.com/alexander-kolodka/crestic/internal/logger"
	"github.com/alexander-kolodka/crestic/internal/restic"
)

type Command struct {
	Repos []*entity.Repository
	Cmd   string
	Args  []string
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
		repoCtx = logger.FromContext(repoCtx).With().
			Str("cmd", cmd.Cmd).
			Strs("args", cmd.Args).
			Logger().WithContext(repoCtx)

		err := h.restic.Exec(repoCtx, repo, cmd.Cmd, cmd.Args)
		if err != nil {
			return err
		}
	}

	return nil
}
