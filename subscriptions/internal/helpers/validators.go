package helpers

import (
	_ "regexp"
	_ "strconv"
	"time"

	"github.com/go-playground/validator/v10"
)

func ValidateDateMMYYYYFormatValidator(fl validator.FieldLevel) bool {
	return ValidateDateMMYYYYFormat(fl.Field().String())
}

func ValidateDateMMYYYYFormat(date string) bool {
	_, err := time.Parse("01-2006", date)
	return err == nil
}

