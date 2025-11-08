package panix

import (
	"fmt"
	"runtime/debug"
)

// PanicError wraps panic information as a standard error.
// Contains the panic value and full stacktrace for diagnostics.
type PanicError struct {
	Value      any
	StackTrace string
}

// NewPanicError creates a PanicError with captured current stacktrace.
func NewPanicError(p any) *PanicError {
	return &PanicError{Value: p, StackTrace: string(debug.Stack())}
}

func (e *PanicError) Error() string {
	return fmt.Sprintf("panic: %v", e.Value)
}
