package errors

import "fmt"

const daggerEngineErrorPrefix = "Dagger engine error: "

type DaggerEngineError struct {
	Details string
	Err     error
}

func (e *DaggerEngineError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("%s: %s: %s", daggerEngineErrorPrefix, e.Details, e.Err.Error())
	}
	return fmt.Sprintf("%s: %s", daggerEngineErrorPrefix, e.Details)
}

func NewDaggerEngineError(details string, err error) *DaggerEngineError {
	return &DaggerEngineError{
		Details: details,
		Err:     err,
	}
}
