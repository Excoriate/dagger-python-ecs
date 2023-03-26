package daggerio

import (
	"context"
	"dagger.io/dagger"
	"os"
)

func NewDaggerClient(workDir string, ctx context.Context) (*dagger.Client, error) {
	if workDir == "" {
		client, err := dagger.Connect(ctx, dagger.WithLogOutput(os.Stdout))
		if err != nil {
			return nil, err
		}
		return client, nil
	}

	client, err := dagger.Connect(ctx, dagger.WithLogOutput(os.Stdout),
		dagger.WithWorkdir(workDir))
	if err != nil {
		return nil, err
	}

	return client, nil
}
