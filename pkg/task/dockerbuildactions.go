package task

type DockerBuildAction struct {
	Task   CoreTasker
	prefix string // How the UX messages should be prefixed
	Id     string // The ID of the task
	Name   string // The name of the task
}

type DockerBuildActions interface {
	BuildFromDockerFile(dockerFile string) (Output, error)
}

func (t *DockerBuildAction) BuildFromDockerFile(dockerFile string) (Output, error) {
	//t.ux.ShowSuccess(t.prefix, fmt.Sprintf("Building Docker image from Dockerfile: %s",
	//	dockerFile))
	return Output{}, nil
}

func NewDockerBuildAction(task CoreTasker) DockerBuildActions {
	return &DockerBuildAction{
		Task:   task,
		prefix: "DOCKER-ACTION:BUILD",
	}
}
