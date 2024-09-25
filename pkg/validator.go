package pkg

import (
	validate "github.com/go-playground/validator/v10"
)

type CustomValidator struct {
	Validator *validate.Validate
}

func NewValidator() *CustomValidator {
	// Create a new CustomValidator instance
	return &CustomValidator{
		Validator: validate.New(),
	}
}

func (cv *CustomValidator) Validate(i interface{}) error {
	if err := cv.Validator.Struct(i); err != nil {
		return err
	}

	return nil
}
