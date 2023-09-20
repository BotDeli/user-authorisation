package app

import (
	"user-authorization/internal/config"
	"user-authorization/internal/server/GRPC"
	"user-authorization/internal/storage/postgres/user"
	"user-authorization/internal/storage/redis/session"
	"user-authorization/pkg/errorHandle"
)

func StartApp() {
	cfg := config.MustReadConfig()

	userDisplay := user.MustInitUserDisplay(cfg.Postgres)
	defer userDisplay.Close()
	sessionDisplay := session.MustInitSessionDisplay(cfg.Redis)
	defer sessionDisplay.Close()

	service := GRPC.InitService(userDisplay, sessionDisplay)

	errorHandle.Fatal(
		"internal/app",
		"app.go",
		"StartApp",
		GRPC.StartGRPC(cfg.Grpc, service),
	)
}
