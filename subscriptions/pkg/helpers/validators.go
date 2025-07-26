package helpers

import (
	"regexp"
	"strconv"

	"github.com/go-playground/validator/v10"
)

func ValidateDateMMYYYYFormatValidator(fl validator.FieldLevel) bool {
	return ValidateDateMMYYYYFormat(fl.Field().String())
}

func ValidateDateMMYYYYFormat(date string) bool {
	var dateRegex = regexp.MustCompile(`^\d{2}-\d{4}$`)
	
	monthInt, _ := strconv.ParseInt(date[:2], 10, 64)
	isRightMonth := true
	if monthInt <= 0 || monthInt > 12 {
		isRightMonth = false
	}

	return dateRegex.MatchString(date) && isRightMonth
}