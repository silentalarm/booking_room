package clubs

import (
	dbh "github.com/silentalarm/booking_room/scr/database"
	ses "github.com/silentalarm/booking_room/scr/sessions"
	"html/template"
	"net/http"
	"time"
)

func RegisterNewClub(w http.ResponseWriter, r *http.Request) {
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

	db := dbh.OpenDB()
	defer db.Close()

	//inClub := dbh.IsUserInClub(db, user.Name, user.ID)
	//if inClub == true {
	//	http.Redirect(w, r, "/", http.StatusFound)
	//	return
	//}

	member := dbh.IsUserClubMember(db, user.Name, user.ID)

	tmpl, _ := template.ParseFiles("static/clubRegistration.html")
	if r.Method != http.MethodPost {
		dataMap := map[string]interface{}{
			"user":   user,
			"member": member,
		}
		_ = tmpl.Execute(w, dataMap)
		return
	}

	clubName := r.FormValue("clubName")

	nameIsExist, _ := dbh.ClubNameIsExist(db, clubName)
	if nameIsExist == true {
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}

	clubAbout := r.FormValue("clubAbout")
	slack := r.FormValue("slack")
	date := time.Now().Format("02.01.2006")
	dbh.InsertNewClub(db, clubAbout, clubName, user.Name, date, slack)
	dbh.UserJoinСlub(db, user.Name, clubName, 3, date, user.ID)
	http.Redirect(w, r, "/club?clubname="+clubName, http.StatusFound)
}
