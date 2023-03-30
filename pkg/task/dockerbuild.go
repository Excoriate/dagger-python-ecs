package task

import (
	"dagger.io/dagger"
	"fmt"
	"github.com/Excoriate/dagger-python-ecs/internal/daggerio"
	"github.com/Excoriate/dagger-python-ecs/internal/filesystem"
	"github.com/Excoriate/dagger-python-ecs/internal/tui"
	"github.com/Excoriate/dagger-python-ecs/pkg/job"
	"github.com/Excoriate/dagger-python-ecs/pkg/pipeline"
)

const taskNamePrefix = "DOCKER-BUILD"

type DockerBuildTask struct {
	Init    *InitOptions
	Cfg     *Task
	Actions []string
}

func (t *DockerBuildTask) MountDir(dir string, container *dagger.Container) (*dagger.Container, error) {
	ux := tui.NewTUIMessage()

	if dir == "" {
		ux.ShowWarning(taskNamePrefix, "An empty directory was passed to MountDir, "+
			"hence the default working directory will be used resolved from the '.' value")

		dir = "."
	}

	dirDagger, err := daggerio.GetDaggerDir(t.GetClient(), dir)

	if err != nil {
		ux.ShowError(taskNamePrefix,
			fmt.Sprintf("Failed to mount directory, failed "+
				"to build a dagger directory from the directory"+
				" passed in: %s",
				dir), err)

		return nil, err
	}

	containerMounted, err := daggerio.MountDir(container, dirDagger, "")

	if err != nil {
		ux.ShowError(taskNamePrefix,
			fmt.Sprintf("Failed to mount directory %s", dir), err)

		return nil, err
	}

	return containerMounted, nil
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

func (t *DockerBuildTask) GetCoreTask() *Task {
	return t.Cfg
}

func (t *DockerBuildTask) GetJobContainerImage() string {
	return t.Cfg.JobCfg.ContainerImageURL
}

func (t *DockerBuildTask) GetJobContainerDefault() *dagger.Container {
	return t.Cfg.JobCfg.ContainerDefault
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

func NewTaskDockerBuild(coreTask *Task, actions []string,
	init *InitOptions) CoreTasker {

	return &DockerBuildTask{
		Init:    init,
		Cfg:     coreTask,
		Actions: actions,
	}
}
