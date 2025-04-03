package db

import (
	"context"

	"github.com/guncv/Poll-Voting-Website/backend/config"
	"github.com/redis/go-redis/v9"
)

type CacheService interface {
	Get(key string) (string, error)
	Set(key string, value string) error
}

type RedisCacheService struct {
	rdb *redis.Client
	ctx context.Context
}

func NewRedisCacheService(cfg config.Config) *RedisCacheService {
	rdb := redis.NewClient(&redis.Options{
		Addr:     cfg.RedisConfig.Host + ":" + cfg.RedisConfig.Port,
		Password: cfg.RedisConfig.Password, // Set password if needed
		DB:       cfg.RedisConfig.DB,      // Use default DB
	})
	if err := rdb.Ping(context.Background()).Err(); err != nil {
		panic("Failed to connect to Redis: " + err.Error())
	}

	return &RedisCacheService{
		rdb: rdb,
		ctx: context.Background(),
	}
}

func (r *RedisCacheService) Get(key string) (string, error) {
	val, err := r.rdb.Get(r.ctx, key).Result()
	if err == redis.Nil {
		return "", nil // Key not found
	}
	return val, err
}

func (r *RedisCacheService) Set(key string, value string) error {
	return r.rdb.Set(r.ctx, key, value, 0).Err()
}
