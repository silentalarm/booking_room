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

	if sumbit == "Удалить клуб" {
		redirect = delete(db, user.Name, user.ID, clubName)
	} else if sumbit == "Сохранить" {
		clubAbout := r.FormValue("clubAbout")

		redirect = save(db, clubAbout, user.Name, user.ID, clubName)
	}

	http.Redirect(w, r, redirect, http.StatusFound)
}

func delete(db *sql.DB, nickName, idIntra, clubName string) string {
	redirect := "/"

	dbh.DeleteClubByOwner(db, nickName, idIntra, clubName)

	return redirect
}

func save(db *sql.DB, newAbout, nickName, idIntra, clubName string) string {
	redirect := "/club?clubname=" + clubName

	dbh.SetAboutClub(db, newAbout, nickName, idIntra, clubName)

	return redirect

}
