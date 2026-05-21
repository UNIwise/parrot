package poedit

import (
	"fmt"

	"github.com/pkg/errors"
)

var (
	ErrFailedToUnmarshalResponse = errors.New("Failed to unmarshal response")
	ErrNotImplemented            = errors.New("Method is not implemented")
)

type ProjectPermissionDeniedError struct {
	ProjectID int
}

func (e *ProjectPermissionDeniedError) Error() string {
	return fmt.Sprintf(
		"You don't have permission to access project %d",
		e.ProjectID,
	)
}

type LanguageNotFoundError struct {
	ProjectID    int
	LanguageCode string
}

func (e *LanguageNotFoundError) Error() string {
	return fmt.Sprintf(
		"Project %d does not contain specified language %s",
		e.ProjectID,
		e.LanguageCode,
	)
}
