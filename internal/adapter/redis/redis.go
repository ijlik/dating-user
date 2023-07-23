package redis

import (
	"time"

	"github.com/go-redis/redis/v8"
	"golang.org/x/net/context"
)

func NewRedisRepository(conn redis.Cmdable) RedisDomain {
	return &rdb{conn}
}

type rdb struct {
	conn redis.Cmdable
}

type RedisDomain interface {
	Get(ctx context.Context, key string) (string, error)
	Set(ctx context.Context, value, key string, interval int) error
	Del(ctx context.Context, key string) error
}

func (r *rdb) Get(ctx context.Context, key string) (string, error) {
	value, err := r.conn.Get(ctx, key).Result()
	if err != nil {
		return "", err
	}

	return value, nil
}

func (r *rdb) Set(ctx context.Context, value, key string, interval int) error {
	if interval == 0 {
		interval = 5
	}

	return r.conn.Set(ctx, key, value, time.Duration(interval*int(time.Minute))).Err()
}

func (r *rdb) Del(ctx context.Context, key string) error {
	return r.conn.Del(ctx, key).Err()
}
