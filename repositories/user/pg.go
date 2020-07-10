package user

import (
	"database/sql"
	"errors"
	"github.com/h3isenbug/url-panel/repositories"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
)

type PostgresUserRepository struct {
	con *sqlx.DB
}

func NewPostgresUserRepository(con *sqlx.DB) (UserRepository, error) {
	_, err := con.Exec(`CREATE TABLE IF NOT EXISTS users(
						username VARCHAR(20) UNIQUE,
						email    varchar(100) PRIMARY KEY,
						password varchar(128))`)
	if err != nil {
		return nil, err
	}
	return &PostgresUserRepository{con: con}, nil
}

func (repo PostgresUserRepository) GetUserByUsername(username string) (*User, error) {
	var user User
	err := repo.con.Get(&user, "SELECT * FROM users WHERE username=$1;", username)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, repositories.ErrNotFound
	}
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (repo PostgresUserRepository) GetUserByEMail(email string) (*User, error) {
	var user User
	err := repo.con.Get(&user, "SELECT * FROM users WHERE email=$1;", email)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, repositories.ErrNotFound
	}
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (repo PostgresUserRepository) AddUser(username, email, password string) error {
	_, err := repo.con.Exec("INSERT INTO users (username, email, password) VALUES ($1, $2, $3)", username, email, password)
	pqError, ok := err.(*pq.Error)
	if !ok {
		return err
	}

	if pqError.Constraint == "users_username_key" {
		return ErrUsernameTaken
	}

	if pqError.Constraint == "users_pkey" {
		return ErrEMailTaken
	}

	return nil
}
