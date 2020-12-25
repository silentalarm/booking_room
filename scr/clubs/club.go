package clubs

import (
	"database/sql"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/silentalarm/booking_room/scr/cloud"
	dbh "github.com/silentalarm/booking_room/scr/database"
	ses "github.com/silentalarm/booking_room/scr/sessions"
	"html/template"
	"net/http"
	"sort"
)

type ByAccess []dbh.ClubMember

func (a ByAccess) Len() int           { return len(a) }
func (a ByAccess) Less(i, j int) bool { return a[i].Access > a[j].Access }
func (a ByAccess) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }

const (
	AWS_S3_ID     = "AKIAX6P2RNC252RG7E3N"
	AWS_S3_SECRET = "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"
	AWS_S3_REGION = "eu-central-1"
	AWS_S3_BUCKET = "21clubs"
)

func connect() *session.Session {
	sess, err := session.NewSession(&aws.Config{
		Region:      aws.String(AWS_S3_REGION),
		Credentials: credentials.NewStaticCredentials(AWS_S3_ID, AWS_S3_SECRET, ""),
	})
	if err != nil {
		panic(err)
	}
	return sess
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
	member := dbh.IsUserClubMember(db, user.Name, user.ID)
	members, _ := dbh.GetClubMembers(db, clubName)
	sort.Sort(ByAccess(members))
	tmpl, _ := template.ParseFiles("static/club.html")
	if r.Method != http.MethodPost {
		dataMap := map[string]interface{}{
			"user":    user,
			"club":    club,
			"members": members,
			"owner":   owner,
			"member":  member,
		}
		_ = tmpl.Execute(w, dataMap)
		return
	}

	sumbit := r.FormValue("sumbit")
	nickName := r.FormValue("nickName")
	intraID := r.FormValue("intraID")

	switch sumbit {
	case "Удалить клуб":
		redirect = delete(db, user.Name, user.ID, clubName)
	case "Сохранить":
		clubAbout := r.FormValue("clubAbout")

		redirect = save(db, r, "file", clubAbout, user.Name, user.ID, clubName)
	case "upload":
		file, header, err := r.FormFile("file")
		if err != nil {
			panic(err)
			return
		}
		defer file.Close()

		sess := connect()

		filename := header.Filename
		uploader := s3manager.NewUploader(sess)

		_, err = uploader.Upload(&s3manager.UploadInput{
			Bucket: aws.String(AWS_S3_BUCKET),
			Key:    aws.String(filename),
			ACL:    aws.String("public-read"),
			Body:   file,
		})
		if err != nil {
			panic(err)
			return
		}
	case "setOwner":
		redirect = setOwner(db, nickName, user.Name, intraID, clubName)
	case "kick":
		redirect = kick(db, nickName, intraID, clubName)
	case "makeModer":
		redirect = giveModerku(db, nickName, clubName)
	}

	http.Redirect(w, r, redirect, http.StatusFound)
}

func delete(db *sql.DB, nickName, idIntra, clubName string) string {
	redirect := "/"

	dbh.DeleteClubByOwner(db, nickName, idIntra, clubName)
	dbh.DeleteUsersFromClub(db, clubName)

	return redirect
}

func save(db *sql.DB, r *http.Request, key, newAbout, nickName, idIntra, clubName string) string {
	redirect := "/club?clubname=" + clubName

	dbh.SetAboutClub(db, newAbout, nickName, idIntra, clubName)

	return redirect
}

func kick(db *sql.DB, nickName, idIntra, clubName string) string {
	redirect := "/club?clubname=" + clubName

	dbh.UserLeaveСlub(db, nickName, idIntra, clubName)

	return redirect
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

func upload(r *http.Request, key, clubName string) string {
	redirect := "/club?clubname=" + clubName

	cloud.Upload(r, key, clubName)

	return redirect
}
