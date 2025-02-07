package errors

import (
	"fmt"
)

// Custom error types
type ErrorType string

const (
	ErrorTypeValidation ErrorType = "VALIDATION_ERROR"
	ErrorTypeEthereum   ErrorType = "ETHEREUM_ERROR"
	ErrorTypeStorage    ErrorType = "STORAGE_ERROR"
	ErrorTypeUnexpected ErrorType = "UNEXPECTED_ERROR"
)

// AppError represents an application-specific error
type AppError struct {
	Type    ErrorType
	Message string
	Err     error
	Meta    map[string]interface{}
}

func (e *AppError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("%s: %s: %v", e.Type, e.Message, e.Err)
	}
	return fmt.Sprintf("%s: %s", e.Type, e.Message)
}

// Error constructors
func NewValidationError(message string, err error) *AppError {
	return &AppError{
		Type:    ErrorTypeValidation,
		Message: message,
		Err:     err,
	}
}

func NewEthereumError(message string, err error) *AppError {
	return &AppError{
		Type:    ErrorTypeEthereum,
		Message: message,
		Err:     err,
	}
}

func NewStorageError(message string, err error) *AppError {
	return &AppError{
		Type:    ErrorTypeStorage,
		Message: message,
		Err:     err,
	}
}

func NewUnexpectedError(message string, err error) *AppError {
	return &AppError{
		Type:    ErrorTypeUnexpected,
		Message: message,
		Err:     err,
	}
}
