package booking

import (
	"database/sql"
	"encoding/json"
	"fmt"
	dbh "github.com/silentalarm/booking_room/scr/database"
	ses "github.com/silentalarm/booking_room/scr/sessions"
	"html/template"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type AClub struct {
	Name string
}

type Groups struct {
	List []string
}

var tableWhiteList = []string{
	"floor_2",
	"floor_3",
}

type ViewData struct {
	TableData []TData
}

type TData struct {
	LineID       int
	Time         string
	NickName     string
	ClubName     string
	PeopleNumber int
	Moder        bool
	Color        string
	GroupName    string
	GroupColor   string
}

func tableInit() ViewData {
	//http.ServeFile(w, r, "static/clubs.html")
	data := ViewData{
		TableData: []TData{
			{LineID: 0, Time: "00:00", NickName: "", ClubName: "", PeopleNumber: 0, Moder: false},
			{LineID: 1, Time: "01:00", NickName: "", ClubName: "", PeopleNumber: 0, Moder: false},
			{LineID: 2, Time: "02:00", NickName: "", ClubName: "", PeopleNumber: 0, Moder: false},
			{LineID: 3, Time: "03:00", NickName: "", ClubName: "", PeopleNumber: 0, Moder: false},
			{LineID: 4, Time: "04:00", NickName: "", ClubName: "", PeopleNumber: 0, Moder: false},
			{LineID: 5, Time: "05:00", NickName: "", ClubName: "", PeopleNumber: 0, Moder: false},
			{LineID: 6, Time: "06:00", NickName: "", ClubName: "", PeopleNumber: 0, Moder: false},
			{LineID: 7, Time: "07:00", NickName: "", ClubName: "", PeopleNumber: 0, Moder: false},
			{LineID: 8, Time: "08:00", NickName: "", ClubName: "", PeopleNumber: 0, Moder: false},
			{LineID: 9, Time: "09:00", NickName: "", ClubName: "", PeopleNumber: 0, Moder: false},
			{LineID: 10, Time: "10:00", NickName: "", ClubName: "", PeopleNumber: 0, Moder: false},
			{LineID: 11, Time: "11:00", NickName: "", ClubName: "", PeopleNumber: 0, Moder: false},
			{LineID: 12, Time: "12:00", NickName: "", ClubName: "", PeopleNumber: 0, Moder: false},
			{LineID: 13, Time: "13:00", NickName: "", ClubName: "", PeopleNumber: 0, Moder: false},
			{LineID: 14, Time: "14:00", NickName: "", ClubName: "", PeopleNumber: 0, Moder: false},
			{LineID: 15, Time: "15:00", NickName: "", ClubName: "", PeopleNumber: 0, Moder: false},
			{LineID: 16, Time: "16:00", NickName: "", ClubName: "", PeopleNumber: 0, Moder: false},
			{LineID: 17, Time: "17:00", NickName: "", ClubName: "", PeopleNumber: 0, Moder: false},
			{LineID: 18, Time: "18:00", NickName: "", ClubName: "", PeopleNumber: 0, Moder: false},
			{LineID: 19, Time: "19:00", NickName: "", ClubName: "", PeopleNumber: 0, Moder: false},
			{LineID: 20, Time: "20:00", NickName: "", ClubName: "", PeopleNumber: 0, Moder: false},
			{LineID: 21, Time: "21:00", NickName: "", ClubName: "", PeopleNumber: 0, Moder: false},
			{LineID: 22, Time: "22:00", NickName: "", ClubName: "", PeopleNumber: 0, Moder: false},
			{LineID: 23, Time: "23:00", NickName: "", ClubName: "", PeopleNumber: 0, Moder: false},
		},
	}
	return data
}

func rebuildTable(rows []dbh.ReserveRow, clubs []string) *ViewData {
	data := tableInit()
	db := dbh.OpenDB()
	defer db.Close()

	for _, row := range rows {
		tableRow := &data.TableData[row.ReserveTime]

		tableRow.NickName = row.NickName
		tableRow.ClubName = row.ClubName
		tableRow.PeopleNumber = row.PeopleNumber
		tableRow.GroupName = row.GroupName
		tableRow.GroupColor = dbh.GetGroupColor(db, tableRow.ClubName, tableRow.GroupName)
		tableRow.Color = dbh.GetClubColor(db, tableRow.ClubName)

		for _, clubName := range clubs {
			if tableRow.ClubName == clubName {
				tableRow.Moder = true
			}
		}
	}
	return &data
}

func SaveReserve(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		panic(err)
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

	db := dbh.OpenDB()
	defer db.Close()

	clubName := r.FormValue("clubName")
	clubIsExist := dbh.ClubIsExist(db, clubName)
	if clubIsExist == false {
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}

	groupName := r.FormValue("groupSelect")
	groupIsExis := dbh.GroupIsExist(db, groupName, clubName)
	if groupIsExis == false {
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}

	lines := r.FormValue("lines")
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
		groupName,
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
	tableIsExist := tableIsCorrect(tableName)
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

	db := dbh.OpenDB()
	defer db.Close()

	clubName := r.URL.Query().Get("clubname")

	clubIsExist := dbh.ClubIsExist(db, clubName)
	if clubIsExist == false {
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}

	dbh.DeleteRowByOwnerOrModer(db, user, tableName, date, clubName, deltime)
	http.Redirect(w, r, "?table="+tableName+"&date="+date, http.StatusFound)
}

func tryInsertLines(user *ses.User, db *sql.DB, table, clubName, peopleNumber, groupName string, date []string, lines []int) (interface{}, interface{}) {
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

		for _, line := range lines {
			empty := dbh.ReserveIsExist(db, table, date_, line)
			strHour := strconv.Itoa(line)
			if empty == false && (line >= 0 && line <= 23) {
				_ = dbh.InsertReserve(db, table, user.Name, clubName, intPeopleNumber, line, date_, groupName)
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

func tableIsCorrect(table string) bool {
	for _, val := range tableWhiteList {
		if table == val {
			return true
		}
	}
	return false
}

func Index_v2(w http.ResponseWriter, r *http.Request) {
	db := dbh.OpenDB()
	defer db.Close()

	tableName := r.URL.Query().Get("table")
	fmt.Printf(tableName)
	tableIsExist := tableIsCorrect(tableName)
	if tableIsExist == false {
		tableName = "floor_2"
	}
	date := r.URL.Query().Get("date")
	if date == "" {
		date = time.Now().Format("02.01.2006")
	}

	tmpl, _ := template.ParseFiles("static/table_v2.html")

	session, err := ses.Store.Get(r, "auth-session")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	user := ses.GetUser(session)

	timeRes, _ := dbh.GetDateReserves(db, tableName, date)
	memberClubsOwner, _ := dbh.GetMemberClubsByAccess(db, user.Name, user.ID, 3)
	memberClubsModer, _ := dbh.GetMemberClubsByAccess(db, user.Name, user.ID, 2)
	memberClubsOwner = append(memberClubsOwner, memberClubsModer...)
	data := rebuildTable(timeRes, memberClubsOwner)
	member := dbh.IsUserClubMember(db, user.Name, user.ID)

	var groupList Groups
	if user.Authenticated == true {
		groupList, _ = getGroupList(memberClubsOwner[0])
	}

	dataMap := map[string]interface{}{
		"data":      data,
		"user":      user,
		"tableName": tableName,
		"date":      date,
		"member":    member,
		"clubs":     memberClubsOwner,
		"groupList": groupList,
	}
	_ = tmpl.Execute(w, dataMap)
}
func ExecuteGroupsByClub(w http.ResponseWriter, r *http.Request) {
	var club AClub
	err := json.NewDecoder(r.Body).Decode(&club)
	if err != nil {
		panic(err)
	}

	groupList, _ := getGroupList(club.Name)
	fmt.Printf("%s\n", groupList)

	a, err := json.Marshal(groupList)
	if err != nil {
		panic(err)
	}
	w.Write(a)
}

func getGroupList(clubName string) (Groups, error) {
	db := dbh.OpenDB()
	defer db.Close()

	clubGroups, err := dbh.GetClubGroups(db, clubName)
	if err != nil {
		panic(err)
	}

	groupList := []string{}
	for _, group := range clubGroups {
		groupList = append(groupList, group.Name)
	}

	groups := Groups{
		List: groupList,
	}

	return groups, nil
}
