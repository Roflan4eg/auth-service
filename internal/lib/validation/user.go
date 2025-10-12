package validation

import (
	"fmt"
	"github.com/go-playground/validator/v10"
	"regexp"
)

var validate *validator.Validate

func validateStrongPassword(fl validator.FieldLevel) bool {
	password := fl.Field().String()

	if len(password) < 8 {
		return false
	}

	hasUpper := regexp.MustCompile(`[A-Z]`).MatchString(password)
	hasDigit := regexp.MustCompile(`[0-9]`).MatchString(password)
	hasSpecial := regexp.MustCompile(`[!@#$%^&*()_+\-=\[\]{};':"\\|,.<>\/?]`).MatchString(password)

	return hasUpper && hasDigit && hasSpecial
}

func init() {
	validate = validator.New()
	err := validate.RegisterValidation("strongPassword", validateStrongPassword)
	if err != nil {
		panic(fmt.Sprintf("failed to register validation: %v", err))
	}
}
func ValidateStruct(s interface{}) error {
	return validate.Struct(s)
}
