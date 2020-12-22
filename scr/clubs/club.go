package clubs

import (
	dbh "github.com/silentalarm/booking_room/scr/database"
	ses "github.com/silentalarm/booking_room/scr/sessions"
	"html/template"
	"net/http"
)

func Club(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		panic(err)
		return
	}

	clubName := r.URL.Query().Get("clubname")
	if clubName == "" {
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}

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

	owner := dbh.IsUserClubOwner(db, user.Name, user.ID, clubName)

	//if owner == false {
	//	http.Redirect(w, r, "/", http.StatusFound)
	//	return
	//}

	approved := dbh.ClubIsApproved(db, clubName)
	if approved == false {
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}

	club, _ := dbh.GetClub(db, clubName, true)
	members, _ := dbh.GetClubMembers(db, clubName)
	tmpl, _ := template.ParseFiles("static/club.html")
	if r.Method != http.MethodPost {
		dataMap := map[string]interface{}{
			"user":    user,
			"club":    club,
			"members": members,
			"owner":   owner,
		}
		_ = tmpl.Execute(w, dataMap)
		return
	}

}
