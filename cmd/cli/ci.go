package cli

import (
	"github.com/spf13/cobra"
)

var (
	CITaskBuild      bool
	CITaskTest       bool
	CITaskLint       bool
	CITaskInstallDep bool
)

var CI = &cobra.Command{
	Version: "v0.0.1",
	Use:     "ci",
	Long: `The 'ci' subcommand automates various Continuous Integration (CI) tasks such as building Docker images, running tests, linting, and installing dependencies.
You can specify the tasks you want to perform using the provided flags.`,
	Example: `
  # Run pipeline CI with build, test, and lint tasks
  pipeline ci --build
  # Run pipeline CI with install dependencies task
  pipeline ci --install`,
	Run: func(cmd *cobra.Command, args []string) {
		_ = cmd.Help()
	},
}

func AddCIArguments() {
	CI.Flags().BoolVarP(&CITaskBuild,
		"build",
		"b", false,
		"If the project includes a DockerFile, it'll build the image.")

	CI.Flags().BoolVarP(&CITaskTest,
		"test",
		"t", false,
		"Run all the existing test.")

	CI.Flags().BoolVarP(&CITaskLint,
		"lint",
		"l", false,
		"Run all the existing linters.")

	CI.Flags().BoolVarP(&CITaskInstallDep,
		"install-dep",
		"i", false,
		"Install all the dependencies.")
}

func init() {
	AddCIArguments()
}
