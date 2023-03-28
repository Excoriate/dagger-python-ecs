package task

import (
	"dagger.io/dagger"
	"github.com/Excoriate/dagger-python-ecs/pkg/job"
	"github.com/Excoriate/dagger-python-ecs/pkg/pipeline"
)

type DockerBuildTask struct {
	Init *InitOptions
	Cfg  *Task
}

func (t *DockerBuildTask) GetClient() *dagger.Client {
	return t.Cfg.JobCfg.Client
}

func (t *DockerBuildTask) SetEnvVars(envVars []map[string]string,
	container *dagger.Container) (*dagger.Container, error) {
	return nil, nil
}

func (t *DockerBuildTask) GetContainer(fromImage string) (*dagger.Container,
	error) {
	if fromImage == "" {
		return t.Cfg.JobCfg.ContainerDefault, nil
	}

	return t.Cfg.JobCfg.Client.Container().From(fromImage), nil
}

func (t *DockerBuildTask) RunTasksDefault(dir string, tasks []string) (Output, error) {
	return Output{}, nil
}

func (t *DockerBuildTask) RunDefault(dir string) (Output, error) {
	return Output{}, nil
}

func NewTaskDockerBuild(p *pipeline.Config, job *job.Job, actions []string,
	init *InitOptions) DockerBuildTask {
	t := NewTask(p, job, actions, init)

	return DockerBuildTask{
		Init: init,
		Cfg:  &t,
	}
}
