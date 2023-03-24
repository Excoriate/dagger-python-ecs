package cli

import (
	"context"
	"github.com/spf13/cobra"
	"os"
)

var Workdir string

var rootCmd = &cobra.Command{
	Version: "v0.0.1",
	Use:     "pipeline",
	Long: `Pipeline is a command-line tool that helps automate the process of Continuous
Integration (CI) and Continuous Deployment (CD), using Dagger (dagger.io) as the workflow engine.
It provides options to manage various tasks such as building Docker images, running tests, linting, and deployment to AWS services.`,
	Example: `
  # Run pipeline with the specified working directory
  pipeline --workdir /path/to/working/directory <subcommand> <task>,
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

func AddRootArguments() {
	rootCmd.Flags().StringVarP(&Workdir,
		"workdir",
		"w", "",
		"Working directory where the pipeline will be executed")

	_ = rootCmd.MarkFlagRequired("workdir")
}

func init() {
	AddRootArguments()
	rootCmd.AddCommand(CI)
	rootCmd.AddCommand(CD)
}
