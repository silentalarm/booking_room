package cloud

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
)

func Delete(fileName string) error {
	sess := connect()
	svc := s3.New(sess)

	params := &s3.DeleteObjectInput{
		Bucket: aws.String(AWS_S3_BUCKET),
		Key:    aws.String(fileName),
	}

	_, err := svc.DeleteObject(params)
	if err != nil {
		panic(err)
	}
	return nil
}
