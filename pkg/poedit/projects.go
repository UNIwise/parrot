package poedit

import (
	"context"
	"fmt"
	"strings"

	"github.com/pkg/errors"
)

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
