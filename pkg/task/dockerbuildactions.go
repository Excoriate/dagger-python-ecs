package task

import (
	"context"
	"github.com/Excoriate/dagger-python-ecs/internal/common"
)

type DockerBuildAction struct {
	Task   CoreTasker
	prefix string // How the UX messages should be prefixed
	Id     string // The ID of the task
	Name   string // The name of the task
	Ctx    context.Context
}

type DockerBuildActions interface {
	BuildFromDockerFile(dockerFile string) (Output, error)
}

func (a *DockerBuildAction) BuildFromDockerFile(dockerFile string) (Output, error) {
	ctx := a.Task.GetJob().Ctx

	container := a.Task.GetJobContainerDefault()
	client := a.Task.GetClient()
	targetDir := a.Task.GetCoreTask().Dirs.TargetDir
	preRequiredFiles := []string{"Dockerfile"}

	mountedContainer, err := a.Task.MountDir(targetDir, client, container, preRequiredFiles, ctx)
	if err != nil {
		return Output{}, err
	}

	_, err = mountedContainer.WithExec([]string{"ls", "-ltrh"}).
		ExitCode(ctx)

	return Output{}, nil
}

func NewDockerBuildAction(task CoreTasker) DockerBuildActions {
	return &DockerBuildAction{
		Task:   task,
		prefix: "DOCKER-ACTION:BUILD",
		Id:     common.GetUUID(),
		Name:   "Build Docker Image",
	}
}
