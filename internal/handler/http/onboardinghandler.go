package http

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/ijlik/dating-user/internal/business/domain"
	ctxsdk "github.com/ijlik/dating-user/pkg/context"
	errpkg "github.com/ijlik/dating-user/pkg/error"
	httppkg "github.com/ijlik/dating-user/pkg/http"
	"mime/multipart"
)

func (rh *requestHandler) UpdatePersonalInfo(c *gin.Context) {
	ctx := c.Request.Context()
	UserID := fmt.Sprintf("%v", ctx.Value(ctxsdk.USER_ID))
	if UserID == "" {
		httppkg.BuildErrorResponse(c, errpkg.ErrUnauthorize, "")
		return
	}

	var request domain.UpdatePersonalInfo
	err := decodeRequest(c, &request)
	if err != nil {
		httppkg.BuildErrorResponse(c, errpkg.ErrBadRequest, err.Error())
		return
	}
	if err := request.Validate(); err != nil {
		httppkg.BuildErrorResponse(c, errpkg.ErrBadRequest, err.Error())
		return
	}

	err = rh.service.UpdatePersonalInfo(ctx, &request, UserID)
	if err != nil {
		httppkg.BuildErrorResponse(c, errpkg.ErrBadRequest, err.Error())
		return
	}

	response := httppkg.DefaultSuccessResponse(nil)
	c.JSON(response.HttpCode, response)
}

func (rh *requestHandler) UpdatePhotos(c *gin.Context) {
	ctx := c.Request.Context()
	UserID := fmt.Sprintf("%v", ctx.Value(ctxsdk.USER_ID))
	if UserID == "" {
		httppkg.BuildErrorResponse(c, errpkg.ErrUnauthorize, "")
		return
	}
	var photos []*multipart.FileHeader
	for _, i := range []string{"0", "1", "2", "3", "4"} {
		formName := fmt.Sprintf("photos[%s]", i)
		file, err := c.FormFile(formName)
		if err != nil {
			break
		}
		photos = append(photos, file)
	}
	request := domain.UpdatePhotos{Photos: photos}
	err := request.Validate()
	if err != nil {
		httppkg.BuildErrorResponse(c, errpkg.ErrBadRequest, err.Error())
		return
	}

	err = rh.service.UpdatePhotos(ctx, &request, UserID)
	if err != nil {
		httppkg.BuildErrorResponse(c, errpkg.ErrBadRequest, err.Error())
		return
	}

	response := httppkg.DefaultSuccessResponse(nil)
	c.JSON(response.HttpCode, response)
}

func (rh *requestHandler) UpdateHobbyAndInterest(c *gin.Context) {
	ctx := c.Request.Context()
	UserID := fmt.Sprintf("%v", ctx.Value(ctxsdk.USER_ID))
	if UserID == "" {
		httppkg.BuildErrorResponse(c, errpkg.ErrUnauthorize, "")
		return
	}

	var request domain.UpdateHobbyAndInterest
	err := decodeRequest(c, &request)
	if err != nil {
		httppkg.BuildErrorResponse(c, errpkg.ErrBadRequest, err.Error())
		return
	}
	if err := request.Validate(); err != nil {
		httppkg.BuildErrorResponse(c, errpkg.ErrBadRequest, err.Error())
		return
	}

	err = rh.service.UpdateHobbyAndInterest(ctx, &request, UserID)
	if err != nil {
		httppkg.BuildErrorResponse(c, errpkg.ErrBadRequest, err.Error())
		return
	}

	response := httppkg.DefaultSuccessResponse(nil)
	c.JSON(response.HttpCode, response)
}

func (rh *requestHandler) UpdateLocation(c *gin.Context) {
	ctx := c.Request.Context()
	UserID := fmt.Sprintf("%v", ctx.Value(ctxsdk.USER_ID))
	if UserID == "" {
		httppkg.BuildErrorResponse(c, errpkg.ErrUnauthorize, "")
		return
	}

	var request domain.Location
	err := decodeRequest(c, &request)
	if err != nil {
		httppkg.BuildErrorResponse(c, errpkg.ErrBadRequest, err.Error())
		return
	}
	if err := request.Validate(); err != nil {
		httppkg.BuildErrorResponse(c, errpkg.ErrBadRequest, err.Error())
		return
	}

	err = rh.service.UpdateLocation(ctx, &request, UserID)
	if err != nil {
		httppkg.BuildErrorResponse(c, errpkg.ErrBadRequest, err.Error())
		return
	}

	response := httppkg.DefaultSuccessResponse(nil)
	c.JSON(response.HttpCode, response)
}
