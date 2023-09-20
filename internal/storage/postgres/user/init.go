package user

import (
	"user-authorization/internal/config"
	"user-authorization/pkg/errorHandle"
)

//go:generate go run github.com/vektra/mockery/v2@v2.32.0 --name=Display
type Display interface {
	NewUser(login, password string) error
	IsUser(login string) bool
	AuthenticationUser(login, password string) error
	ChangePassword(login, password, newPassword string) error
	Close()
}

const path = "internal/storage/postgres/user"

func MustInitUserDisplay(cfg *config.PostgresConfig) Display {
	pg, err := initPostgres(cfg)
	if err != nil {
		errorHandle.Fatal(path, "init.go", "MustInitUserDisplay", err)
	}
	return pg
}
