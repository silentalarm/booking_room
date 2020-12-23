package clubs

import (
	"database/sql"
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
	redirect := "/"

	clubName := r.URL.Query().Get("clubname")
	if clubName == "" {
		http.Redirect(w, r, redirect, http.StatusFound)
		return
	}

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

	db := dbh.OpenDB("postgres")
	defer db.Close()

	owner := dbh.IsUserClubOwner(db, user.Name, user.ID, clubName)

	//if owner == false {
	//	http.Redirect(w, r, "/", http.StatusFound)
	//	return
	//}

	club, _ := dbh.GetClub(db, clubName, true)
	if club.Approved == false {
		http.Redirect(w, r, redirect, http.StatusFound)
		return
	}

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

	sumbit := r.FormValue("sumbit")

	switch sumbit {
	case "Удалить клуб":
		redirect = delete(db, user.Name, user.ID, clubName)
	case "Сохранить":
		clubAbout := r.FormValue("clubAbout")

		redirect = save(db, clubAbout, user.Name, user.ID, clubName)
	case "setOwner":
		nickName := r.FormValue("nickName")
		intraID := r.FormValue("intraID")

		redirect = setOwner(db, nickName, user.Name, intraID, clubName)
	case "kick":

	case "makeModer":

	}

	http.Redirect(w, r, redirect, http.StatusFound)
}

func delete(db *sql.DB, nickName, idIntra, clubName string) string {
	redirect := "/"

	dbh.DeleteClubByOwner(db, nickName, idIntra, clubName)
	dbh.DeleteUsersFromClub(db, clubName)

	return redirect
}

func save(db *sql.DB, newAbout, nickName, idIntra, clubName string) string {
	redirect := "/club?clubname=" + clubName

	dbh.SetAboutClub(db, newAbout, nickName, idIntra, clubName)

	return redirect
}

func kick() {

}

func giveModerku(db *sql.DB, nickName, clubName string) string {
	redirect := "/club?clubname=" + clubName

	_ = dbh.SetAccess(db, nickName, clubName, 2)

	return redirect
}

func setOwner(db *sql.DB, nickName, nickOwner, intraID, clubName string) string {
	redirect := "/club?clubname=" + clubName

	member := dbh.IsUserInClub(db, nickName, intraID, clubName)
	if member == false {
		redirect = "/lol?ss=" + intraID
		return redirect
	}

	_ = dbh.SetAccess(db, nickName, clubName, 3)
	_ = dbh.SetAccess(db, nickOwner, clubName, 0)
	dbh.SetClubOwner(db, nickName, clubName)

	return redirect
}
