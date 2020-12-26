package clubs

import (
	"database/sql"
	dbh "github.com/silentalarm/booking_room/scr/database"
	ses "github.com/silentalarm/booking_room/scr/sessions"
	"html/template"
	"net/http"
	"sort"
	"time"
)

type BySize []dbh.Club

func (a BySize) Len() int           { return len(a) }
func (a BySize) Less(i, j int) bool { return a[i].Size > a[j].Size }
func (a BySize) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }

func Clubs(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		panic(err)
		return
	}
	redirect := "/"

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

	db := dbh.OpenDB()
	defer db.Close()

	member := dbh.IsUserClubMember(db, user.Name, user.ID)
	clubs, _ := dbh.GetClubs(db, true, user.Name, user.ID)
	sort.Sort(BySize(clubs))
	//member := dbh.IsUserClubMember(db, user.Name, user.ID)

	tmpl, _ := template.ParseFiles("static/clubs_v2.html")
	if r.Method != http.MethodPost {
		dataMap := map[string]interface{}{
			"user":   user,
			"clubs":  clubs,
			"member": member,
		}
		_ = tmpl.Execute(w, dataMap)
		return
	}

	clubName := r.FormValue("clubName")
	sumbit := r.FormValue("sumbit")

	switch sumbit {
	case "goToClub":
		redirect = clubPage(clubName)
	case "joinClub":
		redirect = join(db, user.Name, user.ID, clubName)
	case "leaveClub":
		redirect = leave(db, user.Name, user.ID, clubName)
	}

	http.Redirect(w, r, redirect, http.StatusFound)
}

func clubPage(clubName string) string {
	redirect := "/club?clubname=" + clubName

	return redirect
}

func join(db *sql.DB, nickName, idInta, clubName string) string {
	redirect := "/club?clubname=" + clubName

	date := time.Now().Format("02.01.2006")
	dbh.UserJoinСlub(db, nickName, clubName, 0, date, idInta)

	return redirect
}

func leave(db *sql.DB, nickName, idInta, clubName string) string {
	redirect := "/clubs"

	dbh.UserLeaveСlub(db, nickName, idInta, clubName)

	return redirect
}
