package task

import (
	"github.com/Excoriate/dagger-python-ecs/internal/common"
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

func RunTaskDocker(new InitOptions) error {
	taskName := common.NormaliseStringUpper(new.Task)
	p := new.PipelineCfg
	j := new.JobCfg
	actions := new.ActionCommands

	switch taskName {
	case "BUILD":
		t := NewTaskDockerBuild(p, j, actions, &new)
		_, err := t.RunDefault("")
		if err != nil {
			return err
		}
	}
	return nil
}
