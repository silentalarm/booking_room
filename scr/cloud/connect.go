package cloud

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
)

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
