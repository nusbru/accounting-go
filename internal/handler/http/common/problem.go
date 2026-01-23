package common

import "fmt"

// ProblemDetail represents RFC 7807 Problem Details for HTTP APIs
type ProblemDetail struct {
	Type     string `json:"type"`
	Title    string `json:"title"`
	Status   int    `json:"status"`
	Detail   string `json:"detail,omitempty"`
	Instance string `json:"instance,omitempty"`
}

// Problem type URIs
const (
	TypeValidationError  = "https://api.accounting.app/problems/validation-error"
	TypeNotFound         = "https://api.accounting.app/problems/not-found"
	TypeInternalError    = "https://api.accounting.app/problems/internal-error"
	TypeMethodNotAllowed = "https://api.accounting.app/problems/method-not-allowed"
	TypeBadRequest       = "https://api.accounting.app/problems/bad-request"
)

// NewValidationProblem creates a validation error problem detail
func NewValidationProblem(detail, instance string) *ProblemDetail {
	return &ProblemDetail{
		Type:     TypeValidationError,
		Title:    "Validation Error",
		Status:   400,
		Detail:   detail,
		Instance: instance,
	}
}

// NewNotFoundProblem creates a not found problem detail
func NewNotFoundProblem(detail, instance string) *ProblemDetail {
	return &ProblemDetail{
		Type:     TypeNotFound,
		Title:    "Not Found",
		Status:   404,
		Detail:   detail,
		Instance: instance,
	}
}

// NewInternalErrorProblem creates an internal error problem detail
func NewInternalErrorProblem(instance string) *ProblemDetail {
	return &ProblemDetail{
		Type:     TypeInternalError,
		Title:    "Internal Server Error",
		Status:   500,
		Detail:   "An unexpected error occurred. Please try again later.",
		Instance: instance,
	}
}

// NewMethodNotAllowedProblem creates a method not allowed problem detail
func NewMethodNotAllowedProblem(instance string) *ProblemDetail {
	return &ProblemDetail{
		Type:     TypeMethodNotAllowed,
		Title:    "Method Not Allowed",
		Status:   405,
		Detail:   "The HTTP method is not allowed for this endpoint.",
		Instance: instance,
	}
}

// NewBadRequestProblem creates a bad request problem detail
func NewBadRequestProblem(detail, instance string) *ProblemDetail {
	return &ProblemDetail{
		Type:     TypeBadRequest,
		Title:    "Bad Request",
		Status:   400,
		Detail:   detail,
		Instance: instance,
	}
}

// ValidationError represents a field-level validation error
type ValidationError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

// ValidationProblem extends ProblemDetail with validation errors
type ValidationProblem struct {
	ProblemDetail
	Errors []ValidationError `json:"errors,omitempty"`
}

// NewValidationProblemWithErrors creates a validation problem with field errors
func NewValidationProblemWithErrors(instance string, errors []ValidationError) *ValidationProblem {
	detail := "One or more validation errors occurred."
	if len(errors) == 1 {
		detail = fmt.Sprintf("Validation failed: %s", errors[0].Message)
	}

	return &ValidationProblem{
		ProblemDetail: ProblemDetail{
			Type:     TypeValidationError,
			Title:    "Validation Error",
			Status:   400,
			Detail:   detail,
			Instance: instance,
		},
		Errors: errors,
	}
}
