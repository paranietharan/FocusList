package cache

import (
	"context"
	"sync"

	"github.com/redis/go-redis/v9"
)

var (
	rdb  *redis.Client
	once sync.Once
)

type RedisConfig struct {
	Addr     string
	Username string
	Password string
	DB       int
}

func GetRedisClient(cfg RedisConfig) *redis.Client {
	once.Do(func() {
		rdb = redis.NewClient(&redis.Options{
			Addr:     cfg.Addr,
			Username: cfg.Username,
			Password: cfg.Password,
			DB:       cfg.DB,
		})
	})
	return rdb
}

type CacheService struct {
	redisClient *redis.Client
	ctx         context.Context
}

func NewCacheService(cfg RedisConfig, ctx context.Context) *CacheService {
	client := GetRedisClient(cfg)
	return &CacheService{
		redisClient: client,
		ctx:         ctx,
	}
}

func (cs *CacheService) Set(key string, value string) error {
	return cs.redisClient.Set(cs.ctx, key, value, 0).Err()
}

func (cs *CacheService) Get(key string) (string, error) {
	val, err := cs.redisClient.Get(cs.ctx, key).Result()
	if err == redis.Nil {
		return "", nil
	} else if err != nil {
		return "", err
	}
	return val, nil
}
