// Package validator provides custom validation functions without external dependencies.
package validator

import (
	"fmt"
	"regexp"
	"strings"
)

var (
	// emailRegex validates basic email format (RFC 5321 simplified)
	emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)

	// uuidRegex validates UUID v4 format
	uuidRegex = regexp.MustCompile(`^[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}$`)
)

// ValidationError represents a validation failure for a specific field.
type ValidationError struct {
	Field   string
	Message string
}

func (e *ValidationError) Error() string {
	return fmt.Sprintf("%s: %s", e.Field, e.Message)
}

// NewValidationError creates a new ValidationError.
func NewValidationError(field, message string) *ValidationError {
	return &ValidationError{Field: field, Message: message}
}

// ValidateEmail checks if the value is a valid email address.
func ValidateEmail(value string) bool {
	return emailRegex.MatchString(value)
}

// ValidateUUID checks if the value is a valid UUID.
func ValidateUUID(value string) bool {
	return uuidRegex.MatchString(value)
}

// ValidateRequired checks if the value is not empty (after trimming whitespace).
// Returns nil if valid, or a ValidationError if invalid.
func ValidateRequired(value, field string) error {
	if strings.TrimSpace(value) == "" {
		return NewValidationError(field, "is required")
	}
	return nil
}

// ValidateLength checks if the value length is within the specified range.
// Returns nil if valid, or a ValidationError if invalid.
func ValidateLength(value, field string, min, max int) error {
	l := len(value)
	if l < min {
		return NewValidationError(field, fmt.Sprintf("must be at least %d characters", min))
	}
	if l > max {
		return NewValidationError(field, fmt.Sprintf("must be at most %d characters", max))
	}
	return nil
}

// ValidateEmailFormat checks if the value is a valid email format.
// Returns nil if valid, or a ValidationError if invalid.
func ValidateEmailFormat(value, field string) error {
	if !ValidateEmail(value) {
		return NewValidationError(field, "must be a valid email address")
	}
	return nil
}

// ValidateUUIDFormat checks if the value is a valid UUID format.
// Returns nil if valid, or a ValidationError if invalid.
func ValidateUUIDFormat(value, field string) error {
	if !ValidateUUID(value) {
		return NewValidationError(field, "must be a valid UUID")
	}
	return nil
}

// ValidatePositive checks if the value is positive.
// Returns nil if valid, or a ValidationError if invalid.
func ValidatePositive(value float64, field string) error {
	if value <= 0 {
		return NewValidationError(field, "must be greater than zero")
	}
	return nil
}

// ValidateNonNegative checks if the value is non-negative.
// Returns nil if valid, or a ValidationError if invalid.
func ValidateNonNegative(value float64, field string) error {
	if value < 0 {
		return NewValidationError(field, "must not be negative")
	}
	return nil
}

// ValidateOneOf checks if the value is one of the allowed values.
// Returns nil if valid, or a ValidationError if invalid.
func ValidateOneOf(value, field string, allowed []string) error {
	for _, a := range allowed {
		if value == a {
			return nil
		}
	}
	return NewValidationError(field, fmt.Sprintf("must be one of: %s", strings.Join(allowed, ", ")))
}
