package pipeline

import (
	"context"
	"dagger.io/dagger"
	"github.com/Excoriate/dagger-python-ecs/internal/logger"
	"github.com/Excoriate/dagger-python-ecs/internal/tui"
	"github.com/Excoriate/dagger-python-ecs/pkg/config"
)

type Runner struct {
	Logger       logger.Logger
	Dirs         config.DefaultDirs
	UXDisplay    tui.TUIDisplayer
	UXMessage    tui.TUIMessenger
	Platforms    map[dagger.Platform]string
	PipelineOpts *config.PipelineOptions
	Ctx          context.Context
}
