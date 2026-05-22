package ginutils

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Uma resposta a ser enviada
//
// é basicamente uma struct com métodos auxiliares para enviar
// respostas http mais facilmente, usando algo semelhante ao builder pattern
type GinResponse struct {
	c      *gin.Context
	status int
}

// Cria uma nova resposta com status ok
func Res(c *gin.Context) *GinResponse {
	return &GinResponse{
		c,
		http.StatusOK,
	}
}

// Seta o status da resposta a ser enviada
func (r *GinResponse) Status(code int) *GinResponse {
	r.status = code
	return r
}

func (r *GinResponse) StatusBadRequest() *GinResponse {
	return r.Status(http.StatusBadRequest)
}

func (r *GinResponse) StatusForbidden() *GinResponse {
	return r.Status(http.StatusForbidden)
}

func (r *GinResponse) StatusInternalServerError() *GinResponse {
	return r.Status(http.StatusInternalServerError)
}

func (r *GinResponse) StatusNotFound() *GinResponse {
	return r.Status(http.StatusNotFound)
}

func (r *GinResponse) StatusUnauthorized() *GinResponse {
	return r.Status(http.StatusUnauthorized)
}

// Seta o status caso err != nil
func (r *GinResponse) OnErrorSetStatus(code int, err error) *GinResponse {
	if err != nil {
		return r.Status(code)
	}

	return r
}

// Seta o status caso err == nil
func (r *GinResponse) OnNilErrorSetStatus(code int, err error) *GinResponse {
	if err == nil {
		return r.Status(code)
	}

	return r
}

func (r *GinResponse) SendErrorMsg(errorMessage string) {
	r.c.JSON(r.status, gin.H{"error": errorMessage})
}

func (r *GinResponse) SendError(err error) {
	if err != nil {
		r.SendErrorMsg(err.Error())
		return
	}

	r.SendErrorMsg("erro ao lidar com requisição")
}

func (r *GinResponse) SendJson(json any) {
	r.c.JSON(r.status, json)
}

func (r *GinResponse) SendJsonOrError(json any, err error) {
	if err != nil {
		r.SendError(err)
		return
	}

	r.c.JSON(r.status, json)
}

func (r *GinResponse) SendSuccessMessage(msg string) {
	r.c.JSON(r.status, gin.H{"message": msg})
}

func (r *GinResponse) SendSuccessMessageOrError(msg string, err error) {
	if err != nil {
		r.SendError(err)
		return
	}

	r.SendSuccessMessage(msg)
}

func (r *GinResponse) SendJsonOrInternalError(json any, err error) {
	if err != nil {
		r.StatusInternalServerError().SendError(err)
		return
	}

	r.c.JSON(r.status, json)
}

func (r *GinResponse) SendSuccessMessageOrInternalError(msg string, err error) {
	if err != nil {
		r.StatusInternalServerError().SendError(err)
		return
	}

	r.c.JSON(r.status, gin.H{"message": msg})
}
