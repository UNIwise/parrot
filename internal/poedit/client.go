package poedit

import (
	"context"
	"io"
	"io/ioutil"
	"net/http"

	"github.com/pkg/errors"
)

// Client is an interface to poeditors api
type Client interface {
	GetLanguage(ctx context.Context, projectID int, language string) (result string, err error)
}

// ClientImpl is an implementation of the poeditor client interface
type ClientImpl struct {
	APIToken   string
	HTTPClient *http.Client
}

// NewClient creates a new poeditor api client
func NewClient(apiToken string, httpClient *http.Client) *ClientImpl {
	return &ClientImpl{
		APIToken:   apiToken,
		HTTPClient: httpClient,
	}
}

// GetLanguage exports and downloads a language for a project in poeditor
func (c *ClientImpl) GetLanguage(ctx context.Context, projectID int, language string) (io.ReadCloser, error) {
	exportResp, err := ExportLanguage(c.HTTPClient, ExportProjectRequest{
		APIToken: c.APIToken,
		ID:       projectID,
		Language: language,
		Type:     "key_value_json",
	})
	if err != nil {
		return nil, errors.Wrapf(err, "Error exporting language '%s' from project '%d'", language, projectID)
	}

	resp, err := c.HTTPClient.Get(exportResp.Result.URL)
	if err != nil {
		return nil, errors.Wrapf(err, "Error requesting terms from '%s'", exportResp.Result.URL)
	}
	if resp.StatusCode != http.StatusOK {
		b, _ := ioutil.ReadAll(resp.Body)
		resp.Body.Close()
		return nil, errors.Errorf("Unexpected response from poeditor '%d' '%s'", resp.StatusCode, string(b))
	}
	return resp.Body, nil
}
