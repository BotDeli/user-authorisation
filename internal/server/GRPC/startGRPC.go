package GRPC

import (
	"AccountControl/internal/config"
	"AccountControl/internal/server/GRPC/pb"
	"google.golang.org/grpc"
	"net"
)

func StartGRPC(cfg *config.GRPCConfig, service pb.AccountControlServer) error {
	router := grpc.NewServer(grpc.EmptyServerOption{})
	pb.RegisterAccountControlServer(router, service)
	return startListener(router, cfg)
}

func startListener(router *grpc.Server, cfg *config.GRPCConfig) error {
	l, err := net.Listen(cfg.Network, cfg.Address)
	if err != nil {
		return err
	}

	return router.Serve(l)
}
