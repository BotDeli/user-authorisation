package GRPC

import (
	"context"
	"errors"
	"user-authorization/internal/server/GRPC/pb"
	"user-authorization/storage/postgres/session"
	"user-authorization/storage/postgres/user"
)

type Authorization struct {
	user    user.Display
	session session.Display
}

var (
	errUserAlreadyRegistered      = errors.New("пользователь уже зарегистрирован")
	errDontCorrectLoginOrPassword = errors.New("некорректный логин или пароль")
	errUserNotAuthenticated       = errors.New("пользователь не авторизован")
)

func InitService(user user.Display, session session.Display) pb.AuthorizationServer {
	return &Authorization{
		user:    user,
		session: session,
	}
}

func (a *Authorization) Register(_ context.Context, user *pb.User) (*pb.SessionData, error) {
	if a.user.IsUser(user.Login) {
		return nil, errUserAlreadyRegistered
	}

	if err := a.user.NewUser(user.Login, user.Password); err != nil {
		return nil, err
	}

	key, err := a.session.NewSession(user.Login)
	if err != nil {
		return nil, err
	}

	return &pb.SessionData{
		Key: key,
	}, nil
}
func (a *Authorization) LogIn(_ context.Context, user *pb.User) (*pb.SessionData, error) {
	if !a.user.IsUser(user.Login) {
		return nil, errDontCorrectLoginOrPassword
	}

	if err := a.user.AuthenticationUser(user.Login, user.Password); err != nil {
		return nil, errDontCorrectLoginOrPassword
	}

	key, err := a.session.NewSession(user.Login)
	if err != nil {
		return nil, err
	}

	return &pb.SessionData{
		Key: key,
	}, nil
}
func (a *Authorization) IsAuthenticated(_ context.Context, session *pb.SessionData) (*pb.AuthenticatedSession, error) {
	login, err := a.session.GetLoginFromSession(session.Key)
	if err != nil {
		return nil, errUserNotAuthenticated
	}

	a.session.UpdateSessionLifeTime(login)

	return &pb.AuthenticatedSession{
		Login: login,
	}, nil
}

func (a *Authorization) ChangePassword(_ context.Context, data *pb.ChangePasswordData) (*pb.Null, error) {
	err := a.user.ChangePassword(data.Login, data.Password, data.NewPassword)
	if err != nil {
		return nil, err
	}

	return &pb.Null{}, nil
}
