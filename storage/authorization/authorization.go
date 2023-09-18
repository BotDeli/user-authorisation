package authorization

import (
	"context"
	"fmt"
	"user-authorization/internal/server/GRPC/pb"
)

type Display interface {
	Register(ctx context.Context, user *pb.User) (*pb.SessionResponse, error)
	LogIn(ctx context.Context, user *pb.User) (*pb.SessionResponse, error)
	IsAuthenticated(ctx context.Context, session *pb.SessionRequest) (*pb.AuthenticatedSession, error)
}

func InitDisplay() Display {
	return &Authorization{}
}

type Authorization struct {
}

func (a *Authorization) Register(ctx context.Context, user *pb.User) (*pb.SessionResponse, error) {
	fmt.Println("register")
	return &pb.SessionResponse{
		Key:   "",
		Error: "not registered",
	}, nil
}
func (a *Authorization) LogIn(ctx context.Context, user *pb.User) (*pb.SessionResponse, error) {
	fmt.Println("login")
	return &pb.SessionResponse{
		Key:   "",
		Error: "not login",
	}, nil
}
func (a *Authorization) IsAuthenticated(ctx context.Context, session *pb.SessionRequest) (*pb.AuthenticatedSession, error) {
	fmt.Println("isAuth")
	return &pb.AuthenticatedSession{
		Authenticated: true,
		Login:         "test User !!!",
	}, nil
}
