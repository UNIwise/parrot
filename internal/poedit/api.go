package poedit

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/pkg/errors"
)

// Api client

// ExportProjectRequest is the request model for the poeditor project export endpoint
// https://poeditor.com/docs/api#projects_export
type ExportProjectRequest struct {
	APIToken string   `json:"api_token"`
	ID       int      `json:"id"`
	Language string   `json:"language"`
	Type     string   `json:"type"`
	Filters  []string `json:"filters,omitempty"`
	Order    string   `json:"order,omitempty"`
	Tags     []string `json:"tags,omitempty"`
	Options  []string `json:"options,omitempty"`
}

// ExportProjectResponse is the response model for the poeditor project export endpoint
// https://poeditor.com/docs/api#projects_export
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

// ExportLanguage calls the export language endpoint in poeditor
// https://poeditor.com/docs/api#projects_export
func ExportLanguage(h *http.Client, request ExportProjectRequest) (*ExportProjectResponse, error) {
	body, err := json.Marshal(request)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to marshal request body")
	}

	resp, err := h.Post("https://api.poeditor.com/v2/projects/export", "application/json", bytes.NewReader(body))
	if err != nil {
		return nil, errors.Wrap(err, "Failed to request project export")
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		b, _ := ioutil.ReadAll(resp.Body)

		return nil, errors.Errorf("Project export failed. status='%d' body='%s'", resp.StatusCode, string(b))
	}
	decoder := json.NewDecoder(resp.Body)

	response := &ExportProjectResponse{}
	err = decoder.Decode(response)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to decode response body from poeditor")
	}

	return response, nil
}
