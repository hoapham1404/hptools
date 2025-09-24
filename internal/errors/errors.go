package errors

import "fmt"

// ErrorType represents the type of error
type ErrorType string

const (
	// ErrorTypeProcess represents process-related errors
	ErrorTypeProcess ErrorType = "process"
	// ErrorTypeWindow represents window-related errors
	ErrorTypeWindow ErrorType = "window"
	// ErrorTypeAPI represents Windows API errors
	ErrorTypeAPI ErrorType = "api"
	// ErrorTypeConfig represents configuration errors
	ErrorTypeConfig ErrorType = "config"
)

// AppError represents a structured application error
type AppError struct {
	Type    ErrorType `json:"type"`
	Message string    `json:"message"`
	Cause   error     `json:"cause,omitempty"`
}

// Error implements the error interface
func (e *AppError) Error() string {
	if e.Cause != nil {
		return fmt.Sprintf("%s error: %s (caused by: %v)", e.Type, e.Message, e.Cause)
	}
	return fmt.Sprintf("%s error: %s", e.Type, e.Message)
}

// Unwrap returns the underlying error
func (e *AppError) Unwrap() error {
	return e.Cause
}

// NewProcessError creates a new process-related error
func NewProcessError(message string, cause error) *AppError {
	return &AppError{
		Type:    ErrorTypeProcess,
		Message: message,
		Cause:   cause,
	}
}

// NewWindowError creates a new window-related error
func NewWindowError(message string, cause error) *AppError {
	return &AppError{
		Type:    ErrorTypeWindow,
		Message: message,
		Cause:   cause,
	}
}

// NewAPIError creates a new Windows API error
func NewAPIError(message string, cause error) *AppError {
	return &AppError{
		Type:    ErrorTypeAPI,
		Message: message,
		Cause:   cause,
	}
}

// NewConfigError creates a new configuration error
func NewConfigError(message string, cause error) *AppError {
	return &AppError{
		Type:    ErrorTypeConfig,
		Message: message,
		Cause:   cause,
	}
}
