package config

type PipelineOptions struct {
	WorkDir               string
	TargetDir             string
	TaskName              string
	EnvVarsToScanAndSet   []string
	EnvKeyValuePairsToSet map[string]string
	EnvVarsAWSKeysToScan  map[string]string
	IsAWSKeysToScan       bool
}
