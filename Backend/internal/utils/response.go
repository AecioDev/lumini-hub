package utils

import (
	"net/http"
	"simple-erp-service/internal/utils/ginutils"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

// Response representa a estrutura padrão de resposta da API
type Response struct {
	Success bool        `json:"success"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
	Meta    interface{} `json:"meta,omitempty"`
}

// SuccessResponse envia uma resposta de sucesso
func SuccessResponse(c *gin.Context, statusCode int, message string, data interface{}, meta interface{}) {
	c.JSON(statusCode, Response{
		Success: true,
		Message: message,
		Data:    data,
		Meta:    meta,
	})
}

// ErrorResponse envia uma resposta de erro
func ErrorResponse(c *gin.Context, statusCode int, message string, err string) {
	c.JSON(statusCode, Response{
		Success: false,
		Message: message,
		Error:   err,
	})
}

// ValidationErrorResponse envia uma resposta de erro de validação
func ValidationErrorResponse(c *gin.Context, message string, errors interface{}) {
	c.JSON(http.StatusBadRequest, Response{
		Success: false,
		Message: message,
		Data:    errors,
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
