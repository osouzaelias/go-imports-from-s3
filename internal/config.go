package internal

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"log"
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
		//Region: aws.String(os.Getenv("REGION")),
		Region: aws.String("us-west-2"),
	})

	if err != nil {
		log.Fatalln(err)
	}

	return &Config{
		bucket:    "bucketelias",
		backup:    "backup/",
		file:      "myFile0.csv",
		table:     "tb-import-from-s3-v7",
		delimiter: ";",
		hashKey:   "id",
		rangeKey:  "name",
		session:   *sess,
	}

	//return &Config{
	//	bucket:    os.Getenv("BUCKET"),
	//	backup:    os.Getenv("BACKUP"),
	//	file:      os.Getenv("FILE"),
	//	table:     os.Getenv("TABLE"),
	//	delimiter: os.Getenv("DELIMITER"),
	//	hashKey:   os.Getenv("HASH_KEY"),
	//	rangeKey:  os.Getenv("RANGE_KEY"),
	//	session:   *sess,
	//}
}
