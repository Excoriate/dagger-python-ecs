package task

import (
	"fmt"
	"github.com/Excoriate/dagger-python-ecs/internal/common"
	"github.com/Excoriate/dagger-python-ecs/pkg/pipeline"
)

type Job struct {
	Task           string
	CustomCommands []string
	Pipeline       *pipeline.Runner
}

func RunTaskDocker(newJob Job) error {
	task := common.NormaliseStringUpper(newJob.Task)
	p := newJob.Pipeline
	cmds := newJob.CustomCommands

	switch task {
	case "BUILD":
		t := NewTaskDockerBuild(p, cmds)
		out, err := t.Run()
		if err != nil {
			return err
		}

		p.UXMessage.ShowInfo("", fmt.Sprintf("Task %s completed with exit code %s", t.Name,
			out.ExitCode))
	}
	return nil
}
