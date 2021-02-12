package rest

import "github.com/go-playground/validator"

type Validator struct {
	validator *validator.Validate
}

func NewValidator() *Validator {
	return &Validator{
		validator: validator.New(),
	}
}

func (cv *Validator) Validate(i interface{}) error {
	return cv.validator.Struct(i)
}
