package main

import (
	"database/sql"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
)

type ReserveRow struct {
	ID           int
	NickName     string
	ClubName     string
	PeopleNumber int
	ReserveTime  int
	ReserveDate  string
}

func openDB(db_name string, db_path string) *sql.DB {
	db, err := sql.Open(db_name, db_path)
	if err != nil {
		panic(err)
	}
	//defer db.Close()
	return db
}

func insertReserve(db *sql.DB, table string, nickname string, clubname string, people_number int, reserv_time int, reserv_date string) {
	result, err := db.Exec(
		"insert into "+table+" (nickname, clubname, people_number, reserv_time, reserv_date) values ($1, $2, $3, $4, $5)",
		nickname, clubname, people_number, reserv_time, reserv_date)
	if err != nil {
		panic(err)
	}
	fmt.Println(result.LastInsertId())
}

func getDBRows(db *sql.DB, table string) (*sql.Rows, error) {
	rows, err := db.Query("select * from " + table)
	if err != nil {
		return nil, err
	}
	//defer rows.Close()
	return rows, nil
}

func getDateReserves(db *sql.DB, table string, targetDate string) ([]ReserveRow, error) {
	rows, err := getDBRows(db, table)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	timeReserves := []ReserveRow{}

	for rows.Next() {
		Row := ReserveRow{}
		err := rows.Scan(&Row.ID, &Row.NickName, &Row.ClubName, &Row.PeopleNumber, &Row.ReserveTime, &Row.ReserveDate)
		if err != nil {
			fmt.Println(err)
			continue
		}
		if Row.ReserveDate == targetDate {
			timeReserves = append(timeReserves, Row)
		}
	}
	return timeReserves, nil
}

func reserveIsExist(db *sql.DB, table, date string, time int) bool {
	timeRes, _ := getDateReserves(db, date, table)

	for _, value := range timeRes {
		if value.ReserveTime == time {
			return true
		}
	}
	return false
}

func deleteOldReserves(db *sql.DB, table []string, date string) {
	for _, target := range table {
		_, err := db.Exec("delete from " + target + " where reserv_date < '" + date + "'")
		if err != nil {
			panic(err)
		}
	}
}
