package common

import (
	"regexp"
	"strconv"
	"strings"

	"github.com/google/uuid"
)

var emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)

// ValidateRequired checks if a string value is not empty
func ValidateRequired(value, fieldName string) *ValidationError {
	if strings.TrimSpace(value) == "" {
		return &ValidationError{
			Field:   fieldName,
			Message: fieldName + " is required",
		}
	}
	return nil
}

// ValidateEmail checks if a string is a valid email format
func ValidateEmail(email, fieldName string) *ValidationError {
	if email == "" {
		return &ValidationError{
			Field:   fieldName,
			Message: fieldName + " is required",
		}
	}
	if !emailRegex.MatchString(email) {
		return &ValidationError{
			Field:   fieldName,
			Message: fieldName + " must be a valid email address",
		}
	}
	return nil
}

// ValidateStringLength checks if a string length is within bounds
func ValidateStringLength(value, fieldName string, min, max int) *ValidationError {
	length := len(strings.TrimSpace(value))
	if length < min {
		return &ValidationError{
			Field:   fieldName,
			Message: fieldName + " must be at least " + strconv.Itoa(min) + " characters",
		}
	}
	if length > max {
		return &ValidationError{
			Field:   fieldName,
			Message: fieldName + " must be at most " + strconv.Itoa(max) + " characters",
		}
	}
	return nil
}

// ValidateUUID checks if a string is a valid UUID
func ValidateUUID(value, fieldName string) *ValidationError {
	if value == "" {
		return &ValidationError{
			Field:   fieldName,
			Message: fieldName + " is required",
		}
	}
	if _, err := uuid.Parse(value); err != nil {
		return &ValidationError{
			Field:   fieldName,
			Message: fieldName + " must be a valid UUID",
		}
	}
	return nil
}

// ValidateCurrency checks if a string is a valid 3-character currency code
func ValidateCurrency(currency, fieldName string) *ValidationError {
	if currency == "" {
		return &ValidationError{
			Field:   fieldName,
			Message: fieldName + " is required",
		}
	}
	if len(currency) != 3 {
		return &ValidationError{
			Field:   fieldName,
			Message: fieldName + " must be a 3-character ISO currency code",
		}
	}
	// Check if all characters are uppercase letters
	for _, c := range currency {
		if c < 'A' || c > 'Z' {
			return &ValidationError{
				Field:   fieldName,
				Message: fieldName + " must contain only uppercase letters",
			}
		}
	}
	return nil
}

// ValidatePositive checks if a float value is greater than zero
func ValidatePositive(value float64, fieldName string) *ValidationError {
	if value <= 0 {
		return &ValidationError{
			Field:   fieldName,
			Message: fieldName + " must be greater than zero",
		}
	}
	return nil
}

// ValidateEnum checks if a value is in a list of allowed values
func ValidateEnum(value string, allowed []string, fieldName string) *ValidationError {
	if value == "" {
		return &ValidationError{
			Field:   fieldName,
			Message: fieldName + " is required",
		}
	}
	for _, a := range allowed {
		if value == a {
			return nil
		}
	}
	return &ValidationError{
		Field:   fieldName,
		Message: fieldName + " must be one of: " + strings.Join(allowed, ", "),
	}
}

// CollectErrors collects non-nil validation errors
func CollectErrors(errors ...*ValidationError) []ValidationError {
	var result []ValidationError
	for _, err := range errors {
		if err != nil {
			result = append(result, *err)
		}
	}
	return result
}
