package tests

import (
	"AccountControl/internal/storage/postgres/user"
	"AccountControl/pkg/hasher"
	"github.com/DATA-DOG/go-sqlmock"
	"testing"
)

const (
	invalidPassword = "invalid password"
)

var (
	emptyResult = sqlmock.NewErrorResult(nil)
)

func TestErrorGenerateUniqueIDNewUser(t *testing.T) {
	initMocks := func(mock sqlmock.Sqlmock) {
		mock.ExpectQuery(`SELECT`).WithArgs(testID).WillReturnError(testError)
	}

	testMocks := func(t *testing.T, pg user.Postgres) {
		id, err := pg.NewUser(email, password)
		isEmptyStr(t, id)
		errorIsNotNil(t, err)
	}

	testPostgresMock(t, initMocks, testMocks)
}

func isEmptyStr(t *testing.T, str string) {
	if len(str) != 0 {
		t.Error("expected empty string, got", str)
	}
}

func errorIsNotNil(t *testing.T, err error) {
	if err == nil {
		t.Error("expected error, got", err)
	}
}

func TestErrorNewUser(t *testing.T) {
	initMocks := func(mock sqlmock.Sqlmock) {
		rows := getCountRowsNull()
		mock.ExpectQuery(`SELECT`).WithArgs(sqlmock.AnyArg()).WillReturnRows(rows).WillReturnError(nil)
		mock.ExpectExec(`INSERT INTO`).WithArgs(sqlmock.AnyArg(), email, hasher.Hashing(password)).WillReturnError(testError)
	}

	testMocks := func(t *testing.T, pg user.Postgres) {
		id, err := pg.NewUser(email, password)
		isEmptyStr(t, id)
		errorIsTestError(t, err)
	}

	testPostgresMock(t, initMocks, testMocks)
}

func getCountRowsNull() *sqlmock.Rows {
	rows := sqlmock.NewRows([]string{"COUNT(*)"})
	rows.AddRow(0)
	return rows
}

func errorIsTestError(t *testing.T, err error) {
	if err != testError {
		t.Error("expected test error, got", err)
	}
}

func testPostgresMock(t *testing.T, initMocks func(sqlmock.Sqlmock), testMocks func(t *testing.T, pg user.Postgres)) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}

	initMocks(mock)

	pg := user.Postgres{DB: db}

	testMocks(t, pg)
}

func TestSuccessfulNewUser(t *testing.T) {
	initMocks := func(mock sqlmock.Sqlmock) {
		rows := getCountRowsNull()
		mock.ExpectQuery(`SELECT`).WithArgs(sqlmock.AnyArg()).WillReturnRows(rows).WillReturnError(nil)
		mock.ExpectExec(`INSERT INTO`).WithArgs(sqlmock.AnyArg(), email, sqlmock.AnyArg()).WillReturnResult(emptyResult).WillReturnError(nil)
	}

	testMocks := func(t *testing.T, pg user.Postgres) {
		id, err := pg.NewUser(email, password)
		isDontEmptyStr(t, id)
		errorIsNil(t, err)
	}

	testPostgresMock(t, initMocks, testMocks)
}

func isDontEmptyStr(t *testing.T, str string) {
	if len(str) == 0 {
		t.Error("expected dont empty string, got", str)
	}
}

func errorIsNil(t *testing.T, err error) {
	if err != nil {
		t.Error("expected nil, got", err)
	}
}

func TestErrorQueryIsUser(t *testing.T) {
	initMocks := func(mock sqlmock.Sqlmock) {
		mock.ExpectQuery(`SELECT`).WithArgs(email).WillReturnError(testError)
	}

	testMocks := func(t *testing.T, pg user.Postgres) {
		isFalse(t, pg.IsUser(email))
	}

	testPostgresMock(t, initMocks, testMocks)
}

func isFalse(t *testing.T, value bool) {
	if value {
		t.Error("expected false, got ", value)
	}
}

func TestZeroRowsIsUser(t *testing.T) {
	initMocks := func(mock sqlmock.Sqlmock) {
		rows := sqlmock.NewRows([]string{"email", "password"})
		mock.ExpectQuery(`SELECT`).WithArgs(email).WillReturnRows(rows).WillReturnError(nil)
	}

	testMocks := func(t *testing.T, pg user.Postgres) {
		isFalse(t, pg.IsUser(email))
	}

	testPostgresMock(t, initMocks, testMocks)
}

func TestSuccessfulIsUser(t *testing.T) {
	initMocks := func(mock sqlmock.Sqlmock) {
		rows := sqlmock.NewRows([]string{"email", "password"})
		rows.AddRow("testEmail", "testPassword")
		mock.ExpectQuery(`SELECT`).WithArgs(email).WillReturnRows(rows).WillReturnError(nil)
	}

	testMocks := func(t *testing.T, pg user.Postgres) {
		isTrue(t, pg.IsUser(email))
	}

	testPostgresMock(t, initMocks, testMocks)
}

func isTrue(t *testing.T, value bool) {
	if !value {
		t.Error("expected true, got ", value)
	}
}

func TestErrorQueryAuthenticationUser(t *testing.T) {
	initMocks := func(mock sqlmock.Sqlmock) {
		mock.ExpectQuery(`SELECT`).WithArgs(email).WillReturnError(testError)
	}

	testMocks := func(t *testing.T, pg user.Postgres) {
		id, err := pg.AuthenticationUser(email, password)
		isEmptyStr(t, id)
		errorIsTestError(t, err)
	}

	testPostgresMock(t, initMocks, testMocks)
}

func TestNotFoundPasswordAuthenticationUser(t *testing.T) {
	initMocks := func(mock sqlmock.Sqlmock) {
		rows := sqlmock.NewRows([]string{"password"})
		mock.ExpectQuery(`SELECT`).WithArgs(email).WillReturnRows(rows).WillReturnError(nil)
	}

	testMocks := func(t *testing.T, pg user.Postgres) {
		id, err := pg.AuthenticationUser(email, password)
		isEmptyStr(t, id)
		errorIsNotNil(t, err)
	}

	testPostgresMock(t, initMocks, testMocks)
}

func TestInvalidPasswordAuthenticationUser(t *testing.T) {
	initMocks := func(mock sqlmock.Sqlmock) {
		rows := sqlmock.NewRows([]string{"password"})
		rows.AddRow(hasher.Hashing(invalidPassword))
		mock.ExpectQuery(`SELECT`).WithArgs(email).WillReturnRows(rows).WillReturnError(nil)
	}

	testMocks := func(t *testing.T, pg user.Postgres) {
		id, err := pg.AuthenticationUser(email, password)
		isEmptyStr(t, id)
		errorIsNotNil(t, err)
	}

	testPostgresMock(t, initMocks, testMocks)
}

func TestSuccessfulAuthenticationUser(t *testing.T) {
	initMocks := func(mock sqlmock.Sqlmock) {
		initMockSuccessfulAuthenticationUser(mock)
	}

	testMocks := func(t *testing.T, pg user.Postgres) {
		id, err := pg.AuthenticationUser(email, password)
		equalsStr(t, testID, id)
		errorIsNil(t, err)
	}

	testPostgresMock(t, initMocks, testMocks)
}

func initMockSuccessfulAuthenticationUser(mock sqlmock.Sqlmock) {
	rows := sqlmock.NewRows([]string{"password"})
	rows.AddRow(hasher.Hashing(password))
	mock.ExpectQuery(`SELECT`).WithArgs(email).WillReturnRows(rows).WillReturnError(nil)

	rows = sqlmock.NewRows([]string{"id"})
	rows.AddRow(testID)
	mock.ExpectQuery(`SELECT`).WithArgs(email).WillReturnRows(rows).WillReturnError(nil)
}

func equalsStr(t *testing.T, str1, str2 string) {
	if str1 != str2 {
		t.Errorf("%s != %s\n", str1, str2)
	}
}

func TestErrorAuthenticationChangePassword(t *testing.T) {
	initMocks := func(mock sqlmock.Sqlmock) {
		rows := sqlmock.NewRows([]string{"password"})
		rows.AddRow(hasher.Hashing(invalidPassword))
		mock.ExpectQuery(`SELECT`).WithArgs(email).WillReturnRows(rows).WillReturnError(nil)
	}

	testMocks := func(t *testing.T, pg user.Postgres) {
		err := pg.ChangePassword(email, password, newPassword)
		errorIsNotNil(t, err)
	}

	testPostgresMock(t, initMocks, testMocks)
}

func TestErrorExecChangePassword(t *testing.T) {
	initMocks := func(mock sqlmock.Sqlmock) {
		initMockSuccessfulAuthenticationUser(mock)
		mock.ExpectExec(`UPDATE`).WithArgs(email, hasher.Hashing(newPassword)).WillReturnError(testError)
	}

	testMocks := func(t *testing.T, pg user.Postgres) {
		err := pg.ChangePassword(email, password, newPassword)
		errorIsTestError(t, err)
	}

	testPostgresMock(t, initMocks, testMocks)
}

func TestSuccessfulChangePassword(t *testing.T) {
	initMocks := func(mock sqlmock.Sqlmock) {
		initMockSuccessfulAuthenticationUser(mock)
		mock.ExpectExec(`UPDATE`).WithArgs(email, hasher.Hashing(newPassword)).WillReturnResult(emptyResult).WillReturnError(nil)
	}

	testMocks := func(t *testing.T, pg user.Postgres) {
		err := pg.ChangePassword(email, password, newPassword)
		errorIsNil(t, err)
	}

	testPostgresMock(t, initMocks, testMocks)
}

func TestDontCorrectIDDeleteUser(t *testing.T) {
	initMocks := func(mock sqlmock.Sqlmock) {
		initMockSuccessfulAuthenticationUser(mock)
	}

	testMocks := func(t *testing.T, pg user.Postgres) {
		err := pg.DeleteUser(testID[1:], email, password)
		errorIsNotNil(t, err)
	}

	testPostgresMock(t, initMocks, testMocks)
}

func TestErrorDeleteUser(t *testing.T) {
	initMocks := func(mock sqlmock.Sqlmock) {
		initMockSuccessfulAuthenticationUser(mock)
		mock.ExpectExec(`DELETE`).WithArgs(testID, email).WillReturnError(testError)
	}

	testMocks := func(t *testing.T, pg user.Postgres) {
		err := pg.DeleteUser(testID, email, password)
		errorIsTestError(t, err)
	}

	testPostgresMock(t, initMocks, testMocks)
}

func TestSuccessfulDeleteUser(t *testing.T) {
	initMocks := func(mock sqlmock.Sqlmock) {
		initMockSuccessfulAuthenticationUser(mock)
		mock.ExpectExec(`DELETE`).WithArgs(testID, email).WillReturnResult(emptyResult).WillReturnError(nil)
	}

	testMocks := func(t *testing.T, pg user.Postgres) {
		err := pg.DeleteUser(testID, email, password)
		errorIsNil(t, err)
	}

	testPostgresMock(t, initMocks, testMocks)
}

func TestErrorIsVerifiedEmail(t *testing.T) {
	initMocks := func(mock sqlmock.Sqlmock) {
		mock.ExpectQuery(`SELECT verified_email`).WithArgs(email).WillReturnError(testError)
	}

	testMocks := func(t *testing.T, pg user.Postgres) {
		is, err := pg.IsVerifiedEmail(email)
		isFalse(t, is)
		errorIsTestError(t, err)
	}

	testPostgresMock(t, initMocks, testMocks)
}

func TestSuccessfulFalseIsVerifiedEmail(t *testing.T) {
	initMocks := func(mock sqlmock.Sqlmock) {
		rows := getRowsVerifiedEmail(false)
		mock.ExpectQuery(`SELECT verified_email`).WithArgs(email).WillReturnRows(rows)
	}

	testMocks := func(t *testing.T, pg user.Postgres) {
		is, err := pg.IsVerifiedEmail(email)
		isFalse(t, is)
		errorIsNil(t, err)
	}

	testPostgresMock(t, initMocks, testMocks)
}

func getRowsVerifiedEmail(response bool) *sqlmock.Rows {
	rows := sqlmock.NewRows([]string{"verified_email"})
	rows.AddRow(response)
	return rows
}

func TestSuccessfulTrueIsVerifiedEmail(t *testing.T) {
	initMocks := func(mock sqlmock.Sqlmock) {
		rows := getRowsVerifiedEmail(true)
		mock.ExpectQuery(`SELECT verified_email`).WithArgs(email).WillReturnRows(rows)
	}

	testMocks := func(t *testing.T, pg user.Postgres) {
		is, err := pg.IsVerifiedEmail(email)
		isTrue(t, is)
		errorIsNil(t, err)
	}

	testPostgresMock(t, initMocks, testMocks)
}

func TestErrorVerifyEmail(t *testing.T) {
	initMocks := func(mock sqlmock.Sqlmock) {
		mock.ExpectExec(`UPDATE`).WithArgs(email).WillReturnError(testError)
	}

	testMocks := func(t *testing.T, pg user.Postgres) {
		is, err := pg.VerifyEmail(email)
		isFalse(t, is)
		errorIsTestError(t, err)
	}

	testPostgresMock(t, initMocks, testMocks)
}

func TestSuccessfulVerifyEmail(t *testing.T) {
	initMocks := func(mock sqlmock.Sqlmock) {
		mock.ExpectExec(`UPDATE`).WithArgs(email).WillReturnResult(emptyResult)
	}

	testMocks := func(t *testing.T, pg user.Postgres) {
		is, err := pg.VerifyEmail(email)
		isTrue(t, is)
		errorIsNil(t, err)
	}

	testPostgresMock(t, initMocks, testMocks)
}
