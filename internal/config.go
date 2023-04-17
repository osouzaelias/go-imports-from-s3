package internal

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"log"
	"os"
)

type Config struct {
	bucket    string
	backup    string
	file      string
	delimiter string
	table     string
	hashKey   string
	rangeKey  string
	session   session.Session
}

func NewConfig() *Config {
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(os.Getenv("REGION")),
	})

	if err != nil {
		log.Fatalln(err)
	}

	return &Config{
		bucket:    os.Getenv("BUCKET"),
		backup:    os.Getenv("BACKUP"),
		file:      os.Getenv("FILE"),
		table:     os.Getenv("TABLE"),
		delimiter: os.Getenv("DELIMITER"),
		hashKey:   os.Getenv("HASH_KEY"),
		rangeKey:  os.Getenv("RANGE_KEY"),
		session:   *sess,
	}
}
