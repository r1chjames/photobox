package service

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/redis/go-redis/v9"
	"gitlab.com/r1chjames/photobox/api/internal/appconfig"
)

type CacheService struct {
	client  *redis.Client
	enabled bool
	ctx     context.Context
}

func NewCacheService(config appconfig.AppConfig) *CacheService {
	if !config.CacheEnabled {
		return &CacheService{enabled: false, ctx: context.Background()}
	}
	addr := fmt.Sprintf("%s:%s", config.CacheHost, config.CachePort)
	client := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: config.CachePassword,
		DB:       config.CacheDB,
	})

	ctx := context.Background()
	if err := client.Ping(ctx).Err(); err != nil {
		slog.Warn("Cache connection failed, running without cache", "error", err)
		return &CacheService{enabled: false, ctx: ctx}
	}

	slog.Info("Cache connected", "addr", addr)
	return &CacheService{client: client, enabled: true, ctx: ctx}
}

func (cs *CacheService) Get(key string) (string, error) {
	if !cs.enabled {
		return "", fmt.Errorf("cache disabled")
	}
	val, err := cs.client.Get(cs.ctx, key).Result()
	if err == redis.Nil {
		return "", fmt.Errorf("cache miss")
	}
	return val, err
}

func (cs *CacheService) Set(key string, value string, ttl time.Duration) error {
	if !cs.enabled {
		return nil
	}
	return cs.client.Set(cs.ctx, key, value, ttl).Err()
}

func (cs *CacheService) Delete(key string) error {
	if !cs.enabled {
		return nil
	}
	return cs.client.Del(cs.ctx, key).Err()
}

func (cs *CacheService) DeletePattern(pattern string) error {
	if !cs.enabled {
		return nil
	}
	iter := cs.client.Scan(cs.ctx, 0, pattern, 0).Iterator()
	for iter.Next(cs.ctx) {
		cs.client.Del(cs.ctx, iter.Val())
	}
	return iter.Err()
}
