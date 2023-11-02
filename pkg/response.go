package pkg

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type APIResponse struct {
	Code int         `json:"code"`
	Data interface{} `json:"data,omitempty"`
	Msg  string      `json:"message,omitempty"`
}

func SendResponse(c *gin.Context, code int, data interface{}, msg string) {
	resp := APIResponse{
		Code: code,
		Data: data,
		Msg:  msg,
	}
	c.JSON(code, resp)
}

func SendSuccess(c *gin.Context, data interface{}, msg string) {
	SendResponse(c, http.StatusOK, data, msg)
}

func SendError(c *gin.Context, code int, msg string) {
	SendResponse(c, code, nil, msg)
}

func SendBadRequestError(c *gin.Context, msg string) {
	SendError(c, http.StatusBadRequest, msg)
}

func SendInternalServerError(c *gin.Context, msg string) {
	SendError(c, http.StatusInternalServerError, msg)
}

func StatusUnauthorized(c *gin.Context, msg string) {
	SendError(c, http.StatusUnauthorized, msg)
}
