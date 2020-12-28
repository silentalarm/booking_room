package database

import (
	"database/sql"
	"fmt"
	ses "github.com/silentalarm/booking_room/scr/sessions"
)

type ReserveRow struct {
	ID           int
	NickName     string
	ClubName     string
	PeopleNumber int
	ReserveTime  int
	ReserveDate  string
}

func InsertReserve(db *sql.DB, table string, nickname string, clubname string, people_number int, reserv_time int, reserv_date string) error {
	_, err := db.Exec(
		"INSERT INTO "+table+" (nickname, clubname, people_number, reserv_time, reserv_date) values ($1, $2, $3, $4, $5)",
		nickname, clubname, people_number, reserv_time, reserv_date)
	if err != nil {
		panic(err)
		return err
	}
	return nil
}

func GetDateReserves(db *sql.DB, table string, targetDate string) ([]ReserveRow, error) {
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

func ReserveIsExist(db *sql.DB, table, date string, time int) bool {
	timeRes, _ := GetDateReserves(db, table, date)

	for _, value := range timeRes {
		if value.ReserveTime == time {
			return true
		}
	}
	return false
}

func DeleteOldReserves(db *sql.DB, table []string, date string) {
	for _, target := range table {
		_, err := db.Exec("DELETE FROM " + target + " WHERE reserv_date < '" + date + "'")
		if err != nil {
			panic(err)
		}
	}
}

func DeleteClubRow(db *sql.DB, table, date, clubName, time string) {
	_, err := db.Exec("DELETE FROM " + table + " WHERE reserv_date=" + date + " AND reserv_time=" + time + "AND clubname=" + clubName)
	if err != nil {
		panic(err)
	}
}

func DeleteRowByOwnerOrModer(db *sql.DB, user *ses.User, table, date, clubName, time string) {
	isModer := IsUserClubOwnerOrModer(db, user.Name, user.ID, clubName)
	if isModer == false {
		return
	}
	_, err := db.Exec("DELETE FROM " + table + " WHERE reserv_date=" + date + " AND reserv_time=" + time + "AND clubname=" + clubName)
	if err != nil {
		panic(err)
	}
}

func GetOwnerClubs(db *sql.DB, ownerName string, approved bool) ([]Club, error) {
	rows, err := getDBRows(db, "clubs")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	clubs := []Club{}

	for rows.Next() {
		row := Club{}
		err := acceptRow(rows, &row)
		if err != nil {
			continue
		}

		clubsSize, _ := getClubSize(db, row.ClubName)
		row.NickOwner, row.IDOwner, _ = GetMemberByAccess(db, row.ClubName, 3)
		row.Size = clubsSize
		if row.Approved == approved && row.NickOwner == ownerName {
			clubs = append(clubs, row)
		}
	}
	return clubs, nil
}
