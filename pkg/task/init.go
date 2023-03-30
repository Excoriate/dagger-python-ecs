package task

import (
	"github.com/Excoriate/dagger-python-ecs/pkg/job"
	"github.com/Excoriate/dagger-python-ecs/pkg/pipeline"
)

type InitOptions struct {
	Task  string
	Stack string

	PipelineCfg *pipeline.Config
	JobCfg      *job.Job

	// Directories that the task will use.
	WorkDir   string
	MountDir  string
	TargetDir string

	// Behaviour
	ActionCommands []string
}
