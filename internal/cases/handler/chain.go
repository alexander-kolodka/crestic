package handler

// Middleware represents a decorator function that wraps a Handler.
type Middleware[CMD any] func(Handler[CMD]) Handler[CMD]

// Chain wraps the base function with all provided middlewares.
// Middlewares are applied in the same order as given:
// the first middleware becomes the outermost wrapper,
// the last middleware is the closest to the base function.
//
// So if you call Chain(base, m1, m2, m3), the execution flow will be:
// m1 → m2 → m3 → base.
func Chain[CMD any](base Handler[CMD], middlewares ...Middleware[CMD]) Handler[CMD] {
	result := base
	for i := len(middlewares) - 1; i >= 0; i-- {
		result = middlewares[i](result)
	}
	return result
}
