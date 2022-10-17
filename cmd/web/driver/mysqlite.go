package driver

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
)

// ConnectToMySQL takes mysql config, forms the connection string and connects to mysql.
func ConnectToMySQL() (*sql.DB, error) {
	db, err := sql.Open("sqlite3", "./apartments.db")
	if err != nil {
		return nil, err
	}
	return db, nil
}
