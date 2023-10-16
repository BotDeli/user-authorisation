package user

import (
	"database/sql"
	"errors"
	_ "github.com/lib/pq"
	"log"
	"user-authorization/internal/config"
	"user-authorization/pkg/errorHandle"
	"user-authorization/pkg/hasher"
)

type Postgres struct {
	DB *sql.DB
}

var (
	errDontCorrectLoginOrPassword = errors.New("некорректный логин или пароль")
)

func initPostgres(cfg *config.PostgresConfig) (*Postgres, error) {
	db, err := sql.Open("postgres", cfg.GetSourceName())
	if err != nil {
		log.Printf("Error connection, %s\n", cfg.GetSourceName())
		return nil, err
	}

	initTables(db)

	return &Postgres{DB: db}, db.Ping()
}

func initTables(db *sql.DB) {
	if _, err := db.Exec(`CREATE TABLE IF NOT EXISTS users (
		login VARCHAR NOT NULL UNIQUE PRIMARY KEY,
		password VARCHAR NOT NULL
	)`); err != nil {
		log.Printf("Error creating tables: %s\n", err)
	}
}

func (pg *Postgres) Close() {
	if err := pg.DB.Close(); err != nil {
		errorHandle.Commit(path, "user.go", "Close", err)
	}
}

func (pg *Postgres) NewUser(login, password string) error {
	hash := hasher.Hashing(password)
	return insertNewUser(pg, login, hash)
}

func insertNewUser(pg *Postgres, login, hash string) error {
	query := `INSERT INTO users (login, password) VALUES($1, $2)`
	_, err := pg.DB.Exec(query, login, hash)
	return err
}

func (pg *Postgres) IsUser(login string) bool {
	rows, err := getRowsPasswordFromLogin(pg, login)
	return err == nil && rows.Next()
}

func getRowsPasswordFromLogin(pg *Postgres, login string) (*sql.Rows, error) {
	query := `SELECT password FROM users WHERE login = $1`
	return pg.DB.Query(query, login)
}

func (pg *Postgres) AuthenticationUser(login, password string) error {
	rows, err := getRowsPasswordFromLogin(pg, login)
	if err != nil {
		return err
	}

	savedHash, err := scanFirstPasswordFromRows(rows)
	if err != nil {
		return err
	}

	return isEqualsSavedHashAndPassword(savedHash, password)
}

func scanFirstPasswordFromRows(rows *sql.Rows) (string, error) {
	var savedHash string
	rows.Next()
	err := rows.Scan(&savedHash)
	return savedHash, err
}

func isEqualsSavedHashAndPassword(savedHash, password string) error {
	hash := hasher.Hashing(password)
	if savedHash != hash {
		return errDontCorrectLoginOrPassword
	}

	return nil
}

func (pg *Postgres) ChangePassword(login, password, newPassword string) error {
	if err := pg.AuthenticationUser(login, password); err != nil {
		return err
	}

	hash := hasher.Hashing(newPassword)
	return updatePassword(pg, login, hash)
}

func updatePassword(pg *Postgres, login, hash string) error {
	query := `UPDATE users SET password = $2 WHERE login = $1`
	_, err := pg.DB.Exec(query, login, hash)
	return err
}
