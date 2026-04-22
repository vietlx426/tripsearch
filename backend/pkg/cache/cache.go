package cache

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

var client *redis.Client

func Init(redisURL string) error {
	opts, err := redis.ParseURL(redisURL)
	if err != nil {
		return err
	}
	client = redis.NewClient(opts)
	return client.Ping(context.Background()).Err()
}

func Get() *redis.Client {
	return client
}

func Set(ctx context.Context, key string, value any, ttl time.Duration) error {
	return client.Set(ctx, key, value, ttl).Err()
}

func GetValue(ctx context.Context, key string) (string, error) {
	return client.Get(ctx, key).Result()
}

func Delete(ctx context.Context, key string) error {
	return client.Del(ctx, key).Err()
}
