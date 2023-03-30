package daggerio

import (
	"context"
	"dagger.io/dagger"
	"fmt"
	"github.com/Excoriate/dagger-python-ecs/internal/errors"
	"github.com/Excoriate/dagger-python-ecs/internal/filesystem"
)

// GetMountDir returns the working directory of the dagger client.
func GetMountDir(c *dagger.Client, mountDir string) (*dagger.Directory, error) {
	if mountDir == "" {
		return c.Host().Directory("."), nil
	}

	return c.Host().Directory(mountDir), nil
}

// GetDaggerDir returns the working directory of the dagger client.
func GetDaggerDir(c *dagger.Client, dir string) (*dagger.Directory, error) {
	if dir == "" {
		return c.Host().Directory("."), nil // Which will map to the current directory.
	}

	if err := filesystem.DirExist(dir); err != nil {
		return nil, errors.NewDaggerConfigurationError(
			fmt.Sprintf("Failed to create a dagger directory, directory: %s does not exist", dir),
			err)
	}

	if err := filesystem.PathIsADirectory(dir); err != nil {
		return nil, errors.NewDaggerConfigurationError(
			fmt.Sprintf("Failed to create a dagger directory, directory: %s is not a directory",
				dir), err)
	}

	return c.Host().Directory(dir), nil
}

// GetDaggerDirWithEntriesCheck returns the working directory of the dagger client.
func GetDaggerDirWithEntriesCheck(c *dagger.Client, dir string) (*dagger.Directory, error) {
	if dir == "" {
		return c.Host().Directory("."), nil // Which will map to the current directory.
	}

	if err := filesystem.DirExist(dir); err != nil {
		return nil, errors.NewDaggerConfigurationError(
			fmt.Sprintf("Failed to create a dagger directory, directory: %s does not exist", dir),
			err)
	}

	if err := filesystem.PathIsADirectory(dir); err != nil {
		return nil, errors.NewDaggerConfigurationError(
			fmt.Sprintf("Failed to create a dagger directory, directory: %s is not a directory",
				dir), err)
	}

	ctx := context.Background()
	if _, err := ListEntries(c.Host().Directory(dir), true, &ctx); err != nil {
		return nil, err
	}

	return c.Host().Directory(dir), nil
}

// ListEntries lists the entries in a dagger directory.
func ListEntries(d *dagger.Directory, failIsEmpty bool, ctx *context.Context) ([]string, error) {
	entries, err := d.Entries(*ctx)
	if err != nil {
		return nil, errors.NewDaggerConfigurationError(
			"Failed to list entries in dagger directory", err)
	}

	if len(entries) == 0 && failIsEmpty {
		return nil, errors.NewDaggerConfigurationError(
			"The directory passed is empty", nil)
	}

	return entries, nil
}

// MountDir mounts a directory from the host to the container.
func MountDir(c *dagger.Container, workDir *dagger.Directory, execPath string) (*dagger.
Container, error) {
	mountPathInContainer := ContainerMountPathPrefix

	if execPath == "" {
		execPath = mountPathInContainer
	}

	execPathNormalised := NormaliseDaggerPath(execPath)

	// --------------------------------------------------------------
	//_, err = c.Container().From("docker:23.0.1-dind").
	//	WithMountedDirectory("/build", c.Host().Directory(".")).
	//	WithWorkdir("/build/examples/docker").
	//	WithExec([]string{"ls", "-ltrh"}).
	//	ExitCode(t.PipelineCfg.Ctx)
	// --------------------------------------------------------------

	container := c.
		WithMountedDirectory(mountPathInContainer, workDir).WithWorkdir(
		execPathNormalised)

	return container, nil
}
