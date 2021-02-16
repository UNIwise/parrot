package poedit

import (
	"fmt"

	"github.com/pkg/errors"
)

var ErrFailedToUnmarshalResponse = errors.New("Failed to unmarshal response")

type ErrProjectPermissionDenied struct {
	ProjectID int
}

func (e *ErrProjectPermissionDenied) Error() string {
	return fmt.Sprintf(
		"You don't have permission to access project %d",
		e.ProjectID,
	)
}

type ErrLanguageNotFound struct {
	ProjectID    int
	LanguageCode string
}

func (e *ErrLanguageNotFound) Error() string {
	return fmt.Sprintf(
		"Project %d does not contain specified language %s",
		e.ProjectID,
		e.LanguageCode,
	)
}
