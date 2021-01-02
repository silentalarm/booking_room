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
	thisClubMember := dbh.IsUserInClub(db, user.Name, clubName)
	members, _ := dbh.GetClubMembers(db, clubName)
	groups, _ := dbh.GetClubGroups(db, clubName)
	userGroup, _ := dbh.GetUserGroup(db, user.Name, clubName)
	groupOwner := dbh.IsUserGroupsOwner(db, user.Name, clubName)
	sort.Sort(dbh.ByAccess(members))
	sort.Sort(dbh.ByID(groups))
	tmpl, _ := template.ParseFiles("static/club.html")
	if r.Method != http.MethodPost {
		dataMap := map[string]interface{}{
			"user":       user,
			"club":       club,
			"members":    members,
			"owner":      owner,
			"member":     member,
			"groups":     groups,
			"userGroup":  userGroup,
			"thisMember": thisClubMember,
			"groupOwner": groupOwner,
		}
		_ = tmpl.Execute(w, dataMap)
		return
	}

	sumbit := r.FormValue("sumbit")

	if owner == false {
		redirect = stateHandlerUser(r, db, user, sumbit, clubName)
	} else {
		redirect = stateHandlerOwner(r, db, user, sumbit, clubName)
	}

	http.Redirect(w, r, redirect, http.StatusFound)
}

func delete(db *sql.DB, clubName string) string {
	redirect := "/"

	dbh.DeleteClub(db, clubName)
	dbh.DeleteUsersFromClub(db, clubName)
	dbh.DeleteAllClubGroups(db, clubName)
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

	err := dbh.SetAccess(db, nickName, clubName, 2)
	if err != nil {
		return "/"
	}

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

	member := dbh.IsUserInClub(db, nickName, clubName)
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

func addGroup(db *sql.DB, groupName, memberName, clubName string) string {
	redirect := "/club?clubname=" + clubName

	isAlreadyowner := dbh.IsUserGroupsOwner(db, memberName, clubName)
	if isAlreadyowner == true {
		return "/"
	}

	inClub := dbh.IsUserInClub(db, memberName, clubName)
	if inClub == false {
		return "/"
	}

	err := dbh.CreateNewGroup(db, groupName, clubName, memberName)
	if err != nil {
		return "/"
	}

	err = dbh.SetUserGroup(db, groupName, memberName, clubName)
	if err != nil {
		panic(err)
		return "/"
	}

	err = dbh.SetAccess(db, memberName, clubName, 2)
	if err != nil {
		return "/"
	}

	return redirect
}

func deleteGroup(db *sql.DB, groupName, clubName string) string {
	redirect := "/club?clubname=" + clubName

	if groupName == "main" {
		return "/"
	}

	err := dbh.DeleteGroup(db, groupName, clubName)
	if err != nil {
		panic(err)
		return "/"
	}

	err = dbh.SetUsersGroup(db, groupName, "main", clubName)
	if err != nil {
		panic(err)
		return "/"
	}

	return redirect
}

func stateHandlerOwner(r *http.Request, db *sql.DB, user *ses.User, sumbit, clubName string) string {
	var redirect string
	nickName := r.FormValue("nickName")
	intraID := r.FormValue("intraID")

	switch sumbit {
	case "joinGroup":
		groupName := r.FormValue("groupName")

		redirect = userChangeGroup(db, groupName, user.Name, clubName)
	case "deleteGroup":
		deleteGroupName := r.FormValue("deleteGroupName")

		redirect = deleteGroup(db, deleteGroupName, clubName)
	case "addGroup":
		newGroupName := r.FormValue("newGroupName")
		memberName := r.FormValue("memberName")

		redirect = addGroup(db, newGroupName, memberName, clubName)
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

func userChangeGroup(db *sql.DB, groupName, nickName, clubName string) string {
	redirect := "/club?clubname=" + clubName

	err := dbh.ChangeUserGroup(db, groupName, nickName, clubName)
	if err != nil {
		panic(err)
		return "/"
	}

	return redirect
}

func stateHandlerUser(r *http.Request, db *sql.DB, user *ses.User, sumbit, clubName string) string {
	var redirect string

	isMember := dbh.IsUserInClub(db, user.Name, clubName)
	if isMember == false {
		return redirect
	}

	switch sumbit {
	case "joinGroup":
		groupName := r.FormValue("groupName")

		redirect = userChangeGroup(db, groupName, user.Name, clubName)
	}
	return redirect

}
