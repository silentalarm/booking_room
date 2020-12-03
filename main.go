package main

import (
	"fmt"
	"html/template"
	"net/http"
)
type ViewData struct{
	TableData	[]TData
}
type TData struct{
	Time string
	NickName string
	ClubName string
	Status int
}

func index(w http.ResponseWriter, r *http.Request){
	//http.ServeFile(w, r, "static/table.html")
	data := ViewData {
		TableData: []TData{
			{Time: "00:00 - 01:00", NickName: "User1", ClubName: "Robolab1", Status: 1},
			{Time: "01:00 - 02:00", NickName: "User2", ClubName: "Robolab2", Status: 1},
			{Time: "02:00 - 03:00", NickName: "User2", ClubName: "Robolab2", Status: 1},
			{Time: "03:00 - 04:00", NickName: "User3", ClubName: "Robolab2", Status: 1},
			{Time: "04:00 - 05:00", NickName: "", ClubName: "", Status: 0},
			{Time: "05:00 - 06:00", NickName: "", ClubName: "", Status: 0},
			{Time: "06:00 - 07:00", NickName: "", ClubName: "", Status: 0},
			{Time: "07:00 - 08:00", NickName: "", ClubName: "", Status: 0},
			{Time: "08:00 - 09:00", NickName: "", ClubName: "", Status: 0},
			{Time: "09:00 - 10:00", NickName: "User4", ClubName: "Robolab4", Status: 1},
			{Time: "10:00 - 11:00", NickName: "User5", ClubName: "Robolab5", Status: 1},
			{Time: "11:00 - 12:00", NickName: "", ClubName: "", Status: 0},
			{Time: "12:00 - 13:00", NickName: "", ClubName: "", Status: 0},
			{Time: "13:00 - 14:00", NickName: "", ClubName: "", Status: 0},
			{Time: "14:00 - 15:00", NickName: "", ClubName: "", Status: 0},
			{Time: "15:00 - 16:00", NickName: "", ClubName: "", Status: 0},
			{Time: "16:00 - 17:00", NickName: "", ClubName: "", Status: 0},
			{Time: "17:00 - 18:00", NickName: "", ClubName: "", Status: 0},
			{Time: "18:00 - 19:00", NickName: "", ClubName: "", Status: 0},
			{Time: "19:00 - 20:00", NickName: "", ClubName: "", Status: 0},
			{Time: "20:00 - 21:00", NickName: "", ClubName: "", Status: 0},
			{Time: "21:00 - 22:00", NickName: "", ClubName: "", Status: 0},
			{Time: "22:00 - 23:00", NickName: "", ClubName: "", Status: 0},
			{Time: "23:00 - 00:00", NickName: "", ClubName: "", Status: 0},
		},
	}
	tmpl, _ := template.ParseFiles("static/table.html")
	tmpl.Execute(w, data)
}

func about(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "static/calendar.html")
}

func main() {
	fs := http.FileServer(http.Dir("css"))
	http.Handle("/css/", http.StripPrefix("/css/", fs))

	http.HandleFunc("/", index)
	http.HandleFunc("/calendar", about)

	fmt.Println("Server is listening...")
	http.ListenAndServe(":8185", nil)
}

//func main()  {
//	router := mux.NewRouter()
//	router.PathPrefix("/").Handler(http.FileServer(rice.MustFindBox("website").HTTPBox()))
//	http.ListenAndServe(":8081", router)
//}