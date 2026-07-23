package mw

import (
	"context"
	"slices"
)

type (
	Func[T any]       func(ctx context.Context, t T) error
	Middleware[T any] func(Func[T]) Func[T]
)

// Chain wraps the base function with all provided middlewares.
// Middlewares are applied in the same order as given:
// the first middleware becomes the outermost wrapper,
// the last middleware is the closest to the base function.
//
// So if you call Chain(base, m1, m2, m3), the execution flow will be:
// m1 → m2 → m3 → base.
func Chain[T any](base Func[T], middlewares ...Middleware[T]) Func[T] {
	result := base
	for _, mw := range slices.Backward(middlewares) {
		result = mw(result)
	}
	return result
}
