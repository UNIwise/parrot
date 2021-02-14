package poedit

import "context"

type Client interface {
	// Not ExportProject in order to seperate it from the api call
	FetchTerms(ctx context.Context, projectID int, language, format string) (result []byte, err error)
}

type ClientImpl struct{}

func NewClient() *ClientImpl {
	return &ClientImpl{}
}
