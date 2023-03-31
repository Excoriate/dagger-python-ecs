package aws

import (
	"github.com/spf13/cobra"
)

var Cmd = &cobra.Command{
	Version: "v0.0.1",
	Use:     "aws",
	Long: `The 'aws' command automate and perform several aws-related actions (E.
g: push images to ECR, deploy into ECS, etc.).
You can specify the tasks you want to perform using the provided --task flag.`,
	Example: `
  # Push an image into ECR:
  stiletto aws ecr --task=push`,
	Run: func(cmd *cobra.Command, args []string) {
		_ = cmd.Help()
	},
}

func init() {
	Cmd.AddCommand(ECRCmd)
}
