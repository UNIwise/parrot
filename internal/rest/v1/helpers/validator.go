package helpers

import (
	"regexp"

	"github.com/go-playground/validator/v10"
	"github.com/johngb/langreg"
)

type Validator struct {
	validator *validator.Validate
}

func NewValidator() (*Validator, error) {
	v := validator.New()

	if err := v.RegisterValidation("languageCode", validateLanguageCode); err != nil {
		return nil, err
	}

	if err := v.RegisterValidation("version", validateVersion); err != nil {
		return nil, err
	}

	return &Validator{
		validator: v,
	}, nil
}

func (cv *Validator) Validate(i interface{}) error {
	return cv.validator.Struct(i)
}

func validateLanguageCode(fl validator.FieldLevel) bool {
	return langreg.IsValidLanguageCode(fl.Field().String())
}

func validateVersion(fl validator.FieldLevel) bool {
	r := regexp.MustCompile(`^[a-zA-Z0-9.-_]*$`)

	return r.MatchString(fl.Field().String())
}
