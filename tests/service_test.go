package tests

import (
	"context"
	"errors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"testing"
	"time"
	"user-authorization/internal/config"
	"user-authorization/internal/server/GRPC"
	"user-authorization/internal/server/GRPC/pb"
	"user-authorization/tests/mocks/sessionM"
	"user-authorization/tests/mocks/userM"
)

const (
	login    = "login"
	password = "password"
	testKey  = "SESSION_KEY"

	newPassword = "new_password"
)

var (
	cfg = &config.GRPCConfig{
		Network: "tcp",
		Address: ":50050",
	}

	testError = errors.New("test error")
)

func TestStart(t *testing.T) {
	mockUser := userM.NewDisplay(t)
	mockSession := sessionM.NewDisplay(t)
	service := GRPC.InitService(mockUser, mockSession)

	go func() {
		err := GRPC.StartGRPC(cfg, service)
		if err != nil {
			panic(err)
		}
	}()
	time.Sleep(2 * time.Second)

	client := connectToServer(t)

	testsRegistration(t, mockUser, mockSession, client)
	testsLogIn(t, mockUser, mockSession, client)
	testsIsAuthenticated(t, mockSession, client)
	testsChangePassword(t, mockUser, client)
}

func connectToServer(t *testing.T) pb.AuthorizationClient {
	function := "connectToServer"

	conn, err := grpc.Dial(cfg.Address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		t.Fatal(function, err)
	}
	return pb.NewAuthorizationClient(conn)
}

func testsRegistration(t *testing.T, mockUser *userM.Display, mockSession *sessionM.Display, client pb.AuthorizationClient) {
	testAlreadyRegisteredUserRegister(t, mockUser, client)
	testErrorNewUserRegister(t, mockUser, client)
	testErrorNewSessionRegister(t, mockUser, mockSession, client)
	testDontRegisteredUserRegister(t, mockUser, mockSession, client)
}

func testAlreadyRegisteredUserRegister(t *testing.T, mockUser *userM.Display, client pb.AuthorizationClient) {
	function := "testAlreadyRegisteredUserRegister"

	mockUser.On("IsUser", login).Return(true).Once()

	_, err := client.Register(context.Background(), &pb.User{Login: login, Password: password}, grpc.EmptyCallOption{})
	if err == nil {
		t.Error(function, "expected error, got nil")
	}
}

func testErrorNewUserRegister(t *testing.T, mockUser *userM.Display, client pb.AuthorizationClient) {
	function := "testErrorNewUserRegister"

	mockUser.On("IsUser", login).Return(false).Once()
	mockUser.On("NewUser", login, password).Return(testError).Once()

	_, err := client.Register(context.Background(), &pb.User{Login: login, Password: password}, grpc.EmptyCallOption{})
	if err == nil {
		t.Error(function, "expected error, got nil")
	}
}

func testErrorNewSessionRegister(t *testing.T, mockUser *userM.Display, mockSession *sessionM.Display, client pb.AuthorizationClient) {
	function := "testErrorNewSessionRegister"

	mockUser.On("IsUser", login).Return(false).Once()
	mockUser.On("NewUser", login, password).Return(nil).Once()
	mockSession.On("NewSession", login).Return("", testError).Once()

	_, err := client.Register(context.Background(), &pb.User{Login: login, Password: password}, grpc.EmptyCallOption{})
	if err == nil {
		t.Error(function, "expected error, got nil")
	}
}

func testDontRegisteredUserRegister(t *testing.T, mockUser *userM.Display, mockSession *sessionM.Display, client pb.AuthorizationClient) {
	function := "testDontRegisteredUserRegister"

	mockUser.On("IsUser", login).Return(false).Once()
	mockUser.On("NewUser", login, password).Return(nil).Once()
	mockSession.On("NewSession", login).Return(testKey, nil).Once()

	response, err := client.Register(context.Background(), &pb.User{Login: login, Password: password}, grpc.EmptyCallOption{})
	if err != nil {
		t.Error(function, err)
		return
	}

	if response.Key != testKey {
		t.Errorf("%s expected %s, got %s", function, testKey, response.Key)
	}
}

func testsLogIn(t *testing.T, mockUser *userM.Display, mockSession *sessionM.Display, client pb.AuthorizationClient) {
	testDontRegisteredUserLogIn(t, mockUser, client)
	testFalseAuthenticationLogIn(t, mockUser, client)
	testErrorNewSessionLogIn(t, mockUser, mockSession, client)
	testSuccessfulLogIn(t, mockUser, mockSession, client)
}

func testDontRegisteredUserLogIn(t *testing.T, mockUser *userM.Display, client pb.AuthorizationClient) {
	function := "testDontRegisteredUserLogIn"

	mockUser.On("IsUser", login).Return(false).Once()

	_, err := client.LogIn(context.Background(), &pb.User{Login: login, Password: password}, grpc.EmptyCallOption{})
	if err == nil {
		t.Error(function, "expected error, got nil")
	}
}

func testFalseAuthenticationLogIn(t *testing.T, mockUser *userM.Display, client pb.AuthorizationClient) {
	function := "testFalseAuthenticationLogIn"

	mockUser.On("IsUser", login).Return(true).Once()
	mockUser.On("AuthenticationUser", login, password).Return(testError).Once()

	_, err := client.LogIn(context.Background(), &pb.User{Login: login, Password: password}, grpc.EmptyCallOption{})
	if err == nil {
		t.Error(function, "expected error, got nil")
	}
}

func testErrorNewSessionLogIn(t *testing.T, mockUser *userM.Display, mockSession *sessionM.Display, client pb.AuthorizationClient) {
	function := "testErrorNewSessionLogIn"

	mockUser.On("IsUser", login).Return(true).Once()
	mockUser.On("AuthenticationUser", login, password).Return(nil).Once()
	mockSession.On("NewSession", login).Return("", testError).Once()

	_, err := client.LogIn(context.Background(), &pb.User{Login: login, Password: password}, grpc.EmptyCallOption{})
	if err == nil {
		t.Error(function, "expected error, got nil")
	}
}

func testSuccessfulLogIn(t *testing.T, mockUser *userM.Display, mockSession *sessionM.Display, client pb.AuthorizationClient) {
	function := "testErrorNewSessionLogIn"

	mockUser.On("IsUser", login).Return(true).Once()
	mockUser.On("AuthenticationUser", login, password).Return(nil).Once()
	mockSession.On("NewSession", login).Return(testKey, nil).Once()

	response, err := client.LogIn(context.Background(), &pb.User{Login: login, Password: password}, grpc.EmptyCallOption{})
	if err != nil {
		t.Error(function, err)
		return
	}

	if response.Key != testKey {
		t.Errorf("%s expected %s, got %s", function, testKey, response.Key)
	}
}

func testsIsAuthenticated(t *testing.T, mockSession *sessionM.Display, client pb.AuthorizationClient) {
	testDontCorrectSessionIsAuthenticated(t, mockSession, client)
	testCorrectSessionIsAuthenticated(t, mockSession, client)
}

func testDontCorrectSessionIsAuthenticated(t *testing.T, mockSession *sessionM.Display, client pb.AuthorizationClient) {
	function := "testDontCorrectSessionIsAuthenticated"

	mockSession.On("GetLoginFromSession", testKey).Return("", testError).Once()

	_, err := client.IsAuthenticated(context.Background(), &pb.SessionData{Key: testKey}, grpc.EmptyCallOption{})
	if err == nil {
		t.Error(function, "expected error, got nil")
	}
}

func testCorrectSessionIsAuthenticated(t *testing.T, mockSession *sessionM.Display, client pb.AuthorizationClient) {
	function := "testCorrectSessionIsAuthenticated"

	mockSession.On("GetLoginFromSession", testKey).Return(login, nil).Once()
	mockSession.On("UpdateSessionLifeTime", login).Once()

	response, err := client.IsAuthenticated(context.Background(), &pb.SessionData{Key: testKey}, grpc.EmptyCallOption{})
	if err != nil {
		t.Error(function, err)
		return
	}

	if response.Login != login {
		t.Errorf("%s expected %s, got %s", function, login, response.Login)
	}
}

func testsChangePassword(t *testing.T, mockUser *userM.Display, client pb.AuthorizationClient) {
	testErrorChangePassword(t, mockUser, client)
	testSuccessfulChangePassword(t, mockUser, client)
}

func testErrorChangePassword(t *testing.T, mockUser *userM.Display, client pb.AuthorizationClient) {
	function := "testErrorChangePassword"

	mockUser.On("ChangePassword", login, password, newPassword).Return(testError).Once()

	_, err := client.ChangePassword(context.Background(), &pb.ChangePasswordData{Login: login, Password: password, NewPassword: newPassword}, grpc.EmptyCallOption{})
	if err == nil {
		t.Error(function, "expected error, got nil")
	}
}

func testSuccessfulChangePassword(t *testing.T, mockUser *userM.Display, client pb.AuthorizationClient) {
	function := "testErrorChangePassword"

	mockUser.On("ChangePassword", login, password, newPassword).Return(nil).Once()

	_, err := client.ChangePassword(context.Background(), &pb.ChangePasswordData{Login: login, Password: password, NewPassword: newPassword}, grpc.EmptyCallOption{})
	if err != nil {
		t.Error(function, err)
	}
}
