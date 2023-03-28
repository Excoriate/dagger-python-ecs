package cli

import (
	"github.com/Excoriate/dagger-python-ecs/pkg/job"
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
		p, err := pipeline.New(GlobalWorkingDirectory, GlobalMountDir, GlobalTargetDir,
			GlobalTaskName,
			GlobalScanEnvVarKeys,
			GlobalEnvKeyValuePairsToSet, GlobalScanAWSKeys, GlobalScanTFVars, GlobalSetWorkDirDaggerOnInit)

		if err != nil {
			os.Exit(1)
		}

		// 2. Initialising the job.
		j, err := job.NewJob(p, job.InitOptions{
			Name:  GlobalTaskName,
			Stack: "DOCKER",

			// Pipeline reference.
			PipelineCfg: p,

			// Critical directories to be resolved.
			WorkDir:   p.PipelineOpts.WorkDir,
			TargetDir: p.PipelineOpts.TargetDir,
			MountDir:  p.PipelineOpts.MountDir,

			// Environmental configuration
			ScanAWSEnvVars:       GlobalScanAWSKeys,
			ScanTerraformEnvVars: GlobalScanTFVars,
			EnvVarsToSet:         GlobalEnvKeyValuePairsToSet,
			EnvVarsToScan:        GlobalScanEnvVarKeys,
		})

		if err != nil {
			os.Exit(1)
		}

		// 3. Run The (Docker) task
		taskErr := task.RunTaskDocker(task.InitOptions{
			Task:           GlobalTaskName,
			Stack:          "DOCKER",
			PipelineCfg:    p,
			JobCfg:         j,
			WorkDir:        p.PipelineOpts.WorkDir,
			MountDir:       p.PipelineOpts.MountDir,
			TargetDir:      p.PipelineOpts.TargetDir,
			ActionCommands: GlobalCustomCommands,
		})

		if taskErr != nil {
			os.Exit(1)
		}
	},
}

func AddCIArguments() {
	// Add the flags to the root command.
}

func init() {
	AddCIArguments()
}
