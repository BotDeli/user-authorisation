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

/*
	POST register:
		check login in users
		password => hash
		save login, hash
		return error

	POST logIn:
		check login in users
		password => hash
		equals insert hash and saved hash
		create new sessionM in sessions is lifetime 24h
		return sessionM

	POST isAuthenticated:
		check key in sessionM
		get login from sessionM
		update sessionM lifetime to 24h
		return true, login

	UPDATE changePassword:
		check login in users
		lastPassword => hash
		equals last hash and saved hash
		newPassword => hash
		update saved hash to new hash
		return error
*/
