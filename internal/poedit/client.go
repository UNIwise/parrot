package poedit

import (
	"context"
	"net/http"

	"gopkg.in/resty.v1"
)

const (
	HostURL = "https://api.poeditor.com"
)

// Client is an interface to poeditors API.
type Client interface {
	ExportProject(ctx context.Context, req ExportProjectRequest) (result *ExportProjectResponse, err error)
	ListProjectLanguages(ctx context.Context, req ListProjectLanguagesRequest) (result *ListProjectLanguagesResponse, err error)
}

// ClientImpl is an implementation of the poeditor client interface.
type ClientImpl struct {
	r *resty.Client
}

// NewClient creates a new poeditor api client.
func NewClient(apiToken string, httpClient *http.Client) *ClientImpl {
	r := resty.NewWithClient(httpClient)
	r.FormData.Add("api_token", apiToken)
	r.SetHostURL(HostURL)

	return &ClientImpl{
		r: r,
	}
}
