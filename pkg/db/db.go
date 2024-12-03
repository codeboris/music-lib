package db

import "database/sql"

func Connect(databaseURL string) (*sql.DB, error) {
	return sql.Open("postgres", databaseURL)
}
