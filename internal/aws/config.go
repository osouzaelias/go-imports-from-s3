package aws

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"go-import-from-s3/internal/telemetry"
	"go.opentelemetry.io/otel/trace"
	"log"
)

type config struct {
	bucket    string
	backup    string
	file      string
	delimiter string
	table     string
	hashKey   string
	rangeKey  string
	session   session.Session
	tracer    trace.Tracer
}

func newConfigDynamoDb() *config {
	return &config{
		bucket:    "bucketelias",
		file:      "myFile0.csv",
		table:     "tb-import-from-s3-v7",
		delimiter: ";",
		hashKey:   "id",
		rangeKey:  "name",
		session:   *getAWSConfigSession(),
		tracer:    *telemetry.GetTracer(),
	}
}

func newConfigS3() *config {
	return &config{
		bucket:    "bucketelias",
		backup:    "backup/",
		file:      "myFile0.csv",
		delimiter: ";",
		session:   *getAWSConfigSession(),
		tracer:    *telemetry.GetTracer(),
	}
}

var instance *session.Session

func getAWSConfigSession() *session.Session {
	if instance == nil {
		sess, err := session.NewSession(&aws.Config{
			Region: aws.String("us-west-2"),
		})

		if err != nil {
			log.Fatalln(err)
		}

		instance = sess
	}
	return instance
}
