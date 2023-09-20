package user

import (
	"database/sql"
	"errors"
	_ "github.com/lib/pq"
	"user-authorization/internal/config"
	"user-authorization/pkg/errorHandle"
	"user-authorization/pkg/hasher"
)

var (
	errInvalidPassword = errors.New("неверный пароль")
)

type Postgres struct {
	DB *sql.DB
}

func initPostgres(cfg *config.PostgresConfig) (*Postgres, error) {
	db, err := sql.Open("postgres", cfg.GetSourceName())
	if err != nil {
		return nil, err
	}
	return &Postgres{DB: db}, nil
}

func (pg *Postgres) Close() {
	if err := pg.DB.Close(); err != nil {
		errorHandle.Commit(path, "user.go", "Close", err)
	}
}

func (pg *Postgres) NewUser(login, password string) error {
	hash := hasher.Hashing(password)
	query := `INSERT INTO users (login, password) VALUES($1, $2)`
	_, err := pg.DB.Exec(query, login, hash)
	return err
}

func (pg *Postgres) IsUser(login string) bool {
	query := `SELECT COUNT(*) FROM users WHERE login = $1`
	rows, err := pg.DB.Query(query, login)
	return err == nil && rows.Next()
}

func (pg *Postgres) AuthenticationUser(login, password string) error {
	query := `SELECT password FROM users WHERE login = $1`
	rows, err := pg.DB.Query(query, login)
	if err != nil {
		return err
	}

	var savedHash string
	rows.Next()
	err = rows.Scan(&savedHash)
	if err != nil {
		return err
	}

	hash := hasher.Hashing(password)
	if savedHash != hash {
		return errInvalidPassword
	}

	return nil
}

func (pg *Postgres) ChangePassword(login, password, newPassword string) error {
	if err := pg.AuthenticationUser(login, password); err != nil {
		return err
	}
	query := `UPDATE users SET password = $1 WHERE login = $2`
	_, err := pg.DB.Exec(query, newPassword, login)
	return err
}
