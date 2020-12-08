package main

import (
	"fmt"
	"html/template"
	"net/http"
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
	type TableLines struct {
		ID           int
		NickName     int
		ClubName     string
		PeopleNumber int
	}
	// Call ParseForm() to parse the raw query and update r.PostForm and r.Form.
	if err := r.ParseForm(); err != nil {
		fmt.Fprintf(w, "ParseForm() err: %v", err)
		return
	}

	//nickName := make(23, 'string')
	//fmt.Fprintf(w, "Post from website! r.PostFrom = %v\n", r.PostForm)
	r.ParseForm()
	for _, v := range r.Form {
		//fmt.Printf("value: %s\n", v)
		fmt.Printf("value: %s\n", v[0])
		//fmt.Printf("value: %s\n", v[1])
		break
	}
	//fmt.Println("nickName1", r.Form["nickName1"])
	//fmt.Println(r.PostFormValue("nickName1")) //response is empty

	//nickName := r.FormValue("nickName")
	//clubName := r.FormValue("clubName")
	//peopleNumber := r.FormValue("peopleNumber")
	////fmt.Fprintf(w, "nickName = %s\n", nickName)
	////fmt.Fprintf(w, "clubName = %s\n", clubName)
	////fmt.Fprintf(w, "peopleNumber = %s\n", peopleNumber)
	////a, b := nickName[0];
	//for k, v := range nickName {
	//	fmt.Printf("key: %d, value: %d\n", k, v)
	//}
	//for k, v := range clubName {
	//	fmt.Printf("key: %d, value: %d\n", k, v)
	//}
	//for k, v := range peopleNumber {
	//	fmt.Printf("key: %d, value: %d\n", k, v)
	//}

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
