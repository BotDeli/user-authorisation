package user

import (
	"AccountControl/internal/config"
	"AccountControl/pkg/errorHandle"
	"AccountControl/pkg/generator"
	"AccountControl/pkg/hasher"
	"database/sql"
	"errors"
	_ "github.com/lib/pq"
	"log"
)

type Postgres struct {
	DB *sql.DB
}

var (
	errDontCorrectEmailOrPassword = errors.New("некорректный логин или пароль")
	errDontCorrectData            = errors.New("некорректные данные")
	errGenerationUniqueId         = errors.New("ошибка генерации уникального id")
)

func initPostgres(cfg *config.PostgresConfig) (*Postgres, error) {
	db, err := sql.Open("postgres", cfg.GetSourceName())

	if err != nil {
		log.Printf("Error connection, %s\n", cfg.GetSourceName())
		return nil, err
	}

	initTable(db)

	return &Postgres{DB: db}, db.Ping()
}

func initTable(db *sql.DB) {
	if _, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS users (
    	id VARCHAR NOT NULL PRIMARY KEY UNIQUE,
		email VARCHAR NOT NULL UNIQUE,
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

func (pg *Postgres) NewUser(email, password string) (string, error) {
	id, err := getUniqueID(pg)
	if err != nil {
		log.Println(err)
		return "", err
	}

	hash := hasher.Hashing(password)
	err = insertNewUser(pg, id, email, hash)
	if err != nil {
		return "", err
	}

	return id, nil
}

func getUniqueID(pg *Postgres) (string, error) {
	var id string
	for i := 0; i < 10; i++ {
		id = generator.NewUUIDDigits()
		if isUniqueId(pg, id) {
			return id, nil
		}
	}
	return "", errGenerationUniqueId
}

func isUniqueId(pg *Postgres, id string) bool {
	if id == "" {
		return false
	}

	actual, err := ActualID(pg, id)
	if err != nil {
		return false
	}

	return !actual
}

func ActualID(pg *Postgres, id string) (bool, error) {
	query := `SELECT COUNT(*) FROM users WHERE id = $1`
	rows, err := pg.DB.Query(query, id)
	if err != nil {
		return false, err
	}

	count, err := scanFirstOneValueFromRows[int](rows)
	if err != nil {
		return false, err
	}

	return count > 0, nil
}

func scanFirstOneValueFromRows[T any](rows *sql.Rows) (T, error) {
	var value T
	rows.Next()
	err := rows.Scan(&value)
	return value, err
}

func insertNewUser(pg *Postgres, id, email, hash string) error {
	query := `INSERT INTO users (id, email, password) VALUES($1, $2, $3)`
	_, err := pg.DB.Exec(query, id, email, hash)
	return err
}

func (pg *Postgres) IsUser(email string) bool {
	rows, err := getRowsPasswordFromEmail(pg, email)
	return err == nil && rows.Next()
}

func getRowsPasswordFromEmail(pg *Postgres, email string) (*sql.Rows, error) {
	query := `SELECT password FROM users WHERE email = $1`
	return pg.DB.Query(query, email)
}

func (pg *Postgres) AuthenticationUser(email, password string) (string, error) {
	rows, err := getRowsPasswordFromEmail(pg, email)
	if err != nil {
		return "", err
	}

	savedHash, err := scanFirstOneValueFromRows[string](rows)
	if err != nil {
		return "", err
	}

	hash := hasher.Hashing(password)
	if savedHash != hash {
		return "", errDontCorrectEmailOrPassword
	}

	return getIDFromEmail(pg, email)
}

func getIDFromEmail(pg *Postgres, email string) (string, error) {
	rows, err := getRowsIDFromEmail(pg, email)
	if err != nil {
		return "", err
	}
	return scanFirstOneValueFromRows[string](rows)
}

func getRowsIDFromEmail(pg *Postgres, email string) (*sql.Rows, error) {
	query := `SELECT id FROM users WHERE email = $1`
	return pg.DB.Query(query, email)
}

func (pg *Postgres) ChangePassword(email, password, newPassword string) error {
	if _, err := pg.AuthenticationUser(email, password); err != nil {
		return err
	}

	hash := hasher.Hashing(newPassword)
	return updatePassword(pg, email, hash)
}

func updatePassword(pg *Postgres, email, hash string) error {
	query := `UPDATE users SET password = $2 WHERE email = $1`
	_, err := pg.DB.Exec(query, email, hash)
	return err
}

func (pg *Postgres) DeleteUser(id, email, password string) error {
	verifyID, err := pg.AuthenticationUser(email, password)
	if err != nil {
		return err
	}

	if id != verifyID {
		return errDontCorrectData
	}
	return dropUserData(pg, id, email)
}

func dropUserData(pg *Postgres, id, email string) error {
	query := `DELETE FROM users WHERE id = $1 AND email = $2`
	_, err := pg.DB.Exec(query, id, email)
	return err
}
