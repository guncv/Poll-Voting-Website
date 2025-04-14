package db

import (
	"context"
	"fmt"
	"strconv"

	"github.com/guncv/Poll-Voting-Website/backend/config"
	"github.com/redis/go-redis/v9"
)

type CacheService interface {
	Get(key string) (string, error)
	Set(key string, value string) error
	IsSetMember(key, member string) (bool, error)
	AddSetMember(key, member string) error
	IncrementField(key, field string) int64
	GetField(key, field string) (string, error)
	GetFieldInt(key, field string) (int, error)
	SetHash(key string, data map[string]string) error
	GetAllHash(key string) (map[string]string, error)
	AddToSet(key, value string) error
	GetSetMembers(key string) ([]string, error)
	DeleteKey(key string) error
}

type RedisCacheService struct {
	rdb *redis.Client
	ctx context.Context
}

func NewRedisCacheService(cfg config.Config) *RedisCacheService {
	rdb := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", cfg.RedisConfig.Host, cfg.RedisConfig.Port),
		Password: cfg.RedisConfig.Password, // Set password if needed
		DB:       cfg.RedisConfig.DB,       // Use default DB
	})

	ctx := context.Background()

	if err := rdb.Ping(ctx).Err(); err != nil {
		panic("❌ Failed to connect to Redis: " + err.Error())
	}

	fmt.Println("✅ Redis connected successfully at", cfg.RedisConfig.Host+":"+cfg.RedisConfig.Port)

	return &RedisCacheService{
		rdb: rdb,
		ctx: ctx,
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

func (r *RedisCacheService) IsSetMember(key, member string) (bool, error) {
	return r.rdb.SIsMember(r.ctx, key, member).Result()
}

func (r *RedisCacheService) AddSetMember(key, member string) error {
	return r.rdb.SAdd(r.ctx, key, member).Err()
}

func (r *RedisCacheService) IncrementField(key, field string) int64 {
	val, err := r.rdb.HIncrBy(r.ctx, key, field, 1).Result()
	if err != nil {
		return 0
	}
	return val
}

func (r *RedisCacheService) GetField(key, field string) (string, error) {
	val, err := r.rdb.HGet(r.ctx, key, field).Result()
	if err == redis.Nil {
		return "", nil
	}
	return val, err
}

func (r *RedisCacheService) GetFieldInt(key, field string) (int, error) {
	valStr, err := r.GetField(key, field)
	if err != nil {
		return 0, err
	}
	valInt, _ := strconv.Atoi(valStr)
	return valInt, nil
}

func (r *RedisCacheService) SetHash(key string, data map[string]string) error {
	return r.rdb.HSet(r.ctx, key, data).Err()
}

func (r *RedisCacheService) GetAllHash(key string) (map[string]string, error) {
	return r.rdb.HGetAll(r.ctx, key).Result()
}

func (r *RedisCacheService) AddToSet(key, value string) error {
	return r.rdb.SAdd(r.ctx, key, value).Err()
}

func (r *RedisCacheService) GetSetMembers(key string) ([]string, error) {
	return r.rdb.SMembers(r.ctx, key).Result()
}

func (r *RedisCacheService) DeleteKey(key string) error {
	return r.rdb.Del(r.ctx, key).Err()
}
