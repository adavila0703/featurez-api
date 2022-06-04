package clients

import (
	"context"

	redis "github.com/go-redis/redis/v8"
)

var Redis *RedisClient

type RedisClient struct {
	Client *redis.Client
}

func NewRedisClient(address string) *RedisClient {
	return &RedisClient{
		Client: redis.NewClient(&redis.Options{
			Addr: address,
		}),
	}
}

func (r *RedisClient) SetKey(ctx context.Context, key string, value interface{}) {
	r.Client.Set(ctx, key, value, 0)
}

func (r *RedisClient) ScanKeys(ctx context.Context) ([]string, error) {
	var cursor uint64
	keys, _, err := r.Client.Scan(ctx, cursor, "", 0).Result()
	return keys, err
}

func (r *RedisClient) Exists(ctx context.Context, key string) (int64, error) {
	exists, err := r.Client.Exists(ctx, key).Result()
	return exists, err
}

func (r *RedisClient) GetAllKeys(ctx context.Context) ([]string, error) {
	keys, err := r.Client.Keys(ctx, "*").Result()
	return keys, err
}

func (r *RedisClient) GetValues(ctx context.Context, key string) (string, error) {
	value, err := r.Client.Get(ctx, key).Result()
	return value, err
}

func (r *RedisClient) Delete(ctx context.Context, key string) (int64, error) {
	result, err := r.Client.Del(ctx, key).Result()
	return result, err
}
