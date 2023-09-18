package GRPC

import (
	"google.golang.org/grpc"
	"net"
	"user-authorization/internal/config"
	"user-authorization/internal/server/GRPC/pb"
	"user-authorization/storage/authorization"
)

func StartGRPC(cfg config.GRPCConfig) error {
	router := grpc.NewServer(grpc.EmptyServerOption{})
	registerAuthorization(router)
	return startListener(router, cfg)
}

func registerAuthorization(router *grpc.Server) {
	display := authorization.InitDisplay()
	pb.RegisterAuthorizationServer(router, display)
}

func startListener(router *grpc.Server, cfg config.GRPCConfig) error {
	l, err := net.Listen(cfg.Network, cfg.Address)
	if err != nil {
		return err
	}
	return router.Serve(l)
}
