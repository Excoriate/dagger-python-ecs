package cli

import (
	"github.com/spf13/cobra"
)

var (
	CDTaskECRPush   bool
	CDTaskDeployECS bool
)

var CD = &cobra.Command{
	Version: "v0.0.1",
	Use:     "cd",
	Long: `The 'cd' subcommand automates various Continuous Deployment (CD) tasks such as pushing Docker images to Amazon Elastic Container Registry (ECR) and deploying them to Amazon Elastic Container Service (ECS).
You can specify the tasks you want to perform using the provided flags.`,
	Example: `
  # Run pipeline CD with ECR push task
  pipeline cd --ecr-push
  # Run pipeline CD with ECS deploy task
  pipeline cd --deploy-ecs`,
	Run: func(cmd *cobra.Command, args []string) {
		_ = cmd.Help()
	},
}

func AddCDArguments() {
	CD.Flags().BoolVarP(&CDTaskECRPush,
		"ecr-push",
		"e", false,
		"Push the image to ECR.")

	CD.Flags().BoolVarP(&CDTaskDeployECS,
		"deploy-ecs",
		"d", false,
		"Deploy the image to ECS.")
}

func init() {
	AddCDArguments()
}
