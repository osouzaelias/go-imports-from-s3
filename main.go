package main

import (
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

func main() {
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String("us-west-2"),
	})
	if err != nil {
		fmt.Println("Erro ao criar sessão do AWS SDK:", err.Error())
		return
	}

	svc := dynamodb.New(sess)

	importTableInput := getImportTableInput()
	importTableOutput, err := svc.ImportTable(importTableInput)
	if err != nil {
		fmt.Println("Erro ao importar tabela:", err.Error())
		return
	}

	err = waitForImportCompletion(svc, importTableOutput.ImportTableDescription.ImportArn)
	if err != nil {
		fmt.Println("Erro ao importar tabela:", err.Error())
		return
	}

	fmt.Println("Tabela importada com sucesso!")
}

func waitForImportCompletion(svc *dynamodb.DynamoDB, importArn *string) error {
	for {
		describeImportOutput, err := svc.DescribeImport(&dynamodb.DescribeImportInput{
			ImportArn: importArn,
		})

		if err != nil {
			return err
		}

		importStatus := *describeImportOutput.ImportTableDescription.ImportStatus

		switch importStatus {
		case dynamodb.ImportStatusCompleted:
			return nil
		case dynamodb.ImportStatusCancelling, dynamodb.ImportStatusCancelled, dynamodb.ImportStatusFailed:
			return fmt.Errorf("import failed with status %s", importStatus)
		default:
			fmt.Printf("[%s] - Aguardando a importação do arquivo...\n", time.Now().Format("2006-01-02T15:04:05"))
			time.Sleep(time.Second * 5)
		}

		fmt.Printf("[%s] - Aguardando a importação do arquivo...\n", time.Now().Format("2006-01-02T15:04:05"))
		time.Sleep(time.Second * 5)
	}
}

func getImportTableInput() *dynamodb.ImportTableInput {
	importTableInput := &dynamodb.ImportTableInput{
		InputFormat: aws.String("CSV"),
		InputFormatOptions: &dynamodb.InputFormatOptions{
			Csv: &dynamodb.CsvOptions{
				Delimiter: aws.String(";"),
			},
		},
		S3BucketSource: &dynamodb.S3BucketSource{
			S3Bucket:    aws.String("bucketelias"),
			S3KeyPrefix: aws.String("myFile0.csv"),
		},
		TableCreationParameters: &dynamodb.TableCreationParameters{
			AttributeDefinitions: []*dynamodb.AttributeDefinition{
				{
					AttributeName: aws.String("id"),
					AttributeType: aws.String("S"),
				},
				{
					AttributeName: aws.String("name"),
					AttributeType: aws.String("S"),
				},
			},
			KeySchema: []*dynamodb.KeySchemaElement{
				{
					AttributeName: aws.String("id"),
					KeyType:       aws.String("HASH"),
				},
				{
					AttributeName: aws.String("name"),
					KeyType:       aws.String("RANGE"),
				},
			},
			BillingMode: aws.String("PAY_PER_REQUEST"),
			TableName:   aws.String("tb-import-from-s3-v4"),
		},
	}
	return importTableInput
}
