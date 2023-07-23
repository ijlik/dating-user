package http

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/ijlik/dating-user/internal/business/domain"
	ctxsdk "github.com/ijlik/dating-user/pkg/context"
	errpkg "github.com/ijlik/dating-user/pkg/error"
	httppkg "github.com/ijlik/dating-user/pkg/http"
)

func (rh *requestHandler) CreatePayment(c *gin.Context) {
	ctx := c.Request.Context()
	UserID := fmt.Sprintf("%v", ctx.Value(ctxsdk.USER_ID))
	if UserID == "" {
		httppkg.BuildErrorResponse(c, errpkg.ErrUnauthorize, "")
		return
	}
	var request domain.PaymentRequest
	err := decodeRequest(c, &request)
	if err != nil {
		httppkg.BuildErrorResponse(c, errpkg.ErrBadRequest, err.Error())
		return
	}
	if err := request.Validate(); err != nil {
		httppkg.BuildErrorResponse(c, errpkg.ErrBadRequest, err.Error())
		return
	}

	err = rh.service.CreatePayment(ctx, &request, UserID)
	if err != nil {
		httppkg.BuildErrorResponse(c, errpkg.ErrBadRequest, err.Error())
		return
	}
	response := httppkg.DefaultSuccessResponse(nil)
	c.JSON(response.HttpCode, response)
}
