package report

import (
	dbh "github.com/silentalarm/booking_room/scr/database"
	ses "github.com/silentalarm/booking_room/scr/sessions"
	"html/template"
	"net/http"
)

func List(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		panic(err)
		return
	}
	redirect := "/"

	session, err := ses.Store.Get(r, "auth-session")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	user := ses.GetUser(session)

	if user.Authenticated == false {
		http.Redirect(w, r, redirect, http.StatusFound)
		return
	}

	db := dbh.OpenDB()
	defer db.Close()

	isAdmin := ses.IsAdmin(user.Name)
	reports, _ := dbh.GetReportList(db)

	tmpl, _ := template.ParseFiles("static/reportList.html")
	if r.Method != http.MethodPost {
		dataMap := map[string]interface{}{
			"user":    user,
			"reports": reports,
			"isadmin": isAdmin,
		}
		_ = tmpl.Execute(w, dataMap)
		return
	}

	if isAdmin == false {
		return
	}

	reportID := r.FormValue("reportID")
	dbh.ReportReady(db, reportID)
	http.Redirect(w, r, "/report", http.StatusFound)
}
