package redis

import (
	"context"
	"strconv"
	"time"

	redis "github.com/redis/go-redis/v9"
	"github.com/tingin/base/config"
	"github.com/tingin/base/patterns/singleton"
)

var singletonMap = singleton.NewSingletonMap[string, RedisClient]()

var defaultKey = "Default"

func init() {
	singletonMap.AddFactory(defaultKey, defaultInstance)
}

func defaultInstance() *RedisClient {
	addr := config.GetEnv("RedisAddr", "127.0.0.1:6379")
	password := config.GetEnv("RedisPasswd", "")
	dbcfg := config.GetEnv("Redisdb", "")
	db, err := strconv.Atoi(dbcfg)
	if err != nil {
		db = 0
	}
	return NewRedisClient(addr, password, db)
}

func Default() *RedisClient {
	return singletonMap.GetInstance(defaultKey)
}

type RedisClient struct {
	Client *redis.Client
	Ctx    context.Context
}

func NewRedisClient(addr string, password string, db int) *RedisClient {
	rdb := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       db,
	})

	return &RedisClient{
		Client: rdb,
		Ctx:    context.Background(),
	}
}

func (r *RedisClient) Set(key string, value interface{}, expiration time.Duration) error {
	err := r.Client.Set(r.Ctx, key, value, expiration).Err()
	return err
}

func (r *RedisClient) Get(key string) (string, error) {
	val, err := r.Client.Get(r.Ctx, key).Result()
	if err != nil {
		return "", err
	}
	return val, nil
}

func (r *RedisClient) Del(key string) error {
	err := r.Client.Del(r.Ctx, key).Err()
	return err
}
