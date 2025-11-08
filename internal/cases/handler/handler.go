package handler

import "context"

type Handler[CMD any] interface {
	Handle(ctx context.Context, cmd CMD) error
}

type handler[CMD any] struct {
	handle func(ctx context.Context, cmd CMD) error
}

func (h *handler[CMD]) Handle(ctx context.Context, cmd CMD) error {
	return h.handle(ctx, cmd)
}

func NewHandler[CMD any](handle func(ctx context.Context, cmd CMD) error) Handler[CMD] {
	return &handler[CMD]{
		handle: handle,
	}
}
