package service

import (
	configdata "github.com/ijlik/dating-user/pkg/config/data"
	// business package
	"github.com/ijlik/dating-user/internal/adapter/redis"
	"github.com/ijlik/dating-user/internal/adapter/repository"
	"github.com/ijlik/dating-user/internal/business/port"

	mailerpkg "github.com/ijlik/dating-user/pkg/mailer"
	commonmath "github.com/ijlik/dating-user/pkg/math"
	timemachine "github.com/ijlik/dating-user/pkg/timemachine"
)

type service struct {
	repo   repository.UserRepository
	config configdata.Config
	math   commonmath.Math
	mailer mailerpkg.Mail
	time   timemachine.TimeMachine
	redis  redis.RedisDomain
}

func NewUserService(
	repo repository.UserRepository,
	config configdata.Config,
	mail mailerpkg.Mail,
	redis redis.RedisDomain,
) port.UserDomainService {
	dateTime := timemachine.NewTimeMachine()
	math := commonmath.NewMath()
	return &service{
		repo,
		config,
		math,
		mail,
		dateTime,
		redis,
	}
}
