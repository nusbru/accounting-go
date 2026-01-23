// Package middleware provides HTTP middleware for the application.
package middleware

// ContextKey is a type for context keys to avoid collisions.
type ContextKey string

const (
	// RequestIDKey is the context key for request ID.
	RequestIDKey ContextKey = "request-id"
)
