package handler

import "github.com/alexander-kolodka/crestic/internal/pkg/mw"

// Chain applies middleware around a Handler.
func Chain[CMD any](h Handler[CMD], mws ...mw.Middleware[CMD]) Handler[CMD] {
	return NewHandler(mw.Chain(h.Handle, mws...))
}
