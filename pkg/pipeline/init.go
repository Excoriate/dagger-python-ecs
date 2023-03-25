package pipeline

import (
	"dagger.io/dagger"
	"fmt"
	"github.com/Excoriate/dagger-python-ecs/internal/common"
	"github.com/Excoriate/dagger-python-ecs/internal/errors"
	"github.com/Excoriate/dagger-python-ecs/internal/filesystem"
	"github.com/Excoriate/dagger-python-ecs/internal/logger"
	"github.com/Excoriate/dagger-python-ecs/internal/tui"
	"github.com/Excoriate/dagger-python-ecs/pkg/config"
)

func isWorkDirValid(workDir string) error {
	workDirNormalised := common.NormaliseNoSpaces(workDir)
	if _, err := filesystem.PathExist(workDirNormalised); err != nil {
		return errors.NewPipelineConfigurationError("Pipeline cant initialise", err)
	}

	if err := filesystem.PathIsADirectory(workDirNormalised); err != nil {
		return errors.NewPipelineConfigurationError("Pipeline cant initialise", err)
	}

	return nil
}

func isTargetDirValid(targetDir string) error {
	if targetDir == "" {
		return nil
	}

	targetDirNormalised := common.NormaliseNoSpaces(targetDir)
	if _, err := filesystem.PathExist(targetDirNormalised); err != nil {
		return errors.NewPipelineConfigurationError("Pipeline cant initialise", err)
	}

	if err := filesystem.PathIsADirectory(targetDirNormalised); err != nil {
		return errors.NewPipelineConfigurationError("Pipeline cant initialise", err)
	}

	return nil
}

func isTaskNameValid(taskName string) error {
	normalisedTask := common.NormaliseStringUpper(taskName)
	if normalisedTask != "LINT" && normalisedTask != "TEST" && normalisedTask != "BUILD" {
		return errors.NewPipelineConfigurationError("Pipeline cant initialise", fmt.Errorf("task name %s is not valid", taskName))
	}

	return nil
}

func areEnvKeysToScanExported(envKeysToScan []string) error {
	if len(envKeysToScan) == 0 {
		return nil
	}

	err := filesystem.CheckEnvVars(envKeysToScan)
	if err != nil {
		return errors.NewPipelineConfigurationError("Pipeline cant initialise", err)
	}
	return nil
}

func isEnvVarsMapToSetValid(envVarsMapToSet map[string]string) error {
	if len(envVarsMapToSet) == 0 {
		return nil
	}

	for key := range envVarsMapToSet {
		if err := filesystem.IsEnvVarSetOrExported(key); err != nil {
			return errors.NewPipelineConfigurationError("Pipeline cant initialise", fmt.Errorf("env var %s is not set or exported", key))
		}
	}

	return nil
}

func IsAWSKeysExported() error {
	// Look for AWS keys exported, if not, fail with an error
	if err := filesystem.IsEnvVarSetOrExported("AWS_ACCESS_KEY_ID"); err != nil {
		return errors.NewPipelineConfigurationError("Pipeline cant initialise", fmt.Errorf("env var AWS_ACCESS_KEY_ID is not set or exported"))
	}

	if err := filesystem.IsEnvVarSetOrExported("AWS_SECRET_ACCESS_KEY"); err != nil {
		return errors.NewPipelineConfigurationError("Pipeline cant initialise", fmt.Errorf("env var AWS_SECRET_ACCESS_KEY is not set or exported"))
	}

	return nil
}

func CheckPreConditions(args config.PipelineOptions) error {
	ux := tui.TUIMessage{}
	if err := isWorkDirValid(args.WorkDir); err != nil {
		ux.ShowError("VALIDATION", "Preconditions failed", err)
		return err
	}

	if err := isTargetDirValid(args.TargetDir); err != nil {
		ux.ShowError("VALIDATION", "Preconditions failed", err)
		return err
	}

	if err := isTaskNameValid(args.TaskName); err != nil {
		ux.ShowError("VALIDATION", "Preconditions failed", err)
		return err
	}

	if err := areEnvKeysToScanExported(args.EnvVarsToScanAndSet); err != nil {
		ux.ShowError("VALIDATION", "Preconditions failed", err)
		return err
	}

	if err := isEnvVarsMapToSetValid(args.EnvKeyValuePairsToSet); err != nil {
		ux.ShowError("VALIDATION", "Preconditions failed", err)
		return err
	}

	if err := IsAWSKeysExported(); err != nil {
		ux.ShowError("VALIDATION", "Preconditions failed", err)
		return err
	}

	return nil
}

func New(workdir, targetDir, taskName string, envVarKeysToScan []string,
	envVarsMapToSet map[string]string, envVarsAWSKeysToScan []string) (*Runner, error) {
	args := config.PipelineOptions{
		// Specific directories to work with passed.
		WorkDir:   workdir,
		TargetDir: targetDir,

		// Task identifier, that'll be used to determine what to do.
		TaskName: taskName,
		// Specific environmental options passed.
		EnvVarsToScanAndSet:   envVarKeysToScan,
		EnvKeyValuePairsToSet: envVarsMapToSet,
		EnvVarsAWSKeysToScan:  envVarsAWSKeysToScan,
	}

	if err := CheckPreConditions(args); err != nil {
		return nil, err
	}

	logPrinter := logger.NewLogger()
	logPrinter.InitLogger()

	dirs := config.GetDefaultDirs()

	platformToArch := map[dagger.Platform]string{
		"linux/amd64": "amd64",
		"linux/arm64": "arm64",
	}

	return &Runner{
		Logger:       logPrinter,
		Dirs:         *dirs,
		UXDisplay:    tui.NewTitle(),
		Platforms:    platformToArch,
		UXMessage:    tui.NewTUIMessage(),
		PipelineOpts: &args,
	}, nil
}
