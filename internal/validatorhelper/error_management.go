package validatorhelper

import (
	"fmt"

	"github.com/go-playground/validator/v10"
)

func ProcessErrors(validationErrors validator.ValidationErrors) string {
	errorMessage := ""
	for _, err := range validationErrors {
		errorMessage += fmt.Sprintf("Error in field '%s'. ", err.Field())
		switch err.Tag() {
		case "required":
			errorMessage += "Is required"
		case "email":
			errorMessage += "Must be a valid email"
			// TODO: Add all tags. For the moment we only use required and email.
		}
		errorMessage += "\n"
	}
	return errorMessage
}
