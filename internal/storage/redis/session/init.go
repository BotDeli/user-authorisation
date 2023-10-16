package session

import (
	"user-authorization/internal/config"
	"user-authorization/pkg/errorHandle"
)

//go:generate go run github.com/vektra/mockery/v2@v2.32.0 --name=Display
type Display interface {
	NewSession(login string) (string, error)
	GetLoginFromSession(session string) (string, error)
	UpdateSessionLifeTime(key string)
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
