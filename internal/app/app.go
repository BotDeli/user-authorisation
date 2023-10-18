package app

import (
	"AccountControl/internal/config"
	"AccountControl/internal/server/GRPC"
	"AccountControl/internal/storage/postgres/user"
	"AccountControl/internal/storage/redis/session"
	"AccountControl/pkg/errorHandle"
)

func StartApp() {
	cfg := config.MustReadConfig()

	displayU := user.MustInitUserDisplay(cfg.Postgres)
	defer displayU.Close()

	displayS := session.MustInitSessionDisplay(cfg.Redis)
	defer displayS.Close()

	service := GRPC.InitService(displayU, displayS)

	errorHandle.Fatal(
		"internal/app",
		"app.go",
		"StartApp",
		GRPC.StartGRPC(cfg.Grpc, service),
	)
}
