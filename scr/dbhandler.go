package main

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
)

type ReserveRow struct {
	ID           int
	NickName     string
	ClubName     string
	PeopleNumber int
	ReserveTime  int
	ReserveDate  string
}

func openDB(db_name string) *sql.DB {
	connStr := "user=vwemxptfpegfam password=2befa15ec5222f63f97015389b935cd3a4b9ab8bf9b478643b2fdec8e21ab4a5 dbname=deq58l9o8oe3et sslmode=disable"
	db, err := sql.Open(db_name, connStr)
	if err != nil {
		panic(err)
	}
	return db
}

func insertReserve(db *sql.DB, table string, nickname string, clubname string, people_number int, reserv_time int, reserv_date string) error {
	_, err := db.Exec(
		"INSERT INTO "+table+" (nickname, clubname, people_number, reserv_time, reserv_date) values ($1, $2, $3, $4, $5)",
		nickname, clubname, people_number, reserv_time, reserv_date)
	if err != nil {
		panic(err)
		return err
	}
	return nil
}

func getDBRows(db *sql.DB, table string) (*sql.Rows, error) {
	rows, err := db.Query("SELECT * FROM " + table)
	if err != nil {
		return nil, err
	}

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
	timeRes, _ := getDateReserves(db, table, date)

	for _, value := range timeRes {
		if value.ReserveTime == time {
			return true
		}
	}
	return false
}

func deleteOldReserves(db *sql.DB, table []string, date string) {
	for _, target := range table {
		_, err := db.Exec("DELETE FROM " + target + " WHERE reserv_date < '" + date + "'")
		if err != nil {
			panic(err)
		}
	}
}

func tryDeleteRowByOwner(db *sql.DB, table string, date string, userName string, time string) {
	_, err := db.Exec("DELETE FROM " + table + " WHERE nickname='" + userName + "' AND reserv_date='" + date + "' AND reserv_time=" + time)
	if err != nil {
		panic(err)
	}
}
