package main

import (
	"go-import-from-s3/internal"
	"go-import-from-s3/internal/aws"
	"go-import-from-s3/internal/webhook"
	"log"
)

func main() {
	//shutdown := telemetry.InitProvider()
	//defer shutdown()

	cfg := internal.NewConfig()
	s3Client := aws.NewS3Client(cfg)

	if s3Client.FileExists() {

		dynamoDbClient := aws.NewDynamoDbClient(cfg)
		err := dynamoDbClient.Import()

		if err == nil {

			err = s3Client.MoveToBackup()

			if err == nil {
				err = s3Client.DeleteFile()
				if err != nil {
					log.Println("Error > DeleteFile >", err)
				}
			} else {
				log.Println("Error > MoveToBackup >", err)
			}

			err = dynamoDbClient.EnableTimeToLive()
			if err != nil {
				log.Println("Error > EnableTimeToLive >", err)
			}
		} else {
			log.Println("Error > Import >", err)
		}
	} else {
		if cfg.AlwaysDeleteTable() {
			dynamoDbClient := aws.NewDynamoDbClient(cfg)

			err := dynamoDbClient.PrepareForImport()
			if err != nil {
				log.Println("Error > PrepareForImport >", err)
			}
		}
	}

	hook := webhook.NewCompletedHook(cfg)
	if err := hook.NotifyImportCompleted(); err != nil {
		log.Println("Error > NotifyImportCompleted >", err)
	}

	log.Println("Processamento concluído com sucesso.")
}
