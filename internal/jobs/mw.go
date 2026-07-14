package jobs

import (
	"context"
	"slices"

	"github.com/alexander-kolodka/crestic/internal/entity"
)

type (
	mw func(do) do
	do func(ctx context.Context, j entity.Job) error
)

// chain wraps the base function with all provided middlewares.
// Middlewares are applied in the same order as given:
// the first middleware becomes the outermost wrapper,
// the last middleware is the closest to the base function.
//
// So if you call chain(base, m1, m2, m3), the execution flow will be:
// m1 → m2 → m3 → base.
func chain(base func(ctx context.Context, job entity.Job) error, middlewares ...mw) do {
	result := base
	for _, v := range slices.Backward(middlewares) {
		result = v(result)
	}
	return result
}
