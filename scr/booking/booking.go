package booking

import (
	"database/sql"
	"fmt"
	dbh "github.com/silentalarm/booking_room/scr/database"
	ses "github.com/silentalarm/booking_room/scr/sessions"
	"net/http"
	"strconv"
	"strings"
	"time"
)

var tableWhiteList = []string{
	"floor_2",
	"floor_3",
}

func SaveReserve(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		fmt.Fprintf(w, "ParseForm() err: %v", err)
		return
	}

	session, err := ses.Store.Get(r, "auth-session")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	user := ses.GetUser(session)

	if user.Authenticated == false {
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}

	db := dbh.OpenDB("postgres")
	defer db.Close()

	lines := r.FormValue("lines")
	clubName := r.FormValue("clubName")
	tableName := r.FormValue("hero")
	date := r.FormValue("showDate")
	peopleNumber := r.FormValue("peopleNumber")
	splitedLines := strings.Split(lines, ",")
	convertLines := convertArray(splitedLines)

	succ, unSucc := tryInsertLines(
		user,
		db,
		tableName,
		clubName,
		peopleNumber,
		[]string{
			date,
		},
		convertLines,
	)
	fmt.Print(succ, unSucc)
	//insertFromLines(w, r, db, linesToAdd) //дату изменил (надо сделать чтобы смена даты была из HTML)
	http.Redirect(w, r, "?table="+tableName+"&date="+date, http.StatusFound)
}

func DeleteReserveFromUser(w http.ResponseWriter, r *http.Request) {
	tableName := r.URL.Query().Get("table")
	tableIsExist := TableIsCorrect(tableName)
	if tableIsExist == false {
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}
	date := r.URL.Query().Get("date")
	if date == "" {
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}

	deltime := r.URL.Query().Get("deltime")
	convetedDeltime, err := strconv.Atoi(deltime)
	if err != nil {
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}
	if convetedDeltime < 0 || convetedDeltime > 23 {
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}

	session, err := ses.Store.Get(r, "auth-session")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	user := ses.GetUser(session)
	if user.Authenticated == false {
		fmt.Printf("user: %s auth: %t", user.Name, user.Authenticated)
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}

	db := dbh.OpenDB("postgres")
	defer db.Close()

	dbh.TryDeleteRowByOwner(db, tableName, date, user.Name, deltime)
	http.Redirect(w, r, "?table="+tableName+"&date="+date, http.StatusFound)
}

func tryInsertLines(user *ses.User, db *sql.DB, table string, clubName string, peopleNumber string, date []string, lines []int) (interface{}, interface{}) {
	successfullyAdded := make(map[string][]string)
	unSuccessfullyAdded := make(map[string][]string)
	intPeopleNumber, _ := strconv.Atoi(peopleNumber)
	for _, date_ := range date {
		successfullyLines := []string{}
		unSuccessfullyLines := []string{}

		dateIsCorrect := checkDate(date_)
		if dateIsCorrect == false {
			continue
		}

		for _, i := range lines {
			empty := dbh.ReserveIsExist(db, table, date_, i)
			strHour := strconv.Itoa(i)
			if empty == false {
				dbh.InsertReserve(db, table, user.Name, clubName, intPeopleNumber, i, date_)
				successfullyLines = append(successfullyLines, strHour)
			} else {
				unSuccessfullyLines = append(unSuccessfullyLines, strHour)
			}
		}
		successfullyAdded[date_] = successfullyLines
		unSuccessfullyAdded[date_] = unSuccessfullyLines
	}
	return successfullyAdded, unSuccessfullyAdded
}

func convertArray(lines []string) []int {
	convertedArray := make([]int, len(lines))
	for i := range convertedArray {
		convertedArray[i], _ = strconv.Atoi(lines[i])
	}
	return convertedArray
}

func checkDate(date string) bool {
	today := time.Now()
	targetDate, err := time.Parse("02.01.2006", date)
	if err != nil {
		panic(err.Error())
		return false
	}

	days := targetDate.Sub(today).Hours() / 24
	if days > 30 || days < -1 {
		return false
	}
	return true
}

func TableIsCorrect(table string) bool {
	for _, val := range tableWhiteList {
		if table == val {
			return true
		}
	}
	return false
}
