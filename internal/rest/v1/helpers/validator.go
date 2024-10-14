package helpers

import (
	"github.com/go-playground/validator/v10"
	"github.com/johngb/langreg"
)

type Validator struct {
	validator *validator.Validate
}

func NewValidator() *Validator {
	v := validator.New()

	v.RegisterValidation("languageCode", validateLanguageCode)

	return &Validator{
		validator: v,
	}
}

func (cv *Validator) Validate(i interface{}) error {
	return cv.validator.Struct(i)
}

func validateLanguageCode(fl validator.FieldLevel) bool {
	return langreg.IsValidLanguageCode(fl.Field().String())
}
