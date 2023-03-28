package task

import (
	"fmt"
	"github.com/Excoriate/dagger-python-ecs/internal/common"
	"github.com/Excoriate/dagger-python-ecs/internal/filesystem"
	"github.com/Excoriate/dagger-python-ecs/pkg/job"
	"github.com/Excoriate/dagger-python-ecs/pkg/pipeline"
)

func NewTask(p *pipeline.Config, job *job.Job, actions []string,
	init *InitOptions) Task {
	taskId := common.GetUUID()
	taskName := "BUILD"
	stackName := "DOCKER"

	p.UXMessage.ShowInfo("TASK", fmt.Sprintf("Initialising task: %s with id: %s",
		taskName, taskId))

	envVars := filesystem.MergeEnvVars(job.EnvVarsToSet, job.EnvVarsAWSScanned,
		job.EnvVarsCustomScanned, job.EnvVarsTerraformScanned)

	randomContainerName := common.GenerateRandomStringWithPrefix(3, false, true, false,
		"dagtaskdef")

	t := Task{
		// Identifiers
		Id:    taskId,
		Name:  taskName,
		Stack: stackName,
		// Parent objects
		PipelineCfg:           p,
		JobCfg:                job,
		EnvVarsInheritFromJob: envVars,

		// Default inherited container runtime from the job instantiated.
		ContainerImageDefault: job.ContainerImageURL,
		ContainerDefault:      job.ContainerDefault,
		ContainerNameDefault:  randomContainerName,

		Dirs: Dirs{
			RootDir:         ".",
			WorkDir:         init.WorkDir,
			MountDir:        init.MountDir,
			TargetDir:       init.TargetDir,
			RootDirDagger:   job.RootDir,
			WorkDirDagger:   job.WorkDir,
			MountDirDagger:  job.MountDir,
			TargetDirDagger: job.TargetDir,
		},

		PreReqs: PreRequisites{
			Files: []string{"Dockerfile"},
		},

		Actions: Actions{
			CustomCommands:  actions,
			DefaultCommands: []string{"docker", "build", "-t", randomContainerName, "."},
		},
	}

	return t
}
