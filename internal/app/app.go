package app

import (
	"log"
	"user-authorization/internal/config"
	"user-authorization/internal/server/GRPC"
)

func StartApp() {
	cfg := config.MustReadConfig()
	service := GRPC.MustInitService()
	log.Fatal(GRPC.StartGRPC(cfg.Grpc, service))
}
