package task

import (
	"dagger.io/dagger"
	"github.com/Excoriate/dagger-python-ecs/pkg/job"
	"github.com/Excoriate/dagger-python-ecs/pkg/pipeline"
)

type Tasker interface {
	GetClient() *dagger.Client
	SetEnvVars(envVars []map[string]string, container *dagger.Container) (*dagger.Container, error)
	GetContainer(fromImage string) (*dagger.Container, error)
	RunTasksDefault(dir string, tasks []string) (Output, error)
	RunDefault(dir string) (Output, error)
}

type Runner struct {
	Init *InitOptions
	Cfg  *Task
}

type Task struct {
	// Identifiers.
	Id    string
	Name  string
	Stack string

	// Configuration
	PipelineCfg *pipeline.Config
	JobCfg      *job.Job

	// Specific attributes
	EnvVarsInheritFromJob map[string]string
	Dirs                  Dirs
	ContainerImageDefault string
	ContainerNameDefault  string
	ContainerDefault      *dagger.Container

	PreReqs PreRequisites
	Actions Actions

	// Output
	Result Output
}

type Dirs struct {
	RootDir         string
	WorkDir         string
	MountDir        string
	TargetDir       string
	RootDirDagger   *dagger.Directory
	WorkDirDagger   *dagger.Directory
	MountDirDagger  *dagger.Directory
	TargetDirDagger *dagger.Directory
}

type Output struct {
	Files        []*dagger.File
	Directories  []*dagger.Directory
	ExitCode     int
	DaggerOutput interface{}
	IsError      bool
}

type Actions struct {
	CustomCommands  []string
	DefaultCommands []string
}

type PreRequisites struct {
	Files       []string
	Directories []string
}
