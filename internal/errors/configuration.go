package errors

import "fmt"

type PipelineConfigurationError struct {
	Details string
	Err     error
}

func (e *PipelineConfigurationError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("Pipeline configuration error: %s: %s", e.Details, e.Err.Error())
	}
	return fmt.Sprintf("Pipeline configuration error: %s", e.Details)
}

func NewPipelineConfigurationError(details string, err error) *PipelineConfigurationError {
	return &PipelineConfigurationError{
		Details: fmt.Sprintf("Unable to read file %s", details),
		Err:     err,
	}
}
