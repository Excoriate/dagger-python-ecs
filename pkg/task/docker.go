package task

import (
	"fmt"
	"github.com/Excoriate/dagger-python-ecs/internal/common"
)

type DockerAction struct {
	// Set of task types that this task can run.
	Task *DockerBuildTask
}

type DockerActionExecutor interface {
	BuildDockerFile(dockerFile string) (Output, error)
}

func (t *DockerAction) BuildDockerFile(dockerFile string) (Output, error) {
	buildTask := t.Task.Cfg
	ux := buildTask.PipelineCfg.UXMessage

	ux.ShowSuccess("DOCKER-ACTION:BUILD",
		fmt.Sprintf("Building Docker image from Dockerfile: %s",
			dockerFile))
	return Output{}, nil
}

// RunTaskDocker is the entry point for all Docker tasks.
func RunTaskDocker(opt InitOptions) error {
	taskSelector := common.NormaliseStringUpper(opt.Task)

	p := opt.PipelineCfg
	j := opt.JobCfg

	actionCMDs := opt.ActionCommands

	switch taskSelector {
	case "BUILD":
		// New (core) instance of a task
		c := NewTask(p, j, actionCMDs, &opt)

		// New specific instance of a task (E.g.: Docker, AWS, etc.)
		t := NewTaskDockerBuild(c, actionCMDs, &opt)

		// New action to execute (mapped to the --task passed from the command line)
		a := NewDockerBuildAction(t)

		// Run the action
		_, err := a.BuildFromDockerFile("Dockerfile")
		if err != nil {
			return err
		}
	}
	return nil
}
