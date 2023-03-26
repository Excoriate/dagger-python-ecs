package task

import (
	"dagger.io/dagger"
	"fmt"
	"github.com/Excoriate/dagger-python-ecs/internal/common"
	"github.com/Excoriate/dagger-python-ecs/internal/daggerio"
	"github.com/Excoriate/dagger-python-ecs/internal/errors"
	"github.com/Excoriate/dagger-python-ecs/internal/filesystem"
	"github.com/Excoriate/dagger-python-ecs/pkg/pipeline"
)

//func (t Task) Init() (*dagger.Client, error) {
func (t *Task) Init() error {
	t.Runner.UXMessage.ShowInfo("DAGGER", fmt.Sprintf("Initialising dagger client: %s - id: %s",
		t.Name, t.Id))

	c, err := daggerio.NewDaggerClient(t.Config.Workdir, t.Runner.Ctx)

	if err != nil {
		msg := GetErrMsg(t, "Dagger initialisation failed", nil)
		return errors.NewDaggerEngineError(msg, err)
	}

	defer c.Close()

	t.Runner.UXMessage.ShowInfo("DAGGER", GetInfoMsg(t, "Dagger client initialised"))
	t.Client = c

	return nil
}

func (t *Task) Configure() error {
	// Scan AWS credentials from env vars
	if t.Runner.PipelineOpts.IsAWSKeysToScan {
		awsEnvVarsFetched, err := filesystem.ScanAWSCredentialsEnvVars()
		if err != nil {
			msg := GetErrMsg(t, "Failed to scan AWS credentials", nil)
			t.Runner.UXMessage.ShowError("", msg, nil)
			return errors.NewTaskConfigurationError(msg, err)
		}

		t.Config.AWSEnvVars = awsEnvVarsFetched
		t.Runner.UXMessage.ShowInfo("", GetInfoMsg(t, "AWS credentials successfully scanned"))
	}

	// Scan 'passed' by parameter set of environment variables keys.
	if len(t.Runner.PipelineOpts.EnvVarsToScanAndSet) > 0 {
		t.Runner.UXMessage.ShowInfo("", fmt.Sprintf("Scanning environment variables: %v", t.Runner.PipelineOpts.EnvVarsToScanAndSet))
		envVarsFetched, err := filesystem.FetchEnvVarsAsMap(t.Runner.PipelineOpts.EnvVarsToScanAndSet)
		if err != nil {
			msg := GetErrMsg(t, "Failed to scan environment variables", nil)
			t.Runner.UXMessage.ShowError("", msg, nil)
			return errors.NewTaskConfigurationError(msg, err)
		}

		t.Config.EnvVarsScanned = envVarsFetched
		t.Runner.UXMessage.ShowInfo("", GetInfoMsg(t, "Environment variables successfully scanned"))
	}

	envVarsToSet := t.Runner.PipelineOpts.EnvKeyValuePairsToSet
	if !common.MapIsNulOrEmpty(envVarsToSet) {
		if err := filesystem.AreEnvVarsConsistent(envVarsToSet); err != nil {
			msg := GetErrMsg(t, "Environment variables passed are not consistent", nil)
			t.Runner.UXMessage.ShowError("", msg, nil)
			return errors.NewTaskConfigurationError(msg, err)
		}
	}

	// (Dagger) working directory configuration
	workDirDagger, err := daggerio.GetWorkDir(t.Client, t.Config.Workdir)
	if err != nil {
		msg := GetErrMsg(t, "Failed to get and configure Dagger working directory", nil)
		t.Runner.UXMessage.ShowError("", msg, nil)
		return errors.NewTaskConfigurationError(msg, err)
	}

	t.Runner.UXMessage.ShowInfo("",
		GetInfoMsg(t, fmt.Sprintf("Dagger working directory successfully configured: %s", t.Config.Workdir)))
	t.WorkDir = workDirDagger

	return nil
}

func (t *Task) ConfigureContainer(cfg Config, container *dagger.Container) error {
	return nil
}

func (t *Task) Execute(container *dagger.Container) error {
	return nil
}

func (t *Task) BuildContainer(stack, image string) (*dagger.Container, error) {
	if image == "" {
		// If the image isn't passed, it'll resolve by STACK.
		if imageResolved, err := daggerio.GetContainerImagePerStack(stack, ""); err != nil {
			msg := GetErrMsg(t, "Failed to get container image", nil)
			t.Runner.UXMessage.ShowError("", msg, nil)
			return nil, errors.NewTaskConfigurationError(msg, err)
		} else {
			t.ContainerImage = imageResolved
		}
	} else {
		// If the image is passed, it'll fetch it by it.
		if customImageResolved, err := daggerio.GetContainerImageCustom(image, ""); err != nil {
			msg := GetErrMsg(t, "Failed to get container image (custom)", nil)
			t.Runner.UXMessage.ShowError("", msg, nil)
			return nil, errors.NewTaskConfigurationError(msg, err)
		} else {
			t.ContainerImage = fmt.Sprintf("%s:%s", customImageResolved.Image,
				customImageResolved.Version)
		}
	}

	t.Runner.UXMessage.ShowInfo("DAGGER", GetInfoMsg(t, fmt.Sprintf("Container image resolved: %s",
		t.ContainerImage)))

	container, err := daggerio.GetContainer(t.Client, t.ContainerImage)
	if err != nil {
		msg := GetErrMsg(t, fmt.Sprintf("Failed to initialise container with image %s", image),
			nil)
		t.Runner.UXMessage.ShowError("", msg, nil)
		return nil, errors.NewTaskConfigurationError(msg, err)
	}

	t.Container = container

	t.Runner.UXMessage.ShowInfo("DAGGER", GetInfoMsg(t, fmt.Sprintf("Container successfully initialised: %s",
		t.ContainerImage)))

	return container, nil
}

func (t *Task) Run() (Output, error) {
	t.Runner.UXMessage.ShowInfo("", fmt.Sprintf("Running new task name: %s - id: %s", t.Name, t.Id))

	// Initialise the dagger client.
	if err := t.Init(); err != nil {
		return Output{
			ExitCode: 1,
			IsError:  true,
		}, err
	}

	// Configure the task.
	if err := t.Configure(); err != nil {
		return Output{
			ExitCode: 1,
			IsError:  true,
		}, err
	}

	// Build the container. The ref. of the container is store inside the 't' object.
	_, err := t.BuildContainer(t.Stack, t.ContainerImageDefault)
	if err != nil {
		return Output{
			ExitCode: 1,
			IsError:  true,
		}, err
	}

	return Output{
		Files:       []*dagger.File{},
		Directories: []*dagger.Directory{},
		ExitCode:    0,
	}, nil
}

func NewTaskDockerBuild(p *pipeline.Runner, cmds []string) Task {
	targetDir := p.PipelineOpts.TargetDir
	workDir := p.PipelineOpts.WorkDir
	return Task{
		Id:     common.GetUUID(),
		Name:   "Docker Build",
		Stack:  "DOCKER",
		Runner: p,
		Config: Config{Commands: cmds, TargetDir: targetDir, Workdir: workDir},
	}
}
