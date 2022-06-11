package services

import (
	"context"
	"regexp"
	"strings"

	"github.com/go-redis/redis/v8"
	"github.com/pkg/errors"
)

var Redis *RedisService

type RedisServicer interface {
	SetKey(ctx context.Context, key string, value interface{})
	ScanKeys(ctx context.Context) ([]string, error)
	Exists(ctx context.Context, key string) (int64, error)
	GetAllKeys(ctx context.Context) ([]string, error)
	GetValues(ctx context.Context, key string) (string, error)
	Delete(ctx context.Context, key string) (int64, error)
}

type RedisService struct {
	Client *redis.Client
}

var _ RedisServicer = (*RedisService)(nil)

func NewRedisService(address string) *RedisService {
	var addr string
	if matched, _ := regexp.MatchString("http://", address); matched {
		addrSlice := strings.Split(address, "/")
		addr = addrSlice[2]
	} else {
		addr = address
	}
	return &RedisService{
		Client: redis.NewClient(&redis.Options{
			Addr: addr,
		}),
	}
}

func (r *RedisService) SetKey(ctx context.Context, key string, value interface{}) {
	r.Client.Set(ctx, key, value, 0)
}

func (r *RedisService) ScanKeys(ctx context.Context) ([]string, error) {
	var cursor uint64
	keys, _, err := r.Client.Scan(ctx, cursor, "", 0).Result()
	if err != nil {
		return nil, errors.WithStack(err)
	}
	return keys, nil
}

func (r *RedisService) Exists(ctx context.Context, key string) (int64, error) {
	exists, err := r.Client.Exists(ctx, key).Result()
	if err != nil {
		return 0, errors.WithStack(err)
	}
	return exists, nil
}

func (r *RedisService) GetAllKeys(ctx context.Context) ([]string, error) {
	keys, err := r.Client.Keys(ctx, "*").Result()
	if err != nil {
		return nil, errors.WithStack(err)
	}
	return keys, nil
}

func (r *RedisService) GetValues(ctx context.Context, key string) (string, error) {
	value, err := r.Client.Get(ctx, key).Result()
	if err != nil {
		return "", errors.WithStack(err)
	}
	return value, nil
}

func (r *RedisService) Delete(ctx context.Context, key string) (int64, error) {
	result, err := r.Client.Del(ctx, key).Result()
	if err != nil {
		return -1, errors.WithStack(err)
	}
	return result, nil
}
