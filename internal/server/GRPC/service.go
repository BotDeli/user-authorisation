package GRPC

import (
	"AccountControl/internal/server/GRPC/pb"
	"AccountControl/internal/storage/postgres/user"
	"AccountControl/internal/storage/redis/session"
	"context"
	"errors"
)

type AccountControl struct {
	displayU user.Display
	displayS session.Display
}

var (
	errUserAlreadyRegistered      = errors.New("пользователь уже зарегистрирован")
	errDontCorrectEmailOrPassword = errors.New("некорректный логин или пароль")
	errUserNotAuthorized          = errors.New("пользователь не авторизован")
	errServerError                = errors.New("ошибка на сервере, попробуйте позже")
)

func InitService(displayU user.Display, displayS session.Display) pb.AccountControlServer {
	return &AccountControl{
		displayU: displayU,
		displayS: displayS,
	}
}

func (a *AccountControl) RegistrationAccount(_ context.Context, user *pb.User) (*pb.SessionData, error) {
	if a.displayU.IsUser(user.Email) {
		return nil, errUserAlreadyRegistered
	}

	id, err := a.displayU.NewUser(user.Email, user.Password)
	if err != nil {
		return nil, errServerError
	}

	key, err := a.displayS.NewSession(id)
	if err != nil {
		return nil, errServerError
	}

	return &pb.SessionData{
		Id:  id,
		Key: key,
	}, nil
}

func (a *AccountControl) AuthorizationAccount(_ context.Context, user *pb.User) (*pb.SessionData, error) {
	if !a.displayU.IsUser(user.Email) {
		return nil, errDontCorrectEmailOrPassword
	}

	id, err := a.displayU.AuthenticationUser(user.Email, user.Password)
	if err != nil {
		return nil, errDontCorrectEmailOrPassword
	}

	key, err := a.displayS.NewSession(id)
	if err != nil {
		return nil, errServerError
	}

	return &pb.SessionData{
		Id:  id,
		Key: key,
	}, nil
}

func (a *AccountControl) ChangePasswordAccount(_ context.Context, data *pb.ChangePasswordData) (*pb.Null, error) {
	err := a.displayU.ChangePassword(data.Email, data.Password, data.NewPassword)
	if err != nil {
		return nil, errUserNotAuthorized
	}

	return &pb.Null{}, nil
}

func (a *AccountControl) DeleteAccount(_ context.Context, info *pb.FullInfoUser) (*pb.Null, error) {
	err := a.displayU.DeleteUser(info.Id, info.Email, info.Password)
	if err != nil {
		return nil, errUserNotAuthorized
	}

	return &pb.Null{}, nil
}

func (a *AccountControl) IsAuthorizedSessionData(_ context.Context, session *pb.SessionData) (*pb.AccountID, error) {
	id, err := a.displayS.GetIdFromSession(session.Key)
	if err != nil {
		return nil, errUserNotAuthorized
	}

	a.displayS.UpdateSessionLifeTime(id)

	return &pb.AccountID{
		Id: id,
	}, nil
}

func (a *AccountControl) DeleteSessionData(_ context.Context, session *pb.SessionData) (*pb.Null, error) {
	a.displayS.DeleteSession(session.Key)
	return &pb.Null{}, nil
}

func (a *AccountControl) IsVerifiedEmail(_ context.Context, emailData *pb.EmailData) (*pb.VerifiedEmailData, error) {
	is, err := a.displayU.IsVerifiedEmail(emailData.Email)
	if err != nil {
		return nil, errUserNotAuthorized
	}
	return &pb.VerifiedEmailData{
		IsVerified: is,
	}, nil
}

func (a *AccountControl) VerifyEmail(_ context.Context, emailData *pb.EmailData) (*pb.VerifiedEmailData, error) {
	is, err := a.displayU.VerifyEmail(emailData.Email)
	if err != nil {
		return nil, errUserNotAuthorized
	}
	return &pb.VerifiedEmailData{
		IsVerified: is,
	}, nil
}
