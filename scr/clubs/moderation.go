package clubs

import (
	"fmt"
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
	user := ses.GetUser(session)

	if user.Authenticated == false { //  || user.Staff == false
		http.Redirect(w, r, "/", http.StatusFound)
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

	sumbitType := r.FormValue("sumbit")

	fmt.Printf("%s", sumbitType)
	clubName := r.FormValue("clubName")
	dbh.AppproveClub(db, clubName)
	http.Redirect(w, r, "/clubstoapproved", http.StatusFound)
}
