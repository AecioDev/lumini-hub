package utils

import (
	"fmt"
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

// ErrorResponse envia uma resposta de erro
func ErrorResponse(c *gin.Context, statusCode int, message string, err string) {
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
			validationErrors = v
		default:
			validationErrors = []string{fmt.Sprintf("%v", v)}
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

