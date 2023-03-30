package task

import (
	"dagger.io/dagger"
	"github.com/Excoriate/dagger-python-ecs/internal/daggerio"
	"github.com/Excoriate/dagger-python-ecs/internal/filesystem"
	"github.com/Excoriate/dagger-python-ecs/pkg/job"
	"github.com/Excoriate/dagger-python-ecs/pkg/pipeline"
)

const taskNamePrefix = "TASK-DOCKER"

type DockerBuildTask struct {
	Init    *InitOptions
	Cfg     *Task
	Actions []string
}

func (t *DockerBuildTask) MountDir(dir string, container *dagger.Container) (*dagger.Container, error) {
	return nil, nil
}

func (t *DockerBuildTask) GetClient() *dagger.Client {
	return t.Cfg.JobCfg.Client
}

func (t *DockerBuildTask) GetPipeline() *pipeline.Config {
	return t.Cfg.PipelineCfg
}

func (t *DockerBuildTask) GetJob() *job.Job {
	return t.Cfg.JobCfg
}

func (t *DockerBuildTask) GetJobEnvVars() map[string]string {
	return t.Cfg.EnvVarsInheritFromJob
}

func (t *DockerBuildTask) SetEnvVars(envVars []map[string]string,
	container *dagger.Container) (*dagger.Container, error) {
	ux := t.Cfg.PipelineCfg.UXMessage

	if len(envVars) == 0 {
		ux.ShowInfo(taskNamePrefix, "There is no environment variables to be set in the container")
		return container, nil
	}

	var envVarsMerged map[string]string

	for _, envVar := range envVars {
		envVarsMerged = filesystem.MergeEnvVars(envVarsMerged, envVar)
	}

	return daggerio.SetEnvVarsInContainer(container, envVarsMerged)
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
	c := t.GetClient()
	defer c.Close()

	return Output{}, nil
}

func NewTaskDockerBuild(coreTask *Task, actions []string,
	init *InitOptions) CoreTasker {

	return &DockerBuildTask{
		Init:    init,
		Cfg:     coreTask,
		Actions: actions,
	}
}
