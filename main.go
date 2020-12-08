package main

import (
	"database/sql"
	"fmt"
	"html/template"
	"net/http"
	"strconv"
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
			{LineID: 0, Time: "00:00 - 01:00", NickName: "", ClubName: "", PeopleNumber: 0},
			{LineID: 1, Time: "01:00 - 02:00", NickName: "", ClubName: "", PeopleNumber: 0},
			{LineID: 2, Time: "02:00 - 03:00", NickName: "", ClubName: "", PeopleNumber: 0},
			{LineID: 3, Time: "03:00 - 04:00", NickName: "", ClubName: "", PeopleNumber: 0},
			{LineID: 4, Time: "04:00 - 05:00", NickName: "", ClubName: "", PeopleNumber: 0},
			{LineID: 5, Time: "05:00 - 06:00", NickName: "", ClubName: "", PeopleNumber: 0},
			{LineID: 6, Time: "06:00 - 07:00", NickName: "", ClubName: "", PeopleNumber: 0},
			{LineID: 7, Time: "07:00 - 08:00", NickName: "", ClubName: "", PeopleNumber: 0},
			{LineID: 8, Time: "08:00 - 09:00", NickName: "", ClubName: "", PeopleNumber: 0},
			{LineID: 9, Time: "09:00 - 10:00", NickName: "", ClubName: "", PeopleNumber: 0},
			{LineID: 10, Time: "10:00 - 11:00", NickName: "", ClubName: "", PeopleNumber: 0},
			{LineID: 11, Time: "11:00 - 12:00", NickName: "", ClubName: "", PeopleNumber: 0},
			{LineID: 12, Time: "12:00 - 13:00", NickName: "", ClubName: "", PeopleNumber: 0},
			{LineID: 13, Time: "13:00 - 14:00", NickName: "", ClubName: "", PeopleNumber: 0},
			{LineID: 14, Time: "14:00 - 15:00", NickName: "", ClubName: "", PeopleNumber: 0},
			{LineID: 15, Time: "15:00 - 16:00", NickName: "", ClubName: "", PeopleNumber: 0},
			{LineID: 16, Time: "16:00 - 17:00", NickName: "", ClubName: "", PeopleNumber: 0},
			{LineID: 17, Time: "17:00 - 18:00", NickName: "", ClubName: "", PeopleNumber: 0},
			{LineID: 18, Time: "18:00 - 19:00", NickName: "", ClubName: "", PeopleNumber: 0},
			{LineID: 19, Time: "19:00 - 20:00", NickName: "", ClubName: "", PeopleNumber: 0},
			{LineID: 20, Time: "20:00 - 21:00", NickName: "", ClubName: "", PeopleNumber: 0},
			{LineID: 21, Time: "21:00 - 22:00", NickName: "", ClubName: "", PeopleNumber: 0},
			{LineID: 22, Time: "22:00 - 23:00", NickName: "", ClubName: "", PeopleNumber: 0},
			{LineID: 23, Time: "23:00 - 00:00", NickName: "", ClubName: "", PeopleNumber: 0},
		},
	}
	return data
}

func index(w http.ResponseWriter, r *http.Request) {
	//http.ServeFile(w, r, "static/table.html")
	db := openDB("sqlite3", "reserves.db")
	timeRes, _ := getDateReserves(db, "15.12.20")
	data := rebuildTable(timeRes)
	tmpl, _ := template.ParseFiles("static/table.html")
	tmpl.Execute(w, data)
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

	db := openDB("sqlite3", "reserves.db")
	insertFromLines(db, r, "15.12.20", []int{5, 3})
	http.Redirect(w, r, "/", http.StatusFound)
}

func insertFromLines(db *sql.DB, r *http.Request, date string, lines []int) {
	for _, i := range lines {
		fmt.Print(i)
		strNickName := fmt.Sprintf("nickName%d", i)
		strClubName := fmt.Sprintf("clubName%d", i)
		strPeopleNumber := fmt.Sprintf("peopleNumber%d", i)

		nickName := r.FormValue(strNickName)
		clubName := r.FormValue(strClubName)
		peopleNumber := r.FormValue(strPeopleNumber)

		intPeopleNumber, _ := strconv.Atoi(peopleNumber)

		empty := reserveIsExist(db, date, i)
		if empty == false {
			insertReserve(db, nickName, clubName, intPeopleNumber, i, date)
		}
	}
}

func main() {
	fs := http.FileServer(http.Dir("css"))
	http.Handle("/css/", http.StripPrefix("/css/", fs))
	fs2 := http.FileServer(http.Dir("js"))
	http.Handle("/js/", http.StripPrefix("/js/", fs2))

	db := openDB("sqlite3", "reserves.db")
	//timeRes, _ := getDateReserves(db, "15.12.20")
	//for _, p := range timeRes {
	//	fmt.Println(p.ID, p.NickName, p.ClubName, p.PeopleNumber, p.ReserveTime, p.ReserveDate)
	//}
	date_ := "15.12.20"
	time_ := 10
	empty := reserveIsExist(db, date_, time_)
	//	fmt.Printf("%b", empty)
	if empty == false {
		insertReserve(db, "neya", "top", 1, time_, date_)
	}

	http.HandleFunc("/", index)
	http.HandleFunc("/calendar", about)

	http.HandleFunc("/login", authLogin)
	http.HandleFunc("/callback", authCallbackHandler)
	http.HandleFunc("/profile", profileUser)
	http.HandleFunc("/logout", userLogout)
	http.HandleFunc("/saveToDB", saveToDB)

	fmt.Println("Server is listening...")
	http.ListenAndServe(":8185", nil)

}

//func main()  {
//	router := mux.NewRouter()
//	router.PathPrefix("/").Handler(http.FileServer(rice.MustFindBox("website").HTTPBox()))
//	http.ListenAndServe(":8081", router)
//}
