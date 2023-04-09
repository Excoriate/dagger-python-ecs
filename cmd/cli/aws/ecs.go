package aws

import (
	"github.com/Excoriate/dagger-python-ecs/internal/tui"
	"github.com/Excoriate/dagger-python-ecs/pkg/config"
	"github.com/Excoriate/dagger-python-ecs/pkg/job"
	"github.com/Excoriate/dagger-python-ecs/pkg/pipeline"
	"github.com/Excoriate/dagger-python-ecs/pkg/task"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"os"
)

var (
	ecsService     string
	ecsCluster     string
	taskDefinition string // Map to the task definition family name.
	imageURL       string
	releaseVersion string
)

var ECSCmd = &cobra.Command{
	Version: "v0.0.1",
	Use:     "ecs",
	Long: `The 'ecs' command automates and implement several Elastic Container Service actions,
E.g.: 'deploy'`,
	Example: `
  # Deploy a new version of a task running in a ECS service:
  stiletto aws ecs --task=deploy`,
	Run: func(cmd *cobra.Command, args []string) {
		// 1. Instantiate the pipeline runner, which will be used to run the tasks.
		ux := tui.TUITitle{}
		msg := tui.NewTUIMessage()
		config.ShowCLITitle()

		cliGlobalArgs := config.GetCLIGlobalArgs()

		p, err := pipeline.New(cliGlobalArgs.WorkingDir, cliGlobalArgs.MountDir,
			cliGlobalArgs.TargetDir, cliGlobalArgs.TaskName,
			cliGlobalArgs.ScanEnvVarKeys,
			cliGlobalArgs.EnvKeyValuePairsToSetString, cliGlobalArgs.ScanAWSKeys,
			cliGlobalArgs.ScanTerraformVars, cliGlobalArgs.InitDaggerWithWorkDirByDefault)

		if err != nil {
			msg.ShowError("INIT", "Failed pipeline initialization", err)
			os.Exit(1)
		}

		ux.ShowSubTitle("JOB:", "AWS-ECS")
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
			msg.ShowError("INIT", "Failed job initialization", err)
			os.Exit(1)
		}

		// 3. Run The (Docker) task
		ux.ShowSubTitle("TASK:", cliGlobalArgs.TaskName)
		ux.ShowTaskDetails("AWS:ECS", cliGlobalArgs.TaskName, j.WorkDirPath, j.TargetDirPath,
			j.MountDirPath)
		taskErr := task.RunTaskAWSECS(task.InitOptions{
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
			msg.ShowError("", "Failed to run task", taskErr)
			os.Exit(1)
		}
	},
}

func addECSCmdFlags() {
	ECSCmd.Flags().StringVarP(&ecsService, "ecs-service", "", "",
		"The name of the ECS service to be deployed.")

	ECSCmd.Flags().StringVarP(&ecsCluster, "ecs-cluster", "", "",
		"The name of the ECS cluster to be deployed.")

	ECSCmd.Flags().StringVarP(&taskDefinition, "task-definition", "", "",
		"The name of the ECS task definition to be deployed.")

	ECSCmd.Flags().StringVarP(&imageURL, "image-url", "", "",
		"The URL of the image to be deployed.")

	ECSCmd.Flags().StringVarP(&releaseVersion, "release-version", "", "",
		"The tag or version of the (container) image to be deployed. If not specified, "+
			"the default value is 'latest'.")

	err := ECSCmd.MarkFlagRequired("ecs-service")
	if err != nil {
		panic(err)
	}

	err = ECSCmd.MarkFlagRequired("ecs-cluster")
	if err != nil {
		panic(err)
	}

	err = ECSCmd.MarkFlagRequired("task-definition")
	if err != nil {
		panic(err)
	}

	_ = viper.BindPFlag("ecs-service", ECSCmd.Flags().Lookup("ecs-service"))
	_ = viper.BindPFlag("ecs-cluster", ECSCmd.Flags().Lookup("ecs-cluster"))
	_ = viper.BindPFlag("task-definition", ECSCmd.Flags().Lookup("task-definition"))
	_ = viper.BindPFlag("image-url", ECSCmd.Flags().Lookup("image-url"))
	_ = viper.BindPFlag("release-version", ECSCmd.Flags().Lookup("release-version"))
}

func init() {
	addECSCmdFlags()
}
