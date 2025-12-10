package entity

import (
	"fmt"
	"strings"
	"time"

	"github.com/samber/lo"
)

type JobResults struct {
	SuccessJobs []SuccessJob `json:"successJobs,omitempty"`
	FailedJobs  []FailedJob  `json:"failedJobs,omitempty"`
}

type SuccessJob struct {
	Name    string `json:"name"`
	Elapsed string `json:"elapsed"`
}

type FailedJob struct {
	Name    string `json:"name"`
	Elapsed string `json:"elapsed"`
	Error   string `json:"error"`
}

func NewJobResults() *JobResults {
	return &JobResults{
		FailedJobs:  []FailedJob{},
		SuccessJobs: []SuccessJob{},
	}
}

func (r *JobResults) ErrorMsg() string {
	if len(r.FailedJobs) == 0 {
		return ""
	}

	errs := lo.Map(r.FailedJobs, func(job FailedJob, _ int) string {
		return fmt.Sprintf("%s elapsed: %s failed: %s", job.Name, job.Elapsed, job.Error)
	})

	return strings.Join(errs, ";\n")
}

func (r *JobResults) HasErrors() bool {
	return len(r.FailedJobs) > 0
}

func (r *JobResults) Add(jobName string, elapsed time.Duration, err error) {
	if err != nil {
		r.failed(jobName, elapsed, err)
		return
	}

	r.success(jobName, elapsed)
}

func (r *JobResults) success(job string, elapsed time.Duration) {
	r.SuccessJobs = append(r.SuccessJobs, SuccessJob{
		Name:    job,
		Elapsed: elapsed.String(),
	})
}

func (r *JobResults) failed(job string, elapsed time.Duration, err error) {
	r.FailedJobs = append(r.FailedJobs, FailedJob{
		Name:    job,
		Elapsed: elapsed.String(),
		Error:   err.Error(),
	})
}
