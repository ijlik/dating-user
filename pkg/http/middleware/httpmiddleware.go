package middleware

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	ctxsdk "github.com/ijlik/dating-user/pkg/context"
	"github.com/ijlik/dating-user/pkg/jwt"
)

type DatabaseData struct {
	UserName string
	Password string
	Host     string
	Port     int
	Database string
}

type SecretData func(ctx context.Context, Id string) (DatabaseData, error)

const (
	authorization = "Authorization"
	signatureX    = "SignatureX"
	tokenData     = "tokenData"
)

func WithLoginAndRedis(
	pubKey string,
	redis redis.Cmdable,
) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		auth := ctx.GetHeader(authorization)
		ArrAuth := strings.Split(auth, " ")
		if len(ArrAuth) != 2 {
			UnauthorizedResponse(ctx, "Missing authorization header")
			return
		}

		resp, err := jwt.ValidateToken(ArrAuth[1], pubKey)
		if err != nil {
			UnauthorizedResponse(ctx, err.Error())
			return
		}

		data, err := json.Marshal(resp)
		if err != nil {
			UnauthorizedResponse(ctx, err.Error())
			return
		}

		var md map[ctxsdk.ContextMetadata]any
		err = json.Unmarshal([]byte(data), &md)
		if err != nil {
			UnauthorizedResponse(ctx, err.Error())
			return
		}

		// append medata token to context value
		cmd := ctxsdk.SetContext(ctx.Request.Context(), md)
		ctx.Request = ctx.Request.WithContext(cmd)

		// validate user already logout or login
		UserID := ctx.Request.Context().Value(ctxsdk.USER_ID)
		_, err = redis.Get(ctx.Request.Context(), fmt.Sprintf("%s", UserID)).Result()
		if err != nil {
			UnauthorizedResponse(ctx, "User already logged out, please login")
			return
		}

		ctx.Next()
	}
}

func WithAllowedCORS() gin.HandlerFunc {
	return func(c *gin.Context) {
		if origin := c.Request.Header.Get("Origin"); origin != "" {
			c.Header("Access-Control-Allow-Origin", origin)
			c.Header("Access-Control-Allow-Credentials", "true")
			c.Header("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With, X-Menu-Slug, X-Origin-Path, X-Request-Id")
			c.Header("Access-Control-Allow-Methods", "POST, HEAD, PATCH, OPTIONS, GET, PUT")

			if c.Request.Method == "OPTIONS" {
				c.AbortWithStatus(204)
				return
			}
		}
		c.Next()
	}
}
