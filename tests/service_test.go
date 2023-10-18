package tests

import (
	"AccountControl/internal/config"
	"AccountControl/internal/server/GRPC"
	"AccountControl/internal/server/GRPC/pb"
	"AccountControl/tests/mocks/sessionM"
	"AccountControl/tests/mocks/userM"
	"context"
	"errors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"testing"
	"time"
)

const (
	email       = "email"
	password    = "password"
	newPassword = "new_password"
	testKey     = "SESSION_KEY"
	testID      = "1234567"
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
	testsDeleteAccount(t, mockUser, client)
	testDeleteSession(t, mockSession, client)
}

func connectToServer(t *testing.T) pb.AccountControlClient {
	function := "connectToServer"

	conn, err := grpc.Dial(cfg.Address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		t.Fatal(function, err)
	}
	return pb.NewAccountControlClient(conn)
}

func testsRegistration(t *testing.T, mockUser *userM.Display, mockSession *sessionM.Display, client pb.AccountControlClient) {
	testAlreadyRegisteredUserRegister(t, mockUser, client)
	testErrorNewUserRegister(t, mockUser, client)
	testErrorNewSessionRegister(t, mockUser, mockSession, client)
	testDontRegisteredUserRegister(t, mockUser, mockSession, client)
}

func testAlreadyRegisteredUserRegister(t *testing.T, mockUser *userM.Display, client pb.AccountControlClient) {
	function := "testAlreadyRegisteredUserRegister"

	mockUser.On("IsUser", email).Return(true).Once()

	_, err := client.RegistrationAccount(context.Background(), &pb.User{Email: email, Password: password}, grpc.EmptyCallOption{})
	if err == nil {
		t.Error(function, "expected error, got nil")
	}
}

func testErrorNewUserRegister(t *testing.T, mockUser *userM.Display, client pb.AccountControlClient) {
	function := "testErrorNewUserRegister"

	mockUser.On("IsUser", email).Return(false).Once()
	mockUser.On("NewUser", email, password).Return("", testError).Once()

	_, err := client.RegistrationAccount(context.Background(), &pb.User{Email: email, Password: password}, grpc.EmptyCallOption{})
	if err == nil {
		t.Error(function, "expected error, got nil")
	}
}

func testErrorNewSessionRegister(t *testing.T, mockUser *userM.Display, mockSession *sessionM.Display, client pb.AccountControlClient) {
	function := "testErrorNewSessionRegister"

	mockUser.On("IsUser", email).Return(false).Once()
	mockUser.On("NewUser", email, password).Return(testID, nil).Once()
	mockSession.On("NewSession", testID).Return("", testError).Once()

	_, err := client.RegistrationAccount(context.Background(), &pb.User{Email: email, Password: password}, grpc.EmptyCallOption{})
	if err == nil {
		t.Error(function, "expected error, got nil")
	}
}

func testDontRegisteredUserRegister(t *testing.T, mockUser *userM.Display, mockSession *sessionM.Display, client pb.AccountControlClient) {
	function := "testDontRegisteredUserRegister"

	mockUser.On("IsUser", email).Return(false).Once()
	mockUser.On("NewUser", email, password).Return(testID, nil).Once()
	mockSession.On("NewSession", testID).Return(testKey, nil).Once()

	response, err := client.RegistrationAccount(context.Background(), &pb.User{Email: email, Password: password}, grpc.EmptyCallOption{})
	if err != nil {
		t.Error(function, err)
		return
	}

	if response.Key != testKey {
		t.Errorf("%s expected %s, got %s", function, testKey, response.Key)
	}
}

func testsLogIn(t *testing.T, mockUser *userM.Display, mockSession *sessionM.Display, client pb.AccountControlClient) {
	testDontRegisteredUserLogIn(t, mockUser, client)
	testFalseAuthenticationLogIn(t, mockUser, client)
	testErrorNewSessionLogIn(t, mockUser, mockSession, client)
	testSuccessfulLogIn(t, mockUser, mockSession, client)
}

func testDontRegisteredUserLogIn(t *testing.T, mockUser *userM.Display, client pb.AccountControlClient) {
	function := "testDontRegisteredUserLogIn"

	mockUser.On("IsUser", email).Return(false).Once()

	_, err := client.AuthorizationAccount(context.Background(), &pb.User{Email: email, Password: password}, grpc.EmptyCallOption{})
	if err == nil {
		t.Error(function, "expected error, got nil")
	}
}

func testFalseAuthenticationLogIn(t *testing.T, mockUser *userM.Display, client pb.AccountControlClient) {
	function := "testFalseAuthenticationLogIn"

	mockUser.On("IsUser", email).Return(true).Once()
	mockUser.On("AuthenticationUser", email, password).Return("", testError).Once()

	_, err := client.AuthorizationAccount(context.Background(), &pb.User{Email: email, Password: password}, grpc.EmptyCallOption{})
	if err == nil {
		t.Error(function, "expected error, got nil")
	}
}

func testErrorNewSessionLogIn(t *testing.T, mockUser *userM.Display, mockSession *sessionM.Display, client pb.AccountControlClient) {
	function := "testErrorNewSessionLogIn"

	mockUser.On("IsUser", email).Return(true).Once()
	mockUser.On("AuthenticationUser", email, password).Return(testID, nil).Once()
	mockSession.On("NewSession", testID).Return("", testError).Once()

	_, err := client.AuthorizationAccount(context.Background(), &pb.User{Email: email, Password: password}, grpc.EmptyCallOption{})
	if err == nil {
		t.Error(function, "expected error, got nil")
	}
}

func testSuccessfulLogIn(t *testing.T, mockUser *userM.Display, mockSession *sessionM.Display, client pb.AccountControlClient) {
	function := "testErrorNewSessionLogIn"

	mockUser.On("IsUser", email).Return(true).Once()
	mockUser.On("AuthenticationUser", email, password).Return(testID, nil).Once()
	mockSession.On("NewSession", testID).Return(testKey, nil).Once()

	response, err := client.AuthorizationAccount(context.Background(), &pb.User{Email: email, Password: password}, grpc.EmptyCallOption{})
	if err != nil {
		t.Error(function, err)
		return
	}

	if response.Key != testKey {
		t.Errorf("%s expected %s, got %s", function, testKey, response.Key)
	}
}

func testsIsAuthenticated(t *testing.T, mockSession *sessionM.Display, client pb.AccountControlClient) {
	testDontCorrectSessionIsAuthenticated(t, mockSession, client)
	testCorrectSessionIsAuthenticated(t, mockSession, client)
}

func testDontCorrectSessionIsAuthenticated(t *testing.T, mockSession *sessionM.Display, client pb.AccountControlClient) {
	function := "testDontCorrectSessionIsAuthenticated"

	mockSession.On("GetIdFromSession", testKey).Return(testID, testError).Once()

	_, err := client.IsAuthorizedSessionData(context.Background(), &pb.SessionData{Id: testID, Key: testKey}, grpc.EmptyCallOption{})
	if err == nil {
		t.Error(function, "expected error, got nil")
	}
}

func testCorrectSessionIsAuthenticated(t *testing.T, mockSession *sessionM.Display, client pb.AccountControlClient) {
	function := "testCorrectSessionIsAuthenticated"

	mockSession.On("GetIdFromSession", testKey).Return(testID, nil).Once()
	mockSession.On("UpdateSessionLifeTime", testID).Once()

	response, err := client.IsAuthorizedSessionData(context.Background(), &pb.SessionData{Id: testID, Key: testKey}, grpc.EmptyCallOption{})
	if err != nil {
		t.Error(function, err)
		return
	}

	if response.Id != testID {
		t.Errorf("%s expected %s, got %s", function, testID, response.Id)
	}
}

func testsChangePassword(t *testing.T, mockUser *userM.Display, client pb.AccountControlClient) {
	testErrorChangePassword(t, mockUser, client)
	testSuccessfulChangePassword(t, mockUser, client)
}

func testErrorChangePassword(t *testing.T, mockUser *userM.Display, client pb.AccountControlClient) {
	function := "testErrorChangePassword"

	mockUser.On("ChangePassword", email, password, newPassword).Return(testError).Once()

	_, err := client.ChangePasswordAccount(context.Background(), &pb.ChangePasswordData{Email: email, Password: password, NewPassword: newPassword}, grpc.EmptyCallOption{})
	if err == nil {
		t.Error(function, "expected error, got nil")
	}
}

func testSuccessfulChangePassword(t *testing.T, mockUser *userM.Display, client pb.AccountControlClient) {
	function := "testErrorChangePassword"

	mockUser.On("ChangePassword", email, password, newPassword).Return(nil).Once()

	_, err := client.ChangePasswordAccount(context.Background(), &pb.ChangePasswordData{Email: email, Password: password, NewPassword: newPassword}, grpc.EmptyCallOption{})
	if err != nil {
		t.Error(function, err)
	}
}

func testsDeleteAccount(t *testing.T, mockUser *userM.Display, client pb.AccountControlClient) {
	testErrorDeleteAccount(t, mockUser, client)
	testSuccessfulDeleteAccount(t, mockUser, client)
}

func testErrorDeleteAccount(t *testing.T, mockUser *userM.Display, client pb.AccountControlClient) {
	function := "testErrorDeleteAccount"

	mockUser.On("DeleteUser", testID, email, password).Return(testError).Once()

	_, err := client.DeleteAccount(context.Background(), &pb.FullInfoUser{Id: testID, Email: email, Password: password}, grpc.EmptyCallOption{})
	if err == nil {
		t.Error(function, err)
	}
}

func testSuccessfulDeleteAccount(t *testing.T, mockUser *userM.Display, client pb.AccountControlClient) {
	function := "testErrorDeleteAccount"

	mockUser.On("DeleteUser", testID, email, password).Return(nil).Once()

	_, err := client.DeleteAccount(context.Background(), &pb.FullInfoUser{Id: testID, Email: email, Password: password}, grpc.EmptyCallOption{})
	if err != nil {
		t.Error(function, err)
	}
}

func testDeleteSession(t *testing.T, mockSession *sessionM.Display, client pb.AccountControlClient) {
	function := "testDeleteSession"

	mockSession.On("DeleteSession", testKey).Return(nil).Once()

	_, err := client.DeleteSessionData(context.Background(), &pb.SessionData{Id: testID, Key: testKey}, grpc.EmptyCallOption{})
	if err != nil {
		t.Error(function, err)
	}

	mockSession.On("DeleteSession", testKey).Return(testError).Once()

	_, err = client.DeleteSessionData(context.Background(), &pb.SessionData{Id: testID, Key: testKey}, grpc.EmptyCallOption{})
	if err != nil {
		t.Error(function, err)
	}
}
