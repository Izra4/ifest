package helpers

import "github.com/go-playground/validator/v10"

func FormatValidationError(err validator.FieldError) string {
	switch err.Field() {
	case "Name":
		if err.Tag() == "required" {
			return "Name is required"
		} else if err.Tag() == "min" || err.Tag() == "max" {
			return "Name must be between 2 and 100 characters"
		}
	case "Email":
		if err.Tag() == "required" {
			return "Email is required"
		} else if err.Tag() == "email" {
			return "Email must be a valid email address"
		}
	case "Password":
		if err.Tag() == "required" {
			return "Password is required"
		} else if err.Tag() == "min" {
			return "Password must be at least 6 characters long"
		}
	default:
		return "Invalid value for " + err.Field()
	}
	return ""
}
