package GRPC

import (
	"google.golang.org/grpc"
	"net"
	"user-authorization/internal/config"
	"user-authorization/internal/server/GRPC/pb"
)

func StartGRPC(cfg *config.GRPCConfig, service pb.AuthorizationServer) error {
	router := grpc.NewServer(grpc.EmptyServerOption{})
	pb.RegisterAuthorizationServer(router, service)
	return startListener(router, cfg)
}

func startListener(router *grpc.Server, cfg *config.GRPCConfig) error {
	l, err := net.Listen(cfg.Network, cfg.Address)
	if err != nil {
		return err
	}
	return router.Serve(l)
}
