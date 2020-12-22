package clubs

import (
	"database/sql"
	dbh "github.com/silentalarm/booking_room/scr/database"
	ses "github.com/silentalarm/booking_room/scr/sessions"
	"html/template"
	"net/http"
)

func Moderation(w http.ResponseWriter, r *http.Request) {
	session, err := ses.Store.Get(r, "auth-session")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	redirect := "/"
	user := ses.GetUser(session)

	if user.Authenticated == false { //  || user.Staff == false
		http.Redirect(w, r, redirect, http.StatusFound)
		return
	}

	db := dbh.OpenDB("postgres")
	defer db.Close()

	clubs, _ := dbh.GetClubs(db, false)

	tmpl, _ := template.ParseFiles("static/clubsToApproved.html")
	if r.Method != http.MethodPost {
		dataMap := map[string]interface{}{
			"user":  user,
			"clubs": clubs,
		}
		_ = tmpl.Execute(w, dataMap)
		return
	}

	sumbit := r.FormValue("sumbit")
	if sumbit == "Отказать" {
		redirect = cancel()
	} else if sumbit == "Подтвердить" {
		redirect = accept(db, r)
	}

	http.Redirect(w, r, redirect, http.StatusFound)
}

func accept(db *sql.DB, r *http.Request) string {
	redirect := "/clubstoapproved"

	clubName := r.FormValue("clubName")
	dbh.AppproveClub(db, clubName)

	return redirect
}

func cancel() string {
	retirectURL := "/clubstoapproved"

	return retirectURL
}
