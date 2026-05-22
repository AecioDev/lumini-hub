package path

import (
	"fmt"
	"simple-erp-service/internal/utils/ginutils"
	"strconv"

	"github.com/gin-gonic/gin"
)

// busca e converte um uint de um parametro do path
func UintFromPathParam(c *gin.Context, paramName string) (uint, error) {
	paramUint, err := strconv.ParseUint(c.Param(paramName), 10, 64)
	if err != nil {
		return 0, fmt.Errorf("path param: %s deve ser um integer positivo", paramName)
	}

	return uint(paramUint), nil
}

// busca um ID do path, respondendo com BadRequest no caso de erro
func IdFromPathParamOrSendError(c *gin.Context) (uint, error) {
	id, err := UintFromPathParam(c, "id")
	if err != nil {
		ginutils.Res(c).StatusBadRequest().SendError(err)
	}

	return id, err
}
