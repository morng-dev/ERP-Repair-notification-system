package utils

import (
	"regexp"

	"github.com/go-playground/validator/v10"
)

var validate = validator.New()

func ValidateStruct(s interface{}) error {
	return validate.Struct(s)
}

func IsEmail(email string) bool {
	match, err := regexp.MatchString(`^[\w-\.]+@([\w-]+\.)+[\w-]{2,4}$`, email)
	if err != nil {
		return false
	}
	return match
}
