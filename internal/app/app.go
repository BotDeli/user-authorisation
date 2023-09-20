package app

import (
	"user-authorization/internal/config"
	"user-authorization/internal/server/GRPC"
	"user-authorization/pkg/errorHandle"
)

func StartApp() {
	cfg := config.MustReadConfig()
	service := GRPC.InitService()
	errorHandle.Fatal(
		"internal/app",
		"app.go",
		"StartApp",
		GRPC.StartGRPC(cfg.Grpc, service),
	)
}
