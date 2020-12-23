package clubs

import (
	dbh "github.com/silentalarm/booking_room/scr/database"
	ses "github.com/silentalarm/booking_room/scr/sessions"
	"html/template"
	"net/http"
	"time"
)

func MyClubs(w http.ResponseWriter, r *http.Request) {
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

	member := dbh.IsUserClubMember(db, user.Name, user.ID)
	if member == false { //тут внимательнее возможно не фолс а тру
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}

	clubs, _ := dbh.GetUserClubs(db, true, user.Name, user.ID)

	tmpl, _ := template.ParseFiles("static/myClubs.html")
	if r.Method != http.MethodPost {
		dataMap := map[string]interface{}{
			"user":  user,
			"clubs": clubs,
		}
		_ = tmpl.Execute(w, dataMap)
		return
	}

	clubName := r.FormValue("clubName")
	date := time.Now().Format("02.01.2006")
	dbh.UserJoinСlub(db, user.Name, clubName, 0, date, user.ID)
	http.Redirect(w, r, "/", http.StatusFound)
}
