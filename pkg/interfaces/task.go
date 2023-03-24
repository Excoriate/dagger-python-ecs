package interfaces

import (
	"context"
	"dagger.io/dagger"
	"github.com/Excoriate/dagger-python-ecs/internal/logger"
)

type Task struct {
	Id     string
	Logger logger.PipelineLogger
}

type TaskRunner interface {
	InitClient(ctx *context.Context) (*dagger.Client, error)
	InitWorkDir(workdir string) error
	InitContainer() (*dagger.Container, error)
}
