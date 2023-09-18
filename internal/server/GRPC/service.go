package GRPC

import (
	"context"
	"fmt"
	"user-authorization/internal/server/GRPC/pb"
	"user-authorization/storage/postgres/session"
	"user-authorization/storage/postgres/user"
)

type Authorization struct {
	user    user.Display
	session session.Display
}

func MustInitService() pb.AuthorizationServer {
	return &Authorization{}
}

func (a *Authorization) Register(ctx context.Context, user *pb.User) (*pb.SessionResponse, error) {
	fmt.Println(user)
	return &pb.SessionResponse{
		Key:   "",
		Error: "not registered",
	}, nil
}
func (a *Authorization) LogIn(ctx context.Context, user *pb.User) (*pb.SessionResponse, error) {
	fmt.Println(user)
	return &pb.SessionResponse{
		Key:   "",
		Error: "not login",
	}, nil
}
func (a *Authorization) IsAuthenticated(ctx context.Context, session *pb.SessionRequest) (*pb.AuthenticatedSession, error) {
	fmt.Println(session)
	return &pb.AuthenticatedSession{
		Authenticated: true,
		Login:         "test User !!!",
	}, nil
}
