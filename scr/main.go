package main

import (
	"database/sql"
	"fmt"
	auth "github.com/silentalarm/booking_room/scr/authorization"
	dbh "github.com/silentalarm/booking_room/scr/databaseHandler"
	ses "github.com/silentalarm/booking_room/scr/sessions"
	"golang.org/x/oauth2"
	"html/template"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

type ViewData struct {
	TableData []TData
}
type TData struct {
	LineID       int
	Time         string
	NickName     string
	ClubName     string
	PeopleNumber int
}

var tableWhiteList = []string{
	"floor_2",
	"floor_3",
}

func tableInit() ViewData {
	//http.ServeFile(w, r, "static/table.html")
	data := ViewData{
		TableData: []TData{
			{LineID: 0, Time: "00:00", NickName: "", ClubName: "", PeopleNumber: 0},
			{LineID: 1, Time: "01:00", NickName: "", ClubName: "", PeopleNumber: 0},
			{LineID: 2, Time: "02:00", NickName: "", ClubName: "", PeopleNumber: 0},
			{LineID: 3, Time: "03:00", NickName: "", ClubName: "", PeopleNumber: 0},
			{LineID: 4, Time: "04:00", NickName: "", ClubName: "", PeopleNumber: 0},
			{LineID: 5, Time: "05:00", NickName: "", ClubName: "", PeopleNumber: 0},
			{LineID: 6, Time: "06:00", NickName: "", ClubName: "", PeopleNumber: 0},
			{LineID: 7, Time: "07:00", NickName: "", ClubName: "", PeopleNumber: 0},
			{LineID: 8, Time: "08:00", NickName: "", ClubName: "", PeopleNumber: 0},
			{LineID: 9, Time: "09:00", NickName: "", ClubName: "", PeopleNumber: 0},
			{LineID: 10, Time: "10:00", NickName: "", ClubName: "", PeopleNumber: 0},
			{LineID: 11, Time: "11:00", NickName: "", ClubName: "", PeopleNumber: 0},
			{LineID: 12, Time: "12:00", NickName: "", ClubName: "", PeopleNumber: 0},
			{LineID: 13, Time: "13:00", NickName: "", ClubName: "", PeopleNumber: 0},
			{LineID: 14, Time: "14:00", NickName: "", ClubName: "", PeopleNumber: 0},
			{LineID: 15, Time: "15:00", NickName: "", ClubName: "", PeopleNumber: 0},
			{LineID: 16, Time: "16:00", NickName: "", ClubName: "", PeopleNumber: 0},
			{LineID: 17, Time: "17:00", NickName: "", ClubName: "", PeopleNumber: 0},
			{LineID: 18, Time: "18:00", NickName: "", ClubName: "", PeopleNumber: 0},
			{LineID: 19, Time: "19:00", NickName: "", ClubName: "", PeopleNumber: 0},
			{LineID: 20, Time: "20:00", NickName: "", ClubName: "", PeopleNumber: 0},
			{LineID: 21, Time: "21:00", NickName: "", ClubName: "", PeopleNumber: 0},
			{LineID: 22, Time: "22:00", NickName: "", ClubName: "", PeopleNumber: 0},
			{LineID: 23, Time: "23:00", NickName: "", ClubName: "", PeopleNumber: 0},
		},
	}
	return data
}

func index(w http.ResponseWriter, r *http.Request) {
	db := dbh.OpenDB("postgres")
	defer db.Close()

	tableName := r.URL.Query().Get("table")
	fmt.Printf(tableName)
	tableIsExist := tableIsCorrect(tableName, tableWhiteList)
	if tableIsExist == false {
		tableName = "floor_2"
	}
	date := r.URL.Query().Get("date")
	fmt.Printf(date)
	if date == "" {
		date = time.Now().Format("02.01.2006")
	}
	timeRes, _ := dbh.GetDateReserves(db, tableName, date)
	data := rebuildTable(timeRes)
	tmpl, _ := template.ParseFiles("static/table.html")

	session, err := ses.Store.Get(r, "auth-session")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	user := ses.GetUser(session)

	data_map := map[string]interface{}{
		"var1": data,
		"var2": user,
		"var3": tableName,
		"var4": date,
	}
	tmpl.Execute(w, data_map)
}

func tableIsCorrect(table string, whiteList []string) bool {
	for _, val := range whiteList {
		if table == val {
			return true
		}
	}
	return false
}

func rebuildTable(rows []dbh.ReserveRow) *ViewData {
	data := tableInit()

	for _, row := range rows {
		tableRow := &data.TableData[row.ReserveTime]

		tableRow.NickName = row.NickName
		tableRow.ClubName = row.ClubName
		tableRow.PeopleNumber = row.PeopleNumber
	}
	return &data
}
func about(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "static/calendar.html")
}

func saveToDB(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/saveToDB" {
		http.Error(w, "404 not found.", http.StatusNotFound)
		return
	}
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
	http.Redirect(w, r, "/", http.StatusFound)
	http.Redirect(w, r, "?table="+tableName+"&date="+date, http.StatusTemporaryRedirect)
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

func deleteReserveFromUser(w http.ResponseWriter, r *http.Request) {
	tableName := r.URL.Query().Get("table")
	tableIsExist := tableIsCorrect(tableName, tableWhiteList)
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
	http.Redirect(w, r, "?table="+tableName+"&date="+date, http.StatusTemporaryRedirect)
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

func init() {
	auth.AuthConfig = &oauth2.Config{
		RedirectURL:  "https://booking21.herokuapp.com/callback",
		ClientID:     "c7a7c50ad67f03a72f23c77545b25ac48d616bc1e5daef046d956ed55acf95fd",
		ClientSecret: "157505de170d0b275ab4e10041d4dba1f4f90e21bd1ab5567fc9694b1f040716",
		Scopes:       []string{"public"},
		Endpoint: oauth2.Endpoint{
			AuthURL:   "https://api.intra.42.fr/oauth/authorize",
			TokenURL:  "https://api.intra.42.fr/oauth/token",
			AuthStyle: oauth2.AuthStyleInHeader,
		},
	}
}

func main() {
	http.Handle("/css/", http.StripPrefix("/css/", http.FileServer(http.Dir("css"))))
	http.Handle("/js/", http.StripPrefix("/js/", http.FileServer(http.Dir("js"))))
	http.Handle("/webfonts/", http.StripPrefix("/webfonts/", http.FileServer(http.Dir("webfonts"))))

	port := os.Getenv("PORT")
	//port := "8185"

	http.HandleFunc("/", index)
	http.HandleFunc("/calendar", about)
	http.HandleFunc("/login", auth.Login)
	http.HandleFunc("/callback", auth.CallbackHandler)
	http.HandleFunc("/logout", ses.Delete)
	http.HandleFunc("/saveToDB", saveToDB)
	http.HandleFunc("/delreserve", deleteReserveFromUser)

	fmt.Println("Server is listening...")
	http.ListenAndServe(":"+port, nil)

}

//func main()  {
//	router := mux.NewRouter()
//	router.PathPrefix("/").Handler(http.FileServer(rice.MustFindBox("website").HTTPBox()))
//	http.ListenAndServe(":8082", router)
//}
