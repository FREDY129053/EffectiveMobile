package helpers

import (
	"regexp"

	"github.com/go-playground/validator/v10"
)

func ValidateDateMMYYYYFormat(fl validator.FieldLevel) bool {
	date := fl.Field().String()

	// TODO: Check month
	var dateRegex = regexp.MustCompile(`^\d{2}-\d{4}$`)

	return dateRegex.MatchString(date)
}