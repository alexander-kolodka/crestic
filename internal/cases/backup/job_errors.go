package backup

import (
	"fmt"
	"strings"

	"github.com/samber/lo"
)

type JobListError struct {
	errors map[string]error
}

func newJobErrors() *JobListError {
	return &JobListError{errors: map[string]error{}}
}

func (e *JobListError) Error() string {
	errs := lo.MapToSlice(e.errors, func(job string, err error) string {
		return fmt.Sprintf("%s failed: %s", job, err)
	})

	return strings.Join(errs, ";\n")
}

func (e *JobListError) HasErrors() bool {
	return len(e.errors) > 0
}

func (e *JobListError) Add(job string, err error) {
	e.errors[job] = err
}
