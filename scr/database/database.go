package database

import (
	"database/sql"
	_ "github.com/lib/pq"
	"os"
)

func OpenDB() *sql.DB {
	db, err := sql.Open("postgres", os.Getenv("HEROKU_POSTGRESQL_MAROON_URL"))
	if err != nil {
		panic(err)
	}
	return db
}

func getDBRows(db *sql.DB, table string) (*sql.Rows, error) {
	rows, err := db.Query("SELECT * FROM " + table)
	if err != nil {
		return nil, err
	}

	return rows, nil
}
