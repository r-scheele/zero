package services

import (
	"regexp"
	"github.com/go-playground/validator/v10"
)

// Validator provides validation mainly validating structs within the web context
type Validator struct {
	// validator stores the underlying validator
	validator *validator.Validate
}

// NewValidator creats a new Validator
func NewValidator() *Validator {
	v := &Validator{
		validator: validator.New(),
	}
	
	// Register custom validation for E.164 phone number format
	v.validator.RegisterValidation("e164", validateE164)
	
	return v
}

// validateE164 validates phone numbers in E.164 format
func validateE164(fl validator.FieldLevel) bool {
	phoneNumber := fl.Field().String()
	if phoneNumber == "" {
		return true // Allow empty values, use 'required' tag for mandatory fields
	}
	
	// E.164 format: + followed by 1-15 digits
	matched, _ := regexp.MatchString(`^\+[1-9]\d{1,14}$`, phoneNumber)
	return matched
}

// Validate validates a struct
func (v *Validator) Validate(i any) error {
	if err := v.validator.Struct(i); err != nil {
		return err
	}
	return nil
}
