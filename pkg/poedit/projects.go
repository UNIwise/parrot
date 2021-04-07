package poedit

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"strings"

	"github.com/pkg/errors"
)

type ListProjectsResponse struct {
	Response struct {
		Status  string `json:"status"`
		Code    string `json:"code"`
		Message string `json:"message"`
	} `json:"response"`
	Result struct {
		Projects []struct {
			ID      int64  `json:"id"`
			Name    string `json:"name"`
			Public  int64  `json:"public"`
			Open    int64  `json:"open"`
			Created string `json:"created"`
		} `json:"projects"`
	} `json:"result"`
}

func (c *ClientImpl) ListProjects(ctx context.Context) (*ListProjectsResponse, error) {
	req := c.r.R()

	req.SetContext(ctx)

	req.SetResult(&ListProjectsResponse{})

	resp, err := req.Post("/v2/projects/list")
	if err != nil {
		return nil, err
	}

	res, ok := resp.Result().(*ListProjectsResponse)
	if !ok {
		return nil, ErrFailedToUnmarshalResponse
	}

	if res.Response.Code != "200" {
		return nil, errors.New(res.Response.Message)
	}

	return res, nil
}

type ViewProjectRequest struct {
	ID int
}

type ViewProjectResponse struct {
	Response struct {
		Status  string `json:"status"`
		Code    string `json:"code"`
		Message string `json:"message"`
	} `json:"response"`
	Result struct {
		Project struct {
			ID                int64  `json:"id"`
			Name              string `json:"name"`
			Description       string `json:"description"`
			Public            int64  `json:"public"`
			Open              int64  `json:"open"`
			ReferenceLanguage string `json:"reference_language"`
			Terms             int64  `json:"terms"`
			Created           string `json:"created"`
		} `json:"project"`
	} `json:"result"`
}

func (c *ClientImpl) ViewProject(ctx context.Context, r ViewProjectRequest) (*ViewProjectResponse, error) {
	req := c.r.R()

	req.SetFormData(map[string]string{
		"id": fmt.Sprintf("%d", r.ID),
	})

	req.SetContext(ctx)

	req.SetResult(&ViewProjectResponse{})

	resp, err := req.Post("/v2/projects/view")
	if err != nil {
		return nil, err
	}

	res, ok := resp.Result().(*ViewProjectResponse)
	if !ok {
		return nil, ErrFailedToUnmarshalResponse
	}

	if res.Response.Code != "200" {
		return nil, errors.New(res.Response.Message)
	}

	return res, nil
}

type AddProjectRequest struct {
	Name        string
	Description string
}

type AddProjectResponse struct {
	Response struct {
		Status  string `json:"status"`
		Code    string `json:"code"`
		Message string `json:"message"`
	} `json:"response"`
	Result struct {
		Project struct {
			ID                int64  `json:"id"`
			Name              string `json:"name"`
			Description       string `json:"description"`
			Public            int64  `json:"public"`
			Open              int64  `json:"open"`
			ReferenceLanguage string `json:"reference_language"`
			Terms             int64  `json:"terms"`
			Created           string `json:"created"`
		} `json:"project"`
	} `json:"result"`
}

func (c *ClientImpl) AddProject(ctx context.Context, r AddProjectRequest) (*AddProjectResponse, error) {
	req := c.r.R()

	req.SetFormData(map[string]string{
		"name":        r.Name,
		"description": r.Description,
	})

	req.SetContext(ctx)

	req.SetResult(&AddProjectResponse{})

	resp, err := req.Post("/v2/projects/add")
	if err != nil {
		return nil, err
	}

	res, ok := resp.Result().(*AddProjectResponse)
	if !ok {
		return nil, ErrFailedToUnmarshalResponse
	}

	if res.Response.Code != "200" {
		return nil, errors.New(res.Response.Message)
	}

	return res, nil
}

type UpdateProjectSettingsRequest struct {
	ID                int
	Name              string
	Description       string
	ReferenceLanguage string
}

type UpdateProjectSettingsResponse struct {
	Response struct {
		Status  string `json:"status"`
		Code    string `json:"code"`
		Message string `json:"message"`
	} `json:"response"`
	Result struct {
		Project struct {
			ID                int64  `json:"id"`
			Name              string `json:"name"`
			Description       string `json:"description"`
			Public            int64  `json:"public"`
			Open              int64  `json:"open"`
			ReferenceLanguage string `json:"reference_language"`
			Terms             int64  `json:"terms"`
			Created           string `json:"created"`
		} `json:"project"`
	} `json:"result"`
}

func (c *ClientImpl) UpdateProjectSettings(ctx context.Context, r UpdateProjectSettingsRequest) (*UpdateProjectSettingsResponse, error) {
	req := c.r.R()

	req.SetFormData(map[string]string{
		"id":                 fmt.Sprintf("%d", r.ID),
		"name":               r.Name,
		"description":        r.Description,
		"reference_language": r.ReferenceLanguage,
	})

	req.SetContext(ctx)

	req.SetResult(&UpdateProjectSettingsResponse{})

	resp, err := req.Post("/v2/projects/update")
	if err != nil {
		return nil, err
	}

	res, ok := resp.Result().(*UpdateProjectSettingsResponse)
	if !ok {
		return nil, ErrFailedToUnmarshalResponse
	}

	if res.Response.Code != "200" {
		return nil, errors.New(res.Response.Message)
	}

	return res, nil
}

type DeleteProjectRequest struct {
	ID int
}

type DeleteProjectResponse struct {
	Response struct {
		Status  string `json:"status"`
		Code    string `json:"code"`
		Message string `json:"message"`
	} `json:"response"`
}

func (c *ClientImpl) DeleteProject(ctx context.Context, r DeleteProjectRequest) (*DeleteProjectResponse, error) {
	req := c.r.R()

	req.SetFormData(map[string]string{
		"id": fmt.Sprintf("%d", r.ID),
	})

	req.SetContext(ctx)

	req.SetResult(&DeleteProjectResponse{})

	resp, err := req.Post("/v2/projects/delete")
	if err != nil {
		return nil, err
	}

	res, ok := resp.Result().(*DeleteProjectResponse)
	if !ok {
		return nil, ErrFailedToUnmarshalResponse
	}

	if res.Response.Code != "200" {
		return nil, errors.New(res.Response.Message)
	}

	return res, nil
}

type UploadProjectRequest struct {
	ID             int
	Updating       string
	File           io.Reader
	Language       string
	Overwrite      bool
	SyncTerms      bool
	Tags           string
	ReadFromSource bool
	FuzzyTrigger   bool
}

type UploadProjectResponse struct {
	Response struct {
		Status  string `json:"status"`
		Code    string `json:"code"`
		Message string `json:"message"`
	} `json:"response"`
	Result struct {
		Terms struct {
			Parsed  int64 `json:"parsed"`
			Added   int64 `json:"added"`
			Deleted int64 `json:"deleted"`
		} `json:"terms"`
		Translations struct {
			Parsed  int64 `json:"parsed"`
			Added   int64 `json:"added"`
			Updated int64 `json:"updated"`
		} `json:"translations"`
	} `json:"result"`
}

func (c *ClientImpl) UploadProject(ctx context.Context, r UploadProjectRequest) (*UploadProjectResponse, error) {
	return nil, ErrNotImplemented
}

type SyncProjectTermsRequest struct {
	ID   int
	Data []struct {
		Term      string   `json:"term"`
		Context   string   `json:"context"`
		Reference string   `json:"reference"`
		Plural    string   `json:"plural"`
		Comment   string   `json:"comment,omitempty"`
		Tags      []string `json:"tags"`
	}
}

type SyncProjectTermsResponse struct {
	Response struct {
		Status  string `json:"status"`
		Code    string `json:"code"`
		Message string `json:"message"`
	} `json:"response"`
	Result struct {
		Terms struct {
			Parsed  int64 `json:"parsed"`
			Added   int64 `json:"added"`
			Updated int64 `json:"updated"`
			Deleted int64 `json:"deleted"`
		} `json:"terms"`
	} `json:"result"`
}

func (c *ClientImpl) SyncProjectTerms(ctx context.Context, r SyncProjectTermsRequest) (*SyncProjectTermsResponse, error) {
	req := c.r.R()

	data, err := json.Marshal(r.Data)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to marshal data")
	}

	req.SetFormData(map[string]string{
		"id":   fmt.Sprintf("%d", r.ID),
		"data": string(data),
	})

	req.SetContext(ctx)

	req.SetResult(&SyncProjectTermsResponse{})

	resp, err := req.Post("/v2/projects/sync")
	if err != nil {
		return nil, err
	}

	res, ok := resp.Result().(*SyncProjectTermsResponse)
	if !ok {
		return nil, ErrFailedToUnmarshalResponse
	}

	if res.Response.Code != "200" {
		return nil, errors.New(res.Response.Message)
	}

	return res, nil
}

type ExportProjectRequest struct {
	ID       int
	Language string
	Type     string
	Order    string
	Tags     []string
	Filters  []string
	Options  []string
}

type ExportProjectResponse struct {
	Response struct {
		Status  string `json:"status"`
		Code    string `json:"code"`
		Message string `json:"message"`
	} `json:"response"`
	Result struct {
		URL string `json:"url"`
	} `json:"result"`
}

func (c *ClientImpl) ExportProject(ctx context.Context, r ExportProjectRequest) (*ExportProjectResponse, error) {
	req := c.r.R()

	req.SetFormData(map[string]string{
		"id":       fmt.Sprintf("%d", r.ID),
		"language": r.Language,
		"type":     r.Type,
		"order":    r.Order,
		"tags":     formDataArray(r.Tags),
		"filters":  formDataArray(r.Filters),
		"options":  formDataArray(r.Options),
	})

	req.SetContext(ctx)

	req.SetResult(&ExportProjectResponse{})

	resp, err := req.Post("/v2/projects/export")
	if err != nil {
		return nil, err
	}

	res, ok := resp.Result().(*ExportProjectResponse)
	if !ok {
		return nil, ErrFailedToUnmarshalResponse
	}

	if res.Response.Code == "403" {
		return nil, &ErrProjectPermissionDenied{
			ProjectID: r.ID,
		}
	}

	if res.Response.Code == "4044" {
		return nil, &ErrLanguageNotFound{
			ProjectID:    r.ID,
			LanguageCode: r.Language,
		}
	}

	if res.Response.Code != "200" {
		return nil, errors.New(res.Response.Message)
	}

	return res, nil
}

func formDataArray(data []string) string {
	if len(data) > 0 {
		return fmt.Sprintf("[\"%s\"]", strings.Join(data, "\",\""))
	}

	return "[]"
}
