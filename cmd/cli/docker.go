package cli

import (
	"github.com/Excoriate/dagger-python-ecs/pkg/pipeline"
	"github.com/Excoriate/dagger-python-ecs/pkg/task"
	"github.com/spf13/cobra"
	"os"
)

var DockerCMD = &cobra.Command{
	Version: "v0.0.1",
	Use:     "docker",
	Long: `The 'docker' command automates various Continuous Integration (
DockerCMD) tasks such as building lint, build, running tests, linting, and installing dependencies.
You can specify the tasks you want to perform using the provided flags.`,
	Example: `
  # DockerCMD pipeline DockerCMD with build, test, and lint tasks
  pipeline docker --build
  # DockerCMD pipeline DockerCMD with install dependencies task
  pipeline docker --install`,
	Run: func(cmd *cobra.Command, args []string) {
		// 1. Instantiate the pipeline runner, which will be used to run the tasks.
		p, err := pipeline.New(GlobalWorkDir, GlobalTargetDir, GlobalTaskName, GlobalScanEnvVarKeys,
			GlobalEnvKeyValuePairsToSet, GlobalScanEnvVarKeys)

		if err != nil {
			os.Exit(1)
		}

		// 2. Run Docker job.
		taskErr := task.RunTaskDocker(task.Job{
			Stack:          "DOCKER",
			Task:           GlobalTaskName,
			CustomCommands: GlobalCustomCommands,
			Pipeline:       p,
		})

		if taskErr != nil {
			os.Exit(1)
		}
	},
}

func AddCIArguments() {
	DockerCMD.Flags().StringVarP(&GlobalTaskName,
		"task",
		"t", "",
		"Name of the task to run. E.g.: build, test, lint, install")

	DockerCMD.Flags().StringVarP(&GlobalWorkDir,
		"workdir",
		"w", "",
		"Working directory where the pipeline will be executed")

	DockerCMD.Flags().StringVarP(&GlobalTargetDir,
		"target-dir",
		"r", "",
		"Target directory where the pipeline will be executed. If it's not set, it'll use the same value as the working directory.")

	DockerCMD.Flags().StringSliceVarP(&GlobalScanEnvVarKeys,
		"scan-env",
		"s", []string{},
		"List of environment variable keys that are already exported, that'll be scanned and set.")

	DockerCMD.Flags().StringToStringVarP(&GlobalEnvKeyValuePairsToSet,
		"set-env",
		"e", map[string]string{},
		"List of environment variable key-value pairs to set.")

	DockerCMD.Flags().StringSliceVarP(&GlobalCustomCommands,
		"commands",
		"c", []string{},
		"List of custom commands to run.")

	DockerCMD.Flags().BoolVarP(&GlobalScanAWSKeys,
		"scan-aws-keys",
		"a", false,
		"Scan AWS keys and set them as environment variables.")

	_ = DockerCMD.MarkFlagRequired("task")
	_ = DockerCMD.MarkFlagRequired("workdir")
}

func init() {
	AddCIArguments()
}
