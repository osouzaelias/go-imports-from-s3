package main

import (
	"fmt"
	"github.com/aws/aws-sdk-go/service/s3"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

const bucketName = "bucketelias"
const fileName = "myFile0.csv"
const backupFolder = "backup/"

const tableName = "tb-import-from-s3-v6"
const hashKey = "id"
const rangeKey = "name"

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

	if err = moveToBackup(sess); err != nil {
		fmt.Println("Erro ao mover aquivo para backup:", err.Error())
		return
	}

	fmt.Println("Processo concluído com sucesso!")
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
			S3Bucket:    aws.String(bucketName),
			S3KeyPrefix: aws.String(fileName),
		},
		TableCreationParameters: &dynamodb.TableCreationParameters{
			AttributeDefinitions: []*dynamodb.AttributeDefinition{
				{
					AttributeName: aws.String(hashKey),
					AttributeType: aws.String("S"),
				},
				{
					AttributeName: aws.String(rangeKey),
					AttributeType: aws.String("S"),
				},
			},
			KeySchema: []*dynamodb.KeySchemaElement{
				{
					AttributeName: aws.String(hashKey),
					KeyType:       aws.String("HASH"),
				},
				{
					AttributeName: aws.String(rangeKey),
					KeyType:       aws.String("RANGE"),
				},
			},
			BillingMode: aws.String("PAY_PER_REQUEST"),
			TableName:   aws.String(tableName),
		},
	}
	return importTableInput
}

func moveToBackup(sess *session.Session) error {
	svc := s3.New(sess)

	copyInput := &s3.CopyObjectInput{
		Bucket:     aws.String(bucketName),
		CopySource: aws.String(bucketName + "/" + fileName),
		Key:        aws.String(backupFolder + fileName),
	}

	_, err := svc.CopyObject(copyInput)
	if err != nil {
		return err
	}

	deleteInput := &s3.DeleteObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(fileName),
	}
	_, err = svc.DeleteObject(deleteInput)
	return err
}
