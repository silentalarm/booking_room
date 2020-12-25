package uploader

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	dbh "github.com/silentalarm/booking_room/scr/database"
	ses "github.com/silentalarm/booking_room/scr/sessions"
	"html/template"
	"net/http"
)

const (
	AWS_S3_ID1     = "AKIAX6P2RNC252RG7E3N"
	AWS_S3_SECRET1 = "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"
	AWS_S3_REGION1 = "eu-central-1"
	AWS_S3_BUCKET1 = "21clubs"
)

func connect() *session.Session {
	sess, err := session.NewSession(&aws.Config{
		Region:      aws.String(AWS_S3_REGION1),
		Credentials: credentials.NewStaticCredentials(AWS_S3_ID1, AWS_S3_SECRET1, ""),
	})
	if err != nil {
		panic(err)
	}
	return sess
}

func Upload(w http.ResponseWriter, r *http.Request) {
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

	tmpl, _ := template.ParseFiles("static/upload.html")
	if r.Method != http.MethodPost {
		_ = tmpl.Execute(w, nil)
		return
	}

	db := dbh.OpenDB()
	defer db.Close()

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
		Bucket: aws.String(AWS_S3_BUCKET1),
		Key:    aws.String(filename),
		ACL:    aws.String("public-read"),
		Body:   file,
	})
	if err != nil {
		panic(err)
		return
	}
	http.Redirect(w, r, "/", http.StatusFound)
}
