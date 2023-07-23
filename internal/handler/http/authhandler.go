package http

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ijlik/dating-user/internal/business/domain"
	errpkg "github.com/ijlik/dating-user/pkg/error"
	httppkg "github.com/ijlik/dating-user/pkg/http"

	ctxsdk "github.com/ijlik/dating-user/pkg/context"
)

func (rh *requestHandler) ResendOtp(c *gin.Context) {
	ctx := c.Request.Context()
	var request domain.ResendOtpRequest
	err := decodeRequest(c, &request)
	if err != nil {
		httppkg.BuildErrorResponse(c, errpkg.ErrBadRequest, err.Error())
		return
	}

	allowedDisposable := rh.config.GetBool("ALLOWED_DISPOSABLE_EMAIL")

	if err := request.Validate(allowedDisposable); err != nil {
		httppkg.BuildErrorResponse(c, errpkg.ErrBadRequest, err.Error())
		return
	}

	data, err := rh.service.ResendOtp(ctx, &request)
	if err != nil {
		httppkg.BuildErrorResponse(c, errpkg.ErrBadRequest, err.Error())
		return
	}

	response := httppkg.DefaultSuccessResponse(data)
	c.JSON(response.HttpCode, response)
}

func (rh *requestHandler) LoginOrRegister(c *gin.Context) {
	ctx := c.Request.Context()
	var request domain.AuthRequest
	err := decodeRequest(c, &request)
	if err != nil {
		httppkg.BuildErrorResponse(c, errpkg.ErrBadRequest, "")
		return
	}

	if err := request.Validate(); err != nil {
		httppkg.BuildErrorResponse(c, errpkg.ErrBadRequest, err.Error())
		return
	}

	data, err := rh.service.LoginOrRegister(ctx, &request)
	if err != nil {
		httppkg.BuildErrorResponse(c, errpkg.ErrBadRequest, err.Error())
		return
	}

	response := httppkg.DefaultSuccessResponse(data)
	c.JSON(response.HttpCode, response)
}

func (rh *requestHandler) Logout(c *gin.Context) {
	ctx := c.Request.Context()
	UserID := fmt.Sprintf("%v", ctx.Value(ctxsdk.USER_ID))
	if UserID == "" {
		httppkg.BuildErrorResponse(c, errpkg.ErrUnauthorize, "")
		return
	}

	errs := rh.service.Logout(ctx, fmt.Sprintf("%v", UserID))
	if errs != nil {
		httppkg.BuildErrorResponse(c, errs.GetCode(), errs.Error())
		return
	}

	c.JSON(http.StatusOK, httppkg.DefaultSuccessResponse(nil))
}

func (rh *requestHandler) ShowProfile(c *gin.Context) {
	ctx := c.Request.Context()
	UserID := fmt.Sprintf("%v", ctx.Value(ctxsdk.USER_ID))
	if UserID == "" {
		httppkg.BuildErrorResponse(c, errpkg.ErrUnauthorize, "")
		return
	}
	email := fmt.Sprintf("%v", ctx.Value(ctxsdk.EMAIL))
	if email == "" {
		httppkg.BuildErrorResponse(c, errpkg.ErrUnauthorize, "")
		return
	}

	data, errs := rh.service.ShowProfile(ctx, UserID)
	if errs != nil {
		httppkg.BuildErrorResponse(c, errpkg.ErrUnauthorize, errs.Error())
		return
	}

	response := httppkg.DefaultSuccessResponse(data)
	c.JSON(response.HttpCode, response)
}
