package aws

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"go-import-from-s3/internal/telemetry"
	"go.opentelemetry.io/otel/trace"
	"log"
)

type timeToLive struct {
	enabled       bool
	attributeName string
}

type Config struct {
	bucket    string
	backup    string
	file      string
	delimiter string
	table     string
	hashKey   string
	rangeKey  string
	ttl       timeToLive
	session   *session.Session
	tracer    *trace.Tracer
}

func NewConfig() *Config {
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String("us-west-2"),
	})

	if err != nil {
		log.Fatalln(err)
	}

	return &Config{
		bucket:    "bucketelias",
		backup:    "backup/",
		file:      "myFile0.csv",
		table:     "tb-import-from-s3-v8",
		delimiter: ";",
		hashKey:   "id",
		rangeKey:  "firstname",
		ttl: timeToLive{
			enabled:       true,
			attributeName: "ttl",
		},
		session: sess,
		tracer:  telemetry.GetTracer(),
	}
}
