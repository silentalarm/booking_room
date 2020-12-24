package report

import (
	dbh "github.com/silentalarm/booking_room/scr/database"
	ses "github.com/silentalarm/booking_room/scr/sessions"
	"html/template"
	"net/http"
)

func Registration(w http.ResponseWriter, r *http.Request) {
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

	member := dbh.IsUserClubMember(db, user.Name, user.ID)

	tmpl, _ := template.ParseFiles("static/reportRegistration.html")
	if r.Method != http.MethodPost {
		dataMap := map[string]interface{}{
			"user":   user,
			"member": member,
		}
		_ = tmpl.Execute(w, dataMap)
		return
	}

	reportComment := r.FormValue("reportComment")
	if reportComment != "" {
		dbh.InsertNewReport(db, user.Name, user.ID, reportComment)
	}

	http.Redirect(w, r, redirect, http.StatusFound)
}
