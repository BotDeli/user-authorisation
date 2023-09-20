package session

import (
	"user-authorization/internal/config"
)

//go:generate go run github.com/vektra/mockery/v2@v2.32.0 --name=Display
type Display interface {
	NewSession(login string) (string, error)
	GetLoginFromSession(session string) (string, error)
	UpdateSessionLifeTime(key string)
	Close()
}

func MustInitSessionDisplay(cfg *config.RedisConfig) Display {
	return initRedis(cfg)
}
