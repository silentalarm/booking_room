package database

import (
	"database/sql"
	_ "github.com/lib/pq"
	"os"
)

type ReserveRow struct {
	ID           int
	NickName     string
	ClubName     string
	PeopleNumber int
	ReserveTime  int
	ReserveDate  string
}

type Club struct {
	ID           int
	About        string
	NickOwner    string
	IDOwner      string
	ClubName     string
	NickCreator  string
	CreationDate string
	Approved     bool
	Size         int
	Slack        string
	Member       bool
	Owner        bool
}

type ClubMember struct {
	ID       int
	NickName string
	ClubName string
	Access   int
	JoinDate string
	IDIntra  string
}

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
