package daggerio

import (
	"dagger.io/dagger"
	"github.com/Excoriate/dagger-python-ecs/internal/errors"
)

// GetWorkDir returns the working directory of the dagger client.
func GetWorkDir(c *dagger.Client, workdir string) (*dagger.Directory, error) {
	if workdir == "" {
		return nil, errors.NewDaggerEngineError("Working directory is empty", nil)
	}

	return c.Host().Directory(workdir), nil
}
