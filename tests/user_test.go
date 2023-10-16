package tests

import (
	"github.com/DATA-DOG/go-sqlmock"
	"testing"
	"user-authorization/internal/storage/postgres/user"
	"user-authorization/pkg/hasher"
)

const (
	invalidPassword = "invalid password"
)

var (
	emptyResult = sqlmock.NewErrorResult(nil)
)

func TestErrorNewUser(t *testing.T) {
	initMocks := func(mock sqlmock.Sqlmock) {
		mock.ExpectExec(`INSERT INTO`).WithArgs(login, sqlmock.AnyArg()).WillReturnError(testError)
	}

	testMocks := func(t *testing.T, pg user.Postgres) {
		err := pg.NewUser(login, password)
		if err == nil {
			t.Error("expected error, got nil")
		}
	}

	testPostgresMock(t, initMocks, testMocks)
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
		mock.ExpectExec(`INSERT INTO`).WithArgs(login, sqlmock.AnyArg()).WillReturnResult(emptyResult).WillReturnError(nil)
	}

	testMocks := func(t *testing.T, pg user.Postgres) {
		err := pg.NewUser(login, password)
		if err != nil {
			t.Error(err)
		}
	}

	testPostgresMock(t, initMocks, testMocks)
}

func TestErrorQueryIsUser(t *testing.T) {
	initMocks := func(mock sqlmock.Sqlmock) {
		mock.ExpectQuery(`SELECT`).WithArgs(login).WillReturnError(testError)
	}

	testMocks := func(t *testing.T, pg user.Postgres) {
		if pg.IsUser(login) {
			t.Error("expected false, got true")
		}
	}

	testPostgresMock(t, initMocks, testMocks)
}

func TestZeroRowsIsUser(t *testing.T) {
	initMocks := func(mock sqlmock.Sqlmock) {
		rows := sqlmock.NewRows([]string{"login", "password"})
		mock.ExpectQuery(`SELECT`).WithArgs(login).WillReturnRows(rows).WillReturnError(nil)
	}

	testMocks := func(t *testing.T, pg user.Postgres) {
		if pg.IsUser(login) {
			t.Error("expected false, got true")
		}
	}

	testPostgresMock(t, initMocks, testMocks)
}

func TestSuccessfulIsUser(t *testing.T) {
	initMocks := func(mock sqlmock.Sqlmock) {
		rows := sqlmock.NewRows([]string{"login", "password"})
		rows.AddRow("testLogin", "testPassword")
		mock.ExpectQuery(`SELECT`).WithArgs(login).WillReturnRows(rows).WillReturnError(nil)
	}

	testMocks := func(t *testing.T, pg user.Postgres) {
		if !pg.IsUser(login) {
			t.Error("expected true, got false")
		}
	}

	testPostgresMock(t, initMocks, testMocks)
}

func TestErrorQueryAuthenticationUser(t *testing.T) {
	initMocks := func(mock sqlmock.Sqlmock) {
		mock.ExpectQuery(`SELECT`).WithArgs(login).WillReturnError(testError)
	}

	testMocks := func(t *testing.T, pg user.Postgres) {
		err := pg.AuthenticationUser(login, password)
		if err == nil {
			t.Error("expected error, got nil")
		}
	}

	testPostgresMock(t, initMocks, testMocks)
}

func TestNotFoundPasswordAuthenticationUser(t *testing.T) {
	initMocks := func(mock sqlmock.Sqlmock) {
		rows := sqlmock.NewRows([]string{"password"})
		mock.ExpectQuery(`SELECT`).WithArgs(login).WillReturnRows(rows).WillReturnError(nil)
	}

	testMocks := func(t *testing.T, pg user.Postgres) {
		err := pg.AuthenticationUser(login, password)
		if err == nil {
			t.Error("expected error, got nil")
		}
	}

	testPostgresMock(t, initMocks, testMocks)
}

func TestInvalidPasswordAuthenticationUser(t *testing.T) {
	initMocks := func(mock sqlmock.Sqlmock) {
		rows := sqlmock.NewRows([]string{"password"})
		rows.AddRow(hasher.Hashing(invalidPassword))
		mock.ExpectQuery(`SELECT`).WithArgs(login).WillReturnRows(rows).WillReturnError(nil)
	}

	testMocks := func(t *testing.T, pg user.Postgres) {
		err := pg.AuthenticationUser(login, password)
		if err == nil {
			t.Error("expected error, got nil")
		}
	}

	testPostgresMock(t, initMocks, testMocks)
}

func TestSuccessfulAuthenticationUser(t *testing.T) {
	initMocks := func(mock sqlmock.Sqlmock) {
		rows := sqlmock.NewRows([]string{"password"})
		rows.AddRow(hasher.Hashing(password))
		mock.ExpectQuery(`SELECT`).WithArgs(login).WillReturnRows(rows).WillReturnError(nil)
	}

	testMocks := func(t *testing.T, pg user.Postgres) {
		err := pg.AuthenticationUser(login, password)
		if err != nil {
			t.Error(err)
		}
	}

	testPostgresMock(t, initMocks, testMocks)
}

func TestErrorAuthenticationChangePassword(t *testing.T) {
	initMocks := func(mock sqlmock.Sqlmock) {
		// authentication
		rows := sqlmock.NewRows([]string{"password"})
		rows.AddRow(hasher.Hashing(invalidPassword))
		mock.ExpectQuery(`SELECT`).WithArgs(login).WillReturnRows(rows).WillReturnError(nil)

	}

	testMocks := func(t *testing.T, pg user.Postgres) {
		err := pg.ChangePassword(login, password, newPassword)
		if err == nil {
			t.Error("expected error, got nil")
		}
	}

	testPostgresMock(t, initMocks, testMocks)
}

func TestErrorExecChangePassword(t *testing.T) {
	initMocks := func(mock sqlmock.Sqlmock) {
		// authentication
		rows := sqlmock.NewRows([]string{"password"})
		rows.AddRow(hasher.Hashing(password))
		mock.ExpectQuery(`SELECT`).WithArgs(login).WillReturnRows(rows).WillReturnError(nil)

		// change
		mock.ExpectExec(`UPDATE`).WithArgs(login, newPassword).WillReturnError(testError)
	}

	testMocks := func(t *testing.T, pg user.Postgres) {
		err := pg.ChangePassword(login, password, newPassword)
		if err == nil {
			t.Error("expected error, got nil")
		}
	}

	testPostgresMock(t, initMocks, testMocks)
}

func TestSuccessfulChangePassword(t *testing.T) {
	initMocks := func(mock sqlmock.Sqlmock) {
		// authentication
		rows := sqlmock.NewRows([]string{"password"})
		rows.AddRow(hasher.Hashing(password))
		mock.ExpectQuery(`SELECT`).WithArgs(login).WillReturnRows(rows).WillReturnError(nil)

		// change
		mock.ExpectExec(`UPDATE`).WithArgs(login, hasher.Hashing(newPassword)).WillReturnResult(emptyResult).WillReturnError(nil)
	}

	testMocks := func(t *testing.T, pg user.Postgres) {
		err := pg.ChangePassword(login, password, newPassword)
		if err != nil {
			t.Error(err)
		}
	}

	testPostgresMock(t, initMocks, testMocks)
}
