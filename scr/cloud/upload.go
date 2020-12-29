package cloud

import (
	"crypto/rand"
	"encoding/hex"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	dbh "github.com/silentalarm/booking_room/scr/database"
	"net/http"
	"path/filepath"
)

var allowedExtension = []string{
	".jpg",
	".png",
	".jpeg",
}

func Upload(r *http.Request, key, clubName string) error {
	file, handler, err := r.FormFile(key) //file, header, err := r.FormFile(key)
	if err != nil {
		panic(err)
		return err
	}
	defer file.Close()

	fileExtension := filepath.Ext(handler.Filename)

	isAllowed := getExtension(fileExtension)
	if isAllowed == false {
		panic(fileExtension)
		return nil
	}

	db := dbh.OpenDB()
	defer db.Close()

	sess := connect()
	randomName := bytesToHex(16)
	filename := randomName + fileExtension //header.flname...

	dbh.SetImageName(db, clubName, filename)

	uploader := s3manager.NewUploader(sess)
	_, err = uploader.Upload(&s3manager.UploadInput{
		Bucket: aws.String(AWS_S3_BUCKET),
		Key:    aws.String(filename),
		ACL:    aws.String("public-read"),
		Body:   file,
	})
	if err != nil {
		panic(err)
		return err
	}
	return nil
}

func getExtension(extension string) bool {
	for _, ext := range allowedExtension {
		if ext == extension {
			return true
		}
	}
	return false
}

func generateBytes(n int) []byte {
	b := make([]byte, n)
	_, err := rand.Read(b)
	if err != nil {
		panic(err)
	}
	return b
}

func bytesToHex(n int) string {
	return hex.EncodeToString(generateBytes(n))
}
