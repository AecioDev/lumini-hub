package validator

import (
	"errors"
)

// ValidationError representa um erro de validação
type ValidationError struct {
	Field   string
	Message string
}

// ValidationErrors representa uma coleção de erros de validação
type ValidationErrors []ValidationError

// Error implementa a interface error
func (ve ValidationErrors) Error() string {
	if len(ve) == 0 {
		return ""
	}
	return ve[0].Message
}

// GetErrors retorna todos os erros de validação
func (ve ValidationErrors) GetErrors() []ValidationError {
	return ve
}

// AddError adiciona um erro de validação
func (ve *ValidationErrors) AddError(field, message string) {
	*ve = append(*ve, ValidationError{Field: field, Message: message})
}

// HasErrors verifica se existem erros de validação
func (ve ValidationErrors) HasErrors() bool {
	return len(ve) > 0
}

// IsValidationError verifica se um erro é um erro de validação
func IsValidationError(err error) bool {
	var ve ValidationErrors
	return errors.As(err, &ve)
}