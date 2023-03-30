package task

import (
	"context"
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

func (t *DockerBuildTask) MountDir(targetDir string, client *dagger.Client, container *dagger.
Container,
	filesPreRequisites []string, ctx context.Context) (*dagger.Container, error) {
	ux := tui.NewTUIMessage()

	if targetDir == "" {
		ux.ShowWarning(taskNamePrefix, "An empty directory was passed to be a Target directory ("+
			"also known as Execution path), "+
			"hence the default working directory will be used resolved from the '.' value")

		targetDir = "."
	}

	if targetDir != "." && len(filesPreRequisites) > 0 {
		ux.ShowInfo(taskNamePrefix, "The target directory is not the working directory, "+
			"therefore the files pre-requisites will be verified before mounting the directory")

		if err := daggerio.VerifyFileEntriesInMountedDir(client, targetDir,
			filesPreRequisites, ctx); err != nil {
			ux.ShowError(taskNamePrefix, "Failed to mount the directory", err)
			return nil, err
		}
	}

	workDirDagger, err := daggerio.GetDaggerDir(t.GetClient(), ".")

	if err != nil {
		ux.ShowError(taskNamePrefix,
			fmt.Sprintf("Failed to mount the working directory (with value '.'), failed "+
				"to build a dagger directory from the directory"), err)

		return nil, err
	}

	containerMounted, err := daggerio.MountDir(container, workDirDagger, targetDir)

	if err != nil {
		ux.ShowError(taskNamePrefix,
			fmt.Sprintf("Failed to mount directory %s", targetDir), err)

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

func (t *DockerBuildTask) ConvertDir(c *dagger.Client, dir string) (*dagger.Directory, error) {
	return daggerio.GetDaggerDir(c, dir)
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
