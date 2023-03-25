package task

import (
	"dagger.io/dagger"
	"github.com/Excoriate/dagger-python-ecs/pkg/pipeline"
)

type TaskRunner interface {
	Run() (Output, error)
	Init() error
	Configure(opt Config) error
	BuildContainer(image string) (*dagger.Container, error)
	Execute(container *dagger.Container) error
}

type Task struct {
	// Identifiers.
	Id   string
	Name string

	// Pipeline client.
	Runner *pipeline.Runner
	Client *dagger.Client

	// Specific attributes
	DefaultImage string
	Config       Config

	// Output
	Result Output
}

type Config struct {
	EnvVarsScanned map[string]string
	EnvVarsPassed  map[string]string
	AWSEnvVars     map[string]string
	Commands       []string
}

type Output struct {
	Files       []*dagger.File
	Directories []*dagger.Directory
	ExitCode    int
	IsError     bool
}
