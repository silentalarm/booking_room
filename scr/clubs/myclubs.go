package clubs

import (
	dbh "github.com/silentalarm/booking_room/scr/database"
	ses "github.com/silentalarm/booking_room/scr/sessions"
	"html/template"
	"net/http"
	"sort"
)

func MyClubs1(w http.ResponseWriter, r *http.Request) {
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

	member := dbh.IsUserClubMember(db, user.Name, user.ID)
	if member == false { //тут внимательнее возможно не фолс а тру
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}

	clubs, _ := dbh.GetUserClubs(db, true, user.Name, user.ID)
	onConfirmationClubs, _ := dbh.GetUserClubs(db, false, user.Name, user.ID)

	tmpl, _ := template.ParseFiles("static/myClubs.html")
	if r.Method != http.MethodPost {
		dataMap := map[string]interface{}{
			"user":    user,
			"clubs":   clubs,
			"onConCl": onConfirmationClubs,
			"member":  member,
		}
		_ = tmpl.Execute(w, dataMap)
		return
	}

	clubName := r.FormValue("clubName")
	dbh.UserLeaveСlub(db, user.Name, user.ID, clubName)
	http.Redirect(w, r, "/myclubs", http.StatusFound)
}

func MyClubs(w http.ResponseWriter, r *http.Request) {
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
	clubs, _ := dbh.GetUserClubs(db, true, user.Name, user.ID)
	//onConfirmationClubs, _ := dbh.GetUserClubs(db, false, user.Name, user.ID)
	sort.Sort(BySize(clubs))
	//member := dbh.IsUserClubMember(db, user.Name, user.ID)

	tmpl, _ := template.ParseFiles("static/myClubs_v2.html")
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
	case "leaveClub":
		redirect = leave(db, user.Name, user.ID, clubName)
	}

	http.Redirect(w, r, redirect, http.StatusFound)
}
