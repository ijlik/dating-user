package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type DefaultResponse struct {
	Code     string `json:"code"`
	Message  string `json:"message"`
	HttpCode int    `json:"-"`
}

const (
	UnauthorizedCode = "01"
)

func UnauthorizedResponse(ctx *gin.Context, msg string) {
	ctx.JSON(http.StatusUnauthorized, DefaultResponse{
		Code:    UnauthorizedCode,
		Message: msg,
	})
	ctx.Abort()
}
