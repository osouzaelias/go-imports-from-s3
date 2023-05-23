package aws

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"go-import-from-s3/internal/telemetry"
	"go.opentelemetry.io/otel/trace"
	"log"
)

type Config struct {
	bucket            string
	backup            string
	file              string
	delimiter         string
	table             string
	hashKey           string
	rangeKey          string
	ttlName           string
	alwaysDeleteTable bool
	session           *session.Session
	tracer            trace.Tracer
}

func NewConfig() *Config {
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String("us-west-2"),
	})

	if err != nil {
		log.Fatalln(err)
	}

	return &Config{
		bucket:            "bucketelias",
		backup:            "backup",
		file:              "testdata.csv",
		table:             "tb-import-from-s3",
		delimiter:         ",",
		hashKey:           "ID",
		rangeKey:          "FirstName",
		ttlName:           "DataExpirationDate",
		alwaysDeleteTable: true,
		session:           sess,
		tracer:            telemetry.GetTracer(),
	}
}

func (c Config) AlwaysDeleteTable() bool {
	return c.alwaysDeleteTable
}
