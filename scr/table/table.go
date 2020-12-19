package table

import (
	"fmt"
	bk "github.com/silentalarm/booking_room/scr/booking"
	dbh "github.com/silentalarm/booking_room/scr/databaseHandler"
	ses "github.com/silentalarm/booking_room/scr/sessions"
	"html/template"
	"net/http"
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

func Index(w http.ResponseWriter, r *http.Request) {
	db := dbh.OpenDB("postgres")
	defer db.Close()

	tableName := r.URL.Query().Get("table")
	fmt.Printf(tableName)
	tableIsExist := bk.TableIsCorrect(tableName)
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
