package redis

import (
	"context"
	"time"

	redis "github.com/redis/go-redis/v9"
	"github.com/tingin/base/config"
	"github.com/tingin/base/env"
	"github.com/tingin/base/patterns/singleton"
)

var singletonMap = singleton.NewSingletonMap[string, RedisClient]()

var defaultKey = "Default"

func init() {
	singletonMap.AddFactory(defaultKey, defaultInstance)
}

type Redis string

const (
	Addr         Redis = "Addr"
	PassWord     Redis = "PassWord"
	DB           Redis = "DB"
	PoolSize     Redis = "PoolSize"
	MinIdleConns Redis = "MinIdleConns"
	TimeOut      Redis = "TimeOut"
)

func defaultInstance() *RedisClient {
	addr := config.GetEnv(string(Addr), "127.0.0.1:6379")
	password := config.GetEnv(string(PassWord), "")
	db := env.GetEnvAsInt(string(DB), 0)
	poolsize := env.GetEnvAsInt(string(PoolSize), 10)
	minIdleConns := env.GetEnvAsInt(string(MinIdleConns), 10)
	timeout := env.GetEnvAsInt(string(TimeOut), 30)
	return NewRedisClient(addr, password, db, poolsize, minIdleConns, timeout)
}

func Default() *RedisClient {
	return singletonMap.GetInstance(defaultKey)
}

type RedisClient struct {
	Client *redis.Client
	Ctx    context.Context
}

func NewRedisClient(addr string, password string, db int, poolsize int, minIdleConns int, timeout int) *RedisClient {
	rdb := redis.NewClient(&redis.Options{
		Addr:         addr,
		Password:     password,
		DB:           db,
		PoolSize:     poolsize,
		MinIdleConns: minIdleConns,
		PoolTimeout:  time.Duration(timeout) * time.Second,
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
