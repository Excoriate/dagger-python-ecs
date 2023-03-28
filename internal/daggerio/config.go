package daggerio

import (
	"dagger.io/dagger"
	"github.com/Excoriate/dagger-python-ecs/internal/common"
	"github.com/Excoriate/dagger-python-ecs/internal/logger"
)

const ContainerMountPathPrefix = "/build"

func SetEnvVarsInContainer(c *dagger.Container, envVars map[string]string) (*dagger.Container,
	error) {
	logPrinter := logger.PipelineLogger{}
	logPrinter.InitLogger()

	if common.MapIsNulOrEmpty(envVars) {
		logPrinter.LogWarn("Dagger Container Configuration",
			"No environment variables are passed, skipping the environment variable configuration step")
		return c, nil
	}

	for k, v := range envVars {
		c = c.WithEnvVariable(k, v)
	}

	return c, nil
}
