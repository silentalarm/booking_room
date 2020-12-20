package page

import (
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

	tmpl, _ := template.ParseFiles("static/clubRegistration.html")
	if r.Method != http.MethodPost {
		data_map := map[string]interface{}{
			"user": user,
		}
		tmpl.Execute(w, data_map)
		//tmpl.Execute(w, nil)
		return
	}

	clubName := r.FormValue("clubName")
	clubAbout := r.FormValue("clubAbout")
	date := time.Now().Format("02.01.2006")
	dbh.InsertNewClub(db, clubAbout, user.Name, user.ID, clubName, user.Name, date)
	http.Redirect(w, r, "/club?clubname="+clubName, http.StatusFound)
}

func ClubsTable(w http.ResponseWriter, r *http.Request) {
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

	//lubs, _ := dbh.GetClubs(db)
	tmpl, _ := template.ParseFiles("static/clubs.html")
	if r.Method != http.MethodPost {
		//data_map := map[string]interface{}{
		//	"clubs": clubs,
		//}
		tmpl.Execute(w, nil)
		return
	}
}
