package helpers

import (
	"regexp"

	"github.com/go-playground/validator/v10"
)

func ValidateDateMMYYYYFormatValidator(fl validator.FieldLevel) bool {
	return ValidateDateMMYYYYFormat(fl.Field().String())
}

func ValidateDateMMYYYYFormat(date string) bool {
	// TODO: Check month
	var dateRegex = regexp.MustCompile(`^\d{2}-\d{4}$`)

	return dateRegex.MatchString(date)
}