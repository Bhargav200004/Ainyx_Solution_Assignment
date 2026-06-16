package models

import (
	"fmt"
	"strings"
	"time"

	"github.com/go-playground/validator/v10"
)

// validate is the singleton validator instance.
var validate *validator.Validate

func init() {
	validate = validator.New()

	// Register custom "dateformat" tag for YYYY-MM-DD validation.
	_ = validate.RegisterValidation("dateformat", func(fl validator.FieldLevel) bool {
		val := fl.Field().String()
		if val == "" {
			return false
		}
		parsed, err := time.Parse(DOBLayout, val)
		if err != nil {
			return false
		}
		// Ensure the date is not in the future.
		if parsed.After(time.Now()) {
			return false
		}
		return true
	})
}

// Validate validates a struct against its validation tags.
func Validate(s interface{}) error {
	return validate.Struct(s)
}

// FormatValidationErrors converts validator.ValidationErrors into a
// human-readable map of field → message.
func FormatValidationErrors(err error) map[string]string {
	errors := make(map[string]string)

	if validationErrors, ok := err.(validator.ValidationErrors); ok {
		for _, e := range validationErrors {
			field := strings.ToLower(e.Field())
			switch e.Tag() {
			case "required":
				errors[field] = fmt.Sprintf("%s is required", field)
			case "min":
				errors[field] = fmt.Sprintf("%s must be at least %s characters", field, e.Param())
			case "max":
				errors[field] = fmt.Sprintf("%s must be at most %s characters", field, e.Param())
			case "dateformat":
				errors[field] = fmt.Sprintf("%s must be a valid date in YYYY-MM-DD format and not in the future", field)
			default:
				errors[field] = fmt.Sprintf("%s is invalid", field)
			}
		}
	}

	return errors
}
