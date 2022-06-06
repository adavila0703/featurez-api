package services

import (
	"context"

	redis "github.com/go-redis/redis/v8"
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
	return &RedisService{
		Client: redis.NewClient(&redis.Options{
			Addr: address,
		}),
	}
}

func (r *RedisService) SetKey(ctx context.Context, key string, value interface{}) {
	r.Client.Set(ctx, key, value, 0)
}

func (r *RedisService) ScanKeys(ctx context.Context) ([]string, error) {
	var cursor uint64
	keys, _, err := r.Client.Scan(ctx, cursor, "", 0).Result()
	return keys, err
}

func (r *RedisService) Exists(ctx context.Context, key string) (int64, error) {
	exists, err := r.Client.Exists(ctx, key).Result()
	return exists, err
}

func (r *RedisService) GetAllKeys(ctx context.Context) ([]string, error) {
	keys, err := r.Client.Keys(ctx, "*").Result()
	return keys, err
}

func (r *RedisService) GetValues(ctx context.Context, key string) (string, error) {
	value, err := r.Client.Get(ctx, key).Result()
	return value, err
}

func (r *RedisService) Delete(ctx context.Context, key string) (int64, error) {
	result, err := r.Client.Del(ctx, key).Result()
	return result, err
}
