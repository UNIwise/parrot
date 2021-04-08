package poedit

import (
	"context"
	"net/http"

	"gopkg.in/resty.v1"
)

const (
	HostURL    = "https://api.poeditor.com"
	TimeFormat = "2006-01-02T15:04:05+0000"
)

// Client is an interface to poeditors API.
type Client interface {
	ViewProject(ctx context.Context, r ViewProjectRequest) (*ViewProjectResponse, error)
	AddProject(ctx context.Context, r AddProjectRequest) (*AddProjectResponse, error)
	ExportProject(ctx context.Context, req ExportProjectRequest) (result *ExportProjectResponse, err error)

	ListProjectLanguages(ctx context.Context, req ListProjectLanguagesRequest) (result *ListProjectLanguagesResponse, err error)
	AddProjectLanguage(ctx context.Context, r AddProjectLanguageRequest) (*AddProjectLanguageResponse, error)
	UpdateProjectLanguage(ctx context.Context, r UpdateProjectLanguageRequest) (*UpdateProjectLanguageResponse, error)
	DeleteProjectLanguage(ctx context.Context, r DeleteProjectLanguageRequest) (*DeleteProjectLanguageResponse, error)
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
