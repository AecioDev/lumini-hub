package validator

import (
	"errors"
)

// ValidationError representa um erro de validação
type ValidationError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
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

// FieldMessages devolve a mensagem de cada campo com erro (o primeiro erro de cada campo), pelo nome do campo no
// JSON. Implementa utils.FieldMessager: a resposta de erro da API leva esse mapa em `fieldErrors`.
func (ve ValidationErrors) FieldMessages() map[string]string {
	messages := map[string]string{}
	for _, e := range ve {
		if e.Field == "" {
			continue
		}
		if _, exists := messages[e.Field]; !exists {
			messages[e.Field] = e.Message
		}
	}
	return messages
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
