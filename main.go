package main

import (
	"database/sql"
	"fmt"
	"html/template"
	"net/http"
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
	today := time.Now().Format("02.01.2006")
	//http.ServeFile(w, r, "static/table.html")
	db := openDB("sqlite3", "reserves.db")
	defer db.Close()

	r.ParseForm()
	tableName := r.FormValue("hero")
	timeRes, _ := getDateReserves(db, tableName, today)
	data := rebuildTable(timeRes)
	tmpl, _ := template.ParseFiles("static/table.html")

	session, err := store.Get(r, "auth-session")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	user := getUser(session)

	data_map := map[string]interface{}{
		"var1": data,
		"var2": user,
	}

	tmpl.Execute(w, data_map)

}

func rebuildTable(rows []ReserveRow) *ViewData {
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

	session, err := store.Get(r, "auth-session")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	user := getUser(session)

	if user.Authenticated == false {
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}

	db := openDB("sqlite3", "reserves.db")
	defer db.Close()

	r.ParseForm()
	lines := r.FormValue("lines")
	tableName := r.FormValue("hero")
	clubName := r.FormValue("clubName")
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
			"12.12.2020",
			"13.12.2020",
			"14.12.2020",
		},
		convertLines,
	)
	fmt.Print(succ, unSucc)
	//insertFromLines(w, r, db, linesToAdd) //дату изменил (надо сделать чтобы смена даты была из HTML)
	http.Redirect(w, r, "/", http.StatusFound)
}

func tryInsertLines(user *User, db *sql.DB, table string, clubName string, peopleNumber string, date []string, lines []int) (interface{}, interface{}) {
	successfullyAdded := make(map[string][]string)
	unSuccessfullyAdded := make(map[string][]string)
	intPeopleNumber, _ := strconv.Atoi(peopleNumber)
	fmt.Print("sa")
	for _, date_ := range date {
		successfullyLines := []string{}
		unSuccessfullyLines := []string{}
		for _, i := range lines {
			empty := reserveIsExist(db, table, date_, i)
			strHour := strconv.Itoa(i)
			if empty == false {
				insertReserve(db, table, user.Name, clubName, intPeopleNumber, i, date_)
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

func deleteFromMe(w http.ResponseWriter, r *http.Request) {
	session, err := store.Get(r, "auth-session")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	user := getUser(session)

	if user.Authenticated == false {
		fmt.Printf("user: %s auth: %t", user.Name, user.Authenticated)
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}

	db := openDB("sqlite3", "reserves.db")
	defer db.Close()

	tryDeleteRowByOwner(db, "floor_2", "14.12.2020", user.Name, "2")
	http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
}

func convertArray(lines []string) []int {
	convertedArray := make([]int, len(lines))
	for i := range convertedArray {
		convertedArray[i], _ = strconv.Atoi(lines[i])
	}
	return convertedArray
}

func main() {
	http.Handle("/css/", http.StripPrefix("/css/", http.FileServer(http.Dir("css"))))
	http.Handle("/js/", http.StripPrefix("/js/", http.FileServer(http.Dir("js"))))
	http.Handle("/webfonts/", http.StripPrefix("/webfonts/", http.FileServer(http.Dir("webfonts"))))

	//db := openDB("sqlite3", "reserves.db")
	//today := time.Now().Format("02.02.2007")
	//dbTables := []string{
	//	"floor_3",
	//	"floor_2",
	//}
	//deleteOldReserves(db, dbTables, today)

	http.HandleFunc("/", index)
	http.HandleFunc("/calendar", about)

	http.HandleFunc("/login", authLogin)
	http.HandleFunc("/callback", authCallbackHandler)
	http.HandleFunc("/profile", profileUser)
	http.HandleFunc("/logout", userLogout)
	http.HandleFunc("/saveToDB", saveToDB)
	http.HandleFunc("/delete", deleteFromMe)

	fmt.Println("Server is listening...")
	http.ListenAndServe(":8185", nil)

}

//func main()  {
//	router := mux.NewRouter()
//	router.PathPrefix("/").Handler(http.FileServer(rice.MustFindBox("website").HTTPBox()))
//	http.ListenAndServe(":8082", router)
//}
