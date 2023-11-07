package user

import (
	"AccountControl/internal/config"
	"AccountControl/pkg/errorHandle"
)

//go:generate go run github.com/vektra/mockery/v2@v2.32.0 --name=Display
type Display interface {
	Close()
	NewUser(email, password string) (string, error)
	IsUser(email string) bool
	AuthenticationUser(email, password string) (string, error)
	ChangePassword(email, password, newPassword string) error
	DeleteUser(id, email, password string) error
	IsVerifiedEmail(email string) (bool, error)
	VerifyEmail(email string) (bool, error)
}

const path = "internal/storage/postgres/user"

func MustInitUserDisplay(cfg *config.PostgresConfig) Display {
	pg, err := initPostgres(cfg)
	if err != nil {
		errorHandle.Fatal(path, "init.go", "MustInitUserDisplay", err)
	}
	return pg
}
