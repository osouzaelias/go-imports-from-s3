package internal

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
	webhook           string
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
		alwaysDeleteTable: false,
		session:           sess,
		tracer:            telemetry.GetTracer(),
		webhook:           "http://localhost:8080/endpoint",
	}
}

func (c Config) Backup() string {
	return c.backup
}

func (c Config) File() string {
	return c.file
}

func (c Config) Delimiter() string {
	return c.delimiter
}

func (c Config) Table() string {
	return c.table
}

func (c Config) HashKey() string {
	return c.hashKey
}

func (c Config) RangeKey() string {
	return c.rangeKey
}

func (c Config) TtlName() string {
	return c.ttlName
}

func (c Config) Session() *session.Session {
	return c.session
}

func (c Config) Tracer() trace.Tracer {
	return c.tracer
}

func (c Config) Webhook() string {
	return c.webhook
}

func (c Config) Bucket() string {
	return c.bucket
}

func (c Config) AlwaysDeleteTable() bool {
	return c.alwaysDeleteTable
}
