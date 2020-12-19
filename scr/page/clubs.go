package page

import (
	"fmt"
	dbh "github.com/silentalarm/booking_room/scr/database"
	ses "github.com/silentalarm/booking_room/scr/sessions"
	"html/template"
	"net/http"
	"time"
)

func RegistrationPage(w http.ResponseWriter, r *http.Request) {
	tmpl, _ := template.ParseFiles("static/clubRegistration.html")
	tmpl.Execute(w, "")
}

func InsertNewClub(w http.ResponseWriter, r *http.Request) {
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
	fmt.Print(r.Form)
	user := ses.GetUser(session)

	if user.Authenticated == false {
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}

	db := dbh.OpenDB("postgres")
	defer db.Close()

	clubName := r.FormValue("clubName")
	clubAbout := r.FormValue("clubAbout")
	date := time.Now().Format("02.01.2006")
	dbh.InsertNewClub(db, clubAbout, user.Name, user.ID, clubName, user.Name, date)
	http.Redirect(w, r, "/club?clubname="+clubName, http.StatusFound)
}