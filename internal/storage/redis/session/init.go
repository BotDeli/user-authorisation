package session

import (
	"AccountControl/internal/config"
	"AccountControl/pkg/errorHandle"
)

//go:generate go run github.com/vektra/mockery/v2@v2.32.0 --name=Display
type Display interface {
	NewSession(id string) (string, error)
	GetIdFromSession(session string) (string, error)
	UpdateSessionLifeTime(key string)
	DeleteSession(key string)
	Close()
}

const path = "/internal/storage/redis/session"

func MustInitSessionDisplay(cfg *config.RedisConfig) Display {
	r, err := initRedis(cfg)
	if err != nil {
		errorHandle.Fatal(path, "init.go", "MustInitSessionDisplay", err)
	}
	return r
}
