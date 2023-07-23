package main

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/go-redsync/redsync/v4"
	configenv "github.com/ijlik/dating-user/pkg/config"
	configdata "github.com/ijlik/dating-user/pkg/config/data"
	httpmiddlewaresdk "github.com/ijlik/dating-user/pkg/http/middleware"
	"github.com/jmoiron/sqlx"

	rediseight "github.com/go-redis/redis/v8"
	_rsyncpool "github.com/go-redsync/redsync/v4/redis/goredis/v8"
	// internal package
	rdbrepo "github.com/ijlik/dating-user/internal/adapter/redis"
	"github.com/ijlik/dating-user/internal/adapter/repository"
	"github.com/ijlik/dating-user/internal/business/port"
	"github.com/ijlik/dating-user/internal/business/service"
	httpdelivery "github.com/ijlik/dating-user/internal/handler/http"
	mailerpkg "github.com/ijlik/dating-user/pkg/mailer"
	_ "github.com/lib/pq"
)

var config configdata.Config

type RedisModule struct {
	*redsync.Redsync
	rediseight.Cmdable
}

func getRedisClient() rediseight.Cmdable {

	redisHost := config.GetString("REDIS_ADDR")
	redisPort := config.GetString("REDIS_PORT")
	redisAddr := fmt.Sprintf("%s:%s", redisHost, redisPort)

	module := &RedisModule{}
	rc := rediseight.NewClient(&rediseight.Options{
		Addr: redisAddr,
	})

	pool := _rsyncpool.NewPool(rc)
	rs := redsync.New(pool)

	module.Cmdable = rc
	module.Redsync = rs

	if err := module.Ping(context.Background()).Err(); err != nil {
		log.Println("redis connection error", err)
		return nil
	}

	return module
}

func getService(
	db *sqlx.DB,
	rdb rdbrepo.RedisDomain,
) port.UserDomainService {
	mailPort := config.GetInt("MAILER_PORT")
	mailUsername := config.GetString("MAILER_USERNAME")
	mailHost := config.GetString("MAILER_HOST")
	mailPassword := config.GetString("MAILER_PASSWORD")
	mailFrom := config.GetString("MAILER_FROM")

	mailer := mailerpkg.NewMailer(
		mailPort,
		mailUsername,
		mailFrom,
		mailHost,
		mailPassword,
	)

	repo := repository.NewUserRepo(db)
	services := service.NewUserService(
		repo,
		config,
		mailer,
		rdb,
	)

	return services
}

func getConfig() configdata.Config {
	c := configenv.NewConfig("", 5)

	if c == nil {
		panic(errors.New("missing config"))
	}

	return c
}

func getDatabase() (*sqlx.DB, error) {
	var (
		host     = config.GetString("DB_HOST")
		port     = config.GetInt("DB_PORT")
		user     = config.GetString("DB_USER")
		password = config.GetString("DB_PASSWORD")
		dbname   = config.GetString("DB_NAME")
		timeZone = "UTC"
	)

	dsn := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable TimeZone=%s",
		host, port, user, password, dbname, timeZone)

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, err
	}

	if err != nil {
		return nil, err
	}

	return sqlx.NewDb(db, "postgres"), nil
}

func main() {
	// get config
	config = getConfig()

	db, err := getDatabase()
	if err != nil {
		panic(err)
	}

	defer db.Close()

	rdb := getRedisClient()
	rdbConn := rdbrepo.NewRedisRepository(rdb)

	router := gin.Default()
	router.Use(
		httpmiddlewaresdk.WithAllowedCORS(),
	)

	options := []httpmiddlewaresdk.HealthCheckOptions{
		httpmiddlewaresdk.WithRedis(rdb),
		httpmiddlewaresdk.WithDatabase(db),
	}
	httpmiddlewaresdk.HealthCheckHandler("user", router, options...)

	services := getService(db, rdbConn)

	httpdelivery.HandlerHttp(
		router,
		config,
		services,
		rdb,
	)
}
