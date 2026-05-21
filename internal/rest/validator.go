package rest

import (
	"github.com/go-playground/validator"
	"github.com/johngb/langreg"
)

type Validator struct {
	validator *validator.Validate
}

func NewValidator() *Validator {
	v := validator.New()

	_ = v.RegisterValidation("languageCode", validateLanguageCode) //nolint:errcheck

	return &Validator{
		validator: v,
	}
}

func (cv *Validator) Validate(i any) error {
	return cv.validator.Struct(i)
}

func validateLanguageCode(fl validator.FieldLevel) bool {
	return langreg.IsValidLanguageCode(fl.Field().String())
}
