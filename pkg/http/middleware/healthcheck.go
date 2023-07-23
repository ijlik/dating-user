package middleware

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"github.com/jmoiron/sqlx"
	"net/http"
)

type HealthCheckOptions func(hcc *healthCheckConfig)

func WithRedis(rdb redis.Cmdable) HealthCheckOptions {
	return func(hcc *healthCheckConfig) {
		hcc.rdb = rdb
		hcc.withRedis = true
	}
}

func WithDatabase(db *sqlx.DB) HealthCheckOptions {
	return func(hcc *healthCheckConfig) {
		hcc.db = db
		hcc.withDatabase = true
	}
}

type healthCheckConfig struct {
	service      string
	db           *sqlx.DB
	rdb          redis.Cmdable
	withRedis    bool
	withDatabase bool
}

func HealthCheckHandler(service string, r *gin.Engine, ops ...HealthCheckOptions) {
	config := defaultConfig(service)

	for _, v := range ops {
		v(config)
	}

	// add handler here
	r.GET(fmt.Sprintf("%s/health-check", service), config.healthCheckHandler)
}

// default config
func defaultConfig(service string) *healthCheckConfig {
	var hcc = healthCheckConfig{
		withRedis:    false,
		withDatabase: false,
		service:      service,
	}

	return &hcc
}

// validate all health check handler here
func (hcc *healthCheckConfig) healthCheckHandler(c *gin.Context) {
	if hcc.withRedis {
		err := hcc.rdb.Ping(context.Background()).Err()
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"status":  "redis ping failed",
				"message": err.Error(),
			})
			return
		}
	}

	if hcc.withDatabase {
		err := hcc.db.Ping()
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"status":  "postgres ping failed",
				"message": err.Error(),
			})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "service online",
	})
}
