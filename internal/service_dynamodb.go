package internal

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"time"
)

//const tableName = "tb-import-from-s3-v6"
//const hashKey = "id"
//const rangeKey = "name"

type ServiceDynamoDb struct {
	svc dynamodb.DynamoDB
	cfg Config
}

func NewServiceDynamoDb(c Config) *ServiceDynamoDb {
	return &ServiceDynamoDb{
		svc: *dynamodb.New(&c.session),
		cfg: c,
	}
}

func (s ServiceDynamoDb) Import() error {
	input := s.getImportTableInput()
	output, err := s.svc.ImportTable(input)

	if err != nil {
		return err
	}

	err = s.waitForImportCompletion(output.ImportTableDescription.ImportArn)
	if err != nil {
		return err
	}

	return nil
}

func (s ServiceDynamoDb) waitForImportCompletion(importArn *string) error {
	for {
		describeImportOutput, err := s.svc.DescribeImport(&dynamodb.DescribeImportInput{
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
			return err
		default:
			fmt.Printf("[%s] - Aguardando a importação do arquivo...\n", time.Now().Format("2006-01-02T15:04:05"))
			time.Sleep(time.Second * 5)
		}
	}
}

func (s ServiceDynamoDb) getImportTableInput() *dynamodb.ImportTableInput {
	importTableInput := &dynamodb.ImportTableInput{
		InputFormat: aws.String("CSV"),
		InputFormatOptions: &dynamodb.InputFormatOptions{
			Csv: &dynamodb.CsvOptions{
				Delimiter: aws.String(s.cfg.delimiter),
			},
		},
		S3BucketSource: &dynamodb.S3BucketSource{
			S3Bucket:    aws.String(s.cfg.bucket),
			S3KeyPrefix: aws.String(s.cfg.file),
		},
		TableCreationParameters: &dynamodb.TableCreationParameters{
			AttributeDefinitions: []*dynamodb.AttributeDefinition{
				{
					AttributeName: aws.String(s.cfg.hashKey),
					AttributeType: aws.String("S"),
				},
				{
					AttributeName: aws.String(s.cfg.rangeKey),
					AttributeType: aws.String("S"),
				},
			},
			KeySchema: []*dynamodb.KeySchemaElement{
				{
					AttributeName: aws.String(s.cfg.hashKey),
					KeyType:       aws.String("HASH"),
				},
				{
					AttributeName: aws.String(s.cfg.rangeKey),
					KeyType:       aws.String("RANGE"),
				},
			},
			BillingMode: aws.String("PAY_PER_REQUEST"),
			TableName:   aws.String(s.cfg.table),
		},
	}
	return importTableInput
}
