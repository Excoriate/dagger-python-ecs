package task

import (
	"dagger.io/dagger"
	"github.com/Excoriate/dagger-python-ecs/pkg/pipeline"
)

type Executioner interface {
	Run() (Output, error)
	//Init() (*dagger.Client, error)
	Init() error
	Configure() error
	ConfigureContainer(cfg Config, container *dagger.Container) error
	BuildContainer(stack, image string) (*dagger.Container, error)
	Execute(container *dagger.Container) error
}

type Task struct {
	// Identifiers.
	Id    string
	Name  string
	Stack string

	// Pipeline client.
	Runner    *pipeline.Runner
	Client    *dagger.Client
	WorkDir   *dagger.Directory
	TargetDIr *dagger.Directory

	// Specific attributes
	Config                Config
	ContainerImage        string
	ContainerImageDefault string
	Container             *dagger.Container

	// Output
	Result Output
}

type Config struct {
	EnvVarsScanned map[string]string
	EnvVarsPassed  map[string]string
	AWSEnvVars     map[string]string
	Commands       []string
	Workdir        string
	TargetDir      string
	ContainerImage string
}

type Output struct {
	Files       []*dagger.File
	Directories []*dagger.Directory
	ExitCode    int
	IsError     bool
}
