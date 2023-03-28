package task

type DockerTask struct {
	// Identifiers.
	Id    string
	Name  string
	Stack string

	JobCfg *Task
}
