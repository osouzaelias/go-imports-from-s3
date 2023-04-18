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
		bucket:    "bucketelias",          // Nome do Bucket S3
		backup:    "backup/",              // Diretório de backup dentro do bucket
		file:      "myFile0.csv",          // Arquivo que deverá estar na raiz do bucket
		table:     "tb-import-from-s3-v7", // Nome da tabela do dynamodb que será criada
		delimiter: ";",                    // O caracter separador do arquivo
		hashKey:   "id",                   // Nome da partition key da tabela DynamoDB
		rangeKey:  "name",                 // Nome da sort key da tabela DynamoDB
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
