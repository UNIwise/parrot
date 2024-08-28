package poedit

import (
	"context"
	"net/http"

	"gopkg.in/resty.v1"
)

// Client is an interface to poeditors api
type Client interface {
	ExportProject(ctx context.Context, req ExportProjectRequest) (result *ExportProjectResponse, err error)
	ListProjectLanguages(ctx context.Context, r ListProjectLanguagesRequest) (*ListProjectLanguagesResponse, error)
	ViewProject(ctx context.Context, r ViewProjectRequest) (*ViewProjectResponse, error)
	ListProjects(ctx context.Context) (*ListProjectsResponse, error)
}

// ClientImpl is an implementation of the poeditor client interface
type ClientImpl struct {
	r *resty.Client
}

// NewClient creates a new poeditor api client
func NewClient(apiToken string, httpClient *http.Client) *ClientImpl {
	r := resty.NewWithClient(httpClient)
	r.FormData.Add("api_token", apiToken)
	r.SetHostURL("https://api.poeditor.com")

	return &ClientImpl{
		r: r,
	}
}
