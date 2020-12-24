package clubs

import (
	dbh "github.com/silentalarm/booking_room/scr/database"
	ses "github.com/silentalarm/booking_room/scr/sessions"
	"html/template"
	"net/http"
	"sort"
	"time"
)

type BySize []dbh.Club

func (a BySize) Len() int           { return len(a) }
func (a BySize) Less(i, j int) bool { return a[i].Size < a[j].Size }
func (a BySize) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }

func Clubs(w http.ResponseWriter, r *http.Request) {
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

	clubs, _ := dbh.GetClubs(db, true, user.Name, user.ID)
	sort.Sort(BySize(clubs))
	//member := dbh.IsUserClubMember(db, user.Name, user.ID)

	tmpl, _ := template.ParseFiles("static/clubs.html")
	if r.Method != http.MethodPost {
		dataMap := map[string]interface{}{
			"user":  user,
			"clubs": clubs,
			//"inclub": member,
		}
		_ = tmpl.Execute(w, dataMap)
		return
	}

	//if member == false { //тут внимательнее возможно не фолс а тру
	//	http.Redirect(w, r, "/", http.StatusFound)
	//	return
	//}

	clubName := r.FormValue("clubName")
	date := time.Now().Format("02.01.2006")
	dbh.UserJoinСlub(db, user.Name, clubName, 0, date, user.ID)
	http.Redirect(w, r, "/", http.StatusFound)
}
