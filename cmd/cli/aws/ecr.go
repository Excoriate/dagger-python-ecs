package aws

import (
	"github.com/Excoriate/dagger-python-ecs/internal/tui"
	"github.com/Excoriate/dagger-python-ecs/pkg/config"
	"github.com/Excoriate/dagger-python-ecs/pkg/job"
	"github.com/Excoriate/dagger-python-ecs/pkg/pipeline"
	"github.com/Excoriate/dagger-python-ecs/pkg/task"
	"github.com/spf13/cobra"
	"os"
)

var ECRCmd = &cobra.Command{
	Version: "v0.0.1",
	Use:     "ecr",
	Long: `The 'ecr' command automates and implement actions on top of AWS Elastic Container
Registry`,
	Example: `
  # Push an image into ECR:
  stiletto aws ecr --task=push`,
	Run: func(cmd *cobra.Command, args []string) {
		// 1. Instantiate the pipeline runner, which will be used to run the tasks.
		ux := tui.TUITitle{}
		config.ShowCLITitle()

		cliGlobalArgs := config.GetCLIGlobalArgs()

		p, err := pipeline.New(cliGlobalArgs.WorkingDir, cliGlobalArgs.MountDir,
			cliGlobalArgs.TargetDir, cliGlobalArgs.TaskName,
			cliGlobalArgs.ScanEnvVarKeys,
			cliGlobalArgs.EnvKeyValuePairsToSetString, cliGlobalArgs.ScanAWSKeys,
			cliGlobalArgs.ScanTerraformVars, cliGlobalArgs.InitDaggerWithWorkDirByDefault)

		if err != nil {
			os.Exit(1)
		}

		ux.ShowSubTitle("JOB:", "AWS-ECR")
		ux.ShowInitDetails("ECR", cliGlobalArgs.TaskName, p.PipelineOpts.WorkDirPath,
			p.PipelineOpts.TargetDirPath, p.PipelineOpts.MountDirPath)
		// 2. Initialising the job.
		j, err := job.NewJob(p, job.InitOptions{
			Name:  cliGlobalArgs.TaskName,
			Stack: "AWS",

			// Pipeline reference.
			PipelineCfg: p,

			// Critical directories to be resolved.
			WorkDir:   p.PipelineOpts.WorkDir,
			TargetDir: p.PipelineOpts.TargetDir,
			MountDir:  p.PipelineOpts.MountDir,

			// Environmental configuration
			ScanAWSEnvVars:       cliGlobalArgs.ScanAWSKeys,
			ScanTerraformEnvVars: cliGlobalArgs.ScanTerraformVars,
			EnvVarsToSet:         cliGlobalArgs.EnvKeyValuePairsToSetString,
			EnvVarsToScan:        cliGlobalArgs.ScanEnvVarKeys,
		})

		if err != nil {
			os.Exit(1)
		}

		// 3. Run The (Docker) task
		ux.ShowSubTitle("TASK:", cliGlobalArgs.TaskName)
		ux.ShowTaskDetails("DOCKER", cliGlobalArgs.TaskName, j.WorkDirPath, j.TargetDirPath, j.MountDirPath)
		taskErr := task.RunTaskDocker(task.InitOptions{
			//Task:           GlobalTaskName,
			Task:           cliGlobalArgs.TaskName,
			Stack:          "AWS",
			PipelineCfg:    p,
			JobCfg:         j,
			WorkDir:        p.PipelineOpts.WorkDir,
			MountDir:       p.PipelineOpts.MountDir,
			TargetDir:      p.PipelineOpts.TargetDir,
			ActionCommands: cliGlobalArgs.CustomCommands,
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
