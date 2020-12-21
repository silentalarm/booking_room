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
	dbh.UserJoinlub(db, user.Name, clubName, 3, date, user.ID)
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

	clubs, _ := dbh.GetClubs(db)

	for _, key := range clubs {
		fmt.Printf("nick name: %s ", key.ClubName)
	}

	inClub := dbh.IsUserInClub(db, user.Name, user.ID)

	tmpl, _ := template.ParseFiles("static/clubs.html")
	if r.Method != http.MethodPost {
		data_map := map[string]interface{}{
			"user":   user,
			"clubs":  clubs,
			"inclub": inClub,
		}
		tmpl.Execute(w, data_map)
		return
	}

	clubName := r.FormValue("clubName")
	date := time.Now().Format("02.01.2006")
	dbh.UserJoinlub(db, user.Name, clubName, 0, date, user.ID)
	http.Redirect(w, r, "/", http.StatusFound)
}
