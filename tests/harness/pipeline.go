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
