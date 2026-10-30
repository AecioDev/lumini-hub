package utils

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"reflect"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

// O validador do Gin chama o campo pelo nome em Go ("FirstName"). O cliente da API conhece o campo
// pelo nome do JSON ("first_name"), então registramos o nome do JSON como o nome do campo nos erros.
// Roda no init, antes de qualquer validação (o validador guarda em cache o que já viu).
func init() {
	if engine, ok := binding.Validator.Engine().(*validator.Validate); ok {
		engine.RegisterTagNameFunc(func(field reflect.StructField) string {
			name := strings.SplitN(field.Tag.Get("json"), ",", 2)[0]
			if name == "" || name == "-" {
				return field.Name
			}
			return name
		})
	}
}

// describeRequestError traduz o erro de um bind/validação de requisição em frases legíveis, sem
// mostrar nomes internos do código (struct, campo em Go, regra do validador). Devolve false quando o
// erro não é desse tipo (por exemplo, uma regra de negócio nossa), para o chamador usar o texto dele.
func describeRequestError(err error) ([]string, bool) {
	var fieldErrors validator.ValidationErrors
	if errors.As(err, &fieldErrors) {
		messages := make([]string, 0, len(fieldErrors))
		for _, fieldError := range fieldErrors {
			messages = append(messages, fmt.Sprintf("%s: %s", fieldError.Field(), validationRuleMessage(fieldError)))
		}
		return messages, true
	}

	var typeError *json.UnmarshalTypeError
	if errors.As(err, &typeError) {
		if typeError.Field != "" {
			return []string{fmt.Sprintf("%s: tipo de valor inválido", typeError.Field)}, true
		}
		return []string{"tipo de valor inválido"}, true
	}

	var syntaxError *json.SyntaxError
	if errors.As(err, &syntaxError) || errors.Is(err, io.ErrUnexpectedEOF) {
		return []string{"o corpo da requisição não é um JSON válido"}, true
	}
	if errors.Is(err, io.EOF) {
		return []string{"o corpo da requisição está vazio"}, true
	}

	var numError *strconv.NumError
	if errors.As(err, &numError) {
		return []string{"valor numérico inválido"}, true
	}

	return nil, false
}

// validationRuleMessage escreve, em português, o que a regra do validador exige.
func validationRuleMessage(fieldError validator.FieldError) string {
	limit := fieldError.Param()
	isText := fieldError.Kind() == reflect.String
	isList := fieldError.Kind() == reflect.Slice || fieldError.Kind() == reflect.Array || fieldError.Kind() == reflect.Map

	switch fieldError.Tag() {
	case "required":
		return "campo obrigatório"
	case "max":
		switch {
		case isText:
			return "no máximo " + limit + " caracteres"
		case isList:
			return "no máximo " + limit + " itens"
		default:
			return "valor máximo " + limit
		}
	case "min":
		switch {
		case isText:
			return "no mínimo " + limit + " caracteres"
		case isList:
			return "no mínimo " + limit + " itens"
		default:
			return "valor mínimo " + limit
		}
	case "len":
		return "deve ter exatamente " + limit + " caracteres"
	case "email":
		return "e-mail inválido"
	case "oneof":
		return "valor não permitido"
	case "gt", "gte", "lt", "lte":
		return "valor fora do limite permitido"
	default:
		return "valor inválido"
	}
}
