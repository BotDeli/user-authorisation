package session

import (
	"github.com/go-redis/redis"
	"user-authorization/internal/config"
	"user-authorization/pkg/errorHandle"
)

type Redis struct {
	Client *redis.Client
}

const path = "internal/storage/redis/session"

func initRedis(cfg *config.RedisConfig) *Redis {
	client := redis.NewClient(&redis.Options{
		Network:  cfg.Network,
		Addr:     cfg.Address,
		Password: cfg.Password,
		DB:       cfg.DB,
	})
	return &Redis{Client: client}
}

func (r *Redis) Close() {
	if err := r.Client.Close(); err != nil {
		errorHandle.Commit(path, "session.go", "Close", err)
	}
}

func (r *Redis) NewSession(login string) (string, error) {
	return "", nil
}

func (r *Redis) GetLoginFromSession(session string) (string, error) {
	return "", nil
}

func (r *Redis) UpdateSessionLifeTime(login string) {
}
