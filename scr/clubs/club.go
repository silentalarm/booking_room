package clubs

import (
	"database/sql"
	"github.com/silentalarm/booking_room/scr/cloud"
	dbh "github.com/silentalarm/booking_room/scr/database"
	ses "github.com/silentalarm/booking_room/scr/sessions"
	"html/template"
	"net/http"
	"sort"
)

var allowedColors = []string{
	"#808080",
	"#039BE5",
	"#be68be",
	"#c10000",
	"#2bcd01",
	"#dbd400",
	"#ff8316",
	"#1269ff",
	"#a514ff",
	"#0fbe9b",
}

type ByAccess []dbh.ClubMember

func (a ByAccess) Len() int           { return len(a) }
func (a ByAccess) Less(i, j int) bool { return a[i].Access > a[j].Access }
func (a ByAccess) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }

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

	db := dbh.OpenDB()
	defer db.Close()

	isApproved, _ := dbh.GetClubApprove(db, clubName)
	if isApproved == false {
		http.Redirect(w, r, redirect, http.StatusFound)
		return
	}

	owner := dbh.IsUserClubOwner(db, user.Name, user.ID, clubName)
	club, _ := dbh.GetClub(db, clubName, true)
	member := dbh.IsUserClubMember(db, user.Name, user.ID)
	members, _ := dbh.GetClubMembers(db, clubName)
	groups, _ := dbh.GetClubGroups(db, clubName)
	sort.Sort(ByAccess(members))
	tmpl, _ := template.ParseFiles("static/club.html")
	if r.Method != http.MethodPost {
		dataMap := map[string]interface{}{
			"user":    user,
			"club":    club,
			"members": members,
			"owner":   owner,
			"member":  member,
			"groups":  groups,
		}
		_ = tmpl.Execute(w, dataMap)
		return
	}

	if owner == false {
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}

	sumbit := r.FormValue("sumbit")
	nickName := r.FormValue("nickName")
	intraID := r.FormValue("intraID")

	redirect = stateHandler(r, db, user, sumbit, nickName, clubName, intraID)

	http.Redirect(w, r, redirect, http.StatusFound)
}

func delete(db *sql.DB, clubName string) string {
	redirect := "/"

	dbh.DeleteClub(db, clubName)
	dbh.DeleteUsersFromClub(db, clubName)
	dbh.DeleteAllClubReserves(db, clubName)
	return redirect
}

func save(db *sql.DB, color, slack, newAbout, clubName string) string {
	redirect := "/club?clubname=" + clubName

	dbh.SetAboutClub(db, newAbout, clubName)
	dbh.SetSlackClub(db, slack, clubName)
	colorIsAllowed := allowedColor(color)
	if colorIsAllowed == true {
		dbh.SetColorClub(db, color, clubName)
	}

	return redirect
}

func kick(db *sql.DB, nickName, idIntra, clubName string) string {
	redirect := "/club?clubname=" + clubName

	dbh.UserLeaveСlub(db, nickName, idIntra, clubName)

	return redirect
}

func makeModerator(db *sql.DB, nickName, clubName string) string {
	redirect := "/club?clubname=" + clubName

	_ = dbh.SetAccess(db, nickName, clubName, 2)

	return redirect
}

//fixme
func upRank(db *sql.DB, nickName, clubName string) string {
	redirect := "/club?clubname=" + clubName

	access := dbh.GetMemberAccess(db, nickName, clubName)
	if access == 0 {
		_ = dbh.SetAccess(db, nickName, clubName, access+1)
	}
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

func upload(db *sql.DB, r *http.Request, key, nickName, intraID, clubName string) string {
	redirect := "/club?clubname=" + clubName

	oldName, err := dbh.GetNameFile(db, clubName)
	if err != nil {
		return "/"
	}

	success, err := cloud.Upload(r, key, clubName)
	if err != nil || success == false {
		return "/"
	}

	if oldName != "default_logo.png" {
		cloud.Delete(oldName)
	}

	return redirect
}

func allowedColor(colorForCheck string) bool {
	for _, color := range allowedColors {
		if color == colorForCheck {
			return true
		}
	}
	return false
}

func stateHandler(r *http.Request, db *sql.DB, user *ses.User, sumbit, nickName, clubName, intraID string) string {
	var redirect string

	switch sumbit {
	case "Удалить клуб":
		redirect = delete(db, clubName)
	case "Сохранить":
		clubAbout := r.FormValue("clubAbout")
		slack := r.FormValue("slack")
		color := r.FormValue("color")

		redirect = save(db, color, slack, clubAbout, clubName)
	case "setOwner":
		redirect = setOwner(db, nickName, user.Name, intraID, clubName)
	case "kick":
		redirect = kick(db, nickName, intraID, clubName)
	case "makeModer":
		redirect = makeModerator(db, nickName, clubName)
	default:
		redirect = upload(db, r, "file", user.Name, user.ID, clubName)
	}
	return redirect
}
