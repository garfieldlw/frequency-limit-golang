package redis

import (
	"context"
	"errors"
	"github.com/redis/go-redis/v9"
	"sync"
	"time"
)

var redisClient *Redis
var lock *sync.Mutex = &sync.Mutex{}

// Redis provides a cache backed by a Redis server.
type Redis struct {
	Config *redis.Options
	Client *redis.Client
}

// New returns an initialized Redis cache object.
func New(config *redis.Options) *Redis {
	client := redis.NewClient(config)
	return &Redis{Config: config, Client: client}
}

func (r *Redis) TTL(ctx context.Context, key string) (time.Duration, error) {
	return r.Client.TTL(ctx, key).Result()
}

func (r *Redis) HGet(ctx context.Context, key, field string) (string, error) {
	return r.Client.HGet(ctx, key, field).Result()
}

func (r *Redis) HGetAll(ctx context.Context, key string) (map[string]string, error) {
	return r.Client.HGetAll(ctx, key).Result()
}

func (r *Redis) HSet(ctx context.Context, key, field, value string, expire time.Duration) error {
	err := r.Client.HSet(ctx, key, field, value).Err()
	if err == nil && expire > 0 {
		err = r.Client.Expire(ctx, key, expire).Err()
	}
	return err
}

func (r *Redis) HSetNX(ctx context.Context, key, field, value string, expire time.Duration) error {
	err := r.Client.HSetNX(ctx, key, field, value).Err()
	if err == nil && expire > 0 {
		err = r.Client.Expire(ctx, key, expire).Err()
	}
	return err
}

func (r *Redis) HDel(ctx context.Context, key string, field ...string) error {
	return r.Client.HDel(ctx, key, field...).Err()
}

func (r *Redis) Ping(ctx context.Context) error {
	return r.Client.Ping(ctx).Err()
}

func GetRedis() (*Redis, error) {
	if redisClient != nil {
		return redisClient, nil
	}

	lock.Lock()
	defer lock.Unlock()

	conf := getRedisConfig()
	if conf == nil {
		return nil, errors.New("redis config is invalid")
	}

	if redisClient != nil {
		return redisClient, nil
	}

	redisClient = New(&redis.Options{
		Network:     "tcp",
		Password:    conf.Password,
		Addr:        conf.Address,
		DB:          int(conf.DB),
		DialTimeout: time.Second,
		PoolSize:    50,
		PoolTimeout: time.Second,
	})

	return redisClient, nil
}

func InitRedis() error {
	if redisClient != nil {
		return nil
	}

	lock.Lock()
	defer lock.Unlock()

	conf := getRedisConfig()
	if conf == nil {
		return errors.New("redis config is invalid")
	}

	if redisClient != nil {
		return nil
	}

	redisClient = New(&redis.Options{
		Network:     "tcp",
		Password:    conf.Password,
		Addr:        conf.Address,
		DB:          int(conf.DB),
		DialTimeout: time.Second,
		PoolSize:    50,
		PoolTimeout: time.Second,
	})

	return nil
}

func Ping(ctx context.Context) error {
	c, err := GetRedis()
	if err != nil {
		return err
	}

	err = c.Ping(ctx)
	if err != nil {
		return err
	}

	return nil
}

type RedisConfigItem struct {
	Address  string `json:"address"`
	Password string `json:"password"`
	DB       int32  `json:"db"`
}

func getRedisConfig() *RedisConfigItem {
	return &RedisConfigItem{
		Address:  "127.0.0.1:6379",
		Password: "",
		DB:       0,
	}
}
