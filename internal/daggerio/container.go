package daggerio

import (
	"dagger.io/dagger"
	"fmt"
	"github.com/Excoriate/dagger-python-ecs/internal/common"
	"github.com/Excoriate/dagger-python-ecs/internal/errors"
	"github.com/Excoriate/dagger-python-ecs/internal/logger"
	"strings"
)

// GetContainerImageCustom DaggerContainerImage represents the container image of the dagger client.
func GetContainerImageCustom(imageURL, version string) (DaggerContainerImage, error) {
	logPrinter := logger.PipelineLogger{}
	logPrinter.InitLogger()
	var imageNormalised string
	var versionNormalised string

	if imageURL == "" {
		return DaggerContainerImage{}, errors.NewDaggerEngineError("Unable to fetch container image, "+
			"image URL value is empty",
			nil)
	}

	// If the image has ":", it means that the user has passed the version as well.
	if strings.Contains(imageURL, ":") {
		logPrinter.LogWarn("Dagger Image Configuration", "It seems that you have passed the"+
			" version of the image as"+
			" well along"+
			" with the image name. The version will be ignored and the version passed will be used.")

		imageNormalised = common.NormaliseStringLower(strings.Split(imageURL, ":")[0])
		versionNormalised = common.NormaliseStringLower(strings.Split(imageURL, ":")[1])
	} else {
		imageNormalised = imageURL

		if version == "" {
			logPrinter.LogWarn("Dagger Image Configuration",
				fmt.Sprintf("The 'version' option is empty, "+
					"it will use 'latest' as the image version for the image %s passed", imageURL))
			versionNormalised = "latest"
		} else {
			versionNormalised = version
		}
	}

	return DaggerContainerImage{
		Image:   imageNormalised,
		Version: versionNormalised,
	}, nil
}

// GetContainerImagePerStack returns the container image of the dagger client.
func GetContainerImagePerStack(stack string, version string) (string, error) {
	logPrinter := logger.PipelineLogger{}
	logPrinter.InitLogger()

	if stack == "" {
		return "", errors.NewDaggerEngineError("Unable to fetch container image, "+
			"stack value is empty",
			nil)
	}

	stackNormalised := common.NormaliseStringUpper(stack)

	if version == "" {
		logPrinter.LogWarn("Dagger Image Configuration",
			fmt.Sprintf("The 'version' option is empty, "+
				"it will use 'latest' as the image version for the image %s passed", stackNormalised))
		version = "latest"
	}

	if _, ok := StackImagesMap[stackNormalised]; !ok {
		return "", errors.NewDaggerEngineError(fmt.Sprintf("Unable to fetch container image, "+
			"stack %s is not supported or it could not be found", stackNormalised),
			nil)
	}

	return fmt.Sprintf("%s:%s", StackImagesMap[stackNormalised], version), nil
}

// GetContainer returns the container of the dagger client.
func GetContainer(c *dagger.Client, image string) (*dagger.Container, error) {
	if image == "" {
		return nil, errors.NewDaggerEngineError("Unable to fetch container, image value is empty", nil)
	}

	if c == nil {
		return nil, errors.NewDaggerEngineError("Unable to fetch container, dagger client is nil", nil)
	}

	return c.Container().From(common.NormaliseStringLower(image)), nil
}
