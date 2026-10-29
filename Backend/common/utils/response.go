package utils

import (
	"fmt"
	"log"
	"net/http"

	"lumini-hub/common/utils/ginutils"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

// Response representa a estrutura padrão de resposta da API (ResponseDTO)
type Response struct {
	Success          bool        `json:"success"`
	StatusCode       int         `json:"statusCode"`
	Message          string      `json:"message,omitempty"`
	ValidationErrors []string    `json:"validationErrors,omitempty"`
	Data             interface{} `json:"data,omitempty"`
	Error            string      `json:"error,omitempty"`
	Meta             interface{} `json:"meta,omitempty"`
}

// SuccessResponse envia uma resposta de sucesso
func SuccessResponse(c *gin.Context, statusCode int, message string, data interface{}, meta interface{}) {
	c.JSON(statusCode, Response{
		Success:    true,
		StatusCode: statusCode,
		Message:    message,
		Data:       data,
		Meta:       meta,
	})
}

// ErrorResponse envia uma resposta de erro.
//
// Texto de erro do banco de dados (SQLSTATE, falha de conexão) nunca vai ao cliente em nenhum
// status: o original fica no log do servidor e a resposta leva GenericOperationError.
//
// Erro interno (500) nunca devolve o texto do erro ao cliente: os handlers passam o
// err.Error() cru, que num erro de banco traz nome de tabela/constraint e até valores
// (ex.: `invalid input syntax for type integer: "abc"`). No 500 o erro completo vai para o
// log do servidor, com método e caminho, e a resposta leva só a `message` (o frontend cai
// nela quando não há `error`). Os demais status continuam como estão: são mensagens de
// regra de negócio escritas pra o usuário ver (e o 503 do SQL Server legado é um texto fixo).
func ErrorResponse(c *gin.Context, statusCode int, message string, err string) {
	if statusCode == http.StatusInternalServerError {
		logInternalError(c, message, err)
		err = ""
	} else {
		// Em qualquer outro status, erro de banco nunca vai ao cliente (ver response_dberror.go).
		err = hideDatabaseError(c, statusCode, message, err)
	}
	c.JSON(statusCode, Response{
		Success:    false,
		StatusCode: statusCode,
		Message:    message,
		Error:      err,
	})
}

// ValidationErrorResponse envia uma resposta de erro de validação
func ValidationErrorResponse(c *gin.Context, message string, errors interface{}) {
	var validationErrors []string

	if errors != nil {
		switch v := errors.(type) {
		case string:
			validationErrors = []string{v}
		case []string:
			validationErrors = append([]string(nil), v...)
		default:
			validationErrors = []string{fmt.Sprintf("%v", v)}
		}
		for i, text := range validationErrors {
			validationErrors[i] = hideDatabaseError(c, http.StatusBadRequest, message, text)
		}
	}

	c.JSON(http.StatusBadRequest, Response{
		Success:          false,
		StatusCode:       http.StatusBadRequest,
		Message:          message,
		ValidationErrors: validationErrors,
	})
}

func handleBindError(c *gin.Context, err error) error {
	if err != nil {
		ginutils.Res(c).StatusBadRequest().SendError(err)
	}

	return err
}

func BindJsonOrSendErrorRes(c *gin.Context, obj any) error {
	return handleBindError(c, c.ShouldBindJSON(obj))
}

func BindQueryOrSendErrorRes(c *gin.Context, obj any) error {
	return handleBindError(c, c.ShouldBindQuery(obj))
}

func BindUriOrSendErrorRes(c *gin.Context, obj any) error {
	return handleBindError(c, c.ShouldBindUri(obj))
}

func BindFormOrSendErrorRes(c *gin.Context, obj any) error {
	return handleBindError(c, c.ShouldBindWith(obj, binding.Form))
}

// GetUserIDFromContext retorna o ID do usuário autenticado do contexto
func GetUserIDFromContext(c *gin.Context) (uint, bool) {
	userIDRaw, exists := c.Get("userID")
	if !exists {
		return 0, false
	}

	userID, ok := userIDRaw.(uint)
	if !ok {
		return 0, false
	}

	return userID, true
}

// logInternalError grava o erro completo de um 500 no log do servidor.
func logInternalError(c *gin.Context, message string, err string) {
	method, path := requestInfo(c)
	if err == "" {
		log.Printf("[erro 500] %s %s: %s", method, path, message)
		return
	}
	log.Printf("[erro 500] %s %s: %s: %s", method, path, message, err)
}
