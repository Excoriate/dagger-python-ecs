package cli

import (
	"context"
	"github.com/spf13/cobra"
	"os"
)

var (
	GlobalWorkDir               string
	GlobalTargetDir             string
	GlobalTaskName              string
	GlobalScanEnvVarKeys        []string
	GlobalEnvKeyValuePairsToSet map[string]string
	GlobalCustomCommands        []string
	GlobalScanAWSKeys           bool
)

var rootCmd = &cobra.Command{
	Version: "v0.0.1",
	Use:     "pipeline",
	Long: `Pipeline is a command-line tool that helps automate the process of Continuous
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
	rootCmd.AddCommand(DockerCMD)
}
