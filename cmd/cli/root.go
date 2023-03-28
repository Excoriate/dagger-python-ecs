package cli

import (
	"context"
	"github.com/spf13/cobra"
	"os"
)

var (
	GlobalWorkingDirectory       string
	GlobalMountDir               string
	GlobalTargetDir              string
	GlobalTaskName               string
	GlobalScanEnvVarKeys         []string
	GlobalEnvKeyValuePairsToSet  map[string]string
	GlobalCustomCommands         []string
	GlobalScanAWSKeys            bool
	GlobalScanTFVars             bool
	GlobalSetWorkDirDaggerOnInit bool
)

var rootCmd = &cobra.Command{
	Version: "v0.0.1",
	Use:     "pipeline",
	Long: `PipelineCfg is a command-line tool that helps automate the process of Continuous
Integration (DockerCMD) and Continuous Deployment (CD), using Dagger (dagger.io) as the workflow engine.
It provides options to manage various tasks such as building Docker images, running tests, linting, and deployment to AWS services.`,
	Example: `
  # DockerCMD pipeline with the specified working directory
  pipeline <command> --workdir /path/to/working/directory
  E.g.:
  pipeline --workdir /path/to/working/directory ci --build`,
	Run: func(cmd *cobra.Command, args []string) {
		_ = cmd.Help()
	},
}

func Execute() {
	err := rootCmd.ExecuteContext(context.Background())
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().StringVarP(&GlobalTaskName,
		"task",
		"t", "",
		"Name of the task to run. E.g.: build, test, lint, install")

	rootCmd.PersistentFlags().StringVarP(&GlobalWorkingDirectory,
		"work-dir",
		"w", "",
		"Work directory where the pipeline will be executed. If it's not set, "+
			"it'll use '.' value, which represents the current directory.")

	rootCmd.PersistentFlags().StringVarP(&GlobalTargetDir,
		"target-dir",
		"d", "",
		"Target directory represents the subdirectory within the mounted directory where the"+
			" actions (commands) will be executed.")

	rootCmd.PersistentFlags().StringVarP(&GlobalMountDir,
		"mount-dir",
		"m", "",
		"Mount directory represents what subdirectory within the working directory will be used"+
			" to mount into the container while it's performing its actions.")

	rootCmd.PersistentFlags().StringSliceVarP(&GlobalScanEnvVarKeys,
		"scan-env",
		"s", []string{},
		"List of environment variable keys that are already exported, that'll be scanned and set.")

	rootCmd.PersistentFlags().StringToStringVarP(&GlobalEnvKeyValuePairsToSet,
		"set-env",
		"e", map[string]string{},
		"List of environment variable key-value pairs to set.")

	rootCmd.PersistentFlags().StringSliceVarP(&GlobalCustomCommands,
		"commands",
		"c", []string{},
		"List of custom commands to run.")

	rootCmd.PersistentFlags().BoolVarP(&GlobalScanAWSKeys,
		"scan-aws-keys",
		"a", false,
		"Scan AWS keys and set them as environment variables.")

	rootCmd.PersistentFlags().BoolVarP(&GlobalScanTFVars,
		"scan-terraform-vars",
		"f", false,
		"Scan terraform exported environment variables and set it into the generated containers ("+
			"TG_VAR_).")

	rootCmd.PersistentFlags().BoolVarP(&GlobalSetWorkDirDaggerOnInit,
		"set-workdir-dagger-on-init",
		"i", false,
		"Set the working directory to the Dagger's working directory on initialisation of the"+
			" Dagger client. By default, it's 'false'.")

	_ = rootCmd.MarkFlagRequired("task")
	_ = rootCmd.MarkFlagRequired("workdir")

	rootCmd.AddCommand(DockerCMD)
}
