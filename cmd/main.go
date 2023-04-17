package main

import (
	"fmt"
	"go-import-from-s3/internal"
)

func main() {
	conf := internal.NewConfig()

	serviceS3 := internal.NewServiceS3(*conf)
	if serviceS3.FileExists() {
		var err error
		dynamodbSvc := internal.NewServiceDynamoDb(*conf)
		if err = dynamodbSvc.Import(); err != nil {
			fmt.Println("Erro ao importar arquivo:", err.Error())
			return
		}

		if err = serviceS3.MoveToBackup(); err != nil {
			fmt.Println("Erro ao mover aquivo para backup:", err.Error())
			return
		}

		fmt.Println("Processo concluído com sucesso!")
	}
}
