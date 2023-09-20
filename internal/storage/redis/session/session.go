package session

import (
	"github.com/go-redis/redis"
	"time"
	"user-authorization/internal/config"
	"user-authorization/pkg/UUIDGenerator"
	"user-authorization/pkg/errorHandle"
)

type redisClient interface {
	Close() error
	Set(key string, value interface{}, expiration time.Duration) *redis.StatusCmd
	Get(key string) *redis.StringCmd
	Expire(key string, expiration time.Duration) *redis.BoolCmd
}

type Redis struct {
	Client   redisClient
	Lifetime time.Duration
}

const path = "internal/storage/redis/session"

func initRedis(cfg *config.RedisConfig) *Redis {
	client := redis.NewClient(&redis.Options{
		Network:  cfg.Network,
		Addr:     cfg.Address,
		Password: cfg.Password,
		DB:       cfg.DB,
	})
	return &Redis{Client: client, Lifetime: cfg.LifetimeWrite}
}

func (r *Redis) Close() {
	if err := r.Client.Close(); err != nil {
		errorHandle.Commit(path, "session.go", "Close", err)
	}
}

func (r *Redis) NewSession(login string) (string, error) {
	key := UUIDGenerator.NewUUID()
	cmd := r.Client.Set(key, login, r.Lifetime)
	return key, cmd.Err()
}

func (r *Redis) GetLoginFromSession(key string) (string, error) {
	cmd := r.Client.Get(key)
	return cmd.Result()
}

func (r *Redis) UpdateSessionLifeTime(key string) {
	r.Client.Expire(key, r.Lifetime)
}
