package task

import (
	"dagger.io/dagger"
	"fmt"
	"github.com/Excoriate/dagger-python-ecs/internal/common"
	"github.com/Excoriate/dagger-python-ecs/pkg/pipeline"
)

func (t Task) Init() error {
	t.Runner.UXMessage.ShowInfo("Initialising task", fmt.Sprintf("Task name: %s", t.Name))
	return nil
}

func (t Task) Configure(opt Config) error {
	return nil
}

func (t Task) Execute(container *dagger.Container) error {
	return nil
}

func (t Task) BuildContainer(image string) (*dagger.Container, error) {
	return nil, nil
}

func (t Task) Run() (Output, error) {
	t.Runner.UXMessage.ShowInfo("Running task", fmt.Sprintf("Task name: %s", t.Name))
	return Output{
		Files:       []*dagger.File{},
		Directories: []*dagger.Directory{},
		ExitCode:    0,
	}, nil
}

func NewTaskDockerBuild(p *pipeline.Runner, cmds []string) Task {
	return Task{
		Id:           common.GetUUID(),
		Name:         "Docker Build",
		Runner:       p,
		Config:       Config{Commands: cmds}, // The other options are set later in the process.
		DefaultImage: "docker:latest",
	}
}
