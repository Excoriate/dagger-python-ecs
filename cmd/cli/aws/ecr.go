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
	ecrRepositoryName string
	ecrRegistryName   string
	imageTag          string
	dockerFileName    string
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
			msg.ShowError("INIT", "Failed job initialization", err)
			os.Exit(1)
		}

		// 3. Run The (Docker) task
		ux.ShowSubTitle("TASK:", cliGlobalArgs.TaskName)
		ux.ShowTaskDetails("DOCKER", cliGlobalArgs.TaskName, j.WorkDirPath, j.TargetDirPath, j.MountDirPath)
		taskErr := task.RunTaskAWSECR(task.InitOptions{
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

func addECRCmdFlags() {
	ECRCmd.Flags().StringVarP(&ecrRepositoryName, "ecr-repository", "", "",
		"The name of the ECR repository")
	ECRCmd.Flags().StringVarP(&imageTag, "tag", "", "latest",
		"The tag of the image to be pushed. If not specified, it will default to 'latest'")
	ECRCmd.Flags().StringVarP(&dockerFileName, "dockerfile", "", "",
		"The name of the Dockerfile. If not specified, it will default to 'Dockerfile'")
	ECRCmd.Flags().StringVarP(&ecrRegistryName, "ecr-registry", "", "",
		"The name of the ECR registry.")

	err := ECRCmd.MarkFlagRequired("ecr-repository")
	if err != nil {
		panic(err)
	}

	err = ECRCmd.MarkFlagRequired("ecr-registry")
	if err != nil {
		panic(err)
	}

	_ = viper.BindPFlag("ecr-repository", ECRCmd.Flags().Lookup("ecr-repository"))
	_ = viper.BindPFlag("ecr-registry", ECRCmd.Flags().Lookup("ecr-registry"))
	_ = viper.BindPFlag("tag", ECRCmd.Flags().Lookup("tag"))
	_ = viper.BindPFlag("dockerfile", ECRCmd.Flags().Lookup("dockerfile"))
}

func init() {
	addECRCmdFlags()
}
