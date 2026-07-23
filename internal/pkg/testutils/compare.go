package testutils

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

func Equal[T any](t *testing.T, expected, actual T) {
	t.Helper()

	if !cmp.Equal(expected, actual) {
		t.Error(cmp.Diff(expected, actual))
	}
}
