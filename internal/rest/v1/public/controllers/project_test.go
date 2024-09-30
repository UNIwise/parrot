package controllers

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"regexp"
	"testing"
	"time"

	"github.com/go-playground/validator"
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/uniwise/parrot/internal/project"

	gomock "go.uber.org/mock/gomock"
)

var (
	errTest                     = errors.New("test error")
	testID               int64  = 1
	testLanguage         string = "en"
	testFormat           string = "key_value_json"
	testVersion          string = "latest"
)

type CustomValidator struct {
	validator *validator.Validate
}

// Validate implements the echo.Validator interface
func (cv *CustomValidator) Validate(i interface{}) error {
	return cv.validator.Struct(i)
}

// Custom validation function to check if the language code is valid
func languageCodeValidator(fl validator.FieldLevel) bool {
	// Define a regex for language code validation (e.g., 'en', 'fr', 'es')
	languageCodePattern := `^[a-z]{2}$`
	matched, _ := regexp.MatchString(languageCodePattern, fl.Field().String())
	return matched
}

func TestGetProjectLanguage(t *testing.T) {
	t.Parallel()

	e := echo.New()
	validator := validator.New()
	validator.RegisterValidation("languageCode", languageCodeValidator)
	e.Validator = &CustomValidator{validator: validator}

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	resp := httptest.NewRecorder()

	testCtx := e.NewContext(req, resp)

	testCtx.SetPath("/projects/:project/language/:language?format=:format&version=:version")
	testCtx.SetParamNames("project", "language", "format", "version")
	testCtx.SetParamValues(fmt.Sprintf("%d", testID), testLanguage, testFormat, testVersion)

	projectService := project.NewMockService(gomock.NewController(t))

	h := &Handlers{
		ProjectService: projectService,
	}

	t.Run("getProjectLanguage, success", func(t *testing.T) {

		projectService.EXPECT().GetTranslation(context.Background(), int(testID), testLanguage, testFormat, testVersion).Times(1).Return(&project.Translation{
			Data:     []byte("test"),
			Checksum: "test",
			TTL:      time.Second,
		}, nil)

		err := h.getProjectLanguage(testCtx, logrus.NewEntry(logrus.New()))
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.Code)
	})

	t.Run("getProjectLanguage, fail", func(t *testing.T) {

		projectService.EXPECT().GetTranslation(context.Background(), int(testID), testLanguage, testFormat, testVersion).Times(1).Return(nil, errTest)

		err := h.getProjectLanguage(testCtx, logrus.NewEntry(logrus.New()))
		assert.Error(t, err)

	})
}
