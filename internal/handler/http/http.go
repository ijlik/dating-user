package http

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	configdata "github.com/ijlik/dating-user/pkg/config/data"
	httppkg "github.com/ijlik/dating-user/pkg/http"

	"github.com/go-redis/redis/v8"
	"github.com/ijlik/dating-user/internal/business/port"
	httpmiddlewaresdk "github.com/ijlik/dating-user/pkg/http/middleware"
)

type requestHandler struct {
	config  configdata.Config
	service port.UserDomainService
	pubKey  string
	rdb     redis.Cmdable
}

func HandlerHttp(
	router *gin.Engine,
	config configdata.Config,
	service port.UserDomainService,
	rdb redis.Cmdable,
) {
	pubkey := config.GetString("TOKEN_PUBLIC_KEY")

	rh := requestHandler{
		config:  config,
		service: service,
		pubKey:  pubkey,
		rdb:     rdb,
	}

	routeHandler(router, rh)

	addr := config.GetString("HTTP_ADDR")
	if addr == "" {
		addr = ":8080"
	}

	fmt.Println("HTTP Running ON ", addr)

	httppkg.Serve(router, addr)
}

func routeHandler(router *gin.Engine, rh requestHandler) {
	authRoute := router.Group("/auth")
	authRoute.POST("/otp/resend", rh.ResendOtp)
	authRoute.POST("/otp", rh.LoginOrRegister)
	authRoute.Use(
		httpmiddlewaresdk.WithLoginAndRedis(rh.pubKey, rh.rdb),
	).POST("/logout", rh.Logout)
	authRoute.Use(
		httpmiddlewaresdk.WithLoginAndRedis(rh.pubKey, rh.rdb),
	).GET("/me", rh.ShowProfile)

	onboardRoute := router.Group("/on-boarding").Use(
		httpmiddlewaresdk.WithLoginAndRedis(rh.pubKey, rh.rdb),
	)
	onboardRoute.POST("/personal-info", rh.UpdatePersonalInfo)
	onboardRoute.POST("/photos", rh.UpdatePhotos)
	onboardRoute.POST("/hobby-and-interest", rh.UpdateHobbyAndInterest)
	onboardRoute.POST("/location", rh.UpdateLocation)

	feedsRoute := router.Group("/feeds").Use(
		httpmiddlewaresdk.WithLoginAndRedis(rh.pubKey, rh.rdb),
	)
	feedsRoute.GET("", rh.ShowFeeds)
	feedsRoute.POST("", rh.Swipes)

	paymentRoute := router.Group("/payment").Use(
		httpmiddlewaresdk.WithLoginAndRedis(rh.pubKey, rh.rdb),
	)
	paymentRoute.POST("", rh.CreatePayment)
}

func decodeRequest(c *gin.Context, i interface{}) error {
	err := json.NewDecoder(c.Request.Body).Decode(i)
	if err != nil {
		return err
	}

	return nil
}
