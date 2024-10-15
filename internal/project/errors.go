package project

import (
	"fmt"
)

type ErrLanguageNotFoundInStorage struct {
	ProjectID    int
	LanguageCode string
	Version 	string
}

func (e *ErrLanguageNotFoundInStorage) Error() string {
	return fmt.Sprintf(
		"Project %d does not contain specified language %s with version %s in storage",
		e.ProjectID,
		e.LanguageCode,
		e.Version,
	)
}
