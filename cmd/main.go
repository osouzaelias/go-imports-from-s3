package main

import (
	"go-import-from-s3/internal/aws"
	"go-import-from-s3/internal/telemetry"
	"log"
)

func main() {
	shutdown := telemetry.InitProvider()
	defer shutdown()

	//tracer := otel.Tracer("demo-client-tracer")
	//
	//method, _ := baggage.NewMember("method", "repl")
	//client, _ := baggage.NewMember("client", "cli")
	//bag, _ := baggage.New(method, client)
	//
	//defaultCtx := baggage.ContextWithBaggage(context.Background(), bag)

	//for {
	//	_, span := tracer.Start(defaultCtx, "ExecuteRequest")
	//	log.Println("testando")
	//	time.Sleep(5 * time.Second)
	//	span.End()
	//}

	s3Client := aws.NewS3Client()

	if s3Client.FileExists() {

		dynamoDbClient := aws.NewDynamoDbClient()

		if status := dynamoDbClient.Import(); status == aws.ImportStatusCompleted {
			s3Client.MoveToBackup()
		}
	}

	log.Println("Processamento concluído com sucesso.")
}
