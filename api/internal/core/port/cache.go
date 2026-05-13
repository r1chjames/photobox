package port

import "time"

//go:generate mockgen -source=cache.go -destination=mock/cache.go -package=mock

type CacheService interface {
	Get(key string) (string, error)
	Set(key string, value string, ttl time.Duration) error
	Delete(key string) error
	DeletePattern(pattern string) error
	GetBytes(key string) ([]byte, error)
	SetBytes(key string, value []byte, ttl time.Duration) error
}
