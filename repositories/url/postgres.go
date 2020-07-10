package url

import (
	"database/sql"
	"errors"
	"github.com/jmoiron/sqlx"
)

type PostgresURLRepository struct {
	con *sqlx.DB
}

func NewPostgresURLRepository(con *sqlx.DB) (URLRepository, error) {
	_, err := con.Exec(`CREATE TABLE IF NOT EXISTS urls (
		short_path VARCHAR(10) PRIMARY KEY,
		long_url   VARCHAR(2000),
		email      VARCHAR(100),
		created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
	)`)
	if err != nil {
		return nil, err
	}

	return &PostgresURLRepository{
		con: con,
	}, nil
}

func (repo PostgresURLRepository) SaveShortPath(email, longURL, shortPath string) error {
	_, err := repo.con.Exec("INSERT INTO urls (short_path, long_url, email) VALUES ($1, $2, $3)", shortPath, longURL, email)
	return err
}

func (repo PostgresURLRepository) DeleteURL(email, shortPath string) error {
	_, err := repo.con.Exec("DELETE FROM urls WHERE short_path=$1 AND email=$2", shortPath, email)
	return err
}

func (repo PostgresURLRepository) UserOwnsURL(email, shortPath string) (bool, error) {
	var url URL
	err := repo.con.Get(&url, "SELECT * FROM urls WHERE email=$1 AND short_path=$2", email, shortPath)
	if errors.Is(err, sql.ErrNoRows) {
		return false, nil
	}
	if err != nil {
		return false, err
	}

	return true, nil
}

func (repo PostgresURLRepository) GetMyURLS(email string) ([]*URL, error) {
	var urls []*URL
	rows, err := repo.con.Queryx("SELECT * FROM urls WHERE email=$1", email)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var url URL
		if err := rows.StructScan(&url); err != nil {
			return nil, err
		}

		urls = append(urls, &url)
	}

	return urls, nil
}
