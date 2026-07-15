package harness

// PipelineBuilder configures a pipeline in the sandbox config.
type PipelineBuilder struct {
	sandbox  *Sandbox
	pipeline *pipelineEntry
}

// Cron sets an optional cron schedule for the pipeline.
func (b *PipelineBuilder) Cron(expr string) *PipelineBuilder {
	b.pipeline.Cron = expr
	return b
}

// Hooks sets pipeline-level lifecycle hooks.
func (b *PipelineBuilder) Hooks(h Hooks) *PipelineBuilder {
	b.pipeline.Hooks = h
	return b
}

// Backup adds a backup job to the pipeline.
func (b *PipelineBuilder) Backup(name string, from []string, to *RepoDir) *PipelineBuilder {
	b.pipeline.Jobs = append(b.pipeline.Jobs, jobEntry{
		Type:    "backup",
		Name:    name,
		From:    from,
		To:      to.Name,
		Options: map[string]any{"skip-if-unchanged": true},
	})
	return b
}

// Copy adds a copy job to the pipeline.
func (b *PipelineBuilder) Copy(name string, from, to *RepoDir) *PipelineBuilder {
	b.pipeline.Jobs = append(b.pipeline.Jobs, jobEntry{
		Type:     "copy",
		Name:     name,
		FromRepo: from.Name,
		To:       to.Name,
	})
	return b
}

// JobHooks sets lifecycle hooks on the last added job.
func (b *PipelineBuilder) JobHooks(h Hooks) *PipelineBuilder {
	if len(b.pipeline.Jobs) == 0 {
		b.sandbox.t.Fatal("JobHooks called with no jobs on pipeline")
	}
	b.pipeline.Jobs[len(b.pipeline.Jobs)-1].Hooks = h
	return b
}
