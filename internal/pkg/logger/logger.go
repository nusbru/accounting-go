package logger

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"os"
	"strings"
)

// Logger wraps slog.Logger with convenience methods and request tracking
type Logger struct {
	*slog.Logger
}

// New creates a new Logger instance
// logFormat can be "json" or "text"
// logLevel can be "debug", "info", "warn", "error"
func New(logFormat, logLevel string) *Logger {
	var handler slog.Handler

	// Parse log level
	level := slog.LevelInfo
	switch strings.ToLower(logLevel) {
	case "debug":
		level = slog.LevelDebug
	case "info":
		level = slog.LevelInfo
	case "warn":
		level = slog.LevelWarn
	case "error":
		level = slog.LevelError
	}

	opts := &slog.HandlerOptions{
		Level: level,
	}

	// Choose handler based on format
	if strings.ToLower(logFormat) == "json" {
		handler = slog.NewJSONHandler(os.Stderr, opts)
	} else {
		handler = slog.NewTextHandler(os.Stderr, opts)
	}

	slogger := slog.New(handler)
	return &Logger{Logger: slogger}
}

// NewWithWriter creates a new Logger with a custom writer
func NewWithWriter(w io.Writer, logFormat, logLevel string) *Logger {
	var handler slog.Handler

	// Parse log level
	level := slog.LevelInfo
	switch strings.ToLower(logLevel) {
	case "debug":
		level = slog.LevelDebug
	case "info":
		level = slog.LevelInfo
	case "warn":
		level = slog.LevelWarn
	case "error":
		level = slog.LevelError
	}

	opts := &slog.HandlerOptions{
		Level: level,
	}

	// Choose handler based on format
	if strings.ToLower(logFormat) == "json" {
		handler = slog.NewJSONHandler(w, opts)
	} else {
		handler = slog.NewTextHandler(w, opts)
	}

	slogger := slog.New(handler)
	return &Logger{Logger: slogger}
}

// WithContext adds structured context to log messages
func (l *Logger) WithContext(ctx context.Context) *Logger {
	// Extract request ID from context if present
	// Uses the middleware.ContextKey type for type-safe context keys
	type contextKey string
	const requestIDKey contextKey = "request-id"

	reqID := ctx.Value(requestIDKey)
	if reqID != nil {
		return &Logger{Logger: l.Logger.With("request_id", fmt.Sprintf("%v", reqID))}
	}
	return l
}

// Debugf logs a debug message with format string
func (l *Logger) Debugf(msg string, args ...any) {
	l.Debug(fmt.Sprintf(msg, args...))
}

// Infof logs an info message with format string
func (l *Logger) Infof(msg string, args ...any) {
	l.Info(fmt.Sprintf(msg, args...))
}

// Warnf logs a warn message with format string
func (l *Logger) Warnf(msg string, args ...any) {
	l.Warn(fmt.Sprintf(msg, args...))
}

// Errorf logs an error message with format string
func (l *Logger) Errorf(msg string, args ...any) {
	l.Error(fmt.Sprintf(msg, args...))
}

// WithError adds an error to the context
func (l *Logger) WithError(err error) *Logger {
	return &Logger{Logger: l.Logger.With("error", err.Error())}
}

// WithField adds a single field to the context
func (l *Logger) WithField(key string, value any) *Logger {
	return &Logger{Logger: l.Logger.With(key, value)}
}

// WithFields adds multiple fields to the context
func (l *Logger) WithFields(fields map[string]any) *Logger {
	var args []any
	for k, v := range fields {
		args = append(args, k, v)
	}
	return &Logger{Logger: l.Logger.With(args...)}
}
