package session

import (
	"AccountControl/internal/config"
	"AccountControl/pkg/errorHandle"
	"AccountControl/pkg/generator"
	"github.com/go-redis/redis"
	"time"
)

type redisClient interface {
	Close() error
	Set(key string, value interface{}, expiration time.Duration) *redis.StatusCmd
	Get(key string) *redis.StringCmd
	Expire(key string, expiration time.Duration) *redis.BoolCmd
	Del(keys ...string) *redis.IntCmd
}

type Redis struct {
	Client   redisClient
	Lifetime time.Duration
}

func initRedis(cfg *config.RedisConfig) (*Redis, error) {
	client := getNewClient(cfg)
	return &Redis{Client: client, Lifetime: cfg.LifetimeWrite}, client.Ping().Err()
}

func getNewClient(cfg *config.RedisConfig) *redis.Client {
	return redis.NewClient(&redis.Options{
		Network:  cfg.Network,
		Addr:     cfg.Address,
		Password: cfg.Password,
		DB:       cfg.DB,
	})
}

func (r *Redis) Close() {
	if err := r.Client.Close(); err != nil {
		errorHandle.Commit(path, "session.go", "Close", err)
	}
}

func (r *Redis) NewSession(id string) (string, error) {
	key := generator.NewUUIDDigitsLetters()
	cmd := r.Client.Set(key, id, r.Lifetime)
	return key, cmd.Err()
}

func (r *Redis) GetIdFromSession(key string) (string, error) {
	cmd := r.Client.Get(key)
	return cmd.Result()
}

func (r *Redis) UpdateSessionLifeTime(key string) {
	r.Client.Expire(key, r.Lifetime)
}

func (r *Redis) DeleteSession(key string) {
	r.Client.Del(key)
}
