package database

import (
	"os"
)

type config struct {
	region    string
	bucket    string
	table     string
	delimiter string
	hashKey   string
	rangeKey  string
}

func newConfigDynamoDB() *config {
	return &config{
		region:    os.Getenv("REGION"),
		bucket:    os.Getenv("BUCKET"),
		table:     os.Getenv("TABLE"),
		delimiter: os.Getenv("DELIMITER"),
		hashKey:   os.Getenv("HASH_KEY"),
		rangeKey:  os.Getenv("RANGE_KEY"),
	}
}
