package poedit

import (
	"context"
	"fmt"

	"github.com/pkg/errors"
)

type ListAvailableLanguagesResponse struct {
	Response struct {
		Status  string `json:"status"`
		Code    string `json:"code"`
		Message string `json:"message"`
	} `json:"response"`
	Result struct {
		Languages []struct {
			Name string `json:"name"`
			Code string `json:"code"`
		} `json:"languages"`
	} `json:"result"`
}

func (c *ClientImpl) ListAvailableLanguages(ctx context.Context) (*ListAvailableLanguagesResponse, error) {
	req := c.r.R()

	req.SetContext(ctx)

	req.SetResult(&ListAvailableLanguagesResponse{})

	resp, err := req.Post("/v2/languages/available")
	if err != nil {
		return nil, err
	}

	res, ok := resp.Result().(*ListAvailableLanguagesResponse)
	if !ok {
		return nil, ErrFailedToUnmarshalResponse
	}

	if res.Response.Code != "200" {
		return nil, errors.New(res.Response.Message)
	}

	return res, nil
}

type ListProjectLanguagesRequest struct {
	ID int
}

type ListProjectLanguagesResponse struct {
	Response struct {
		Status  string `json:"status"`
		Code    string `json:"code"`
		Message string `json:"message"`
	} `json:"response"`
	Result struct {
		Languages []struct {
			Name         string  `json:"name"`
			Code         string  `json:"code"`
			Translations int64   `json:"translations"`
			Percentage   float64 `json:"percentage"`
			Updated      string  `json:"updated"`
		} `json:"languages"`
	} `json:"result"`
}

func (c *ClientImpl) ListProjectLanguages(ctx context.Context, r ListProjectLanguagesRequest) (*ListProjectLanguagesResponse, error) {
	req := c.r.R()

	req.SetFormData(map[string]string{
		"id": fmt.Sprintf("%d", r.ID),
	})

	req.SetContext(ctx)

	req.SetResult(&ListProjectLanguagesResponse{})

	resp, err := req.Post("/v2/languages/list")
	if err != nil {
		return nil, err
	}

	res, ok := resp.Result().(*ListProjectLanguagesResponse)
	if !ok {
		return nil, ErrFailedToUnmarshalResponse
	}

	if res.Response.Code != "200" {
		return nil, errors.New(res.Response.Message)
	}

	return res, nil
}

type AddProjectLanguageRequest struct {
	ID       int
	Language string
}

type AddProjectLanguageResponse struct {
	Response struct {
		Status  string `json:"status"`
		Code    string `json:"code"`
		Message string `json:"message"`
	} `json:"response"`
}

func (c *ClientImpl) AddProjectlanguage(ctx context.Context, r AddProjectLanguageRequest) (*AddProjectLanguageResponse, error) {
	req := c.r.R()

	req.SetFormData(map[string]string{
		"id":       fmt.Sprintf("%d", r.ID),
		"language": r.Language,
	})

	req.SetContext(ctx)

	req.SetResult(&AddProjectLanguageResponse{})

	resp, err := req.Post("/v2/languages/add")
	if err != nil {
		return nil, err
	}

	res, ok := resp.Result().(*AddProjectLanguageResponse)
	if !ok {
		return nil, ErrFailedToUnmarshalResponse
	}

	if res.Response.Code != "200" {
		return nil, errors.New(res.Response.Message)
	}

	return res, nil
}

type UpdateProjectLanguageRequest struct {
	ID   int
	Data []struct{}
}

type UpdateProjectLanguageResponse struct {
	Response struct {
		Status  string `json:"status"`
		Code    string `json:"code"`
		Message string `json:"message"`
	} `json:"response"`
	Result struct {
		Translations struct {
			Parsed  int64 `json:"parsed"`
			Added   int64 `json:"added"`
			Updated int64 `json:"updated"`
		} `json:"translations"`
	} `json:"result"`
}

// UpdateProjectLanguage Inserts / overwrites translations.
//
// https://poeditor.com/docs/api#languages_update
//
// NOT IMPLEMENTED in this sdk
func (c *ClientImpl) UpdateProjectLanguage(ctx context.Context, r UpdateProjectLanguageRequest) (*UpdateProjectLanguageResponse, error) {
	return nil, ErrNotImplemented
}

type DeleteProjectLanguageRequest struct {
	ID       int
	Language string
}

type DeleteProjectLanguageResponse struct {
	Response struct {
		Status  string `json:"status"`
		Code    string `json:"code"`
		Message string `json:"message"`
	} `json:"response"`
}

func (c *ClientImpl) DeleteProjectLanguage(ctx context.Context, r DeleteProjectLanguageRequest) (*DeleteProjectLanguageResponse, error) {
	req := c.r.R()

	req.SetFormData(map[string]string{
		"id":       fmt.Sprintf("%d", r.ID),
		"language": r.Language,
	})

	req.SetContext(ctx)

	req.SetResult(&DeleteProjectLanguageResponse{})

	resp, err := req.Post("/v2/languages/delete")
	if err != nil {
		return nil, err
	}

	res, ok := resp.Result().(*DeleteProjectLanguageResponse)
	if !ok {
		return nil, ErrFailedToUnmarshalResponse
	}

	if res.Response.Code != "200" {
		return nil, errors.New(res.Response.Message)
	}

	return res, nil
}
