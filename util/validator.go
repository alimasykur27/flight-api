package util

import (
	"flight-api/internal/enum"

	"github.com/go-playground/validator"
)

func NewValidator() *validator.Validate {
	var validate *validator.Validate = validator.New()

	//Register custom validation for FacilityTYpe
	validate.RegisterValidation("facility", func(fl validator.FieldLevel) bool {
		val := fl.Field().String()
		return val == "" || val == string(*enum.AIRPORT) || val == string(*enum.HELIPORT)
	})

	//Register custom validation for Ownership
	validate.RegisterValidation("ownership", func(fl validator.FieldLevel) bool {
		val := fl.Field().String()
		return val == "" || val == string(*enum.OWN_PUBLIC) || val == string(*enum.OWN_PRIVATE)
	})

	//Register custom validation for Use
	validate.RegisterValidation("use", func(fl validator.FieldLevel) bool {
		val := fl.Field().String()
		return val == "" || val == string(*enum.USE_PUBLIC) || val == string(*enum.USE_PRIVATE)
	})

	return validate
}
