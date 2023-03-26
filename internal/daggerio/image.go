package daggerio

var StackImagesMap = map[string]string{
	"PYTHON":     "python:3.8.5-slim-buster",
	"DOCKER":     "docker:23.0.1-dind",
	"TERRAFORM":  "hashicorp/terraform:1.3.9",
	"TERRAGRUNT": "alpine/terragrunt",
	"AWSCLI":     "amazon/aws-cli:2.1.19",
}

type DaggerContainerImage struct {
	Image   string
	Version string
}
