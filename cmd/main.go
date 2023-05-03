package main

import (
	"go-import-from-s3/internal/aws"
	"log"
)

func main() {
	//shutdown := telemetry.InitProvider()
	//defer shutdown()

	cfg := aws.NewConfig()
	s3Client := aws.NewS3Client(cfg)

	if s3Client.FileExists() {

		dynamoDbClient := aws.NewDynamoDbClient(cfg)

		if status := dynamoDbClient.Import(); status == aws.ImportStatusCompleted {
			s3Client.MoveToBackup()

			// Session for extra settings
			dynamoDbClient.EnableTimeToLive()
		}
	}

	log.Println("Processamento concluído com sucesso.")
}
