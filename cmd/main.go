package main

import (
	"go-import-from-s3/internal"
	"log"
)

func main() {
	conf := internal.NewConfig()
	s3Svc := internal.NewServiceS3(*conf)
	if s3Svc.FileExists() {
		dynamodbSvc := internal.NewServiceDynamoDb(*conf)
		dynamodbSvc.Import()
		s3Svc.MoveToBackup()
		log.Println("Processo concluído com sucesso.")
	}
}
