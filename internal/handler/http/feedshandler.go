package http

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/ijlik/dating-user/internal/business/domain"
	ctxsdk "github.com/ijlik/dating-user/pkg/context"
	errpkg "github.com/ijlik/dating-user/pkg/error"
	httppkg "github.com/ijlik/dating-user/pkg/http"
)

func (rh *requestHandler) ShowFeeds(c *gin.Context) {
	ctx := c.Request.Context()
	swiperId := fmt.Sprintf("%v", ctx.Value(ctxsdk.PROFILE_ID))
	if swiperId == "" {
		httppkg.BuildErrorResponse(c, errpkg.ErrUnauthorize, "")
		return
	}

	profileId := c.Query("profile_id")

	data, err := rh.service.ShowFeeds(ctx, swiperId, profileId)
	if err != nil {
		httppkg.BuildErrorResponse(c, errpkg.ErrBadRequest, err.Error())
		return
	}

	response := httppkg.DefaultSuccessResponse(data)
	c.JSON(response.HttpCode, response)
}

func (rh *requestHandler) Swipes(c *gin.Context) {
	ctx := c.Request.Context()
	UserID := fmt.Sprintf("%v", ctx.Value(ctxsdk.USER_ID))
	if UserID == "" {
		httppkg.BuildErrorResponse(c, errpkg.ErrUnauthorize, "")
		return
	}

	var request domain.SwipeRequest
	err := decodeRequest(c, &request)
	if err != nil {
		httppkg.BuildErrorResponse(c, errpkg.ErrBadRequest, err.Error())
		return
	}
	if err := request.Validate(); err != nil {
		httppkg.BuildErrorResponse(c, errpkg.ErrBadRequest, err.Error())
		return
	}

	err = rh.service.Swipes(ctx, &request, UserID)
	if err != nil {
		httppkg.BuildErrorResponse(c, errpkg.ErrBadRequest, err.Error())
		return
	}
	response := httppkg.DefaultSuccessResponse(nil)
	c.JSON(response.HttpCode, response)
}
